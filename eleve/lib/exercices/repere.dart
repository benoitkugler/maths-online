import 'dart:math';

import 'package:eleve/exercices/fields.dart';
import 'package:eleve/exercices/repere.gen.dart';
import 'package:eleve/exercices/types.gen.dart';
import 'package:flutter/material.dart';

class StaticRepere extends StatelessWidget {
  final Figure figure;
  const StaticRepere(this.figure, {Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    final metrics = RepereMetrics(figure.bounds, context);
    final painter = DrawingsPainter(metrics, figure.drawings);
    final texts = painter.extractTexts();
    return BaseRepere(
      metrics,
      figure.showGrid,
      [
        // custom drawing
        CustomPaint(
          size: metrics.size,
          painter: painter,
        ),
      ],
      texts,
    );
  }
}

class Tick {
  final int logical;
  final double visual;
  const Tick(this.logical, this.visual);
}

class Ticks extends Iterable<Tick> {
  final List<Tick> ticks;
  const Ticks(this.ticks);

  Iterable<double> get visual => ticks.map((e) => e.visual);
  Iterable<int> get logical => ticks.map((e) => e.logical);

  @override
  Iterator<Tick> get iterator => ticks.iterator;
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

  /// [size] is the actual widget size used
  Size get size => Size(canvasWidth, canvasHeight);

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

  Ticks buildXTicks({int logicalStep = 1}) {
    final xOrigin = figure.origin.x.ceil();
    final firstLogical = -xOrigin ~/ logicalStep * logicalStep;
    final ticks = <Tick>[];
    for (var i = 0;
        firstLogical + i + xOrigin <= figure.width;
        i += logicalStep) {
      final logical = IntCoord(firstLogical + i, 0);
      final offset = logicalIntToVisual(logical);
      ticks.add(Tick(firstLogical + i, offset.dx));
    }
    return Ticks(ticks);
  }

