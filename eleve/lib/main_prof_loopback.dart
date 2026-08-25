import 'dart:async';
import 'dart:convert';
import 'dart:js_interop' as js;
import 'package:web/web.dart' as web;

import 'package:eleve/build_mode.dart';
import 'package:eleve/loopback/loopback.dart';
import 'package:eleve/types/src_prof_preview.dart';
import 'package:eleve/main_shared.dart';
import 'package:eleve/shared/errors.dart';
import 'package:flutter/material.dart';

void main() {
  // the static app is called via an url setting the session ID
  // note that the MaterialApp routing erase these parameters,
  // so that we need to fetch it early
  final uri = Uri.parse(web.window.location.href);
  // final id = uri.queryParameters["sessionID"]!;
  final mode = uri.queryParameters["mode"];
  final bm = APISetting.fromString(mode ?? "");

  runApp(LoopbackApp(bm));
}

/// [LoopbackApp] show the content of a question or an exercice instance
/// being edited, as it will be displayed to the student.
/// It is meant to be embedded in a Web page, not used as a mobile app.
class LoopbackApp extends StatelessWidget {
  final BuildMode buildMode;

  const LoopbackApp(this.buildMode, {super.key});

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

/// listen to HTML events and switch between question
/// or exercice mode accordingly
class _EditorLoopback extends StatefulWidget {
  final BuildMode buildMode;

  const _EditorLoopback(this.buildMode);

  @override
  State<_EditorLoopback> createState() => _EditorLoopbackState();
}

class _EditorLoopbackState extends State<_EditorLoopback> {
  late final StreamSubscription<web.MessageEvent> subs;
  late final LoopbackServerAPI _api;

  LoopbackServerEvent event = const LoopbackPaused();

  @override
  void initState() {
    subs = web.window.onMessage.listen((event) {
      listen((event.data as js.JSString).toDart);
    });
    _api = LoopbackServerAPI(widget.buildMode);
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
        event = loopbackServerEventFromJson(jsonDecode(jsonEvent));
      });
    } catch (e) {
      _showError(e);
      return;
    }
  }

  @override
  Widget build(BuildContext context) {
    return EditorLoopback(event, _api);
  }
}
