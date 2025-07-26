import 'package:eleve/activities/ceintures/api.dart';
import 'package:eleve/exercice/exercice.dart';
import 'package:eleve/questions/question.dart';
import 'package:eleve/quotes.dart';
import 'package:eleve/shared/errors.dart';
import 'package:eleve/shared/title.dart';
import 'package:eleve/types/src_maths_questions_client.dart';
import 'package:eleve/types/src_prof_ceintures.dart';
import 'package:eleve/types/src_sql_ceintures.dart';
import 'package:eleve/types/src_tasks.dart';
import 'package:flutter/material.dart';

class SeanceAnswers {
  final List<IdBeltquestion> questions;
  final List<AnswerP> answers;
  const SeanceAnswers(this.questions, this.answers);
}

class Seance extends StatefulWidget {
  final CeinturesAPI api;
  final StudentTokens tokens;
  final Stage stage;

  final void Function(bool isSuccess, StudentEvolution newEvolution) onValid;

  const Seance(this.api, this.tokens, this.stage, this.onValid, {super.key});

  @override
  State<Seance> createState() => SeanceState();
}

class SeanceState extends State<Seance> {
  late Future<SelectQuestionsOut> loader;
  // SeanceController? controller;

  @override
  void initState() {
    loader = widget.api
        .selectQuestions(SelectQuestionsIn(widget.tokens, widget.stage));
    super.initState();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: Text(domainLabel(widget.stage.domain))),
      body: FutureBuilder<SelectQuestionsOut>(
          future: loader,
          builder: (context, snapshot) {
            final data = snapshot.data;
            return snapshot.error != null
                ? ErrorCard(
                    "Impossible de charger les questions.", snapshot.error)
                : data == null
                    ? const Center(child: CircularProgressIndicator())
                    : CeinturesQuestionsW(
                        SeanceController(data.questions), _evaluate, _reset, (
                        api: widget.api,
                        tokens: widget.tokens,
                        stage: widget.stage
                      ));
          }),
    );
  }

  void _reset() async {
    setState(() {
      loader = widget.api
          .selectQuestions(SelectQuestionsIn(widget.tokens, widget.stage));
    });
  }

  Future<List<QuestionAnswersOut>> _evaluate(SeanceAnswers answers) async {
    final res = await widget.api.evaluateAnswers(EvaluateAnswersIn(
        widget.tokens, widget.stage, answers.questions, answers.answers));

    // notify the parent
    widget.onValid(res.success, res.evolution);

    return res.answers;
  }
}

class SeanceController {
  final List<InstantiatedBeltQuestion> questions;
  bool showCorrection;

  var _state = _State.answering;
  final PageController _pageC;

  final List<QuestionController> _controllers;

  SeanceController(this.questions,
      {int initialPage = 0, this.showCorrection = false})
      : _pageC = PageController(initialPage: initialPage),
        _controllers = questions.map(_buildQuestionController).toList();

  /// currentQuestion returns the question currently visible.
  int get currentQuestion => _pageC.hasClients ? _pageC.page?.round() ?? 0 : 0;

  /// [answers] returns the current answers for all the questions.
  SeanceAnswers answers() => SeanceAnswers(
      questions.map((qu) => qu.id).toList(),
      List.generate(
          questions.length,
          (index) =>
              AnswerP(questions[index].params, _controllers[index].answers())));

  /// [notAnswered] returns the indices of the question not done yet.
  List<int> notAnswered() {
    final out = <int>[];
    for (var i = 0; i < _controllers.length; i++) {
      if (!_controllers[i].buttonEnabled) {
        out.add(i);
      }
    }
    return out;
  }

  /// [setFeedback] displays the given feedback
  void setFeedback(List<QuestionAnswersOut> res) {
    _state = _State.displayingFeedback;
    for (var i = 0; i < res.length; i++) {
      final answer = res[i];
      _controllers[i].setFeedback(answer.results);
      _controllers[i].buttonEnabled = true;
      _controllers[i].buttonLabel = "Essayer à nouveau...";
    }
  }

  // only makes sense when displaying feedback
  bool hasError(int index) {
    return _controllers[index].feedback().values.any((b) => b);
  }

  /// setQuestionAnswers show the answers for the current question
  void setQuestionAnswers(Answers answers) {
    _state = _State.answering;
    _controllers[currentQuestion].setAnswers(answers);
    _controllers[currentQuestion].buttonEnabled = true;
    _controllers[currentQuestion].buttonLabel = "Valider";
  }
}

