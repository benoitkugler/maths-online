import 'dart:math';
import 'dart:ui';

import 'package:eleve/activities/trivialpoursuit/pawn.dart';
import 'package:eleve/types/src_trivial.dart';
import 'package:flutter/material.dart';

import 'categories.dart';

typedef OnTapTile = void Function(int);

/// Board is a squared widget with a fixed side length
class Board extends StatelessWidget {
  static const middle = _RP(50, 50);

  final double sideLength;
  final OnTapTile onTap;
  final Set<int> highlights;
  final int pawnTile;

  const Board(this.sideLength, this.onTap, this.highlights, this.pawnTile,
      {Key? key})
      : super(key: key);

  static const innerRingRadius = 34;
  static const outerRingRadius = innerRingRadius + 14;
  static const angularSection = 180 / 6;

  /// graphical description of the board
  static const shapes = [
    // center, start
    _ShapeDescriptor(Categorie.purple, _Circle(middle, _RL(9))),
    // two vertical tiles
    _ShapeDescriptor(Categorie.yellow,
        _RoundedTrapezoide(middle, _RL(9), 180 + 20, 180 - 40, _RP(59, 30))),
    _ShapeDescriptor(
        Categorie.blue,
        _RoundedTrapezoide(middle, _RL(innerRingRadius),
            270 - angularSection / 2, angularSection, _RP(40 + 19, 30))),
    // cross
    _ShapeDescriptor(
        Categorie.purple,
        _ArcSection(middle, _RL(innerRingRadius), _RL(outerRingRadius),
            270 - angularSection / 2, angularSection)),
    // 5 regular sections
    _ShapeDescriptor(
        Categorie.green,
        _ArcSection(middle, _RL(innerRingRadius), _RL(outerRingRadius),
            270 + angularSection / 2, angularSection)),
    _ShapeDescriptor(
        Categorie.blue,
        _ArcSection(middle, _RL(innerRingRadius), _RL(outerRingRadius),
            270 + angularSection / 2 + angularSection, angularSection)),
    _ShapeDescriptor(
        Categorie.orange,
        _ArcSection(middle, _RL(innerRingRadius), _RL(outerRingRadius),
            270 + angularSection / 2 + 2 * angularSection, angularSection)),
    _ShapeDescriptor(
        Categorie.yellow,
        _ArcSection(middle, _RL(innerRingRadius), _RL(outerRingRadius),
            270 + angularSection / 2 + 3 * angularSection, angularSection)),
    _ShapeDescriptor(
        Categorie.purple,
        _ArcSection(middle, _RL(innerRingRadius), _RL(outerRingRadius),
            270 + angularSection / 2 + 4 * angularSection, angularSection)),
    // cross
    _ShapeDescriptor(
        Categorie.green,
        _ArcSection(middle, _RL(innerRingRadius), _RL(outerRingRadius),
            270 + angularSection / 2 + 5 * angularSection, angularSection)),
    // 5 regular sections
    _ShapeDescriptor(
        Categorie.blue,
        _ArcSection(middle, _RL(innerRingRadius), _RL(outerRingRadius),
            90 + angularSection / 2 + 0 * angularSection, angularSection)),
    _ShapeDescriptor(
        Categorie.orange,
        _ArcSection(middle, _RL(innerRingRadius), _RL(outerRingRadius),
            90 + angularSection / 2 + 1 * angularSection, angularSection)),
    _ShapeDescriptor(
        Categorie.yellow,
        _ArcSection(middle, _RL(innerRingRadius), _RL(outerRingRadius),
            90 + angularSection / 2 + 2 * angularSection, angularSection)),
    _ShapeDescriptor(
        Categorie.blue,
        _ArcSection(middle, _RL(innerRingRadius), _RL(outerRingRadius),
            90 + angularSection / 2 + 3 * angularSection, angularSection)),
    _ShapeDescriptor(
        Categorie.orange,
        _ArcSection(middle, _RL(innerRingRadius), _RL(outerRingRadius),
            90 + angularSection / 2 + 4 * angularSection, angularSection)),
    // two last vertical tiles
    _ShapeDescriptor(
        Categorie.orange,
        _RoundedTrapezoide(middle, _RL(innerRingRadius),
            90 - angularSection / 2, angularSection, _RP(41, 70))),
    _ShapeDescriptor(Categorie.green,
        _RoundedTrapezoide(middle, _RL(9), 20, 180 - 40, _RP(41, 70))),
  ];

