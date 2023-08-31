import 'package:eleve/exercice/congratulations.dart';
import 'package:eleve/exercice/home.dart';
import 'package:eleve/questions/fields.dart';
import 'package:eleve/questions/question.dart';
import 'package:eleve/quotes.dart';
import 'package:eleve/shared/errors.dart';
import 'package:eleve/types/src.dart';
import 'package:eleve/types/src_maths_questions_client.dart';
import 'package:eleve/types/src_sql_editor.dart';
import 'package:eleve/types/src_tasks.dart';
import 'package:flutter/material.dart' hide Flow;

class NotificationExerciceDone extends Notification {}

abstract class ExerciceAPI extends FieldAPI {
  /// [evaluate] must evaluate the exercice answers, according
  /// to the exercice mode, and returns the feedback and the new version
  /// of the questions if needed
  Future<EvaluateWorkOut> evaluate(EvaluateWorkIn params);
}

class _ExerciceQuestionController extends BaseQuestionController {
  void Function() onClick;

  _ExerciceQuestionController(Question question, FieldAPI api, this.onClick)
      : super(question, api);

  void markDone() {
    state.buttonLabel = "Question terminée";
    state.buttonEnabled = false;
    state.footerQuote = pickQuote();
    disableAllFields();
  }

  // delegate to the exercice widget
  @override
  void onPrimaryButtonClick() {
    onClick();
  }
}

enum ExerciceStatus {
  /// the student is working on a question, validation triggers evaluation
  answering,

  /// the wrong answers are displayed in red, validation triggers retry
  displayingFeedback,
}

/// [ExerciceController] exposes the state
/// which will change during an exercice,
/// so that parent widget may react properly.
///
class ExerciceController {
  /// [exeAndProg] stores the server instantiated exercice with
  /// the progression state.
  StudentWork exeAndProg;

  /// [questionIndex] is the current question, or null for the summary
  int? questionIndex;

  ExerciceStatus status = ExerciceStatus.answering;

  List<QuestionStatus> questionStates = [];

  List<_ExerciceQuestionController> _questions;

  /// onValid is the callback triggered when validatin a question
  /// and is filled by the widget using the controller
  void Function()? onValid;

  final FieldAPI api;

  ExerciceController(this.exeAndProg, this.questionIndex, this.api)
      : _questions = [] {
    _questions = _createControllers(api);
    _refreshStates();
  }

  Enonce get currentCorrection {
    if (questionIndex == null) return [];
    return exeAndProg.exercice.questions[questionIndex!].question.correction;
  }

  /// checks if the student as already try at least once the current question
  bool get isCurrentCorrectionEnabled {
    if (questionIndex == null) return false;
    return exeAndProg.progression.questions[questionIndex!].isNotEmpty;
  }

  List<_ExerciceQuestionController> _createControllers(FieldAPI api) {
    return exeAndProg.exercice.questions
        .map((qu) =>
            _ExerciceQuestionController(qu.question, api, _onQuestionValid))
        .toList();
  }

  void _onQuestionValid() {
    switch (status) {
      case ExerciceStatus.answering:
        _questions[questionIndex!].state.buttonEnabled = false;
        _questions[questionIndex!].state.buttonLabel = "Correction...";
        break;
      case ExerciceStatus.displayingFeedback:
        _questions[questionIndex!].state.buttonEnabled = true;
        _questions[questionIndex!].state.buttonLabel = "Valider...";
        break;
    }
    _refreshStates();
    if (onValid != null) onValid!();
  }

  /// [showFeedback] set the given feedback (for many questions)
  /// and set the state to displayingFeedback
  void showFeedback(Map<int, QuestionAnswersOut> feedbacks) {
    status = ExerciceStatus.displayingFeedback;
    for (var question in feedbacks.entries) {
      final index = question.key;
      final feedback = question.value.results;
      _questions[index].setFeedback(feedback);
    }
    for (var question in _questions) {
      question.state.buttonEnabled = true;
      question.state.buttonLabel = "Essayer à nouveau...";
    }
    _refreshStates();
  }

