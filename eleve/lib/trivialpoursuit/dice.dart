import 'dart:math';

import 'package:flutter/material.dart';

/// Face is the face of a dice.
enum Face { one, two, three }

class _DiceDots extends StatelessWidget {
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

  const _DiceDots({
    Key? key,
    required this.faceSize,
    required this.face,
  }) : super(key: key);

  final double faceSize;
  final Face face;

  double get bulletSize => faceSize * 0.15;

  @override
  Widget build(BuildContext context) {
    return Stack(
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
    for (var i = 25; i < 45; i++) {
      currentFace = choices[(currentFace.index + 1) % choices.length];
      yield currentFace;
      await Future<void>.delayed(Duration(
          milliseconds: 20 + exp(i.toDouble() * log(400) / 50).round()));
    }
  }

  double get faceSize => animation == null ? 60 : 80;

  @override
  Widget build(BuildContext context) {
    return AnimatedContainer(
      duration: const Duration(milliseconds: 400),
      decoration: BoxDecoration(
        color: isDisabled ? Colors.grey : Colors.white,
        borderRadius: const BorderRadiusDirectional.all(Radius.circular(10)),
        boxShadow: [
          BoxShadow(color: shadow, blurRadius: 10),
        ],
      ),
      width: faceSize,
      height: faceSize,
      child: StreamBuilder<Face>(
          stream: animation,
          builder: (context, snapshot) {
            if (snapshot.hasData) {
              return _DiceDots(faceSize: faceSize, face: snapshot.data!);
            }
            return _DiceDots(faceSize: faceSize, face: Face.one);
          }),
    );
  }

  // @override
  // Widget build(BuildContext context) {
  //   final animation = this.animation;
  //   if (animation != null) {
  //     return StreamBuilder<Face>(
  //         stream: animation,
  //         builder: (context, snapshot) {
  //           if (snapshot.hasData) {
  //             return _DiceFace(80, snapshot.data!, true, false);
  //           } else if (snapshot.connectionState == ConnectionState.done) {
  //             return const _DiceFace(60, Face.one, false, false);
  //           }
  //           return const _DiceFace(60, Face.one, false, false);
  //         });
  //   }
  //   return _DiceFace(60, lastResult, false, isDisabled);
  // }

  Color get shadow {
    if (isDisabled) {
      return Colors.blueGrey;
    }
    return animation != null ? Colors.red : Colors.blue;
  }

  // @override
  // Widget build(BuildContext context) {
  //   return AnimatedContainer(
  //     duration: const Duration(seconds: 1),
  //     decoration: BoxDecoration(
  //       color: isDisabled ? Colors.grey : Colors.white,
  //       borderRadius: const BorderRadiusDirectional.all(Radius.circular(10)),
  //       boxShadow: [
  //         BoxShadow(color: shadow, blurRadius: 10),
  //       ],
  //     ),
  //     width: faceSize,
  //     height: faceSize,
  //     child: _DiceDots(faceSize: faceSize, face: face),
  //   );
  // }
}
