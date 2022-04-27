import 'package:eleve/exercices/fields.dart';
import 'package:eleve/exercices/repere.dart';
import 'package:eleve/exercices/repere.gen.dart';
import 'package:eleve/exercices/types.gen.dart';
import 'package:flutter/material.dart';

enum VectorPairPointID { from1, to1, from2, to2 }

extension _Equal on IntCoord {
  bool equals(IntCoord other) {
    return x == other.x && y == other.y;
  }
}

class FigureVectorPairController extends FieldController {
  final Figure figure;

  IntCoord from1;
  IntCoord to1;
  IntCoord from2;
  IntCoord to2;
  bool _hasData = false;

  FigureVectorPairController(this.figure, void Function() onChange)
      : from1 =
            IntCoord(figure.bounds.width ~/ 4 - 2, figure.bounds.height ~/ 4),
        to1 = IntCoord(figure.bounds.width ~/ 4, figure.bounds.height ~/ 4),
        from2 =
            IntCoord(figure.bounds.width ~/ 4 + 3, figure.bounds.height ~/ 4),
        to2 = IntCoord(figure.bounds.width ~/ 4, figure.bounds.height ~/ 4 + 3),
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
}

class FigureVectorPairField extends StatefulWidget {
  final FigureVectorPairController controller;

  const FigureVectorPairField(this.controller, {Key? key}) : super(key: key);

  @override
  State<FigureVectorPairField> createState() => _FigureVectorPairFieldState();
}

class _FigureVectorPairFieldState extends State<FigureVectorPairField> {
  final _zoomController = TransformationController();

  @override
  void initState() {
    _zoomController.addListener(() => setState(() {}));
    super.initState();
  }

  @override
  Widget build(BuildContext context) {
    final figure = widget.controller.figure;
    final metrics = RepereMetrics(figure.bounds, context);
    final from1 = widget.controller.from1;
    final to1 = widget.controller.to1;
    final from2 = widget.controller.from2;
    final to2 = widget.controller.to2;
    final zoomFactor = _zoomController.value.getMaxScaleOnAxis();
    return InteractiveViewer(
      transformationController: _zoomController,
      maxScale: 5,
      child: NotificationListener<PointMovedNotification<VectorPairPointID>>(
        onNotification: (event) {
          setState(() {
            widget.controller.setPoint(event.logicalPos, event.id);
          });
          return true;
        },
        child: BaseRepere<VectorPairPointID>(
          metrics,
          figure.showGrid,
          [
            DraggableGridPoint(from1, metrics.logicalIntToVisual(from1),
                VectorPairPointID.from1, zoomFactor,
                color: Colors.orange, disabled: !widget.controller.enabled),
            DraggableGridPoint(to1, metrics.logicalIntToVisual(to1),
                VectorPairPointID.to1, zoomFactor,
                color: Colors.orange, disabled: !widget.controller.enabled),
            DraggableGridPoint(from2, metrics.logicalIntToVisual(from2),
                VectorPairPointID.from2, zoomFactor,
                color: Colors.deepOrange, disabled: !widget.controller.enabled),
            DraggableGridPoint(to2, metrics.logicalIntToVisual(to2),
                VectorPairPointID.to2, zoomFactor,
                color: Colors.deepOrange, disabled: !widget.controller.enabled),
            CustomPaint(
              size: Size(metrics.canvasWidth, metrics.canvasHeight),
              painter: VectorPainter(metrics.logicalIntToVisual(from1),
                  metrics.logicalIntToVisual(to1),
                  color: Colors.blue),
            ),
            CustomPaint(
              size: Size(metrics.canvasWidth, metrics.canvasHeight),
              painter: VectorPainter(metrics.logicalIntToVisual(from2),
                  metrics.logicalIntToVisual(to2),
                  color: Colors.purple),
            ),
          ],
        ),
      ),
    );
  }
}
