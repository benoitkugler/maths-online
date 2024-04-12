import 'dart:math';

import 'package:flutter/material.dart';

/// Face is the face of a dice.
enum Face { one, two, three, four, five, six }

class _DiceDots extends StatelessWidget {
  /// returns top left offset for the CENTER of positionned bullets
  /// square size = faceSize
  /// dot size = bulletSize
  List<List<Offset>> get faces => [
        // one
        [Offset(faceSize / 2, faceSize / 2)],
        // two
        [
          Offset(faceSize / 4, faceSize / 4),
          Offset(faceSize * 3 / 4, faceSize * 3 / 4)
        ],
        // three
        [
          Offset(faceSize / 2, faceSize / 2), // center
          Offset(faceSize / 4, faceSize / 4),
          Offset(faceSize * 3 / 4, faceSize * 3 / 4)
        ],
        // four
        [
          Offset(faceSize / 4, faceSize / 4),
          Offset(faceSize / 4, 3 * faceSize / 4),
          Offset(3 * faceSize / 4, faceSize / 4),
          Offset(3 * faceSize / 4, 3 * faceSize / 4),
        ],
        // five
        [
          Offset(faceSize / 2, faceSize / 2), // center
          Offset(faceSize / 4, faceSize / 4),
          Offset(faceSize / 4, 3 * faceSize / 4),
          Offset(faceSize * 3 / 4, faceSize / 4),
          Offset(faceSize * 3 / 4, 3 * faceSize / 4),
        ],
        // six
        [
          Offset(faceSize / 3, faceSize / 4),
          Offset(faceSize / 3, faceSize * 2 / 4),
          Offset(faceSize / 3, faceSize * 3 / 4),
          Offset(faceSize * 2 / 3, faceSize / 4),
          Offset(faceSize * 2 / 3, faceSize * 2 / 4),
          Offset(faceSize * 2 / 3, faceSize * 3 / 4),
        ]
      ];

  const _DiceDots({
    Key? key,
    required this.faceSize,
    required this.face,
  }) : super(key: key);

  final double faceSize;
  final Face face;

  @override
  Widget build(BuildContext context) {
    final bulletSize = faceSize * 0.15;
    final halfBulletSize = bulletSize / 2;
    final dots = faces[face.index];
    return Stack(
      children: dots
          .map((e) => Positioned(
                left: e.dx - halfBulletSize,
                top: e.dy - halfBulletSize,
                width: bulletSize,
                height: bulletSize,
                child: Container(
                  decoration: const BoxDecoration(
                    shape: BoxShape.circle,
                    color: Colors.black,
                  ),
                ),
              ))
          .toList(),
    );
  }
}

/// Dice presents a dice roll, with three states :
///   - in animation
///   - static and enabled
///   - disabled
class Dice extends StatefulWidget {
  /// [onTap] is ignored if [isDisabled] is true
  final void Function() onTap;

  /// If non null, provides the faces to animate the roll.
  final Stream<Face>? animation;

  final bool isDisabled;

  const Dice(this.onTap, this.animation, this.isDisabled, {Key? key})
      : super(key: key);

  /// rollDice returns a stream animating a rolling dice,
  /// to be used as input of DiceRoll
  static Stream<Face> rollDice(Face lastFace) async* {
    const choices = Face.values;
    var currentFace = Face.one;
    for (var i = 25; i < 45; i++) {
      currentFace = choices[(currentFace.index + 1) % choices.length];
      yield currentFace;
      await Future<void>.delayed(Duration(
          milliseconds: 20 + exp(i.toDouble() * log(400) / 50).round()));
    }
    yield lastFace;
  }

  @override
  State<Dice> createState() => _DiceState();
}

class _DiceState extends State<Dice> with SingleTickerProviderStateMixin {
  double get faceSize => widget.animation == null ? 60 : 80;

  late AnimationController _animationController;
  late Animation<double> _animation;

  @override
  void initState() {
    _animationController = AnimationController(
        vsync: this, duration: const Duration(milliseconds: 500));
    _animation = Tween(begin: 1.0, end: 20.0).animate(
        CurvedAnimation(parent: _animationController, curve: Curves.easeIn))
      ..addListener(() {
        setState(() {});
      });
    if (!widget.isDisabled) {
      _animationController.repeat(reverse: true);
    }
    super.initState();
  }

  @override
  void dispose() {
    _animationController.dispose();
    super.dispose();
  }

  @override
  void didUpdateWidget(covariant Dice oldWidget) {
    if (!widget.isDisabled) {
      _animationController.repeat(reverse: true);
    } else {
      _animationController.reset();
    }
    super.didUpdateWidget(oldWidget);
  }

  @override
  Widget build(BuildContext context) {
    return SizedBox(
      // wrap the actual button into a box large enough to contain the animated version
      // so that the layout is not shaken
      height: 110,
      child: Center(
        child: RawMaterialButton(
          shape:
              RoundedRectangleBorder(borderRadius: BorderRadius.circular(10)),
          padding: const EdgeInsets.symmetric(vertical: 15),
          elevation: 2,
          onPressed: widget.isDisabled ? null : widget.onTap,
          child: AnimatedContainer(
            duration: const Duration(milliseconds: 400),
            decoration: BoxDecoration(
              color: widget.isDisabled ? Colors.grey : Colors.white,
              borderRadius: BorderRadius.circular(10),
              boxShadow: [
                BoxShadow(
                    color: shadow,
                    blurRadius: _animation.value,
                    spreadRadius: _animation.value),
              ],
            ),
            width: faceSize,
            height: faceSize,
            child: StreamBuilder<Face>(
                stream: widget.animation,
                builder: (context, snapshot) {
                  if (snapshot.hasData) {
                    return _DiceDots(faceSize: faceSize, face: snapshot.data!);
                  }
                  // this is actually never reached
                  return _DiceDots(faceSize: faceSize, face: Face.one);
                }),
          ),
        ),
      ),
    );
  }

  Color get shadow {
    if (widget.isDisabled) {
      return Colors.blueGrey;
    }
    return widget.animation != null
        ? Colors.red
        : const Color.fromARGB(255, 33, 243, 208);
  }
}
