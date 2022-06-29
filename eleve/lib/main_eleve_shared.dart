import 'package:eleve/audio.dart';
import 'package:eleve/build_mode.dart';
import 'package:eleve/exercice/exercice.dart';
import 'package:eleve/main_shared.dart';
import 'package:eleve/questions/question_gallery.dart';
import 'package:eleve/questions/types.gen.dart';
import 'package:eleve/settings.dart';
import 'package:eleve/shared_gen.dart';
import 'package:eleve/trivialpoursuit/controller.dart';
import 'package:eleve/trivialpoursuit/login.dart';
import 'package:flutter/material.dart' hide Flow;

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

  const EleveApp(this.audioPlayer, this.buildMode, {Key? key})
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
    final _settings = await loadUserSettings();
    setState(() {
      settings = _settings;
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
          // body: _HomePage(widget.audioPlayer, widget.buildMode, settings),
          body: ExerciceW(
              widget.buildMode,
              Exercice(
                  InstantiatedExercice(
                      0, "REcherche de primitives", Flow.sequencial, [
                    InstantiatedQuestion(
                        0, Question("", [NumberFieldBlock(0)]), []),
                    InstantiatedQuestion(
                        0,
                        Question("", [
                          NumberFieldBlock(0),
                          DropDownFieldBlock([
                            [TextOrMath("AAA", true)],
                            [TextOrMath("BB", true)],
                          ], 1)
                        ]),
                        []),
                    InstantiatedQuestion(
                        0, Question("", [NumberFieldBlock(0)]), []),
                  ], [
                    1,
                    2,
                    3
                  ]),
                  ProgressionExt(
                      Progression(0, 0),
                      [
                        [true, true],
                        [true, false],
                        []
                      ],
                      1)),
              (dataIn) async => EvaluateExerciceOut(
                      dataIn.answers.map((index, qu) => MapEntry(
                          index,
                          QuestionAnswersOut(
                              qu.answer.data
                                  .map((key, value) => MapEntry(key, false)),
                              {}))),
                      ProgressionExt(Progression(0, 0), [], 1),
                      [
                        InstantiatedQuestion(
                            0,
                            Question("", [
                              TextBlock([TextOrMath("text", true)], true, false,
                                  true),
                              NumberFieldBlock(0)
                            ]),
                            []),
                        InstantiatedQuestion(
                            0,
                            Question("", [
                              TextBlock([TextOrMath("text", true)], true, false,
                                  true),
                              NumberFieldBlock(0)
                            ]),
                            []),
                        InstantiatedQuestion(
                            0,
                            Question("", [
                              TextBlock([TextOrMath("text", true)], true, false,
                                  true),
                              NumberFieldBlock(0)
                            ]),
                            []),
                      ])),
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
