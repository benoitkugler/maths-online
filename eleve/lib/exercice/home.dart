import 'package:eleve/questions/types.gen.dart';
import 'package:eleve/shared/title.dart';
import 'package:eleve/shared_gen.dart';
import 'package:flutter/material.dart' hide Flow;

extension IsCorrect on QuestionAnswersOut {
  /// isCorrect is true if every fields are correct
  bool get isCorrect {
    return results.values.every((success) => success);
  }
}

const assignementIcon =
    IconData(0xf587, fontFamily: 'MaterialIcons', matchTextDirection: true);
const completedIcon = IconData(0xe156, fontFamily: 'MaterialIcons');

extension ProgressionExtension on ProgressionExt {
  /// [getQuestion] returns an empty list if progression is empty
  QuestionHistory getQuestion(int index) {
    if (questions.length <= index) {
      return [];
    }
    return questions[index];
  }

  bool _isQuestionCompleted(List<bool> history) {
    return history.isNotEmpty && history.last;
  }

  /// returns `true` if the question at [index] is completed
  bool isQuestionCompleted(int index) {
    return _isQuestionCompleted(getQuestion(index));
  }

  /// returns `true` if all the questions of the exercice are completed
  bool isCompleted() {
    return questions.every(_isQuestionCompleted);
  }
}

/// ExerciceHome shows a welcome screen when opening an exercice,
/// with its questions and bareme
class ExerciceHome extends StatelessWidget {
  final StudentExerciceInst data;
  final Set<int> validatedQuestions;
  final Set<int> incorrectQuestions;
  final void Function(int index) onSelectQuestion;

  const ExerciceHome(this.data, this.validatedQuestions,
      this.incorrectQuestions, this.onSelectQuestion,
      {Key? key})
      : super(key: key);

  @override
  Widget build(BuildContext context) {
    print(incorrectQuestions);
    print(validatedQuestions);
    return Scaffold(
        body: Center(
      child: Column(children: [
        Padding(
          padding: const EdgeInsets.symmetric(horizontal: 4.0, vertical: 20),
          child: ColoredTitle(data.exercice.exercice.title, Colors.purple),
        ),
        Expanded(
            child: _QuestionList(
                data, validatedQuestions, incorrectQuestions, onSelectQuestion))
      ]),
    ));
  }
}

class _SuccessSquare extends StatelessWidget {
  final bool success;
  const _SuccessSquare(this.success, {Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.all(8.0),
      child: Container(
        height: 30,
        width: 30,
        color: success ? Colors.green : Colors.red,
      ),
    );
  }
}

class MarkBareme {
  final int mark;
  final int bareme;
  MarkBareme(this.mark, this.bareme);
}

class _QuestionList extends StatelessWidget {
  final StudentExerciceInst data;
  final Set<int> validatedQuestions;
  final Set<int> incorrectQuestions; // temporary indication
  final void Function(int index) onSelectQuestion;

  const _QuestionList(this.data, this.validatedQuestions,
      this.incorrectQuestions, this.onSelectQuestion,
      {Key? key})
      : super(key: key);

  MarkBareme get mark {
    int mark = 0;
    int bareme = 0;
    for (var i = 0; i < data.exercice.baremes.length; i++) {
      bareme += data.exercice.baremes[i];
      if (data.progression.isQuestionCompleted(i)) {
        mark += data.exercice.baremes[i];
      }
    }
    return MarkBareme(mark, bareme);
  }

  void _showProgressionDetails(BuildContext context, int index) {
    showDialog<void>(
        context: context,
        builder: (context) => Dialog(
              child: Card(
                child: Column(
                  mainAxisSize: MainAxisSize.min,
                  children: [
                    const Padding(
                      padding: EdgeInsets.all(20.0),
                      child: Text(
                        "Historique de tes tentatives",
                        style: TextStyle(fontSize: 20),
                      ),
                    ),
                    Padding(
                      padding: const EdgeInsets.all(8.0),
                      child: Wrap(
                        children: data.progression
                            .getQuestion(index)
                            .map((e) => _SuccessSquare(e))
                            .toList(),
                      ),
                    )
                  ],
                ),
              ),
            ));
  }

  bool allowDoQuestion(int questionIndex) {
    switch (data.exercice.exercice.flow) {
      case Flow.sequencial:
        return data.progression.nextQuestion == questionIndex;
      case Flow.parallel:
        return true;
    }
  }

  _QuestionState state(int questionIndex) {
    if (data.progression.isQuestionCompleted(questionIndex)) {
      return _QuestionState.checked;
    }

    if (data.exercice.exercice.flow == Flow.sequencial &&
        data.progression.nextQuestion < questionIndex) {
      return _QuestionState.locked;
    }

    // after validating, both validatedQuestions and incorrectQuestions
    // may contain the same index : give the priority to incorrectQuestions
    if (incorrectQuestions.contains(questionIndex)) {
      return _QuestionState.incorrect;
    } else if (validatedQuestions.contains(questionIndex)) {
      return _QuestionState.waitingCorrection;
    }
    return _QuestionState.toDo;
  }

  @override
  Widget build(BuildContext context) {
    final mb = mark;
    return ListView(
      children: [
        ...List<Widget>.generate(
          data.exercice.questions.length,
          (index) => _QuestionRow(
            state(index),
            "Question ${index + 1}",
            data.exercice.baremes[index],
            showDetails: () => _showProgressionDetails(context, index),
            onClick:
                allowDoQuestion(index) ? () => onSelectQuestion(index) : null,
          ),
        ),
        if (data.progression.questions.isNotEmpty)
          ListTile(
            title: const Text("Total"),
            trailing: Text("${mb.mark} / ${mb.bareme}",
                style: TextStyle(fontSize: 14)),
          )
      ],
    );
  }
}

enum _QuestionState { locked, checked, toDo, waitingCorrection, incorrect }

extension _Icon on _QuestionState {
  Icon get icon {
    switch (this) {
      case _QuestionState.locked:
        return const Icon(
          IconData(0xf889, fontFamily: 'MaterialIcons'),
          color: Colors.grey,
        );
      case _QuestionState.checked:
        return const Icon(IconData(0xe156, fontFamily: 'MaterialIcons'),
            color: Colors.green);
      case _QuestionState.toDo:
        return const Icon(
          assignementIcon,
          color: Colors.purpleAccent,
        );
      case _QuestionState.waitingCorrection:
        return const Icon(IconData(0xf51a, fontFamily: 'MaterialIcons'),
            color: Colors.orange);
      case _QuestionState.incorrect:
        return const Icon(IconData(0xf647, fontFamily: 'MaterialIcons'),
            color: Colors.red);
    }
  }
}

class _QuestionRow extends StatelessWidget {
  final _QuestionState state;
  final String title;
  final int bareme;

  final void Function() showDetails;
  final void Function()? onClick;
  const _QuestionRow(this.state, this.title, this.bareme,
      {required this.showDetails, required this.onClick, Key? key})
      : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.all(4.0),
      child: ListTile(
        shape: const RoundedRectangleBorder(
            borderRadius: BorderRadius.all(Radius.circular(4))),
        tileColor: state == _QuestionState.toDo
            ? Colors.purple.withOpacity(0.5)
            : null,
        leading: OutlinedButton(onPressed: showDetails, child: state.icon),
        title: Text(title),
        trailing: Text("/ $bareme"),
        onTap: onClick,
      ),
    );
  }
}
