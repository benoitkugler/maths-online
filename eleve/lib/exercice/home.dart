import 'package:eleve/shared/title.dart';
import 'package:eleve/types/src.dart';
import 'package:eleve/types/src_maths_questions_client.dart';
import 'package:eleve/types/src_sql_editor.dart';
import 'package:eleve/types/src_sql_tasks.dart';
import 'package:eleve/types/src_tasks.dart';
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
  final StudentWork data;
  final List<QuestionStatus> states;
  final void Function(int index) onSelectQuestion;
  final bool noticeSandbox;

  const ExerciceHome(
      this.data, this.states, this.onSelectQuestion, this.noticeSandbox,
      {Key? key})
      : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Center(
      child: Column(children: [
        Padding(
          padding: const EdgeInsets.symmetric(horizontal: 6.0, vertical: 10),
          child: ColoredTitle(data.exercice.title, Colors.purple),
        ),
        if (noticeSandbox)
          const Card(
            margin: EdgeInsets.only(bottom: 10),
            child: Padding(
              padding: EdgeInsets.all(8.0),
              child: Text("Ta progression n'est pas enregistrée."),
            ),
          ),
        Expanded(
          child: _QuestionList(data, states, onSelectQuestion),
        )
      ]),
    );
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
  final StudentWork data;
  final List<QuestionStatus> states;

  final void Function(int index) onSelectQuestion;

  const _QuestionList(this.data, this.states, this.onSelectQuestion, {Key? key})
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

  void _showProgressionDetails(BuildContext context, int questionIndex) {
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
                            .getQuestion(questionIndex)
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
    // if the question has been validated, always allow access
    if (states[questionIndex] == QuestionStatus.checked) {
      return true;
    }

    switch (data.exercice.flow) {
      case Flow.sequencial:
        return data.progression.nextQuestion == questionIndex;
      case Flow.parallel:
        return true;
    }
  }

  @override
  Widget build(BuildContext context) {
    final mb = mark;
    return ListView(
      children: [
        ...List<Widget>.generate(
          data.exercice.questions.length,
          (index) => _QuestionRow(
            states[index],
            "Question ${index + 1}",
            data.exercice.questions[index].difficulty,
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
                style: const TextStyle(fontSize: 14)),
          )
      ],
    );
  }
}

enum QuestionStatus { locked, checked, toDo, incorrect }

extension _Icon on QuestionStatus {
  Icon get icon {
    switch (this) {
      case QuestionStatus.locked:
        return const Icon(
          IconData(0xf889, fontFamily: 'MaterialIcons'),
          color: Colors.grey,
        );
      case QuestionStatus.checked:
        return const Icon(IconData(0xe156, fontFamily: 'MaterialIcons'),
            color: Colors.green);
      case QuestionStatus.toDo:
        return const Icon(
          assignementIcon,
          color: Colors.purpleAccent,
        );
      case QuestionStatus.incorrect:
        return const Icon(IconData(0xf647, fontFamily: 'MaterialIcons'),
            color: Colors.red);
    }
  }
}

const _difficulties = {
  DifficultyTag.diff1: "★",
  DifficultyTag.diff2: "★★",
  DifficultyTag.diff3: "★★★",
  DifficultyTag.diffEmpty: ""
};

class _QuestionRow extends StatelessWidget {
  final QuestionStatus state;
  final String title;
  final DifficultyTag difficultyTag;
  final int bareme;
  final void Function() showDetails;
  final void Function()? onClick;

  const _QuestionRow(this.state, this.title, this.difficultyTag, this.bareme,
      {required this.showDetails, required this.onClick, Key? key})
      : super(key: key);

  @override
  Widget build(BuildContext context) {
    final diff = _difficulties[difficultyTag] ?? "";
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 4.0),
      child: ListTile(
        shape: const RoundedRectangleBorder(
            borderRadius: BorderRadius.all(Radius.circular(4))),
        tileColor: state == QuestionStatus.toDo
            ? Colors.purple.shade400.withOpacity(0.5)
            : null,
        leading: OutlinedButton(onPressed: showDetails, child: state.icon),
        title: Text(title),
        subtitle: diff.isEmpty ? null : Text(diff),
        trailing: Text("/ $bareme"),
        onTap: onClick,
      ),
    );
  }
}
