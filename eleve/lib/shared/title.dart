import 'package:flutter/material.dart';

class ColoredTitle extends StatelessWidget {
  final String title;
  final Color color;
  const ColoredTitle(this.title, this.color, {Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    final shadows = [
      Shadow(
          color: color.withValues(alpha: 0.9),
          offset: const Offset(1, -2),
          blurRadius: 1.3)
    ];
    return DecoratedBox(
      decoration: BoxDecoration(
        border: Border.all(color: color),
        borderRadius: BorderRadius.circular(4),
      ),
      child: Padding(
        padding: const EdgeInsets.all(8.0),
        child: Text(
          title,
          textAlign: TextAlign.center,
          style: TextStyle(fontSize: 22, shadows: shadows),
        ),
      ),
    );
  }
}
