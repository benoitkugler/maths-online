import 'package:eleve/activities/automatismes/automatismes.dart';
import 'package:eleve/activities/homework/exercice.dart';
import 'package:eleve/questions/question.dart';
import 'package:eleve/quotes.dart';
import 'package:eleve/shared/errors.dart';
import 'package:eleve/types/src_automatismes.dart';
import 'package:eleve/types/src_maths_questions_client.dart';
import 'package:eleve/types/src_tasks.dart';
import 'package:flutter/material.dart';

class AutomatismeQuestionW extends StatefulWidget {
  final InstantiatedAutomatismeQuestion question;
  final AutomatismesAPI api;

  final void Function() startAgain;

  const AutomatismeQuestionW(
    this.question,
    this.api,
    this.startAgain, {
    super.key,
  });

  @override
  State<AutomatismeQuestionW> createState() => _AutomatismeQuestionWState();
}

enum _State { answering, displayingFeedback }

class _AutomatismeQuestionWState extends State<AutomatismeQuestionW> {
  _State state = .answering;
  late QuestionController ct;

  @override
  void initState() {
    ct = _buildQuestionController(widget.question);
    super.initState();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(),
      body: QuestionView(
        widget.question.question,
        ct,
        _onValidQuestion,
        Colors.teal,
      ),
    );
  }

  void _showCorrection() {
    Navigator.of(context).push(
      MaterialPageRoute<void>(
        builder: (context) => Scaffold(
          appBar: AppBar(),
          body: CorrectionView(
            widget.question.question.correction,
            Colors.greenAccent,
            pickQuote(),
          ),
        ),
      ),
    );
  }

  void _onValidQuestion() {
    if (state == .displayingFeedback) {
      // reset
      widget.startAgain();
      return;
    }

    _submitAnswers();
  }

  void _submitAnswers() async {
    final QuestionAnswersOut res;
    try {
      res = await widget.api.evaluateQuestion(
        EvaluateAutomatismeIn(
          widget.question.id,
          AnswerP(widget.question.params, ct.answers()),
        ),
      );
    } catch (e) {
      if (!mounted) return;
      showError("Impossible d'évaluer la réponse.", e, context);
      return;
    }
    if (!mounted) return;

    setState(() {
      state = _State.displayingFeedback;
      ct.buttonLabel = "Recommencer";
      ct.buttonEnabled = true;
      ct.setFeedback(res.isCorrect ? null : res.results);
    });

    if (res.isCorrect) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(
          content: const Text("Bonne réponse !"),
          backgroundColor: Colors.lightGreen,
        ),
      );
    } else {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(
          content: const Text("Réponse incorrecte."),
          backgroundColor: Colors.red.shade300,
        ),
      );
    }
  }
}

QuestionController _buildQuestionController(
  InstantiatedAutomatismeQuestion question,
) {
  final out = QuestionController.fromQuestion(question.question);
  out.buttonLabel = "Valider";
  return out;
}
