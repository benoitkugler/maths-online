import 'package:eleve/exercice/congratulations.dart';
import 'package:eleve/exercice/home.dart';
import 'package:eleve/questions/fields.dart';
import 'package:eleve/questions/question.dart';
import 'package:eleve/questions/types.gen.dart' hide Answer;
import 'package:eleve/shared/errors.dart';
import 'package:eleve/shared_gen.dart';
import 'package:flutter/material.dart' hide Flow;

class NotificationExerciceDone extends Notification {}

abstract class ExerciceAPI extends FieldAPI {
  Future<EvaluateExerciceOut> evaluate(EvaluateExerciceIn params);
}

/// ExerciceW is the widget providing one exercice to
/// the student.
/// It is used in the editor loopback, and as the base for
/// at home training activity
class ExerciceW extends StatefulWidget {
  final ExerciceAPI api;

  /// [data] stores the server instantiated exercice with
  /// the initial progression state.
  final StudentExerciceInst data;

  const ExerciceW(this.api, this.data, {Key? key}) : super(key: key);

  @override
  State<ExerciceW> createState() => _ExerciceWState();
}

class _ExerciceWState extends State<ExerciceW> {
  late List<InstantiatedQuestion> questions; // will change upon wrong answer
  late ProgressionExt progression; // will change on validate

  // the questions to display when trying again
  List<InstantiatedQuestion> nextQuestions = [];

  // the currrent answsers of the student, filled
  // when validating a question
  Map<int, QuestionAnswersIn> currentAnswers = {};
  Map<int, QuestionAnswersOut> results = {};

  int? questionIndex; // null means summary

  @override
  void initState() {
    questions = widget.data.exercice.questions;
    progression = widget.data.progression;
    super.initState();
  }

  @override
  void didUpdateWidget(covariant ExerciceW oldWidget) {
    questions = widget.data.exercice.questions;
    progression = widget.data.progression;
    super.didUpdateWidget(oldWidget);
  }

  // handle the errors
  Future<EvaluateExerciceOut?> _evaluate(EvaluateExerciceIn params) async {
    try {
      final res = await widget.api.evaluate(params);
      return res;
    } catch (error) {
      showError("Impossible d'évaluer l'exercice", error, context);
      return null;
    }
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
    final resp = await _evaluate(EvaluateExerciceIn(
        widget.data.exercice.exercice.id,
        {index: Answer(questions[index].params, currentAnswers[index]!)},
        progression));
    if (resp == null) {
      return;
    }

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
        results = {index: resp.results[index]!};
      });
    }
  }

  // only validate if all the questions have been completed
  void onValidQuestionParallel() async {
    // check if all the questions are done
    final toSend = <int, Answer>{};
    int? goToQuestion;
    for (var index = 0; index < questions.length; index++) {
      final history = progression.getQuestion(index);
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
    final resp = await _evaluate(EvaluateExerciceIn(
        widget.data.exercice.exercice.id, toSend, progression));
    if (resp == null) {
      return;
    }

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

    // display the errors and go to the menu
    setState(() {
      results = resp.results;
      questionIndex = null;
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
    switch (widget.data.exercice.exercice.flow) {
      case Flow.sequencial:
        return onValidQuestionSequential();
      case Flow.parallel:
        return onValidQuestionParallel();
    }
  }

  bool get goToPreviousEnabled => questionIndex != null;

  bool get goToNextEnabled {
    switch (widget.data.exercice.exercice.flow) {
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
                StudentExerciceInst(
                  InstantiatedExercice(
                      ex.exercice,
                      // replace the questions
                      questions,
                      ex.baremes),
                  progression,
                ),
                currentAnswers.keys.toSet(),
                results.keys
                    .where((index) => !results[index]!.isCorrect)
                    .toSet(),
                (index) => setState(() {
                      questionIndex = index;
                    }))
            : QuestionW(
                widget.api,
                questions[questionIndex!].question,
                Colors.purpleAccent,
                onValideQuestion,
                title: "Question ${questionIndex! + 1}",
                timeout: null,
                blockOnSubmit: true,
                feedback: results[questionIndex!]?.results,
                answer: currentAnswers[questionIndex!]?.data,
                onRetry: onRetryQuestion,
              ));
  }
}
