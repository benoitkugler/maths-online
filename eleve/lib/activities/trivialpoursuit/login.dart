import 'dart:convert';

import 'package:eleve/activities/trivialpoursuit/controller.dart';
import 'package:eleve/build_mode.dart';
import 'package:eleve/settings.dart';
import 'package:eleve/shared/activity_start.dart';
import 'package:eleve/shared/errors.dart';
import 'package:eleve/shared/pin.dart';
import 'package:eleve/shared/students.gen.dart';
import 'package:eleve/types/src_prof_trivial.dart';
import 'package:eleve/types/src_sql_editor.dart';
import 'package:eleve/types/src_sql_trivial.dart';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:http/http.dart' as http;

class SaveGameMetaNotification extends Notification {
  final String gameCode;
  final String gameMeta;

  const SaveGameMetaNotification(this.gameCode, this.gameMeta);
}

class TrivialSettings {
  final BuildMode buildMode;
  final UserSettings settings;

  const TrivialSettings(this.buildMode, this.settings);

  // dispatch is called to trigger game meta on-disk save
  Future<GameAcces> _login(
      String code, void Function(String, String) saveGameMeta) async {
    const sessionIDKey = "session-id";
    // we assume that the time to type the code is enough to load the settings
    final uri = Uri.parse(buildMode.serverURL("/trivial/game/setup", query: {
      sessionIDKey: code,
      studentIDKey: settings.studentID,
      // send (optional) meta so that we may reconnect
      TrivialPoursuitController.gameMetaKey:
          settings.trivialGameMetas[code] ?? "",
    }));

    final resp = await http.get(uri);
    final body = jsonDecode(resp.body) as JSON;
    // body is either the expect GameMeta or an error
    if (body.containsKey("GameMeta")) {
      final gameMeta = body["GameMeta"] as String;
      // save in memory...
      settings.trivialGameMetas[code] = gameMeta;
      // ... and trigger on disk save as well, so
      // that reconnection is possible accross app restart
      saveGameMeta(code, gameMeta);

      return GameAcces(
          code, settings.studentID, settings.studentPseudo, gameMeta);
    }

    throw body["message"] as String;
  }

// the returned Future completes when the route is popped
  Future<void> _showGameBoard(
      GameAcces data, BuildContext context, bool isSelfLaunched) async {
    final route = Navigator.of(context).push(MaterialPageRoute<void>(
      settings: const RouteSettings(name: "/board"),
      builder: (_) => Scaffold(
          appBar: AppBar(
              automaticallyImplyLeading: true, title: const Text("Triv'Maths")),
          body: NotificationListener<GameTerminatedNotification>(
              onNotification: (n) {
                settings.trivialGameMetas.remove(data.code);
                return true;
              },
              child:
                  TrivialPoursuitController(buildMode, data, isSelfLaunched))),
    ));

    return route;
  }
}

/// [TrivialGameSelect] is a home screen allowing the
/// user to choose between in classroom/self access games
class TrivialGameSelect extends StatelessWidget {
  final TrivialSettings settings;
  final void Function(String gameCode, String gameMeta) saveMeta;

  const TrivialGameSelect(this.settings, this.saveMeta, {super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text("Jouer à TrivMaths"),
      ),
      body: Column(mainAxisAlignment: MainAxisAlignment.spaceEvenly, children: [
        LaunchCard(
            "Accéder à une partie",
            "J'ai un code et je veux rejoindre une partie existante.",
            const Icon(Icons.login_outlined), () {
          Navigator.of(context).push(MaterialPageRoute<void>(
              builder: (_) => Scaffold(body: _Loggin(settings, saveMeta))));
        }),
        const Divider(thickness: 4),
        LaunchCard(
            "Créer une partie",
            "Je veux démarrer une partie et partager le code avec des amis.",
            const Icon(Icons.add_box_outlined),
            settings.settings.studentID.isEmpty
                ? null
                : () {
                    Navigator.of(context).push(MaterialPageRoute<void>(
                        builder: (_) => Scaffold(
                            appBar: AppBar(
                                title: const Text("Démarrer une partie")),
                            body: _SelfaccessList(settings, saveMeta))));
                  }),
      ]),
    );
  }
}

class _SelfaccessList extends StatefulWidget {
  final TrivialSettings settings;
  final void Function(String gameCode, String gameMeta) saveMeta;

  const _SelfaccessList(this.settings, this.saveMeta, {super.key});

  @override
  State<_SelfaccessList> createState() => __SelfaccessListState();
}

class __SelfaccessListState extends State<_SelfaccessList> {
  List<Trivial>? trivials;

  @override
  void initState() {
    _fetchTrivials();
    super.initState();
  }

  void _fetchTrivials() async {
    final uri = Uri.parse(widget.settings.buildMode
        .serverURL("/api/student/trivial/selfaccess", query: {
      studentIDKey: widget.settings.settings.studentID,
    }));

    try {
      final resp = await http.get(uri);
      final body = checkServerError(resp.body);
      final data = getSelfaccessOutFromJson(body);
      setState(() {
        trivials = data.trivials;
      });
    } catch (e) {
      showError("Erreur", e, context);
    }
  }

