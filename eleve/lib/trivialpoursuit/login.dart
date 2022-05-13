import 'package:eleve/build_mode.dart';
import 'package:eleve/settings.dart';
import 'package:eleve/shared/pin.dart';
import 'package:eleve/trivialpoursuit/game.dart';
import 'package:flutter/material.dart';

/// Loggin is an introduction screen to access
/// a TrivialPoursuit game
class TrivialPoursuitLoggin extends StatefulWidget {
  final BuildMode buildMode;
  const TrivialPoursuitLoggin(this.buildMode, {Key? key}) : super(key: key);

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

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text("Rejoindre une partie"),
      ),
      body: Pin("Code de la partie", pinController, _launchTrivialPoursuit),
    );
  }

  void _launchTrivialPoursuit(String code) {
    // we assume that the time to type the code is enough to load the settings
    final student = StudentMeta(
        settings[studentIDKey] ?? "", settings[studentPseudoKey] ?? "", code);
    final route = Navigator.of(context).push(MaterialPageRoute<void>(
      settings: const RouteSettings(name: "/board"),
      builder: (_) => Scaffold(
          appBar: AppBar(
              automaticallyImplyLeading: false,
              title: const Text("Trivial Poursuit")),
          body: TrivialPoursuitController(widget.buildMode, student)),
    ));

    route.then((value) {
      setState(() {
        pinController.clear();
      });
    });
  }
}
