import 'package:eleve/activities/trivialpoursuit/question.dart';
import 'package:eleve/decrassage/decrassage.dart';
import 'package:eleve/exercice/exercice.dart';
import 'package:eleve/loopback/question.dart';
import 'package:eleve/main_shared.dart';
import 'package:eleve/questions/debug.dart';
import 'package:eleve/types/src.dart';
import 'package:eleve/types/src_prof_preview.dart';
import 'package:eleve/types/src_maths_questions.dart' as server_questions;
import 'package:eleve/types/src_maths_questions_client.dart';
import 'package:eleve/types/src_maths_repere.dart';
import 'package:eleve/types/src_sql_editor.dart';
import 'package:eleve/types/src_tasks.dart';
import 'package:eleve/types/src_trivial.dart';
import 'package:flutter/material.dart' hide Flow;

void main() async {
  runApp(const _QuestionTestApp());
}

final questionComplexe = Question([
  TextBlock([T("A question to exercice many fields")], false, false, false),
  const NumberFieldBlock(0, 10),
  const ExpressionFieldBlock("x=", "", 10, true, 1),
  const ExpressionFieldBlock("", " = 0", 10, false, 2),
  TreeBlock(
      [
        [T("A")],
        [T("B")],
        [T("C")]
      ],
      const TreeNodeAnswer([
        TreeNodeAnswer([
          TreeNodeAnswer([], [], 0),
          TreeNodeAnswer([], [], 0),
        ], [
          "",
          "7"
        ], 1),
        TreeNodeAnswer([
          TreeNodeAnswer([], [], 0),
          TreeNodeAnswer([], [], 0),
        ], [
          "0.1",
          "x"
        ], 0),
      ], [
        "x",
        "1/3"
      ], 0)),
  TreeFieldBlock([
    [2, 2],
    [2, 3],
  ], [
    [T("A")],
    [T("B")],
    [T("C")]
  ], 10),
  const GeometricConstructionFieldBlock(
      3,
      GFPoint(),
      FunctionsGraphBlock([], [], [],
          [FunctionPoint("#FF0000", "A point", Coord(5, 5))], bounds)),
  const GeometricConstructionFieldBlock(4, GFPoint(),
      FigureBlock(Figure(Drawings({}, [], [], [], []), bounds, true, true))),
  const GeometricConstructionFieldBlock(5, GFVector("test", true),
      FigureBlock(Figure(Drawings({}, [], [], [], []), bounds, true, true))),
  const GeometricConstructionFieldBlock(6, GFVectorPair(),
      FigureBlock(Figure(Drawings({}, [], [], [], []), bounds, true, true))),
], []);

const questionComplexeAnswers = {
  0: NumberAnswer(11.5),
  1: ExpressionAnswer("4 / (2x)"),
  2: ExpressionAnswer("x^2 + 4 /8 "),
  3: PointAnswer(IntCoord(3, 8)),
  4: PointAnswer(IntCoord(3, 8)),
  5: DoublePointAnswer(IntCoord(3, 8), IntCoord(3, 8)),
  6: DoublePointPairAnswer(
      IntCoord(3, 8), IntCoord(3, 8), IntCoord(3, 8), IntCoord(3, 8)),
  10: TreeAnswer(TreeNodeAnswer([], [], 0))
};

Question numberQuestion(String title, {bool withCorrection = true}) {
  return Question(
      [
        TextBlock([T(title)], false, false, false),
        const NumberFieldBlock(0, 10)
      ],
      withCorrection
          ? [
              TextBlock(
                  List.filled(40, T("Une très belle correction : $title")),
                  false,
                  true,
                  true)
            ]
          : []);
}

final qu1 = numberQuestion("Test 1");
final qu2 = numberQuestion("Test 2", withCorrection: false);
final qu3 = numberQuestion("Test 3");

final quI1 = InstantiatedQuestion(1, qu1, DifficultyTag.diff1, []);
final quI2 = InstantiatedQuestion(2, qu2, DifficultyTag.diff2, []);
final quI3 = InstantiatedQuestion(3, qu3, DifficultyTag.diffEmpty, []);

