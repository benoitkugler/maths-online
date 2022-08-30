import 'dart:async';
import 'dart:convert';
import 'dart:js' as js;
import 'dart:math';

import 'package:eleve/build_mode.dart';
import 'package:eleve/exercice/exercice.dart';
import 'package:eleve/loopback_types.gen.dart';
import 'package:eleve/main_shared.dart';
import 'package:eleve/questions/fields.dart';
import 'package:eleve/questions/question.dart';
import 'package:eleve/questions/types.gen.dart' hide Answer;
import 'package:eleve/shared/errors.dart';
import 'package:eleve/shared_gen.dart';
import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;
import 'package:web_socket_channel/web_socket_channel.dart';

void main() {
  // the static app is called via an url setting the session ID
  // note that the MaterialApp routing erase these parameters,
  // so that we need to fetch it early
  final uri = Uri.parse(js.context['location']['href'] as String);
  final id = uri.queryParameters["sessionID"]!;
  final mode = uri.queryParameters["mode"];
  final bm = APISetting.fromString(mode ?? "");

  runApp(LoopbackApp(id, bm));
}

/// [LoopbackApp] show the content of a question or an exercice instance
/// being edited, as it will be displayed to the student
/// It is meant to be embedded in a Web page, not used as a mobile app.
class LoopbackApp extends StatelessWidget {
  final String sessionID;
  final BuildMode buildMode;

  const LoopbackApp(this.sessionID, this.buildMode, {Key? key})
      : super(key: key);

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Isyro',
      theme: theme,
      debugShowCheckedModeBanner: false,
      localizationsDelegates: localizations,
      supportedLocales: locales,
      home: _EditorLoopback(sessionID, buildMode),
    );
  }
}

enum _Mode { paused, question, exercice }

class _QuestionData {
  final Question question;
  final Answers? serverAnswer;

  const _QuestionData(this.question, this.serverAnswer);
}

class _ExerciceData {
  final StudentWork exercice;
  const _ExerciceData(this.exercice);
}

/// owns the websocket connection and switch between question
/// or exercice mode
class _EditorLoopback extends StatefulWidget {
  final String sessionID;
  final BuildMode buildMode;

  const _EditorLoopback(this.sessionID, this.buildMode, {Key? key})
      : super(key: key);

  @override
  State<_EditorLoopback> createState() => _EditorLoopbackState();
}

class _EditorLoopbackState extends State<_EditorLoopback> {
  late WebSocketChannel channel;
  late Timer _keepAliveTimmer;

  _Mode mode = _Mode.paused;
  _QuestionData? questionData;
  _ExerciceData? exerciceData;

  @override
  void initState() {
    final url =
        widget.buildMode.websocketURL("/prof-loopback/${widget.sessionID}");

    // API connection
    channel = WebSocketChannel.connect(Uri.parse(url));
    channel.stream.listen(listen, onError: _showError);

    // websocket is closed in case of inactivity
    // prevent it by sending pings
    _keepAliveTimmer = Timer.periodic(const Duration(seconds: 50), (_) {
      _send(const LoopbackPing());
    });

    super.initState();
  }

  @override
  void dispose() {
    channel.sink.close(1000, "Bye bye");
    _keepAliveTimmer.cancel();
    super.dispose();
  }

  void _send(LoopbackClientEvent event) {
    channel.sink.add(jsonEncode(loopbackClientEventToJson(event)));
  }

  void _showError(dynamic error) {
    showError("Une erreur est survenue ", error, context);
  }

  void _onServerValidAnswer(QuestionAnswersOut rep) {
    final crible = rep.results;
    final errors = crible.values.where((value) => !value).toList();
    final isValid = errors.isEmpty;
    final errorMessage = errors.length >= 2
        ? "${errors.length} champs sont incorrects"
        : "Un champ est incorrect";
    ScaffoldMessenger.of(context).showSnackBar(SnackBar(
      backgroundColor: isValid ? Colors.lightGreen : Colors.red,
      content: Text(isValid ? "Bonne réponse" : errorMessage),
    ));
  }

  void listen(dynamic data) {
    final LoopbackServerEvent event;
    try {
      event = loopbackServerEventFromJson(jsonDecode(data as String));
    } catch (e) {
      _showError(e);
      return;
    }

    if (event is LoopbackPaused) {
      setState(() {
        mode = _Mode.paused;
      });
    } else if (event is LoopbackQuestion) {
      final qu = event.question;
      setState(() {
        mode = _Mode.question;
        questionData = _QuestionData(qu, null);
      });
    } else if (event is LoopbackQuestionValidOut) {
      _onServerValidAnswer(event.answers);
    } else if (event is LoopbackQuestionCorrectAnswersOut) {
      final ans = event.answers;
      setState(() {
        questionData = _QuestionData(questionData!.question, ans.data);
      });
    } else if (event is LoopbackShowExercice) {
      final ex = StudentWork(event.exercice, event.progression);
      setState(() {
        mode = _Mode.exercice;
        exerciceData = _ExerciceData(ex);
      });
    }
  }

  void evaluateQuestionAnswer(QuestionAnswersIn data) {
    _send(LoopbackQuestionValidIn(data));
  }

  void _showCorrectAnswer() {
    if (mode == _Mode.question) {
      _send(const LoopbackQuestionCorrectAnswersIn());
    } else if (mode == _Mode.exercice) {
      // TODO:
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
        return Scaffold(
          appBar: AppBar(actions: [
            TextButton(
                onPressed: _showCorrectAnswer,
                child: const Text("Afficher la réponse."))
          ]),
          body: _QuestionLoopback(
            widget.buildMode,
            questionData!.question,
            questionData!.serverAnswer,
            evaluateQuestionAnswer,
          ),
        );
      case _Mode.exercice:
        return _ExerciceLoopback(
            _LoopbackServerAPI(widget.buildMode), exerciceData!.exercice);
    }
  }
}

class _QuestionLoopback extends StatelessWidget {
  final BuildMode buildMode;
  final Question question;
  final Answers? answers;
  final void Function(QuestionAnswersIn) onValid;

  const _QuestionLoopback(
      this.buildMode, this.question, this.answers, this.onValid,
      {Key? key})
      : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 8.0),
      child: QuestionW(
        ServerFieldAPI(buildMode),
        question,
        Color.fromARGB(255, 150 + Random().nextInt(100),
            150 + Random().nextInt(100), Random().nextInt(256)),
        onValid,
        timeout: null,
        blockOnSubmit: false,
        answer: answers,
      ),
    );
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

class _ExerciceLoopback extends StatelessWidget {
  final _LoopbackServerAPI api;
  final StudentWork exercice;

  const _ExerciceLoopback(this.api, this.exercice, {Key? key})
      : super(key: key);

  @override
  Widget build(BuildContext context) {
    return ExerciceW(api, exercice);
  }
}
