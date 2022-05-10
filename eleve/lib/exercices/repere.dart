import 'dart:math';

import 'package:eleve/exercices/repere.gen.dart';
import 'package:eleve/exercices/types.gen.dart';
import 'package:flutter/material.dart';

class StaticRepere extends StatelessWidget {
  final Figure figure;
  const StaticRepere(this.figure, {Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    final metrics = RepereMetrics(figure.bounds, context);
    return BaseRepere(metrics, figure.showGrid, [
      // custom drawing
      CustomPaint(
        size: Size(metrics.canvasWidth, metrics.canvasHeight),
        painter: DrawingsPainter(metrics, figure.drawings),
      ),
    ]);
  }
}

class RepereMetrics {
  final double
      _displayLength; // displayLength is the length of the largest size of the figure
  final RepereBounds figure;

  RepereMetrics(this.figure, BuildContext context)
      : _displayLength = MediaQuery.of(context).size.shortestSide * 0.95;

  double get resolution => max(figure.width, figure.height).toDouble();

  double get canvasWidth => _displayLength * figure.width / resolution;
  double get canvasHeight => _displayLength * figure.height / resolution;

  Offset logicalToVisual(Coord point) {
    // shift by the origin
    final x = figure.origin.x + point.x;
    final y = figure.origin.y + point.y;

    return Offset(_displayLength * x / resolution,
        canvasHeight - _displayLength * y / resolution);
  }

  Offset logicalIntToVisual(IntCoord point) {
    return logicalToVisual(Coord(point.x.toDouble(), point.y.toDouble()));
  }

  IntCoord visualToLogical(Offset offset) {
    final x = offset.dx * resolution / _displayLength;
    final y = -(offset.dy - canvasHeight) * resolution / _displayLength;
    // shift back by the origin
    return IntCoord(
        (x - figure.origin.x).round(), (y - figure.origin.y).round());
  }

  List<double> buildXTicks({int logicalStep = 1}) {
    final firstLogical = -figure.origin.x.ceil() ~/ logicalStep * logicalStep;
    final out = <double>[];
    for (var i = 0; i <= figure.width; i += logicalStep) {
      final logical = IntCoord(firstLogical + i, 0);
      final offset = logicalIntToVisual(logical);
      out.add(offset.dx);
    }
    return out;
  }

  List<double> buildYTicks({int logicalStep = 1}) {
    final firstLogical = -figure.origin.y.ceil() ~/ logicalStep * logicalStep;
    final out = <double>[];
    for (var i = 0; i <= figure.height; i += logicalStep) {
      final logical = IntCoord(0, firstLogical + i);
      final offset = logicalIntToVisual(logical);
      out.add(offset.dy);
    }
    return out;
  }
}

/// [PointMovedNotification] is emitted when a point
/// is moved by drag and drop
class PointMovedNotification<T> extends Notification {
  final T id;
  final IntCoord logicalPos;
  PointMovedNotification(this.id, this.logicalPos);
}

class BaseRepere<PointIDType extends Object> extends StatelessWidget {
  final RepereMetrics metrics;

  final bool showGrid;

  /// [layers] are added in the stack
  final List<Widget> layers;

  const BaseRepere(this.metrics, this.showGrid, this.layers, {Key? key})
      : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Container(
      decoration: const BoxDecoration(
        // color: Colors.white.withOpacity(0.8),
        boxShadow: [BoxShadow(color: Colors.white, blurRadius: 5)],
      ),
      child: DragTarget<PointIDType>(
        builder: (context, candidateData, rejectedData) {
          final hasDropOver = candidateData.isNotEmpty;
          return Stack(
            clipBehavior: Clip.none,
            children: [
              _OriginPainter.asCustomPaint(metrics),
              // grid
              if (showGrid) ...[
                _GridPainter.asCustomPaint(metrics, hasDropOver),
                _AxisPainter.asCustomPaint(metrics),
              ],
              ...layers
            ],
          );
        },
        onAcceptWithDetails: (details) {
          final box = context.findRenderObject();
          if (box is! RenderBox) {
            return;
          }
          final localOffset = box.globalToLocal(details.offset);
          final logicalCoord = metrics.visualToLogical(localOffset);
          PointMovedNotification(details.data, logicalCoord).dispatch(context);
        },
      ),
    );
  }
}

class VectorPainter extends CustomPainter {
  final Offset from;
  final Offset to;
  final Color color;
  VectorPainter(this.from, this.to, {this.color = Colors.blue});

