import 'dart:math';

import 'package:flutter/material.dart';

class UnlockAnimation extends StatefulWidget {
  final Color color;
  const UnlockAnimation(this.color, {super.key});

  @override
  State<UnlockAnimation> createState() => _UnlockAnimationState();
}

class _UnlockAnimationState extends State<UnlockAnimation>
    with SingleTickerProviderStateMixin {
  static const duration = Duration(milliseconds: 1200);

  late AnimationController animation;
  List<_Particle> particles = [];

  @override
  void initState() {
    super.initState();

    particles = List.generate(40, (_) => _Particle.random());

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
            _UnlockSnapshot(widget.color, particles, animation));
  }
}

class _UnlockSnapshot extends StatelessWidget {
  final Color color;
  final List<_Particle> particles;

  final AnimationController controller;
  final Animation<double> riseAnimation;
  final Animation<double> particleAnimation;
  final Animation<double> opacityAnimation;

  _UnlockSnapshot(this.color, this.particles, this.controller, {super.key})
      : riseAnimation = Tween<double>(begin: 150, end: 0).animate(
          CurvedAnimation(
              parent: controller,
              curve: const Interval(0, 0.3, curve: Curves.easeInOut)),
        ),
        particleAnimation = Tween<double>(begin: 0.2, end: 1).animate(
          CurvedAnimation(
              parent: controller,
              curve: const Interval(0.15, 1, curve: Curves.easeOut)),
        ),
        opacityAnimation = Tween<double>(begin: 0, end: 1).animate(
          CurvedAnimation(
              parent: controller,
              curve: const Interval(0.15, 1, curve: Curves.easeInOut)),
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
            width: 300,
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
          child: Text("Nouvelle ceinture !",
              style: Theme.of(context).textTheme.headlineMedium),
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
