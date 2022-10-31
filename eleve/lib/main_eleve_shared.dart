import 'package:eleve/activities/homework/homework.dart';
import 'package:eleve/activities/trivialpoursuit/controller.dart';
import 'package:eleve/activities/trivialpoursuit/login.dart';
import 'package:eleve/audio.dart';
import 'package:eleve/build_mode.dart';
import 'package:eleve/main_shared.dart';
import 'package:eleve/settings.dart';
import 'package:eleve/shared/activity_start.dart';
import 'package:flutter/material.dart' hide Flow;
import 'package:upgrader/upgrader.dart';

Future<Audio> loadAudioFromSettings() async {
  WidgetsFlutterBinding
      .ensureInitialized(); // required to load the settings path

  final audio = Audio();
  final settings = await loadUserSettings();
  audio.setSongs(settings.songs);
  return audio;
}

class EleveApp extends StatefulWidget {
  final Audio audioPlayer;
  final BuildMode buildMode;
  final Upgrader? checkUprades;

  const EleveApp(this.audioPlayer, this.buildMode,
      {Key? key, this.checkUprades})
      : super(key: key);

  @override
  State<EleveApp> createState() => _EleveAppState();
}

class _EleveAppState extends State<EleveApp> {
  UserSettings settings = UserSettings();

  @override
  void initState() {
    _loadSettings();
    super.initState();
  }

  void _loadSettings() async {
    final set = await loadUserSettings();
    setState(() {
      settings = set;
    });
  }

  void _showAudioSettings(BuildContext context) {
    final ct = widget.audioPlayer.playlist.toList();
    final onPop = Navigator.of(context)
        .push<void>(MaterialPageRoute<void>(builder: (_) => Playlist(ct)));
    onPop.then((_) async {
      widget.audioPlayer.setSongs(ct);
      settings.songs = ct;
      await saveUserSettings(settings); // commit on disk

      ScaffoldMessenger.of(context).showSnackBar(SnackBar(
        backgroundColor: Theme.of(context).colorScheme.secondary,
        content: const Text("Playlist mise à jour."),
      ));
    });
  }

  void _showAppSettings(BuildContext context) async {
    final newSettings = await Navigator.of(context).push(
        MaterialPageRoute<UserSettings>(
            builder: (_) => Settings(widget.buildMode)));
    if (newSettings != null) {
      setState(() {
        settings = newSettings;
      });
    }
  }

  @override
  Widget build(BuildContext context) {
    final body = _AppBody(widget.audioPlayer, widget.buildMode, settings);
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
          body: widget.checkUprades == null
              ? body
              : UpgradeAlert(
                  upgrader: widget.checkUprades,
                  child: body,
                ),
        ));
  }
}

class _AppBody extends StatefulWidget {
  final Audio audioPlayer;
  final BuildMode buildMode;
  final UserSettings settings;

  const _AppBody(this.audioPlayer, this.buildMode, this.settings, {Key? key})
      : super(key: key);

  @override
  State<_AppBody> createState() => _AppBodyState();
}

class _AppBodyState extends State<_AppBody> {
  /// [trivialMetaCache] stores the credentials needed
  /// to reconnect in game.
  Map<String, String> trivialMetaCache = {};

  void _launchTrivialPoursuit() async {
    final onDone = await Navigator.of(context).push(MaterialPageRoute<bool>(
        builder: (context) =>
            ActivityStart(() => Navigator.of(context).pop(true))));
    if (onDone == null) {
      return;
    }

    widget.audioPlayer.run();
    final onPop = Navigator.of(context).push(MaterialPageRoute<void>(
        builder: (_) => Scaffold(
            body: TrivialPoursuitLoggin(
                widget.buildMode, trivialMetaCache, widget.settings))));
    onPop.then((value) => widget.audioPlayer.pause());
  }

  void _launchHomework() async {
    final onDone = await Navigator.of(context).push(MaterialPageRoute<bool>(
        builder: (context) =>
            ActivityStart(() => Navigator.of(context).pop(true))));
    if (onDone == null) {
      return;
    }

    widget.audioPlayer.run();
    final onPop = Navigator.of(context).push(MaterialPageRoute<void>(
        builder: (_) => Scaffold(
            body: HomeworkW(ServerHomeworkAPI(
                widget.buildMode, widget.settings.studentID)))));
    onPop.then((value) => widget.audioPlayer.pause());
  }

  @override
  Widget build(BuildContext context) {
    return Card(
      child: Center(
        child: Column(
          // crossAxisAlignment: CrossAxisAlignment.stretch,
          mainAxisAlignment: MainAxisAlignment.spaceEvenly,
          children: [
            const Padding(
              padding: EdgeInsets.all(8.0),
              child: Text(
                "Bienvenue sur Isyro !",
                style: TextStyle(fontSize: 25),
              ),
            ),
            Column(
              mainAxisSize: MainAxisSize.min,
              children: [
                const Text(
                  "Activités disponibles",
                  style: TextStyle(fontSize: 20),
                ),
                const SizedBox(height: 20),
                Padding(
                  padding: const EdgeInsets.all(8.0),
                  child: Wrap(
                    spacing: 16,
                    runSpacing: 16,
                    alignment: WrapAlignment.spaceEvenly,
                    children: [
                      TrivialActivityIcon(_launchTrivialPoursuit),
                      HomeworkActivityIcon(_launchHomework),
                    ],
                  ),
                )
              ],
            )
          ],
        ),
      ),
    );
  }
}
