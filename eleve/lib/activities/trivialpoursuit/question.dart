import 'package:eleve/questions/question.dart';
import 'package:eleve/quotes.dart';
import 'package:eleve/types/src_maths_questions_client.dart';
import 'package:eleve/types/src_trivial.dart' hide Answer;
import 'package:flutter/material.dart';

import 'categories.dart';

class InGameQuestionRoute extends StatefulWidget {
  final ShowQuestion question;
  final void Function(QuestionAnswersIn) onValid;

  const InGameQuestionRoute(this.question, this.onValid, {super.key});

  @override
  State<InGameQuestionRoute> createState() => _InGameQuestionRouteState();
}

class _InGameQuestionRouteState extends State<InGameQuestionRoute> {
  late final QuestionController controller;

  @override
  void initState() {
    controller = QuestionController.fromQuestion(widget.question.question);
    controller.timeout = Duration(
        seconds: widget.question.timeoutSeconds); // show the timeout bar
    super.initState();
  }

  @override
  void didUpdateWidget(covariant InGameQuestionRoute oldWidget) {
    controller = QuestionController.fromQuestion(widget.question.question);
    controller.timeout = Duration(
        seconds: widget.question.timeoutSeconds); // show the timeout bar
    super.didUpdateWidget(oldWidget);
  }

  Future<bool> _confirmCancel(BuildContext context) async {
    final res = await showDialog<bool?>(
        context: context,
        builder: (BuildContext context) {
          return AlertDialog(
            title: const Text("Abandonner la question"),
            content: const Text("Es-tu sûr d'abandonner la question ?"),
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

  void onValid() {
    setState(() {
      controller.buttonEnabled = false;
      controller.buttonLabel = "En attente des autres joueurs...";
      controller.footerQuote = pickQuote();
      controller.setFieldsEnabled(false);
    });
    // propagate the event
    widget.onValid(controller.answers());
  }

  @override
  Widget build(BuildContext context) {
    // make the route block until validated
    return WillPopScope(
      onWillPop: () async {
        final cancel = await _confirmCancel(context);
        if (cancel) {
          // send an empty response
          widget.onValid(const QuestionAnswersIn({}));
        }
        return cancel;
      },
      child: Scaffold(
        resizeToAvoidBottomInset: false,
        appBar: AppBar(
          automaticallyImplyLeading: true,
        ),
        body: QuestionView(
          widget.question.question,
          controller,
          onValid,
          widget.question.categorie.color,
        ),
      ),
    );
  }
}

class LastQuestionRoute extends StatelessWidget {
  final ShowQuestion question;
  final void Function() onDone;

  final Answers answers;

  const LastQuestionRoute(this.question, this.onDone, this.answers,
      {super.key});

  @override
  Widget build(BuildContext context) {
    final ct = QuestionController.fromQuestion(question.question);
    ct.setAnswers(answers);
    ct.setFieldsEnabled(false);
    ct.buttonEnabled = true;
    ct.buttonLabel = "Retour";

    return Scaffold(
      resizeToAvoidBottomInset: false,
      appBar: AppBar(
        automaticallyImplyLeading: true,
      ),
      body: QuestionView(
        question.question,
        ct,
        onDone,
        question.categorie.color,
        title: "Dernière question",
      ),
    );
  }
}