  @override
  Widget build(BuildContext context) {
    final size = Size(sideLength, sideLength);

    final pawnShape = shapes[pawnTile];
    // adjust the pawn center for better visual result
    final center = pawnShape.builder.visualCenter(size) ??
        pawnShape.builder.buildPath(size).getBounds().center;

    final configs = List<_TileConfig>.generate(shapes.length, (index) {
      final shape = shapes[index];
      final path = shape.builder.buildPath(size);
      return _TileConfig(size, path, () => onTap(index), shape.categorie.color);
    });

    final List<Widget> regular = [];
    final List<Widget> highligthed = []; // wrap with a glow effect
    for (var i = 0; i < configs.length; i++) {
      if (highlights.contains(i)) {
        highligthed.add(_HightightedTile(configs[i]));
      } else {
        regular.add(_RegularTile(configs[i]));
      }
    }

    return Container(
      decoration: BoxDecoration(shape: BoxShape.circle, boxShadow: [
        BoxShadow(
            color: Colors.white.withOpacity(0.4),
            blurRadius: 5,
            spreadRadius: 2)
      ]),
      child: SizedBox.square(
          dimension: sideLength,
          child: Stack(
            children: // place highligths over regular
                regular + highligthed + [PawnImage(center, sideLength * 0.05)],
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

class _HightightedTile extends StatelessWidget {
  final _TileConfig tile;

  const _HightightedTile(this.tile, {Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Stack(
      children: [
        CustomPaint(
          size: tile.size,
          painter: _TilePainter(tile, true),
        ),
        GestureDetector(
          onTap: tile.onTap,
          child: _AnimatedGlow(tile),
        ),
      ],
    );
  }
}

// add a glow animation to a tile
class _AnimatedGlow extends StatefulWidget {
  final _TileConfig tile;

  const _AnimatedGlow(this.tile, {Key? key}) : super(key: key);

  @override
  __AnimatedGlowState createState() => __AnimatedGlowState();
}

class __AnimatedGlowState extends State<_AnimatedGlow>
    with SingleTickerProviderStateMixin {
  late AnimationController controller;
  static const radiusFactor = 5;
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
    return AnimatedBuilder(
        animation: controller,
        builder: (_, __) {
          final radius = radiusFactor * controller.value + 4;
          return CustomPaint(
            size: widget.tile.size, // required to work around pointer issues
            painter: _TileGlow(widget.tile.color, widget.tile.path, radius),
          );
        });
  }
}

class _TileGlow extends CustomPainter {
  final Color insideColor;
  final Path _path;
  final double blurRadius;

  static const highlightWidth = 15.0;

  const _TileGlow(this.insideColor, this._path, this.blurRadius);

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
        ..strokeWidth = highlightWidth
        ..color = Colors.white
        ..imageFilter =
            ImageFilter.blur(sigmaX: blurRadius, sigmaY: blurRadius),
    );

    canvas.restore();
  }

  @override
  bool shouldRepaint(_TileGlow oldDelegate) {
    return oldDelegate.blurRadius != blurRadius;
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
      ..strokeWidth = isHighlighted ? 4 : 6
      ..color = desc.color;
  }

  Color get darkened {
    const amount = 0.1;

    final hsl = HSLColor.fromColor(desc.color);
    final hslDark = hsl.withLightness((hsl.lightness - amount).clamp(0.0, 1.0));

    return hslDark.toColor();
  }

  Paint _fillStyle() {
    return Paint()
      ..style = PaintingStyle.fill
      ..color = isHighlighted ? darkened : desc.color.withOpacity(0.6);
  }

  @override
  bool? hitTest(Offset position) {
    return desc.path.contains(position);
  }

  @override
  void paint(Canvas canvas, Size size) {
    final path = desc.path;

    if (isHighlighted) {
      canvas.drawPath(path, _strokeStyle());
      canvas.drawPath(path, _fillStyle());
    } else {
      canvas.save();
      canvas.clipPath(path); // no to stroke outside the path
      canvas.drawPath(path, _strokeStyle());
      canvas.drawPath(path, _fillStyle());
      canvas.restore();
    }
  }

  @override
  bool shouldRepaint(_TilePainter oldDelegate) {
    return oldDelegate.isHighlighted != isHighlighted;
  }
}

class _ShapeDescriptor {
  final Categorie categorie;
  final _PathBuilder builder;
  const _ShapeDescriptor(this.categorie, this.builder);
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
