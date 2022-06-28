import 'package:eleve/build_mode.dart';
import 'package:eleve/questions/question.dart';
import 'package:eleve/questions/types.gen.dart';
import 'package:flutter/material.dart';

class QuestionPage extends StatelessWidget {
  final BuildMode buildMode;
  final Question question;
  final void Function(ValidQuestionNotification) onValid;

  const QuestionPage(this.buildMode, this.question, this.onValid, {Key? key})
      : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(),
      body: QuestionW(
        buildMode,
        question,
        Colors.purpleAccent,
        onValid,
        timeout: null,
        blockOnSubmit: true,
      ),
    );
  }
}
