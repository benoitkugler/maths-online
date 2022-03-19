import 'package:eleve/exercices/repere.gen.dart';
import 'package:flutter/material.dart';

class Repere extends StatelessWidget {
  final Figure spec;
  const Repere(this.spec, {Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    final displayLength = MediaQuery.of(context).size.shortestSide * 0.7;

    return Container(
      decoration: BoxDecoration(
          color: Colors.white.withOpacity(0.7),
          boxShadow: [BoxShadow(color: Colors.white, blurRadius: 5)]),
      child: CustomPaint(
        size: Size(displayLength * spec.width / 1000,
            displayLength * spec.height / 1000),
        painter: _ReperePainter(spec, displayLength),
      ),
    );
  }
}

extension _OffsetLabel on LabelPos {
  Offset offset(double textWidth, double textHeight) {
    switch (this) {
      case LabelPos.top:
        return Offset(-textWidth / 2, -textHeight);
      case LabelPos.bottom:
        return Offset(-textWidth / 2, 0);
      case LabelPos.left:
        return Offset(-textWidth, -textHeight / 2);
      case LabelPos.right: // nothing to do
        return Offset(0, -textHeight / 2);
      case LabelPos.topLeft:
        return Offset(-textWidth, -textHeight);
      case LabelPos.topRight:
        return Offset(0, -textHeight);
      case LabelPos.bottomLeft:
        return Offset(-textWidth, 0);
      case LabelPos.bottomRight:
        return const Offset(0, 0);
    }
  }
}

class _ReperePainter extends CustomPainter {
  final Figure spec;
  final double displayLength;

  _ReperePainter(this.spec, this.displayLength);

  /// also adjust from math convention to flutter convention
  Offset scale(Coord point, double canvasHeight) {
    return Offset(displayLength * point.x / 1000,
        canvasHeight - displayLength * point.y / 1000);
  }

  void _paintText(Canvas canvas, Size size, LabeledPoint point, String name,
      Color color, FontWeight weight) {
    final pt = TextPainter(
        text: TextSpan(
          text: " " + name + " ",
          style: TextStyle(fontWeight: weight, color: color),
        ),
        textDirection: TextDirection.ltr);
    pt.layout();

    final textWidth = pt.width;
    final textHeight = pt.height;

    final originalPos = scale(point.point, size.height);

    final offset = point.pos.offset(textWidth, textHeight);

    pt.paint(canvas, originalPos.translate(offset.dx, offset.dy));
  }

  void _paintPoint(Canvas canvas, Size size, LabeledPoint point, String name) {
    _paintText(canvas, size, point, name, Colors.white, FontWeight.w900);
    _paintText(
        canvas, size, point, name, Colors.blue.shade800, FontWeight.normal);
  }

  void _paintPoints(Canvas canvas, Size size) {
    spec.points.forEach((key, value) {
      _paintPoint(canvas, size, value, key);
    });
  }

  void _paintLine(Canvas canvas, Size size, Line line) {
    final from = spec.points[line.from]!.point;
    final to = spec.points[line.to]!.point;
    canvas.drawLine(scale(from, size.height), scale(to, size.height), Paint());

    if (line.labelName.isNotEmpty) {
      _paintPoint(
          canvas,
          size,
          LabeledPoint(
              Coord((from.x + to.x) ~/ 2, (from.y + to.y) ~/ 2), line.labelPos),
          line.labelName);
    }
  }

  void _paintLines(Canvas canvas, Size size) {
    for (var element in spec.lines) {
      _paintLine(canvas, size, element);
    }
  }

  @override
  void paint(Canvas canvas, Size size) {
    _paintLines(canvas, size);

    _paintPoints(canvas, size);
  }

  @override
  bool shouldRepaint(covariant CustomPainter oldDelegate) {
    return false;
  }
}
