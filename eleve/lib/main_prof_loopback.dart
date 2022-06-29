import 'dart:async';
import 'dart:convert';
import 'dart:js' as js;
import 'dart:math';

import 'package:eleve/build_mode.dart';
import 'package:eleve/loopback_types.gen.dart';
import 'package:eleve/main_shared.dart';
import 'package:eleve/questions/question.dart';
import 'package:eleve/questions/types.gen.dart';
import 'package:flutter/material.dart';
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

/// [LoopbackApp] show the content of a question instance
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
      home: _QuestionLoopback(sessionID, buildMode),
    );
  }
}

class _QuestionLoopback extends StatefulWidget {
  final String sessionID;
  final BuildMode buildMode;

  const _QuestionLoopback(this.sessionID, this.buildMode, {Key? key})
      : super(key: key);

  @override
  State<_QuestionLoopback> createState() => _QuestionLoopbackState();
}

class _QuestionLoopbackState extends State<_QuestionLoopback> {
  late WebSocketChannel channel;
  late Timer _keepAliveTimmer;

  LoopbackState? question;
  QuestionAnswersIn? serverAnswer;

  @override
  void initState() {
    final url =
        widget.buildMode.websocketURL("/prof-loopback/${widget.sessionID}");

    // API connection
    channel = WebSocketChannel.connect(Uri.parse(url));
    channel.stream.listen(listen, onError: showError);

    // websocket is closed in case of inactivity
    // prevent it by sending pings
    _keepAliveTimmer = Timer.periodic(const Duration(seconds: 50), (timer) {
      channel.sink.add(jsonEncode(
          {"Kind": loopbackClientDataKindToJson(LoopbackClientDataKind.ping)}));
    });

    super.initState();
  }

  @override
  void dispose() {
    channel.sink.close(1000, "Bye bye");
    _keepAliveTimmer.cancel();
    super.dispose();
  }

  void showError(dynamic error) {
    ScaffoldMessenger.of(context).showSnackBar(SnackBar(
      duration: const Duration(seconds: 5),
      backgroundColor: Theme.of(context).colorScheme.error,
      content: Text("Une erreur est survenue : $error"),
    ));
  }

  void _onServerState(LoopbackState state) {
    setState(() {
      question = state;
      serverAnswer = null;
    });
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

  void _onServerShowCorrectAnswer(QuestionAnswersIn answer) {
    setState(() {
      serverAnswer = answer;
    });
  }

  void listen(dynamic event) {
    try {
      final data = jsonDecode(event as String) as Map<String, dynamic>;
      final kind = loopbackServerDataKindFromJson(data["Kind"]);
      switch (kind) {
        case LoopbackServerDataKind.state:
          return _onServerState(loopbackStateFromJson(data["Data"]));
        case LoopbackServerDataKind.validAnswerOut:
          return _onServerValidAnswer(questionAnswersOutFromJson(data["Data"]));
        case LoopbackServerDataKind.showCorrectAnswerOut:
          return _onServerShowCorrectAnswer(
              questionAnswersInFromJson(data["Data"]));
      }
    } catch (e) {
      showError(e);
    }
  }

  void _validAnswer(QuestionAnswersIn data) {
    channel.sink.add(jsonEncode({
      "Kind":
          loopbackClientDataKindToJson(LoopbackClientDataKind.validAnswerIn),
      "Data": questionAnswersInToJson(data)
    }));
  }

  void _showCorrectAnswer() {
    channel.sink.add(jsonEncode({
      "Kind": loopbackClientDataKindToJson(
          LoopbackClientDataKind.showCorrectAnswerIn),
      "Data": null,
    }));
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(actions: [
        TextButton(
            onPressed: _showCorrectAnswer,
            child: const Text("Afficher la réponse."))
      ]),
      body: question == null || question!.isPaused
          ? Center(
              child: Column(
              mainAxisAlignment: MainAxisAlignment.center,
              children: const [
                CircularProgressIndicator(),
                Padding(
                  padding: EdgeInsets.symmetric(vertical: 8.0),
                  child: Text("En attente de prévisualisation..."),
                ),
              ],
            ))
          : Padding(
              padding: const EdgeInsets.all(8.0),
              child: QuestionW(
                widget.buildMode,
                question!.question,
                Color.fromARGB(255, 150 + Random().nextInt(100),
                    150 + Random().nextInt(100), Random().nextInt(256)),
                _validAnswer,
                timeout: null,
                blockOnSubmit: false,
                answer: serverAnswer?.data,
              ),
            ),
    );
  }
}
