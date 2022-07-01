import 'package:eleve/build_mode.dart';
import 'package:eleve/exercice/congratulations.dart';
import 'package:eleve/exercice/home.dart';
import 'package:eleve/questions/question.dart';
import 'package:eleve/questions/types.gen.dart' hide Answer;
import 'package:eleve/shared_gen.dart';
import 'package:flutter/material.dart' hide Flow;

class NotificationExerciceDone extends Notification {}

extension IsCorrect on QuestionAnswersOut {
  bool get isCorrect {
    return results.values.every((success) => success);
  }
}

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

// the questions to display when trying again
  List<InstantiatedQuestion> nextQuestions = [];

  Map<int, QuestionAnswersIn> currentAnswers = {};
  Map<int, Map<int, bool>> results = {};

  int? questionIndex; // null means summary

  @override
  void initState() {
    questions = widget.data.exercice.questions;
    progression = widget.data.progression;
    super.initState();
  }

  void onDone() async {
    setState(() {
      questionIndex = null;
    });

    // exercice is over
    await showDialog<void>(
        context: context,
        builder: (context) => const Dialog(child: Congrats()));
    NotificationExerciceDone().dispatch(context);
  }

  void onValidQuestionSequential() async {
    final index = questionIndex!;

    // if we are not at the current question, just go to it
    if (index != progression.nextQuestion) {
      setState(() {
        questionIndex = progression.nextQuestion;
      });
      return;
    }

    // validate the given answer
    final resp = await widget.onEvaluate(EvaluateExerciceIn(
        widget.data.exercice.id,
        {index: Answer(questions[index].params, currentAnswers[index]!)},
        progression));

    progression = resp.progression; // update the progression
    nextQuestions = resp.newQuestions; // buffer until retry

    final isCorrect = resp.results[index]!.isCorrect;

    ScaffoldMessenger.of(context).showSnackBar(SnackBar(
      backgroundColor: isCorrect ? Colors.lightGreen : Colors.red.shade200,
      duration: Duration(seconds: isCorrect ? 2 : 4),
      content:
          Text(isCorrect ? "Bonne réponse ! Bravo." : "Réponse incorrecte"),
    ));

    if (progression.nextQuestion == -1) {
      onDone();
      return;
    }

    if (isCorrect) {
      // go to next question
      setState(() {
        results.clear();
        questionIndex = progression.nextQuestion;
      });
    } else {
      // show errors and ask for retry
      setState(() {
        results = {index: resp.results[index]!.results};
      });
    }
  }

  // only validate if all the questions have been completed
  void onValidQuestionParallel() async {
    // check if all the questions are done
    final toSend = <int, Answer>{};
    int? goToQuestion;
    for (var index = 0; index < questions.length; index++) {
      final history = progression.questions[index];
      if (history.isEmpty || !history.last) {
        // the question must be answered
        if (!currentAnswers.containsKey(index)) {
          // go to this question
          goToQuestion = index;
          break;
        } else {
          // add it to the send answsers
          toSend[index] =
              Answer(questions[index].params, currentAnswers[index]!);
        }
      }
    }

    if (goToQuestion != null) {
      // there are still some questions to to
      setState(() {
        questionIndex = goToQuestion;
      });
      return;
    }

    // all good, lets send the results
    final resp = await widget.onEvaluate(
        EvaluateExerciceIn(widget.data.exercice.id, toSend, progression));

    progression = resp.progression; // update the progression
    nextQuestions = resp.newQuestions; // buffer until retry

    if (progression.nextQuestion == -1) {
      onDone();
      return;
    }

    // display the incorrect answers
    final wrongAnswersPlus1 = resp.results.keys
        .where((index) => !resp.results[index]!.isCorrect)
        .map((e) => e + 1);
    final text = (wrongAnswersPlus1.length > 1)
        ? "Les questions ${wrongAnswersPlus1.join(', ')} sont incorrectes."
        : "La question ${wrongAnswersPlus1.first} est incorrecte.";

    ScaffoldMessenger.of(context).showSnackBar(SnackBar(
      backgroundColor: Colors.red.shade200,
      duration: const Duration(seconds: 4),
      content: Text(text),
    ));

    // display the errors and go to the first wrong question
    setState(() {
      results =
          resp.results.map((index, value) => MapEntry(index, value.results));
      questionIndex = progression.nextQuestion;
    });
  }

  void onRetryQuestion() {
    setState(() {
      questions = nextQuestions;
      currentAnswers.clear();
      results.clear();
    });
  }

  void onValideQuestion(QuestionAnswersIn answer) async {
    currentAnswers[questionIndex!] = answer;
    switch (widget.data.exercice.flow) {
      case Flow.sequencial:
        return onValidQuestionSequential();
      case Flow.parallel:
        return onValidQuestionParallel();
    }
  }

  bool get goToPreviousEnabled => questionIndex != null;

  bool get goToNextEnabled {
    switch (widget.data.exercice.flow) {
      case Flow.sequencial: // do not show locked questions
        return questionIndex == null ||
            questionIndex! < progression.nextQuestion;
      case Flow.parallel: // no restriction:
        return questionIndex != questions.length - 1;
    }
  }

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
                feedback: results[questionIndex!],
                answer: currentAnswers[questionIndex!]?.data,
                onRetry: onRetryQuestion,
              ));
  }
}
