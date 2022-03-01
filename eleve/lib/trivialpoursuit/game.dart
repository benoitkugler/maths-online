import 'dart:async';
import 'dart:convert';

import 'package:eleve/trivialpoursuit/board.dart';
import 'package:eleve/trivialpoursuit/dice.dart' as dice;
import 'package:eleve/trivialpoursuit/events.gen.dart';
import 'package:eleve/trivialpoursuit/pie.dart';
import 'package:eleve/trivialpoursuit/question.dart';
import 'package:flutter/material.dart';
import 'package:web_socket_channel/web_socket_channel.dart';

class GameController extends StatefulWidget {
  const GameController({Key? key}) : super(key: key);

  @override
  _GameControllerState createState() => _GameControllerState();
}

class _GameControllerState extends State<GameController> {
  late WebSocketChannel channel;

  int playerID = 0;
  bool hasGameStarted = false;

  GameState state = const GameState({
    0: [false, true, false, true, true]
  }, 0, 0);
  Set<int> highligthedTiles = {};

  bool showDiceRoll = false;
  dice.Face diceResult = dice.Face.one;
  Completer<void> diceCompleter = Completer();

  @override
  void initState() {
    /// API connection
    channel = WebSocketChannel.connect(
      Uri.parse('ws://localhost:8080/trivial-poursuit'),
    );

    channel.stream.listen(listen, onError: showError);
    super.initState();
  }

  void showError(dynamic error) {
    print("ERROR: $error");
  }

  void listen(dynamic event) {
    try {
      final events = listGameEventsFromJson(jsonDecode(event as String));
      processEvents(events);
    } catch (e) {
      showError(e);
    }
  }

  void _sendEvent(ClientEventData event) {
    channel.sink
        .add(jsonEncode(clientEventToJson(ClientEvent(event, playerID))));
  }

  void _onPlayerJoin(PlayerJoin event) {
    // PlayerJoin is only emitted to the player who actually joined
    setState(() {
      playerID = event.player;
    });

    ScaffoldMessenger.of(context).showSnackBar(SnackBar(
      duration: const Duration(seconds: 3),
      backgroundColor: Theme.of(context).colorScheme.primary,
      content: const Text("Connection réussie !"),
    ));
  }

  void _onGameStart() {
    setState(() {
      hasGameStarted = true;
    });
  }

  void onTapTile(int tile) {
    // send the move request to the server,
    // ignore event if it not our turn to play
    // or if the tile is not selected
    if (state.player != playerID) {
      return;
    }
    if (!highligthedTiles.contains(tile)) {
      return;
    }

    _sendEvent(Move(tile));
  }

  Future<void> _onPlayerTurn(PlayerTurn event) async {
    ScaffoldMessenger.of(context).showSnackBar(SnackBar(
      duration: const Duration(seconds: 3),
      backgroundColor: Theme.of(context).colorScheme.secondary,
      content: event.player == playerID
          ? const Text("C'est à toi !")
          : Text("Au tour de ${event.playerName}"),
    ));
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

    setState(() {
      diceCompleter = Completer<void>();
      diceResult = face;
      showDiceRoll = true;
    });

    diceCompleter.future.then((value) {
      setState(() {
        showDiceRoll = false;
      });
    });
    return diceCompleter.future;
  }

  Future<void> _onPossibleMoves(PossibleMoves event) async {
    // only the current player may choose the tile to move
    if (playerID != event.currentPlayer) {
      return;
    }
    setState(() {
      highligthedTiles = event.tiles.toSet();
    });
  }

  void _onMove(Move event) {
    // TODO: animated
    setState(() {
      highligthedTiles.clear();
    });
  }

  Future<void> _onShowQuestion(ShowQuestion event) async {
    // TODO: handle timeout and cancelation correctly;
    final widget = NotificationListener<SubmitAnswerNotification>(
        onNotification: (notification) {
          Navigator.pop(context);
          _sendEvent(Answer(notification.answer));
          return true;
        },
        child: QuestionRoute(event));
    return _showRoute(widget);
  }

  void _onPlayerAnswerResult(PlayerAnswerResult event) {
    // for now, we simply ignore other player success
    if (event.player != playerID) {
      return;
    }

    ScaffoldMessenger.of(context).showSnackBar(SnackBar(
      duration: const Duration(seconds: 3),
      backgroundColor: event.success ? Colors.lightGreen : Colors.orange,
      content: Text(event.success ? "Bravo !" : "Dommage..."),
    ));
  }

  // process the given event
  Future<void> _processEvent(GameEvent event) async {
    if (event is PlayerJoin) {
      _onPlayerJoin(event);
    } else if (event is GameStart) {
      _onGameStart();
    } else if (event is PlayerTurn) {
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
      throw Exception("unexpected event type ${event.runtimeType}");
    }
  }

  void processEvents(List<GameEvents> eventList) async {
    for (var events in eventList) {
      for (var event in events.events) {
        await _processEvent(event);
      }
      setState(() {
        state = events.state;
      });
    }
  }

  Widget get _game {
    return hasGameStarted
        ? _GameStarted(state.successes[playerID]!, showDiceRoll, diceResult,
            diceCompleter, onTapTile, highligthedTiles, state.pawnTile)
        : const _GameLobby();
  }

  @override
  Widget build(BuildContext context) {
    return AnimatedSwitcher(
      child: _game,
      duration: const Duration(seconds: 3),
    );
  }
}

class _GameLobby extends StatelessWidget {
  const _GameLobby({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Center(
      child: Column(
        mainAxisAlignment: MainAxisAlignment.spaceEvenly,
        children: const [
          Text(
            "En attente d'autres joueurs...",
            style: TextStyle(fontSize: 20),
          ),
          CircularProgressIndicator(),
        ],
      ),
    );
  }
}

class _GameStarted extends StatelessWidget {
  final Success success;

  final bool showDiceRoll;
  final dice.Face diceResult;
  final Completer<void> diceCompleter;

  final OnTapTile onTapTile;
  final Set<int> availableTiles;
  final int pawnTile;

  const _GameStarted(this.success, this.showDiceRoll, this.diceResult,
      this.diceCompleter, this.onTapTile, this.availableTiles, this.pawnTile,
      {Key? key})
      : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Container(
      color: Colors.grey,
      child: Column(
        children: [
          Padding(
            padding: const EdgeInsets.only(top: 25),
            child: Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                Padding(
                  padding: const EdgeInsets.only(top: 10),
                  child: Pie(success),
                ),
                const Spacer(),
                Padding(
                  padding: const EdgeInsets.only(right: 30),
                  child: showDiceRoll
                      ? NotificationListener<dice.DoneRolling>(
                          child: dice.DiceRoll(diceResult),
                          onNotification: (n) {
                            diceCompleter.complete();
                            return true;
                          },
                        )
                      : dice.DisabledDice(diceResult),
                ),
              ],
            ),
          ),
          Expanded(
            flex: 1,
            child: Center(
              child: AspectRatio(
                aspectRatio: 1,
                child: Board(onTapTile, availableTiles, pawnTile),
              ),
            ),
          ),
        ],
      ),
    );
  }
}
