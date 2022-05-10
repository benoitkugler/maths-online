import 'package:eleve/exercices/repere.dart';
import 'package:eleve/exercices/repere.gen.dart';
import 'package:eleve/exercices/types.gen.dart';
import 'package:flutter/material.dart';

class FunctionGraphW extends StatelessWidget {
  final FunctionGraphBlock function;

  const FunctionGraphW(this.function, {Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    final metrics = RepereMetrics(function.graph.bounds, context);
    return BaseRepere(metrics, true, [
      CustomPaint(
        size: Size(metrics.canvasWidth, metrics.canvasHeight),
        painter: BezierCurvesPainter(metrics, function.graph.functions),
      ),
    ]);
  }
}

class BezierCurvesPainter extends CustomPainter {
  final RepereMetrics metrics;
  final List<FunctionGraph> functions;
  BezierCurvesPainter(this.metrics, this.functions);

  void _paintFunction(Canvas canvas, Size size, FunctionGraph fn) {
    if (fn.segments.isEmpty) {
      return;
    }

    final color = fromHex(fn.decoration.color);
    final paint = Paint()
      ..color = color
      ..strokeWidth = 1
      ..style = PaintingStyle.stroke;

    final path = Path();

    final start = metrics.logicalToVisual(fn.segments[0].p0);
    path.moveTo(start.dx, start.dy);
    for (var segment in fn.segments) {
      final p1 = metrics.logicalToVisual(segment.p1);
      final p2 = metrics.logicalToVisual(segment.p2);
      path.quadraticBezierTo(p1.dx, p1.dy, p2.dx, p2.dy);
    }
    canvas.drawPath(path, paint);

    final labelIndex = fn.segments.length * 3 ~/ 4;
    final labelPos = fn.segments[labelIndex].p0;
    // adjust the position based on space available
    final visualLabelPos = metrics.logicalToVisual(labelPos);
    final putTop = (visualLabelPos.dy > size.height / 2);
    final putLeft = (visualLabelPos.dx > size.width / 2);
    LabelPos pos;
    if (putTop) {
      pos = putLeft ? LabelPos.topLeft : LabelPos.topRight;
    } else {
      pos = putLeft ? LabelPos.bottomLeft : LabelPos.bottomRight;
    }

    DrawingsPainter.paintText(
        metrics,
        canvas,
        LabeledPoint(Coord(labelPos.x, labelPos.y + 1), pos),
        fn.decoration.label,
        color: color);
  }

  @override
  void paint(Canvas canvas, Size size) {
    for (var element in functions) {
      _paintFunction(canvas, size, element);
    }
  }

  @override
  bool shouldRepaint(covariant BezierCurvesPainter oldDelegate) {
    return metrics != oldDelegate.metrics || functions != oldDelegate.functions;
  }
}
