import 'package:eleve/trivialpoursuit/game.dart';
import 'package:flutter/material.dart';

void main() {
  runApp(const MyApp());
}

final Color darkBlue = const Color.fromARGB(255, 18, 32, 47);

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
          ),
          // Outer white container with padding
          body: Game()),
    );
  }
}
