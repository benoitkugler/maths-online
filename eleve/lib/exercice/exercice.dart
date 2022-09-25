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
  Future<EvaluateWorkOut> evaluate(EvaluateWorkIn params);
}

/// [ExerciceController] exposes the state
/// which will change during an exercice,
/// so that parent widget may react properly.
class ExerciceController {
  /// [exeAndProg] stores the server instantiated exercice with
  /// the progression state.
  StudentWork exeAndProg;

  /// [questionIndex] is the current question, or null for the summary
  int? questionIndex;

  ExerciceController(this.exeAndProg, this.questionIndex);
}

extension _CopyIW on InstantiatedWork {
  InstantiatedWork copyWithQuestions(List<InstantiatedQuestion> questions) {
    return InstantiatedWork(
        iD,
        title,
        flow,
        // replace the questions
        questions,
        baremes);
  }
}

extension _CopySW on StudentWork {
  StudentWork copyWithQuestions(List<InstantiatedQuestion> questions) {
    return StudentWork(exercice.copyWithQuestions(questions), progression);
  }
}

/// ExerciceW is the widget providing one exercice to
/// the student.
/// It is used in the editor loopback, and as the base for
/// at home training activity
class ExerciceW extends StatefulWidget {
  final ExerciceAPI api;

  final ExerciceController controller;

  /// If not null, [questionAnswers] overide the user answers
  /// for the current question
  final Answers? questionAnswers;

  final void Function()? onShowCorrectAnswer;

  const ExerciceW(this.api, this.controller,
      {Key? key, this.questionAnswers, this.onShowCorrectAnswer})
      : super(key: key);

  @override
  State<ExerciceW> createState() => _ExerciceWState();
}

class _ExerciceWState extends State<ExerciceW> {
  // the questions to display when trying again
  // this is a temporay slice, affected to the controller on user validation
  List<InstantiatedQuestion> nextQuestions = [];

  // the currrent answsers of the student, filled
  // when validating a question
  Map<int, QuestionAnswersIn> currentAnswers = {};

  // the feedback to display
  Map<int, QuestionAnswersOut> feedback = {};

  @override
  void initState() {
    reset();
    super.initState();
  }

  @override
  void didUpdateWidget(covariant ExerciceW oldWidget) {
    reset();
    super.didUpdateWidget(oldWidget);
  }

  void reset() {
    nextQuestions = [];
    currentAnswers.clear();
    feedback.clear();
    // show answers from the server
    if (widget.controller.questionIndex != null &&
        widget.questionAnswers != null) {
      currentAnswers[widget.controller.questionIndex!] =
          QuestionAnswersIn(widget.questionAnswers!);
    }
  }

  // handle the errors
  Future<EvaluateWorkOut?> _evaluate(EvaluateWorkIn params) async {
    try {
      final res = await widget.api.evaluate(params);
      return res;
    } catch (error) {
      showError("Impossible d'évaluer l'exercice", error, context);
      return null;
    }
  }

  void onExerciceOver() async {
    setState(() {
      widget.controller.questionIndex = null;
    });

    // exercice is over
    await showDialog<void>(
        context: context,
        builder: (context) => const Dialog(child: Congrats()));
    NotificationExerciceDone().dispatch(context);
  }

