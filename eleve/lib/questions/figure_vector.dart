import 'package:eleve/questions/fields.dart';
import 'package:eleve/questions/repere.dart';
import 'package:eleve/questions/repere.gen.dart';
import 'package:eleve/questions/types.gen.dart';
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

  void setPoint(IntCoord point, VectorPointID id) {
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

  @override
  void setData(Answer answer) {
    final ans = answer as DoublePointAnswer;
    from = ans.from;
    to = ans.to;
    _hasData = true;
  }
}

enum VectorPointID { from, to }

class FigureVectorField extends StatefulWidget {
  final FigureVectorController controller;
  final TransformationController zoom;

  const FigureVectorField(this.controller, this.zoom, {Key? key})
      : super(key: key);

  @override
  State<FigureVectorField> createState() => _FigureVectorFieldState();
}

class _FigureVectorFieldState extends State<FigureVectorField> {
  @override
  void initState() {
    widget.zoom.addListener(onZoomUpdate);
    super.initState();
  }

  @override
  void dispose() {
    widget.zoom.removeListener(onZoomUpdate);
    super.dispose();
  }

  void onZoomUpdate() {
    setState(() {});
  }

  @override
  Widget build(BuildContext context) {
    final figure = widget.controller.data.figure;
    final metrics = RepereMetrics(figure.bounds, context);
    final from = widget.controller.from;
    final to = widget.controller.to;
    final zoomFactor = widget.zoom.value.getMaxScaleOnAxis();
    final color = widget.controller.hasError ? Colors.red : null;

    final List<PositionnedText> texts = [];
    final CustomPainter linePainter;
    if (widget.controller.data.asLine) {
      linePainter = _AffineLinePainter(
        metrics,
        from,
        to,
        widget.controller.data.lineLabel,
        color: color,
      );
      texts.add((linePainter as _AffineLinePainter).positionnedLabel);
    } else {
      linePainter = VectorPainter(
          metrics.logicalIntToVisual(from), metrics.logicalIntToVisual(to),
          color: color);
    }

    final figurePainter = DrawingsPainter(metrics, figure.drawings);
    texts.addAll(figurePainter.extractTexts());

    return NotificationListener<PointMovedNotification<VectorPointID>>(
      onNotification: (event) {
        setState(() {
          widget.controller.setPoint(event.logicalPos, event.id);
        });
        return true;
      },
      child: BaseRepere<VectorPointID>(
          metrics,
          figure.showGrid,
          figure.showOrigin,
          [
            // custom drawing
            CustomPaint(
              size: metrics.size,
              painter: figurePainter,
            ),
            CustomPaint(
              size: metrics.size,
              painter: linePainter,
            ),
            DraggableGridPoint(
              from,
              metrics.logicalIntToVisual(from),
              VectorPointID.from,
              zoomFactor,
              disabled: !widget.controller.isEnabled,
              color: color,
            ),
            DraggableGridPoint(
              to,
              metrics.logicalIntToVisual(to),
              VectorPointID.to,
              zoomFactor,
              disabled: !widget.controller.isEnabled,
              color: color,
            ),
          ],
          texts,
          color: color),
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
  final Color color;

  _AffineLinePainter(this.metrics, this.from, this.to, this.label,
      {Color? color})
      : color = color ?? Colors.teal;

  PositionnedText get positionnedLabel =>
      DrawingsPainter.lineLabel(metrics, _line, label);

  /// infer line from the points
  AffineLine get _line {
    final from = Coord(this.from.x.toDouble(), this.from.y.toDouble());
    final to = Coord(this.to.x.toDouble(), this.to.y.toDouble());
    return DrawingsPainter.inferLine(from, to, color);
  }

  @override
  void paint(Canvas canvas, Size size) {
    DrawingsPainter.paintAffineLine(metrics, canvas, _line);
  }

  @override
  bool shouldRepaint(covariant _AffineLinePainter oldDelegate) {
    return metrics != oldDelegate.metrics ||
        !from.equals(oldDelegate.from) ||
        !to.equals(oldDelegate.to);
  }
}
