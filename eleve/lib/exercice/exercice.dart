import 'package:eleve/build_mode.dart';
import 'package:eleve/exercice/congratulations.dart';
import 'package:eleve/exercice/home.dart';
import 'package:eleve/questions/question.dart';
import 'package:eleve/questions/types.gen.dart' hide Answer;
import 'package:eleve/shared_gen.dart';
import 'package:flutter/material.dart' hide Flow;

class NotificationExerciceDone extends Notification {}

/// ExerciceW is the widget providing one exercice to
/// the student.
/// It is used in the editor loopback, and as the base for
/// at home training activity
class ExerciceW extends StatefulWidget {
  final BuildMode buildMode;

  /// [data] stores the server instantiated exercice with
  /// the initial progression state.
  final Exercice data;

  final Future<EvaluateExerciceOut> Function(EvaluateExerciceIn) onEvaluate;

  const ExerciceW(this.buildMode, this.data, this.onEvaluate, {Key? key})
      : super(key: key);

  @override
  State<ExerciceW> createState() => _ExerciceWState();
}

class _ExerciceWState extends State<ExerciceW> {
  late List<InstantiatedQuestion> questions; // will change upon wrong answer
  late ProgressionExt progression; // will change on validate

  Map<int, QuestionAnswersIn> currentAnswers = {};

  int? questionIndex; // null means summary
  Map<int, bool>? results;

  @override
  void initState() {
    questions = widget.data.exercice.questions;
    progression = widget.data.progression;
    super.initState();
  }

  EvaluateExerciceIn get currentEvaluate => EvaluateExerciceIn(
      widget.data.exercice.id,
      currentAnswers.map((index, value) =>
          MapEntry(index, Answer(questions[index].params, value))),
      progression);

  void onValideQuestion(QuestionAnswersIn answer) async {
    final index = questionIndex!;
    currentAnswers[index] = answer;
    switch (widget.data.exercice.flow) {
      case Flow.sequencial:
        // validate the given answer
        final resp = await widget.onEvaluate(currentEvaluate);
        progression = resp.progression; // update the progression
        final isCorrect =
            resp.results[index]!.results.values.every((success) => success);
        ScaffoldMessenger.of(context).showSnackBar(SnackBar(
          backgroundColor: isCorrect ? Colors.lightGreen : Colors.red.shade200,
          duration: Duration(seconds: isCorrect ? 2 : 4),
          content: Text(isCorrect ? "Bonne réponse" : "Réponse incorrecte"),
        ));

        if (progression.nextQuestion == -1) {
          // exercice is over
          await showDialog<void>(
              context: context,
              builder: (context) => const Dialog(child: Congrats()));
          NotificationExerciceDone().dispatch(context);
          return;
        }

        if (isCorrect) {
          // go to next question
          setState(() {
            results = null;
            questionIndex = progression.nextQuestion;
          });
        } else {
          // show errors and ask for retry
          setState(() {
            results = resp.results[index]!.results;
          });
        }
        break;
      case Flow.parallel:
      // only validate if all the questions have been completed

    }
  }

  bool get goToPreviousEnabled => questionIndex != null;
  bool get goToNextEnabled => questionIndex != questions.length - 1;

  void goToPrevious() {
    final newIndex = questionIndex == 0 ? null : questionIndex! - 1;
    setState(() {
      questionIndex = newIndex;
    });
  }

  void goToNext() {
    final newIndex = questionIndex == null ? 0 : questionIndex! + 1;
    setState(() {
      questionIndex = newIndex;
    });
  }

  @override
  Widget build(BuildContext context) {
    final ex = widget.data.exercice;
    return Scaffold(
        appBar: AppBar(
          title: const Text("Exercice"),
          actions: [
            IconButton(
                onPressed: goToPreviousEnabled ? goToPrevious : null,
                icon: Icon(IconData(0xf572,
                    fontFamily: 'MaterialIcons', matchTextDirection: true))),
            IconButton(
                onPressed: goToNextEnabled ? goToNext : null,
                icon: Icon(IconData(0xf57a,
                    fontFamily: 'MaterialIcons', matchTextDirection: true))),
          ],
        ),
        body: questionIndex == null
            ? ExerciceHome(
                Exercice(
                    InstantiatedExercice(
                        ex.id, ex.title, ex.flow, questions, ex.baremes),
                    progression),
                (index) => setState(() {
                      questionIndex = index;
                    }))
            : QuestionW(
                widget.buildMode,
                questions[questionIndex!].question,
                Colors.purpleAccent,
                onValideQuestion,
                title: "Question ${questionIndex! + 1}",
                timeout: null,
                blockOnSubmit: true,
                results: results,
                answer: currentAnswers[questionIndex!]?.data,
              ));
  }
}
