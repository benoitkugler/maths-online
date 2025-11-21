import 'dart:async';

import 'package:eleve/activities/homework/congratulations.dart';
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

abstract class ExerciceAPI {
  /// [evaluate] must evaluate the exercice answers and return the feedback and the new version
  /// of the questions if needed
  Future<EvaluateWorkOut> evaluate(EvaluateWorkIn params);
}

class ExerciceController {
  /// [exeAndProg] stores the server instantiated exercice with
  /// the progression state.
  StudentWork exeAndProg;

  final QuestionRepeat questionRepeat;
  final int questionTimeLimit;

  /// [_questionIndex] is the current question, or -1 for the summary
  int _questionIndex = -1;
  ExerciceStep step = .answering;
  List<QuestionController> _questions = [];

  /// [_newQuestions] is returned from the server when evaluating
  /// and stored while in feedback step.
  /// When going back to summary or retrying, these questions must replace
  /// the one in use
  List<InstantiatedQuestion> _newQuestions = [];

  ExerciceController(
    this.exeAndProg,
    this.questionRepeat,
    this.questionTimeLimit,
  ) {
    _questions = _createControllers();
    _refreshQuestions();
  }

  List<QuestionController> _createControllers() {
    return exeAndProg.exercice.questions.map((qu) {
      final out = QuestionController.fromQuestion(qu.question);
      out.timeout = questionTimeLimit == 0
          ? null
          : Duration(seconds: questionTimeLimit);
      return out;
    }).toList();
  }

  void _refreshQuestions() {
    final status = getQuestionsStatus();
    for (var index = 0; index < _questions.length; index++) {
      final qu = _questions[index];
      final quStatut = status[index];
      if (quStatut.visibility == .notAccessible ||
          quStatut.visibility == .disabled) {
        // mark done
        qu.buttonLabel = quStatut.success == .right
            ? "Question terminée"
            : "Question verrouillée";
        qu.buttonEnabled = false;
        qu.footerQuote = pickQuote();
        qu.setFieldsEnabled(false);
        qu.timeout = null;
      } else if (step == .answering) {
        qu.buttonEnabled = false;
        qu.buttonLabel = "Valider";
        qu.setFieldsEnabled(true);
        qu.timeout = questionTimeLimit == 0
            ? null
            : Duration(seconds: questionTimeLimit);
      }
    }
  }

  // also apply answering step
  void applyNewQuestions() {
    exeAndProg = exeAndProg.copyWithQuestions(_newQuestions);
    _newQuestions = [];
    step = .answering;
    _questions = _createControllers();
    _refreshQuestions();
  }

  /// [updateFromEvaluate] is called after receiving the server
  /// answer. It removes waiting questions and update the questions status.
  void updateFromEvaluate(EvaluateWorkOut resp) {
    exeAndProg = StudentWork(exeAndProg.exercice, resp.progression);
    step = .displayingFeedback;
    _questions[_questionIndex].timeout = null;
    _newQuestions = resp.newQuestions;

    final isCorrect = resp.result.isCorrect;
    if (!isCorrect) {
      // show errors and ask for retry
      /// [_showFeedback] set the given feedback (for the current question)
      /// and set the state to [displayingFeedback]
      _questions[_questionIndex].setFeedback(resp.result.results);

      final isOneTry = questionRepeat == QuestionRepeat.oneTry;
      _questions[_questionIndex].buttonEnabled = !isOneTry;
      _questions[_questionIndex].buttonLabel = isOneTry
          ? "Réponse incorrecte"
          : "Essayer à nouveau...";
      _questions[_questionIndex].setFieldsEnabled(false);
    }

    _refreshQuestions();
  }

  /// setQuestionAnswers show the answers for the current question
  /// and enable the question (only used in the prof preview for now)
  void setQuestionAnswers(Answers answers) {
    _questions[_questionIndex].setAnswers(answers);
    _questions[_questionIndex].buttonEnabled = true;
    _questions[_questionIndex].buttonLabel = "Valider";
  }

  bool isExerciceOver() => exeAndProg.progression.nextQuestion == -1;

  bool get inQuestion => _questionIndex != -1;

  /// valid only in [inQuestion] is true;
  QuestionController get currentQuestion => _questions[_questionIndex];

  bool get goToPreviousEnabled {
    if (_questionIndex <= 0) return false;
    final status = getQuestionsStatus();
    if (status[_questionIndex - 1].visibility == .notAccessible) {
      return false;
    }
    return true;
  }

