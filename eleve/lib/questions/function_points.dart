import 'package:eleve/questions/fields.dart';
import 'package:eleve/questions/function_graph.dart';
import 'package:eleve/questions/repere.dart';
import 'package:eleve/types/src_maths_functiongrapher.dart';
import 'package:eleve/types/src_maths_questions_client.dart';
import 'package:eleve/types/src_maths_repere.dart';
import 'package:flutter/material.dart';

class FunctionPointsController extends FieldController {
  final FunctionPointsFieldBlock data;
  final List<int?> fxs;

  FunctionPointsController(this.data, void Function() onChange)
      : fxs = List<int?>.generate(data.xs.length, (index) => null),
        super(onChange);

  @override
  bool hasValidData() {
    return fxs.every((element) => element != null);
  }

  @override
  Answer getData() {
    return FunctionPointsAnswer(fxs.map((e) => e!).toList());
  }

  @override
  void setData(Answer answer) {
    final ans = (answer as FunctionPointsAnswer).fxs;
    for (var i = 0; i < fxs.length; i++) {
      fxs[i] = ans[i];
    }
  }
}

class FunctionPoints extends StatefulWidget {
  final FunctionPointsController controller;
  final TransformationController zoom;

  const FunctionPoints(this.controller, this.zoom, {Key? key})
      : super(key: key);

  @override
  State<FunctionPoints> createState() => _FunctionPointsState();
}

typedef _PointID = int;

class _FunctionPointsState extends State<FunctionPoints> {
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

  static Coord _controlFromDerivatives(
      Coord p0, Coord p2, double dFrom, double dTo) {
    if (dFrom == dTo) {
      return Coord((p0.x + p2.x) / 2, (p0.y + p2.y) / 2);
    }
    final xIntersec = (p2.y - p0.y + dFrom * p0.x - dTo * p2.x) / (dFrom - dTo);
    final yIntersec = dFrom * (xIntersec - p0.x) + p0.y;
    return Coord(xIntersec, yIntersec);
  }

  List<Coord> points() {
    final ct = widget.controller;
    final out = <Coord>[];
    for (var index = 0; index < ct.fxs.length; index++) {
      final y = ct.fxs[index];
      if (y == null) {
        continue;
      }
      out.add(Coord(ct.data.xs[index].toDouble(), y.toDouble()));
    }
    return out;
  }

  List<BezierCurve> segments() {
    final ct = widget.controller;
    final derivatives = ct.data.dfxs;
    final List<BezierCurve> out = [];
    for (var index = 0; index < ct.fxs.length - 1; index++) {
      final from = ct.fxs[index];
      final to = ct.fxs[index + 1];
      if (from == null || to == null) {
        continue;
      }
      final p0 = Coord(ct.data.xs[index].toDouble(), from.toDouble());
      final p2 = Coord(ct.data.xs[index + 1].toDouble(), to.toDouble());
      final p1 = _controlFromDerivatives(
          p0, p2, derivatives[index], derivatives[index + 1]);
      out.add(BezierCurve(p0, p1, p2));
    }
    return out;
  }

  @override
  Widget build(BuildContext context) {
    final ct = widget.controller;
    final metrics = RepereMetrics(ct.data.bounds, context);
    final decoration = FunctionDecoration(ct.data.label, "");
    final painter = BezierCurvesPainter(
      metrics,
      functions:
          ct.data.isDiscrete ? [] : [FunctionGraph(decoration, segments())],
      sequences:
          ct.data.isDiscrete ? [SequenceGraph(decoration, points())] : [],
    );
    final texts = painter.extractTexts();
    final color = widget.controller.hasError ? Colors.red : Colors.purple;
    return NotificationListener<PointMovedNotification<_PointID>>(
      onNotification: (notification) {
        setState(() {
          ct.fxs[notification.id] = notification.logicalPos.y;
          ct.onChange();
        });
        return true;
      },
      child: BaseRepere<_PointID>(
        metrics,
        true,
        true,
        [
          CustomPaint(
            size: metrics.size,
            painter: painter,
          ),
          ...List<Widget>.generate(ct.fxs.length, (index) {
            final logical = IntCoord(ct.data.xs[index], ct.fxs[index] ?? 0);
            return DraggableGridPoint<_PointID>(
              logical,
              metrics.logicalIntToVisual(logical),
              index,
              widget.zoom.value.getMaxScaleOnAxis(),
              disabled: !widget.controller.isEnabled,
              color: color,
            );
          }),
        ],
        texts,
        color: color,
      ),
    );
  }
}