final quI1bis = InstantiatedQuestion(
    1, numberQuestion("Variante 1"), DifficultyTag.diff3, []);
final quI2bis = InstantiatedQuestion(
    2, numberQuestion("Variante 2"), DifficultyTag.diff2, []);
final quI3bis = InstantiatedQuestion(
    3, numberQuestion("Variante 3"), DifficultyTag.diff2, []);

const qu1Answer = {0: NumberAnswer(0)};
const qu2Answer = {0: NumberAnswer(1)};
const qu3Answer = {0: NumberAnswer(2)};

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
                  onPressed: () => showLoopackExerciceSequencial(context),
                  child: const Text("Loopack: Exercice Sequencial")),
              ElevatedButton(
                  onPressed: () => showLoopackExerciceParallel(context),
                  child: const Text("Loopack: Exercice Parallel")),
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
        builder: (context) =>
            _LoopbackQuestion(() => Navigator.of(context).pop())));
  }

  void showExerciceSequencial(BuildContext context) async {
    await Navigator.of(context).push(MaterialPageRoute<void>(
        builder: (context) => const _ExerciceSequential()));
  }

  void showExerciceParallel(BuildContext context) async {
    await Navigator.of(context).push(MaterialPageRoute<void>(
        builder: (context) => const _ExerciceParallel()));
  }

  void showLoopackExerciceSequencial(BuildContext context) async {
    await Navigator.of(context).push(MaterialPageRoute<void>(
        builder: (context) => const _LoopbackExerciceSequential()));
  }

  void showLoopackExerciceParallel(BuildContext context) async {
    await Navigator.of(context).push(MaterialPageRoute<void>(
        builder: (context) => const _LoopbackExerciceParallel()));
  }
}

class _TrivialInGame extends StatelessWidget {
  const _TrivialInGame();

  @override
  Widget build(BuildContext context) {
    return InGameQuestionRoute(
        ShowQuestion(60, Categorie.blue, 0, questionComplexe),
        (a) => onValid(a, context));
  }

  void onValid(QuestionAnswersIn answers, BuildContext context) {
    showDialog<void>(
        context: context,
        builder: (context) => Dialog(child: Text("$answers")));
  }
}

class _TrivialLast extends StatelessWidget {
  final void Function() onDone;

  const _TrivialLast(this.onDone);

  @override
  Widget build(BuildContext context) {
    return LastQuestionRoute(
      ShowQuestion(60, Categorie.blue, 0, questionComplexe),
      onDone,
      questionComplexeAnswers,
    );
  }
}