  Enonce get currentCorrection {
    return exeAndProg.exercice.questions[_questionIndex].question.correction;
  }

  /// checks if the student as already try at least once the current question
  bool get isCurrentCorrectionEnabled {
    return exeAndProg.progression.questions[_questionIndex].isNotEmpty;
  }

  bool get goToNextEnabled {
    final currentIndex = _questionIndex;
    final hasNextQuestion =
        exeAndProg.exercice.questions.isNotEmpty &&
        currentIndex < exeAndProg.exercice.questions.length - 1;
    final status = getQuestionsStatus();
    return hasNextQuestion &&
        status[currentIndex + 1].visibility != .notAccessible;
  }

  // computed from the progression
  List<QuestionStatus> getQuestionsStatus() {
    final exercice = exeAndProg.exercice;
    final progression = exeAndProg.progression;
    final isOneTry = questionRepeat == QuestionRepeat.oneTry;

    return List<QuestionStatus>.generate(exercice.questions.length, (
      questionIndex,
    ) {
      final success = progression.questionSuccess(questionIndex);
      final blockedByFlow =
          exercice.flow == .sequencial &&
          progression.nextQuestion < questionIndex;

      final QVisibility visibility;
      switch (success) {
        case .neverTried:
          visibility = blockedByFlow ? .notAccessible : .enabled;
        case .wrong:
          visibility = (isOneTry || blockedByFlow) ? .notAccessible : .enabled;
        case .right:
          visibility = isOneTry ? .notAccessible : .disabled;
      }

      return (visibility: visibility, success: success);
    });
  }
}

enum ExerciceStep {
  /// the student is working on a question, validation triggers evaluation
  answering,

  /// the wrong answers are displayed in red, validation (if enabled) triggers retry
  displayingFeedback,
}

enum QVisibility { notAccessible, disabled, enabled }

enum QSuccess { neverTried, wrong, right }

// QuestionStatus define two caracteristics about one question:
//  is it "hidden", "accessible but disabled" or "enabled"
//  is it "not answered", "wrongly answered", "correctly answered"
typedef QuestionStatus = ({QVisibility visibility, QSuccess success});

extension on QuestionStatus {
  (Icon, Icon?) get icons {
    const lock = Icon(Icons.lock, color: Colors.grey);
    const right = Icon(Icons.check, color: Colors.green);
    const wrong = Icon(Icons.clear, color: Colors.red);
    switch (visibility) {
      case .notAccessible:
        switch (success) {
          case .neverTried:
            return (lock, null);
          case .right:
            return (right, lock);
          case .wrong:
            return (wrong, lock);
        }
      // icons wise, we don't distinguish the two states
      case .disabled:
      case .enabled:
        switch (success) {
          case .neverTried:
            return (
              const Icon(Icons.assignment, color: Colors.purpleAccent),
              null,
            );
          case .right:
            return (right, null);
          case .wrong:
            return (wrong, null);
        }
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
      baremes,
    );
  }
}

extension on StudentWork {
  StudentWork copyWithQuestions(List<InstantiatedQuestion> questions) {
    return StudentWork(exercice.copyWithQuestions(questions), progression);
  }
}

extension on QuestionAnswersOut {
  /// isCorrect is true if every fields are correct
  bool get isCorrect {
    return results.values.every((success) => success);
  }
}

extension on ProgressionExt {
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
  QSuccess questionSuccess(int index) {
    final qu = getQuestion(index);
    if (_isQuestionCompleted(qu)) return .right;
    return qu.isEmpty ? .neverTried : .wrong;
  }

  /// returns `true` if all the questions of the exercice are completed
  bool isCompleted() {
    return questions.every(_isQuestionCompleted);
  }
}

// Widget

typedef EditorPreviewParams = ({
  void Function() onShowCorrectAnswer,

  /// if true, displays the correction root instead of
  /// the enonce
  bool instantShowCorrection,

  /// if given, skip the summary
  int? initialQuestion,
});

/// ExerciceStartRoute is the widget providing one exercice to
/// the student.
/// It displays an outline of the questions and opens
/// a new route for each question started.
/// It is used in the prof editor preview, and as the base for
/// at home training activity
class ExerciceStartRoute extends StatefulWidget {
  final ExerciceAPI api;

  final ExerciceController controller;

  /// if true, and if the exercice has a correction,
  /// show a button to access the correction after on try
  final bool showCorrectionButtonOnFail;

  /// if true, show an alert about progression not being updated
  final bool noticeSandbox;

