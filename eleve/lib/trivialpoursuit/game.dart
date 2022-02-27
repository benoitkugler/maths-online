import 'dart:async';

import 'package:eleve/trivialpoursuit/board.dart';
import 'package:eleve/trivialpoursuit/dice.dart' as dice;
import 'package:eleve/trivialpoursuit/events.gen.dart';
import 'package:eleve/trivialpoursuit/pie.dart';
import 'package:eleve/trivialpoursuit/question.dart';
import 'package:flutter/material.dart';

class GameController extends StatefulWidget {
  final int playerID; // the index of the client player

  const GameController(this.playerID, {Key? key}) : super(key: key);

  @override
  _GameControllerState createState() => _GameControllerState();
}

class _GameControllerState extends State<GameController> {
  GameState state = const GameState([
    [false, true, false, true, true]
  ], 0, 0);

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

  Future<void> _onPlayerTurn(PlayerTurn event) async {
    if (event.player != widget.playerID) {
      return;
    }
    ScaffoldMessenger.of(context).showSnackBar(SnackBar(
      duration: const Duration(seconds: 3),
      backgroundColor: Theme.of(context).colorScheme.secondary,
      content: const Text("C'est à toi !"),
    ));
    return Future.delayed(const Duration(seconds: 2));
  }

  Future<void> _showRoute(Widget content) {
    return Navigator.push<void>(context, MaterialPageRoute(builder: (context) {
      return content;
    }));
  }

  // triggers and wait for a dice roll
  // with the given value
  Future<void> _onDiceThrow(DiceThrow event) async {
    final face =
        dice.Face.values[event.face - 1]; // event.face is the "natural" face

    final widget = dice.DiceRoll(face);

    final completer = Completer<void>();

    final entry = OverlayEntry(
      builder: (context) => Container(
        color: Colors.grey.withOpacity(0.5),
        child: Center(
          child: NotificationListener<dice.DoneRolling>(
            child: widget,
            onNotification: (n) {
              completer.complete();
              return true;
            },
          ),
        ),
      ),
    );

    final overlayState = Overlay.of(context)!;
    overlayState.insert(entry);

    final out = completer.future;
    out.then((value) => entry.remove());
    return out;
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
      return _onPlayerTurn(event);
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

  void processDebugEvents() {
    const List<Event> events = [
      PlayerTurn(0),
      DiceThrow(3),
      PossibleMoves([3, 16]),
      Move(16),
      ShowQuestion("Ma belle question", 0),
      PlayerAnswerResult(0, true)
    ];
    processEventRange(const EventRange(
        events,
        GameState([
          [true, true, true, true, true]
        ], 1, 0),
        0));
  }

  void processEventRange(EventRange events) async {
    for (var event in events.events) {
      await _processEvent(event);
    }
    setState(() {
      state = events.state;
    });
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
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              Pie(state.successes[widget.playerID]),
              ElevatedButton(
                  onPressed: processDebugEvents, child: const Text("Lancer !"))
            ],
          )
        ],
      ),
      // child: TextButton(
      //     onPressed: () => _rollDice(dice.Face.two), child: const Text("Dice")),
    );
  }
}
