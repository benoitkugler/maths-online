import 'package:eleve/loopback/loopback.dart';
import 'package:eleve/main_shared.dart';
import 'package:eleve/questions/debug.dart';
import 'package:eleve/types/src.dart';
import 'package:eleve/types/src_maths_questions.dart' as server_questions;
import 'package:eleve/types/src_maths_questions_client.dart';
import 'package:eleve/types/src_maths_repere.dart';
import 'package:eleve/types/src_prof_editor.dart';
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
  const FigurePointFieldBlock(
      Figure(Drawings({}, [], [], [], []), bounds, true, true), 3)
]);

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
  ]);
}

const origin = server_questions.QuestionPage(null, null);

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
            ],
          ),
        ));
  }

  void showPaused(BuildContext context) async {
    await Navigator.of(context).push(MaterialPageRoute<void>(
        builder: (context) =>
            EditorLoopback(const LoopbackPaused(), _LoopbackAPI())));
  }

  void showQuestion(BuildContext context) async {
    await Navigator.of(context).push(MaterialPageRoute<void>(
        builder: (context) => EditorLoopback(
              LoopbackShowQuestion(qu1, [], origin),
              _LoopbackAPI(),
            )));
  }

  void showExercice(BuildContext context) async {
    await Navigator.of(context).push(MaterialPageRoute<void>(
        builder: (context) => EditorLoopback(
            LoopbackShowExercice(workSequencial,
                ProgressionExt([[], [], []], 0), [origin, origin, origin]),
            _LoopbackAPI())));
  }
}

// final workParallel = StudentWork(
//   InstantiatedWork(const WorkID(0, true), "Identités remarquables (parallèle)",
//       Flow.parallel, [quI1, quI2, quI3], [1, 1, 2]),
//   ProgressionExt([[], [], []], 0),
// );

final workSequencial = InstantiatedWork(
    const WorkID(0, true),
    "Identités remarquables (séquentiel)",
    Flow.sequencial,
    [quI1, quI2, quI3],
    [1, 1, 2]);

class _LoopbackAPI implements LoopbackAPI {
  _LoopbackAPI();

  @override
  Future<CheckExpressionOut> checkExpressionSyntax(String expression) async {
    return const CheckExpressionOut("", true);
  }

  @override
  Future<EvaluateWorkOut> evaluate(EvaluateWorkIn params) async {
    final questionIndex = params.answers.keys.first;
    final answer = params.answers[questionIndex]!;
    final isCorrect =
        questionIndex == (answer.answer.data[0] as NumberAnswer).value;
    final res = {
      questionIndex: QuestionAnswersOut({0: isCorrect}, {})
    };
    params.progression.questions[questionIndex].add(isCorrect);
    return EvaluateWorkOut(
        res,
        ProgressionExt(
            params.progression.questions,
            isCorrect
                ? (questionIndex == 2 ? -1 : questionIndex + 1)
                : questionIndex),
        [quI1bis, quI2bis, quI3bis]);
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
}
