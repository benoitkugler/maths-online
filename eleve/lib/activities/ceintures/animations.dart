import 'dart:math';

import 'package:eleve/types/src_sql_ceintures.dart';
import 'package:flutter/material.dart';

class UnlockAnimation extends StatefulWidget {
  final Color color;
  final String title;
  const UnlockAnimation(this.color, this.title, {super.key});

  @override
  State<UnlockAnimation> createState() => _UnlockAnimationState();
}

class _UnlockAnimationState extends State<UnlockAnimation>
    with SingleTickerProviderStateMixin {
  static const duration = Duration(milliseconds: 1500);

  late AnimationController animation;
  List<_Particle> particles = [];

  @override
  void initState() {
    super.initState();

    particles = List.generate(30, (_) => _Particle.random());

    animation = AnimationController(
      vsync: this,
      duration: duration,
      reverseDuration: duration,
    );

    animation.forward();
  }

  @override
  void dispose() {
    animation.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return AnimatedBuilder(
        animation: animation,
        builder: (context, child) =>
            _UnlockSnapshot(widget.color, widget.title, particles, animation));
  }
}

class _UnlockSnapshot extends StatelessWidget {
  final Color color;
  final String title;
  final List<_Particle> particles;

  final AnimationController controller;
  final Animation<double> riseAnimation;
  final Animation<double> particleAnimation;
  final Animation<double> opacityAnimation;

  _UnlockSnapshot(this.color, this.title, this.particles, this.controller,
      {super.key})
      : riseAnimation = Tween<double>(begin: 150, end: 0).animate(
          CurvedAnimation(
              parent: controller,
              curve: const Interval(0, 0.6, curve: Curves.easeInOut)),
        ),
        particleAnimation = Tween<double>(begin: 0.5, end: 1).animate(
          CurvedAnimation(
              parent: controller,
              curve: const Interval(0, 0.5, curve: Curves.easeOut)),
        ),
        opacityAnimation = Tween<double>(begin: 0, end: 1).animate(
          CurvedAnimation(
              parent: controller,
              curve: const Interval(0.15, 0.5, curve: Curves.easeInOut)),
        );

  @override
  Widget build(BuildContext context) {
    return Column(
      mainAxisSize: MainAxisSize.min,
      children: [
        CustomPaint(
          painter: _Painter(
            color,
            particles,
            particleAnimation.value,
          ),
          child: SizedBox(
            width: 400,
            height: 300,
            child: Center(
                child: Padding(
              padding: EdgeInsets.only(top: riseAnimation.value),
              child: Container(
                padding: const EdgeInsets.all(12),
                width: 120,
                height: 120,
                decoration: BoxDecoration(
                    boxShadow: [
                      BoxShadow(blurRadius: 3, spreadRadius: 5, color: color)
                    ],
                    shape: BoxShape.circle,
                    color: color == Colors.white
                        ? Colors.grey.shade300
                        : Colors.white),
                child: Image.asset(
                  "assets/images/yellow-belt.png",
                  color: color,
                ),
              ),
            )),
          ),
        ),
        Opacity(
          opacity: opacityAnimation.value,
          child: Card(
              child: Padding(
            padding: const EdgeInsets.all(8.0),
            child: Text(title, style: Theme.of(context).textTheme.titleLarge),
          )),
        )
      ],
    );
  }
}

class _Particle {
  final double size;
  final double rotation;
  // in radians
  final double angularPosition;
  final double finalDistanceRatio; // in [0;1]
  _Particle(
      this.size, this.rotation, this.angularPosition, this.finalDistanceRatio);

  factory _Particle.random() => _Particle(
      1 + Random().nextDouble() * 4,
      Random().nextDouble() * pi,
      Random().nextDouble() * 2 * pi,
      0.4 + Random().nextDouble() * 0.6);
}

class _Painter extends CustomPainter {
  final Color color;
  final List<_Particle> particles;
  final double progression; // in [0; 1]

  _Painter(this.color, this.particles, this.progression);

  @override
  void paint(Canvas canvas, Size size) {
    if (progression <= 0.2) return;
    final middle = size.center(Offset.zero);
    final length = size.shortestSide / 2;
    final paint = Paint()
      ..style = PaintingStyle.fill
      ..color = color;
    final border = Paint()
      ..style = PaintingStyle.stroke
      ..color = Colors.white.withOpacity(0.6);
    for (var particle in particles) {
      final particulePosition = middle +
          Offset(cos(particle.angularPosition), sin(particle.angularPosition)) *
              length *
              progression *
              particle.finalDistanceRatio;
      canvas.save();
      canvas.translate(particulePosition.dx, particulePosition.dy);
      canvas.rotate(particle.rotation);
      canvas.drawRect(
          Rect.fromCircle(center: Offset.zero, radius: particle.size), paint);
      canvas.drawRect(
          Rect.fromCircle(center: Offset.zero, radius: particle.size), border);
      canvas.restore();
    }
  }

  @override
  bool shouldRepaint(covariant _Painter oldDelegate) {
    return oldDelegate.progression != progression;
  }
}

extension RC on Rank {
  Color get color {
    switch (this) {
      case Rank.startRank:
        return Colors.transparent;
      case Rank.blanche:
        return Colors.white;
      case Rank.jaune:
        return Colors.yellow;
      case Rank.orange:
        return Colors.orange;
      case Rank.verteI:
        return Colors.lime;
      case Rank.verteII:
        return Colors.green;
      case Rank.bleue:
        return Colors.blue;
      case Rank.violet:
        return Colors.purple;
      case Rank.rouge:
        return Colors.red;
      case Rank.marron:
        return Colors.brown;
      case Rank.noire:
        return Colors.black87;
    }
  }

  Rank? get next => this == Rank.noire ? null : Rank.values[index + 1];
}

class CeintureIcon extends StatelessWidget {
  final Rank rank;
  final double scale;
  final bool withBackground;
  const CeintureIcon(this.rank,
      {super.key, this.scale = 1, this.withBackground = false});

  @override
  Widget build(BuildContext context) {
    return Container(
        height: 48 * scale,
        width: 48 * scale,
        decoration: BoxDecoration(
            boxShadow: withBackground
                ? [BoxShadow(color: rank.color, blurRadius: 2, spreadRadius: 2)]
                : null,
            shape: BoxShape.circle,
            color: !withBackground
                ? Colors.transparent
                : rank == Rank.blanche
                    ? Colors.grey
                    : Colors.white),
        child: Padding(
          padding: const EdgeInsets.all(2.0),
          child: Image.asset(
            "assets/images/yellow-belt.png",
            color: rank.color,
          ),
        ));
  }
}
