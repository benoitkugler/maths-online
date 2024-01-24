import 'dart:convert';
import 'package:eleve/activities/ceintures/seance.dart';
import 'package:eleve/exercice/home.dart';
import 'package:eleve/loopback/ceintures.dart';
import 'package:eleve/questions/question.dart';
import 'package:eleve/quotes.dart';
import 'package:http/http.dart' as http;

import 'package:eleve/build_mode.dart';
import 'package:eleve/exercice/exercice.dart';
import 'package:eleve/loopback/exercice.dart';
import 'package:eleve/loopback/question.dart';
import 'package:eleve/shared/errors.dart';
import 'package:eleve/types/src_maths_questions.dart';
import 'package:eleve/types/src_maths_questions_client.dart';
import 'package:eleve/types/src_tasks.dart';
import 'package:eleve/types/src_prof_preview.dart';
import 'package:flutter/material.dart';

/// [EditorLoopback] switch between pause, question or exercice mode according
/// to the [event] parameter.
class EditorLoopback extends StatefulWidget {
  final LoopbackServerEvent event;
  final LoopbackAPI api;

  final String? rootRoute;

  const EditorLoopback(this.event, this.api, {this.rootRoute, Key? key})
      : super(key: key);

  @override
  State<EditorLoopback> createState() => _EditorLoopbackState();
}

class _EditorLoopbackState extends State<EditorLoopback> {
  LoopbackController? controller;

  @override
  void initState() {
    _initData();
    super.initState();
  }

  @override
  void didUpdateWidget(covariant EditorLoopback oldWidget) {
    _initData();
    super.didUpdateWidget(oldWidget);
  }

  void _initData() {
    final event = widget.event;
    if (event is LoopbackPaused) {
      controller = null;
    } else if (event is LoopbackShowQuestion) {
      controller = LoopackQuestionController(event, evaluateQuestionAnswer);
    } else if (event is LoopbackShowExercice) {
      controller = LoopbackExerciceController(event);
    } else if (event is LoopbackShowCeinture) {
      controller = LoopbackCeinturesController(event);
    }

    WidgetsBinding.instance.addPostFrameCallback((_) => _handleCorrection());
  }

  // properly show correction or hide potential old route
  void _handleCorrection() {
    final controller = this.controller;
    if (controller is LoopackQuestionController) {
      if (controller.data.showCorrection) {
        _showQuestionCorrection();
      } else {
        _popToRoot();
      }
      return;
    } else if (controller is LoopbackExerciceController) {
      if (!controller.data.showCorrection) _popToRoot();
      return;
    } else if (controller is LoopbackCeinturesController) {
      // TODO:
      return _popToRoot();
    } else {
      return _popToRoot();
    }
  }

  void _popToRoot() {
    Navigator.popUntil(
        context,
        (Route<dynamic> route) => widget.rootRoute == null
            ? route.isFirst
            : route.settings.name == widget.rootRoute);
  }

  void _showError(dynamic error) {
    showError("Une erreur est survenue ", error, context);
  }

  void evaluateQuestionAnswer(QuestionAnswersIn data) async {
    final ct = controller as LoopackQuestionController;
    try {
      final res = await widget.api.evaluateQuestionAnswer(data, ct.data);
      if (!mounted) {
        return;
      }
      final snack = LoopbackQuestionW.serverValidation(
          res.answers, _showQuestionCorrection);
      ScaffoldMessenger.of(context).showSnackBar(snack);
      setState(() {
        ct.setFeedback(res.answers.results);
      });
    } catch (e) {
      _showError(e);
    }
  }

  Future<List<QuestionAnswersOut>> _evaluateCeinture(
      SeanceAnswers answers) async {
    final res = await widget.api.evaluateCeinture(
        LoopbackEvaluateCeintureIn(answers.questions, answers.answers));
    if (!mounted) return [];

    if (res.answers.every((element) => element.isCorrect)) {
      ScaffoldMessenger.of(context).showSnackBar(const SnackBar(
        content: Text("Bonnes réponses !"),
        backgroundColor: Colors.lightGreen,
      ));
    }

    return res.answers;
  }

  void _resetCeinture() {
    final c = controller as LoopbackCeinturesController;
    setState(() {
      controller = LoopbackCeinturesController(c.data);
    });
  }

  void _showQuestionCorrection() {
    final controller = this.controller;
    if (controller is! LoopackQuestionController) return;
    Navigator.of(context).push(MaterialPageRoute<void>(
        builder: (context) => Scaffold(
              appBar: AppBar(),
              body: CorrectionW(controller.data.question.correction,
                  randColor(), pickQuote()),
            )));
  }

