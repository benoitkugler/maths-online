import 'package:eleve/exercice/home.dart';
import 'package:eleve/homework/homework.dart';
import 'package:eleve/homework/types.gen.dart';
import 'package:eleve/homework/utils.dart';
import 'package:eleve/shared/title.dart';
import 'package:flutter/material.dart';

extension _ext on SheetProgression {
  // get bareme {
  //   this.exercices.map((e) => e.exercice.) reduce((value, element) => null)
  // }
}

class SheetW extends StatelessWidget {
  final HomeworkAPI api;
  final SheetProgression sheet;

  const SheetW(this.api, this.sheet, {Key? key}) : super(key: key);

  void _startExercice(
      ExerciceProgressionHeader ex, BuildContext context) async {
    showDialog(
        barrierDismissible: false,
        context: context,
        builder: (context) => Dialog(
              child: Padding(
                padding:
                    const EdgeInsets.symmetric(vertical: 16.0, horizontal: 8),
                child: Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: const [
                    Text(
                      "Chargement de l'exercice...",
                      style: TextStyle(fontSize: 16),
                    ),
                    CircularProgressIndicator(),
                  ],
                ),
              ),
            ));
    final instantiatedExercice = await api.loadExercice(ex.exercice);
    Navigator.of(context).pop(); // remove the dialog
    print(instantiatedExercice); // TODO: show the actual exercice route

    // Navigator.of(context).push(MaterialPageRoute<void>(
    //     builder: (context) => ExerciceW(buildMode, ex)));
  }

  @override
  Widget build(BuildContext context) {
    final hasNotation = sheet.sheet.notation != Notation.noNotation;
    return Scaffold(
      appBar: AppBar(title: const Text("Fiche de travail")),
      body: Padding(
        padding: const EdgeInsets.all(8.0),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Padding(
              padding: const EdgeInsets.all(8.0),
              child: Center(
                  child: ColoredTitle(sheet.sheet.title, Colors.blueAccent)),
            ),
            if (hasNotation)
              Padding(
                padding: const EdgeInsets.symmetric(vertical: 8.0),
                child: Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    const Text("Travail noté", style: TextStyle(fontSize: 18)),
                    Text("pour le ${formatTime(sheet.sheet.deadline)}",
                        style: const TextStyle(fontSize: 18)),
                  ],
                ),
              ),
            Expanded(
                child: _ExerciceList(sheet.exercices, hasNotation,
                    (ex) => _startExercice(ex, context))),
          ],
        ),
      ),
    );
  }
}

class _ExerciceList extends StatelessWidget {
  final List<ExerciceProgressionHeader> exercices;
  final bool hasNotation;
  final void Function(ExerciceProgressionHeader) onStart;

  const _ExerciceList(this.exercices, this.hasNotation, this.onStart,
      {Key? key})
      : super(key: key);

  @override
  Widget build(BuildContext context) {
    final children = exercices
        .map((ex) => ListTile(
              onTap: () => onStart(ex),
              leading: getCompletion(ex).icon,
              title: Text(ex.exercice.title),
              trailing: hasNotation ? Text("TODO: note") : null,
            ))
        .toList();
    if (hasNotation) {
      // add the total score of the sheet
      children.add(ListTile(
        title: const Text("Total"),
        trailing: Text("TODO: total"),
      ));
    }
    return ListView(children: children);
  }
}

enum ExerciceCompletion { notStarted, started, completed }

ExerciceCompletion getCompletion(ExerciceProgressionHeader ex) {
  if (!ex.hasProgression) {
    return ExerciceCompletion.notStarted;
  }
  if (ex.progression.isCompleted()) {
    return ExerciceCompletion.completed;
  }
  return ExerciceCompletion.started;
}

extension IconE on ExerciceCompletion {
  Icon get icon {
    switch (this) {
      case ExerciceCompletion.notStarted:
        return const Icon(IconData(0xf587,
            fontFamily: 'MaterialIcons', matchTextDirection: true));
      case ExerciceCompletion.started:
        return const Icon(
            IconData(0xf587,
                fontFamily: 'MaterialIcons', matchTextDirection: true),
            color: Colors.orange);
      case ExerciceCompletion.completed:
        return const Icon(IconData(0xe156, fontFamily: 'MaterialIcons'),
            color: Colors.green);
    }
  }
}
