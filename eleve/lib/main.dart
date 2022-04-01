import 'package:eleve/audio.dart';
import 'package:eleve/build_mode.dart';
import 'package:eleve/exercices/question_gallery.dart';
import 'package:eleve/trivialpoursuit/game.dart';
import 'package:eleve/trivialpoursuit/login.dart';
import 'package:flutter/material.dart';
import 'package:flutter_localizations/flutter_localizations.dart';

const Color darkBlue = Color.fromARGB(255, 27, 54, 82);

final bm = buildMode();
// final bm = BuildMode.dev;

void main() {
  final audio = Audio();
  audio.setSongs(["GrooveBow.mp3", "NouvelleTrajectoire.mp3"]);
  runApp(MyApp(audio));
}

class MyApp extends StatelessWidget {
  final Audio audioPlayer;

  const MyApp(this.audioPlayer, {Key? key}) : super(key: key);

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
        body: _HomePage(audioPlayer),
        // body: Center(
        // child: Tree(
        //     Colors.blue,
        //     TreeController(
        //         TreeFieldBlock([
        //           [2, 2],
        //           [2, 3],
        //           [3, 2],
        //           [2, 2, 2],
        //         ], [
        //           TextOrMath("A", false),
        //           TextOrMath("x", true),
        //           TextOrMath("C", false),
        //         ], 0),
        //         () {})))
      ),
    );
  }
}

class _HomePage extends StatefulWidget {
  final Audio audioPlayer;

  const _HomePage(this.audioPlayer, {Key? key}) : super(key: key);

  @override
  State<_HomePage> createState() => _HomePageState();
}

class _HomePageState extends State<_HomePage> {
  void _launchTrivialPoursuit() {
    widget.audioPlayer.run();
    final onPop = Navigator.of(context).push<void>(MaterialPageRoute<void>(
        builder: (_) => Scaffold(body: TrivialPoursuitLoggin(bm))));
    onPop.then((value) => widget.audioPlayer.pause());
  }

  void _launchQuestionGallery() {
    widget.audioPlayer.run();
    final onPop = Navigator.of(context)
        .push(MaterialPageRoute<void>(builder: (_) => QuestionGallery(bm)));
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
