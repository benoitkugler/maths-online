import 'package:eleve/questions/fields.dart';
import 'package:eleve/questions/repere.dart';
import 'package:eleve/questions/repere.gen.dart';
import 'package:eleve/questions/types.gen.dart';
import 'package:flutter/material.dart';

class FigurePointController extends FieldController {
  IntCoord? point;
  FigurePointController(void Function() onChange) : super(onChange);

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

class FigurePointField extends StatefulWidget {
  final Figure figure;
  final FigurePointController controller;

  const FigurePointField(this.figure, this.controller, {Key? key})
      : super(key: key);

  @override
  _FigurePointFieldState createState() => _FigurePointFieldState();
}

class _FigurePointFieldState extends State<FigurePointField> {
  //  @override
  // void initState() {
  //   widget.zoom.addListener(onZoomUpdate);
  //   super.initState();
  // }

  // @override
  // void dispose() {
  //   widget.zoom.removeListener(onZoomUpdate);
  //   super.dispose();
  // }

  // void onZoomUpdate() {
  //   setState(() {});
  // }

  void _setCurrentPoint(Offset visual, RepereMetrics metrics) {
    final logical = metrics.visualToLogical(visual);
    if (isInBounds(logical, metrics.figure)) {
      setState(() {
        widget.controller.setPoint(logical);
      });
    }
  }

  @override
  Widget build(BuildContext context) {
    final metrics = RepereMetrics(widget.figure.bounds, context);
    final point = widget.controller.point;
    final painter = DrawingsPainter(metrics, widget.figure.drawings);
    final texts = painter.extractTexts();
    return GestureDetector(
      onTapUp: widget.controller.enabled
          ? (details) => _setCurrentPoint(details.localPosition, metrics)
          : null,
      child: BaseRepere(
        metrics,
        widget.figure.showGrid,
        widget.figure.showOrigin,
        [
          // custom drawing
          CustomPaint(
            size: metrics.size,
            painter: painter,
          ),
          if (point != null)
            GridPoint(point, metrics.logicalIntToVisual(point),
                color: widget.controller.hasError ? Colors.red : null),
          // if (showTooltip && point != null)
          //   GridPointHighlight(point, metrics.logicalIntToVisual(point)),
        ],
        texts,
        color: widget.controller.hasError ? Colors.red : null,
      ),
    );
  }
}
