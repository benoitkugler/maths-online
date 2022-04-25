import 'package:eleve/exercices/question.dart';
import 'package:eleve/quotes.dart';
import 'package:flutter/material.dart';

import 'categories.dart';
import 'events.gen.dart';

class QuestionRoute extends StatelessWidget {
  final ShowQuestion question;
  final void Function(ValidQuestionNotification) onValid;
  final void Function(CheckQuestionSyntaxeNotification) onCheckSyntax;

  const QuestionRoute(this.question, this.onValid, this.onCheckSyntax,
      {Key? key})
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
          question.question,
          question.categorie.color,
          onCheckSyntax,
          onValid,
          footerQuote: pickQuote(),
          timeout: Duration(seconds: question.timeoutSeconds),
        ),
      ),
    );
  }
}
