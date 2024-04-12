import 'package:eleve/loopback/loopback.dart';
import 'package:eleve/main_shared.dart';
import 'package:eleve/main_test_questions.dart';
import 'package:eleve/questions/debug.dart';
import 'package:eleve/types/src_prof_preview.dart';
import 'package:eleve/types/src_maths_questions.dart' as server_questions;
import 'package:eleve/types/src_maths_questions_client.dart';
import 'package:eleve/types/src_maths_repere.dart';
import 'package:eleve/types/src_sql_editor.dart';
import 'package:eleve/types/src_tasks.dart';
import 'package:flutter/material.dart' hide Flow;

// test application for the editor loopback app,
// using local data
void main() async {
  runApp(const _LoopbackTestApp());
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
  1: ExpressionAnswer("4 / (2x)"),
  2: ExpressionAnswer("x^2 + 4 /8 "),
  3: PointAnswer(IntCoord(3, 8)),
};

Question numberQuestion(String title) {
  return Question([
    TextBlock([T(title)], false, false, false),
    const NumberFieldBlock(0, 10)
  ], [
    TextBlock([T("Une jolie correction")], true, true, false),
    TextBlock(
        [T("Avec un super lien https://www.google.com")], true, true, false)
  ]);
}

const origin = server_questions.QuestionPage(null, null, null);

class _LoopbackTestApp extends StatelessWidget {
  const _LoopbackTestApp({Key? key}) : super(key: key);

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
                  onPressed: () => showPaused(context),
                  child: const Text("Loopack: Paused")),
              ElevatedButton(
                  onPressed: () => showQuestion(context),
                  child: const Text("Loopack: Question")),
              ElevatedButton(
                  onPressed: () => showExercice(context),
                  child: const Text("Loopack: Exercice")),
              ElevatedButton(
                  onPressed: () => showCeintures(context),
                  child: const Text("Loopack: Ceintures")),
            ],
          ),
        ));
  }

  void _showRoute(LoopbackServerEvent event, BuildContext context) async {
    await Navigator.of(context).push(MaterialPageRoute<void>(
      settings: const RouteSettings(name: "rootLoopback"),
      builder: (context) =>
          EditorLoopback(event, _LoopbackAPI(), rootRoute: "rootLoopback"),
    ));
  }

  void showPaused(BuildContext context) async {
    return _showRoute(const LoopbackPaused(), context);
  }

  void showQuestion(BuildContext context) async {
    return _showRoute(LoopbackShowQuestion(qu1, [], true, origin), context);
  }

  void showExercice(BuildContext context) async {
    return _showRoute(
        LoopbackShowExercice(workSequencial, ProgressionExt([[], [], []], 0),
            true, [origin, origin, origin]),
        context);
  }

  void showCeintures(BuildContext context) async {
    return _showRoute(
        LoopbackShowCeinture(
            [
              InstantiatedBeltQuestion(1, qu1, []),
              InstantiatedBeltQuestion(2, qu1, []),
              InstantiatedBeltQuestion(3, qu1, []),
            ],
            0,
            [origin, origin, origin],
            false),
        context);
  }
}

// final workParallel = StudentWork(
//   InstantiatedWork(const WorkID(0, true), "Identités remarquables (parallèle)",
//       Flow.parallel, [quI1, quI2, quI3], [1, 1, 2]),
//   ProgressionExt([[], [], []], 0),
// );

final workSequencial = InstantiatedWork(
    const WorkID(0, WorkKind.workExercice, true),
    "Identités remarquables (séquentiel)",
    Flow.sequencial,
    [quI1, quI2, quI3],
    [1, 1, 2]);

class _LoopbackAPI implements LoopbackAPI {
  _LoopbackAPI();

  @override
  Future<EvaluateWorkOut> evaluate(EvaluateWorkIn params) async {
    final questionIndex = params.answerIndex;
    final answer = params.answer;
    final isCorrect =
        questionIndex == (answer.answer.data[0] as NumberAnswer).value;
    params.progression.questions[questionIndex].add(isCorrect);
    return EvaluateWorkOut(
        ProgressionExt(
            params.progression.questions,
            isCorrect
                ? (questionIndex == 2 ? -1 : questionIndex + 1)
                : questionIndex),
        [quI1bis, quI2bis, quI3bis],
        questionIndex,
        QuestionAnswersOut({0: isCorrect}, {}));
  }

  @override
  Future<LoopbackEvaluateQuestionOut> evaluateQuestionAnswer(
      QuestionAnswersIn data, LoopbackShowQuestion origin) async {
    await Future<void>.delayed(const Duration(milliseconds: 200));

    return LoopbackEvaluateQuestionOut(QuestionAnswersOut(
        {0: (data.data[0] as NumberAnswer).value == qu1Answer[0]!.value},
        qu1Answer));
  }

  @override
  Future<LoopbackShowQuestionAnswerOut> showQuestionAnswer(
      server_questions.QuestionPage originPage, Params originParams) async {
    await Future<void>.delayed(const Duration(milliseconds: 200));
    return const LoopbackShowQuestionAnswerOut(QuestionAnswersIn(qu1Answer));
  }

  @override
  Future<LoopbackEvaluateCeintureOut> evaluateCeinture(
      LoopbackEvaluateCeintureIn args) async {
    await Future<void>.delayed(const Duration(milliseconds: 100));
    return LoopbackEvaluateCeintureOut(args.answers.map((an) {
      final isCorrect = 0 == (an.answer.data[0] as NumberAnswer).value;
      return QuestionAnswersOut({0: isCorrect}, {});
    }).toList());
  }
}
