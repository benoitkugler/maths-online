import 'package:eleve/trivialpoursuit/game.dart';
import 'package:flutter/material.dart';

void main() {
  runApp(const MyApp());
}

const Color darkBlue = Color.fromARGB(255, 27, 54, 82);

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
        body: const _HomePage(),
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
        builder: (_) => const Scaffold(body: TrivialPoursuitController(60))));
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