  void onValidQuestionSequential() async {
    final ct = widget.controller;
    final index = ct.questionIndex!;

    // if we are not at the current question, just go to it
    // and return
    if (index != ct.exeAndProg.progression.nextQuestion) {
      setState(() {
        ct.questionIndex = ct.exeAndProg.progression.nextQuestion;
      });
      return;
    }

    // validate the given answer
    final resp = await _evaluate(EvaluateWorkIn(
        ct.exeAndProg.exercice.iD,
        {
          index: Answer(ct.exeAndProg.exercice.questions[index].params,
              currentAnswers[index]!)
        },
        ct.exeAndProg.progression));
    if (resp == null) {
      return;
    }

    ct.exeAndProg = StudentWork(
        ct.exeAndProg.exercice, resp.progression); // update the progression
    nextQuestions = resp.newQuestions; // buffer until retry

    final isCorrect = resp.results[index]!.isCorrect;
    final hasNextQuestion =
        isCorrect && ct.exeAndProg.progression.nextQuestion != -1;
    ScaffoldMessenger.of(context).showSnackBar(SnackBar(
      backgroundColor: isCorrect ? Colors.lightGreen : Colors.red.shade200,
      duration: Duration(
          seconds: hasNextQuestion ? 10000 : 4), // block on correct answer
      content:
          Text(isCorrect ? "Bonne réponse ! Bravo." : "Réponse incorrecte"),
      action: hasNextQuestion
          ? SnackBarAction(
              label: "Continuer l'exercice",
              textColor: Colors.black,
              onPressed: () {
                setState(() {
                  ct.questionIndex = ct.exeAndProg.progression.nextQuestion;
                });
              },
            )
          : null,
    ));

    if (ct.exeAndProg.progression.nextQuestion == -1) {
      onExerciceOver();
      return;
    }

    if (isCorrect) {
      setState(() {
        feedback.clear();
      });
      // wait for the user to go to the next question
    } else {
      // show errors and ask for retry
      setState(() {
        feedback = {index: resp.results[index]!};
      });
    }
  }

