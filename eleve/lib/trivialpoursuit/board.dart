import 'dart:math';
import 'dart:ui';

import 'package:flutter/material.dart';
import 'package:flutter/rendering.dart';

typedef OnTapTile = void Function(int);

Color darken(Color color, [double amount = .1]) {
  assert(amount >= 0 && amount <= 1);

  final hsl = HSLColor.fromColor(color);
  final hslDark = hsl.withLightness((hsl.lightness - amount).clamp(0.0, 1.0));

  return hslDark.toColor();
}

class Board extends StatelessWidget {
  static const middle = _RP(50, 50);

  final OnTapTile onTap;
  final Set<int> highlights;

  const Board(this.onTap, this.highlights, {Key? key}) : super(key: key);

  static const innerRingRadius = 35;
  static const outerRingRadius = innerRingRadius + 10;
  static const crossAngularSection = 40.0;
  static const angularSection = (180 - crossAngularSection) / 5;
  static const shapes = [
    // center, start
    _ShapeDescriptor(Colors.purple, _Circle(middle, _RL(8))),
    // three vertical tiles
    _ShapeDescriptor(Colors.green,
        _RoundedTrapezoide(middle, _RL(8), 180 + 20, 180 - 40, _RP(60, 35))),
    _ShapeDescriptor(
        Colors.blue, _Trapeze(_RP(40, 25), _RP(20, 0), _RP(0, 10))),
    _ShapeDescriptor(
        Colors.green,
        _RoundedTrapezoide(
            middle, _RL(innerRingRadius), 270 - 20, 20 * 2, _RP(40 + 20, 25))),
    // cross
    _ShapeDescriptor(
        Colors.blue,
        _ArcSection(middle, _RL(innerRingRadius), _RL(outerRingRadius),
            270 - crossAngularSection / 2, crossAngularSection)),
    // 5 regular sections
    _ShapeDescriptor(
        Colors.green,
        _ArcSection(middle, _RL(innerRingRadius), _RL(outerRingRadius), 290,
            angularSection)),
    _ShapeDescriptor(
        Colors.orange,
        _ArcSection(middle, _RL(innerRingRadius), _RL(outerRingRadius),
            290 + angularSection, angularSection)),
    _ShapeDescriptor(
        Colors.purple,
        _ArcSection(middle, _RL(innerRingRadius), _RL(outerRingRadius),
            290 + 2 * angularSection, angularSection)),
    _ShapeDescriptor(
        Color(0xFFFBC02D),
        _ArcSection(middle, _RL(innerRingRadius), _RL(outerRingRadius),
            290 + 3 * angularSection, angularSection)),
    _ShapeDescriptor(
        Colors.orange,
        _ArcSection(middle, _RL(innerRingRadius), _RL(outerRingRadius),
            290 + 4 * angularSection, angularSection)),
    // cross
    _ShapeDescriptor(
        Colors.yellow,
        _ArcSection(middle, _RL(innerRingRadius), _RL(outerRingRadius),
            290 + 5 * angularSection, crossAngularSection)),
    // 5 regular sections
    _ShapeDescriptor(
        Colors.orange,
        _ArcSection(middle, _RL(innerRingRadius), _RL(outerRingRadius),
            90 + crossAngularSection / 2 + 0 * angularSection, angularSection)),
    _ShapeDescriptor(
        Colors.yellow,
        _ArcSection(middle, _RL(innerRingRadius), _RL(outerRingRadius),
            90 + crossAngularSection / 2 + 1 * angularSection, angularSection)),
    _ShapeDescriptor(
        Colors.purple,
        _ArcSection(middle, _RL(innerRingRadius), _RL(outerRingRadius),
            90 + crossAngularSection / 2 + 2 * angularSection, angularSection)),
    _ShapeDescriptor(
        Colors.orange,
        _ArcSection(middle, _RL(innerRingRadius), _RL(outerRingRadius),
            90 + crossAngularSection / 2 + 3 * angularSection, angularSection)),
    _ShapeDescriptor(
        Colors.green,
        _ArcSection(middle, _RL(innerRingRadius), _RL(outerRingRadius),
            90 + crossAngularSection / 2 + 4 * angularSection, angularSection)),
    // three last vertical tiles
    _ShapeDescriptor(
        Colors.blue,
        _RoundedTrapezoide(
            middle, _RL(innerRingRadius), 90 - 20, 20 * 2, _RP(40, 75))),
    _ShapeDescriptor(
        Colors.yellow, _Trapeze(_RP(40, 75), _RP(20, 0), _RP(0, -10))),
    _ShapeDescriptor(Colors.blue,
        _RoundedTrapezoide(middle, _RL(8), 20, 180 - 40, _RP(40, 65))),
  ];

