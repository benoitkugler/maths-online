import 'package:eleve/exercice/exercice.dart';
import 'package:eleve/main_shared.dart';
import 'package:eleve/questions/debug.dart';
import 'package:eleve/questions/question.dart';
import 'package:eleve/questions/repere.gen.dart';
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
  NumberFieldBlock(0, 10)
]);
final qu2 = Question("", [
  TextBlock([T("Test 2")], false, false, false),
  NumberFieldBlock(0, 10)
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
        home: Scaffold(
          body: QuestionW(
            _API(),
            Question("", [
              TextBlock([T("Un long text et suffisant")], true, false, false),
              FigureBlock(Figure(
                  Drawings({}, [], [], [
                    Circle(Coord(0, 0), 4, "#FF0000", "invalid", "C_f"),
                    Circle(Coord(5, 5), 2, "#FF00FF", "#55FFFF00", "C_g"),
                  ], []),
                  bounds,
                  true,
                  true)),
            ]),
            Colors.green,
            (_) {},
          ),
        ));
  }
}
