import 'package:eleve/build_mode.dart';
import 'package:eleve/homework/types.gen.dart';
import 'package:eleve/shared/title.dart';
import 'package:flutter/material.dart';

extension _ext on SheetProgression {
  // get bareme {
  //   this.exercices.map((e) => e.exercice.) reduce((value, element) => null)
  // }
}

String formatTime(DateTime time) {
  return "${time.day}/${time.month}/${time.year}, à ${time.hour}h";
}

class SheetW extends StatelessWidget {
  final BuildMode buildMode;
  final SheetProgression sheet;
  const SheetW(this.buildMode, this.sheet, {Key? key}) : super(key: key);

  void _startExercice(ExerciceProgressionHeader ex, BuildContext context) {
    showDialog(
        context: context,
        builder: (context) => Dialog(
              child: Text("Chargement de l'exercice..."),
            ));
    // Navigator.of(context).push(MaterialPageRoute<void>(
    //     builder: (context) => ExerciceLauncher(buildMode, ex)));
  }

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.all(8.0),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Padding(
            padding: const EdgeInsets.all(8.0),
            child:
                Center(child: ColoredTitle(sheet.sheet.title, Colors.purple)),
          ),
          if (sheet.sheet.notation != Notation.noNotation)
            const Padding(
              padding: EdgeInsets.all(8.0),
              child: Text("Noté", style: TextStyle(fontSize: 18)),
            ),
          Padding(
            padding: const EdgeInsets.all(8.0),
            child: Text("A faire pour le : ${formatTime(sheet.sheet.deadline)}",
                style: const TextStyle(fontSize: 18)),
          ),
          Expanded(
              child: _ExerciceList(
                  sheet.exercices, (ex) => _startExercice(ex, context))),
        ],
      ),
    );
  }
}

class _ExerciceList extends StatelessWidget {
  final List<ExerciceProgressionHeader> exercices;
  final void Function(ExerciceProgressionHeader) onStart;

  const _ExerciceList(this.exercices, this.onStart, {Key? key})
      : super(key: key);

  @override
  Widget build(BuildContext context) {
    return ListView(
        children: exercices
            .map((ex) => ListTile(
                  onTap: () => onStart(ex),
                  title: Text(ex.exercice.title),
                ))
            .toList());
  }
}