  void reset() {
    _refreshStates();
    status = ExerciceStatus.answering;
  }

  void resetWithNextQuestions(List<InstantiatedQuestion> nextQuestions) {
    exeAndProg = exeAndProg.copyWithQuestions(nextQuestions);
    _questions = _createControllers(api);
    reset();
  }

  void _refreshStates() {
    questionStates = _computeQuestionStates();
    for (var index = 0; index < _questions.length; index++) {
      final qu = _questions[index];
      if (questionStates[index] == QuestionStatus.checked) {
        qu.markDone();
      }
    }
  }

  /// [updateProgression] is called after receiving the server
  /// answer. It removes waiting questions and update the questions status.
  void updateProgression(ProgressionExt progression) {
    exeAndProg = StudentWork(exeAndProg.exercice, progression);
    _refreshStates();
  }

  List<QuestionStatus> _computeQuestionStates() {
    final exP = exeAndProg;
    return List<QuestionStatus>.generate(exP.exercice.questions.length,
        (questionIndex) {
      if (exP.progression.isQuestionCompleted(questionIndex)) {
        return QuestionStatus.checked;
      }

      if (exP.exercice.flow == Flow.sequencial &&
          exP.progression.nextQuestion < questionIndex) {
        return QuestionStatus.locked;
      }

      final qu = _questions[questionIndex];
      if (qu.feedback().values.any((error) => error)) {
        return QuestionStatus.incorrect;
      }
      return QuestionStatus.toDo;
    });
  }

  /// setQuestionAnswers show the answers for the current question
  void setQuestionAnswers(Answers answers) {
    if (questionIndex == null) return;
    _questions[questionIndex!].setAnswers(answers);
    _questions[questionIndex!].state.buttonEnabled = true;
    _questions[questionIndex!].state.buttonLabel = "Valider";
  }

  bool isExerciceOver() => exeAndProg.progression.nextQuestion == -1;

  bool get goToPreviousEnabled => questionIndex != null;

  bool get goToNextEnabled {
    final currentIndex = questionIndex;
    final hasNextQuestion = exeAndProg.exercice.questions.isNotEmpty &&
        currentIndex != exeAndProg.exercice.questions.length - 1;
    // always allow to go to already done question
    if (hasNextQuestion &&
        questionStates[(currentIndex ?? -1) + 1] == QuestionStatus.checked) {
      return true;
    }
    switch (exeAndProg.exercice.flow) {
      case Flow
            .sequencial: // do not show locked questions when exercice is not over
        return currentIndex == null ||
            (isExerciceOver() && hasNextQuestion) ||
            currentIndex < exeAndProg.progression.nextQuestion;
      case Flow.parallel: // no restriction:
        return hasNextQuestion;
    }
  }
}

extension _CopyIW on InstantiatedWork {
  InstantiatedWork copyWithQuestions(List<InstantiatedQuestion> questions) {
    return InstantiatedWork(
        iD,
        title,
        flow,
        // replace the questions
        questions,
        baremes);
  }
}

extension _CopySW on StudentWork {
  StudentWork copyWithQuestions(List<InstantiatedQuestion> questions) {
    return StudentWork(exercice.copyWithQuestions(questions), progression);
  }
}

/// ExerciceW is the widget providing one exercice to
/// the student.
/// It is used in the editor loopback, and as the base for
/// at home training activity
class ExerciceW extends StatefulWidget {
  final ExerciceAPI api;

  final ExerciceController controller;

  final void Function()? onShowCorrectAnswer;

  /// if true, and if the exercice has a correction,
  /// show a button to access the correction after on try
  final bool showCorrectionButtonOnFail;

  /// if true, displays the correction root instead of
  /// the enonce
  /// Only used in loopback
  final bool instantShowCorrection;