  void _showCorrectAnswer() async {
    final QuestionPage originPage;
    final Params originParams;

    final controller = this.controller;
    if (controller is LoopackQuestionController) {
      originPage = controller.data.origin;
      originParams = controller.data.params;
    } else if (controller is LoopbackExerciceController) {
      final index = controller.controller.questionIndex!;
      originPage = controller.data.origin[index];
      originParams =
          controller.controller.exeAndProg.exercice.questions[index].params;
    } else if (controller is LoopbackCeinturesController) {
      final index = controller.controller.currentQuestion;
      originPage = controller.data.origin[index];
      originParams = controller.data.questions[index].params;
    } else {
      return;
    }

    final QuestionAnswersIn ans;
    try {
      final res = await widget.api.showQuestionAnswer(originPage, originParams);
      ans = res.answers;
    } catch (e) {
      _showError(e);
      return;
    }

    if (controller is LoopackQuestionController) {
      setState(() {
        controller.setAnswers(ans.data);
      });
    } else if (controller is LoopbackExerciceController) {
      setState(() {
        controller.controller.setQuestionAnswers(ans.data);
        controller.instantShowCorrection = false;
      });
    } else if (controller is LoopbackCeinturesController) {
      setState(() {
        controller.controller.setQuestionAnswers(ans.data);
        controller.instantShowCorrection = false;
      });
    } else {
      return;
    }
  }

  @override
  Widget build(BuildContext context) {
    final controller = this.controller;
    if (controller is LoopackQuestionController) {
      return LoopbackQuestionW(controller, _showCorrectAnswer);
    } else if (controller is LoopbackExerciceController) {
      return ExerciceW(widget.api, controller.controller,
          onShowCorrectAnswer: _showCorrectAnswer,
          showCorrectionButtonOnFail: true,
          instantShowCorrection: controller.instantShowCorrection);
    } else if (controller is LoopbackCeinturesController) {
      return CeintureW(
          controller, _evaluateCeinture, _resetCeinture, _showCorrectAnswer);
    } else {
      return const Scaffold(
        body: Center(
            child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            CircularProgressIndicator(),
            Padding(
              padding: EdgeInsets.symmetric(vertical: 8.0),
              child: Text("En attente de prévisualisation..."),
            ),
          ],
        )),
      );
    }
  }
}

/// [LoopbackAPI] is the interface defining the
/// methods required by the loopback widget.
abstract class LoopbackAPI extends ExerciceAPI {
  Future<LoopbackEvaluateQuestionOut> evaluateQuestionAnswer(
      QuestionAnswersIn data, LoopbackShowQuestion origin);

  Future<LoopbackShowQuestionAnswerOut> showQuestionAnswer(
      QuestionPage originPage, Params originParams);

  Future<LoopbackEvaluateCeintureOut> evaluateCeinture(
      LoopbackEvaluateCeintureIn args);
}

/// [LoopbackServerAPI] is the defaut implementation for [LoopbackAPI]
/// using a backend and http calls.
class LoopbackServerAPI implements LoopbackAPI {
  final BuildMode buildMode;
  const LoopbackServerAPI(this.buildMode);

  @override
  Future<LoopbackEvaluateQuestionOut> evaluateQuestionAnswer(
      QuestionAnswersIn data, LoopbackShowQuestion origin) async {
    final uri = buildMode.serverURL("/api/loopack/evaluate-question");
    final params =
        LoopackEvaluateQuestionIn(origin.origin, AnswerP(origin.params, data));
    final resp = await http.post(uri,
        body: jsonEncode(loopackEvaluateQuestionInToJson(params)),
        headers: {
          'Content-type': 'application/json',
        });
    return loopbackEvaluateQuestionOutFromJson(checkServerError(resp.body));
  }

  @override
  Future<LoopbackShowQuestionAnswerOut> showQuestionAnswer(
      QuestionPage originPage, Params originParams) async {
    final uri = buildMode.serverURL("/api/loopack/question-answer");
    final params = LoopbackShowQuestionAnswerIn(originPage, originParams);
    final resp = await http.post(uri,
        body: jsonEncode(loopbackShowQuestionAnswerInToJson(params)),
        headers: {
          'Content-type': 'application/json',
        });
    return loopbackShowQuestionAnswerOutFromJson(checkServerError(resp.body));
  }

  @override
  Future<EvaluateWorkOut> evaluate(EvaluateWorkIn params) async {
    final uri = buildMode.serverURL("/api/exercices/evaluate");
    final resp = await http
        .post(uri, body: jsonEncode(evaluateWorkInToJson(params)), headers: {
      'Content-type': 'application/json',
    });
    return evaluateWorkOutFromJson(checkServerError(resp.body));
  }

  @override
  Future<LoopbackEvaluateCeintureOut> evaluateCeinture(
      LoopbackEvaluateCeintureIn params) async {
    final uri = buildMode.serverURL("/api/loopback/evaluate-ceinture");
    final resp = await http.post(uri,
        body: jsonEncode(loopbackEvaluateCeintureInToJson(params)),
        headers: {
          'Content-type': 'application/json',
        });
    return loopbackEvaluateCeintureOutFromJson(checkServerError(resp.body));
  }
}
