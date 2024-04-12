import 'package:eleve/questions/question.dart';
import 'package:eleve/quotes.dart';
import 'package:eleve/types/src_maths_questions_client.dart';
import 'package:eleve/types/src_trivial.dart' hide Answer;
import 'package:flutter/material.dart';

import 'categories.dart';

class InGameQuestionController extends BaseQuestionController {
  final void Function(QuestionAnswersIn) onValid;

  InGameQuestionController(ShowQuestion question, this.onValid)
      : super(question.question) {
    state.timeout =
        Duration(seconds: question.timeoutSeconds); // show the timeout bar
    state.footerQuote = null;
  }

  @override
  void onPrimaryButtonClick() {
    state.buttonEnabled = false;
    state.buttonLabel = "En attente des autres joueurs...";
    state.footerQuote = pickQuote();
    disableAllFields();
    // propagate the event
    onValid(answers());
  }
}

class LastQuestionController extends BaseQuestionController {
  final void Function() onClose;
  LastQuestionController(
      Question question, Map<int, Answer> answers, this.onClose)
      : super(question) {
    setAnswers(answers);
    disableAllFields();
    state.buttonEnabled = true;
    state.buttonLabel = "Retour";
  }

  @override
  void onPrimaryButtonClick() {
    // simple go back
    onClose();
  }
}

class InGameQuestionRoute extends StatelessWidget {
  final ShowQuestion question;
  final void Function(QuestionAnswersIn) onValid;

  const InGameQuestionRoute(this.question, this.onValid, {Key? key})
      : super(key: key);

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

  @override
  Widget build(BuildContext context) {
    final ct = InGameQuestionController(question, onValid);

    // make the route block until validated
    return WillPopScope(
      onWillPop: () async {
        final cancel = await _confirmCancel(context);
        if (cancel) {
          // send an empty response
          onValid(const QuestionAnswersIn({}));
        }
        return cancel;
      },
      child: Scaffold(
        resizeToAvoidBottomInset: false,
        appBar: AppBar(
          automaticallyImplyLeading: true,
        ),
        body: QuestionW(
          ct,
          question.categorie.color,
        ),
      ),
    );
  }
}

class LastQuestionRoute extends StatelessWidget {
  final ShowQuestion question;
  final void Function() onDone;

  final Answers answers;

  const LastQuestionRoute(this.question, this.onDone, this.answers, {Key? key})
      : super(key: key);

  @override
  Widget build(BuildContext context) {
    final ct = LastQuestionController(question.question, answers, onDone);

    return Scaffold(
      resizeToAvoidBottomInset: false,
      appBar: AppBar(
        automaticallyImplyLeading: true,
      ),
      body: QuestionW(
        ct,
        question.categorie.color,
        title: "Dernière question",
      ),
    );
  }
}
