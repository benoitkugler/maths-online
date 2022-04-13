import 'package:eleve/exercices/fields.dart';
import 'package:eleve/exercices/repere.dart';
import 'package:eleve/exercices/repere.gen.dart';
import 'package:eleve/exercices/types.gen.dart';
import 'package:flutter/material.dart';

class FigureVectorController extends FieldController {
  final FigureVectorFieldBlock data;

  IntCoord from;
  IntCoord to;
  bool _hasData = false;

  FigureVectorController(this.data, void Function() onChange)
      : from = IntCoord(
            data.figure.bounds.width ~/ 4 - 2, data.figure.bounds.height ~/ 4),
        to = IntCoord(
            data.figure.bounds.width ~/ 4 + 2, data.figure.bounds.height ~/ 4),
        super(onChange);

  bool _checkSameX(IntCoord point, VectorPointID id) {
    switch (id) {
      case VectorPointID.from:
        return point.x != to.x;
      case VectorPointID.to:
        return point.x != from.x;
    }
  }

  void setPoint(IntCoord point, VectorPointID id) {
    if (data.asLine && !_checkSameX(point, id)) {
      return;
    }

    switch (id) {
      case VectorPointID.from:
        from = point;
        break;
      case VectorPointID.to:
        to = point;
        break;
    }
    _hasData = true;
    onChange();
  }

  @override
  bool hasValidData() {
    return _hasData;
  }

  @override
  Answer getData() {
    return DoublePointAnswer(from, to);
  }
}

enum VectorPointID { from, to }

class FigureVectorField extends StatefulWidget {
  final FigureVectorController controller;

  const FigureVectorField(this.controller, {Key? key}) : super(key: key);

  @override
  State<FigureVectorField> createState() => _FigureVectorFieldState();
}

class _FigureVectorFieldState extends State<FigureVectorField> {
  final _zoomController = TransformationController();

  @override
  void initState() {
    _zoomController.addListener(() => setState(() {}));
    super.initState();
  }

  @override
  Widget build(BuildContext context) {
    final figure = widget.controller.data.figure;
    final metrics = RepereMetrics(figure.bounds, context);
    final from = widget.controller.from;
    final to = widget.controller.to;
    final zoomFactor = _zoomController.value.getMaxScaleOnAxis();
    return InteractiveViewer(
      transformationController: _zoomController,
      maxScale: 5,
      child: NotificationListener<PointMovedNotification<VectorPointID>>(
        onNotification: (event) {
          setState(() {
            widget.controller.setPoint(event.logicalPos, event.id);
          });
          return true;
        },
        child: BaseRepere<VectorPointID>(
          metrics,
          figure.showGrid,
          [
            // custom drawing
            CustomPaint(
              size: Size(metrics.canvasWidth, metrics.canvasHeight),
              painter: DrawingsPainter(metrics, figure.drawings),
            ),
            CustomPaint(
              size: Size(metrics.canvasWidth, metrics.canvasHeight),
              painter: widget.controller.data.asLine
                  ? _AffineLinePainter(
                      metrics, from, to, widget.controller.data.lineLabel)
                  : VectorPainter(metrics.logicalIntToVisual(from),
                      metrics.logicalIntToVisual(to)),
            ),
            DraggableGridPoint(from, metrics.logicalIntToVisual(from),
                VectorPointID.from, zoomFactor),
            DraggableGridPoint(to, metrics.logicalIntToVisual(to),
                VectorPointID.to, zoomFactor),
          ],
        ),
      ),
    );
  }
}

extension _Equals on IntCoord {
  bool equals(IntCoord other) {
    return x == other.x && y == other.y;
  }
}

class _AffineLinePainter extends CustomPainter {
  final RepereMetrics metrics;
  final IntCoord from;
  final IntCoord to;
  final String label;

  _AffineLinePainter(this.metrics, this.from, this.to, this.label);

  @override
  void paint(Canvas canvas, Size size) {
    final a = (to.y - from.y).toDouble() / (to.x - from.x);
    final b = from.y - a * from.x;
    DrawingsPainter.paintAffineLine(
        metrics,
        canvas,
        Line(
          label,
          "#a832a2",
          a,
          b,
        ),
        size);
  }

  @override
  bool shouldRepaint(covariant _AffineLinePainter oldDelegate) {
    return metrics != oldDelegate.metrics ||
        !from.equals(oldDelegate.from) ||
        !to.equals(oldDelegate.to);
  }
}