  /// optional, if given and reached, shows a dialog
  final DateTime? deadline;

  // Only used in prof editor preview
  final EditorPreviewParams? editorPreviewParams;

  const ExerciceStartRoute(
    this.api,
    this.controller, {
    super.key,
    this.showCorrectionButtonOnFail = true,
    this.noticeSandbox = false,
    this.deadline,
    this.editorPreviewParams,
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
    final qu = widget.editorPreviewParams?.initialQuestion;
    if (qu == null || widget.controller.inQuestion) {
      return;
    }
    WidgetsBinding.instance.addPostFrameCallback((_) => _goToQuestion(qu));
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
        child: Column(
          children: [
            Padding(
              padding: const EdgeInsets.symmetric(
                horizontal: 6.0,
                vertical: 10,
              ),
              child: ColoredTitle(exP.exercice.title, Colors.purple),
            ),
            if (ct.questionRepeat == QuestionRepeat.oneTry)
              Card(
                margin: EdgeInsets.all(8),
                child: ListTile(
                  title: Text("Un seul essai par question !"),
                  trailing: Badge.count(count: 1),
                ),
              ),
            if (ct.questionTimeLimit != 0)
              Card(
                margin: EdgeInsets.all(8),
                child: ListTile(
                  title: Text(
                    "Temps limité à ${ct.questionTimeLimit} sec. par question !",
                  ),
                  trailing: const Icon(Icons.timer),
                ),
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
              child: _QuestionList(exP, ct.getQuestionsStatus(), _goToQuestion),
            ),
          ],
        ),
      ),
    );
  }

  void _goToQuestion(int questionIndex) async {
    final ct = widget.controller;
    ct._questionIndex = questionIndex;
    ct.step = .answering;
    ct._refreshQuestions();

    final goToQuestion = await Navigator.of(context).push(
      MaterialPageRoute<int?>(
        builder: (context) => _QuestionsRoute(
          widget.api,
          ct,
          widget.showCorrectionButtonOnFail,
          noticeSandbox: widget.noticeSandbox,
          deadline: widget.deadline,
          onShowCorrectAnswer: widget.editorPreviewParams?.onShowCorrectAnswer,
          instantShowCorrection:
              widget.editorPreviewParams?.instantShowCorrection ?? false,
        ),
      ),
    );
    if (!mounted) return;

    if (goToQuestion != null) {
      _goToQuestion(goToQuestion);
    } else {
      setState(() {
        // properly apply new questions
        ct.applyNewQuestions();
        // go to summary
        ct._questionIndex = -1;
      });
    }
  }
}

class MarkBareme {
  final int mark;
  final int bareme;
  MarkBareme(this.mark, this.bareme);
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

class _QuestionList extends StatelessWidget {
  final StudentWork data;
  final List<QuestionStatus> status;

  final void Function(int index) onSelectQuestion;

  const _QuestionList(this.data, this.status, this.onSelectQuestion);

  MarkBareme get mark {
    int mark = 0;
    int bareme = 0;
    for (var i = 0; i < data.exercice.baremes.length; i++) {
      bareme += data.exercice.baremes[i];
      if (data.progression.questionSuccess(i) == .right) {
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
              ),
            ],
          ),
        ),
      ),
    );
  }

  bool allowDoQuestion(int questionIndex) {
    return status[questionIndex].visibility != .notAccessible;
  }

  @override
  Widget build(BuildContext context) {
    final mb = mark;
    return ListView(
      children: [
        ...List<Widget>.generate(
          data.exercice.questions.length,
          (index) => _QuestionRow(
            status[index],
            "Question ${index + 1}",
            data.exercice.questions[index].difficulty,
            data.exercice.baremes[index],
            showDetails: () => _showProgressionDetails(context, index),
            onClick: allowDoQuestion(index)
                ? () => onSelectQuestion(index)
                : null,
          ),
        ),
        if (data.progression.questions.isNotEmpty)
          ListTile(
            title: const Text("Total"),
            trailing: Text(
              "${mb.mark} / ${mb.bareme}",
              style: const TextStyle(fontSize: 14),
            ),
          ),
      ],
    );
  }
}

// display one question among the list.
class _QuestionsRoute extends StatefulWidget {
  final ExerciceAPI api;
  final ExerciceController controller;

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

  /// Only used in loopback
  final void Function()? onShowCorrectAnswer;

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
  Timer? timeLimitTimer;

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
    timeLimitTimer?.cancel();
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

    _initQuestionTimer();