  /// if true, show an alter about progression not being updated
  final bool noticeSandbox;

  const ExerciceW(this.api, this.controller,
      {Key? key,
      this.onShowCorrectAnswer,
      this.showCorrectionButtonOnFail = true,
      this.instantShowCorrection = false,
      this.noticeSandbox = false})
      : super(key: key);

  @override
  State<ExerciceW> createState() => _ExerciceWState();
}

class _ExerciceWState extends State<ExerciceW> {
  // the questions to display when trying again
  // this is a temporay slice, affected to the controller on user validation
  List<InstantiatedQuestion> nextQuestions = [];

  @override
  void initState() {
    _init();
    super.initState();
  }

  @override
  void didUpdateWidget(covariant ExerciceW oldWidget) {
    _init();
    super.didUpdateWidget(oldWidget);
  }

  void _init() {
    nextQuestions = [];
    widget.controller.reset();
    widget.controller.onValid = onQuestionButtonClick;

    if (widget.instantShowCorrection) {
      WidgetsBinding.instance.addPostFrameCallback((_) => _showCorrection());
    }
  }

  @override
  Widget build(BuildContext context) {
    final ct = widget.controller;
    final exP = widget.controller.exeAndProg;
    return Scaffold(
        appBar: AppBar(
          title: const Text("Exercice"),
          actions: [
            if (widget.onShowCorrectAnswer != null && ct.questionIndex != null)
              TextButton(
                  onPressed: widget.onShowCorrectAnswer,
                  child: const Text("Afficher la réponse")),
            IconButton(
                onPressed: ct.goToPreviousEnabled ? goToPrevious : null,
                icon: const Icon(IconData(0xf572,
                    fontFamily: 'MaterialIcons', matchTextDirection: true))),
            IconButton(
                onPressed: ct.goToNextEnabled ? goToNext : null,
                icon: const Icon(IconData(0xf57a,
                    fontFamily: 'MaterialIcons', matchTextDirection: true))),
          ],
        ),
        body: ct.questionIndex == null
            ? ExerciceHome(
                exP,
                ct.questionStates,
                (index) => setState(() {
                      ct.questionIndex = index;
                    }),
                widget.noticeSandbox)
            : QuestionW(
                widget.controller._questions[widget.controller.questionIndex!],
                Colors.purpleAccent,
                title: "Question ${widget.controller.questionIndex! + 1}",
              ));
  }

  // handle the errors
  Future<EvaluateWorkOut?> _evaluate(EvaluateWorkIn params) async {
    try {
      final res = await widget.api.evaluate(params);
      return res;
    } catch (error) {
      showError("Impossible d'évaluer l'exercice", error, context);
      return null;
    }
  }

  void onExerciceOver() async {
    setState(() {
      widget.controller.questionIndex = null;
    });

    // exercice is over
    await showDialog<void>(
        context: context,
        builder: (context) => const Dialog(child: Congrats()));
    if (!mounted) return;
    NotificationExerciceDone().dispatch(context);
  }

  void _onValideQuestion() async {
    final ct = widget.controller;
    final index = ct.questionIndex!;

    // for sequencial exercices, if we are not at the current question, just go to it
    // and return
    if (ct.exeAndProg.exercice.flow == Flow.sequencial &&
        index != ct.exeAndProg.progression.nextQuestion) {
      setState(() {
        ct.questionIndex = ct.exeAndProg.progression.nextQuestion;
      });
      return;
    }

    // validate the given answer
    final resp = await _evaluate(EvaluateWorkIn(
        ct.exeAndProg.exercice.iD,
        {
          index: AnswerP(ct.exeAndProg.exercice.questions[index].params,
              ct._questions[index].answers())
        },
        ct.exeAndProg.progression));
    if (resp == null || !mounted) {
      return;
    }

    ct.updateProgression(resp.progression); // update the progression
    nextQuestions = resp.newQuestions; // buffer until retry

    if (ct.exeAndProg.progression.nextQuestion == -1) {
      onExerciceOver();
      return;
    }

    final isCorrect = resp.results[index]!.isCorrect;
    if (!isCorrect) {
      // show errors and ask for retry
      setState(() {
        ct.showFeedback({index: resp.results[index]!});
      });
    } else {
      // update the menu status
      setState(() {});
    }

    _showValidDialogOrSnack(isCorrect);
  }

