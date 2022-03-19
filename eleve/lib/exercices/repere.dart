import 'package:eleve/exercices/repere.gen.dart';
import 'package:flutter/material.dart';

class Repere extends StatelessWidget {
  final Figure spec;
  const Repere(this.spec, {Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    final displayLength = MediaQuery.of(context).size.shortestSide * 0.7;
    return _Board(displayLength, spec);
  }
}

class _Board extends StatefulWidget {
  final double displayLength;

  static const gridResolution = 50.0;
  final Figure spec;

  const _Board(this.displayLength, this.spec, {Key? key}) : super(key: key);

  @override
  __BoardState createState() => __BoardState();

  double get canvasWidth => displayLength * spec.width / gridResolution;
  double get canvasHeight => displayLength * spec.height / gridResolution;

  Offset _logicalToVisual(Coord point) {
    return Offset(displayLength * point.x / gridResolution,
        canvasHeight - displayLength * point.y / gridResolution);
  }

  Coord _visualToLogical(Offset offset) {
    return Coord(offset.dx * gridResolution ~/ displayLength,
        -(offset.dy - canvasHeight) * gridResolution ~/ displayLength);
  }

  List<double> _buildXTicks() {
    final out = <double>[];
    for (var i = 0; i <= spec.width; i += 1) {
      final logical = Coord(i, 0);
      final offset = _logicalToVisual(logical);
      out.add(offset.dx);
    }
    return out;
  }

  List<double> _buildYTicks() {
    final out = <double>[];
    for (var i = 0; i <= spec.height; i += 1) {
      final logical = Coord(0, i);
      final offset = _logicalToVisual(logical);
      out.add(offset.dy);
    }
    return out;
  }
}

class __BoardState extends State<_Board> {
  Coord? logicalPoint;
  bool showTooltip = false;

  void _setCurrentPoint(Offset offset) {
    setState(() {
      logicalPoint = widget._visualToLogical(offset);
    });
  }

  @override
  Widget build(BuildContext context) {
    final visualPos =
        logicalPoint == null ? null : widget._logicalToVisual(logicalPoint!);
    return GestureDetector(
      onTapUp: (details) => _setCurrentPoint(details.localPosition),
      onPanUpdate: (details) => _setCurrentPoint(details.localPosition),
      child: Container(
        decoration: BoxDecoration(
            color: Colors.white.withOpacity(0.7),
            boxShadow: const [BoxShadow(color: Colors.white, blurRadius: 5)]),
        child: Stack(
          children: [
            // grid
            CustomPaint(
              size: Size(widget.canvasWidth, widget.canvasHeight),
              painter:
                  _GridPainter(widget._buildXTicks(), widget._buildYTicks()),
            ),
            // custom drawing
            CustomPaint(
              size: Size(widget.canvasWidth, widget.canvasHeight),
              painter: _ReperePainter(widget.spec, widget._logicalToVisual),
            ),
            // board
            if (logicalPoint != null)
              _GridPoint(logicalPoint!, visualPos!, () {
                setState(() {
                  showTooltip = true;
                });
              }, () {
                setState(() {
                  showTooltip = false;
                });
              }),
            if (showTooltip)
              Positioned(
                child: Container(
                    decoration: BoxDecoration(
                      borderRadius: BorderRadius.circular(10),
                      color: Colors.grey,
                    ),
                    padding: const EdgeInsets.all(8),
                    child: Text("(${logicalPoint!.x};${logicalPoint!.y})")),
                left: visualPos!.dx,
                top: visualPos.dy + 20,
              ),
          ],
        ),
      ),
    );
  }
}

class _GridPainter extends CustomPainter {
  final List<double> xs;
  final List<double> ys;

  _GridPainter(this.xs, this.ys);

  @override
  void paint(Canvas canvas, Size size) {
    final paint = Paint()..color = Colors.grey.withOpacity(0.7);
    for (var x in xs) {
      canvas.drawLine(Offset(x, 0), Offset(x, size.height), paint);
    }
    for (var y in ys) {
      canvas.drawLine(Offset(0, y), Offset(size.width, y), paint);
    }
  }

  @override
  bool shouldRepaint(covariant CustomPainter oldDelegate) {
    return false;
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
  final Figure spec;
  // final double displayLength;
  // final double gridResolution;
  final Offset Function(Coord) logicalToVisual;

  _ReperePainter(this.spec, this.logicalToVisual);

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

    final originalPos = logicalToVisual(point.point);

    final offset = point.pos.offset(textWidth, textHeight);

    pt.paint(canvas, originalPos.translate(offset.dx, offset.dy));
  }

  void _paintPoint(Canvas canvas, LabeledPoint point, String name) {
    _paintText(canvas, point, name, Colors.blue.shade800, FontWeight.bold);
  }

  void _paintPoints(Canvas canvas) {
    spec.points.forEach((key, value) {
      _paintPoint(canvas, value, key);
    });
  }

  void _paintLine(Canvas canvas, Line line) {
    final from = spec.points[line.from]!.point;
    final to = spec.points[line.to]!.point;
    canvas.drawLine(logicalToVisual(from), logicalToVisual(to), Paint());

    if (line.labelName.isNotEmpty) {
      _paintPoint(
          canvas,
          LabeledPoint(
              Coord((from.x + to.x) ~/ 2, (from.y + to.y) ~/ 2), line.labelPos),
          line.labelName);
    }
  }

  void _paintLines(Canvas canvas) {
    for (var element in spec.lines) {
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

class _GridPoint extends StatelessWidget {
  final Coord logical;
  final Offset pos;
  final Function onLongPress;
  final Function onLongPressEnd;
  const _GridPoint(
      this.logical, this.pos, this.onLongPress, this.onLongPressEnd,
      {Key? key})
      : super(key: key);

  static const radius = 10.0;

  @override
  Widget build(BuildContext context) {
    return Positioned(
      left: pos.dx - radius / 2,
      top: pos.dy - radius / 2,
      child: GestureDetector(
        onLongPress: () => onLongPress(),
        onLongPressEnd: (_) => onLongPressEnd(),
        child: Container(
            width: radius,
            height: radius,
            decoration: BoxDecoration(
                color: Colors.orange.withOpacity(0.8), shape: BoxShape.circle)),
      ),
    );
  }
}
