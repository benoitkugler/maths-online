import 'package:flutter/material.dart';

class Congrats extends StatelessWidget {
  const Congrats({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Card(
      shadowColor: Colors.green,
      color: Colors.greenAccent.shade400,
      child: const Padding(
        padding: EdgeInsets.all(8.0),
        child: Text(
          "Exercice terminé. \nBravo !",
          style: TextStyle(fontSize: 24),
          textAlign: TextAlign.center,
        ),
      ),
    );
  }
}
