import 'dart:math';

import 'package:eleve/exercices/repere.gen.dart';
import 'package:eleve/exercices/types.gen.dart';
import 'package:flutter/material.dart';

class StaticRepere extends StatelessWidget {
  final Figure spec;
  const StaticRepere(this.spec, {Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    final metrics = RepereMetrics(spec, context);
    return BaseRepere(metrics, const []);
  }
}

class RepereMetrics {
  final double
      _displayLength; // displayLength is the length of the largest size of the figure
  final Figure figure;

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

  List<double> buildXTicks() {
    final out = <double>[];
    for (var i = 0; i <= figure.width; i += 1) {
      final logical = Coord(i.toDouble() - figure.origin.x, 0);
      final offset = logicalToVisual(logical);
      out.add(offset.dx);
    }
    return out;
  }

  List<double> buildYTicks() {
    final out = <double>[];
    for (var i = 0; i <= figure.height; i += 1) {
      final logical = Coord(0, i.toDouble() - figure.origin.y);
      final offset = logicalToVisual(logical);
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

  /// [layers] are added in the stack
  final List<Widget> layers;

  const BaseRepere(this.metrics, this.layers, {Key? key}) : super(key: key);

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
              // grid
              if (metrics.figure.showGrid)
                CustomPaint(
                  size: Size(metrics.canvasWidth, metrics.canvasHeight),
                  painter: _GridPainter(metrics.buildXTicks(),
                      metrics.buildYTicks(), hasDropOver),
                ),
              // custom drawing
              CustomPaint(
                size: Size(metrics.canvasWidth, metrics.canvasHeight),
                painter: _ReperePainter(metrics),
              ),
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

class _GridPainter extends CustomPainter {
  final List<double> xs;
  final List<double> ys;
  final bool isHighlighted;

  _GridPainter(this.xs, this.ys, this.isHighlighted);

  @override
  void paint(Canvas canvas, Size size) {
    final paint = Paint()
      ..color =
          isHighlighted ? Colors.deepOrange : Colors.grey.withOpacity(0.7);
    for (var x in xs) {
      canvas.drawLine(Offset(x, 0), Offset(x, size.height), paint);
    }
    for (var y in ys) {
      canvas.drawLine(Offset(0, y), Offset(size.width, y), paint);
    }
  }

  @override
  bool shouldRepaint(_GridPainter oldDelegate) {
    return isHighlighted != oldDelegate.isHighlighted;
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

class _ReperePainter extends CustomPainter {
  final RepereMetrics metrics;

  _ReperePainter(this.metrics);

  @override
  bool? hitTest(_) {
    return false;
  }

  void _paintText(Canvas canvas, LabeledPoint point, String name,
      {Color? color}) {
    color = color ?? Colors.blue.shade800;
    const weight = FontWeight.bold;
    final pt = TextPainter(
        text: TextSpan(
          text: name,
          style: TextStyle(
            fontWeight: weight,
            color: color,
            backgroundColor: Colors.white.withOpacity(0.5),
          ),
        ),
        textDirection: TextDirection.ltr);
    pt.layout();

    final textWidth = pt.width;
    final textHeight = pt.height;

    final originalPos = metrics.logicalToVisual(point.point);

    final offset = point.pos.offset(textWidth, textHeight);

    pt.paint(canvas, originalPos.translate(offset.dx, offset.dy));
  }

  void _paintPoint(Canvas canvas, LabeledPoint point, String name,
      {Color color = Colors.blue}) {
    _paintText(canvas, point, name, color: color);
    canvas.drawCircle(
        metrics.logicalToVisual(point.point),
        2,
        Paint()
          ..style = PaintingStyle.fill
          ..color = color);
  }

  void _paintSegment(Canvas canvas, Segment line) {
    final from = metrics.figure.points[line.from]!.point;
    final to = metrics.figure.points[line.to]!.point;
    final visualFrom = metrics.logicalToVisual(from);
    final visualTo = metrics.logicalToVisual(to);
    canvas.drawLine(visualFrom, visualTo, Paint());

    if (line.asVector) {
      final path = VectorPainter.arrowPath(visualFrom, visualTo);
      canvas.drawPath(path, Paint()..style = PaintingStyle.fill);

      // _paintText(
      //     canvas,
      //     LabeledPoint(
      //         Coord((from.x + to.x) / 2, (from.y + to.y) / 2), line.labelPos),
      //     line.from + line.to);
    }

    if (line.labelName.isNotEmpty) {
      _paintText(
          canvas,
          LabeledPoint(
              Coord((from.x + to.x) / 2, (from.y + to.y) / 2), line.labelPos),
          line.labelName);
    }
  }

  void _paintLine(Canvas canvas, Line line) {
    final origin = metrics.figure.origin;
    // start point
    final logicalStart = Coord(-origin.x, line.a * (-origin.x) + line.b);

    // end point
    final logicalEnd = Coord(metrics.figure.width - origin.x,
        line.a * (metrics.figure.width - origin.x) + line.b);

    canvas.save();
    canvas.clipRect(
        Rect.fromLTWH(0, 0, metrics.canvasWidth, metrics.canvasHeight));
    canvas.drawLine(
        metrics.logicalToVisual(logicalStart),
        metrics.logicalToVisual(logicalEnd),
        Paint()
          ..color = Colors.purple
          ..strokeWidth = 2);

    // label position
    final logicalLabelX = (metrics.figure.width - origin.x) / 2;
    final logicalLabel = Coord(logicalLabelX, line.a * origin.x + line.b);
    _paintText(
        canvas,
        LabeledPoint(
            logicalLabel, line.a > 0 ? LabelPos.bottomRight : LabelPos.topLeft),
        line.label,
        color: Colors.purple);

    canvas.restore();
  }

  @override
  void paint(Canvas canvas, Size size) {
    for (var segment in metrics.figure.segments) {
      _paintSegment(canvas, segment);
    }

    for (var line in metrics.figure.lines) {
      _paintLine(canvas, line);
    }

    metrics.figure.points.forEach((key, value) {
      _paintPoint(canvas, value, key);
    });

    // paint the origin if not implicit
    final origin = metrics.figure.origin;
    if (origin.x != 0 || origin.y != 0) {
      // note that logicalToVisual already shift by the origin
      final pos = origin.y == 0 ? LabelPos.top : LabelPos.bottomRight;
      _paintPoint(canvas, LabeledPoint(const Coord(0, 0), pos), "O",
          color: Colors.black);
    }
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

  const DraggableGridPoint(this.logical, this.visual, this.id, this.zoomFactor,
      {Key? key, this.color = Colors.orange})
      : super(key: key);

  static const outerRadius = 15.0;
  static const radius = 8.0;

  @override
  Widget build(BuildContext context) {
    return Positioned(
        left: visual.dx - outerRadius / 2,
        top: visual.dy - outerRadius / 2,
        child: Draggable<T>(
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

bool isInBounds(IntCoord point, Figure figure) {
  return 0 <= point.x &&
      point.x <= figure.width &&
      0 <= point.y &&
      point.y <= figure.height;
}
