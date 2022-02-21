import 'package:flutter/material.dart';
import 'package:flutter_math_fork/flutter_math.dart';

class Exercice extends StatefulWidget {
  const Exercice({Key? key}) : super(key: key);

  @override
  _ExerciceState createState() => _ExerciceState();
}

class _ExerciceState extends State<Exercice> {
  @override
  Widget build(BuildContext context) {
    return Column(children: [
      const Text("Factoriser l'expression suivante :"),
      Math.tex(
        r"f(x) = u_{n+1} * 3 + \frac{2x - 5}{4}",
        textStyle: const TextStyle(fontSize: 20),
      ),
      const TextField(),
      TextButton(onPressed: () {}, child: const Text("Valider"))
    ]);
  }
}