typedef TrainingMeta = ({
  CeinturesTrainingAPI api,
  StudentTokens tokens,
  Stage stage
});

class CeinturesQuestionsW extends StatefulWidget {
  final SeanceController controller;
  final Future<List<QuestionAnswersOut>> Function(SeanceAnswers)
      evaluateAnswers;

  final void Function() onReset;

  final TrainingMeta? trainingMeta;

  const CeinturesQuestionsW(
      this.controller, this.evaluateAnswers, this.onReset, this.trainingMeta,
      {super.key});

  @override
  State<CeinturesQuestionsW> createState() => _CeinturesQuestionsWState();
}

enum _State { answering, displayingFeedback }

class _CeinturesQuestionsWState extends State<CeinturesQuestionsW> {
  @override
  void initState() {
    if (widget.controller.showCorrection) {
      WidgetsBinding.instance.addPostFrameCallback((_) => _showCorrection());
    }
    super.initState();
  }

  @override
  void didUpdateWidget(covariant CeinturesQuestionsW oldWidget) {
    final c = widget.controller;

    WidgetsBinding.instance.addPostFrameCallback((_) {
      c._pageC.animateToPage(c._pageC.initialPage,
          duration: const Duration(milliseconds: 750), curve: Curves.easeInOut);
    });

    if (c.showCorrection) {
      WidgetsBinding.instance.addPostFrameCallback((_) => _showCorrection());
    }

    super.didUpdateWidget(oldWidget);
  }

  @override
  Widget build(BuildContext context) {
    final c = widget.controller;
    return PageView(
        controller: c._pageC,
        children: List.generate(
            c._controllers.length,
            (index) => QuestionView(
                  c.questions[index].question,
                  c._controllers[index],
                  _onValidQuestion,
                  Colors.teal,
                  title: "Question ${index + 1}/${c._controllers.length}",
                  leadingButton:
                      c._state == _State.displayingFeedback && c.hasError(index)
                          ? Padding(
                              padding: const EdgeInsets.only(right: 8.0),
                              child: OutlinedButton(
                                onPressed: () => _startTraining(index),
                                style: OutlinedButton.styleFrom(
                                    foregroundColor: Colors.yellow),
                                child: const Text("S'entrainer"),
                              ),
                            )
                          : null,
                )).toList());
  }

  void _showCorrection() {
    final c = widget.controller;
    Navigator.of(context).push(MaterialPageRoute<void>(
      builder: (context) => Scaffold(
        appBar: AppBar(),
        body: CorrectionView(c.questions[c.currentQuestion].question.correction,
            Colors.greenAccent, pickQuote()),
      ),
    ));
  }

  void _onValidQuestion() {
    final c = widget.controller;
    if (c._state == _State.displayingFeedback) {
      // reset
      widget.onReset();
      return;
    }

    // find the next not answered question
    final notAnswered = c.notAnswered();
    if (notAnswered.isEmpty) {
      _submitAnswers();
    } else {
      final goTo = notAnswered.first;
      c._pageC.animateToPage(goTo,
          duration: const Duration(milliseconds: 750), curve: Curves.easeInOut);
      // update button label for the last question
      if (notAnswered.length == 1) {
        setState(() {
          c._controllers[goTo].buttonLabel = "Soumettre !";
        });
      }
    }
  }

  void _submitAnswers() async {
    final List<QuestionAnswersOut> res;
    try {
      res = await widget.evaluateAnswers(widget.controller.answers());
    } catch (e) {
      if (!mounted) return;
      showError("Impossible d'évaluer les réponses.", e, context);
      return;
    }
    if (!mounted) return;

    if (res.every((element) => element.isCorrect)) {
      return;
    }

    // always display the feedback and show a summary
    setState(() {
      widget.controller.setFeedback(res);
    });

    // display a recap of errors
    final goTo =
        await Navigator.of(context).push(MaterialPageRoute<_ResultAction>(
      builder: (context) => _ResultsPage(res),
    ));
    if (goTo == null) return;

    // if asked, go to a question
    if (goTo.startTraining) {
      _startTraining(goTo.index);
    } else {
      setState(() {
        widget.controller._pageC.animateToPage(goTo.index,
            duration: const Duration(milliseconds: 750),
            curve: Curves.easeInOut);
      });
    }
  }

