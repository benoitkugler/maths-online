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
///   - static and enabled
///   - disabled
class Dice extends StatelessWidget {
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

  double get faceSize => animation == null ? 60 : 80;

  @override
  Widget build(BuildContext context) {
    return SizedBox(
      // wrap the actual button into a box large enough to contain the animated version
      // so thaht the layout is not shaken
      height: 110,
      child: Center(
        child: RawMaterialButton(
          shape:
              RoundedRectangleBorder(borderRadius: BorderRadius.circular(10)),
          padding: const EdgeInsets.symmetric(vertical: 15),
          elevation: 2,
          onPressed: isDisabled ? null : onTap,
          child: AnimatedContainer(
            duration: const Duration(milliseconds: 400),
            decoration: BoxDecoration(
              color: isDisabled ? Colors.grey : Colors.white,
              borderRadius: BorderRadius.circular(10),
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
                  // this is actually never reached
                  return _DiceDots(faceSize: faceSize, face: Face.one);
                }),
          ),
        ),
      ),
    );
  }

  Color get shadow {
    if (isDisabled) {
      return Colors.blueGrey;
    }
    return animation != null ? Colors.red : Colors.blue;
  }
}
