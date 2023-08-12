import 'package:eleve/activities/homework/homework.dart';
import 'package:eleve/main_shared.dart';
import 'package:eleve/questions/debug.dart';
import 'package:eleve/questions/fields.dart';
import 'package:eleve/types/src.dart';
import 'package:eleve/types/src_maths_questions_client.dart';
import 'package:eleve/types/src_maths_repere.dart';
import 'package:eleve/types/src_prof_homework.dart';
import 'package:eleve/types/src_sql_editor.dart';
import 'package:eleve/types/src_sql_homework.dart';
import 'package:eleve/types/src_sql_tasks.dart';
import 'package:eleve/types/src_tasks.dart';
import 'package:flutter/material.dart' hide Flow;

void main() async {
  runApp(const _HomeworkTestApp());
}

class _FieldAPI implements FieldAPI {
  _FieldAPI();

  @override
  Future<CheckExpressionOut> checkExpressionSyntax(String expression) async {
    return const CheckExpressionOut("", true);
  }
}

final questionComplexe = Question([
  TextBlock([T("A question to exercice many fields")], false, false, false),
  const NumberFieldBlock(0, 10),
  const ExpressionFieldBlock("x=", "", 10, true, 1),
  const ExpressionFieldBlock("", " = 0", 10, false, 2),
  const GeometricConstructionFieldBlock(3, GFPoint(),
      FigureBlock(Figure(Drawings({}, [], [], [], []), bounds, true, true))),
], []);

const questionComplexeAnswers = {
  0: NumberAnswer(11.5),
  1: ExpressionAnswer("x^2 + 4 /8 "),
  2: PointAnswer(IntCoord(3, 8)),
};

Question numberQuestion(String title) {
  return Question([
    TextBlock([T(title)], false, false, false),
    const NumberFieldBlock(0, 10)
  ], []);
}

final qu1 = numberQuestion("Test 1");
final qu2 = numberQuestion("Test 2");
final qu3 = numberQuestion("Test 3");

final quI1 = InstantiatedQuestion(1, qu1, DifficultyTag.diff1, []);
final quI2 = InstantiatedQuestion(2, qu2, DifficultyTag.diff2, []);
final quI3 = InstantiatedQuestion(3, qu3, DifficultyTag.diffEmpty, []);

final quI1bis = InstantiatedQuestion(
    1, numberQuestion("Variante 1"), DifficultyTag.diff1, []);
final quI2bis = InstantiatedQuestion(
    2, numberQuestion("Variante 2"), DifficultyTag.diff2, []);
final quI3bis = InstantiatedQuestion(
    3, numberQuestion("Variante 3"), DifficultyTag.diffEmpty, []);

const qu1Answer = {0: NumberAnswer(0)};
const qu2Answer = {0: NumberAnswer(1)};
const qu3Answer = {0: NumberAnswer(2)};

/// a dev widget testing the behavior of the homework widget
class _HomeworkTestApp extends StatelessWidget {
  const _HomeworkTestApp({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
        title: 'Isyro',
        theme: theme,
        debugShowCheckedModeBanner: false,
        localizationsDelegates: localizations,
        supportedLocales: locales,
        home: Builder(
          builder: (context) => Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              ElevatedButton(
                  onPressed: () => showHomework(context),
                  child: const Text("Homework")),
            ],
          ),
        ));
  }

  void showHomework(BuildContext context) async {
    await Navigator.of(context)
        .push(MaterialPageRoute<void>(builder: (context) => const _Homework()));
  }
}