  static Path arrowPath(Offset from, Offset to) {
    final arrowHead = to;

    final dir = Offset(from.dx - to.dx, from.dy - to.dy);
    final scale = 8 / dir.distance;
    final shift = dir.scale(scale, scale);
    final arrowBasePoint = to.translate(shift.dx, shift.dy);

    final normal = Offset(dir.dy, -dir.dx).scale(scale / 3, scale / 3);
    final p1 = arrowBasePoint.translate(normal.dx, normal.dy);
    final p2 = arrowBasePoint.translate(-normal.dx, -normal.dy);

    final path = Path();
    path.moveTo(arrowHead.dx, arrowHead.dy);
    path.lineTo(p1.dx, p1.dy);
    path.lineTo(p2.dx, p2.dy);
    path.close();
    return path;
  }

  @override
  void paint(Canvas canvas, Size size) {
    canvas.drawLine(
        from,
        to,
        Paint()
          ..color = color
          ..strokeWidth = 0.8);

    final path = arrowPath(from, to);
    canvas.drawPath(
        path,
        Paint()
          ..style = PaintingStyle.fill
          ..color = color);
  }

  @override
  bool shouldRepaint(covariant VectorPainter oldDelegate) {
    return from != oldDelegate.from || to != oldDelegate.to;
  }

  @override
  bool? hitTest(Offset position) {
    return false;
  }
}

class _AxisPainter extends CustomPainter {
  final RepereMetrics metrics;
  _AxisPainter(this.metrics);

  static CustomPaint asCustomPaint(RepereMetrics metrics) {
    return CustomPaint(
      size: Size(metrics.canvasWidth, metrics.canvasHeight),
      painter: _AxisPainter(metrics),
    );
  }

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

class _GridPainter extends CustomPainter {
  final RepereMetrics metrics;
  final bool isHighlighted;

  _GridPainter(this.metrics, this.isHighlighted);

  static CustomPaint asCustomPaint(RepereMetrics metrics, bool isHighlighted) {
    return CustomPaint(
      size: Size(metrics.canvasWidth, metrics.canvasHeight),
      painter: _GridPainter(metrics, isHighlighted),
    );
  }

  @override
  void paint(Canvas canvas, Size size) {
    // minor grid
    final minorPaint = Paint()
      ..color =
          isHighlighted ? Colors.deepOrange : Colors.grey.withOpacity(0.7);
    for (var x in metrics.buildXTicks()) {
      canvas.drawLine(Offset(x, 0), Offset(x, size.height), minorPaint);
    }
    for (var y in metrics.buildYTicks()) {
      canvas.drawLine(Offset(0, y), Offset(size.width, y), minorPaint);
    }

    // ticks
    final tickLinePaint = Paint()
      ..strokeWidth = 1
      ..color =
          isHighlighted ? Colors.deepOrange : Colors.grey.withOpacity(0.7);
    final ticksPaint = Paint()..strokeWidth = 1;
    final visualOrigin = metrics.logicalIntToVisual(const IntCoord(0, 0));
    for (var x in metrics.buildXTicks(logicalStep: 5)) {
      canvas.drawLine(Offset(x, 0), Offset(x, size.height), tickLinePaint);

      canvas.drawLine(Offset(x, visualOrigin.dy - 5),
          Offset(x, visualOrigin.dy + 5), ticksPaint);
    }
    for (var y in metrics.buildYTicks(logicalStep: 5)) {
      canvas.drawLine(Offset(0, y), Offset(size.width, y), tickLinePaint);
      canvas.drawLine(Offset(visualOrigin.dx - 5, y),
          Offset(visualOrigin.dx + 5, y), ticksPaint);
    }
  }

  @override
  bool shouldRepaint(_GridPainter oldDelegate) {
    return isHighlighted != oldDelegate.isHighlighted;
  }
}

class _OriginPainter extends CustomPainter {
  final RepereMetrics metrics;

  _OriginPainter(this.metrics);

  static CustomPaint asCustomPaint(RepereMetrics metrics) {
    return CustomPaint(
      size: Size(metrics.canvasWidth, metrics.canvasHeight),
      painter: _OriginPainter(metrics),
    );
  }

  @override
  void paint(Canvas canvas, Size size) {
    // paint the origin if not implicit
    final origin = metrics.figure.origin;
    if (origin.x != 0 || origin.y != 0) {
      // note that logicalToVisual already shift by the origin
      final pos = origin.y == 0 ? LabelPos.top : LabelPos.bottomRight;
      DrawingsPainter.paintPoint(
          metrics, canvas, LabeledPoint(const Coord(0, 0), pos), "O",
          color: Colors.black);
    }
  }

