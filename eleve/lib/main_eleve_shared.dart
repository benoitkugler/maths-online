import 'package:eleve/audio.dart';
import 'package:eleve/build_mode.dart';
import 'package:eleve/exercices/question_gallery.dart';
import 'package:eleve/main_shared.dart';
import 'package:eleve/settings.dart';
import 'package:eleve/trivialpoursuit/controller.dart';
import 'package:eleve/trivialpoursuit/login.dart';
import 'package:flutter/material.dart';

class EleveApp extends StatefulWidget {
  final Audio audioPlayer;
  final BuildMode buildMode;

  const EleveApp(this.audioPlayer, this.buildMode, {Key? key})
      : super(key: key);

  @override
  State<EleveApp> createState() => _EleveAppState();
}

class _EleveAppState extends State<EleveApp> {
  UserSettings settings = {};

  @override
  void initState() {
    _loadSettings();
    super.initState();
  }

  void _loadSettings() async {
    final _settings = await loadUserSettings();
    setState(() {
      settings = _settings;
    });
  }

  void _showAudioSettings(BuildContext context) {
    final ct = widget.audioPlayer.playlist.toList();
    final onPop = Navigator.of(context)
        .push<void>(MaterialPageRoute<void>(builder: (_) => Playlist(ct)));
    onPop.then((_) => widget.audioPlayer.setSongs(ct));
  }

  void _showAppSettings(BuildContext context) async {
    final newSettings = await Navigator.of(context).push(
        MaterialPageRoute<UserSettings>(builder: (_) => const Settings()));
    if (newSettings != null) {
      setState(() {
        settings = newSettings;
      });
    }
  }

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
        title: 'Isyro',
        theme: theme,
        debugShowCheckedModeBanner: false,
        localizationsDelegates: localizations,
        supportedLocales: locales,
        home: Scaffold(
          appBar: AppBar(
            title: const Text('Isyro'),
            actions: [
              Builder(
                builder: (context) => IconButton(
                  onPressed: () => _showAudioSettings(context),
                  icon:
                      const Icon(IconData(0xe378, fontFamily: 'MaterialIcons')),
                  tooltip: "Choisir la musique",
                ),
              ),
              Builder(
                builder: (context) => IconButton(
                  onPressed: () => _showAppSettings(context),
                  icon:
                      const Icon(IconData(0xe57f, fontFamily: 'MaterialIcons')),
                  tooltip: "Paramètres",
                ),
              )
            ],
          ),
          body: _HomePage(widget.audioPlayer, widget.buildMode, settings),
        ));
  }
}

class _HomePage extends StatefulWidget {
  final Audio audioPlayer;
  final BuildMode buildMode;
  final UserSettings settings;

  const _HomePage(this.audioPlayer, this.buildMode, this.settings, {Key? key})
      : super(key: key);

  @override
  State<_HomePage> createState() => _HomePageState();
}

class _HomePageState extends State<_HomePage> {
  /// [trivialMetaCache] stores the credentials needed
  /// to reconnect in game.
  Map<String, String> trivialMetaCache = {};

  void _launchTrivialPoursuit() async {
    widget.audioPlayer.run();
    final onPop = Navigator.of(context).push<void>(MaterialPageRoute<void>(
        builder: (_) => Scaffold(
            body: TrivialPoursuitLoggin(
                widget.buildMode, trivialMetaCache, widget.settings))));
    onPop.then((value) => widget.audioPlayer.pause());
  }

  void _launchQuestionGallery() {
    widget.audioPlayer.run();
    final onPop = Navigator.of(context).push(MaterialPageRoute<void>(
        builder: (_) => QuestionGallery(widget.buildMode)));
    onPop.then((value) => widget.audioPlayer.pause());
  }

  @override
  Widget build(BuildContext context) {
    return Center(
      child: Card(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.spaceEvenly,
          children: [
            const Padding(
              padding: EdgeInsets.all(8.0),
              child: Text(
                "Bienvenue dans Isyro",
                style: TextStyle(fontSize: 25),
              ),
            ),
            const Text(
              "Activités disponibles",
              style: TextStyle(fontSize: 20),
            ),
            Padding(
              padding: const EdgeInsets.all(8.0),
              child: Row(
                mainAxisAlignment: MainAxisAlignment.spaceEvenly,
                children: [
                  GameIcon(_launchTrivialPoursuit),
                  if (widget.buildMode != BuildMode.production)
                    ElevatedButton(
                        onPressed: _launchQuestionGallery,
                        child: const Text("Galerie de questions")),
                ],
              ),
            )
          ],
        ),
      ),
    );
  }
}