class _API extends _FieldAPI implements HomeworkAPI {
  @override
  Future<Sheets> loadSheets(bool loadNonNoted) async {
    await Future<void>.delayed(const Duration(milliseconds: 200));
    return [
      SheetProgression(
          1,
          Sheet(1, "Feuille en cours", !loadNonNoted,
              DateTime.now().add(const Duration(days: 3)), 0, true, 1),
          [
            const TaskProgressionHeader(
                1,
                "Ex 1",
                "Nombres complexes",
                true,
                ProgressionExt([
                  [true],
                  [true],
                  [true],
                ], 0),
                2,
                6),
            const TaskProgressionHeader(
                2,
                "Ex 2",
                "Nombres complexes",
                true,
                ProgressionExt([
                  [false]
                ], 0),
                0,
                5),
          ]),
      SheetProgression(
          3,
          Sheet(3, "Autre feuille en cours", !loadNonNoted,
              DateTime.now().add(const Duration(days: 4)), 0, true, 1),
          [
            const TaskProgressionHeader(1, "Ex 1", "Nombres complexes", false,
                ProgressionExt([], 0), 0, 6),
            const TaskProgressionHeader(
                2, "Ex 2", "Entiers", false, ProgressionExt([], 0), 0, 5),
          ]),
      SheetProgression(
          2,
          Sheet(2, "Feuille périmée", !loadNonNoted,
              DateTime.now().subtract(const Duration(days: 3)), 0, true, 1),
          [
            const TaskProgressionHeader(
                1,
                "Ex 1",
                "Patholoigcal loooooooonggggg tittllle .. ..",
                false,
                ProgressionExt([], 0),
                0,
                6),
            const TaskProgressionHeader(
                2, "Ex 2", "", false, ProgressionExt([], 0), 0, 5),
          ]),
      SheetProgression(
          2,
          Sheet(4, "Feuille périmée", !loadNonNoted,
              DateTime.now().subtract(const Duration(days: 3)), 0, true, 1),
          [
            const TaskProgressionHeader(
                1, "Ex 1", "", false, ProgressionExt([], 0), 0, 6),
            const TaskProgressionHeader(
                2, "Ex 2", "", false, ProgressionExt([], 0), 0, 5),
          ]),
      SheetProgression(
          2,
          Sheet(5, "Feuille périmée", !loadNonNoted,
              DateTime.now().subtract(const Duration(days: 3)), 0, true, 1),
          [
            const TaskProgressionHeader(
                1, "Ex 1", "", false, ProgressionExt([], 0), 0, 6),
            const TaskProgressionHeader(
                2, "Ex 2", "", false, ProgressionExt([], 0), 0, 5),
          ]),
      SheetProgression(
          2,
          Sheet(6, "Feuille périmée", !loadNonNoted,
              DateTime.now().subtract(const Duration(days: 3)), 0, true, 1),
          [
            const TaskProgressionHeader(
                1, "Ex 1", "", false, ProgressionExt([], 0), 0, 6),
            const TaskProgressionHeader(
                2, "Ex 2", "", false, ProgressionExt([], 0), 0, 5),
          ]),
      SheetProgression(
          2,
          Sheet(7, "Feuille périmée", !loadNonNoted,
              DateTime.now().subtract(const Duration(days: 3)), 0, true, 1),
          [
            const TaskProgressionHeader(
                1, "Ex 1", "", false, ProgressionExt([], 0), 0, 6),
            const TaskProgressionHeader(
                2, "Ex 2", "", false, ProgressionExt([], 0), 0, 5),
          ]),
    ];
  }

  @override
  Future<StudentEvaluateTaskOut> evaluateExercice(
      IdTask idTask, IdTravail idTravail, EvaluateWorkIn ex) {
    // TODO: implement evaluateExercice
    throw UnimplementedError();
  }

  @override
  Future<InstantiatedWork> loadWork(IdTask id) async {
    return InstantiatedWork(WorkID(0, WorkKind.workExercice, false),
        "Exo de Test", Flow.sequencial, [quI1, quI2, quI3], [2, 3, 1]);
  }

  @override
  Future<void> resetTask(IdTravail idTravail, IdTask idTask) async {
    return;
  }
}

class _Homework extends StatelessWidget {
  const _Homework();

  @override
  Widget build(BuildContext context) {
    return HomeworkStart(_API());
  }
}
