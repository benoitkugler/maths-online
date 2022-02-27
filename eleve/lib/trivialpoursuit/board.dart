import 'dart:math';
import 'dart:ui';

import 'package:eleve/trivialpoursuit/pawn.dart';
import 'package:flutter/material.dart';

typedef OnTapTile = void Function(int);

class Board extends StatelessWidget {
  static const middle = _RP(50, 50);

  final OnTapTile onTap;
  final Set<int> highlights;
  final int pawnTile;

  const Board(this.onTap, this.highlights, this.pawnTile, {Key? key})
      : super(key: key);

  static const innerRingRadius = 35;
  static const outerRingRadius = innerRingRadius + 10;
  static const crossAngularSection = 40.0;
  static const angularSection = (180 - crossAngularSection) / 5;

  /// graphical description of the board
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
    final sideLength = MediaQuery.of(context).size.shortestSide;
    final size = Size(sideLength, sideLength);

    final pawnShape = shapes[pawnTile];
    final center = pawnShape.builder.visualCenter(size) ??
        pawnShape.builder.buildPath(size).getBounds().center;

    return Center(
      child: SizedBox.square(
          dimension: sideLength,
          child: Stack(
            children: List<Widget>.generate(shapes.length, (index) {
                  final shape = shapes[index];
                  final path = shape.builder.buildPath(size);

                  final isHighlighted = highlights.contains(index);
                  final tile =
                      _TileConfig(size, path, () => onTap(index), shape.color);
                  if (isHighlighted) {
                    // wrap with a glow effect
                    return _HightightedTile(tile);
                  }
                  return _RegularTile(tile);
                }) +
                [PawnImage(center, sideLength * 0.05)],
          )),
    );
  }
}

class _RP {
  final double x; // in [-100, 100]
  final double y; // in [-100, 100]

  const _RP(this.x, this.y);

  Offset resolve(Size size) {
    return Offset(x / 100 * size.shortestSide, y / 100 * size.height);
  }
}

class _RL {
  final int l; // in[-100, 100]
  const _RL(this.l);

  double resolve(Size size) {
    return l.toDouble() / 100 * size.shortestSide;
  }
}

class _TileConfig {
  final Size size; // used to handle gesture detector issues
  final Path path;
  final void Function() onTap;
  final Color color;

  const _TileConfig(this.size, this.path, this.onTap, this.color);
}

class _RegularTile extends StatelessWidget {
  final _TileConfig tile;

  const _RegularTile(this.tile, {Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return GestureDetector(
      onTap: tile.onTap,
      child: CustomPaint(
        size: tile.size,
        painter: _TilePainter(tile, false),
      ),
    );
  }
}

/// add a glow animation to a tile
class _HightightedTile extends StatefulWidget {
  final _TileConfig tile;

  const _HightightedTile(this.tile, {Key? key}) : super(key: key);

  @override
  __HightightedTileState createState() => __HightightedTileState();
}

class __HightightedTileState extends State<_HightightedTile>
    with SingleTickerProviderStateMixin {
  late AnimationController controller;
  static const radiusFactor = 8;
  static const duration = Duration(milliseconds: 800);

  @override
  void initState() {
    super.initState();

    controller = AnimationController(
      vsync: this,
      duration: duration,
      reverseDuration: duration,
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
    return Stack(
      children: [
        CustomPaint(
          size: widget.tile.size,
          painter: _TilePainter(widget.tile, true),
        ),
        GestureDetector(
          onTap: widget.tile.onTap,
          child: AnimatedBuilder(
              animation: controller,
              builder: (_, __) {
                final radius = radiusFactor * controller.value + 10;
                return CustomPaint(
                  size: widget
                      .tile.size, // required to work around pointer issues
                  painter: _TileGlow(widget.tile.path, radius),
                );
              }),
        ),
      ],
    );
  }
}

class _TileGlow extends CustomPainter {
  static final color = Colors.white.withOpacity(0.8);

  final Path _path;
  final double highlightRadius;
  _TileGlow(this._path, this.highlightRadius);

  @override
  bool? hitTest(Offset position) {
    return _path.contains(position);
  }

  @override
  void paint(Canvas canvas, Size size) {
    canvas.save();
    canvas.clipPath(_path);
    canvas.drawPath(
      _path,
      Paint()
        ..style = PaintingStyle.stroke
        ..strokeWidth = highlightRadius
        ..color = color
        ..imageFilter = ImageFilter.blur(sigmaX: 5, sigmaY: 5),
    );

    canvas.restore();
  }

  @override
  bool shouldRepaint(_TileGlow oldDelegate) {
    return oldDelegate.highlightRadius != highlightRadius;
  }
}

/// _TilePainter is one board tile
class _TilePainter extends CustomPainter {
  final _TileConfig desc;
  final bool isHighlighted;

  _TilePainter(this.desc, this.isHighlighted);

  Paint _strokeStyle() {
    return Paint()
      ..style = PaintingStyle.stroke
      ..strokeWidth = 6
      ..color = desc.color;
  }

  Paint _fillStyle() {
    return Paint()
      ..style = PaintingStyle.fill
      ..color = isHighlighted ? desc.color : desc.color.withOpacity(0.6);
  }

  @override
  bool? hitTest(Offset position) {
    return desc.path.contains(position);
  }

  @override
  void paint(Canvas canvas, Size size) {
    final _path = desc.path;

    canvas.save();
    canvas.clipPath(_path); // no to stroke outside the path
    canvas.drawPath(_path, _strokeStyle());
    canvas.drawPath(_path, _fillStyle());

    canvas.restore();
  }

  @override
  bool shouldRepaint(_TilePainter oldDelegate) {
    return oldDelegate.isHighlighted != isHighlighted;
  }
}

class _ShapeDescriptor {
  final Color color;
  final _PathBuilder builder;
  const _ShapeDescriptor(this.color, this.builder);
}

abstract class _PathBuilder {
  const _PathBuilder();

  /// returns the path scaled to the given size
  Path buildPath(Size size);

  /// returns the esthetic (scaled) center of the shape
  /// if Null, the center of the path is used
  Offset? visualCenter(Size size);
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

  @override
  Offset? visualCenter(Size size) {
    return null;
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

  double get startRadians => arcStartAngle * pi / 180;
  double get sweepRadians => arcSweepAngle * pi / 180;

  @override
  Path buildPath(Size size) {
    final path = Path();
    path.arcTo(
        Rect.fromCircle(
            center: arcCenter.resolve(size), radius: arcRadius.resolve(size)),
        startRadians,
        sweepRadians,
        true);
    final pt = point.resolve(size);
    // infer dxBottom
    final dxBottom = 2 * (size.width / 2 - pt.dx);
    path.lineTo(pt.dx, pt.dy);
    path.relativeLineTo(dxBottom, 0);
    path.close();

    return path;
  }

  @override
  Offset? visualCenter(Size size) {
    final path = Path();
    path.arcTo(
        Rect.fromCircle(
            center: arcCenter.resolve(size), radius: arcRadius.resolve(size)),
        startRadians,
        sweepRadians,
        true);
    path.close();
    final pointY = point.resolve(size).dy;
    final roundBounds = path.getBounds();
    final outX = roundBounds.center.dx;
    if (pointY > roundBounds.bottom) {
      return Offset(outX, (pointY + roundBounds.bottom) / 2 - 5);
    } else {
      return Offset(outX, (pointY + roundBounds.top) / 2 + 5);
    }
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

  @override
  Offset? visualCenter(Size size) {
    return null;
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

  @override
  Offset? visualCenter(Size size) {
    return null;
  }
}