  @override
  Widget build(BuildContext context) {
    return Stack(
      children: List<Widget>.generate(shapes.length, (index) {
        final shape = shapes[index];
        final isHighlighted = highlights.contains(index);
        if (isHighlighted) {
          // wrap with a glow effect
          return _HightightedTile(shape, () => onTap(index));
        }
        return _RegularTile(shape, () => onTap(index));
      }),
    );
  }
}

class _RP {
  final double x; // in [-100, 100]
  final double y; // in [-100, 100]

  const _RP(this.x, this.y);

  Offset resolve(Size size) {
    return Offset(x / 100 * size.width, y / 100 * size.height);
  }
}

class _RL {
  final int l; // in[-100, 100]
  const _RL(this.l);

  double resolve(Size size) {
    return l.toDouble() / 100 * size.shortestSide;
  }
}

class _RegularTile extends StatelessWidget {
  final _ShapeDescriptor shape;
  final void Function() onTap;

  const _RegularTile(this.shape, this.onTap, {Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return SizedBox(
      // pass double.infinity to prevent shrinking of the painter area to 0.
      width: double.infinity,
      height: double.infinity,
      child: GestureDetector(
        onTap: onTap,
        child: CustomPaint(
          painter: _TilePainter(shape, 0),
        ),
      ),
    );
  }
}

/// add a glow animation to a tile
class _HightightedTile extends StatefulWidget {
  final _ShapeDescriptor shape;
  final void Function() onTap;

  const _HightightedTile(this.shape, this.onTap, {Key? key}) : super(key: key);

  @override
  __HightightedTileState createState() => __HightightedTileState();
}

class __HightightedTileState extends State<_HightightedTile>
    with SingleTickerProviderStateMixin {
  late AnimationController controller;
  static const radiusFactor = 20;
  static const duration = 1;

  @override
  void initState() {
    super.initState();

    controller = AnimationController(
      vsync: this,
      duration: const Duration(seconds: duration),
      reverseDuration: const Duration(seconds: duration),
    );

    controller.repeat(reverse: true);
  }

  @override
  void dispose() {
    controller.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return SizedBox(
      // pass double.infinity to prevent shrinking of the painter area to 0.
      width: double.infinity,
      height: double.infinity,
      child: AnimatedBuilder(
          animation: controller,
          builder: (_, __) {
            return GestureDetector(
              onTap: widget.onTap,
              child: CustomPaint(
                painter: _TilePainter(
                    widget.shape, radiusFactor * controller.value + 10),
              ),
            );
          }),
    );
  }
}

/// _TilePainter is one board tile
class _TilePainter extends CustomPainter {
  final _ShapeDescriptor desc;
  final double highlightRadius; // 0 for no highight

  bool get isHighlighted => highlightRadius != 0;

  Path _path = Path(); // used in hitTesting

  _TilePainter(this.desc, this.highlightRadius);

  Paint _strokeStyle() {
    return Paint()
      ..style = PaintingStyle.stroke
      ..strokeWidth = 4
      ..color = desc.color;
  }

  Paint _fillStyle() {
    return Paint()
      ..style = PaintingStyle.fill
      ..color = isHighlighted ? desc.color : desc.color.withOpacity(0.6);
  }

  @override
  bool? hitTest(Offset position) {
    return _path.contains(position);
  }

  @override
  void paint(Canvas canvas, Size size) {
    _path = desc.builder.buildPath(size);

    final bounds = _path.getBounds();

    canvas.saveLayer(bounds, Paint());

    canvas.drawPath(_path, _strokeStyle());
    canvas.drawPath(_path, _fillStyle());

    if (isHighlighted) {
      const color = Colors.white;

      final Rect shadowRect = bounds.inflate(highlightRadius);
      final shadowPaint = Paint()
        ..blendMode = BlendMode.srcATop
        ..colorFilter = const ColorFilter.mode(color, BlendMode.srcOut)
        ..imageFilter = ImageFilter.blur(
            sigmaX: Shadow.convertRadiusToSigma(highlightRadius),
            sigmaY: Shadow.convertRadiusToSigma(highlightRadius));
      canvas.saveLayer(shadowRect, shadowPaint);

      canvas.drawPath(_path, _fillStyle());

      canvas.restore();
    }

    canvas.restore();
  }

