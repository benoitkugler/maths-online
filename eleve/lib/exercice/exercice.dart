import 'dart:async';

import 'package:eleve/exercice/congratulations.dart';
import 'package:eleve/questions/question.dart';
import 'package:eleve/quotes.dart';
import 'package:eleve/shared/errors.dart';
import 'package:eleve/shared/title.dart';
import 'package:eleve/types/src.dart';
import 'package:eleve/types/src_maths_questions_client.dart';
import 'package:eleve/types/src_sql_editor.dart';
import 'package:eleve/types/src_sql_homework.dart';
import 'package:eleve/types/src_sql_tasks.dart';
import 'package:eleve/types/src_tasks.dart';
import 'package:flutter/material.dart' hide Flow;

class NotificationExerciceDone extends Notification {}

abstract class ExerciceAPI {
  /// [evaluate] must evaluate the exercice answers, according
  /// to the exercice mode, and returns the feedback and the new version
  /// of the questions if needed
  Future<EvaluateWorkOut> evaluate(EvaluateWorkIn params);
}

enum ExerciceStatus {
  /// the student is working on a question, validation triggers evaluation
  answering,

  /// the wrong answers are displayed in red, validation (if enabled) triggers retry
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

  List<QuestionStatus> questionStatus =
      []; // computed from the saved progression and the potential current error

  final QuestionRepeat questionRepeat;
  final int questionTimeLimit;

  /// [_questionIndex] is the current question, or -1 for the summary
  int _questionIndex = -1;

  int get questionIndex => _questionIndex;

  bool get inQuestion => _questionIndex != -1;

  ExerciceStatus status = ExerciceStatus.answering;

  List<QuestionController> _questions;

  ExerciceController(
      this.exeAndProg, this.questionRepeat, this.questionTimeLimit)
      : _questions = [] {
    _questions = _createControllers();
    _refreshStates();
  }

  Enonce get currentCorrection {
    return exeAndProg.exercice.questions[_questionIndex].question.correction;
  }

  /// checks if the student as already try at least once the current question
  bool get isCurrentCorrectionEnabled {
    return exeAndProg.progression.questions[_questionIndex].isNotEmpty;
  }

  List<QuestionController> _createControllers() {
    return exeAndProg.exercice.questions.map((qu) {
      final out = QuestionController.fromQuestion(qu.question);
      out.timeout =
          questionTimeLimit == 0 ? null : Duration(seconds: questionTimeLimit);
      return out;
    }).toList();
  }

  void onQuestionValid() {
    switch (status) {
      case ExerciceStatus.answering:
        _questions[_questionIndex].buttonEnabled = false;
        _questions[_questionIndex].buttonLabel = "Correction...";
        break;
      case ExerciceStatus.displayingFeedback:
        _questions[_questionIndex].buttonEnabled = true;
        _questions[_questionIndex].buttonLabel = "Valider...";
        break;
    }
    _refreshStates();
  }

  /// [showFeedback] set the given feedback (for the current question)
  /// and set the state to displayingFeedback
  void showFeedback(QuestionAnswersOut feedback) {
    status = ExerciceStatus.displayingFeedback;
    _questions[_questionIndex].setFeedback(feedback.results);

    final isOneTry = questionRepeat == QuestionRepeat.oneTry;
    _questions[_questionIndex].buttonEnabled = !isOneTry;
    _questions[_questionIndex].buttonLabel =
        isOneTry ? "Réponse incorrecte" : "Essayer à nouveau...";
    _questions[_questionIndex].setFieldsEnabled(false);
    _refreshStates();
  }

  void reset() {
    _refreshStates();
    status = ExerciceStatus.answering;
  }

  void resetWithNextQuestions(List<InstantiatedQuestion> nextQuestions) {
    exeAndProg = exeAndProg.copyWithQuestions(nextQuestions);
    _questions = _createControllers();
    reset();
  }