class _DecrassageAPI implements DecrassageAPI {
  @override
  Future<InstantiateQuestionsOut> loadQuestions(List<int> ids) async {
    await Future<void>.delayed(const Duration(seconds: 1));
    return [
      InstantiatedQuestion(1, qu1, DifficultyTag.diff1, []),
      InstantiatedQuestion(2, qu2, DifficultyTag.diff2, []),
      InstantiatedQuestion(3, qu3, DifficultyTag.diffEmpty, []),
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
  const _Decrassage(this.onDone);

  @override
  Widget build(BuildContext context) {
    return Decrassage(_DecrassageAPI(), const [1, 2, 3]);
  }
}

class _LoopbackQuestion extends StatefulWidget {
  final void Function() onDone;

  const _LoopbackQuestion(this.onDone);

  @override
  State<_LoopbackQuestion> createState() => _LoopbackQuestionState();
}

class _LoopbackQuestionState extends State<_LoopbackQuestion> {
  late final LoopackQuestionController controller;

  @override
  void initState() {
    controller = LoopackQuestionController(
        LoopbackShowQuestion(questionComplexe, [], false,
            const server_questions.QuestionPage(null, null, null)),
        onValid);
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

  // micmic an http call
  void onValid(QuestionAnswersIn a) async {
    await Future<void>.delayed(const Duration(milliseconds: 200));

    const rep = {0: true, 1: false, 2: true, 3: true};

    final snack = LoopbackQuestionW.serverValidation(
        const QuestionAnswersOut(rep, {}), () {});
    ScaffoldMessenger.of(context).showSnackBar(snack);
    setState(() {
      controller.setFeedback(rep);
    });
  }
}

final workParallel = StudentWork(
  InstantiatedWork(
      const WorkID(0, WorkKind.workExercice, true),
      "Identités remarquables (parallèle)",
      Flow.parallel,
      [quI1, quI2, quI3],
      [1, 1, 2]),
  ProgressionExt([[], [], []], 0),
);

final workSequencial = StudentWork(
  InstantiatedWork(
      const WorkID(0, WorkKind.workExercice, true),
      "Identités remarquables (séquentiel)",
      Flow.sequencial,
      [quI1, quI2, quI3],
      [1, 1, 2]),
  ProgressionExt([[], [], []], 0),
);

class _ExerciceSequentialAPI implements ExerciceAPI {
  _ExerciceSequentialAPI();

  @override
  Future<EvaluateWorkOut> evaluate(EvaluateWorkIn params) async {
    final questionIndex = params.answerIndex;
    final answer = params.answer;
    // the correct answer is the index of the question
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
}

class _ExerciceSequential extends StatelessWidget {
  const _ExerciceSequential();

  @override
  Widget build(BuildContext context) {
    return ExerciceW(
        _ExerciceSequentialAPI(), ExerciceController(workSequencial, null));
  }
}

class _LoopbackExerciceSequential extends StatefulWidget {
  const _LoopbackExerciceSequential();

  @override
  State<_LoopbackExerciceSequential> createState() =>
      _LoopbackExerciceSequentialState();
}

class _LoopbackExerciceSequentialState
    extends State<_LoopbackExerciceSequential> {
  ExerciceController ct = ExerciceController(workSequencial, null);

  @override
  Widget build(BuildContext context) {
    return ExerciceW(_ExerciceSequentialAPI(), ct,
        onShowCorrectAnswer: onShowCorrectAnswer);
  }

  void onShowCorrectAnswer() async {
    // mimic server send and receive
    await Future<void>.delayed(const Duration(milliseconds: 200));
    final ans = {0: qu1Answer, 1: qu2Answer, 2: qu3Answer}[ct.questionIndex!]!;
    setState(() {
      ct.setQuestionAnswers(ans);
    });
  }
}

class _ExerciceParallelAPI implements ExerciceAPI {
  _ExerciceParallelAPI();

  @override
  Future<EvaluateWorkOut> evaluate(EvaluateWorkIn params) async {
    final questionIndex = params.answerIndex;
    final answer = params.answer.answer;
    final isCorrect = questionIndex == (answer.data[0] as NumberAnswer).value;
    params.progression.questions[questionIndex].add(isCorrect);
    return EvaluateWorkOut(
        ProgressionExt(
            params.progression.questions,
            params.progression.questions
                .indexWhere((l) => l.every((sucess) => !sucess))),
        [quI1bis, quI2bis, quI3bis],
        questionIndex,
        QuestionAnswersOut({0: isCorrect}, {}));
  }
}

class _ExerciceParallel extends StatelessWidget {
  const _ExerciceParallel();

  @override
  Widget build(BuildContext context) {
    return ExerciceW(
        _ExerciceParallelAPI(), ExerciceController(workParallel, null));
  }
}

class _LoopbackExerciceParallel extends StatefulWidget {
  const _LoopbackExerciceParallel();

  @override
  State<_LoopbackExerciceParallel> createState() =>
      _LoopbackExerciceParallelState();
}

class _LoopbackExerciceParallelState extends State<_LoopbackExerciceParallel> {
  ExerciceController ct = ExerciceController(workParallel, null);

  @override
  Widget build(BuildContext context) {
    return ExerciceW(
      _ExerciceParallelAPI(),
      ct,
      onShowCorrectAnswer: onShowCorrectAnswer,
    );
  }

  void onShowCorrectAnswer() async {
    // mimic server send and receive
    await Future<void>.delayed(const Duration(milliseconds: 200));
    final ans = {0: qu1Answer, 1: qu2Answer, 2: qu3Answer}[ct.questionIndex!]!;
    setState(() {
      ct.setQuestionAnswers(ans);
    });
  }
}
