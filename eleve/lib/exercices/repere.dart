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
  final Figure spec;

  RepereMetrics(this.spec, BuildContext context)
      : _displayLength = MediaQuery.of(context).size.shortestSide * 0.95;

  double get resolution => max(spec.width, spec.height).toDouble();

  double get canvasWidth => _displayLength * spec.width / resolution;
  double get canvasHeight => _displayLength * spec.height / resolution;

  Offset logicalToVisual(Coord point) {
    return Offset(_displayLength * point.x / resolution,
        canvasHeight - _displayLength * point.y / resolution);
  }

  Offset logicalIntToVisual(IntCoord point) {
    return logicalToVisual(Coord(point.x.toDouble(), point.y.toDouble()));
  }

  IntCoord visualToLogical(Offset offset) {
    return IntCoord((offset.dx * resolution / _displayLength).round(),
        (-(offset.dy - canvasHeight) * resolution / _displayLength).round());
  }

  List<double> buildXTicks() {
    final out = <double>[];
    for (var i = 0; i <= spec.width; i += 1) {
      final logical = Coord(i.toDouble(), 0);
      final offset = logicalToVisual(logical);
      out.add(offset.dx);
    }
    return out;
  }

  List<double> buildYTicks() {
    final out = <double>[];
    for (var i = 0; i <= spec.height; i += 1) {
      final logical = Coord(0, i.toDouble());
      final offset = logicalToVisual(logical);
      out.add(offset.dy);
    }
    return out;
  }
}

enum PointID { from, to }

/// [PointMovedNotification] is emitted when a point
/// is moved by drag and drop
class PointMovedNotification extends Notification {
  final PointID id;
  final IntCoord logicalPos;
  PointMovedNotification(this.id, this.logicalPos);
}

class BaseRepere extends StatelessWidget {
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
      child: DragTarget<PointID>(
        builder: (context, candidateData, rejectedData) {
          final hasDropOver = candidateData.isNotEmpty;
          return Stack(
            clipBehavior: Clip.none,
            children: [
              // grid
              CustomPaint(
                size: Size(metrics.canvasWidth, metrics.canvasHeight),
                painter: _GridPainter(
                    metrics.buildXTicks(), metrics.buildYTicks(), hasDropOver),
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

  void _paintText(Canvas canvas, LabeledPoint point, String name, Color color,
      FontWeight weight) {
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

  void _paintPoint(Canvas canvas, LabeledPoint point, String name) {
    _paintText(canvas, point, name, Colors.blue.shade800, FontWeight.bold);
    canvas.drawCircle(
        metrics.logicalToVisual(point.point),
        2,
        Paint()
          ..style = PaintingStyle.fill
          ..color = Colors.blue);
  }

  void _paintPoints(Canvas canvas) {
    metrics.spec.points.forEach((key, value) {
      _paintPoint(canvas, value, key);
    });
  }

  void _paintLine(Canvas canvas, Line line) {
    final from = metrics.spec.points[line.from]!.point;
    final to = metrics.spec.points[line.to]!.point;
    canvas.drawLine(
        metrics.logicalToVisual(from), metrics.logicalToVisual(to), Paint());

    if (line.labelName.isNotEmpty) {
      _paintPoint(
          canvas,
          LabeledPoint(
              Coord((from.x + to.x) / 2, (from.y + to.y) / 2), line.labelPos),
          line.labelName);
    }
  }

  void _paintLines(Canvas canvas) {
    for (var element in metrics.spec.lines) {
      _paintLine(canvas, element);
    }
  }

  @override
  void paint(Canvas canvas, Size size) {
    _paintLines(canvas);

    _paintPoints(canvas);
  }

  @override
  bool shouldRepaint(covariant CustomPainter oldDelegate) {
    return false;
  }
}

class GridPoint extends StatelessWidget {
  final IntCoord logical;
  final Offset visual;
  final bool isDraggable;

  // final Function onLongPress;
  // final Function onLongPressEnd;

  const GridPoint(this.logical, this.visual,
      {Key? key, this.isDraggable = false})
      : super(key: key);

  static const radius = 10.0;

  @override
  Widget build(BuildContext context) {
    final shape = Container(
      width: radius,
      height: radius,
      decoration: BoxDecoration(
          color: Colors.orange.withOpacity(0.8), shape: BoxShape.circle),
    );
    return Positioned(
      left: visual.dx - radius / 2,
      top: visual.dy - radius / 2,
      child: isDraggable ? Draggable(child: shape, feedback: shape) : shape,
    );
  }
}

class DraggableGridPoint extends StatelessWidget {
  final IntCoord logical;
  final Offset visual;
  final PointID id;

  /// used to adjust feedback when used in an InteractiveViewer
  final double zoomFactor;

  const DraggableGridPoint(this.logical, this.visual, this.id, this.zoomFactor,
      {Key? key})
      : super(key: key);

  static const outerRadius = 15.0;
  static const radius = 8.0;

  @override
  Widget build(BuildContext context) {
    return Positioned(
        left: visual.dx - outerRadius / 2,
        top: visual.dy - outerRadius / 2,
        child: Draggable<PointID>(
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
                    color: Colors.orange.withOpacity(0.8),
                    shape: BoxShape.circle),
              ),
            ),
          ),
          feedback: Transform.translate(
            offset: Offset(-radius * zoomFactor / 2, -radius * zoomFactor / 2),
            child: Container(
              width: radius * zoomFactor,
              height: radius * zoomFactor,
              decoration: BoxDecoration(
                  color: Colors.orange.withOpacity(0.8),
                  shape: BoxShape.circle),
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
