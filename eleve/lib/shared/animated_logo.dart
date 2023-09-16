import 'dart:math';

import 'package:flutter/material.dart';

class AnimatedLogo extends StatefulWidget {
  final double successPercentage;

  const AnimatedLogo(this.successPercentage, {Key? key}) : super(key: key);

  @override
  State<AnimatedLogo> createState() => _AnimatedLogoState();
}

class _AnimatedLogoState extends State<AnimatedLogo>
    with SingleTickerProviderStateMixin {
  late Animation<double> animation;
  late AnimationController controller;

  @override
  initState() {
    super.initState();
    controller = AnimationController(
        duration: const Duration(milliseconds: 1500), vsync: this);
    animation = Tween(begin: 0.0, end: widget.successPercentage.clamp(0.0, 1.0))
        .chain(CurveTween(curve: Curves.bounceOut))
        .animate(controller);

    controller.forward();
  }

  @override
  void dispose() {
    controller.stop();
    controller.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return AnimatedBuilder(
        animation: animation,
        builder: (context, child) => CustomPaint(
              painter: _LogoPainter(animation.value),
            ));
  }
}

class _LogoPainter extends CustomPainter {
  final double progressPercentage; // in [0,1]
  _LogoPainter(this.progressPercentage);

  @override
  void paint(Canvas canvas, Size size) {
    const bgColor = Color(0xFF4c4c54);
    const color = Color(0xFF14dbb4);

    final s = size.shortestSide;
    final center = Offset(size.width / 2, size.height / 2);
    final circleCrownRadius = s / 30;
    final outerCircleRadius = s * 0.48;
    final innerCircleRadius = s * 0.35;
    canvas.drawCircle(center, s * 0.5, Paint()..color = Colors.white);
    canvas.drawCircle(center, outerCircleRadius, Paint()..color = bgColor);
    canvas.drawCircle(center, innerCircleRadius, Paint()..color = Colors.white);

    final painter = Paint()..color = color;
    final c1 = Offset(center.dx - s * 0.22, center.dy - s * 0.12);
    final c2 = Offset(center.dx, center.dy - s * 0.2);
    final c3 = Offset(center.dx + s * 0.22, center.dy - s * 0.12);
    canvas.drawCircle(c1, circleCrownRadius, painter);
    canvas.drawCircle(c2, circleCrownRadius, painter);
    canvas.drawCircle(c3, circleCrownRadius, painter);

    final baseWidth = s * 0.3;
    final path = Path();
    final startCrown = Offset(center.dx - baseWidth / 2, center.dy + s * 0.1);
    path.moveTo(startCrown.dx, startCrown.dy);
    path.lineTo(center.dx + baseWidth / 2, center.dy + s * 0.1);
    path.lineTo(c3.dx, c3.dy);
    final midCrown = (c1.dy + startCrown.dy) / 2;
    path.quadraticBezierTo(
        (c2.dx + c3.dx) / 2, midCrown + s * 0.1, c2.dx, c2.dy);
    path.quadraticBezierTo(
        (c1.dx + c2.dx) / 2, midCrown + s * 0.1, c1.dx, c1.dy);
    path.close();

    canvas.drawPath(path, painter);

    canvas.drawRRect(
        RRect.fromRectAndRadius(
            Rect.fromCenter(
                center: Offset(center.dx, center.dy + s * 0.18),
                width: baseWidth,
                height: s * 0.05),
            const Radius.circular(4)),
        painter);

    final mathCoventionAngle = -pi / 2 + progressPercentage * 2 * pi;

    final outerEndProgress = center.translate(
        outerCircleRadius * cos(mathCoventionAngle),
        outerCircleRadius * sin(mathCoventionAngle));

    final progressPath = Path();
    progressPath.moveTo(center.dx, s * 0.15);
    progressPath.arcTo(
        Rect.fromCircle(center: center, radius: innerCircleRadius),
        -pi / 2,
        progressPercentage * 2 * pi,
        false);
    progressPath.arcTo(
        Rect.fromCircle(center: center, radius: outerCircleRadius),
        mathCoventionAngle,
        -progressPercentage * 2 * pi,
        false);
    progressPath.close();
    canvas.drawPath(progressPath, painter);

    final normal = (outerEndProgress - center);
    final progressEndCenter = center +
        normal / normal.distance * (outerCircleRadius + innerCircleRadius) / 2;
    canvas.drawCircle(
        Offset(center.dx,
            (s / 2 - outerCircleRadius + s / 2 - innerCircleRadius) / 2),
        (outerCircleRadius - innerCircleRadius) / 2,
        painter);
    canvas.drawCircle(progressEndCenter,
        (outerCircleRadius - innerCircleRadius) / 2, painter);
  }

  @override
  bool shouldRepaint(covariant _LogoPainter oldDelegate) {
    return oldDelegate.progressPercentage != progressPercentage;
  }
}
