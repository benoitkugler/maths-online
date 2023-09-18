import 'dart:math';

import 'package:flutter/material.dart';

class Congrats extends StatefulWidget {
  const Congrats({Key? key}) : super(key: key);

  @override
  State<Congrats> createState() => _CongratsState();
}

class _CongratsState extends State<Congrats> {
  double opacity = 0;

  @override
  void initState() {
    _triggerAnimation();
    super.initState();
  }

  void _triggerAnimation() async {
    await Future<void>.delayed(const Duration(milliseconds: 100));
    setState(() {
      opacity = 1;
    });
  }

  @override
  Widget build(BuildContext context) {
    return Card(
      elevation: 8,
      shadowColor: Colors.white,
      color: Colors.lightGreen.shade200,
      child: Padding(
        padding: const EdgeInsets.all(8.0),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Text(
              "Exercice terminé. \n\nBravo !",
              style: TextStyle(fontSize: 24, color: Colors.green.shade700),
              textAlign: TextAlign.center,
            ),
            const SizedBox(height: 20),
            AnimatedOpacity(
              duration: const Duration(seconds: 3),
              opacity: opacity,
              child: Image.asset(
                "assets/images/confetti-icon.png",
                width: 120,
              ),
            ),
          ],
        ),
      ),
    );
  }
}

const congratsMessages = [
  "Yes !",
  "Félicitations !",
  "Bravissimo !",
  "Oui oui oui !",
  "Mais comment fais-tu ??",
  "Tu n'es pas là pour beurrer des tartines",
  "Tu n'es pas là pour enfiler des perles",
  "C'est qui le meilleur ? ",
  "The best for ever...",
];

String _randCongrats() {
  return congratsMessages[Random().nextInt(congratsMessages.length)];
}

class CorrectAnswerDialog extends StatelessWidget {
  final void Function() onValid;
  const CorrectAnswerDialog(this.onValid, {super.key});

  @override
  Widget build(BuildContext context) {
    return AlertDialog(
      backgroundColor: Colors.lightGreen.shade400,
      title: Text(_randCongrats()),
      content: const Text("Ta réponse est correcte, bravo !"),
      actions: [
        OutlinedButton(
            onPressed: onValid,
            style: OutlinedButton.styleFrom(foregroundColor: Colors.black54),
            child: const Text("Continuer l'exercice"))
      ],
    );
  }
}
