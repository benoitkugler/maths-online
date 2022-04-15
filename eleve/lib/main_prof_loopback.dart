import 'dart:async';
import 'dart:convert';
import 'dart:js' as js;
import 'dart:math';

import 'package:eleve/build_mode.dart';
import 'package:eleve/exercices/question.dart';
import 'package:eleve/exercices/types.gen.dart';
import 'package:eleve/loopback_types.gen.dart';
import 'package:eleve/main_shared.dart';
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
      title: 'Isiro',
      theme: theme,
      debugShowCheckedModeBanner: false,
      localizationsDelegates: localizations,
      supportedLocales: locales,
      home: Scaffold(
        body: _QuestionLoopback(sessionID, buildMode),
      ),
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
      channel.sink.add(jsonEncode({"Kind": LoopbackClientDataKind.ping}));
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
    });
  }

  void _onServerCheckSyntax(QuestionSyntaxCheckOut rep) {
    if (rep.isValid) {
      return;
    }
    final reason = rep.reason;
    ScaffoldMessenger.of(context).showSnackBar(SnackBar(
      backgroundColor: Colors.red,
      content: Text.rich(TextSpan(children: [
        const TextSpan(text: "Syntaxe invalide: "),
        TextSpan(
            text: reason, style: const TextStyle(fontWeight: FontWeight.bold)),
      ])),
    ));
  }

  void _onServerValidAnswer(QuestionAnswersOut rep) {
    final crible = rep.data;
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

  void listen(dynamic event) {
    try {
      final data = jsonDecode(event as String) as Map<String, dynamic>;
      final kind = loopbackServerDataKindFromJson(data["Kind"]);
      switch (kind) {
        case LoopbackServerDataKind.state:
          return _onServerState(loopbackStateFromJson(data["Data"]));
        case LoopbackServerDataKind.checkSyntaxeOut:
          return _onServerCheckSyntax(
              questionSyntaxCheckOutFromJson(data["Data"]));
        case LoopbackServerDataKind.validAnswerOut:
          return _onServerValidAnswer(questionAnswersOutFromJson(data["Data"]));
      }
    } catch (e) {
      showError(e);
    }
  }

  void _checkSyntax(CheckQuestionSyntaxeNotification notif) {
    channel.sink.add(jsonEncode({
      "Kind":
          loopbackClientDataKindToJson(LoopbackClientDataKind.checkSyntaxIn),
      "Data": questionSyntaxCheckInToJson(notif.data)
    }));
  }

  void _validAnswer(ValidQuestionNotification notif) {
    channel.sink.add(jsonEncode({
      "Kind":
          loopbackClientDataKindToJson(LoopbackClientDataKind.validAnswerIn),
      "Data": questionAnswersInToJson(notif.data)
    }));
  }

  @override
  Widget build(BuildContext context) {
    return question == null || question!.isPaused
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
            child: QuestionPage(
              question!.question,
              Color.fromARGB(255, Random().nextInt(256), Random().nextInt(256),
                  Random().nextInt(256)),
              _checkSyntax,
              _validAnswer,
              showTimeout: false,
            ),
          );
  }
}
