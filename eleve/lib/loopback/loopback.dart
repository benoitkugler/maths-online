import 'dart:convert';
import 'package:http/http.dart' as http;

import 'package:eleve/build_mode.dart';
import 'package:eleve/exercice/exercice.dart';
import 'package:eleve/loopback/exercice.dart';
import 'package:eleve/loopback/question.dart';
import 'package:eleve/questions/fields.dart';
import 'package:eleve/shared/errors.dart';
import 'package:eleve/types/src_maths_questions.dart';
import 'package:eleve/types/src_maths_questions_client.dart';
import 'package:eleve/types/src_prof_editor.dart';
import 'package:eleve/types/src_tasks.dart';
import 'package:flutter/material.dart';

enum _Mode { paused, question, exercice }

/// [EditorLoopback] switch between pause, question or exercice mode according
/// to the [event] parameter.
class EditorLoopback extends StatefulWidget {
  final LoopbackServerEvent event;
  final LoopbackAPI api;

  const EditorLoopback(this.event, this.api, {Key? key}) : super(key: key);

  @override
  State<EditorLoopback> createState() => _EditorLoopbackState();
}

class _EditorLoopbackState extends State<EditorLoopback> {
  _Mode get mode => questionData != null
      ? _Mode.question
      : (exerciceData != null ? _Mode.exercice : _Mode.paused);

  LoopackQuestionController? questionData;
  LoopbackExerciceController? exerciceData;

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
      questionData = null;
      exerciceData = null;
    } else if (event is LoopbackShowQuestion) {
      questionData =
          LoopackQuestionController(event, widget.api, evaluateQuestionAnswer);
      exerciceData = null;
    } else if (event is LoopbackShowExercice) {
      questionData = null;
      exerciceData = LoopbackExerciceController(event, widget.api);
    }
  }

  void _showError(dynamic error) {
    showError("Une erreur est survenue ", error, context);
  }

  void evaluateQuestionAnswer(QuestionAnswersIn data) async {
    try {
      final res =
          await widget.api.evaluateQuestionAnswer(data, questionData!.data);
      if (!mounted) {
        return;
      }
      LoopbackQuestionW.showServerValidation(res.answers, context);
      setState(() {
        questionData!.setFeedback(res.answers.results);
      });
    } catch (e) {
      _showError(e);
    }
  }

  void _showCorrectAnswer() async {
    final QuestionPage originPage;
    final Params originParams;
    switch (mode) {
      case _Mode.question:
        originPage = questionData!.data.origin;
        originParams = questionData!.data.params;
        break;
      case _Mode.exercice:
        final index = exerciceData!.controller.questionIndex!;
        originPage = exerciceData!.data.origin[index];
        originParams = exerciceData!
            .controller.exeAndProg.exercice.questions[index].params;
        break;
      case _Mode.paused:
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

    switch (mode) {
      case _Mode.question:
        setState(() {
          questionData!.setAnswers(ans.data);
        });
        return;
      case _Mode.exercice:
        setState(() {
          exerciceData!.controller.setQuestionAnswers(ans.data);
        });
        return;
      case _Mode.paused:
        return;
    }
  }

  @override
  Widget build(BuildContext context) {
    switch (mode) {
      case _Mode.paused:
        return Scaffold(
          body: Center(
              child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: const [
              CircularProgressIndicator(),
              Padding(
                padding: EdgeInsets.symmetric(vertical: 8.0),
                child: Text("En attente de prévisualisation..."),
              ),
            ],
          )),
        );
      case _Mode.question:
        return LoopbackQuestionW(questionData!, _showCorrectAnswer);
      case _Mode.exercice:
        return ExerciceW(widget.api, exerciceData!.controller,
            onShowCorrectAnswer: _showCorrectAnswer);
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
}

/// [LoopbackServerAPI] is the defaut implementation for [LoopbackAPI]
/// using a backend and http calls.
class LoopbackServerAPI extends ServerFieldAPI implements LoopbackAPI {
  const LoopbackServerAPI(super.buildMode);

  @override
  Future<LoopbackEvaluateQuestionOut> evaluateQuestionAnswer(
      QuestionAnswersIn data, LoopbackShowQuestion origin) async {
    final uri =
        Uri.parse(buildMode.serverURL("/api/loopack/evaluate-question"));
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
    final uri = Uri.parse(buildMode.serverURL("/api/loopack/question-answer"));
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
    final uri = Uri.parse(buildMode.serverURL("/api/exercices/evaluate"));
    final resp = await http
        .post(uri, body: jsonEncode(evaluateWorkInToJson(params)), headers: {
      'Content-type': 'application/json',
    });
    return evaluateWorkOutFromJson(jsonDecode(resp.body));
  }
}
