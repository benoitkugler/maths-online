import 'dart:math';

import 'package:material_ui/material_ui.dart';

class Congrats extends StatefulWidget {
  /// [allCorrect] is true if every question was correctly answered.
  final bool allCorrect;

  /// [mark] and [bareme] are the exercice score, displayed when
  /// there are some mistakes.
  final int mark;
  final int bareme;

  const Congrats({
    super.key,
    required this.allCorrect,
    required this.mark,
    required this.bareme,
  });

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
    if (!mounted) return;
    setState(() {
      opacity = 1;
    });
  }

  @override
  Widget build(BuildContext context) {
    final allCorrect = widget.allCorrect;
    return Card(
      elevation: 8,
      shadowColor: Colors.white,
      color: allCorrect ? Colors.lightGreen.shade500 : Colors.orange.shade300,
      child: Padding(
        padding: const EdgeInsets.all(8.0),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Text(
              allCorrect ? "Exercice terminé. \nBravo !" : "Exercice terminé",
              style: const TextStyle(fontSize: 24),
              textAlign: TextAlign.center,
            ),
            const SizedBox(height: 20),
            if (allCorrect)
              AnimatedOpacity(
                duration: const Duration(seconds: 3),
                opacity: opacity,
                child: Image.asset(
                  "assets/images/confetti-icon.png",
                  width: 120,
                ),
              )
            else ...[
              const Text(
                "Dommage, certaines réponses sont incorrectes.",
                textAlign: TextAlign.center,
              ),
              const SizedBox(height: 12),
              Text(
                "${widget.mark} / ${widget.bareme}",
                style: const TextStyle(
                  fontSize: 24,
                  fontWeight: FontWeight.bold,
                ),
              ),
            ],
            const SizedBox(height: 20),
            OutlinedButton(
              onPressed: () => Navigator.of(context).pop(true),
              style: OutlinedButton.styleFrom(foregroundColor: Colors.black54),
              child: const Text("Retour au sommaire"),
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

String pickCongratulationMessage() {
  return congratsMessages[Random().nextInt(congratsMessages.length)];
}

class CorrectAnswerDialog extends StatelessWidget {
  final void Function() onValid;
  const CorrectAnswerDialog(this.onValid, {super.key});

  @override
  Widget build(BuildContext context) {
    return AlertDialog(
      backgroundColor: Colors.lightGreen.shade500,
      title: Text(pickCongratulationMessage()),
      content: const Text("Ta réponse est correcte, bravo !"),
      actions: [
        OutlinedButton(
          onPressed: onValid,
          style: OutlinedButton.styleFrom(foregroundColor: Colors.black54),
          child: const Text("Continuer l'exercice"),
        ),
      ],
    );
  }
}
