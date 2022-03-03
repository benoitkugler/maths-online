import 'dart:async';
import 'dart:convert';

import 'package:eleve/trivialpoursuit/board.dart';
import 'package:eleve/trivialpoursuit/dice.dart' as dice;
import 'package:eleve/trivialpoursuit/events.gen.dart';
import 'package:eleve/trivialpoursuit/lobby.dart';
import 'package:eleve/trivialpoursuit/pie.dart';
import 'package:eleve/trivialpoursuit/question.dart';
import 'package:flutter/material.dart';
import 'package:web_socket_channel/web_socket_channel.dart';

const devMode = bool.fromEnvironment("dev");

class TrivialPoursuitController extends StatefulWidget {
  final int questionTimeout; // in seconds

  const TrivialPoursuitController(this.questionTimeout, {Key? key})
      : super(key: key);

  @override
  _TrivialPoursuitControllerState createState() =>
      _TrivialPoursuitControllerState();
}

final _wsApi = Uri.parse('ws://localhost:8080/trivial-poursuit');

class _TrivialPoursuitControllerState extends State<TrivialPoursuitController> {
  late WebSocketChannel channel;

  int playerID = 0;
  bool hasGameStarted = false;

  Map<int, String> lobby = {};

  GameState state = const GameState({
    0: [false, false, false, false, false]
  }, 0, 0);
  Set<int> highligthedTiles = {};

  /// null when no animation is displayed
  Stream<dice.Face>? diceRollAnimation;
  dice.Face diceResult = dice.Face.one;
  bool diceDisabled = true;

  bool hasQuestion = false;

  /// empty until game end
  List<int> winners = [];
  List<String> winnerNames = [];

  @override
  void initState() {
    if (devMode) {
      // debug only
      Future.delayed(const Duration(milliseconds: 200), processEventsDebug);
    } else {
      /// API connection
      channel = WebSocketChannel.connect(_wsApi);
      channel.stream.listen(listen, onError: showError);
    }

    super.initState();
  }

  void processEventsDebug() async {
    await processEvents([
      GameEvents(const [
        PlayerJoin(0),
        GameStart(),
        DiceThrow(2),
        PossibleMoves(0, [1, 2, 3]),
        Move([0, 1, 2], 2),
        ShowQuestion("test", Categorie.blue),
      ], state),
    ]);

    await Future<void>.delayed(const Duration(seconds: 1));

    await processEvents([
      GameEvents(const [
        PlayerAnswerResult(0, true),
        GameEnd([0], ["Pierre"])
      ], state),
    ]);
  }

  @override
  void dispose() {
    if (!devMode) {
      channel.sink.close(1000, "Bye bye");
    }
    diceRollAnimation = null;
    super.dispose();
  }

