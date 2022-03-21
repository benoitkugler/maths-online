import 'package:eleve/exercices/fields.dart';
import 'package:eleve/exercices/repere.dart';
import 'package:eleve/exercices/repere.gen.dart';
import 'package:eleve/exercices/types.gen.dart';
import 'package:flutter/material.dart';

class FigureVectorController extends FieldController {
  IntCoord from = const IntCoord(2, 2);
  IntCoord to = const IntCoord(4, 4);
  bool _hasData = false;

  FigureVectorController(void Function() onChange) : super(onChange);

  void setPoint(IntCoord point, PointID id) {
    switch (id) {
      case PointID.from:
        from = point;
        break;
      case PointID.to:
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

class FigureVectorField extends StatefulWidget {
  final Figure spec;
  final FigureVectorController controller;

  const FigureVectorField(this.spec, this.controller, {Key? key})
      : super(key: key);

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
    final metrics = RepereMetrics(widget.spec, context);
    final from = widget.controller.from;
    final to = widget.controller.to;
    final zoomFactor = _zoomController.value.getMaxScaleOnAxis();
    return InteractiveViewer(
      transformationController: _zoomController,
      maxScale: 5,
      child: NotificationListener<PointMovedNotification>(
        onNotification: (event) {
          setState(() {
            widget.controller.setPoint(event.logicalPos, event.id);
          });
          return true;
        },
        child: BaseRepere(
          metrics,
          [
            DraggableGridPoint(from, metrics.logicalIntToVisual(from),
                PointID.from, zoomFactor),
            DraggableGridPoint(
                to, metrics.logicalIntToVisual(to), PointID.to, zoomFactor),
            CustomPaint(
              size: Size(metrics.canvasWidth, metrics.canvasHeight),
              painter: _VectorPainter(metrics.logicalIntToVisual(from),
                  metrics.logicalIntToVisual(to)),
            ),
          ],
        ),
      ),
    );
  }
}

class _VectorPainter extends CustomPainter {
  final Offset from;
  final Offset to;
  _VectorPainter(this.from, this.to);

  @override
  void paint(Canvas canvas, Size size) {
    canvas.drawLine(
        from,
        to,
        Paint()
          ..color = Colors.blue
          ..strokeWidth = 0.8);

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
    canvas.drawPath(
        path,
        Paint()
          ..style = PaintingStyle.fill
          ..color = Colors.blue);
  }

  @override
  bool shouldRepaint(covariant _VectorPainter oldDelegate) {
    return from != oldDelegate.from || to != oldDelegate.to;
  }

  @override
  bool? hitTest(Offset position) {
    return false;
  }
}
