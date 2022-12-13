import 'dart:async';
import 'dart:convert';
import 'dart:js' as js;
import 'dart:html' as html;

import 'package:eleve/build_mode.dart';
import 'package:eleve/exercice/exercice.dart';
import 'package:eleve/loopback/question.dart';
import 'package:eleve/loopback/types.gen.dart';
import 'package:eleve/main_shared.dart';
import 'package:eleve/questions/fields.dart';
import 'package:eleve/questions/types.gen.dart';
import 'package:eleve/shared/errors.dart';
import 'package:eleve/shared_gen.dart' hide Answer;
import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;

void main() {
  // the static app is called via an url setting the session ID
  // note that the MaterialApp routing erase these parameters,
  // so that we need to fetch it early
  final uri = Uri.parse(js.context['location']['href'] as String);
  // final id = uri.queryParameters["sessionID"]!;
  final mode = uri.queryParameters["mode"];
  final bm = APISetting.fromString(mode ?? "");

  runApp(LoopbackApp(bm));
}

/// [LoopbackApp] show the content of a question or an exercice instance
/// being edited, as it will be displayed to the student
/// It is meant to be embedded in a Web page, not used as a mobile app.
class LoopbackApp extends StatelessWidget {
  final BuildMode buildMode;

  const LoopbackApp(this.buildMode, {Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Isyro',
      theme: theme,
      debugShowCheckedModeBanner: false,
      localizationsDelegates: localizations,
      supportedLocales: locales,
      home: _EditorLoopback(buildMode),
    );
  }
}

enum _Mode { paused, question, exercice }

/// owns the websocket connection and switch between question
/// or exercice mode
class _EditorLoopback extends StatefulWidget {
  final BuildMode buildMode;

  const _EditorLoopback(this.buildMode, {Key? key}) : super(key: key);

  @override
  State<_EditorLoopback> createState() => _EditorLoopbackState();
}

class _EditorLoopbackState extends State<_EditorLoopback> {
  late StreamSubscription<html.MessageEvent> subs;

  _Mode get mode => questionData != null
      ? _Mode.question
      : (exerciceData != null ? _Mode.exercice : _Mode.paused);
  LoopackQuestionController? questionData;
  ExerciceController? exerciceData;

  @override
  void initState() {
    subs = html.window.onMessage.listen((event) {
      listen(event.data as String);
    });
    super.initState();
  }

  @override
  void dispose() {
    subs.cancel();
    super.dispose();
  }

  // void _send(LoopbackClientEvent event) {
  //   channel.sink.add(jsonEncode(loopbackClientEventToJson(event)));
  // }

  void _showError(dynamic error) {
    showError("Une erreur est survenue ", error, context);
  }

  void _onServerValidAnswer(QuestionAnswersOut rep) {
    LoopbackQuestionW.showServerValidation(rep, context);
    setState(() {
      questionData!.setFeedback(rep.results);
    });
  }

  void listen(String jsonEvent) {
    final LoopbackServerEvent event;
    try {
      event = loopbackServerEventFromJson(jsonDecode(jsonEvent));
    } catch (e) {
      _showError(e);
      return;
    }

    final api = ServerFieldAPI(widget.buildMode);
    if (event is LoopbackPaused) {
      setState(() {
        questionData = null;
        exerciceData = null;
      });
    } else if (event is LoopbackShowQuestion) {
      final qu = event.question;
      setState(() {
        questionData =
            LoopackQuestionController(qu, api, evaluateQuestionAnswer);
      });
      // } else if (event is LoopbackQuestionValidOut) {
      //   _onServerValidAnswer(event.answers);
      // } else if (event is LoopbackQuestionCorrectAnswersOut) {
      //   final ans = event.answers;
      //   setState(() {
      //     questionData!.setAnswers(ans.data);
      //   });
    } else if (event is LoopbackShowExercice) {
      final ex = StudentWork(event.exercice, event.progression);
      setState(() {
        exerciceData = ExerciceController(ex, null, api);
      });
      // } else if (event is LoopbackExerciceCorrectAnswersOut) {
      //   final ans = event.answers;
      //   setState(() {
      //     exerciceData!.setQuestionAnswers(ans.data);
      //   });
    }
  }

  void evaluateQuestionAnswer(QuestionAnswersIn data) {
    // _send(LoopbackQuestionValidIn(data));
  }

  void _showCorrectAnswer() {
    // if (mode == _Mode.question) {
    //   _send(const LoopbackQuestionCorrectAnswersIn());
    // } else if (mode == _Mode.exercice) {
    //   _send(LoopbackExerciceCorrectAnswsersIn(exerciceData!.questionIndex!));
    // }
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
        return ExerciceW(_LoopbackServerAPI(widget.buildMode), exerciceData!,
            onShowCorrectAnswer: _showCorrectAnswer);
    }
  }
}

class _LoopbackServerAPI implements ExerciceAPI {
  final BuildMode buildMode;
  const _LoopbackServerAPI(this.buildMode);

  @override
  Future<CheckExpressionOut> checkExpressionSyntax(String expression) {
    return ServerFieldAPI(buildMode).checkExpressionSyntax(expression);
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