  // handle the following cases :
  //  - the answer is correct and there is more to do
  //  - the answer is incorrect and there is a correction to display
  //  - the answer is incorrect and there is no correction to display
  void _showValidDialogOrSnack(bool isAnswerCorrect) async {
    final ct = widget.controller;
    final hasNextQuestion =
        isAnswerCorrect && ct.exeAndProg.progression.nextQuestion != -1;
    final showButtonCorrection = widget.showCorrectionButtonOnFail &&
        ct.currentCorrection.isNotEmpty &&
        !isAnswerCorrect;

    if (hasNextQuestion) {
      // show a dialog with next button
      await showDialog<void>(
          context: context,
          builder: (context) => AlertDialog(
                backgroundColor: Colors.green.shade400,
                title: const Text("Yes !"),
                content: const Text("Ta réponse est correcte, bravo !"),
                actions: [
                  OutlinedButton(
                      onPressed: () {
                        Navigator.of(context).pop();
                        setState(() {
                          ct.questionIndex =
                              ct.exeAndProg.progression.nextQuestion;
                        });
                      },
                      style: OutlinedButton.styleFrom(
                          foregroundColor: Colors.black87),
                      child: const Text("Continuer l'exercice"))
                ],
              ));
    } else if (showButtonCorrection) {
      // show a "persitent" snackbar
      ScaffoldMessenger.of(context).showSnackBar(SnackBar(
        backgroundColor: Colors.red.shade200,
        duration: const Duration(seconds: 120),
        content: const Text("Dommage, la réponse est incorrecte."),
        action: SnackBarAction(
            backgroundColor: Colors.white,
            label: "Correction",
            onPressed: _showCorrection),
        actionOverflowThreshold: 0.5,
      ));
    } else if (!isAnswerCorrect) {
      // just show a short snackbar
      ScaffoldMessenger.of(context).showSnackBar(SnackBar(
        backgroundColor: Colors.red.shade200,
        duration: const Duration(seconds: 4),
        content: const Text("Dommage, la réponse est incorrecte."),
      ));
    }
  }

  void _showCorrection() {
    Navigator.of(context).push(MaterialPageRoute<void>(
      builder: (context) => Scaffold(
        appBar: AppBar(),
        body: CorrectionW(widget.controller.currentCorrection,
            Colors.greenAccent, pickQuote()),
      ),
    ));
  }

  void onQuestionButtonClick() {
    // cleanup potential snackbar
    ScaffoldMessenger.of(context).hideCurrentSnackBar();
    switch (widget.controller.status) {
      case ExerciceStatus.answering:
        return _onValideQuestion();
      case ExerciceStatus.displayingFeedback:
        return _onRetryQuestion();
    }
  }

  void _onRetryQuestion() {
    final ct = widget.controller;
    setState(() {
      ct.resetWithNextQuestions(nextQuestions);
    });
  }

  void goToPrevious() {
    // remove potential snackbar
    ScaffoldMessenger.of(context).hideCurrentSnackBar();

    final newIndex = widget.controller.questionIndex == 0
        ? null
        : widget.controller.questionIndex! - 1;
    setState(() {
      widget.controller.questionIndex = newIndex;
    });
  }

  void goToNext() {
    // remove potential snackbar
    ScaffoldMessenger.of(context).hideCurrentSnackBar();

    final newIndex = widget.controller.questionIndex == null
        ? 0
        : widget.controller.questionIndex! + 1;
    setState(() {
      widget.controller.questionIndex = newIndex;
    });
  }
}
