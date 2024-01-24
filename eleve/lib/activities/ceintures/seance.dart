import 'package:eleve/activities/ceintures/api.dart';
import 'package:eleve/exercice/home.dart';
import 'package:eleve/questions/question.dart';
import 'package:eleve/shared/errors.dart';
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
                        SeanceController(data.questions),
                        _evaluate,
                        _reset,
                      );
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

  var _state = _State.answering;
  final _pageC = PageController();

  final List<_QuestionController> _controllers;

  SeanceController(this.questions)
      : _controllers = List.generate(questions.length,
            (index) => _QuestionController(questions[index].question));

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
      if (!_controllers[i].state.buttonEnabled) {
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
      _controllers[i].state.buttonEnabled = true;
      _controllers[i].state.buttonLabel = "Essayer à nouveau...";
    }
  }

  /// setQuestionAnswers show the answers for the current question
  void setQuestionAnswers(Answers answers) {
    _state = _State.answering;
    _controllers[currentQuestion].setAnswers(answers);
    _controllers[currentQuestion].state.buttonEnabled = true;
    _controllers[currentQuestion].state.buttonLabel = "Valider";
  }
}

class CeinturesQuestionsW extends StatefulWidget {
  final SeanceController controller;
  final Future<List<QuestionAnswersOut>> Function(SeanceAnswers)
      evaluateAnswers;

  final void Function() onReset;

  const CeinturesQuestionsW(this.controller, this.evaluateAnswers, this.onReset,
      {super.key});

  @override
  State<CeinturesQuestionsW> createState() => _CeinturesQuestionsWState();
}

enum _State { answering, displayingFeedback }

class _CeinturesQuestionsWState extends State<CeinturesQuestionsW> {
  @override
  Widget build(BuildContext context) {
    final c = widget.controller;
    return PageView(
        controller: c._pageC,
        children: List.generate(
            c._controllers.length,
            (index) => QuestionW(
                  c._controllers[index],
                  Colors.teal,
                  title: "Question ${index + 1}/${c._controllers.length}",
                  onButtonClick: _onValidQuestion,
                )).toList());
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
          c._controllers[goTo].state.buttonLabel = "Soumettre !";
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

    final goTo = res.indexWhere((r) => !r.isCorrect);

    // display feedback and go to the "first" error
    setState(() {
      widget.controller.setFeedback(res);
      widget.controller._pageC.animateToPage(goTo,
          duration: const Duration(milliseconds: 750), curve: Curves.easeInOut);
    });
    ScaffoldMessenger.of(context).showSnackBar(SnackBar(
        backgroundColor: Colors.red.shade400,
        content: const Text("Dommage, il y a des réponses incorrectes...")));
  }
}

class _QuestionController extends BaseQuestionController {
  _QuestionController(super.question) {
    state.buttonLabel = "Continuer";
  }

  @override
  void onPrimaryButtonClick() {}
}

extension on EvaluateAnswersOut {
  bool get success => answers.every((element) => element.isCorrect);
}
