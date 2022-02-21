import 'package:eleve/maths/maths.dart';
import 'package:flutter/material.dart';

void main() {
  runApp(const MyApp());
}

class MyApp extends StatelessWidget {
  const MyApp({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Maths online',
      home: Scaffold(
        appBar: AppBar(
          title: const Text('Bienvenue !'),
        ),
        body: const Padding(
          padding: EdgeInsets.all(8.0),
          child: Exercice(),
        ),
      ),
      theme: ThemeData(
        appBarTheme: const AppBarTheme(
          backgroundColor: Colors.lightGreen,
          foregroundColor: Colors.black,
        ),
        primaryColor: Colors.purple,
      ),
    );
  }
}
