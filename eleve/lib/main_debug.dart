import 'package:eleve/homework/homework.dart';
import 'package:eleve/homework/types.gen.dart';
import 'package:eleve/main_shared.dart';
import 'package:eleve/shared_gen.dart';
import 'package:flutter/material.dart' hide Flow;

void main() async {
  runApp(const _App());
}

const ex = Exercice(
    1, "Pythagore", "Facile !", Parameters([], []), Flow.sequencial, 1, false);
const exNotStarted = ExerciceProgressionHeader(
    ex, false, ProgressionExt(Progression(0, 0), [[], []], -1), 2, 10);
const exStarted = ExerciceProgressionHeader(
    ex,
    true,
    ProgressionExt(
        Progression(0, 0),
        [
          [false]
        ],
        -1),
    0,
    10);
const exCompleted = ExerciceProgressionHeader(
    ex,
    true,
    ProgressionExt(
        Progression(0, 0),
        [
          [false, true]
        ],
        -1),
    8,
    10);

var sheets = [
  SheetProgression(
      Sheet(
          1, 1, "Fonctions affines", Notation.noNotation, true, DateTime.now()),
      [exStarted, exStarted, exCompleted, exNotStarted]),
  SheetProgression(
      Sheet(2, 1, "Fonctions affines", Notation.successNotation, true,
          DateTime.now().subtract(const Duration(days: 1))),
      [exStarted, exStarted]),
  SheetProgression(
      Sheet(3, 1, "Fonctions affines", Notation.successNotation, true,
          DateTime.now().subtract(const Duration(days: 2))),
      [exCompleted, exCompleted, exNotStarted])
];

class _TestHomeworkAPI implements HomeworkAPI {
  @override
  Future<Sheets> loadSheets() async {
    await Future<void>.delayed(const Duration(seconds: 1));
    return sheets;
    // return [];
  }

  @override
  Future<InstantiatedExercice> loadExercice(int idExercice) async {
    await Future<void>.delayed(const Duration(seconds: 2));

    return InstantiatedExercice(ex, [], []);
  }

  @override
  Future<StudentEvaluateExerciceOut> evaluateExercice(
      int idSheet, int index, EvaluateExerciceIn ex) async {
    await Future<void>.delayed(const Duration(seconds: 2));
    return StudentEvaluateExerciceOut(
        EvaluateExerciceOut({}, ProgressionExt(Progression(0, 0), [], -1), []),
        8);
  }
}

class _App extends StatelessWidget {
  const _App({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Isyro',
      theme: theme,
      debugShowCheckedModeBanner: false,
      localizationsDelegates: localizations,
      supportedLocales: locales,
      home: Homework(_TestHomeworkAPI()),
    );
  }
}
