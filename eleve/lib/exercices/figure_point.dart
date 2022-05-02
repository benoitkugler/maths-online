import 'package:eleve/exercices/fields.dart';
import 'package:eleve/exercices/repere.dart';
import 'package:eleve/exercices/repere.gen.dart';
import 'package:eleve/exercices/types.gen.dart';
import 'package:flutter/widgets.dart';

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

class FigurePointField extends StatelessWidget {
  final Figure spec;
  final FigurePointController controller;

  const FigurePointField(this.spec, this.controller, {Key? key})
      : super(key: key);

  @override
  Widget build(BuildContext context) {
    return InteractiveViewer(
      child: _FigurePrivate(spec, controller),
      maxScale: 5,
    );
  }
}

class _FigurePrivate extends StatefulWidget {
  final Figure figure;
  final FigurePointController controller;

  const _FigurePrivate(this.figure, this.controller, {Key? key})
      : super(key: key);

  @override
  _FigurePrivateState createState() => _FigurePrivateState();
}

class _FigurePrivateState extends State<_FigurePrivate> {
  // var showTooltip = false;

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
    return GestureDetector(
      onTapUp: widget.controller.enabled
          ? (details) => _setCurrentPoint(details.localPosition, metrics)
          : null,
      child: BaseRepere(
        metrics,
        widget.figure.showGrid,
        [
          // custom drawing
          CustomPaint(
            size: Size(metrics.canvasWidth, metrics.canvasHeight),
            painter: DrawingsPainter(metrics, widget.figure.drawings),
          ),
          if (point != null)
            GridPoint(point, metrics.logicalIntToVisual(point)),
          // if (showTooltip && point != null)
          //   GridPointHighlight(point, metrics.logicalIntToVisual(point)),
        ],
      ),
    );
  }
}
