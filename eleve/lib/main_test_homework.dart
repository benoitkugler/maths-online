import 'package:eleve/activities/homework/homework.dart';
import 'package:eleve/activities/homework/types.gen.dart';
import 'package:eleve/main_shared.dart';
import 'package:eleve/questions/debug.dart';
import 'package:eleve/questions/fields.dart';
import 'package:eleve/questions/repere.gen.dart';
import 'package:eleve/questions/types.gen.dart';
import 'package:eleve/shared_gen.dart';
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
  const FigurePointFieldBlock(
      Figure(Drawings({}, [], [], [], []), bounds, true, true), 3)
]);

const questionComplexeAnswers = {
  0: NumberAnswer(11.5),
  1: ExpressionAnswer("x^2 + 4 /8 "),
  2: PointAnswer(IntCoord(3, 8)),
};

Question numberQuestion(String title) {
  return Question([
    TextBlock([T(title)], false, false, false),
    const NumberFieldBlock(0, 10)
  ]);
}

final qu1 = numberQuestion("Test 1");
final qu2 = numberQuestion("Test 2");
final qu3 = numberQuestion("Test 3");

final quI1 = InstantiatedQuestion(1, qu1, []);
final quI2 = InstantiatedQuestion(2, qu2, []);
final quI3 = InstantiatedQuestion(3, qu3, []);

final quI1bis = InstantiatedQuestion(1, numberQuestion("Variante 1"), []);
final quI2bis = InstantiatedQuestion(2, numberQuestion("Variante 2"), []);
final quI3bis = InstantiatedQuestion(3, numberQuestion("Variante 3"), []);

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
  Future<Sheets> loadSheets() async {
    await Future<void>.delayed(const Duration(seconds: 1));
    return [
      SheetProgression(
          Sheet(1, 1, "Feuille de test", Notation.successNotation, true,
              DateTime.now().add(const Duration(days: 3))),
          [
            const TaskProgressionHeader(
                1, "Ex 1", false, ProgressionExt([], 0), 0, 6),
            const TaskProgressionHeader(
                2, "Ex 2", false, ProgressionExt([], 0), 0, 5),
          ])
    ];
  }

  @override
  Future<StudentEvaluateTaskOut> evaluateExercice(
      IdTask idTask, EvaluateWorkIn ex) {
    // TODO: implement evaluateExercice
    throw UnimplementedError();
  }

  @override
  Future<InstantiatedWork> loadWork(IdTask id) {
    // TODO: implement loadWork
    throw UnimplementedError();
  }
}

class _Homework extends StatelessWidget {
  const _Homework({super.key});

  @override
  Widget build(BuildContext context) {
    return HomeworkW(_API());
  }
}