  @override
  bool shouldRepaint(_TilePainter oldDelegate) {
    return oldDelegate.highlightRadius != highlightRadius;
  }

  bool hitTesting(Offset position) {
    return _path.contains(position);
  }
}

class _ShapeDescriptor {
  final Color color;
  final _PathBuilder builder;
  const _ShapeDescriptor(this.color, this.builder);
}

abstract class _PathBuilder {
  const _PathBuilder();
  Path buildPath(Size size);
}

class _ArcSection extends _PathBuilder {
  final _RP center;
  final _RL radiusMin;
  final _RL radiusMax;
  final double startAngle; // stored in degrees
  final double sweepAngle; // stored in degrees

  const _ArcSection(this.center, this.radiusMin, this.radiusMax,
      this.startAngle, this.sweepAngle);

  @override
  Path buildPath(Size size) {
    final path = Path();
    final startAngle = this.startAngle * pi / 180;
    final sweepAngle = this.sweepAngle * pi / 180;
    final center = this.center.resolve(size);
    path.arcTo(Rect.fromCircle(center: center, radius: radiusMax.resolve(size)),
        startAngle, sweepAngle, true);
    path.arcTo(Rect.fromCircle(center: center, radius: radiusMin.resolve(size)),
        startAngle + sweepAngle, -sweepAngle, false);
    path.close();
    return path;
  }
}

// one "horizontal" arc, and three straigth lines
class _RoundedTrapezoide extends _PathBuilder {
  final _RP arcCenter;
  final _RL arcRadius;
  final double arcStartAngle; // stored in degrees
  final double arcSweepAngle; // stored in degrees
  final _RP point; // absolute

  const _RoundedTrapezoide(this.arcCenter, this.arcRadius, this.arcStartAngle,
      this.arcSweepAngle, this.point);

  @override
  Path buildPath(Size size) {
    final path = Path();
    final startAngle = arcStartAngle * pi / 180;
    final sweepAngle = arcSweepAngle * pi / 180;
    path.arcTo(
        Rect.fromCircle(
            center: arcCenter.resolve(size), radius: arcRadius.resolve(size)),
        startAngle,
        sweepAngle,
        true);
    final pt = point.resolve(size);
    // infer dxBottom
    final dxBottom = 2 * (size.width / 2 - pt.dx);
    path.lineTo(pt.dx, pt.dy);
    path.relativeLineTo(dxBottom, 0);
    path.close();

    return path;
  }
}

class _Trapeze extends _PathBuilder {
  final _RP topLeft;
  final _RP topLeftToTopRight;
  final _RP topToBottom;

  const _Trapeze(this.topLeft, this.topLeftToTopRight, this.topToBottom);

  @override
  Path buildPath(Size size) {
    // infer dxBottom by symmetry
    final topLeft = this.topLeft.resolve(size);
    final topLeftToTopRight = this.topLeftToTopRight.resolve(size);
    final topToBottom = this.topToBottom.resolve(size);
    // infer dxBottom by symmetry
    final dxBottom = -(topLeftToTopRight.dx + 2 * topToBottom.dx);

    final path = Path();
    path.moveTo(topLeft.dx, topLeft.dy);
    path.relativeLineTo(topLeftToTopRight.dx, topLeftToTopRight.dy);
    path.relativeLineTo(topToBottom.dx, topToBottom.dy);
    path.relativeLineTo(dxBottom, 0);
    path.close();

    return path;
  }
}

class _Circle extends _PathBuilder {
  final _RP center;
  final _RL radius;

  const _Circle(this.center, this.radius);

  @override
  Path buildPath(Size size) {
    final center = this.center.resolve(size);
    final radius = this.radius.resolve(size);
    final path = Path();
    path.addOval(Rect.fromCircle(center: center, radius: radius));
    return path;
  }
}