  void _startTraining(int questionIndex) {
    if (widget.trainingMeta == null) return;

    Navigator.of(context).push(MaterialPageRoute<void>(
      builder: (context) => _TrainingView(
          widget.trainingMeta!, widget.controller.questions[questionIndex].id),
    ));
  }
}

QuestionController _buildQuestionController(InstantiatedBeltQuestion question) {
  final out = QuestionController.fromQuestion(question.question);
  out.buttonLabel = "Continuer";
  return out;
}

extension on EvaluateAnswersOut {
  bool get success => answers.every((element) => element.isCorrect);
}

typedef _ResultAction = ({int index, bool startTraining});

// displays a summary of errors and
// propose the training mode
class _ResultsPage extends StatelessWidget {
  final List<QuestionAnswersOut> resultats;
  const _ResultsPage(this.resultats);

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text("Résultats"),
      ),
      body: Padding(
        padding: const EdgeInsets.all(8),
        child: Column(
          children: [
            const SizedBox(height: 12),
            const ColoredTitle("Voici ton bilan", Colors.orange),
            const SizedBox(height: 12),
            Expanded(
              child: ListView(
                  children: List.generate(resultats.length, (index) {
                final ok = resultats[index].isCorrect;
                return ListTile(
                  shape: RoundedRectangleBorder(
                      borderRadius: BorderRadius.circular(8)),
                  title: Text("Question ${index + 1}"),
                  leading: Icon(
                    ok ? Icons.check : Icons.clear,
                    color: ok ? Colors.green : Colors.red,
                    size: 32,
                  ),
                  trailing: ok
                      ? null
                      : OutlinedButton(
                          style: OutlinedButton.styleFrom(
                              foregroundColor: Colors.yellow),
                          child: const Text("S'entrainer"),
                          onPressed: () => Navigator.of(context)
                              .pop<_ResultAction>(
                                  (index: index, startTraining: true)),
                        ),
                  onTap: () => Navigator.of(context)
                      .pop<_ResultAction>((index: index, startTraining: false)),
                );
              })),
            ),
            Quote(pickQuote()),
          ],
        ),
      ),
    );
  }
}

class _TrainingView extends StatefulWidget {
  final TrainingMeta trainingMeta;
  final IdBeltquestion idQuestion;

  const _TrainingView(this.trainingMeta, this.idQuestion);

  @override
  State<_TrainingView> createState() => __TrainingViewState();
}

class __TrainingViewState extends State<_TrainingView> {
  InstantiatedBeltQuestion? question;
  QuestionController? controller;
  _State state = _State.answering;

  @override
  void didUpdateWidget(covariant _TrainingView oldWidget) {
    _loadQuestion();
    super.didUpdateWidget(oldWidget);
  }

  @override
  void initState() {
    _loadQuestion();
    super.initState();
  }

  void _loadQuestion() async {
    final res = await widget.trainingMeta.api.instantiateTraining(
        InstantiateTrainingQuestionIn(
            widget.trainingMeta.tokens, widget.idQuestion));
    setState(() {
      state = _State.answering;
      question = res;
      controller = QuestionController.fromQuestion(res.question);
    });
  }

  @override
  Widget build(BuildContext context) {
    final question = this.question;
    return Scaffold(
      appBar: AppBar(
        title: const Text("Entraînement"),
      ),
      body: question == null
          ? CircularProgressIndicator()
          : QuestionView(
              question.question, controller!, _onButtonClick, Colors.yellow),
    );
  }

  void _onButtonClick() async {
    switch (state) {
      case _State.answering:
        _evaluate();
      case _State.displayingFeedback:
        _loadQuestion();
    }
  }

  void _evaluate() async {
    if (question == null || controller == null) return;

    final result = await widget.trainingMeta.api.evaluateTraining(
        EvaluateAnswerTrainingIn(
            widget.trainingMeta.tokens,
            widget.trainingMeta.stage,
            widget.idQuestion,
            AnswerP(question!.params, controller!.answers())));
    if (!mounted) return;

    setState(() {
      state = _State.displayingFeedback;
      controller?.buttonLabel = "Recommencer";
      controller?.buttonEnabled = true;
      controller?.setFeedback(result.isCorrect ? null : result.results);
    });

    if (result.isCorrect) {
      ScaffoldMessenger.of(context).showSnackBar(SnackBar(
        content: const Text("Bonne réponse !"),
        backgroundColor: Colors.lightGreen,
      ));
    }
  }
}
