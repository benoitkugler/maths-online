import 'package:eleve/build_mode.dart';
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

  @override
  void initState() {
    if (widget.buildMode == BuildMode.debug) {
      // skip loggin screen
      Future.delayed(
          const Duration(milliseconds: 50), () => _launchTrivialPoursuit(""));
    }

    super.initState();
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
    final route = Navigator.of(context).push(MaterialPageRoute<void>(
      settings: const RouteSettings(name: "/board"),
      builder: (_) => Scaffold(
          appBar: AppBar(),
          body: TrivialPoursuitController(widget.buildMode, code)),
    ));

    route.then((value) {
      setState(() {
        pinController.clear();
      });
    });
  }
}
