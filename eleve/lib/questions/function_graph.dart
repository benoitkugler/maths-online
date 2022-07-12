import 'package:eleve/questions/repere.dart';
import 'package:eleve/questions/repere.gen.dart';
import 'package:eleve/questions/types.gen.dart';
import 'package:flutter/material.dart';

class FunctionGraphW extends StatelessWidget {
  final FunctionGraphBlock function;

  const FunctionGraphW(this.function, {Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    final metrics = RepereMetrics(function.graph.bounds, context);
    final painter = BezierCurvesPainter(metrics, function.graph.functions);
    final texts = painter.extractTexts();
    return BaseRepere(
      metrics,
      true,
      true,
      [
        CustomPaint(
          size: metrics.size,
          painter: painter,
        ),
      ],
      texts,
    );
  }
}

class BezierCurvesPainter extends CustomPainter {
  final RepereMetrics metrics;
  final List<FunctionGraph> functions;
  BezierCurvesPainter(this.metrics, this.functions);

  List<PositionnedText> extractTexts() {
    final out = <PositionnedText>[];
    for (var fn in functions) {
      final text = _functionText(fn);
      if (text != null) {
        out.add(text);
      }
    }
    return out;
  }

  PositionnedText? _functionText(FunctionGraph fn) {
    if (fn.segments.isEmpty || fn.decoration.label.isEmpty) {
      return null;
    }

    final labelIndex = fn.segments.length * 3 ~/ 4;
    final labelPos = fn.segments[labelIndex].p0;
    // adjust the position based on space available
    final visualLabelPos = metrics.logicalToVisual(labelPos);
    final putBottom = (visualLabelPos.dy >
        metrics.size.height / 2); // visual y grows from the top
    final putLeft = (visualLabelPos.dx > metrics.size.width / 2);
    LabelPos pos;
    if (putBottom) {
      pos = putLeft ? LabelPos.bottomLeft : LabelPos.bottomRight;
    } else {
      pos = putLeft ? LabelPos.topLeft : LabelPos.topRight;
    }
    return PositionnedText(
        fn.decoration.label, PosPoint(Coord(labelPos.x, labelPos.y + 1), pos),
        color: fromHex(fn.decoration.color));
  }

  void _paintFunction(Canvas canvas, FunctionGraph fn) {
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
  }

  @override
  void paint(Canvas canvas, Size size) {
    for (var element in functions) {
      _paintFunction(canvas, element);
    }
  }

  @override
  bool shouldRepaint(covariant BezierCurvesPainter oldDelegate) {
    return metrics != oldDelegate.metrics || functions != oldDelegate.functions;
  }
}
