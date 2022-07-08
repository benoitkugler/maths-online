import 'package:eleve/shared/title.dart';
import 'package:eleve/shared_gen.dart';
import 'package:flutter/material.dart' hide Flow;

/// ExerciceHome shows a welcome screen when opening an exercice,
/// with its questions and bareme
class ExerciceHome extends StatelessWidget {
  final Exercice data;
  final void Function(int index) onSelectQuestion;

  const ExerciceHome(this.data, this.onSelectQuestion, {Key? key})
      : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Scaffold(
        body: Center(
      child: Column(children: [
        Padding(
          padding: const EdgeInsets.symmetric(horizontal: 4.0, vertical: 20),
          child: ColoredTitle(data.exercice.title, Colors.purple),
        ),
        Expanded(child: _QuestionList(data, onSelectQuestion))
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

class _QuestionList extends StatelessWidget {
  final Exercice data;
  final void Function(int index) onSelectQuestion;

  const _QuestionList(this.data, this.onSelectQuestion, {Key? key})
      : super(key: key);

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
                        children: data.progression.questions[index]
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
    switch (data.exercice.flow) {
      case Flow.sequencial:
        return data.progression.nextQuestion == questionIndex;
      case Flow.parallel:
        return true;
    }
  }

  _QuestionState state(int questionIndex) {
    final history = data.progression.questions[questionIndex];
    if (history.isNotEmpty && history.last) {
      return _QuestionState.checked;
    }

    if (data.exercice.flow == Flow.sequencial &&
        data.progression.nextQuestion < questionIndex) {
      return _QuestionState.locked;
    }
    return _QuestionState.toDo;
  }

  @override
  Widget build(BuildContext context) {
    return ListView(
      children: List<_QuestionRow>.generate(
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
    );
  }
}

enum _QuestionState { locked, checked, toDo }

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
          IconData(0xe09f, fontFamily: 'MaterialIcons'),
          color: Colors.purpleAccent,
        );
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
    return ListTile(
      title: Text(title),
      leading: OutlinedButton(onPressed: showDetails, child: state.icon),
      trailing: Text("/ $bareme"),
      onTap: onClick,
    );
  }
}
