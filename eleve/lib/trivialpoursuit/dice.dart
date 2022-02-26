import 'dart:math';

import 'package:flutter/material.dart';

/// Face is the face of a dice.
enum Face { one, two, three }

class _Dice extends StatelessWidget {
  final Face face;
  final bool isRolling;
  const _Dice(this.face, this.isRolling, {Key? key}) : super(key: key);

  static const faceSize = 100.0;
  static const bulletSize = faceSize * 0.15;

  /// returns top left offset for positionned bullets
  static const faces = [
    [Offset(faceSize / 2 - bulletSize / 2, faceSize / 2 - bulletSize / 2)],
    [
      Offset(faceSize / 4 - bulletSize / 2, faceSize / 4 - bulletSize / 2),
      Offset(
          faceSize * 3 / 4 - bulletSize / 2, faceSize * 3 / 4 - bulletSize / 2)
    ],
    [
      Offset(faceSize / 2 - bulletSize / 2,
          faceSize / 2 - bulletSize / 2), // center
      Offset(faceSize / 4 - bulletSize / 2, faceSize / 4 - bulletSize / 2),
      Offset(
          faceSize * 3 / 4 - bulletSize / 2, faceSize * 3 / 4 - bulletSize / 2)
    ]
  ];

  @override
  Widget build(BuildContext context) {
    return Container(
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: const BorderRadiusDirectional.all(Radius.circular(10)),
        boxShadow: [
          BoxShadow(
              color: isRolling ? Colors.red : Colors.blue, blurRadius: 10),
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

class DiceRoll extends StatefulWidget {
  final void Function() onDone;
  final Face _result;
  const DiceRoll(this._result, this.onDone, {Key? key}) : super(key: key);

  @override
  _DiceRollState createState() => _DiceRollState();
}

class _DiceRollState extends State<DiceRoll> {
  Face currentFace = Face.one;
  bool isRolling = true;

  void roll() async {
    const choices = Face.values;
    for (var i = 10; i < 50; i++) {
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
    await Future<void>.delayed(const Duration(seconds: 1));
    widget.onDone();
  }

  @override
  void initState() {
    super.initState();
    roll();
  }

  @override
  Widget build(BuildContext context) {
    return _Dice(currentFace, isRolling);
  }
}