    if (widget.instantShowCorrection) {
      WidgetsBinding.instance.addPostFrameCallback((_) => _showCorrection());
    }
  }

  // handle time limit, disabling it when needed
  void _initQuestionTimer() {
    timeLimitTimer?.cancel();
    final timeout = widget.controller.currentQuestion.timeout;
    if (timeout != null) {
      timeLimitTimer = Timer(timeout, _onTimeLimitReached);
    }
  }

  @override
  Widget build(BuildContext context) {
    final ct = widget.controller;
    // disable exiting the question for time limited
    // to prevent abuses
    final canPop =
        ct.step == .displayingFeedback || ct.currentQuestion.timeout == null;
    return Scaffold(
      appBar: AppBar(
        title: Text("Question ${ct._questionIndex + 1}"),
        actions: [
          if (widget.onShowCorrectAnswer != null)
            TextButton(
              onPressed: widget.onShowCorrectAnswer,
              child: const Text("Afficher la réponse"),
            ),
          if (widget.deadline == null) ...[
            // disable when time limit is on
            IconButton(
              onPressed: ct.goToPreviousEnabled ? goToPrevious : null,
              icon: const Icon(Icons.arrow_back),
            ),
            IconButton(
              onPressed: ct.goToNextEnabled ? goToNext : null,
              icon: const Icon(Icons.arrow_forward),
            ),
          ],
        ],
      ),
      body: PopScope(
        canPop: canPop,
        child: QuestionView(
          ct.exeAndProg.exercice.questions[ct._questionIndex].question,
          ct.currentQuestion,
          onQuestionButtonClick,
          Colors.purpleAccent,
        ),
      ),
    );
  }

  void _onTimeLimitReached() async {
    // send an empty response to the server in background
    final evalFuture = _evaluate(QuestionAnswersIn({}));

    await showDialog<void>(
      context: context,
      builder: (context) => AlertDialog(
        backgroundColor: Colors.orange,
        icon: const Icon(Icons.timer),
        title: const Text("Limite dépassée"),
        content: Text(
          "Tu as dépassé la limite de ${widget.controller.questionTimeLimit} sec.",
        ),
      ),
    );
    final resp = await evalFuture;
    if (!mounted) return;

    // go back to exercice home
    if (resp != null) {
      widget.controller.updateFromEvaluate(resp);
    }
    Navigator.of(context).pop();
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
          "Attention, la date de rendu est maintenant dépassée. Ta progression n'est plus enregistrée.",
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.of(context).pop(),
            child: const Text("OK"),
          ),
        ],
      ),
    );
  }

  void onQuestionButtonClick() {
    // cleanup potential snackbar
    ScaffoldMessenger.of(context).hideCurrentSnackBar();

    switch (widget.controller.step) {
      // update UI before http call
      case .answering:
        final question = widget.controller.currentQuestion;
        setState(() {
          question.buttonEnabled = false;
          question.buttonLabel = "Correction...";
        });
        return _onValideQuestion();
      case .displayingFeedback:
        return _onRetryQuestion();
    }
  }

  void _goToQuestion(int questionIndex) {
    Navigator.of(context).pop(questionIndex);
  }

  void _onExerciceOver() async {
    // exercice is over
    await showDialog<bool>(
      context: context,
      builder: (context) => const Dialog(child: Congrats()),
    );
    if (!mounted) return;

    // always go to summary to prevent abuse in oneTry mode
    Navigator.of(context).pop();
  }

  // pure function wrapping the API call with error handling
  Future<EvaluateWorkOut?> _evaluate(QuestionAnswersIn answsers) async {
    final ct = widget.controller;
    final index = ct._questionIndex;

    // validate the given answer
    final params = EvaluateWorkIn(
      ct.exeAndProg.exercice.iD,
      ct.exeAndProg.progression,
      index,
      AnswerP(ct.exeAndProg.exercice.questions[index].params, answsers),
    );

    final EvaluateWorkOut resp;
    try {
      resp = await widget.api.evaluate(params);
    } catch (error) {
      if (!mounted) return null;
      showError("Impossible d'évaluer la question", error, context);
      return null;
    }

    return resp;
  }

  void _onValideQuestion() async {
    final ct = widget.controller;
    final index = ct._questionIndex;

    // always cancel the timeout, if any
    timeLimitTimer?.cancel();

    // for sequencial exercices, if we are not at the current question, just go to it
    // and return
    if (ct.exeAndProg.exercice.flow == Flow.sequencial &&
        index != ct.exeAndProg.progression.nextQuestion) {
      _goToQuestion(ct.exeAndProg.progression.nextQuestion);
      return;
    }

    final resp = await _evaluate(ct.currentQuestion.answers());
    if (resp == null) return; // network error

    if (resp.progression.nextQuestion == -1) {
      _onExerciceOver();
      return;
    }

    setState(() {
      ct.updateFromEvaluate(resp);
    });

    _showValidDialogOrSnack(resp.result.isCorrect);
  }

  // handle the following cases :
  //  - the answer is correct and there is more to do
  //  - the answer is incorrect and there is a correction to display
  //  - the answer is incorrect and there is no correction to display
  void _showValidDialogOrSnack(bool isAnswerCorrect) async {
    final ct = widget.controller;
    final hasNextQuestion =
        isAnswerCorrect && ct.exeAndProg.progression.nextQuestion != -1;
    final showButtonCorrection =
        widget.showCorrectionButtonOnFail &&
        ct.currentCorrection.isNotEmpty &&
        !isAnswerCorrect;

    if (hasNextQuestion) {
      // show a dialog with next button
      final goToNext = await showDialog<bool>(
        context: context,
        builder: (context) => CorrectAnswerDialog(() {
          Navigator.of(context).pop(true);
        }),
      );
      if (!mounted) return;
      if (goToNext ?? false) {
        _goToQuestion(ct.exeAndProg.progression.nextQuestion);
      }
    } else if (showButtonCorrection) {
      // show a "persitent" snackbar
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(
          backgroundColor: Colors.red.shade200,
          duration: const Duration(seconds: 120),
          content: const Text("Dommage, la réponse est incorrecte."),
          action: SnackBarAction(
            backgroundColor: Colors.white,
            label: "Correction",
            onPressed: _showCorrection,
          ),
          actionOverflowThreshold: 0.5,
        ),
      );
    } else if (!isAnswerCorrect) {
      // just show a short snackbar
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(
          backgroundColor: Colors.red.shade200,
          duration: const Duration(seconds: 4),
          content: const Text("Dommage, la réponse est incorrecte."),
        ),
      );
    }
  }

  void _showCorrection() {
    Navigator.of(context).push(
      MaterialPageRoute<void>(
        builder: (context) => Scaffold(
          appBar: AppBar(),
          body: CorrectionView(
            widget.controller.currentCorrection,
            Colors.greenAccent,
            pickQuote(),
          ),
        ),
      ),
    );
  }

  void _onRetryQuestion() {
    setState(() {
      widget.controller.applyNewQuestions();
      _initQuestionTimer(); // after building the questions again to have the proper timeout
    });
  }

  void goToPrevious() {
    // remove potential snackbar
    ScaffoldMessenger.of(context).hideCurrentSnackBar();

    _goToQuestion(widget.controller._questionIndex - 1);
  }

  void goToNext() {
    // remove potential snackbar
    ScaffoldMessenger.of(context).hideCurrentSnackBar();

    _goToQuestion(widget.controller._questionIndex + 1);
  }
}

