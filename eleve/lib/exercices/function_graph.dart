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
        painter: _AxisPainter(metrics),
      ),
      CustomPaint(
        size: Size(metrics.canvasWidth, metrics.canvasHeight),
        painter: _FuncPainter(metrics, function.label, function.graph.segments),
      ),
    ]);
  }
}

class _FuncPainter extends CustomPainter {
  final RepereMetrics metrics;
  final String label;
  final List<BezierCurve> segments;
  _FuncPainter(this.metrics, this.label, this.segments);

  @override
  void paint(Canvas canvas, Size size) {
    if (segments.isEmpty) {
      return;
    }

    final paint = Paint()
      ..color = Colors.purple
      ..strokeWidth = 2
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

    DrawingsPainter.paintText(
        metrics, canvas, LabeledPoint(segments[5].p0, LabelPos.top), label,
        color: Colors.purple);
  }

  @override
  bool shouldRepaint(covariant _FuncPainter oldDelegate) {
    return metrics != oldDelegate.metrics || segments != oldDelegate.segments;
  }
}

class _AxisPainter extends CustomPainter {
  final RepereMetrics metrics;
  _AxisPainter(this.metrics);

  @override
  void paint(Canvas canvas, Size size) {
    final paint = Paint()
      ..color = Colors.black
      ..strokeWidth = 1
      ..style = PaintingStyle.stroke;

    final visualOrigin = metrics.logicalToVisual(const Coord(0, 0));
    canvas.drawLine(Offset(visualOrigin.dx, 0),
        Offset(visualOrigin.dx, size.height), paint);
    canvas.drawLine(
        Offset(0, visualOrigin.dy), Offset(size.width, visualOrigin.dy), paint);
  }

  @override
  bool shouldRepaint(covariant _AxisPainter oldDelegate) {
    return metrics != oldDelegate.metrics;
  }
}
