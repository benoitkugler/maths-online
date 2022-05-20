import 'dart:math';

import 'package:eleve/trivialpoursuit/events.gen.dart';
import 'package:flutter/material.dart';

import "categories.dart";

/// Pie displays the current sucesses of the player,
/// using a pie chart.
class Pie extends StatelessWidget {
  static const size = 60.0;

  final double glowWidth;
  final Success success;
  final Color backgroundColor;
  const Pie(this.glowWidth, this.success,
      {Key? key, this.backgroundColor = Colors.white})
      : super(key: key);

  static RawMaterialButton asButton(
      void Function() onPressed, double glowWidth, Success success) {
    return RawMaterialButton(
      onPressed: onPressed,
      elevation: 2.0,
      child: Pie(glowWidth, success),
      padding: const EdgeInsets.all(5.0),
      shape: const CircleBorder(),
    );
  }

  @override
  Widget build(BuildContext context) {
    final nbSections = Categorie.values.length;
    final angularSection = 2 * pi / nbSections;

    return Stack(
      children: [
        Container(
          width: size,
          height: size,
          decoration: BoxDecoration(shape: BoxShape.circle, boxShadow: [
            BoxShadow(
              color: backgroundColor,
              blurRadius: glowWidth,
            )
          ]),
        ),
        ...Categorie.values.map((cat) => AnimatedRotation(
              turns: success[cat.index] ? 3 : 0.5,
              curve: Curves.easeOut,
              duration: const Duration(milliseconds: 3000),
              child: AnimatedScale(
                duration: const Duration(milliseconds: 3000),
                curve: Curves.easeOut,
                scale: success[cat.index] ? 1 : 1.5,
                child: AnimatedOpacity(
                  duration: const Duration(milliseconds: 3000),
                  curve: Curves.easeOut,
                  opacity: success[cat.index] ? 1 : 0,
                  child: _PieFraction(
                      cat.color, (cat.index + 0.5) * angularSection, size),
                ),
              ),
            )),
        CustomPaint(
          size: const Size(size, size),
          painter: _PieBackgroundPainter(),
        ),
      ],
    );
  }
}

class _PieBackgroundPainter extends CustomPainter {
  static const radiusRatio = 0.45;
  _PieBackgroundPainter();

  @override
  void paint(Canvas canvas, Size size) {
    final radius = size.shortestSide * radiusRatio;
    final nbSections = Categorie.values.length;
    final angularSection = 2 * pi / nbSections;
    final center = size.center(Offset.zero);

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
  bool shouldRepaint(covariant _PieBackgroundPainter oldDelegate) {
    return false;
  }
}

class _PieFraction extends StatelessWidget {
  final Color color;
  final double angle;
  final double size;

  const _PieFraction(this.color, this.angle, this.size, {Key? key})
      : super(key: key);

  @override
  Widget build(BuildContext context) {
    return CustomPaint(
      size: Size(size, size),
      painter: _PiePartPainter(color, angle),
    );
  }
}

class _PiePartPainter extends CustomPainter {
  final Color color;
  final double angle;

  _PiePartPainter(this.color, this.angle);

  @override
  void paint(Canvas canvas, Size size) {
    final radius = size.shortestSide * _PieBackgroundPainter.radiusRatio;
    final nbSections = Categorie.values.length;
    final angularSection = 2 * pi / nbSections;
    final center = size.center(Offset.zero);

    final arcRect = Rect.fromCircle(center: center, radius: radius);
    canvas.drawArc(
        arcRect,
        angle - angularSection / 2,
        angularSection,
        true,
        Paint()
          ..style = PaintingStyle.stroke
          ..color = Colors.black.withOpacity(0.5)
          ..strokeWidth = 1);
    canvas.drawArc(
        arcRect,
        angle - angularSection / 2,
        angularSection,
        true,
        Paint()
          ..style = PaintingStyle.fill
          ..color = color);
  }

  @override
  bool shouldRepaint(covariant _PiePartPainter oldDelegate) {
    return color != oldDelegate.color || angle != oldDelegate.angle;
  }
}
