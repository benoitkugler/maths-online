import 'dart:math';

import 'package:flutter/material.dart';

/// Face is the face of a dice.
enum Face { one, two, three }

class _DiceFace extends StatelessWidget {
  final Face face;
  final bool isRolling;
  final bool isDisabled;
  final double faceSize;
  const _DiceFace(this.faceSize, this.face, this.isRolling, this.isDisabled,
      {Key? key})
      : super(key: key);

  double get bulletSize => faceSize * 0.15;

  /// returns top left offset for positionned bullets
  List<List<Offset>> get faces => [
        [Offset(faceSize / 2 - bulletSize / 2, faceSize / 2 - bulletSize / 2)],
        [
          Offset(faceSize / 4 - bulletSize / 2, faceSize / 4 - bulletSize / 2),
          Offset(faceSize * 3 / 4 - bulletSize / 2,
              faceSize * 3 / 4 - bulletSize / 2)
        ],
        [
          Offset(faceSize / 2 - bulletSize / 2,
              faceSize / 2 - bulletSize / 2), // center
          Offset(faceSize / 4 - bulletSize / 2, faceSize / 4 - bulletSize / 2),
          Offset(faceSize * 3 / 4 - bulletSize / 2,
              faceSize * 3 / 4 - bulletSize / 2)
        ]
      ];

  Color get shadow {
    if (isDisabled) {
      return Colors.blueGrey;
    }
    return isRolling ? Colors.red : Colors.blue;
  }

  @override
  Widget build(BuildContext context) {
    return Container(
      decoration: BoxDecoration(
        color: isDisabled ? Colors.grey : Colors.white,
        borderRadius: const BorderRadiusDirectional.all(Radius.circular(10)),
        boxShadow: [
          BoxShadow(color: shadow, blurRadius: 10),
        ],
      ),
      width: faceSize,
      height: faceSize,
      child: Stack(
        children: faces[face.index]
            .map((e) => Positioned(
                  left: e.dx,
                  width: bulletSize,
                  top: e.dy,
                  height: bulletSize,
                  child: Container(
                    decoration: const BoxDecoration(
                      shape: BoxShape.circle,
                      color: Colors.black,
                    ),
                  ),
                ))
            .toList(),
      ),
    );
  }
}

/// Dice presents a dice roll, with three states :
///   - in animation
///   - paused on the last result
///   - disabled
class Dice extends StatelessWidget {
  /// If non null, provides the faces to animate the roll.
  final Stream<Face>? animation;

  /// Show the face with a different highligth color.
  final Face lastResult;
  final bool isDisabled;

  const Dice(this.animation, this.lastResult, this.isDisabled, {Key? key})
      : super(key: key);

  /// rollDice returns a stream animating a rolling dice,
  /// to be used as input of DiceRoll
  static Stream<Face> rollDice() async* {
    const choices = Face.values;
    var currentFace = Face.one;
    for (var i = 10; i < 40; i++) {
      await Future<void>.delayed(Duration(
          milliseconds: 20 + exp(i.toDouble() * log(400) / 50).round()));
      currentFace = choices[(currentFace.index + 1) % choices.length];
      yield currentFace;
    }
  }

  @override
  Widget build(BuildContext context) {
    final animation = this.animation;
    if (animation != null) {
      return StreamBuilder<Face>(
          stream: animation,
          builder: (context, snapshot) {
            if (snapshot.hasData) {
              return _DiceFace(60, snapshot.data!, true, false);
            } else if (snapshot.connectionState == ConnectionState.done) {
              return const _DiceFace(60, Face.one, false, false);
            }
            return const _DiceFace(60, Face.one, false, false);
          });
    }
    return _DiceFace(60, lastResult, false, isDisabled);
  }
}
