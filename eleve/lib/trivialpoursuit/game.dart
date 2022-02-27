import 'package:eleve/trivialpoursuit/board.dart';
import 'package:eleve/trivialpoursuit/dice.dart' as dice;
import 'package:eleve/trivialpoursuit/events.gen.dart';
import 'package:eleve/trivialpoursuit/question.dart';
import 'package:flutter/material.dart';

class GameController extends StatefulWidget {
  final int playerID; // the index of the client player

  const GameController(this.playerID, {Key? key}) : super(key: key);

  @override
  _GameControllerState createState() => _GameControllerState();
}

class _GameControllerState extends State<GameController> {
  GameState state = const GameState([], 0, 0);

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

  void _onPlayerTurn(PlayerTurn event) {
    // TODO: animation
    // the actual state modification is handled in TODO
  }

  Future<void> _showRoute(Widget content) {
    return Navigator.push<void>(context, MaterialPageRoute(builder: (context) {
      return Container(
        color: Colors.grey,
        child: Center(
          child: content,
        ),
      );
    }));
  }

  // triggers and wait for a dice roll
  // with the given value
  Future<void> _onDiceThrow(DiceThrow event) async {
    final face =
        dice.Face.values[event.face - 1]; // event.face is the "natural" face
    final widget = dice.DiceRoll(face, () => Navigator.pop(context));
    return _showRoute(widget);
  }

  Future<void> _onPossibleMoves(PossibleMoves event) async {
    // only the current player may choose the tile to move
    if (widget.playerID != state.player) {
      return;
    }
    setState(() {
      highligthedTiles = event.tiles.toSet();
    });
    return Future.delayed(const Duration(seconds: 2));
  }

  void _onMove(Move event) {
    setState(() {
      highligthedTiles.clear();
    });
  }

  Future<void> _onShowQuestion(ShowQuestion event) async {
    final widget = QuestionRoute(event);
    return _showRoute(widget);
  }

  void _onPlayerAnswerResult(PlayerAnswerResult event) {
    // for now, we simply ignore other player success
    if (event.player != widget.playerID) {
      return;
    }

    ScaffoldMessenger.of(context).showSnackBar(SnackBar(
      duration: const Duration(seconds: 3),
      backgroundColor: event.success ? Colors.lightGreen : Colors.orange,
      content: Text(event.success ? "Bravo !" : "Dommage..."),
    ));
  }

  // process the given event
  Future<void> _processEvent(Event event) async {
    if (event is PlayerTurn) {
      _onPlayerTurn(event);
    } else if (event is DiceThrow) {
      return _onDiceThrow(event);
    } else if (event is PossibleMoves) {
      return _onPossibleMoves(event);
    } else if (event is Move) {
      _onMove(event);
    } else if (event is ShowQuestion) {
      return _onShowQuestion(event);
    } else if (event is PlayerAnswerResult) {
      _onPlayerAnswerResult(event);
    } else {
      print(event);
    }
  }

  void processEvents() async {
    const List<Event> events = [
      PlayerTurn(0),
      DiceThrow(3),
      PossibleMoves([3, 16]),
      Move(16),
      ShowQuestion("Ma belle question", 0),
      PlayerAnswerResult(0, true)
    ];

    for (var event in events) {
      await _processEvent(event);
    }
  }

  @override
  Widget build(BuildContext context) {
    return Container(
      color: Colors.grey,
      child: Column(
        children: [
          Center(
            child: AspectRatio(
              aspectRatio: 1,
              child: Board(onTapTile, highligthedTiles, state.pawnTile),
            ),
          ),
          ElevatedButton(
              onPressed: processEvents, child: const Text("Lancer !"))
        ],
      ),
      // child: TextButton(
      //     onPressed: () => _rollDice(dice.Face.two), child: const Text("Dice")),
    );
  }
}
