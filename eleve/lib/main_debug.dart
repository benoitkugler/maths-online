import 'package:eleve/exercice/exercice.dart';
import 'package:eleve/main_shared.dart';
import 'package:eleve/questions/debug.dart';
import 'package:eleve/questions/types.gen.dart';
import 'package:eleve/shared_gen.dart';
import 'package:flutter/material.dart' hide Flow;

void main() async {
  runApp(const _App());
}

class _API implements ExerciceAPI {
  _API();

  @override
  Future<CheckExpressionOut> checkExpressionSyntax(String expression) {
    // TODO: implement checkExpressionSyntax
    throw UnimplementedError();
  }

  @override
  Future<EvaluateExerciceOut> evaluate(EvaluateExerciceIn params) async {
    return EvaluateExerciceOut(
        params.answers.map(
            (key, value) => MapEntry(key, QuestionAnswersOut({0: true}, {}))),
        ProgressionExt([], 1),
        [quI2, quI2, quI2]);
  }
}

final qu1 = Question("", [
  TextBlock([T("Test 1")], false, false, false),
  NumberFieldBlock(0)
]);
final qu2 = Question("", [
  TextBlock([T("Test 2")], false, false, false),
  NumberFieldBlock(0)
]);

final quI1 = InstantiatedQuestion(0, qu1, []);
final quI2 = InstantiatedQuestion(0, qu2, []);

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
      home: ExerciceW(
          _API(),
          StudentExerciceInst(
              InstantiatedExercice(
                  Exercice(1, "Ex", "", Parameters([], []), Flow.sequencial, 1,
                      true),
                  [quI1, quI1, quI1],
                  [1, 2, 3]),
              ProgressionExt([
                [],
                [],
                [],
              ], 0))),
    );
  }
}
