import 'package:eleve/activities/ceintures/api.dart';
import 'package:eleve/exercice/home.dart';
import 'package:eleve/questions/question.dart';
import 'package:eleve/shared/errors.dart';
import 'package:eleve/types/src_prof_ceintures.dart';
import 'package:eleve/types/src_sql_ceintures.dart';
import 'package:eleve/types/src_tasks.dart';
import 'package:flutter/material.dart';

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
                    : _QuestionsW(widget.api, widget.tokens, widget.stage, data,
                        widget.onValid, reset);
          }),
    );
  }

  void reset() async {
    setState(() {
      loader = widget.api
          .selectQuestions(SelectQuestionsIn(widget.tokens, widget.stage));
    });
  }
}

class _QuestionsW extends StatefulWidget {
  final CeinturesAPI api;
  final StudentTokens tokens;
  final Stage stage;
  final SelectQuestionsOut questions;
  final void Function(bool isSuccess, StudentEvolution newEvolution) onValid;
  final void Function() onReset;
  const _QuestionsW(this.api, this.tokens, this.stage, this.questions,
      this.onValid, this.onReset,
      {super.key});

  @override
  State<_QuestionsW> createState() => __QuestionsWState();
}

enum _State { answering, displayingFeedback }

class __QuestionsWState extends State<_QuestionsW> {
  var state = _State.answering;
  var pageC = PageController();

  List<_QuestionController> controllers = [];

  @override
  void initState() {
    _initControllers();
    super.initState();
  }

  @override
  void didUpdateWidget(covariant _QuestionsW oldWidget) {
    _initControllers();
    pageC.jumpToPage(0);
    super.didUpdateWidget(oldWidget);
  }

  void _initControllers() {
    state = _State.answering;
    final l = widget.questions.questions;
    controllers = List.generate(
        l.length,
        (index) => _QuestionController(
            l[index].params, l[index].question, () => _onValidQuestion(index)));
  }

  @override
  Widget build(BuildContext context) {
    return PageView(
        controller: pageC,
        children: List.generate(
            controllers.length,
            (index) => QuestionW(
                  controllers[index],
                  Colors.teal,
                  title: "Question ${index + 1}/${controllers.length}",
                )).toList());
  }

  void _onValidQuestion(int index) {
    if (state == _State.displayingFeedback) {
      // reset
      widget.onReset();
      return;
    }

    // find the next not answered question
    final goTo = controllers.indexWhere((ct) => !ct.state.buttonEnabled);
    if (goTo == -1) {
      _submitAnswers();
    } else {
      pageC.animateToPage(goTo,
          duration: const Duration(milliseconds: 750), curve: Curves.easeInOut);
      // update button label for the last question
      if (controllers.where((ct) => !ct.state.buttonEnabled).length == 1) {
        setState(() {
          controllers[goTo].state.buttonLabel = "Soumettre !";
        });
      }
    }
  }

  void _submitAnswers() async {
    final EvaluateAnswersOut res;
    try {
      res = await widget.api.evaluateAnswers(EvaluateAnswersIn(
          widget.tokens,
          widget.stage,
          widget.questions.questions.map((qu) => qu.id).toList(),
          controllers.map((ct) => AnswerP(ct.params, ct.answers())).toList()));
    } catch (e) {
      if (!mounted) return;
      showError("Impossible d'évaluer les réponses.", e, context);
      return;
    }
    if (!mounted) return;

    if (res.success) {
      widget.onValid(true, res.evolution);
      return;
    }

    // notify the parent
    widget.onValid(false, res.evolution);

    // display feedback ...
    for (var i = 0; i < res.answers.length; i++) {
      final answer = res.answers[i];
      controllers[i].setFeedback(answer.results);
      controllers[i].state.buttonEnabled = true;
      controllers[i].state.buttonLabel = "Essayer à nouveau...";
    }
    // .. and go to the "first" error
    final goTo = res.answers.indexWhere((r) => !r.isCorrect);
    setState(() {
      state = _State.displayingFeedback;
      pageC.animateToPage(goTo,
          duration: const Duration(milliseconds: 750), curve: Curves.easeInOut);
    });
    ScaffoldMessenger.of(context).showSnackBar(SnackBar(
        backgroundColor: Colors.red.shade400,
        content: const Text("Dommage, il y a des réponses incorrectes...")));
  }
}

class _QuestionController extends BaseQuestionController {
  final Params params;
  void Function() onValid;

  _QuestionController(this.params, super.question, this.onValid) {
    state.buttonLabel = "Continuer";
  }

  @override
  void onPrimaryButtonClick() {
    onValid();
  }
}

extension on EvaluateAnswersOut {
  bool get success => answers.every((element) => element.isCorrect);
}
