import 'package:eleve/activities/ceintures/seance.dart';
import 'package:eleve/loopback/loopback.dart';
import 'package:eleve/types/src_maths_questions_client.dart';
import 'package:eleve/types/src_prof_preview.dart';
import 'package:flutter/material.dart';

class LoopbackCeinturesController implements LoopbackController {
  final LoopbackShowCeinture data;
  final SeanceController controller;

  LoopbackCeinturesController(this.data)
      : controller = SeanceController(data.questions,
            initialPage: data.questionIndex,
            showCorrection: data.showCorrection);
}

class CeintureW extends StatelessWidget {
  final LoopbackCeinturesController controller;

  final Future<List<QuestionAnswersOut>> Function(SeanceAnswers) evaluate;
  final void Function() reset;
  final void Function() loadCorrectAnswer;

  const CeintureW(
      this.controller, this.evaluate, this.reset, this.loadCorrectAnswer,
      {super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text("Ceinture"), actions: [
        TextButton(
            onPressed: loadCorrectAnswer,
            child: const Text("Afficher la réponse"))
      ]),
      body: CeinturesQuestionsW(
        controller.controller,
        evaluate,
        reset,
        null,
      ),
    );
  }
}
