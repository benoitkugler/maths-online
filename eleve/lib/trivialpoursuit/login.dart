import 'dart:convert';

import 'package:eleve/build_mode.dart';
import 'package:eleve/exercices/types.gen.dart';
import 'package:eleve/settings.dart';
import 'package:eleve/shared/pin.dart';
import 'package:eleve/trivialpoursuit/controller.dart';
import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;

/// Loggin is an introduction screen to access
/// a TrivialPoursuit game
class TrivialPoursuitLoggin extends StatefulWidget {
  final BuildMode buildMode;
  final Map<String, String> gameMetaCache;

  const TrivialPoursuitLoggin(this.buildMode, this.gameMetaCache, {Key? key})
      : super(key: key);

  @override
  _TrivialPoursuitLogginState createState() => _TrivialPoursuitLogginState();
}

class _TrivialPoursuitLogginState extends State<TrivialPoursuitLoggin> {
  final pinController = TextEditingController();
  UserSettings settings = {};

  @override
  void initState() {
    _loadSettings();
    if (widget.buildMode == BuildMode.debug) {
      // skip loggin screen
      Future.delayed(
          const Duration(milliseconds: 50), () => _launchTrivialPoursuit(""));
    }

    super.initState();
  }

  void _loadSettings() async {
    settings = await loadUserSettings();
  }

  void _showError(dynamic error) {
    ScaffoldMessenger.of(context).showSnackBar(SnackBar(
        duration: const Duration(seconds: 6),
        backgroundColor: Theme.of(context).colorScheme.error,
        content: RichText(
            text: TextSpan(children: [
          const TextSpan(
              text: "Impossible de se connecter. \n",
              style: TextStyle(fontWeight: FontWeight.bold)),
          const TextSpan(text: "Détails : "),
          TextSpan(
              text: "$error",
              style: const TextStyle(fontStyle: FontStyle.italic)),
        ]))));
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text("Rejoindre une partie"),
      ),
      body: Pin("Code de la partie", pinController, _launchTrivialPoursuit),
    );
  }

  Future<String> _setupGame(String code) async {
    const sessionIDKey = "session-id";
    // we assume that the time to type the code is enough to load the settings
    final uri =
        Uri.parse(widget.buildMode.serverURL("/trivial/game/setup", query: {
      sessionIDKey: code,
      studentIDKey: settings[studentIDKey],
      // send (optional) meta so that we may reconnect
      TrivialPoursuitController.gameMetaKey: widget.gameMetaCache[code] ?? "",
    }));

    try {
      final resp = await http.get(uri);
      final body = jsonDecode(resp.body) as JSON;
      // body is either the expect GameMeta or an error
      if (body.containsKey("GameMeta")) {
        final gameMeta = body["GameMeta"] as String;
        widget.gameMetaCache[code] = gameMeta;
        return gameMeta;
      }
      throw body["message"] as String;
    } catch (e) {
      widget.gameMetaCache.remove(code);
      _showError(e);
      return "";
    }
  }

  void _launchTrivialPoursuit(String code) async {
    final gameMeta = await _setupGame(code);
    if (gameMeta.isEmpty) {
      return;
    }

    final student = GameAcces(settings[studentIDKey] ?? "",
        settings[studentPseudoKey] ?? "", gameMeta);

    final route = Navigator.of(context).push(MaterialPageRoute<void>(
      settings: const RouteSettings(name: "/board"),
      builder: (_) => Scaffold(
          appBar: AppBar(
              automaticallyImplyLeading: true, title: const Text("Triv'Maths")),
          body: NotificationListener<GameTerminatedNotification>(
              onNotification: (n) {
                widget.gameMetaCache.remove(code);
                return true;
              },
              child: TrivialPoursuitController(widget.buildMode, student))),
    ));

    route.then((value) {
      setState(() {
        pinController.clear();
      });
    });
  }
}
