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
            actions: [
              TextButton(
                onPressed: () {},
                child: const Text("DEBUG"),
              )
            ],
          ),
          body: const Center(
            child: GameController(60),
          )),
    );
  }
}
