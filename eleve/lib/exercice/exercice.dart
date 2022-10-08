import 'package:eleve/exercice/congratulations.dart';
import 'package:eleve/exercice/home.dart';
import 'package:eleve/questions/fields.dart';
import 'package:eleve/questions/question.dart';
import 'package:eleve/questions/types.gen.dart' hide Answer;
import 'package:eleve/quotes.dart';
import 'package:eleve/shared/errors.dart';
import 'package:eleve/shared_gen.dart';
import 'package:flutter/material.dart' hide Flow;

class NotificationExerciceDone extends Notification {}

abstract class ExerciceAPI extends FieldAPI {
  /// [evaluate] must evaluate the exercice answers, according
  /// to the exercice mode, and returns the feedback and the new version
  /// of the questions if needed
  Future<EvaluateWorkOut> evaluate(EvaluateWorkIn params);
}

class ExerciceQuestionController extends BaseQuestionController {
  void Function() onClick;

  ExerciceQuestionController(Question question, FieldAPI api, this.onClick)
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

enum ExerciceState {
  /// the student is working on a question, validation triggers evaluation
  answering,

  /// the wrong answers are displayed in red, validation triggers retry
  displayingFeedback
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

  ExerciceState state = ExerciceState.answering;

  List<QuestionStatus> questionStates = [];

  List<ExerciceQuestionController> _questions;

  // questions validated but not submitted
  Set<int> waitingQuestions = {};

  /// onValid is the callback triggered when validatin a question
  /// and is filled by the widget using the controller
  void Function() onValid;

  ExerciceController(this.exeAndProg, this.questionIndex, FieldAPI api)
      : onValid = (() => {}),
        _questions = [] {
    _questions = exeAndProg.exercice.questions
        .map((qu) => ExerciceQuestionController(qu.question, api, onValid))
        .toList();
    _refreshStates();
  }

  /// [showFeedback] set the given feedback (for many questions)
  /// and set the state to displayingFeedback
  void showFeedback(Map<int, QuestionAnswersOut> feedbacks) {
    state = ExerciceState.displayingFeedback;
    for (var question in feedbacks.entries) {
      final index = question.key;
      final feedback = question.value.results;
      _questions[index].setFeedback(feedback);
    }
    for (var question in _questions) {
      question.state.buttonEnabled = true;
      question.state.buttonLabel = "Essayer à nouveau...";
    }
  }

