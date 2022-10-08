import 'package:eleve/activities/trivialpoursuit/events.gen.dart';
import 'package:eleve/activities/trivialpoursuit/question.dart';
import 'package:eleve/decrassage/decrassage.dart';
import 'package:eleve/exercice/exercice.dart';
import 'package:eleve/loopback/question.dart';
import 'package:eleve/main_shared.dart';
import 'package:eleve/questions/debug.dart';
import 'package:eleve/questions/repere.gen.dart';
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

final questionComplexe = Question([
  TextBlock([T("A question to exercice many fields")], false, false, false),
  const NumberFieldBlock(0, 10),
  const ExpressionFieldBlock("x=", 10, 1),
  const FigurePointFieldBlock(
      Figure(Drawings({}, [], [], [], []), bounds, true, true), 2)
]);

const questionComplexeAnswers = {
  0: NumberAnswer(11.5),
  1: ExpressionAnswer("x^2 + 4 /8 "),
  2: PointAnswer(IntCoord(3, 8)),
};

final qu1 = Question([
  TextBlock([T("Test 1")], false, false, false),
  const NumberFieldBlock(0, 10)
]);

final qu2 = Question([
  TextBlock([T("Test 2")], false, false, false),
  const NumberFieldBlock(0, 10)
]);
final qu3 = Question([
  TextBlock([T("Test 3")], false, false, false),
  const NumberFieldBlock(0, 10)
]);

final quI1 = InstantiatedQuestion(0, qu1, []);
final quI2 = InstantiatedQuestion(0, qu2, []);

final exercice = ExerciceController(
    StudentWork(
      InstantiatedWork(
          WorkID(0, true), "", Flow.parallel, [quI1, quI2], [1, 1]),
      ProgressionExt([], 0),
    ),
    null,
    _API());

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
        home: Builder(
          builder: (context) => Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              ElevatedButton(
                  onPressed: () => showTrivialInGame(context),
                  child: const Text("Trivial: InGame")),
              ElevatedButton(
                  onPressed: () => showTrivialLast(context),
                  child: const Text("Trivial: Last")),
              ElevatedButton(
                  onPressed: () => showDecrassage(context),
                  child: const Text("Decrassage")),
              ElevatedButton(
                  onPressed: () => showLoopackQuestion(context),
                  child: const Text("Loopack: Question")),
              ElevatedButton(
                  onPressed: () => showExerciceSequencial(context),
                  child: const Text("Homework: Sequencial")),
              ElevatedButton(
                  onPressed: () => showExerciceParallel(context),
                  child: const Text("Homework: Parallel")),
              ElevatedButton(
                  onPressed: () => showLoopackExercice(context),
                  child: const Text("Loopack: Exercice")),
            ],
          ),
        ));
  }

  void showTrivialInGame(BuildContext context) async {
    await Navigator.of(context).push(
        MaterialPageRoute<void>(builder: (context) => const _TrivialInGame()));
  }

  void showTrivialLast(BuildContext context) async {
    await Navigator.of(context).push(MaterialPageRoute<void>(
        builder: (context) => _TrivialLast(() => Navigator.of(context).pop())));
  }

  void showDecrassage(BuildContext context) async {
    await Navigator.of(context).push(MaterialPageRoute<void>(
        builder: (context) => _Decrassage(() => Navigator.of(context).pop())));
  }

  void showLoopackQuestion(BuildContext context) async {
    await Navigator.of(context).push(MaterialPageRoute<void>(
        builder: (context) => _Loopback(() => Navigator.of(context).pop())));
  }

  void showExerciceSequencial(BuildContext context) async {}

  void showExerciceParallel(BuildContext context) async {}

  void showLoopackExercice(BuildContext context) async {}
}

class _TrivialInGame extends StatelessWidget {
  const _TrivialInGame({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: InGameQuestionRoute(
          _API(),
          ShowQuestion(60, Categorie.blue, 0, questionComplexe),
          (a) => onValid(a, context)),
    );
  }

  void onValid(QuestionAnswersIn answers, BuildContext context) {
    showDialog<void>(
        context: context,
        builder: (context) => Dialog(child: Text("$answers")));
  }
}

class _TrivialLast extends StatelessWidget {
  final void Function() onDone;

  const _TrivialLast(this.onDone, {super.key});

  @override
  Widget build(BuildContext context) {
    return LastQuestionRoute(
      _API(),
      ShowQuestion(60, Categorie.blue, 0, questionComplexe),
      onDone,
      questionComplexeAnswers,
    );
  }
}

class _DecrassageAPI implements DecrassageAPI {
  @override
  Future<CheckExpressionOut> checkExpressionSyntax(String expression) async {
    return const CheckExpressionOut("", true);
  }

  @override
  Future<InstantiateQuestionsOut> loadQuestions(List<int> ids) async {
    await Future<void>.delayed(const Duration(seconds: 1));
    return [
      InstantiatedQuestion(1, qu1, []),
      InstantiatedQuestion(2, qu2, []),
      InstantiatedQuestion(3, qu3, []),
    ];
  }

  @override
  Future<QuestionAnswersOut> evaluateQuestion(EvaluateQuestionIn answer) async {
    return QuestionAnswersOut({
      0: (answer.answer.answer.data[0] as NumberAnswer).value ==
          answer.idQuestion.toDouble()
    }, {
      0: NumberAnswer(answer.idQuestion.toDouble())
    });
  }
}

class _Decrassage extends StatelessWidget {
  final void Function() onDone;
  const _Decrassage(this.onDone, {super.key});

  @override
  Widget build(BuildContext context) {
    return Decrassage(_DecrassageAPI(), const [1, 2, 3]);
  }
}

class _Loopback extends StatefulWidget {
  final void Function() onDone;

  const _Loopback(this.onDone, {super.key});

  @override
  State<_Loopback> createState() => _LoopbackState();
}

class _LoopbackState extends State<_Loopback> {
  late final LoopackQuestionController controller;

  @override
  void initState() {
    controller = LoopackQuestionController(questionComplexe, _API(), onValid);
    super.initState();
  }

  @override
  Widget build(BuildContext context) {
    return LoopbackQuestionW(controller, loadCorrectAnswers);
  }

  void loadCorrectAnswers() {
    setState(() {
      controller.setAnswers(questionComplexeAnswers);
    });
  }

  // micmic a socket send and receive
  void onValid(QuestionAnswersIn a) async {
    await Future<void>.delayed(const Duration(milliseconds: 200));

    const rep = {0: true, 1: false, 2: true};
    LoopbackQuestionW.showServerValidation(
        const QuestionAnswersOut(rep, {}), context);
    setState(() {
      controller.setFeedback(rep);
    });
  }
}
