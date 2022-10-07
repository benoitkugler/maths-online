import 'package:eleve/exercice/exercice.dart';
import 'package:eleve/main_shared.dart';
import 'package:eleve/questions/debug.dart';
import 'package:eleve/questions/types.gen.dart';
import 'package:eleve/shared_gen.dart';
import 'package:flutter/material.dart' hide Flow;

void main() async {
  runApp(const _QuestionTestApp());
}

class _API implements ExerciceAPI {
  _API();

  @override
  Future<CheckExpressionOut> checkExpressionSyntax(String expression) async {
    return CheckExpressionOut("", true);
  }

  @override
  Future<EvaluateWorkOut> evaluate(EvaluateWorkIn params) async {
    final indexQuestion = params.answers.keys.first;
    return EvaluateWorkOut(
        params.answers.map((key, value) =>
            MapEntry(key, QuestionAnswersOut({0: key == 1}, {}))),
        ProgressionExt([
          [false],
          [true]
        ], indexQuestion == 0 ? 1 : -1),
        [quI1, quI2]);
  }
}

final qu1 = Question([
  TextBlock([T("Test 1")], false, false, false),
  NumberFieldBlock(0, 10)
]);
final qu2 = Question([
  TextBlock([T("Test 2")], false, false, false),
  NumberFieldBlock(0, 10)
]);

final quI1 = InstantiatedQuestion(0, qu1, []);
final quI2 = InstantiatedQuestion(0, qu2, []);

final exercice = ExerciceController(
    StudentWork(
      InstantiatedWork(
          WorkID(0, true), "", Flow.parallel, [quI1, quI2], [1, 1]),
      ProgressionExt([], 0),
    ),
    null);

/// a dev widget testing the behavior of the question/exercice
/// widgets for each context
class _QuestionTestApp extends StatelessWidget {
  const _QuestionTestApp({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
        title: 'Isyro',
        theme: theme,
        debugShowCheckedModeBanner: false,
        localizationsDelegates: localizations,
        supportedLocales: locales,
        home: Row(
          children: [
            ElevatedButton(
                onPressed: showTrivialInGame,
                child: const Text("Trivial: InGame")),
            ElevatedButton(
                onPressed: showTrivialLast, child: const Text("Trivial: Last")),
            ElevatedButton(
                onPressed: showDecrassage, child: const Text("Decrassage")),
            ElevatedButton(
                onPressed: showLoopackQuestion,
                child: const Text("Loopack: Question")),
            ElevatedButton(
                onPressed: showExerciceSequencial,
                child: const Text("Homework: Sequencial")),
            ElevatedButton(
                onPressed: showExerciceParallel,
                child: const Text("Homework: Parallel")),
            ElevatedButton(
                onPressed: showLoopackExercice,
                child: const Text("Loopack: Exercice")),
          ],
        ));
  }

  showTrivialInGame() {}
  showTrivialLast() {}
  showDecrassage() {}
  showLoopackQuestion() {}
  showExerciceSequencial() {}
  showExerciceParallel() {}
  showLoopackExercice() {}
}