  // only validate if all the questions have been completed
  void onValidQuestionParallel() async {
    final ct = widget.controller;
    final questions = ct.exeAndProg.exercice.questions;

    // check if all the questions are done
    final toSend = <int, Answer>{};
    int? goToQuestion;
    for (var index = 0; index < questions.length; index++) {
      final history = ct.exeAndProg.progression.getQuestion(index);
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
        widget.controller.questionIndex = goToQuestion;
      });
      return;
    }

    // all good, lets send the results
    final resp = await _evaluate(EvaluateWorkIn(
        ct.exeAndProg.exercice.iD, toSend, ct.exeAndProg.progression));
    if (resp == null) {
      return;
    }

    ct.exeAndProg = StudentWork(
        ct.exeAndProg.exercice, resp.progression); // update the progression
    nextQuestions = resp.newQuestions; // buffer until retry

    if (isExerciceOver) {
      onExerciceOver();
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
      feedback = resp.results;
      widget.controller.questionIndex = null;
    });
  }

  void onRetryQuestion() {
    final ct = widget.controller;
    setState(() {
      ct.exeAndProg = ct.exeAndProg.copyWithQuestions(nextQuestions);
      currentAnswers.clear();
      feedback.clear();
    });
  }

  void onValideQuestion(QuestionAnswersIn answer) async {
    currentAnswers[widget.controller.questionIndex!] = answer;
    switch (widget.controller.exeAndProg.exercice.flow) {
      case Flow.sequencial:
        return onValidQuestionSequential();
      case Flow.parallel:
        return onValidQuestionParallel();
    }
  }

  bool get isExerciceOver =>
      widget.controller.exeAndProg.progression.nextQuestion == -1;

  bool get goToPreviousEnabled => widget.controller.questionIndex != null;

  bool get goToNextEnabled {
    final ex = widget.controller.exeAndProg;
    final currentIndex = widget.controller.questionIndex;
    final hasNextQuestion = currentIndex != ex.exercice.questions.length - 1;
    switch (ex.exercice.flow) {
      case Flow
          .sequencial: // do not show locked questions when exercice is not over
        return currentIndex == null ||
            (isExerciceOver && hasNextQuestion) ||
            currentIndex < ex.progression.nextQuestion;
      case Flow.parallel: // no restriction:
        return hasNextQuestion;
    }
  }

  void goToPrevious() {
    final newIndex = widget.controller.questionIndex == 0
        ? null
        : widget.controller.questionIndex! - 1;
    setState(() {
      widget.controller.questionIndex = newIndex;
    });
  }

  void goToNext() {
    // remove potential snackbar
    ScaffoldMessenger.of(context).hideCurrentSnackBar();

    final newIndex = widget.controller.questionIndex == null
        ? 0
        : widget.controller.questionIndex! + 1;
    setState(() {
      widget.controller.questionIndex = newIndex;
    });
  }

  List<QuestionState> get questionStates {
    final exP = widget.controller.exeAndProg;
    final validatedQuestions = currentAnswers.keys.toSet();
    final incorrectQuestions =
        feedback.keys.where((index) => !feedback[index]!.isCorrect).toSet();

    return List<QuestionState>.generate(exP.exercice.questions.length,
        (questionIndex) {
      if (exP.progression.isQuestionCompleted(questionIndex)) {
        return QuestionState.checked;
      }

      if (exP.exercice.flow == Flow.sequencial &&
          exP.progression.nextQuestion < questionIndex) {
        return QuestionState.locked;
      }

      // after validating, both validatedQuestions and incorrectQuestions
      // may contain the same index : give the priority to incorrectQuestions
      if (incorrectQuestions.contains(questionIndex)) {
        return QuestionState.incorrect;
      } else if (validatedQuestions.contains(questionIndex)) {
        return QuestionState.waitingCorrection;
      }
      return QuestionState.toDo;
    });
  }

  Future<bool> confirmeLeave() async {
    final res = await showDialog<bool?>(
        context: context,
        builder: (BuildContext context) {
          return AlertDialog(
            title: const Text("Abandonner l'exercice"),
            content: const Text(
                "Es-tu sûr d'abandonner l'exercice ? \nTes réponses en attente de correction seront effacées."),
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

  QuestionController questionController() {
    final exP = widget.controller.exeAndProg;
    final out = QuestionController(
        exP.exercice.questions[widget.controller.questionIndex!].question,
        widget.api,
        true);
    out.setAnswers(currentAnswers[widget.controller.questionIndex!]?.data);
    // feedback has the priority against answers
    out.setFeedback(feedback[widget.controller.questionIndex!]?.results);
    if (questionStates[widget.controller.questionIndex!] ==
        QuestionState.checked) {
      out.setDone();
    }
    return out;
  }

  @override
  Widget build(BuildContext context) {
    final exP = widget.controller.exeAndProg;
    return Scaffold(
        appBar: AppBar(
          title: const Text("Exercice"),
          actions: [
            if (widget.onShowCorrectAnswer != null &&
                widget.controller.questionIndex != null)
              TextButton(
                  onPressed: widget.onShowCorrectAnswer,
                  child: const Text("Afficher la réponse.")),
            IconButton(
                onPressed: goToPreviousEnabled ? goToPrevious : null,
                icon: const Icon(IconData(0xf572,
                    fontFamily: 'MaterialIcons', matchTextDirection: true))),
            IconButton(
                onPressed: goToNextEnabled ? goToNext : null,
                icon: const Icon(IconData(0xf57a,
                    fontFamily: 'MaterialIcons', matchTextDirection: true))),
          ],
        ),
        body: WillPopScope(
          onWillPop: () async {
            final isDirty = questionStates
                .any((st) => st == QuestionState.waitingCorrection);
            if (isDirty) {
              final ok = await confirmeLeave();
              return ok;
            }
            return true; // nothing to loose
          },
          child: widget.controller.questionIndex == null
              ? ExerciceHome(
                  exP,
                  questionStates,
                  (index) => setState(() {
                        widget.controller.questionIndex = index;
                      }))
              : QuestionW(
                  questionController(),
                  Colors.purpleAccent,
                  onValideQuestion,
                  title: "Question ${widget.controller.questionIndex! + 1}",
                  timeout: null,
                  onRetry: onRetryQuestion,
                ),
        ));
  }
}
