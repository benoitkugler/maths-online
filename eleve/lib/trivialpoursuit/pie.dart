import 'dart:math';
import 'dart:ui';

import 'package:eleve/trivialpoursuit/events.gen.dart';
import 'package:flutter/foundation.dart';
import 'package:flutter/material.dart';

/// Pie displays the current sucesses of the player,
/// using a pie chart.
class Pie extends StatelessWidget {
  final Success success;
  const Pie(this.success, {Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.only(left: 40),
      child: CustomPaint(
        size: const Size(40, 40),
        painter: _PiePainter(success),
      ),
    );
  }
}

class _PiePainter extends CustomPainter {
  final Success success;

  /// map question categories to colors
  static const categoriesColors = [
    Colors.purple,
    Colors.green,
    Colors.orange,
    Colors.yellow,
    Colors.blue,
  ];
  static const blurRadius = 6.0;

  _PiePainter(this.success);

  @override
  void paint(Canvas canvas, Size size) {
    final radius = size.shortestSide;
    final nbSections = categoriesColors.length;
    final angularSection = 2 * pi / nbSections;
    final center = size.center(Offset.zero);

    // outer glow
    canvas.drawCircle(
      center,
      radius + 2,
      Paint()
        ..color = Colors.white
        ..style = PaintingStyle.stroke
        ..strokeWidth = 10
        ..imageFilter =
            ImageFilter.blur(sigmaX: blurRadius, sigmaY: blurRadius),
    );

    // white background
    canvas.drawCircle(
        center,
        radius,
        Paint()
          ..color = Colors.white.withOpacity(0.8)
          ..style = PaintingStyle.fill);

    final arcRect = Rect.fromCircle(center: center, radius: radius);
    for (var i = 0; i < nbSections; i++) {
      canvas.drawArc(
          arcRect,
          i * angularSection,
          angularSection,
          true,
          Paint()
            ..style = PaintingStyle.stroke
            ..color = Colors.black.withOpacity(0.5)
            ..strokeWidth = 1);
      if (success[i]) {
        canvas.drawArc(
            arcRect,
            i * angularSection,
            angularSection,
            true,
            Paint()
              ..style = PaintingStyle.fill
              ..color = categoriesColors[i].withOpacity(0.8));
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
