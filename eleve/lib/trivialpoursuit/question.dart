import 'package:eleve/build_mode.dart';
import 'package:eleve/exercices/question.dart';
import 'package:eleve/quotes.dart';
import 'package:flutter/material.dart';

import 'categories.dart';
import 'events.gen.dart';

class QuestionRoute extends StatelessWidget {
  final BuildMode buildMode;
  final ShowQuestion question;
  final void Function(ValidQuestionNotification) onValid;

  const QuestionRoute(this.buildMode, this.question, this.onValid, {Key? key})
      : super(key: key);

  @override
  Widget build(BuildContext context) {
    // make the route block until validated
    return WillPopScope(
      onWillPop: () async => false,
      child: Scaffold(
        appBar: AppBar(
          automaticallyImplyLeading: false,
        ),
        body: QuestionW(
          buildMode,
          question.question,
          question.categorie.color,
          onValid,
          timeout: Duration(seconds: question.timeoutSeconds),
          footerQuote: pickQuote(),
        ),
      ),
    );
  }
}