  @override
  bool shouldRepaint(_OriginPainter oldDelegate) {
    return metrics != oldDelegate.metrics;
  }
}

extension _OffsetLabel on LabelPos {
  Offset offset(double textWidth, double textHeight) {
    const padding = 3.0;
    switch (this) {
      case LabelPos.top:
        return Offset(-textWidth / 2, -(textHeight + padding));
      case LabelPos.bottom:
        return Offset(-textWidth / 2, padding);
      case LabelPos.left:
        return Offset(-(textWidth + padding), -textHeight / 2);
      case LabelPos.right: // nothing to do
        return Offset(padding, -textHeight / 2);
      case LabelPos.topLeft:
        return Offset(-(textWidth + padding), -textHeight);
      case LabelPos.topRight:
        return Offset(padding, -(textHeight + padding));
      case LabelPos.bottomLeft:
        return Offset(-(textWidth + padding), padding);
      case LabelPos.bottomRight:
        return const Offset(padding, padding);
    }
  }
}

/// fromHex expected a #FFFFFF string
Color fromHex(String color, {Color onEmpty = Colors.purple}) {
  if (color.isEmpty) {
    return onEmpty;
  }
  color = "FF" + color.replaceAll("#", "");
  return Color(int.parse(color, radix: 16));
}

/// fontSize must be set
List<InlineSpan> parseSubscript(String text, TextStyle regularStyle) {
  final out = <InlineSpan>[];
  while (true) {
    final index = text.indexOf("_");
    if (index == -1) {
      if (text.isNotEmpty) {
        out.add(TextSpan(text: text, style: regularStyle));
      }
      break;
    }

    if (index > 0) {
      // add the regular chunk
      out.add(TextSpan(text: text.substring(0, index), style: regularStyle));
    }

    // skip the _ ...
    text = text.substring(index + 1);
    // then end the subscript at the next white space
    var end = text.indexOf(RegExp(r'\s'));
    if (end == -1) {
      end = text.length;
    }

    final subscript = text.substring(0, end);
    text = text.substring(end);
    if (subscript.isNotEmpty) {
      const textScaleFactor = 0.7; // superscript is usually smaller in size
      // add the subscript chunck
      out.add(TextSpan(
        text: subscript,
        style: regularStyle.copyWith(
          fontSize: regularStyle.fontSize! * textScaleFactor,
        ),
      ));
    }
  }

  return out;
}

class DrawingsPainter extends CustomPainter {
  final RepereMetrics metrics;
  final Drawings drawings;

  DrawingsPainter(this.metrics, this.drawings);

  @override
  bool? hitTest(Offset position) {
    return false;
  }

  static void paintPoint(
      RepereMetrics metrics, Canvas canvas, LabeledPoint point, String name,
      {Color color = Colors.blue}) {
    paintText(metrics, canvas, point, name, color: color);
    canvas.drawCircle(
        metrics.logicalToVisual(point.point),
        2,
        Paint()
          ..style = PaintingStyle.fill
          ..color = color);
  }

  void _paintSegment(Canvas canvas, Segment line) {
    final from = drawings.points[line.from]!.point;
    final to = drawings.points[line.to]!.point;
    final visualFrom = metrics.logicalToVisual(from);
    final visualTo = metrics.logicalToVisual(to);
    canvas.drawLine(visualFrom, visualTo, Paint());

    if (line.asVector) {
      final path = VectorPainter.arrowPath(visualFrom, visualTo);
      canvas.drawPath(path, Paint()..style = PaintingStyle.fill);
    }

    if (line.labelName.isNotEmpty) {
      paintText(
          metrics,
          canvas,
          LabeledPoint(
              Coord((from.x + to.x) / 2, (from.y + to.y) / 2), line.labelPos),
          line.labelName);
    }
  }

  static void paintText(
      RepereMetrics metrics, Canvas canvas, LabeledPoint point, String name,
      {Color? color}) {
    color = color ?? Colors.blue.shade800;
    const weight = FontWeight.bold;
    final spans = parseSubscript(
        name,
        TextStyle(
          fontSize: 16,
          fontWeight: weight,
          color: color,
          backgroundColor: Colors.white.withOpacity(0.5),
        ));

    final pt = TextPainter(
        text: TextSpan(children: spans), textDirection: TextDirection.ltr);
    pt.layout();

    final textWidth = pt.width;
    final textHeight = pt.height;

    final originalPos = metrics.logicalToVisual(point.point);

    final offset = point.pos.offset(textWidth, textHeight);

    pt.paint(canvas, originalPos.translate(offset.dx, offset.dy));
  }

