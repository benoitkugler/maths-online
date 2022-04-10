import 'dart:async';
import 'dart:convert';
import 'dart:js' as js;
import 'dart:math';

import 'package:eleve/build_mode.dart';
import 'package:eleve/exercices/question.dart';
import 'package:eleve/exercices/types.gen.dart';
import 'package:eleve/main.dart';
import 'package:flutter/material.dart';
import 'package:web_socket_channel/web_socket_channel.dart';

const Color darkBlue = Color.fromARGB(255, 27, 54, 82);

final bm = buildMode();
// final bm = BuildMode.dev;

void main() {
  // the static app is called via an url setting the session ID
  // note that the MaterialApp routing erase these parameters,
  // so that we need to fetch it early
  final uri = Uri.parse(js.context['location']['href'] as String);
  final id = uri.queryParameters["sessionID"]!;

  runApp(LoopbackApp(id));
}

/// [LoopbackApp] show the content of a question instance
/// being edited, as it will be displayed to the student
/// It is meant to be embedded in a Web page, not used as a mobile app.
class LoopbackApp extends StatelessWidget {
  final String sessionID;

  const LoopbackApp(this.sessionID, {Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Isiro',
      theme: theme,
      debugShowCheckedModeBanner: false,
      localizationsDelegates: localizations,
      supportedLocales: locales,
      home: Scaffold(
        body: _QuestionLoopback(sessionID),
      ),
    );
  }
}

class _QuestionLoopback extends StatefulWidget {
  final String sessionID;

  const _QuestionLoopback(this.sessionID, {Key? key}) : super(key: key);

  @override
  State<_QuestionLoopback> createState() => _QuestionLoopbackState();
}

class _QuestionLoopbackState extends State<_QuestionLoopback> {
  late WebSocketChannel channel;
  late Timer _keepAliveTimmer;

  Question? question;

  @override
  void initState() {
    final url = bm.websocketURL("/prof-loopback/${widget.sessionID}");

    // API connection
    channel = WebSocketChannel.connect(Uri.parse(url));
    channel.stream.listen(listen, onError: showError);
    print("connected !");

    // websocket is closed in case of inactivity
    // prevent it by sending pings
    _keepAliveTimmer = Timer.periodic(const Duration(seconds: 50), (timer) {
      channel.sink.add("Ping");
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

  void listen(dynamic event) {
    try {
      setState(() {
        question = questionFromJson(jsonDecode(event as String));
      });
    } catch (e) {
      showError(e);
    }
  }

  @override
  Widget build(BuildContext context) {
    return question == null
        ? Center(
            child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: const [
              CircularProgressIndicator(),
              Text("En attente de prévisualisation..."),
            ],
          ))
        : Padding(
            padding: const EdgeInsets.all(8.0),
            child: QuestionPage(
              question!,
              Color.fromARGB(255, Random().nextInt(256), Random().nextInt(256),
                  Random().nextInt(256)),
              (p0) {},
              (p0) {},
              showTimeout: false,
            ),
          );
  }
}
