import 'package:eleve/build_mode.dart';
import 'package:eleve/questions/question.dart';
import 'package:eleve/questions/types.gen.dart';
import 'package:eleve/quotes.dart';
import 'package:flutter/material.dart';

import 'categories.dart';
import 'events.gen.dart';

class QuestionRoute extends StatelessWidget {
  final BuildMode buildMode;
  final ShowQuestion question;
  final void Function(ValidQuestionNotification) onValid;

  /// if [readonly] is true, do no block when leaving the question,
  /// and do not display the timer
  final bool readonly;
  final Answers? answer;

  const QuestionRoute(this.buildMode, this.question, this.onValid,
      {Key? key, this.readonly = false, this.answer})
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
    // make the route block until validated
    return WillPopScope(
      onWillPop: () async {
        if (readonly) {
          return true;
        }

        final cancel = await _confirmCancel(context);
        if (cancel) {
          // send an empty response
          onValid(ValidQuestionNotification(const QuestionAnswersIn({})));
        }
        return cancel;
      },
      child: Scaffold(
        resizeToAvoidBottomInset: false,
        appBar: AppBar(
          automaticallyImplyLeading: true,
        ),
        body: QuestionW(
          buildMode,
          question.question,
          question.categorie.color,
          onValid,
          timeout: readonly ? null : Duration(seconds: question.timeoutSeconds),
          footerQuote: pickQuote(),
          answer: answer,
        ),
      ),
    );
  }
}
