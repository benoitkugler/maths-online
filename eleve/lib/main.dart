import 'dart:convert';

import 'package:eleve/build_mode.dart';
import 'package:eleve/exercices/types.gen.dart';
import 'package:eleve/trivialpoursuit/board.dart';
import 'package:eleve/trivialpoursuit/game.dart';
import 'package:eleve/trivialpoursuit/login.dart';
import 'package:flutter/material.dart';

const Color darkBlue = Color.fromARGB(255, 27, 54, 82);

// final bm = buildMode();
final bm = BuildMode.debug;

void main() {
  runApp(const MyApp());

  // start music as background
  // final player = AudioCache(prefix: "lib/music/");
  // player.loop("DontLetMeGo.mp3");
}

const input = """
{"Title":"Très longue question horizontale","Content":[{"Data":{"Text":"Écrire sous une seule fraction : "},"Kind":3},{"Data":{"Content":"\\\\frac{1}{3} + \\\\frac{2}{5}","IsInline":true},"Kind":2},{"Data":{"Text":"Écrire sous une seule fraction : "},"Kind":3},{"Data":{"Content":"\\\\frac{1}{3} + \\\\frac{2}{5}","IsInline":true},"Kind":2},{"Data":{"Text":"Écrire sous une seule fraction : "},"Kind":3},{"Data":{"Content":"\\\\frac{1}{3} + \\\\frac{2}{5}","IsInline":true},"Kind":2},{"Data":{"Text":"Écrire sous une seule fraction : "},"Kind":3},{"Data":{"Content":"\\\\frac{1}{3} + \\\\frac{2}{5}","IsInline":true},"Kind":2},{"Data":{"Text":"Écrire sous une seule fraction : "},"Kind":3},{"Data":{"Content":"\\\\frac{1}{3} + \\\\frac{2}{5}","IsInline":true},"Kind":2},{"Data":{"Text":"Écrire sous une seule fraction : "},"Kind":3},{"Data":{"Content":"\\\\frac{1}{3} + \\\\frac{2}{5}","IsInline":true},"Kind":2},{"Data":{"Text":"Écrire sous une seule fraction : "},"Kind":3},{"Data":{"Content":"\\\\frac{1}{3} + \\\\frac{2}{5}","IsInline":true
},"Kind":2}]}
""";

final question = clientQuestionFromJson(jsonDecode(input));

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
      home: Scaffold(
        appBar: AppBar(
          title: const Text('Isiro'),
        ),
        // body: const _HomePage(),
        body: Container(color: Colors.grey, child: Board(print, {}, 3)),
        // body: Padding(
        //   padding:
        //       const EdgeInsets.only(top: 10, bottom: 10, left: 20, right: 20),
        //   child: Question(question),
        // ),
      ),
    );
  }
}

class _HomePage extends StatefulWidget {
  const _HomePage({Key? key}) : super(key: key);

  @override
  State<_HomePage> createState() => _HomePageState();
}

class _HomePageState extends State<_HomePage> {
  void _launchTrivialPoursuit() {
    Navigator.of(context).push(MaterialPageRoute<void>(
        builder: (_) => Scaffold(body: TrivialPoursuitLoggin(bm))));
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
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  GameIcon(() => _launchTrivialPoursuit()),
                ],
              ),
            )
          ],
        ),
      ),
    );
  }
}
