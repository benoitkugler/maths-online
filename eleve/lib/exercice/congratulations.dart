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
      shadowColor: Colors.green,
      color: Colors.greenAccent.shade200,
      child: Padding(
        padding: const EdgeInsets.all(8.0),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Text(
              "Exercice terminé. \nBravo !",
              style: TextStyle(fontSize: 24, color: Colors.green.shade700),
              textAlign: TextAlign.center,
            ),
            const SizedBox(height: 20),
            AnimatedOpacity(
              duration: const Duration(seconds: 3),
              opacity: opacity,
              child: Image.asset(
                "lib/images/confetti-icon.png",
                width: 120,
              ),
            ),
          ],
        ),
      ),
    );
  }
}
