import 'package:eleve/exercices/fields.dart';
import 'package:eleve/exercices/repere.dart';
import 'package:eleve/exercices/repere.gen.dart';
import 'package:eleve/exercices/types.gen.dart';
import 'package:flutter/material.dart';

class FigureVectorController extends FieldController {
  final Figure figure;

  IntCoord from;
  IntCoord to;
  bool _hasData = false;

  FigureVectorController(this.figure, void Function() onChange)
      : from = IntCoord(figure.width ~/ 4 - 2, figure.height ~/ 4),
        to = IntCoord(figure.width ~/ 4 + 2, figure.height ~/ 4),
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
    final metrics = RepereMetrics(widget.controller.figure, context);
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
          [
            DraggableGridPoint(from, metrics.logicalIntToVisual(from),
                VectorPointID.from, zoomFactor),
            DraggableGridPoint(to, metrics.logicalIntToVisual(to),
                VectorPointID.to, zoomFactor),
            CustomPaint(
              size: Size(metrics.canvasWidth, metrics.canvasHeight),
              painter: VectorPainter(metrics.logicalIntToVisual(from),
                  metrics.logicalIntToVisual(to)),
            ),
          ],
        ),
      ),
    );
  }
}
