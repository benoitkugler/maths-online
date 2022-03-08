import 'dart:math';

import 'package:eleve/trivialpoursuit/events.gen.dart';
import 'package:flutter/foundation.dart';
import 'package:flutter/material.dart';

/// Pie displays the current sucesses of the player,
/// using a pie chart.
class Pie extends StatelessWidget {
  final double glowWidth;
  final Success success;
  const Pie(this.glowWidth, this.success, {Key? key}) : super(key: key);

  static RawMaterialButton asButton(
      void Function() onPressed, double glowWidth, Success success) {
    return RawMaterialButton(
      onPressed: onPressed,
      elevation: 2.0,
      child: Pie(glowWidth, success),
      padding: const EdgeInsets.all(10.0),
      shape: const CircleBorder(),
    );
  }

  @override
  Widget build(BuildContext context) {
    return Stack(
      children: [
        Container(
          width: 80,
          height: 80,
          decoration: BoxDecoration(shape: BoxShape.circle, boxShadow: [
            BoxShadow(
              color: Colors.white,
              blurRadius: glowWidth,
            )
          ]),
        ),
        CustomPaint(
          size: const Size(80, 80),
          painter: _PiePainter(success),
        ),
      ],
    );
  }
}

/// map question categories to colors

extension CategorieColor on Categorie {
  Color get color {
    switch (this) {
      case Categorie.purple:
        return Colors.purple;
      case Categorie.green:
        return Colors.green;
      case Categorie.orange:
        return Colors.orange;
      case Categorie.yellow:
        return Colors.yellow;
      case Categorie.blue:
        return Colors.blue;
    }
  }
}

class _PiePainter extends CustomPainter {
  final Success success;

  _PiePainter(this.success);

  @override
  void paint(Canvas canvas, Size size) {
    final radius = size.shortestSide * 0.45;
    final nbSections = Categorie.values.length;
    final angularSection = 2 * pi / nbSections;
    final center = size.center(Offset.zero);

    // white background
    canvas.drawCircle(
        center,
        radius,
        Paint()
          ..color = Colors.white.withOpacity(0.8)
          ..style = PaintingStyle.fill);

    final arcRect = Rect.fromCircle(center: center, radius: radius);
    for (var i in Categorie.values) {
      canvas.drawArc(
          arcRect,
          i.index * angularSection,
          angularSection,
          true,
          Paint()
            ..style = PaintingStyle.stroke
            ..color = Colors.black.withOpacity(0.5)
            ..strokeWidth = 1);
      if (success[i.index]) {
        canvas.drawArc(
            arcRect,
            i.index * angularSection,
            angularSection,
            true,
            Paint()
              ..style = PaintingStyle.fill
              ..color = i.color.withOpacity(0.8));
      }
    }

    // external circle drawn on top
    canvas.drawCircle(
        center,
        radius,
        Paint()
          ..color = Colors.black
          ..style = PaintingStyle.stroke
          ..strokeWidth = 2);
  }

  @override
  bool shouldRepaint(covariant _PiePainter oldDelegate) {
    return !listEquals(oldDelegate.success, success);
  }
}