  static void paintAffineLine(
      RepereMetrics metrics, Canvas canvas, Line line, Size size) {
    final origin = metrics.figure.origin;
    // start point
    final logicalStart = Coord(-origin.x, line.a * (-origin.x) + line.b);

    // end point
    final logicalEnd = Coord(metrics.figure.width - origin.x,
        line.a * (metrics.figure.width - origin.x) + line.b);

    final lineColor = fromHex(line.color);
    canvas.save();
    canvas.clipRect(Rect.fromLTWH(0, 0, size.width, size.height));
    canvas.drawLine(
        metrics.logicalToVisual(logicalStart),
        metrics.logicalToVisual(logicalEnd),
        Paint()
          ..color = lineColor
          ..strokeWidth = 2);

    // label position
    final logicalLabelX = (metrics.figure.width - origin.x) / 2;
    final logicalLabel = Coord(logicalLabelX, line.a * origin.x + line.b);
    paintText(
        metrics,
        canvas,
        LabeledPoint(
            logicalLabel, line.a > 0 ? LabelPos.bottomRight : LabelPos.topLeft),
        line.label,
        color: lineColor);

    canvas.restore();
  }

  @override
  void paint(Canvas canvas, Size size) {
    for (var segment in drawings.segments) {
      _paintSegment(canvas, segment);
    }

    for (var line in drawings.lines) {
      paintAffineLine(metrics, canvas, line, size);
    }

    drawings.points.forEach((key, value) {
      paintPoint(metrics, canvas, value, key);
    });
  }

  @override
  bool shouldRepaint(covariant CustomPainter oldDelegate) {
    return false;
  }
}

class GridPoint extends StatelessWidget {
  final IntCoord logical;
  final Offset visual;

  const GridPoint(this.logical, this.visual, {Key? key}) : super(key: key);

  static const radius = 10.0;

  @override
  Widget build(BuildContext context) {
    return Positioned(
      left: visual.dx - radius / 2,
      top: visual.dy - radius / 2,
      child: Container(
        width: radius,
        height: radius,
        decoration: BoxDecoration(
            color: Colors.orange.withOpacity(0.8), shape: BoxShape.circle),
      ),
    );
  }
}

class DraggableGridPoint<T extends Object> extends StatelessWidget {
  final IntCoord logical;
  final Offset visual;
  final T id;
  final Color color;

  /// used to adjust feedback when used in an InteractiveViewer
  final double zoomFactor;
  final bool disabled;

  const DraggableGridPoint(this.logical, this.visual, this.id, this.zoomFactor,
      {Key? key, this.color = Colors.orange, this.disabled = false})
      : super(key: key);

  static const outerRadius = 20.0;
  static const radius = 8.0;

  @override
  Widget build(BuildContext context) {
    return Positioned(
        left: visual.dx - outerRadius / 2,
        top: visual.dy - outerRadius / 2,
        child: Draggable<T>(
          maxSimultaneousDrags: disabled ? 0 : null,
          data: id,
          child: Container(
            color: Colors.transparent,
            width: outerRadius,
            height: outerRadius,
            child: Center(
              child: Container(
                width: radius,
                height: radius,
                decoration: BoxDecoration(
                    color: color.withOpacity(0.8), shape: BoxShape.circle),
              ),
            ),
          ),
          feedback: Transform.translate(
            offset: Offset(-radius * zoomFactor / 2, -radius * zoomFactor / 2),
            child: Container(
              width: radius * zoomFactor,
              height: radius * zoomFactor,
              decoration: BoxDecoration(
                  color: color.withOpacity(0.8), shape: BoxShape.circle),
            ),
          ),
        ));
  }
}

class GridPointHighlight extends StatelessWidget {
  final IntCoord logical;
  final Offset visual;

  const GridPointHighlight(this.logical, this.visual, {Key? key})
      : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Positioned(
      child: Container(
          decoration: BoxDecoration(
            borderRadius: BorderRadius.circular(8),
            color: Colors.grey,
          ),
          padding: const EdgeInsets.all(8),
          child: Text("( ${logical.x} ; ${logical.y} )")),
      left: visual.dx - 20,
      top: visual.dy - 50,
    );
  }
}

bool isInBounds(IntCoord point, RepereBounds figure) {
  final sx = point.x + figure.origin.x;
  final sy = point.y + figure.origin.y;
  return 0 <= sx && sx <= figure.width && 0 <= sy && sy <= figure.height;
}
