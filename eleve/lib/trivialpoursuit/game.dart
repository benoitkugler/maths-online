import 'package:eleve/trivialpoursuit/board.dart';
import 'package:flutter/material.dart';

class Game extends StatefulWidget {
  const Game({Key? key}) : super(key: key);

  @override
  _GameState createState() => _GameState();
}

class _GameState extends State<Game> {
  Set<int> highligthedTiles = {};

  void onTapTile(int tile) {
    setState(() {
      if (highligthedTiles.contains(tile)) {
        highligthedTiles.remove(tile);
      } else {
        highligthedTiles.add(tile);
      }
    });
  }

  @override
  Widget build(BuildContext context) {
    return Container(
      color: Colors.grey,
      child: Center(
        child: AspectRatio(
          aspectRatio: 1,
          child: Board(onTapTile, highligthedTiles),
        ),
      ),
    );
  }
}
