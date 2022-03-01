import 'dart:math';

import 'package:flutter/material.dart';

/// Face is the face of a dice.
enum Face { one, two, three }

class _Dice extends StatelessWidget {
  final Face face;
  final bool isRolling;
  final bool isDisabled;
  final double faceSize;
  const _Dice(this.faceSize, this.face, this.isRolling, this.isDisabled,
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

class DisabledDice extends StatelessWidget {
  final Face face;

  const DisabledDice(this.face, {Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return _Dice(60, face, false, true);
  }
}

class DiceRoll extends StatefulWidget {
  final Face _result;
  const DiceRoll(this._result, {Key? key}) : super(key: key);

  @override
  _DiceRollState createState() => _DiceRollState();
}

class DoneRolling extends Notification {}

class _DiceRollState extends State<DiceRoll> {
  Face currentFace = Face.one;
  bool isRolling = true;

  void roll() async {
    const choices = Face.values;
    for (var i = 10; i < 40; i++) {
      await Future<void>.delayed(Duration(
          milliseconds: 20 + exp(i.toDouble() * log(400) / 50).round()));
      setState(() {
        currentFace = choices[(currentFace.index + 1) % choices.length];
      });
    }
    setState(() {
      currentFace = widget._result;
      isRolling = false;
    });

    // let the result be seen
    await Future<void>.delayed(const Duration(seconds: 1));
    DoneRolling().dispatch(context);
  }

  @override
  void initState() {
    super.initState();
    roll();
  }

  @override
  Widget build(BuildContext context) {
    return _Dice(60, currentFace, isRolling, false);
  }
}
