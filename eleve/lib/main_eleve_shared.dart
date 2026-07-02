import 'package:eleve/activities/automatismes/automatismes.dart';
import 'package:eleve/activities/ceintures/api.dart';
import 'package:eleve/activities/ceintures/ceintures.dart';
import 'package:eleve/activities/homework/homework.dart';
import 'package:eleve/activities/trivialpoursuit/controller.dart';
import 'package:eleve/activities/trivialpoursuit/login.dart';
import 'package:eleve/shared/audio.dart';
import 'package:eleve/build_mode.dart';
import 'package:eleve/main_shared.dart';
import 'package:eleve/settings.dart';
import 'package:eleve/shared/activity_start.dart';
import 'package:eleve/shared/errors.dart';
import 'package:eleve/shared/settings_shared.dart';
import 'package:eleve/types/src_sql_events.dart';
import 'package:flutter/material.dart' hide Flow;
import 'package:http/http.dart' as http;
import 'package:upgrader/upgrader.dart';

Future<Audio> loadAudioFromSettings(SettingsStorage handler) async {
  WidgetsFlutterBinding.ensureInitialized(); // required to load the settings path

  final audio = Audio();
  final settings = await handler.load();
  audio.setSongs(settings.songs);
  return audio;
}

class EleveApp extends StatelessWidget {
  final Audio audioPlayer;
  final SettingsStorage settingsHandler;
  final BuildMode buildMode;
  final Upgrader? checkUprades;

  const EleveApp(
    this.audioPlayer,
    this.settingsHandler,
    this.buildMode, {
    super.key,
    this.checkUprades,
  });

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
        settingsHandler,
        buildMode,
        checkUprades: checkUprades,
      ),
    );
  }
}

class _AppScaffold extends StatefulWidget {
  final Audio audioPlayer;
  final SettingsStorage storage;
  final BuildMode buildMode;
  final Upgrader? checkUprades;

  const _AppScaffold(
    this.audioPlayer,
    this.storage,
    this.buildMode, {
    Key? key,
    this.checkUprades,
  }) : super(key: key);

  @override
  State<_AppScaffold> createState() => __AppScaffoldState();
}

class __AppScaffoldState extends State<_AppScaffold> {
  late final SettingsHandler settings;

  @override
  void initState() {
    _loadSettings();
    super.initState();
  }

  void _loadSettings() async {
    settings = SettingsHandler(widget.storage);
    await settings.init();
    setState(() {});

    if (!settings.settings.hasBeenLaunched) {
      WidgetsBinding.instance.addPostFrameCallback((_) => _showWelcomeScreen());
    }
  }

  void _showAudioSettings() {
    final ct = widget.audioPlayer.playlist;
    final onPop = Navigator.of(
      context,
    ).push<void>(MaterialPageRoute<void>(builder: (_) => Playlist(ct)));
    onPop.then((_) async {
      widget.audioPlayer.setSongs(ct);

      await settings.saveSongs(ct);

      if (!mounted) return;
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(
          backgroundColor: Theme.of(context).colorScheme.secondary,
          content: const Text("Playlist mise à jour."),
        ),
      );

      // notify the server and show event
      final studentID = settings.settings.studentID;
      if (studentID.isNotEmpty) {
        final resp = await http.get(
          widget.buildMode.serverURL(
            "/api/student/set-playlist",
            query: {studentIDKey: studentID},
          ),
        );
        try {
          final notif = eventNotificationFromJson(checkServerError(resp.body));
        } catch (e) {
          // silently fail
        }
      }
    });
  }

  void _showProfile() async {
    Navigator.of(context).push(
      MaterialPageRoute<UserSettings>(
        builder: (_) => Settings(widget.buildMode, settings),
      ),
    );
  }

  void _showWelcomeScreen() async {
    final goTo = await showDialog<bool>(
      context: context,
      builder: (context) => _WelcomeDialog(() {
        Navigator.of(context).pop(true);
      }),
    );

    // in any case, register the screen has been seen
    settings.saveHasBeenLaunched();

    if (goTo != null && goTo) _showProfile();
  }

  void _launchTrivialPoursuit() async {
    final onDone = await Navigator.of(context).push(
      MaterialPageRoute<bool>(
        builder: (context) =>
            MathActivityStart(() => Navigator.of(context).pop(true)),
      ),
    );
    if (onDone == null) return;
    if (!mounted) return;

    widget.audioPlayer.run();
    await Navigator.of(context).push(
      MaterialPageRoute<void>(
        builder: (_) => Scaffold(
          body: TrivialGameSelect(TrivialSettings(widget.buildMode, settings)),
        ),
      ),
    );
    widget.audioPlayer.pause();
  }

  void _launchHomework() async {
    widget.audioPlayer.run();
    final isIdentified = settings.settings.studentID.isNotEmpty;
    await Navigator.of(context).push(
      MaterialPageRoute<void>(
        builder: (_) => isIdentified
            ? HomeworkStart(
                ServerHomeworkAPI(
                  widget.buildMode,
                  settings.settings.studentID,
                ),
              )
            : const HomeworkDisabled(),
      ),
    );
    widget.audioPlayer.pause();
  }

  void _launchCeintures() async {
    widget.audioPlayer.run();
    await Navigator.of(context).push(
      MaterialPageRoute<void>(
        builder: (_) =>
            CeinturesStart(ServerCeinturesAPI(widget.buildMode), settings),
      ),
    );
    widget.audioPlayer.pause();
  }

  void _launchAutomatismes() async {
    widget.audioPlayer.run();
    await Navigator.of(context).push(
      MaterialPageRoute<void>(
        builder: (_) => AutomatismesStart(
          AutomatismesServerAPI(widget.buildMode),
          TrivialSettings(widget.buildMode, settings),
        ),
      ),
    );
    widget.audioPlayer.pause();
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
                      CeinturesActivityIcon(_launchCeintures),
                      AutomatismesActivityIcon(_launchAutomatismes),
                    ],
                  ),
                ),
              ],
            ),
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
            icon: const Icon(Icons.library_music),
            tooltip: "Choisir la musique",
          ),
          IconButton(
            onPressed: () => _showProfile(),
            icon: const Icon(Icons.account_circle),
            tooltip: "Afficher ton profil",
          ),
        ],
      ),
      body: widget.checkUprades == null
          ? body
          : UpgradeAlert(
              upgrader: widget.checkUprades,
              showIgnore: false,
              showLater: false,
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
      content: const Text(
        "Pour commencer, et si tu personnalisais ton appli ?",
      ),
      actions: [
        TextButton(
          onPressed: goToSettings,
          child: const Text("Editer mon profil"),
        ),
      ],
    );
  }
}