  void showError(dynamic error) {
    ScaffoldMessenger.of(context).showSnackBar(SnackBar(
      duration: const Duration(seconds: 5),
      backgroundColor: Theme.of(context).colorScheme.error,
      content: Text("Une erreur est survenue : $error"),
    ));
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
      content: const Text("Connecté au serveur."),
    ));
  }

  void _onLobbyUpdate(LobbyUpdate event) {
    setState(() {
      lobby = event.names;
    });

    if (event.player == playerID) {
      // do not notify our own connection
      return;
    }

    ScaffoldMessenger.of(context).showSnackBar(SnackBar(
      duration: const Duration(seconds: 2),
      backgroundColor: Theme.of(context).colorScheme.primary,
      content: event.isJoining
          ? Text("${event.playerName} a rejoint la partie !")
          : Text("${event.playerName} a quitté la partie."),
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

    _sendEvent(Move([], tile));
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

  // triggers and wait for a dice roll
  // with the given value
  Future<void> _onDiceThrow(DiceThrow event) async {
    final face =
        dice.Face.values[event.face - 1]; // event.face is the "human" face

    final completer = Completer<void>();
    final diceRoll = dice.Dice.rollDice().asBroadcastStream();
    diceRoll.listen(null, onDone: completer.complete);

    setState(() {
      diceRollAnimation = diceRoll;
      diceDisabled = false;
    });

    await completer.future;

    // make the dice result more visible
    setState(() {
      diceRollAnimation = null;
      diceResult = face;
    });
    await Future<void>.delayed(const Duration(seconds: 1));

    setState(() {
      diceDisabled = true;
    });

    return;
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

  void _onMove(Move event) async {
    setState(() {
      highligthedTiles.clear();
    });

    for (var tile in event.path) {
      setState(() {
        state = GameState(state.successes, tile, state.player);
      });

      await Future<void>.delayed(const Duration(milliseconds: 500));
    }
  }

  Future<void> _onShowQuestion(ShowQuestion event) async {
    Navigator.of(context).push(MaterialPageRoute<void>(
      builder: (context) => NotificationListener<SubmitAnswerNotification>(
        onNotification: (notification) {
          if (hasQuestion) {
            Navigator.pop(context);
            hasQuestion = false;
          }
          _sendEvent(Answer(notification.answer));
          return true;
        },
        child: QuestionRoute(event, Duration(seconds: widget.questionTimeout)),
      ),
    ));
    hasQuestion = true;
  }

  void _onPlayerAnswerResult(PlayerAnswerResult event) {
    // for now, we simply ignore other player success
    if (event.player != playerID) {
      return;
    }

    // close the question on timeout
    if (hasQuestion) {
      Navigator.of(context).pop();
      hasQuestion = false;
    }

    ScaffoldMessenger.of(context).showSnackBar(SnackBar(
      duration: const Duration(seconds: 3),
      backgroundColor: event.success ? Colors.lightGreen : Colors.orange,
      content: Text(event.success ? "Bravo !" : "Dommage..."),
    ));
  }

  void _onGameEnd(GameEnd event) {
    setState(() {
      hasGameStarted = false;
      winners = event.winners;
      winnerNames = event.winnerNames;
    });
  }

  // process the given event
  Future<void> _processEvent(GameEvent event) async {
    if (event is PlayerJoin) {
      _onPlayerJoin(event);
    } else if (event is LobbyUpdate) {
      _onLobbyUpdate(event);
    } else if (event is GameStart) {
      _onGameStart();
    } else if (event is PlayerTurn) {
      return _onPlayerTurn(event);
    } else if (event is DiceThrow) {
      return _onDiceThrow(event);
    } else if (event is PossibleMoves) {
      return _onPossibleMoves(event);
    } else if (event is Move) {
      return _onMove(event);
    } else if (event is ShowQuestion) {
      return _onShowQuestion(event);
    } else if (event is PlayerAnswerResult) {
      _onPlayerAnswerResult(event);
    } else if (event is GameEnd) {
      _onGameEnd(event);
    } else {
      throw Exception("unexpected event type ${event.runtimeType}");
    }
  }

  Future<void> processEvents(List<GameEvents> eventList) async {
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
    if (winners.isNotEmpty) {
      return _GameEnd(winnerNames, winners.contains(playerID));
    }

    return hasGameStarted
        ? _GameStarted(
            state.successes[playerID]!,
            diceRollAnimation,
            diceResult,
            diceDisabled,
            onTapTile,
            highligthedTiles,
            state.pawnTile)
        : GameLobby(lobby, playerID);
  }

  @override
  Widget build(BuildContext context) {
    return AnimatedSwitcher(
      child: _game,
      duration: const Duration(seconds: 3),
    );
  }
}

class _GameStarted extends StatelessWidget {
  final Success success;

  final Stream<dice.Face>? diceRollAnimation;
  final dice.Face diceResult;
  final bool diceDisabled;

  final OnTapTile onTapTile;
  final Set<int> availableTiles;
  final int pawnTile;

  const _GameStarted(this.success, this.diceRollAnimation, this.diceResult,
      this.diceDisabled, this.onTapTile, this.availableTiles, this.pawnTile,
      {Key? key})
      : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Container(
      color: Colors.grey,
      child: Column(
        mainAxisAlignment: MainAxisAlignment.spaceEvenly,
        children: [
          Padding(
            padding: const EdgeInsets.only(top: 10),
            child: Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                Padding(
                  padding: const EdgeInsets.only(top: 10, left: 40),
                  child: Pie(10, success),
                ),
                const Spacer(),
                Padding(
                    padding: const EdgeInsets.only(right: 30),
                    child:
                        dice.Dice(diceRollAnimation, diceResult, diceDisabled)),
              ],
            ),
          ),
          Center(
            child: AspectRatio(
              aspectRatio: 1,
              child: Board(onTapTile, availableTiles, pawnTile),
            ),
          ),
        ],
      ),
    );
  }
}

class _GameEnd extends StatelessWidget {
  final List<String> winners;
  final bool hasWon;
  const _GameEnd(this.winners, this.hasWon, {Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    final List<Widget> congrats = [];
    if (hasWon) {
      congrats.add(const Text("Vous avez gagné, bravo !",
          style: TextStyle(
            color: Colors.yellow,
            fontSize: 25,
          )));
    }
    return Center(
      child: Column(
        mainAxisAlignment: MainAxisAlignment.spaceEvenly,
        children: [
          const Text(
            "Partie terminée",
            style: TextStyle(fontSize: 25),
          ),
          ...congrats,
          Column(
            children: [
              const Padding(
                padding: EdgeInsets.only(bottom: 20),
                child: Text(
                  "Les gagnants sont :",
                  style: TextStyle(fontSize: 20),
                ),
              ),
              Row(
                mainAxisAlignment: MainAxisAlignment.spaceEvenly,
                children: winners
                    .map((e) => DecoratedBox(
                        decoration: const BoxDecoration(boxShadow: [
                          BoxShadow(
                            color: Colors.yellow,
                            blurRadius: 5,
                          )
                        ], borderRadius: BorderRadius.all(Radius.circular(10))),
                        child: Card(
                          child: Padding(
                            padding: const EdgeInsets.all(12.0),
                            child: Text(e),
                          ),
                        )))
                    .toList(),
              ),
            ],
          )
        ],
      ),
    );
  }
}

class GameIcon extends StatelessWidget {
  final void Function() onTap;

  const GameIcon(this.onTap, {Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Column(
      children: [
        RawMaterialButton(
          onPressed: onTap,
          elevation: 2.0,
          // fillColor: Colors.white,
          child: Pie(2, Categorie.values.map((e) => true).toList()),
          padding: const EdgeInsets.all(10.0),
          shape: const CircleBorder(),
        ),
        const Padding(
          padding: EdgeInsets.only(top: 6, bottom: 6),
          child: Text("Trivial poursuit"),
        ),
      ],
    );
  }
}
