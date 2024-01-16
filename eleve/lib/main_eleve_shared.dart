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
  WidgetsFlutterBinding
      .ensureInitialized(); // required to load the settings path

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

  const EleveApp(this.audioPlayer, this.settingsHandler, this.buildMode,
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
          settingsHandler,
          buildMode,
          checkUprades: checkUprades,
        ));
  }
}

class _AppScaffold extends StatefulWidget {
  final Audio audioPlayer;
  final SettingsStorage handler;
  final BuildMode buildMode;
  final Upgrader? checkUprades;

  const _AppScaffold(this.audioPlayer, this.handler, this.buildMode,
      {Key? key, this.checkUprades})
      : super(key: key);

  @override
  State<_AppScaffold> createState() => __AppScaffoldState();
}

class __AppScaffoldState extends State<_AppScaffold> {
  UserSettings settings = UserSettings();

  @override
  void initState() {
    _loadSettings();
    super.initState();
  }

  void _loadSettings() async {
    final set = await widget.handler.load();
    setState(() {
      settings = set;
    });

    if (!settings.hasBeenLaunched) {
      WidgetsBinding.instance.addPostFrameCallback((_) => _showWelcomeScreen());
    }
  }

  void _showAudioSettings() {
    final ct = widget.audioPlayer.playlist;
    final onPop = Navigator.of(context)
        .push<void>(MaterialPageRoute<void>(builder: (_) => Playlist(ct)));
    onPop.then((_) async {
      widget.audioPlayer.setSongs(ct);
      settings.songs = ct;
      await widget.handler.save(settings); // commit on disk

      if (!mounted) return;
      ScaffoldMessenger.of(context).showSnackBar(SnackBar(
        backgroundColor: Theme.of(context).colorScheme.secondary,
        content: const Text("Playlist mise à jour."),
      ));

      // notify the server and show event
      final studentID = settings.studentID;
      if (studentID.isNotEmpty) {
        final resp = await http.get(widget.buildMode.serverURL(
            "/api/student/set-playlist",
            query: {studentIDKey: studentID}));
        try {
          final notif = eventNotificationFromJson(checkServerError(resp.body));
          print("${notif.events}  ${notif.points}");
        } catch (e) {
          // silently fail
        }
      }
    });
  }

  void _showAppSettings() async {
    final newSettings = await Navigator.of(context).push(
        MaterialPageRoute<UserSettings>(
            builder: (_) => Settings(widget.buildMode, widget.handler)));
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
    widget.handler.save(settings);

    if (goTo != null && goTo) _showAppSettings();
  }

  void _saveTrivialMeta(String gameCode, String gameMeta) async {
    settings.trivialGameMetas[gameCode] = gameMeta;
    await widget.handler.save(settings);
  }

  void _saveCeinturesAnonymousID(String id) async {
    settings.ceinturesAnonymousID = id;
    await widget.handler.save(settings);
  }

  void _launchTrivialPoursuit() async {
    final onDone = await Navigator.of(context).push(MaterialPageRoute<bool>(
        builder: (context) =>
            ActivityStart(() => Navigator.of(context).pop(true))));
    if (onDone == null) return;
    if (!mounted) return;

    widget.audioPlayer.run();
    await Navigator.of(context).push(MaterialPageRoute<void>(
        builder: (_) => Scaffold(
            body: TrivialGameSelect(TrivialSettings(widget.buildMode, settings),
                _saveTrivialMeta))));
    widget.audioPlayer.pause();
  }

  void _launchHomework() async {
    final onDone = await Navigator.of(context).push(MaterialPageRoute<bool>(
        builder: (context) =>
            ActivityStart(() => Navigator.of(context).pop(true))));
    if (onDone == null) return;
    if (!mounted) return;

    widget.audioPlayer.run();
    final isIdentified = settings.studentID.isNotEmpty;
    await Navigator.of(context).push(MaterialPageRoute<void>(
        builder: (_) => isIdentified
            ? HomeworkStart(
                ServerHomeworkAPI(widget.buildMode, settings.studentID))
            : const HomeworkDisabled()));
    widget.audioPlayer.pause();
  }

  void _launchCeintures() async {
    widget.audioPlayer.run();
    await Navigator.of(context).push(MaterialPageRoute<void>(
        builder: (_) => CeinturesStart(ServerCeinturesAPI(widget.buildMode),
            settings, _saveCeinturesAnonymousID)));
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
