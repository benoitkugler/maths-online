import 'package:eleve/questions/fields.dart';
import 'package:eleve/questions/function_graph.dart';
import 'package:eleve/questions/repere.dart';
import 'package:eleve/types/src_maths_questions_client.dart';
import 'package:eleve/types/src_maths_repere.dart';
import 'package:flutter/material.dart';

class _PointController extends GeometricConstructionController {
  IntCoord? point;
  _PointController(void Function() onChange) : super(onChange);

  void setPoint(IntCoord point) {
    this.point = point;
    onChange();
  }

  @override
  bool hasValidData() {
    return point != null;
  }

  @override
  Answer getData() {
    return PointAnswer(point!);
  }

  @override
  void setData(Answer answer) {
    point = (answer as PointAnswer).point;
  }
}

enum VectorPointID { from, to }

class _VectorController extends GeometricConstructionController {
  IntCoord from;
  IntCoord to;
  bool _hasData = false;

  _VectorController(void Function() onChange, RepereBounds bounds)
      : from = IntCoord(bounds.width ~/ 4 - 2, bounds.height ~/ 4),
        to = IntCoord(bounds.width ~/ 4 + 2, bounds.height ~/ 4),
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

enum VectorPairPointID { from1, to1, from2, to2 }

extension on IntCoord {
  bool equals(IntCoord other) {
    return x == other.x && y == other.y;
  }
}

class _VectorPairController extends GeometricConstructionController {
  IntCoord from1;
  IntCoord to1;
  IntCoord from2;
  IntCoord to2;
  bool _hasData = false;

  _VectorPairController(void Function() onChange, RepereBounds bounds)
      : from1 = IntCoord(bounds.width ~/ 4 - 2, bounds.height ~/ 4),
        to1 = IntCoord(bounds.width ~/ 4, bounds.height ~/ 4),
        from2 = IntCoord(bounds.width ~/ 4 + 3, bounds.height ~/ 4),
        to2 = IntCoord(bounds.width ~/ 4, bounds.height ~/ 4 + 3),
        super(onChange);

  void setPoint(IntCoord point, VectorPairPointID id) {
    switch (id) {
      case VectorPairPointID.from1:
        from1 = point;
        break;
      case VectorPairPointID.to1:
        to1 = point;
        break;
      case VectorPairPointID.from2:
        from2 = point;
        break;
      case VectorPairPointID.to2:
        to2 = point;
        break;
    }
    _hasData = true;
    onChange();
  }

  @override
  bool hasValidData() {
    // disable equals points, which are too trivial
    if (from1.equals(from2) && to1.equals(to2)) {
      return false;
    }
    return _hasData;
  }

  @override
  Answer getData() {
    return DoublePointPairAnswer(from1, to1, from2, to2);
  }

  @override
  void setData(Answer answer) {
    final ans = (answer as DoublePointPairAnswer);
    from1 = ans.from1;
    to1 = ans.to1;
    from2 = ans.from2;
    to2 = ans.to2;
    _hasData = true;
  }
}

extension on GeometricConstructionFieldBlock {
  RepereBounds bounds() {
    final bg = background;
    if (bg is FigureBlock) {
      return bg.figure.bounds;
    } else if (bg is FunctionsGraphBlock) {
      return bg.bounds;
    } else {
      throw "unsupported background type";
    }
  }
}

abstract class GeometricConstructionController extends FieldController {
  GeometricConstructionController(void Function() onChange) : super(onChange);

  factory GeometricConstructionController.fromBlock(
      void Function() onChange, GeometricConstructionFieldBlock block) {
    if (block.field is GFPoint) {
      return _PointController(onChange);
    } else if (block.field is GFVector) {
      return _VectorController(onChange, block.bounds());
    } else if (block.field is GFVectorPair) {
      return _VectorPairController(onChange, block.bounds());
    } else {
      throw "unsupported field type";
    }
  }
}

extension on FigureOrGraph {
  bool showGrid() {
    final fg = this;
    return (fg is FigureBlock) ? (fg.figure.showGrid) : true;
  }

  bool showOrigin() {
    final fg = this;
    return (fg is FigureBlock) ? (fg.figure.showOrigin) : true;
  }

  CustomPainterText painter(RepereMetrics metrics) {
    final fg = this;
    if (fg is FigureBlock) {
      return DrawingsPainter(metrics, fg.figure.drawings);
    } else if (fg is FunctionsGraphBlock) {
      return BezierCurvesPainter(metrics,
          functions: fg.functions,
          sequences: fg.sequences,
          areas: fg.areas,
          points: fg.points);
    } else {
      throw "unsupported background type";
    }
  }
}

class GeometricConstructionFieldW extends StatefulWidget {
  final GeometricConstructionFieldBlock data;
  final GeometricConstructionController controller;
  final TransformationController zoom;

  const GeometricConstructionFieldW(this.data, this.controller, this.zoom,
      {Key? key})
      : super(key: key);