  void _refreshStates() {
    questionStatus = _computeQuestionStates();
    for (var index = 0; index < _questions.length; index++) {
      final qu = _questions[index];
      if (questionStatus[index] == QuestionStatus.checked ||
          questionStatus[index] == QuestionStatus.incorrectAndLocked) {
        // mark done
        qu.buttonLabel = "Question terminée";
        qu.buttonEnabled = false;
        qu.footerQuote = pickQuote();
        qu.setFieldsEnabled(false);
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

      final alreadyTried =
          exP.progression.getQuestion(questionIndex).isNotEmpty;
      if (questionRepeat == QuestionRepeat.oneTry && alreadyTried) {
        // here we now the status is failed
        return QuestionStatus.incorrectAndLocked;
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
  /// and enable the question
  void setQuestionAnswers(Answers answers) {
    _questions[_questionIndex].setAnswers(answers);
    _questions[_questionIndex].buttonEnabled = true;
    _questions[_questionIndex].buttonLabel = "Valider";
  }

  bool isExerciceOver() => exeAndProg.progression.nextQuestion == -1;

  void setQuestionIndex(int newIndex) {
    _questionIndex = newIndex;
  }

  bool get goToPreviousEnabled {
    if (_questionIndex == 0) return false;
    if (questionStatus[_questionIndex - 1] == QuestionStatus.locked) {
      return false;
    }
    return true;
  }

  bool get goToNextEnabled {
    final currentIndex = _questionIndex;
    final hasNextQuestion = exeAndProg.exercice.questions.isNotEmpty &&
        currentIndex != exeAndProg.exercice.questions.length - 1;
    // always allow to go to already done question
    if (hasNextQuestion &&
        questionStatus[currentIndex + 1] == QuestionStatus.checked) {
      return true;
    }
    switch (exeAndProg.exercice.flow) {
      case Flow
          .sequencial: // do not show locked questions when exercice is not over
        return (isExerciceOver() && hasNextQuestion) ||
            currentIndex < exeAndProg.progression.nextQuestion;
      case Flow.parallel: // no restriction:
        return hasNextQuestion;
    }
  }
}

extension on InstantiatedWork {
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

extension on StudentWork {
  StudentWork copyWithQuestions(List<InstantiatedQuestion> questions) {
    return StudentWork(exercice.copyWithQuestions(questions), progression);
  }
}

/// ExerciceW is the widget providing one exercice to
/// the student.
/// It displays an outline of the questions and opens
/// a new route for each question started.
/// It is used in the editor loopback, and as the base for
/// at home training activity
class ExerciceStartRoute extends StatefulWidget {
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

  /// if true, show an alert about progression not being updated
  final bool noticeSandbox;

  /// optional, if given and reached, shows a dialog
  final DateTime? deadline;

  /// if given, skip the summary
  final int? initialQuestion;

  const ExerciceStartRoute(
    this.api,
    this.controller, {
    super.key,
    this.initialQuestion,
    this.onShowCorrectAnswer,
    this.showCorrectionButtonOnFail = true,
    this.instantShowCorrection = false,
    this.noticeSandbox = false,
    this.deadline,
  });

  @override
  State<ExerciceStartRoute> createState() => _ExerciceStartRouteState();
}

class _ExerciceStartRouteState extends State<ExerciceStartRoute> {
  @override
  void didUpdateWidget(covariant ExerciceStartRoute oldWidget) {
    _handleInitialQuestion();
    super.didUpdateWidget(oldWidget);
  }

  @override
  void initState() {
    _handleInitialQuestion();
    super.initState();
  }

  void _handleInitialQuestion() {
    if (widget.initialQuestion == null || widget.controller.inQuestion) return;
    print("handle initiatl");
    WidgetsBinding.instance
        .addPostFrameCallback((_) => _goToQuestion(widget.initialQuestion!));
  }

  @override
  Widget build(BuildContext context) {
    final ct = widget.controller;
    final exP = widget.controller.exeAndProg;

    /// show a welcome screen when opening an exercice,
    /// with its questions and bareme
    return Scaffold(
      appBar: AppBar(title: const Text("Exercice")),
      body: Center(
        child: Column(children: [
          Padding(
            padding: const EdgeInsets.symmetric(horizontal: 6.0, vertical: 10),
            child: ColoredTitle(exP.exercice.title, Colors.purple),
          ),
          if (ct.questionRepeat == QuestionRepeat.oneTry)
            ListTile(
                title: Text("Un seul essai par question !"),
                trailing: Badge.count(count: 1)),
          if (ct.questionTimeLimit != 0)
            ListTile(
              title: Text(
                  "Temps limité à $ct.questionTimeLimit sec. par question !"),
              trailing: const Icon(Icons.timer),
            ),
          if (widget.noticeSandbox)
            const Card(
              margin: EdgeInsets.only(bottom: 10),
              child: Padding(
                padding: EdgeInsets.all(8.0),
                child: Text("Ta progression n'est pas enregistrée."),
              ),
            ),
          Expanded(
            child: _QuestionList(exP, ct.questionStatus, _goToQuestion),
          )
        ]),
      ),
    );
  }

  void _goToQuestion(int questionIndex) async {
    final ct = widget.controller;
    ct.setQuestionIndex(questionIndex);
    await Navigator.of(context).push(MaterialPageRoute<void>(
        builder: (context) => _QuestionsRoute(
              widget.api,
              ct,
              widget.showCorrectionButtonOnFail,
              onShowCorrectAnswer: widget.onShowCorrectAnswer,
              noticeSandbox: widget.noticeSandbox,
              deadline: widget.deadline,
              instantShowCorrection: widget.instantShowCorrection,
            )));
    if (!mounted) return;
    ct.setQuestionIndex(-1);
    // controller has need updated
    setState(() {});
  }
}

// display one question among the list.
// include the list so that question may be
// switched
class _QuestionsRoute extends StatefulWidget {
  final ExerciceAPI api;
  final ExerciceController controller;

  final void Function()? onShowCorrectAnswer;

  /// if true, and if the exercice has a correction,
  /// show a button to access the correction after on try
  final bool showCorrectionButtonOnFail;

  /// if true, show an alert about progression not being updated
  final bool noticeSandbox;

  /// optional, if given and reached, shows a dialog
  final DateTime? deadline;

  /// if true, displays the correction root instead of
  /// the enonce
  /// Only used in loopback
  final bool instantShowCorrection;

  const _QuestionsRoute(
    this.api,
    this.controller,
    this.showCorrectionButtonOnFail, {
    required this.onShowCorrectAnswer,
    required this.noticeSandbox,
    required this.deadline,
    this.instantShowCorrection = false,
  });

  @override
  State<_QuestionsRoute> createState() => _QuestionsRouteState();
}

class _QuestionsRouteState extends State<_QuestionsRoute> {
  bool noticeSandbox = false;
  Timer? deadlineTimer;

  // the questions to display when trying again
  // this is a temporay slice, affected to the controller on user validation
  List<InstantiatedQuestion> nextQuestions = [];

  @override
  void initState() {
    _init();
    super.initState();
  }

  @override
  void didUpdateWidget(covariant _QuestionsRoute oldWidget) {
    _init();
    super.didUpdateWidget(oldWidget);
  }

  @override
  void dispose() {
    deadlineTimer?.cancel();
    super.dispose();
  }

  void _init() {
    noticeSandbox = widget.noticeSandbox;

    // cancel the timer
    deadlineTimer?.cancel();
    final d = widget.deadline;
    if (d != null) {
      final now = DateTime.now();
      if (d.isAfter(now)) {
        deadlineTimer = Timer(d.difference(DateTime.now()), _onDeadlineReached);
      }
    }

    widget.controller.reset();
    if (widget.instantShowCorrection) {
      WidgetsBinding.instance.addPostFrameCallback((_) => _showCorrection());
    }
  }

  @override
  Widget build(BuildContext context) {
    final ct = widget.controller;
    return Scaffold(
      appBar: AppBar(
        title: Text("Question ${ct.questionIndex + 1}"),
        actions: [
          if (widget.onShowCorrectAnswer != null)
            TextButton(
                onPressed: widget.onShowCorrectAnswer,
                child: const Text("Afficher la réponse")),
          IconButton(
              onPressed: ct.goToPreviousEnabled ? goToPrevious : null,
              icon: const Icon(Icons.arrow_back)),
          IconButton(
              onPressed: ct.goToNextEnabled ? goToNext : null,
              icon: const Icon(Icons.arrow_forward)),
        ],
      ),
      body: QuestionView(
        ct.exeAndProg.exercice.questions[ct.questionIndex].question,
        ct._questions[ct.questionIndex],
        onQuestionButtonClick,
        Colors.purpleAccent,
      ),
    );
  }

  void _onDeadlineReached() {
    if (!mounted) return;

    setState(() {
      noticeSandbox = true;
    });
    showDialog<void>(
        context: context,
        builder: (context) => AlertDialog(
              title: const Text("Date de rendu"),
              content: const Text(
                  "Attention, la date de rendu est maintenant dépassée. Ta progression n'est plus enregistrée."),
              actions: [
                TextButton(
                    onPressed: () => Navigator.of(context).pop(),
                    child: const Text("OK"))
              ],
            ));
  }

  void _goToQuestion(int questionIndex) {
    final ct = widget.controller;
    setState(() {
      ct.setQuestionIndex(questionIndex);
    });
  }

  // handle the errors
  Future<EvaluateWorkOut?> _evaluate(EvaluateWorkIn params) async {
    try {
      final res = await widget.api.evaluate(params);
      return res;
    } catch (error) {
      if (!mounted) return null;
      showError("Impossible d'évaluer l'exercice", error, context);
      return null;
    }
  }

  void onExerciceOver() async {
    // exercice is over
    final goBack = await showDialog<bool>(
        context: context,
        builder: (context) => const Dialog(child: Congrats()));
    if (!mounted) return;
    NotificationExerciceDone().dispatch(context);
    if (goBack ?? false) {
      // go to summary
      Navigator.of(context).pop();
    }
  }

  void _onValideQuestion() async {
    final ct = widget.controller;
    final index = ct.questionIndex;

    // for sequencial exercices, if we are not at the current question, just go to it
    // and return
    if (ct.exeAndProg.exercice.flow == Flow.sequencial &&
        index != ct.exeAndProg.progression.nextQuestion) {
      _goToQuestion(ct.exeAndProg.progression.nextQuestion);
      return;
    }

    // validate the given answer
    final resp = await _evaluate(EvaluateWorkIn(
      ct.exeAndProg.exercice.iD,
      ct.exeAndProg.progression,
      index,
      AnswerP(ct.exeAndProg.exercice.questions[index].params,
          ct._questions[index].answers()),
    ));
    if (resp == null || !mounted) {
      return;
    }

    ct.updateProgression(resp.progression); // update the progression
    nextQuestions = resp.newQuestions; // buffer until retry

    if (ct.exeAndProg.progression.nextQuestion == -1) {
      onExerciceOver();
      return;
    }

    final isCorrect = resp.result.isCorrect;
    if (!isCorrect) {
      // show errors and ask for retry
      setState(() {
        ct.showFeedback(resp.result);
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
      final goToNext = await showDialog<bool>(
          context: context,
          builder: (context) => CorrectAnswerDialog(() {
                Navigator.of(context).pop(true);
              }));
      if (!mounted) return;
      if (goToNext ?? false) {
        _goToQuestion(ct.exeAndProg.progression.nextQuestion);
      }
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
        body: CorrectionView(widget.controller.currentCorrection,
            Colors.greenAccent, pickQuote()),
      ),
    ));
  }

  void onQuestionButtonClick() {
    // cleanup potential snackbar
    ScaffoldMessenger.of(context).hideCurrentSnackBar();

    setState(() {
      widget.controller.onQuestionValid();
    });
    switch (widget.controller.status) {
      case ExerciceStatus.answering:
        return _onValideQuestion();
      case ExerciceStatus.displayingFeedback:
        return _onRetryQuestion();
    }
  }

  void _onRetryQuestion() {
    setState(() {
      widget.controller.resetWithNextQuestions(nextQuestions);
    });
  }

  void goToPrevious() {
    // remove potential snackbar
    ScaffoldMessenger.of(context).hideCurrentSnackBar();

    _goToQuestion(widget.controller.questionIndex - 1);
  }

  void goToNext() {
    // remove potential snackbar
    ScaffoldMessenger.of(context).hideCurrentSnackBar();

    _goToQuestion(widget.controller.questionIndex + 1);
  }
}

extension IsCorrect on QuestionAnswersOut {
  /// isCorrect is true if every fields are correct
  bool get isCorrect {
    return results.values.every((success) => success);
  }
}

extension ProgressionExtension on ProgressionExt {
  /// [getQuestion] returns an empty list if progression is empty
  QuestionHistory getQuestion(int index) {
    if (questions.length <= index) {
      return [];
    }
    return questions[index];
  }

  bool _isQuestionCompleted(List<bool> history) {
    return history.isNotEmpty && history.last;
  }

  /// returns `true` if the question at [index] is completed
  bool isQuestionCompleted(int index) {
    return _isQuestionCompleted(getQuestion(index));
  }

  /// returns `true` if all the questions of the exercice are completed
  bool isCompleted() {
    return questions.every(_isQuestionCompleted);
  }
}

class _SuccessSquare extends StatelessWidget {
  final bool success;
  const _SuccessSquare(this.success);

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.all(8.0),
      child: Container(
        height: 30,
        width: 30,
        color: success ? Colors.green : Colors.red,
      ),
    );
  }
}

class MarkBareme {
  final int mark;
  final int bareme;
  MarkBareme(this.mark, this.bareme);
}

class _QuestionList extends StatelessWidget {
  final StudentWork data;
  final List<QuestionStatus> states;

  final void Function(int index) onSelectQuestion;

  const _QuestionList(this.data, this.states, this.onSelectQuestion);

  MarkBareme get mark {
    int mark = 0;
    int bareme = 0;
    for (var i = 0; i < data.exercice.baremes.length; i++) {
      bareme += data.exercice.baremes[i];
      if (data.progression.isQuestionCompleted(i)) {
        mark += data.exercice.baremes[i];
      }
    }
    return MarkBareme(mark, bareme);
  }

  void _showProgressionDetails(BuildContext context, int questionIndex) {
    showDialog<void>(
        context: context,
        builder: (context) => Dialog(
              child: Card(
                child: Column(
                  mainAxisSize: MainAxisSize.min,
                  children: [
                    const Padding(
                      padding: EdgeInsets.all(20.0),
                      child: Text(
                        "Historique de tes tentatives",
                        style: TextStyle(fontSize: 20),
                      ),
                    ),
                    Padding(
                      padding: const EdgeInsets.all(8.0),
                      child: Wrap(
                        children: data.progression
                            .getQuestion(questionIndex)
                            .map((e) => _SuccessSquare(e))
                            .toList(),
                      ),
                    )
                  ],
                ),
              ),
            ));
  }

  bool allowDoQuestion(int questionIndex) {
    // if the question has been validated, always allow access
    if (states[questionIndex] == QuestionStatus.checked) {
      return true;
    }
    if (states[questionIndex] == QuestionStatus.incorrectAndLocked) {
      return false;
    }

    switch (data.exercice.flow) {
      case Flow.sequencial:
        return data.progression.nextQuestion == questionIndex;
      case Flow.parallel:
        return true;
    }
  }

  @override
  Widget build(BuildContext context) {
    final mb = mark;
    return ListView(
      children: [
        ...List<Widget>.generate(
          data.exercice.questions.length,
          (index) => _QuestionRow(
            states[index],
            "Question ${index + 1}",
            data.exercice.questions[index].difficulty,
            data.exercice.baremes[index],
            showDetails: () => _showProgressionDetails(context, index),
            onClick:
                allowDoQuestion(index) ? () => onSelectQuestion(index) : null,
          ),
        ),
        if (data.progression.questions.isNotEmpty)
          ListTile(
            title: const Text("Total"),
            trailing: Text("${mb.mark} / ${mb.bareme}",
                style: const TextStyle(fontSize: 14)),
          )
      ],
    );
  }
}

enum QuestionStatus { locked, checked, toDo, incorrect, incorrectAndLocked }

extension _Icon on QuestionStatus {
  Icon get icon {
    switch (this) {
      case QuestionStatus.locked:
        return const Icon(Icons.lock, color: Colors.grey);
      case QuestionStatus.checked:
        return const Icon(Icons.check, color: Colors.green);
      case QuestionStatus.toDo:
        return const Icon(Icons.assignment, color: Colors.purpleAccent);
      case QuestionStatus.incorrect:
        return const Icon(IconData(0xf647, fontFamily: 'MaterialIcons'),
            color: Colors.red);
      case QuestionStatus.incorrectAndLocked:
        return const Icon(Icons.lock, color: Colors.red);
    }
  }
}

const _difficulties = {
  DifficultyTag.diff1: "★",
  DifficultyTag.diff2: "★★",
  DifficultyTag.diff3: "★★★",
  DifficultyTag.diffEmpty: ""
};

class _QuestionRow extends StatelessWidget {
  final QuestionStatus state;
  final String title;
  final DifficultyTag difficultyTag;
  final int bareme;
  final void Function() showDetails;
  final void Function()? onClick;

  const _QuestionRow(this.state, this.title, this.difficultyTag, this.bareme,
      {required this.showDetails, required this.onClick});

  @override
  Widget build(BuildContext context) {
    final diff = _difficulties[difficultyTag] ?? "";
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 4.0),
      child: ListTile(
        shape: const RoundedRectangleBorder(
            borderRadius: BorderRadius.all(Radius.circular(4))),
        tileColor: state == QuestionStatus.toDo
            ? Colors.purple.shade400.withValues(alpha: 0.5)
            : null,
        leading: OutlinedButton(onPressed: showDetails, child: state.icon),
        title: Text(title),
        subtitle: diff.isEmpty ? null : Text(diff),
        trailing: Text("/ $bareme"),
        onTap: onClick,
      ),
    );
  }
}
