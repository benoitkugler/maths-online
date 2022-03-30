import 'package:audioplayers/audioplayers.dart';
import 'package:eleve/build_mode.dart';
import 'package:eleve/exercices/question_gallery.dart';
import 'package:eleve/exercices/tree.dart';
import 'package:eleve/exercices/types.gen.dart';
import 'package:eleve/trivialpoursuit/game.dart';
import 'package:eleve/trivialpoursuit/login.dart';
import 'package:flutter/material.dart';
import 'package:flutter_localizations/flutter_localizations.dart';

const Color darkBlue = Color.fromARGB(255, 27, 54, 82);

final bm = buildMode();
// final bm = BuildMode.dev;

void main() async {
  runApp(const MyApp());
}

class Audio {
  final AudioCache _cache = AudioCache(prefix: "lib/music/");
  final List<String> songs;

  int _currentSong = -1;
  AudioPlayer? _player;

  Audio(this.songs);

  void run() {
    _startNextSong();
  }

  void pause() {
    if (_player == null) {
      return;
    }
    _player!.stop();
  }

  void _startNextSong() async {
    _currentSong++;
    _player = await _cache.play(songs[_currentSong % songs.length]);
    _player!.onPlayerCompletion.listen((event) {
      _startNextSong();
    });
  }
}

class MyApp extends StatelessWidget {
  const MyApp({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Isiro',
      theme: ThemeData.dark().copyWith(
        scaffoldBackgroundColor: darkBlue,
        cardTheme: ThemeData.dark()
            .cardTheme
            .copyWith(color: darkBlue.withOpacity(0.5)),
      ),
      debugShowCheckedModeBanner: false,
      localizationsDelegates: const [
        GlobalMaterialLocalizations.delegate,
        GlobalWidgetsLocalizations.delegate,
        GlobalCupertinoLocalizations.delegate,
      ],
      supportedLocales: const [
        Locale('fr', ''), // French, no country code
        Locale('en', ''), // English, no country code
      ],
      home: Scaffold(
          appBar: AppBar(
            title: const Text('Isiro'),
          ),
          // body: _HomePage(),
          body: Center(
              child: Tree(
                  Colors.blue,
                  TreeController(
                      TreeFieldBlock([
                        [2, 2],
                        [2, 3],
                        [3, 2],
                        [2, 2, 2],
                      ], [
                        TextOrMath("A", false),
                        TextOrMath("x", true),
                        TextOrMath("C", false),
                      ], 0),
                      () {})))),
    );
  }
}

class _HomePage extends StatefulWidget {
  final player = Audio(["GrooveBow.mp3", "NouvelleTrajectoire.mp3"]);

  _HomePage({Key? key}) : super(key: key);

  @override
  State<_HomePage> createState() => _HomePageState();
}

class _HomePageState extends State<_HomePage> {
  void _launchTrivialPoursuit() {
    widget.player.run();
    final onPop = Navigator.of(context).push<void>(MaterialPageRoute<void>(
        builder: (_) => Scaffold(body: TrivialPoursuitLoggin(bm))));
    onPop.then((value) => widget.player.pause());
  }

  void _launchQuestionGallery() {
    widget.player.run();
    final onPop = Navigator.of(context)
        .push(MaterialPageRoute<void>(builder: (_) => QuestionGallery(bm)));
    onPop.then((value) => widget.player.pause());
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
                "Bienvenue dans Isiro",
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
                  // TODO: polish
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
