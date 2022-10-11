import 'dart:math';

import 'package:eleve/questions/fields.dart';
import 'package:eleve/questions/question.dart';
import 'package:eleve/questions/types.gen.dart';
import 'package:flutter/material.dart';

class LoopackQuestionController extends BaseQuestionController {
  final void Function(QuestionAnswersIn) onValid;

  LoopackQuestionController(Question question, FieldAPI api, this.onValid)
      : super(question, api);

  @override
  void onPrimaryButtonClick() {
    state.buttonEnabled = false;
    state.buttonLabel = "Correction...";
    onValid(answers());
  }

  @override
  void setAnswers(Map<int, Answer> answers) {
    super.setAnswers(answers); // this trigger onFieldChange
    state.buttonEnabled = true;
    state.buttonLabel = "Valider";
  }

  @override
  void setFeedback(Map<int, bool>? feedback) {
    super.setFeedback(feedback);
    state.buttonEnabled = true;
    state.buttonLabel = "Valider";
    setFieldsEnabled(true);
  }
}

class LoopbackQuestionW extends StatefulWidget {
  final LoopackQuestionController controller;

  final void Function() loadCorrectAnswer;

  const LoopbackQuestionW(this.controller, this.loadCorrectAnswer, {super.key});

  @override
  State<LoopbackQuestionW> createState() => _LoopbackQuestionWState();

  static void showServerValidation(
      QuestionAnswersOut rep, BuildContext context) {
    final crible = rep.results;
    final errors = crible.values.where((value) => !value).toList();
    final isValid = errors.isEmpty;
    final errorMessage = errors.length >= 2
        ? "${errors.length} champs sont incorrects"
        : "Un champ est incorrect";
    ScaffoldMessenger.of(context).showSnackBar(SnackBar(
      backgroundColor: isValid ? Colors.lightGreen : Colors.red,
      content: Text(isValid ? "Bonne réponse" : errorMessage),
    ));
  }
}

class _LoopbackQuestionWState extends State<LoopbackQuestionW> {
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
          child: QuestionW(
            widget.controller,
            Color.fromARGB(255, 150 + Random().nextInt(100),
                150 + Random().nextInt(100), Random().nextInt(256)),
          ),
        ));
  }
}
