import 'package:eleve/trivialpoursuit/dice.dart' as dice;
import 'package:eleve/trivialpoursuit/events.gen.dart';
import 'package:flutter/material.dart';

class Game extends StatefulWidget {
  const Game({Key? key}) : super(key: key);

  @override
  _GameState createState() => _GameState();
}

class _GameState extends State<Game> {
  Set<int> highligthedTiles = {};
  int pawnTile = 0;

  void onTapTile(int tile) {
    setState(() {
      if (highligthedTiles.contains(tile)) {
        highligthedTiles.remove(tile);
      } else {
        highligthedTiles.add(tile);
      }
      pawnTile = tile;
    });
  }

  // triggers and wait for a dice roll
  // with the given value
  void _rollDice(dice.Face face) async {
    await Navigator.push(context, MaterialPageRoute<void>(builder: (context) {
      return Container(
        color: Colors.grey,
        child: Center(
          child: dice.DiceRoll(face, () => Navigator.pop(context)),
        ),
      );
    }));
  }

  // process the given event
  void _processEvent(Event event) {
    if (event is PlayerTurn) {
    } else if (event is DiceThrow) {
      _rollDice(dice.Face.values[event.face]);
    } else if (event is Move) {
      setState(() {
        pawnTile = event.tile;
      });
    } else {
      print(event);
    }
  }

  @override
  Widget build(BuildContext context) {
    return Container(
      color: Colors.grey,
      // child: Center(
      //   child: AspectRatio(
      //     aspectRatio: 1,
      //     child: Board(onTapTile, highligthedTiles, pawnTile),
      //   ),
      // ),
      child: TextButton(
          onPressed: () => _rollDice(dice.Face.two), child: const Text("Dice")),
    );
  }
}
