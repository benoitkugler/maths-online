import 'dart:math';

import 'package:flutter/material.dart';

const _variants = [
  "Pour réaliser une bonne recette de mathématiques, prendre un zeste de _stylo_, une _demi-feuille_ de brouillon et une _calculatrice_ bien chargée, et manier le tout une bonne quinzaine de minutes.",
  "As-tu bien ton triplé _Stylo-Feuille-Calculette_ pour tenter de scorer ?",
  "STOP. Vos _papiers_ s'il vous plait... mes _calculs_ sont formels : sans _stylo_, vous êtes en infraction au code des maths !",
];

String _pickVariant() {
  final index = Random().nextInt(_variants.length);
  return _variants[index];
}

class _Text extends StatelessWidget {
  final String text;
  const _Text(this.text, {Key? key}) : super(key: key);

  static List<_TextRun> split(String text, {String separator = "_"}) {
    var isSpecial = true;
    return text.split(separator).map((e) {
      isSpecial = !isSpecial;
      return _TextRun(e, isSpecial);
    }).toList();
  }

  @override
  Widget build(BuildContext context) {
    final parts = split(text);
    return Text.rich(
      TextSpan(
          children: parts
              .map((e) => TextSpan(
                    text: e.text,
                    style: TextStyle(
                        fontWeight:
                            e.isSpecial ? FontWeight.bold : FontWeight.normal),
                  ))
              .toList()),
      style: const TextStyle(fontSize: 22, shadows: [
        Shadow(
          color: Colors.lightBlueAccent,
          offset: Offset(0, 1),
          blurRadius: 6,
        ),
      ]),
    );
  }
}

class _TextRun {
  final String text;
  final bool isSpecial;
  const _TextRun(this.text, this.isSpecial);
}

class ActivityStart extends StatelessWidget {
  final void Function() onContinue;
  const ActivityStart(this.onContinue, {Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Column(
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: [
          Flexible(
            flex: 2,
            child: Image.asset(
              "assets/images/cahier-stylo-calc.png",
            ),
          ),
          Expanded(
            child: Center(
              child: Padding(
                padding: const EdgeInsets.all(16.0),
                child: _Text(_pickVariant()),
              ),
            ),
          ),
          Padding(
            padding: const EdgeInsets.symmetric(horizontal: 30, vertical: 30),
            child: DecoratedBox(
              decoration: const BoxDecoration(boxShadow: [
                BoxShadow(
                  color: Colors.lightGreen,
                  spreadRadius: 2,
                  blurRadius: 10,
                ),
              ]),
              child: ElevatedButton(
                onPressed: onContinue,
                style: ElevatedButton.styleFrom(
                    backgroundColor: Colors.lightGreen,
                    padding: const EdgeInsets.all(16)),
                child: const Text(
                  "J'ai mon matériel !",
                  style: TextStyle(fontSize: 18),
                ),
              ),
            ),
          ),
        ],
      ),
    );
  }
}
