import 'package:eleve/trivialpoursuit/game.dart';
import 'package:flutter/material.dart';

void main() {
  runApp(const MyApp());
}

const Color darkBlue = Color.fromARGB(255, 18, 32, 47);

class MyApp extends StatelessWidget {
  const MyApp({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Isiro',
      theme: ThemeData.dark().copyWith(scaffoldBackgroundColor: darkBlue),
      debugShowCheckedModeBanner: false,
      home: Scaffold(
          appBar: AppBar(
            title: const Text('Bienvenue !'),
            actions: [
              TextButton(
                onPressed: () {},
                child: const Text("DEBUG"),
              )
            ],
          ),
          body: const Center(
            child: GameController(),
            // child: QuestionRoute(ShowQuestion("Test", Categorie.orange)),
          )),
    );
  }
}