  Ticks buildYTicks({int logicalStep = 1}) {
    final yOrigin = figure.origin.y.ceil();
    final firstLogical = -yOrigin ~/ logicalStep * logicalStep;
    final ticks = <Tick>[];

    for (var i = 0;
        firstLogical + i + yOrigin <= figure.height;
        i += logicalStep) {
      final logical = IntCoord(0, firstLogical + i);
      final offset = logicalIntToVisual(logical);
      ticks.add(Tick(firstLogical + i, offset.dy));
    }
    return Ticks(ticks);
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
  final List<PositionnedText> texts;

  const BaseRepere(this.metrics, this.showGrid, this.layers, this.texts,
      {Key? key})
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
              ...layers,
              ...texts.map((text) => _PositionnedTextW(metrics, text))
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

/// [PositionnedText] is an instruction to draw text on
/// a figure, with logical coordinates
/// In order to support latex, it must be extracted in its own widget.
class PositionnedText {
  final String text;

  /// [pos] is the logical position, with a relative offset hint.
  final PosPoint pos;

  final Color color;

  const PositionnedText(this.text, this.pos, {this.color = Colors.blue});
}

// to be used in a stack
class _PositionnedTextW extends StatelessWidget {
  final RepereMetrics metrics;
  final PositionnedText text;

  const _PositionnedTextW(this.metrics, this.text, {Key? key})
      : super(key: key);

  @override
  Widget build(BuildContext context) {
    // position without offset adjustement, which is handled by _PosSingleChildLayoutDelegate
    final originalPos = metrics.logicalToVisual(text.pos.point);
    return Positioned(
      left: originalPos.dx,
      top: originalPos.dy,
      child: CustomSingleChildLayout(
        delegate: _PosSingleChildLayoutDelegate(metrics, text.pos.pos),
        child: Container(
          child: textMath(text.text, TextStyle(color: text.color)),
          color: Colors.white.withOpacity(0.85),
          padding: const EdgeInsets.all(1),
        ),
      ),
    );
  }
}

class _PosSingleChildLayoutDelegate extends SingleChildLayoutDelegate {
  final RepereMetrics metrics;
  final LabelPos pos;

  _PosSingleChildLayoutDelegate(this.metrics, this.pos);

  @override
  bool shouldRelayout(covariant SingleChildLayoutDelegate oldDelegate) {
    return this != oldDelegate;
  }

  @override
  Size getSize(BoxConstraints constraints) {
    return metrics.size;
  }

  @override
  Offset getPositionForChild(Size size, Size childSize) {
    final textWidth = childSize.width;
    final textHeight = childSize.height;
    return pos.offset(textWidth, textHeight);
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
      size: metrics.size,
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
      size: metrics.size,
      painter: _GridPainter(metrics, isHighlighted),
    );
  }

  @override
  void paint(Canvas canvas, Size size) {
    // minor grid
    final minorPaint = Paint()
      ..color =
          isHighlighted ? Colors.deepOrange : Colors.grey.withOpacity(0.7);
    for (var x in metrics.buildXTicks().visual) {
      canvas.drawLine(Offset(x, 0), Offset(x, size.height), minorPaint);
    }
    for (var y in metrics.buildYTicks().visual) {
      canvas.drawLine(Offset(0, y), Offset(size.width, y), minorPaint);
    }

    // ticks
    final majorGridPaint = Paint()
      ..strokeWidth = 1.5
      ..color =
          isHighlighted ? Colors.deepOrange : Colors.grey.withOpacity(0.7);
    final shortTickPaint = Paint()..strokeWidth = 1;
    final visualOrigin = metrics.logicalIntToVisual(const IntCoord(0, 0));
    for (var xTick in metrics.buildXTicks(logicalStep: 5)) {
      final x = xTick.visual;
      canvas.drawLine(
          Offset(x, 0), Offset(x, size.height), majorGridPaint); // major grid
      canvas.drawLine(Offset(x, visualOrigin.dy - 5),
          Offset(x, visualOrigin.dy + 5), shortTickPaint); // tick line

      // paint tick legend, expect for 0,0
      if (xTick.logical == 0) {
        continue;
      }
      DrawingsPainter._paintText(
          metrics,
          canvas,
          PosPoint(Coord(xTick.logical.toDouble(), 0), LabelPos.bottom),
          "${xTick.logical}",
          color: Colors.black);
    }
    for (var yTick in metrics.buildYTicks(logicalStep: 5)) {
      final y = yTick.visual;
      canvas.drawLine(
          Offset(0, y), Offset(size.width, y), majorGridPaint); // major grid
      canvas.drawLine(Offset(visualOrigin.dx - 5, y),
          Offset(visualOrigin.dx + 5, y), shortTickPaint); // ticks

      if (yTick.logical == 0) {
        continue;
      }
      DrawingsPainter._paintText(
          metrics,
          canvas,
          PosPoint(Coord(0, yTick.logical.toDouble()), LabelPos.left),
          "${yTick.logical} ",
          color: Colors.black);
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
      size: metrics.size,
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
      const point = Coord(0, 0);
      DrawingsPainter.paintPoint(metrics, canvas, point, color: Colors.black);
      DrawingsPainter._paintText(metrics, canvas, PosPoint(point, pos), "O");
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
      case LabelPos.hide:
        return Offset.zero;
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

class DrawingsPainter extends CustomPainter {
  final RepereMetrics metrics;
  final Drawings drawings;

  DrawingsPainter(this.metrics, this.drawings);

  /// [extractTexts] returns all the text chunks used in
  /// various items of [drawings], so that they can be rendered
  /// as LaTeX
  /// The painter itself then ignore text information.
  List<PositionnedText> extractTexts() {
    final List<PositionnedText> out = [];

    drawings.points.forEach((key, value) {
      if (value.point.pos == LabelPos.hide) {
        return;
      }
      out.add(PositionnedText(key, value.point,
          color: fromHex(value.color, onEmpty: Colors.blue)));
    });

    for (final segment in drawings.segments) {
      final from = drawings.points[segment.from]!.point.point;
      final to = drawings.points[segment.to]!.point.point;
      if (segment.labelName.isNotEmpty && segment.labelPos != LabelPos.hide) {
        out.add(
          PositionnedText(
              segment.labelName,
              PosPoint(Coord((from.x + to.x) / 2, (from.y + to.y) / 2),
                  segment.labelPos),
              color: fromHex(segment.color)),
        );
      }
    }

    for (final line in drawings.lines) {
      out.add(lineLabel(metrics, line));
    }

    return out;
  }

  @override
  bool? hitTest(Offset position) {
    return false;
  }

  static PositionnedText lineLabel(RepereMetrics metrics, Line line) {
    final origin = metrics.figure.origin;

    // label position
    final Coord logical;
    final LabelPos pos;
    if (line.a.isInfinite) {
      logical = Coord(line.b, 1.5);
      pos = LabelPos.left;
    } else if (line.a == 0) {
      // special case for horizontal lines
      final logicalLabelX = metrics.figure.width * 2 / 3 - origin.x;
      logical = Coord(logicalLabelX, line.b);
      pos = line.b > 0 ? LabelPos.bottom : LabelPos.top;
    } else {
      double logicalLabelX = metrics.figure.width - origin.x;
      double y = line.a * logicalLabelX + line.b;
      // crop to the visible part
      y = min(y, metrics.figure.height - origin.y);
      y = max(y, -origin.y);
      // adujst the x back
      logicalLabelX = (y - line.b) / line.a * 0.9;
      y = line.a * logicalLabelX + line.b;

      logical = Coord(logicalLabelX, y);
      pos = line.a > 0 ? LabelPos.bottomRight : LabelPos.topRight;
    }

    return PositionnedText(
      line.label,
      PosPoint(logical, pos),
      color: fromHex(line.color),
    );
  }

  /// infer line from the points
  static Line inferLine(Coord from, Coord to, String label, String color) {
    double a, b;
    if (to.x == from.x) {
      a = double.infinity;
      b = to.x;
    } else {
      a = (to.y - from.y) / (to.x - from.x);
      b = from.y - a * from.x;
    }
    return Line(
      label,
      color,
      a,
      b,
    );
  }

  static void paintPoint(RepereMetrics metrics, Canvas canvas, Coord point,
      {Color color = Colors.blue}) {
    canvas.drawCircle(
        metrics.logicalToVisual(point),
        2,
        Paint()
          ..style = PaintingStyle.fill
          ..color = color);
  }

  void _paintSegment(Canvas canvas, Segment line) {
    final from = drawings.points[line.from]!.point.point;
    final to = drawings.points[line.to]!.point.point;
    final visualFrom = metrics.logicalToVisual(from);
    final visualTo = metrics.logicalToVisual(to);
    final color = fromHex(line.color);
    switch (line.kind) {
      case SegmentKind.sKLine: // use affine painter
        paintAffineLine(
            metrics, canvas, inferLine(from, to, line.labelName, line.color),
            width: 1);
        break;
      case SegmentKind.sKSegment:
        canvas.drawLine(visualFrom, visualTo, Paint()..color = color);
        break;
      case SegmentKind.sKVector:
        canvas.drawLine(visualFrom, visualTo, Paint()..color = color);
        // add arrow
        final path = VectorPainter.arrowPath(visualFrom, visualTo);
        canvas.drawPath(
            path,
            Paint()
              ..style = PaintingStyle.fill
              ..color = color);
        break;
      default:
    }
  }

  // helper method for regular text, not supporting LaTeX
  static void _paintText(
      RepereMetrics metrics, Canvas canvas, PosPoint point, String text,
      {Color? color}) {
    color = color ?? Colors.blue.shade800;
    const weight = FontWeight.bold;
    final style = TextStyle(
      fontSize: 14,
      fontWeight: weight,
      color: color,
      backgroundColor: Colors.white.withOpacity(0.5),
    );

    final pt = TextPainter(
        text: TextSpan(text: text, style: style),
        textDirection: TextDirection.ltr);
    pt.layout();

    final textWidth = pt.width;
    final textHeight = pt.height;
    final offset = point.pos.offset(textWidth, textHeight);

    final originalPos = metrics.logicalToVisual(point.point);

    pt.paint(canvas, originalPos.translate(offset.dx, offset.dy));
  }

  /// if line.a is infinite, then line.b is interpreted as the abscisse
  /// of a vertical line
  /// the line label is ignored; use [lineLabel] to get it
  static void paintAffineLine(RepereMetrics metrics, Canvas canvas, Line line,
      {double width = 2}) {
    final origin = metrics.figure.origin;

    Coord logicalStart, logicalEnd;
    if (line.a.isInfinite) {
      // start point
      logicalStart = Coord(line.b, -origin.y);

      // end point
      logicalEnd = Coord(line.b, metrics.figure.height - origin.y);
    } else {
      // start point
      logicalStart = Coord(-origin.x, line.a * (-origin.x) + line.b);

      // end point
      logicalEnd = Coord(metrics.figure.width - origin.x,
          line.a * (metrics.figure.width - origin.x) + line.b);
    }

    final lineColor = fromHex(line.color);
    canvas.save();
    canvas
        .clipRect(Rect.fromLTWH(0, 0, metrics.size.width, metrics.size.height));
    canvas.drawLine(
        metrics.logicalToVisual(logicalStart),
        metrics.logicalToVisual(logicalEnd),
        Paint()
          ..color = lineColor
          ..strokeWidth = width);

    canvas.restore();
  }

  @override
  void paint(Canvas canvas, Size size) {
    for (var segment in drawings.segments) {
      _paintSegment(canvas, segment);
    }

    for (var line in drawings.lines) {
      paintAffineLine(metrics, canvas, line);
    }

    drawings.points.forEach((key, value) {
      paintPoint(metrics, canvas, value.point.point,
          color: fromHex(value.color));
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