  @override
  _GeometricConstructionFieldWState createState() =>
      _GeometricConstructionFieldWState();
}

class _GeometricConstructionFieldWState
    extends State<GeometricConstructionFieldW> {
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

  void _setCurrentPoint(Offset visual, RepereMetrics metrics) {
    final logical = metrics.visualToLogical(visual);
    if (isInBounds(logical, metrics.figure)) {
      setState(() {
        (widget.controller as _PointController).setPoint(logical);
      });
    }
  }

  Widget _build(
      {required BuildContext context,
      required RepereMetrics metrics,
      required CustomPainterText background,
      required bool showGrid,
      required bool showOrigin}) {
    // custom drawing for background
    final texts = background.extractTexts();
    final backgroundPaint =
        CustomPaint(size: metrics.size, painter: background);
    final ct = widget.controller;
    final color = ct.hasError ? Colors.red : null;
    final zoomFactor = widget.zoom.value.getMaxScaleOnAxis();
    //
    if (ct is _PointController) {
      final point = ct.point;
      return GestureDetector(
        onTapUp: ct.isEnabled
            ? (details) => _setCurrentPoint(details.localPosition, metrics)
            : null,
        child: BaseRepere(
          metrics,
          showGrid,
          showOrigin,
          [
            backgroundPaint,
            if (point != null)
              GridPoint(point, metrics.logicalIntToVisual(point), color: color),
          ],
          texts,
          color: color,
        ),
      );
      //
    } else if (ct is _VectorController) {
      final data = (widget.data.field as GFVector);
      final from = ct.from;
      final to = ct.to;
      final CustomPainter linePainter;
      if (data.asLine) {
        linePainter = _AffineLinePainter(
          metrics,
          from,
          to,
          data.lineLabel,
          color: color,
        );
        texts.add((linePainter as _AffineLinePainter).positionnedLabel);
      } else {
        linePainter = VectorPainter(
            metrics.logicalIntToVisual(from), metrics.logicalIntToVisual(to),
            color: color);
      }
      return NotificationListener<PointMovedNotification<VectorPointID>>(
        onNotification: (event) {
          setState(() {
            ct.setPoint(event.logicalPos, event.id);
          });
          return true;
        },
        child: BaseRepere<VectorPointID>(
            metrics,
            showGrid,
            showOrigin,
            [
              // custom drawing
              backgroundPaint,
              CustomPaint(
                size: metrics.size,
                painter: linePainter,
              ),
              DraggableGridPoint(
                from,
                metrics.logicalIntToVisual(from),
                VectorPointID.from,
                zoomFactor,
                disabled: !ct.isEnabled,
                color: color,
              ),
              DraggableGridPoint(
                to,
                metrics.logicalIntToVisual(to),
                VectorPointID.to,
                zoomFactor,
                disabled: !ct.isEnabled,
                color: color,
              ),
            ],
            texts,
            color: color),
      );
    } else if (ct is _VectorPairController) {
      final from1 = ct.from1;
      final to1 = ct.to1;
      final from2 = ct.from2;
      final to2 = ct.to2;
      final hasError = ct.hasError;
      return NotificationListener<PointMovedNotification<VectorPairPointID>>(
        onNotification: (event) {
          setState(() {
            ct.setPoint(event.logicalPos, event.id);
          });
          return true;
        },
        child: BaseRepere<VectorPairPointID>(
          metrics,
          showGrid,
          showOrigin,
          [
            // static figure
            backgroundPaint,
            DraggableGridPoint(from1, metrics.logicalIntToVisual(from1),
                VectorPairPointID.from1, zoomFactor,
                color: hasError ? Colors.red : Colors.yellow,
                disabled: !ct.isEnabled),
            DraggableGridPoint(to1, metrics.logicalIntToVisual(to1),
                VectorPairPointID.to1, zoomFactor,
                color: hasError ? Colors.red : Colors.yellow,
                disabled: !ct.isEnabled),
            DraggableGridPoint(from2, metrics.logicalIntToVisual(from2),
                VectorPairPointID.from2, zoomFactor,
                color: hasError ? Colors.red : Colors.teal,
                disabled: !ct.isEnabled),
            DraggableGridPoint(to2, metrics.logicalIntToVisual(to2),
                VectorPairPointID.to2, zoomFactor,
                color: hasError ? Colors.red : Colors.teal,
                disabled: !ct.isEnabled),
            CustomPaint(
              size: metrics.size,
              painter: VectorPainter(metrics.logicalIntToVisual(from1),
                  metrics.logicalIntToVisual(to1),
                  color: hasError ? Colors.redAccent : Colors.yellowAccent),
            ),
            CustomPaint(
              size: metrics.size,
              painter: VectorPainter(metrics.logicalIntToVisual(from2),
                  metrics.logicalIntToVisual(to2),
                  color: hasError ? Colors.redAccent : Colors.tealAccent),
            ),
          ],
          texts,
          color: color,
        ),
      );
    } else {
      throw "unsupported field type";
    }
  }

  @override
  Widget build(BuildContext context) {
    final metrics = RepereMetrics(widget.data.bounds(), context);
    final painter = widget.data.background.painter(metrics);
    return _build(
        context: context,
        metrics: metrics,
        background: painter,
        showGrid: widget.data.background.showGrid(),
        showOrigin: widget.data.background.showOrigin());
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
