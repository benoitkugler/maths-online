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

class EleveApp extends StatelessWidget {
  final Audio audioPlayer;
  final BuildMode buildMode;
  final Upgrader? checkUprades;

  const EleveApp(this.audioPlayer, this.buildMode,
      {super.key, this.checkUprades});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
        title: 'Isyro',
        theme: theme,
        debugShowCheckedModeBanner: false,
        localizationsDelegates: localizations,
        supportedLocales: locales,
        home: _AppScaffold(
          audioPlayer,
          buildMode,
          checkUprades: checkUprades,
        ));
  }
}

class _AppScaffold extends StatefulWidget {
  final Audio audioPlayer;
  final BuildMode buildMode;
  final Upgrader? checkUprades;

  const _AppScaffold(this.audioPlayer, this.buildMode,
      {Key? key, this.checkUprades})
      : super(key: key);

  @override
  State<_AppScaffold> createState() => __AppScaffoldState();
}

class __AppScaffoldState extends State<_AppScaffold> {
  UserSettings settings = UserSettings();

  /// [trivialMetaCache] stores the credentials needed
  /// to reconnect in game.
  Map<String, String> trivialMetaCache = {};

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

    if (!settings.hasBeenLaunched) {
      WidgetsBinding.instance.addPostFrameCallback((_) => _showWelcomeScreen());
    }
  }

  void _showAudioSettings() {
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

  void _showAppSettings() async {
    final newSettings = await Navigator.of(context).push(
        MaterialPageRoute<UserSettings>(
            builder: (_) => Settings(widget.buildMode)));
    if (newSettings != null) {
      setState(() {
        settings = newSettings;
      });
    }
  }

  void _showWelcomeScreen() async {
    final goTo = await showDialog<bool>(
      context: context,
      builder: (context) => _WelcomeDialog(() {
        Navigator.of(context).pop(true);
      }),
    );

    // in any case, register the screen has been seen
    settings.hasBeenLaunched = true;
    saveUserSettings(settings);

    if (goTo != null && goTo) _showAppSettings();
  }

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
            body: TrivialGameSelect(TrivialSettings(
                widget.buildMode, trivialMetaCache, settings)))));
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
    final isIdentified = settings.studentID.isNotEmpty;
    final onPop = Navigator.of(context).push(MaterialPageRoute<void>(
        builder: (_) => isIdentified
            ? HomeworkW(ServerHomeworkAPI(widget.buildMode, settings.studentID))
            : const HomeworkDisabled()));
    onPop.then((value) => widget.audioPlayer.pause());
  }

  @override
  Widget build(BuildContext context) {
    final body = Card(
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
    return Scaffold(
      appBar: AppBar(
        title: const Text('Isyro'),
        actions: [
          IconButton(
            onPressed: () => _showAudioSettings(),
            icon: const Icon(IconData(0xe378, fontFamily: 'MaterialIcons')),
            tooltip: "Choisir la musique",
          ),
          IconButton(
            onPressed: () => _showAppSettings(),
            icon: const Icon(IconData(0xe57f, fontFamily: 'MaterialIcons')),
            tooltip: "Paramètres",
          ),
        ],
      ),
      body: widget.checkUprades == null
          ? body
          : UpgradeAlert(
              upgrader: widget.checkUprades,
              child: body,
            ),
    );
  }
}

class _WelcomeDialog extends StatelessWidget {
  final void Function() goToSettings;
  const _WelcomeDialog(this.goToSettings);

  @override
  Widget build(BuildContext context) {
    return AlertDialog(
      title: const Text("Bienvenue sur Isyro !"),
      content:
          const Text("Pour commencer, et si tu personnalisais ton appli ?"),
      actions: [
        TextButton(
            onPressed: goToSettings, child: const Text("Editer mon profil"))
      ],
    );
  }
}