  void reset() {
    for (var question in _questions) {
      question.setFeedback(null);
    }
    waitingQuestions.clear();
    _refreshStates();
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

      // after validating, a question may be both incorrect and waiting submit:
      // give the priority to incorrectQuestions
      final qu = _questions[questionIndex];
      if (!qu.feedback().values.every((success) => success)) {
        return QuestionStatus.incorrect;
      } else if (waitingQuestions.contains(questionIndex)) {
        return QuestionStatus.waitingCorrection;
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

  const ExerciceW(this.api, this.controller,
      {Key? key, this.onShowCorrectAnswer})
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
    reset();
    super.initState();
  }

  @override
  void didUpdateWidget(covariant ExerciceW oldWidget) {
    reset();
    super.didUpdateWidget(oldWidget);
  }

  void reset() {
    nextQuestions = [];
    widget.controller.reset();
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
    NotificationExerciceDone().dispatch(context);
  }

  void onValidQuestionSequential() async {
    final ct = widget.controller;
    final index = ct.questionIndex!;

    // if we are not at the current question, just go to it
    // and return
    if (index != ct.exeAndProg.progression.nextQuestion) {
      setState(() {
        ct.questionIndex = ct.exeAndProg.progression.nextQuestion;
      });
      return;
    }

    // validate the given answer
    final resp = await _evaluate(EvaluateWorkIn(
        ct.exeAndProg.exercice.iD,
        {
          index: Answer(ct.exeAndProg.exercice.questions[index].params,
              ct._questions[index].answers())
        },
        ct.exeAndProg.progression));
    if (resp == null) {
      return;
    }

    ct.exeAndProg = StudentWork(
        ct.exeAndProg.exercice, resp.progression); // update the progression
    nextQuestions = resp.newQuestions; // buffer until retry

    final isCorrect = resp.results[index]!.isCorrect;
    final hasNextQuestion =
        isCorrect && ct.exeAndProg.progression.nextQuestion != -1;
    ScaffoldMessenger.of(context).showSnackBar(SnackBar(
      backgroundColor: isCorrect ? Colors.lightGreen : Colors.red.shade200,
      duration: Duration(
          seconds: hasNextQuestion ? 10000 : 4), // block on correct answer
      content:
          Text(isCorrect ? "Bonne réponse ! Bravo." : "Réponse incorrecte"),
      action: hasNextQuestion
          ? SnackBarAction(
              label: "Continuer l'exercice",
              textColor: Colors.black,
              onPressed: () {
                setState(() {
                  ct.questionIndex = ct.exeAndProg.progression.nextQuestion;
                });
              },
            )
          : null,
    ));

    if (ct.exeAndProg.progression.nextQuestion == -1) {
      onExerciceOver();
      return;
    }

    if (isCorrect) {
      // wait for the user to go to the next question
    } else {
      // show errors and ask for retry
      setState(() {
        ct.showFeedback({index: resp.results[index]!});
      });
    }
  }

  // only validate if all the questions have been completed
  void onValidQuestionParallel() async {
    final ct = widget.controller;
    final questions = ct.exeAndProg.exercice.questions;

    // check if all the questions are done
    final toSend = <int, Answer>{};
    int? goToQuestion;
    for (var index = 0; index < questions.length; index++) {
      final history = ct.exeAndProg.progression.getQuestion(index);
      if (history.isEmpty || !history.last) {
        // the question must be answered
        if (!ct.waitingQuestions.contains(index)) {
          // go to this question
          goToQuestion = index;
          break;
        } else {
          // add it to the send answsers
          toSend[index] =
              Answer(questions[index].params, ct._questions[index].answers());
        }
      }
    }

    if (goToQuestion != null) {
      // there are still some questions to to
      setState(() {
        widget.controller.questionIndex = goToQuestion;
      });
      return;
    }

    // all good, lets send the results
    final resp = await _evaluate(EvaluateWorkIn(
        ct.exeAndProg.exercice.iD, toSend, ct.exeAndProg.progression));
    if (resp == null) {
      return;
    }

    ct.exeAndProg = StudentWork(
        ct.exeAndProg.exercice, resp.progression); // update the progression
    nextQuestions = resp.newQuestions; // buffer until retry

    if (isExerciceOver) {
      onExerciceOver();
      return;
    }

    // display the incorrect answers
    final wrongAnswersPlus1 = resp.results.keys
        .where((index) => !resp.results[index]!.isCorrect)
        .map((e) => e + 1);
    final text = (wrongAnswersPlus1.length > 1)
        ? "Les questions ${wrongAnswersPlus1.join(', ')} sont incorrectes."
        : "La question ${wrongAnswersPlus1.first} est incorrecte.";

    ScaffoldMessenger.of(context).showSnackBar(SnackBar(
      backgroundColor: Colors.red.shade200,
      duration: const Duration(seconds: 4),
      content: Text(text),
    ));

    // display the errors and go to the menu
    setState(() {
      ct.showFeedback(resp.results);
      widget.controller.questionIndex = null;
    });
  }

  void onQuestionClick() {
    switch (widget.controller.state) {
      case ExerciceState.answering:
        return _onValideQuestion();
      case ExerciceState.displayingFeedback:
        return _onRetryQuestion();
    }
  }

  void _onRetryQuestion() {
    final ct = widget.controller;
    setState(() {
      ct.reset();
      ct.waitingQuestions.clear();
      ct.exeAndProg = ct.exeAndProg.copyWithQuestions(nextQuestions);
    });
  }

  void _onValideQuestion() async {
    widget.controller.waitingQuestions.add(widget.controller.questionIndex!);
    switch (widget.controller.exeAndProg.exercice.flow) {
      case Flow.sequencial:
        return onValidQuestionSequential();
      case Flow.parallel:
        return onValidQuestionParallel();
    }
  }

  bool get isExerciceOver =>
      widget.controller.exeAndProg.progression.nextQuestion == -1;

  bool get goToPreviousEnabled => widget.controller.questionIndex != null;

  bool get goToNextEnabled {
    final ex = widget.controller.exeAndProg;
    final currentIndex = widget.controller.questionIndex;
    final hasNextQuestion = currentIndex != ex.exercice.questions.length - 1;
    switch (ex.exercice.flow) {
      case Flow
          .sequencial: // do not show locked questions when exercice is not over
        return currentIndex == null ||
            (isExerciceOver && hasNextQuestion) ||
            currentIndex < ex.progression.nextQuestion;
      case Flow.parallel: // no restriction:
        return hasNextQuestion;
    }
  }

  void goToPrevious() {
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

  Future<bool> confirmeLeave() async {
    final res = await showDialog<bool?>(
        context: context,
        builder: (BuildContext context) {
          return AlertDialog(
            title: const Text("Abandonner l'exercice"),
            content: const Text(
                "Es-tu sûr d'abandonner l'exercice ? \nTes réponses en attente de correction seront effacées."),
            actions: [
              TextButton(
                  child: const Text("Abandonner"),
                  onPressed: () {
                    Navigator.pop(context, true);
                  })
            ],
          );
        });
    return res ?? false;
  }

  @override
  Widget build(BuildContext context) {
    final ct = widget.controller;
    final exP = widget.controller.exeAndProg;
    return Scaffold(
        appBar: AppBar(
          title: const Text("Exercice"),
          actions: [
            if (widget.onShowCorrectAnswer != null &&
                widget.controller.questionIndex != null)
              TextButton(
                  onPressed: widget.onShowCorrectAnswer,
                  child: const Text("Afficher la réponse.")),
            IconButton(
                onPressed: goToPreviousEnabled ? goToPrevious : null,
                icon: const Icon(IconData(0xf572,
                    fontFamily: 'MaterialIcons', matchTextDirection: true))),
            IconButton(
                onPressed: goToNextEnabled ? goToNext : null,
                icon: const Icon(IconData(0xf57a,
                    fontFamily: 'MaterialIcons', matchTextDirection: true))),
          ],
        ),
        body: WillPopScope(
          onWillPop: () async {
            final isDirty = ct.questionStates
                .any((st) => st == QuestionStatus.waitingCorrection);
            if (isDirty) {
              final ok = await confirmeLeave();
              return ok;
            }
            return true; // nothing to loose
          },
          child: widget.controller.questionIndex == null
              ? ExerciceHome(
                  exP,
                  ct.questionStates,
                  (index) => setState(() {
                        widget.controller.questionIndex = index;
                      }))
              : QuestionW(
                  widget
                      .controller._questions[widget.controller.questionIndex!],
                  Colors.purpleAccent,
                  title: "Question ${widget.controller.questionIndex! + 1}",
                ),
        ));
  }
}
