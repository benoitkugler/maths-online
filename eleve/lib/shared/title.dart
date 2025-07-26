import 'package:flutter/material.dart';

class ColoredTitle extends StatelessWidget {
  final String title;
  final Color color;
  const ColoredTitle(this.title, this.color, {super.key});

  @override
  Widget build(BuildContext context) {
    return Card(
        elevation: 4,
        shadowColor: color,
        child: DecoratedBox(
          decoration: BoxDecoration(
            border: Border.all(color: color),
            borderRadius: BorderRadius.circular(4),
          ),
          child: Padding(
            padding: const EdgeInsets.all(8.0),
            child: Text(
              title,
              textAlign: TextAlign.center,
              style: TextStyle(fontSize: 22),
            ),
          ),
        ));
  }
}