const _difficulties = {
  DifficultyTag.diff1: "★",
  DifficultyTag.diff2: "★★",
  DifficultyTag.diff3: "★★★",
  DifficultyTag.diffEmpty: "",
};

class _QuestionRow extends StatelessWidget {
  final QuestionStatus state;
  final String title;
  final DifficultyTag difficultyTag;
  final int bareme;
  final void Function() showDetails;
  final void Function()? onClick;

  const _QuestionRow(
    this.state,
    this.title,
    this.difficultyTag,
    this.bareme, {
    required this.showDetails,
    required this.onClick,
  });

  @override
  Widget build(BuildContext context) {
    final diff = _difficulties[difficultyTag] ?? "";
    final (icon1, icon2) = state.icons;
    final btn1 = OutlinedButton(onPressed: showDetails, child: icon1);
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 4.0, vertical: 2),
      child: ListTile(
        shape: const RoundedRectangleBorder(
          borderRadius: BorderRadius.all(Radius.circular(4)),
        ),
        tileColor: state.visibility != .notAccessible && state.success != .right
            ? Colors.purple.shade400.withValues(alpha: 0.5)
            : null,
        leading: icon2 == null
            ? btn1
            : Column(mainAxisSize: .min, children: [btn1, icon2]),
        title: Text(title),
        subtitle: diff.isEmpty ? null : Text(diff),
        trailing: Text("/ $bareme"),
        onTap: onClick,
      ),
    );
  }
}
