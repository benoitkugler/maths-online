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
        painter: BezierCurvesPainter(
            metrics, function.label, function.graph.segments),
      ),
    ]);
  }
}

class BezierCurvesPainter extends CustomPainter {
  final RepereMetrics metrics;
  final String label;
  final List<BezierCurve> segments;
  BezierCurvesPainter(this.metrics, this.label, this.segments);

  @override
  void paint(Canvas canvas, Size size) {
    if (segments.isEmpty) {
      return;
    }

    final paint = Paint()
      ..color = Colors.purple
      ..strokeWidth = 1
      ..style = PaintingStyle.stroke;

    final path = Path();

    final start = metrics.logicalToVisual(segments[0].p0);
    path.moveTo(start.dx, start.dy);
    for (var segment in segments) {
      final p1 = metrics.logicalToVisual(segment.p1);
      final p2 = metrics.logicalToVisual(segment.p2);
      path.quadraticBezierTo(p1.dx, p1.dy, p2.dx, p2.dy);
    }
    canvas.drawPath(path, paint);

    // TODO: better choose the label position to avoid painting outside
    // the figure
    final labelIndex = segments.length * 3 ~/ 4;
    final labelPos = segments[labelIndex].p0;
    DrawingsPainter.paintText(metrics, canvas,
        LabeledPoint(Coord(labelPos.x, labelPos.y + 1), LabelPos.top), label,
        color: Colors.purple);
  }

  @override
  bool shouldRepaint(covariant BezierCurvesPainter oldDelegate) {
    return metrics != oldDelegate.metrics || segments != oldDelegate.segments;
  }
}
