import 'dart:async';
import 'dart:convert';
import 'dart:js' as js;
import 'dart:html' as html;

import 'package:eleve/activities/trivialpoursuit/categories.dart';
import 'package:eleve/build_mode.dart';
import 'package:eleve/main_shared.dart';
import 'package:eleve/questions/question.dart';
import 'package:eleve/shared/errors.dart';
import 'package:eleve/types/src_maths_questions_client.dart';
import 'package:eleve/types/src_prof_trivial.dart';
import 'package:flutter/material.dart';

void main() {
  // the static app is called via an url setting the session ID
  // note that the MaterialApp routing erase these parameters,
  // so that we need to fetch it early
  final uri = Uri.parse(js.context['location']['href'] as String);
  // final id = uri.queryParameters["sessionID"]!;
  final mode = uri.queryParameters["mode"];
  final bm = APISetting.fromString(mode ?? "");

  runApp(MonitorApp(bm));
}

/// [MonitorApp] show the content of a question (used in a IsyTriv game).
/// It is meant to be embedded in a Web page, not used as a mobile app.
class MonitorApp extends StatelessWidget {
  final BuildMode buildMode;

  const MonitorApp(this.buildMode, {Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Isyro',
      theme: theme,
      debugShowCheckedModeBanner: false,
      localizationsDelegates: localizations,
      supportedLocales: locales,
      home: _Monitor(buildMode),
    );
  }
}

/// listen to an HTML event containing the question content
class _Monitor extends StatefulWidget {
  final BuildMode buildMode;

  const _Monitor(this.buildMode, {Key? key}) : super(key: key);

  @override
  State<_Monitor> createState() => _MonitorState();
}

class _MonitorState extends State<_Monitor> {
  late final StreamSubscription<html.MessageEvent> subs;

  QuestionContent? event;

  @override
  void initState() {
    subs = html.window.onMessage.listen((event) {
      listen(event.data as String);
    });
    html.window.parent?.postMessage(jsonEncode({"PREVIEW_READY": true}), "*");
    super.initState();
  }

  @override
  void dispose() {
    subs.cancel();
    super.dispose();
  }

  void _showError(dynamic error) {
    showError("Une erreur est survenue ", error, context);
  }

  void listen(String jsonEvent) {
    try {
      setState(() {
        event = questionContentFromJson(jsonDecode(jsonEvent));
      });
    } catch (e) {
      _showError("$e ($jsonEvent)");
      return;
    }
  }

  @override
  Widget build(BuildContext context) {
    final content = event;
    return Scaffold(
      appBar: AppBar(
        title: const Text("Question en cours"),
      ),
      body: content == null
          ? const Center(child: Text("Chargement..."))
          : QuestionW(
              _QuestionController(
                content.question,
              ),
              content.categorie.color),
    );
  }
}

class _QuestionController extends BaseQuestionController {
  _QuestionController(Question question) : super(question) {
    state.buttonEnabled = false;
  }

  @override
  void onPrimaryButtonClick() {}
}
