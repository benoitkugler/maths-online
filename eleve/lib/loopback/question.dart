import 'dart:math';

import 'package:eleve/loopback/loopback.dart';
import 'package:eleve/questions/question.dart';
import 'package:eleve/types/src_maths_questions_client.dart';
import 'package:eleve/types/src_prof_preview.dart';
import 'package:flutter/material.dart';

class LoopackQuestionController implements LoopbackController {
  final LoopbackShowQuestion data;
  final QuestionController controller;

  LoopackQuestionController(this.data)
      : controller = QuestionController.fromQuestion(data.question);

  void setAnswers(QuestionAnswers answers) {
    controller.setAnswers(answers); // this trigger onFieldChange
    controller.buttonEnabled = true;
    controller.buttonLabel = "Valider";
  }

  void setFeedback(QuestionFeedback feedback) {
    controller.setFeedback(feedback);
    controller.buttonEnabled = true;
    controller.buttonLabel = "Valider";
    controller.setFieldsEnabled(true); // reactivate button for convenience
  }
}

class LoopbackQuestionW extends StatefulWidget {
  final LoopackQuestionController controller;

  final void Function(QuestionAnswersIn) onValid;
  final void Function() loadCorrectAnswer;

  const LoopbackQuestionW(this.controller, this.onValid, this.loadCorrectAnswer,
      {super.key});

  @override
  State<LoopbackQuestionW> createState() => _LoopbackQuestionWState();

  static SnackBar serverValidation(
      QuestionAnswersOut rep, void Function() onShowCorrection) {
    final crible = rep.results;
    final errors = crible.values.where((value) => !value).toList();
    final isValid = errors.isEmpty;
    final errorMessage = errors.length >= 2
        ? "${errors.length} champs sont incorrects."
        : "Un champ est incorrect.";
    return SnackBar(
      backgroundColor: isValid ? Colors.lightGreen : Colors.red.shade400,
      content: Text(isValid ? "Bonne réponse" : errorMessage),
      action: SnackBarAction(
          textColor: isValid ? Colors.black : Colors.white,
          label: "Afficher la correction",
          onPressed: onShowCorrection),
    );
  }
}

class _LoopbackQuestionWState extends State<LoopbackQuestionW> {
  void onValid() {
    setState(() {
      widget.controller.controller.buttonEnabled = false;
      widget.controller.controller.buttonLabel = "Correction...";
    });
    widget.onValid(widget.controller.controller.answers());
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
        appBar: AppBar(actions: [
          TextButton(
              onPressed: widget.loadCorrectAnswer,
              child: const Text("Afficher la réponse"))
        ]),
        body: Padding(
          padding: const EdgeInsets.symmetric(vertical: 8.0),
          child: QuestionView(
            widget.controller.data.question,
            widget.controller.controller,
            onValid,
            Colors.yellow.shade700,
          ),
        ));
  }
}

Color randColor() {
  return Color.fromARGB(255, 150 + Random().nextInt(100),
      150 + Random().nextInt(100), Random().nextInt(256));
}