  void _launchTrivial(Trivial trivial) async {
    final uri = Uri.parse(widget.settings.buildMode
        .serverURL("/api/student/trivial/selfaccess/launch", query: {
      studentIDKey: widget.settings.settings.studentID,
      "trivial-id": trivial.id.toString(),
    }));

    final LaunchSelfaccessOut data;
    try {
      final resp = await http.get(uri);
      final body = checkServerError(resp.body);
      data = launchSelfaccessOutFromJson(body);
    } catch (e) {
      showError("Erreur", e, context);
      return;
    }

    if (!mounted) return;
    Navigator.of(context).push(MaterialPageRoute<void>(
        builder: (context) =>
            _GameLaunchedScreen(data.gameID, () => _joinGame(data.gameID))));
  }

  void _joinGame(String code) async {
    try {
      final data = await widget.settings._login(code, widget.saveMeta);
      if (!mounted) return;
      widget.settings._showGameBoard(data, context, true);
    } catch (e) {
      showError("Impossible de se connecter", e, context);
      return;
    }
  }

  @override
  Widget build(BuildContext context) {
    return trivials == null
        ? const Center(
            child: Row(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                CircularProgressIndicator(),
                SizedBox(width: 20),
                Text("Chargement..."),
              ],
            ),
          )
        : trivials!.isEmpty
            ? const Center(child: Text("Aucune partie n'est disponible."))
            : ListView(
                children: trivials!
                    .map((e) => _TrivialRow(e, () => _launchTrivial(e)))
                    .toList(),
              );
  }
}

class _GameLaunchedScreen extends StatelessWidget {
  final String gameCode;
  final void Function() onJoin;

  const _GameLaunchedScreen(this.gameCode, this.onJoin, {super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text("Rejoindre la partie")),
      body: Padding(
        padding: const EdgeInsets.all(8.0),
        child: Center(
          child: Column(mainAxisAlignment: MainAxisAlignment.center, children: [
            const Text(
                "La partie a bien été lancée ! Voici le code d'accès à partager :"),
            const SizedBox(height: 30),
            ElevatedButton.icon(
                onPressed: () async {
                  await Clipboard.setData(ClipboardData(text: gameCode));
                  ScaffoldMessenger.of(context).showSnackBar(const SnackBar(
                      backgroundColor: Colors.lightGreen,
                      content: Text("Code copié dans le presse-papier.")));
                },
                icon: const Icon(Icons.copy),
                label: Text(
                  gameCode,
                  style: const TextStyle(fontSize: 16),
                )),
            const SizedBox(height: 40),
            ElevatedButton(
                onPressed: onJoin, child: const Text("Rejoindre la partie"))
          ]),
        ),
      ),
    );
  }
}

class _TrivialRow extends StatelessWidget {
  final Trivial trivial;
  final void Function() onTap;
  const _TrivialRow(this.trivial, this.onTap, {super.key});

  @override
  Widget build(BuildContext context) {
    return ListTile(
      onTap: onTap,
      title: Text(trivial.name),
      subtitle: Text(_categories),
    );
  }

  String get _categories {
    final cats = trivial.questions.tags;
    final allChapters = <String>{};
    for (var element in cats) {
      for (var element in element) {
        for (var element in element) {
          if (element.section == Section.chapter ||
              element.section == Section.level) {
            allChapters.add(element.tag);
          }
        }
      }
    }
    // restrict to the chapters found in every questions
    return allChapters
        .where((ch) => cats
            .every((l) => l.every((inter) => inter.any((ts) => ts.tag == ch))))
        .join(", ");
  }
}

/// Loggin is an introduction screen to access
/// a TrivialPoursuit game
class _Loggin extends StatefulWidget {
  final TrivialSettings settings;
  final void Function(String gameCode, String gameMeta) saveMeta;

  const _Loggin(this.settings, this.saveMeta, {Key? key}) : super(key: key);

  @override
  _LogginState createState() => _LogginState();
}

class _LogginState extends State<_Loggin> {
  final pinController = TextEditingController();

  @override
  void initState() {
    if (widget.settings.buildMode == BuildMode.debug) {
      // skip loggin screen
      WidgetsBinding.instance.addPostFrameCallback((_) => widget.settings
          ._showGameBoard(const GameAcces("", "", "", ""), context, false));
    }

    super.initState();
  }

  void _showError(dynamic error) {
    showError("Impossible de se connecter.", error, context);
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text("Rejoindre une partie"),
      ),
      body: Pin(
        "Code de la partie",
        _launchTrivialPoursuit,
        controller: pinController,
      ),
    );
  }

  void _launchTrivialPoursuit(String code) async {
    try {
      final data = await widget.settings._login(code, widget.saveMeta);
      if (!mounted) return;
      final route = widget.settings._showGameBoard(data, context, false);
      route.then((value) {
        setState(() {
          pinController.clear();
        });
      });
    } catch (e) {
      _showError(e);
      return;
    }
  }
}
