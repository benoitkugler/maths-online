import 'dart:math';

import 'package:eleve/questions/repere.dart';
import 'package:eleve/types/src_maths_functiongrapher.dart';
import 'package:eleve/types/src_maths_questions_client.dart';
import 'package:eleve/types/src_maths_repere.dart';
import 'package:flutter/material.dart';

class FunctionsGraphW extends StatelessWidget {
  final FunctionsGraphBlock graphs;

  const FunctionsGraphW(this.graphs, {Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    final metrics = RepereMetrics(graphs.bounds, context);
    final painter = BezierCurvesPainter(metrics,
        functions: graphs.functions,
        sequences: graphs.sequences,
        areas: graphs.areas,
        points: graphs.points);
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

class BezierCurvesPainter extends CustomPainterText {
  final RepereMetrics metrics;
  final List<FunctionGraph> functions;
  final List<SequenceGraph> sequences;
  final List<FunctionArea> areas;
  final List<FunctionPoint> points;
  BezierCurvesPainter(
    this.metrics, {
    this.functions = const [],
    this.sequences = const [],
    this.areas = const [],
    this.points = const [],
  });

  /// [extractTexts] returns the positionned text that
  /// must be diplayed, to be handled in a separate widget.
  @override
  List<PositionnedText> extractTexts() {
    final out = <PositionnedText>[];
    for (var fn in functions) {
      final text = _functionText(fn);
      if (text != null) {
        out.add(text);
      }
    }
    for (var fn in sequences) {
      final text = _sequenceText(fn);
      if (text != null) {
        out.add(text);
      }
    }
    for (var point in points) {
      final text = _pointText(point);
      if (text != null) {
        out.add(text);
      }
    }
    return out;
  }

  // returns null if the function has an empty label
  // or has no segments
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

  PositionnedText? _sequenceText(SequenceGraph seq) {
    final color = fromHex(seq.decoration.color);
    if (seq.points.isEmpty) return null;
    // display the label at the right of the last point
    final last = seq.points.last;
    return PositionnedText(seq.decoration.label,
        PosPoint(Coord(last.x + 0.3, last.y), LabelPos.right),
        color: color);
  }

  // returns null if the legend is empty
  PositionnedText? _pointText(FunctionPoint point) {
    if (point.legend.isEmpty) {
      return null;
    }

    // adjust the position based on space available
    final visualLabelPos = metrics.logicalToVisual(point.coord);
    final putBottom = (visualLabelPos.dy >
        metrics.size.height / 2); // visual y grows from the top
    final putLeft = (visualLabelPos.dx > metrics.size.width / 2);
    LabelPos pos;
    if (putBottom) {
      pos = putLeft ? LabelPos.bottomLeft : LabelPos.bottomRight;
    } else {
      pos = putLeft ? LabelPos.topLeft : LabelPos.topRight;
    }
    return PositionnedText(point.legend, PosPoint(point.coord, pos),
        color: fromHex(point.color));
  }

  Path _buildPath(List<BezierCurve> segments) {
    final path = Path();
    final start = metrics.logicalToVisual(segments[0].p0);
    path.moveTo(start.dx, start.dy);
    for (var segment in segments) {
      final p1 = metrics.logicalToVisual(segment.p1);
      final p2 = metrics.logicalToVisual(segment.p2);
      path.quadraticBezierTo(p1.dx, p1.dy, p2.dx, p2.dy);
    }
    return path;
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

    final path = _buildPath(fn.segments);
    canvas.drawPath(path, paint);
  }

  void _paintSequence(Canvas canvas, SequenceGraph fn) {
    final color = fromHex(fn.decoration.color);
    final paint = Paint()
      ..style = PaintingStyle.stroke
      ..strokeWidth = 1.5
      ..color = color;

    for (var point in fn.points) {
      final center = metrics.logicalToVisual(point);
      // draw a rotated cross
      canvas.save();
      canvas.translate(center.dx, center.dy);
      canvas.rotate(pi / 4);
      canvas.drawLine(const Offset(-5, 0), const Offset(5, 0), paint);
      canvas.drawLine(const Offset(0, -5), const Offset(0, 5), paint);
      canvas.restore();
    }
  }

  void _paintArea(Canvas canvas, FunctionArea area) {
    if (area.path.isEmpty) {
      return;
    }

    final color = fromHex(area.color);
    final path = _buildPath(area.path);
    canvas.drawPath(
        path,
        Paint()
          ..style = PaintingStyle.fill
          ..color = color);
  }

  void _paintPoint(Canvas canvas, FunctionPoint point) {
    final color = fromHex(point.color);
    final center = metrics.logicalToVisual(point.coord);
    canvas.drawCircle(
        center,
        2,
        Paint()
          ..style = PaintingStyle.fill
          ..color = color);
  }

  @override
  void paint(Canvas canvas, Size size) {
    for (var area in areas) {
      _paintArea(canvas, area);
    }
    for (var element in functions) {
      _paintFunction(canvas, element);
    }
    for (var element in sequences) {
      _paintSequence(canvas, element);
    }
    // paint the point on top of the rest
    for (var point in points) {
      _paintPoint(canvas, point);
    }
  }

  @override
  bool shouldRepaint(covariant BezierCurvesPainter oldDelegate) {
    return metrics != oldDelegate.metrics ||
        functions != oldDelegate.functions ||
        sequences != oldDelegate.sequences;
  }
}
