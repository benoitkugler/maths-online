import 'dart:async';
import 'dart:convert';

import 'package:eleve/exercices/types.gen.dart' as exercices;
import 'package:eleve/trivialpoursuit/board.dart';
import 'package:eleve/trivialpoursuit/dice.dart' as dice;
import 'package:eleve/trivialpoursuit/events.gen.dart';
import 'package:eleve/trivialpoursuit/game_debug.dart';
import 'package:eleve/trivialpoursuit/game_end.dart';
import 'package:eleve/trivialpoursuit/lobby.dart';
import 'package:eleve/trivialpoursuit/pie.dart';
import 'package:eleve/trivialpoursuit/question.dart';
import 'package:eleve/trivialpoursuit/question_result.dart';
import 'package:eleve/trivialpoursuit/success_recap.dart';
import 'package:flutter/material.dart';
import 'package:web_socket_channel/web_socket_channel.dart';

class TrivialPoursuitController extends StatefulWidget {
  /// empty for no remote connection
  final String apiURL;

  const TrivialPoursuitController(this.apiURL, {Key? key}) : super(key: key);

  @override
  _TrivialPoursuitControllerState createState() =>
      _TrivialPoursuitControllerState();
}

class _TrivialPoursuitControllerState extends State<TrivialPoursuitController> {
  late WebSocketChannel channel;
  late Timer _keepAliveTimmer;

  int playerID = 0;
  bool hasGameStarted = false;

  Map<int, String> lobby = {};

  GameState state = const GameState({
    0: PlayerStatus("", [false, false, false, false, false])
  }, 0, 0);
  Set<int> highligthedTiles = {};

  /// null when no animation is displayed
  Stream<dice.Face>? diceRollAnimation;
  bool diceDisabled = true;

  /// empty until game end
  List<int> winners = [];
  List<String> winnerNames = [];

  double pieGlowWidth = 10;

  @override
  void initState() {
    if (widget.apiURL == "") {
      // debug only
      Future.delayed(const Duration(milliseconds: 200), processEventsDebug);
    } else {
      /// API connection
      channel = WebSocketChannel.connect(Uri.parse(widget.apiURL));
      channel.stream.listen(listen, onError: showError);

      /// websocket is close in case of inactivity
      /// prevent it by sending pings
      _keepAliveTimmer = Timer.periodic(const Duration(seconds: 50), (timer) {
        _sendEvent(const Ping("keeping alive"));
      });
    }

    super.initState();
  }

  void processEventsDebug() async {
    for (var update in updates) {
      await processEvents(update);

      await Future<void>.delayed(const Duration(seconds: 3));
    }
  }

  @override
  void dispose() {
    if (widget.apiURL != "") {
      channel.sink.close(1000, "Bye bye");
      _keepAliveTimmer.cancel();
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
    Navigator.of(context).pop();
  }

  void listen(dynamic event) {
    try {
      final events = listStateUpdateFromJson(jsonDecode(event as String));
      processEvents(events);
    } catch (e) {
      showError(e);
    }
  }

  void _sendEvent(ClientEventData event) {
    if (widget.apiURL != "") {
      channel.sink
          .add(jsonEncode(clientEventToJson(ClientEvent(event, playerID))));
    }
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

  void onTapDice() {
    if (state.player != playerID) {
      // ignore if this it not our turn
      return;
    }

    _sendEvent(const DiceClicked());
  }

  Future<void> _onPlayerTurn(PlayerTurn event) async {
    final isOwnTurn = event.player == playerID;
    ScaffoldMessenger.of(context).showSnackBar(SnackBar(
      duration: const Duration(seconds: 4),
      backgroundColor: Theme.of(context).colorScheme.secondary,
      content: isOwnTurn
          ? const Text("C'est à toi de lancer le dé !")
          : Text("Au tour de ${event.playerName}"),
    ));

    setState(() {
      diceDisabled = !isOwnTurn;
    });
  }

  // triggers and wait for a dice roll
  // with the given value
  Future<void> _onDiceThrow(DiceThrow event) async {
    final face =
        dice.Face.values[event.face - 1]; // event.face is the "human" face

    final completer = Completer<void>();
    final diceRoll = dice.Dice.rollDice(face).asBroadcastStream();
    diceRoll.listen(null, onDone: completer.complete);

    setState(() {
      diceRollAnimation = diceRoll;
      diceDisabled = false;
    });

    await completer.future;

    // make the dice result more visible
    setState(() {
      diceRollAnimation = null;
    });
    await Future<void>.delayed(const Duration(seconds: 1));

    setState(() {
      diceDisabled = true;
    });

    return;
  }

  Future<void> _onPossibleMoves(PossibleMoves event) async {
    final isOwnTurn = event.player == playerID;
    ScaffoldMessenger.of(context).showSnackBar(SnackBar(
      duration: const Duration(seconds: 5),
      backgroundColor: Theme.of(context).colorScheme.secondary,
      content: isOwnTurn
          ? const Text("Choisis où déplacer le pion.")
          : Text("${event.playerName} est en train de choisir la case..."),
    ));

    // only the current player may choose the tile to move
    setState(() {
      diceDisabled = true;
      highligthedTiles = event.tiles.toSet();
    });
  }

  void _onMove(Move event) async {
    ScaffoldMessenger.of(context).hideCurrentSnackBar();

    setState(() {
      highligthedTiles.clear();
    });

    for (var tile in event.path) {
      setState(() {
        state = GameState(state.players, tile, state.player);
      });

      await Future<void>.delayed(const Duration(milliseconds: 500));
    }
  }

  Future<void> _onShowQuestion(ShowQuestion event) async {
    Navigator.of(context).push(MaterialPageRoute<void>(
      settings: const RouteSettings(name: "/question"),
      builder: (context) => NotificationListener<SubmitAnswerNotification>(
        onNotification: (notification) {
          // do not close the page now, it is handled when receiving result
          _sendEvent(Answer(
              const exercices.QuestionAnswersIn({}), notification.answer));
          return true;
        },
        child: QuestionRoute(event, Duration(seconds: event.timeoutSeconds)),
      ),
    ));
  }

  Future<void> _onPlayerAnswerResult(PlayerAnswerResult event) async {
    // for now, we simply ignore other player success
    if (event.player != playerID) {
      return;
    }

    // close the additional routes (question or recap)
    // until the "main" board
    Navigator.of(context).popUntil((route) {
      if (route.settings.name == "/board") {
        return true;
      }
      return false;
    });

    Navigator.of(context).push(MaterialPageRoute<void>(
      settings: const RouteSettings(name: "/answer"),
      builder: (context) => NotificationListener<WantNextTurnNotification>(
        onNotification: (notification) {
          Navigator.pop(context);
          _sendEvent(notification.event);
          return true;
        },
        child: QuestionResult(event),
      ),
    ));
  }

  void _onGameEnd(GameEnd event) {
    setState(() {
      hasGameStarted = false;
      winners = event.winners;
      winnerNames = event.winnerNames;
    });
  }

  void _showSuccessRecap() {
    Navigator.of(context).push(
      MaterialPageRoute<void>(
          builder: (context) => SuccessRecapScaffold(state.players)),
    );
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
      return _onPlayerAnswerResult(event);
    } else if (event is GameEnd) {
      _onGameEnd(event);
    } else {
      // exhaustive switch
      throw Exception("unexpected event type ${event.runtimeType}");
    }
  }

  Future<void> processEvents(List<StateUpdate> eventList) async {
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
      return GameEndPannel(winners, state.players, playerID);
    }

    return hasGameStarted
        ? _GameStarted(
            Hero(
              tag: "recap_$playerID",
              child: Pie.asButton(_showSuccessRecap, pieGlowWidth,
                  state.players[playerID]!.success),
            ),
            SuccessRecapRow(
              playerID,
              state.players,
            ),
            onTapDice,
            diceRollAnimation,
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
  final Widget pie;
  final SuccessRecapRow recapRow;

  final void Function() onTapDice;
  final Stream<dice.Face>? diceRollAnimation;
  final bool diceDisabled;

  final OnTapTile onTapTile;
  final Set<int> availableTiles;
  final int pawnTile;

  const _GameStarted(
      this.pie,
      this.recapRow,
      this.onTapDice,
      this.diceRollAnimation,
      this.diceDisabled,
      this.onTapTile,
      this.availableTiles,
      this.pawnTile,
      {Key? key})
      : super(key: key);

  Text _ruleRow(String content, double fontSize) {
    return Text(content, style: TextStyle(fontSize: fontSize, height: 2));
  }

  void _showRules(BuildContext context) {
    Navigator.of(context).push(
      MaterialPageRoute<void>(
          builder: (context) => Scaffold(
                appBar: AppBar(title: const Text("Règles du jeu")),
                body: Container(
                  padding: const EdgeInsets.symmetric(horizontal: 20),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.stretch,
                    mainAxisAlignment: MainAxisAlignment.center,
                    children: [
                      Card(
                        color: Colors.lightGreen,
                        child: Padding(
                          padding: const EdgeInsets.all(8.0),
                          child: Column(
                              crossAxisAlignment: CrossAxisAlignment.start,
                              children: [
                                _ruleRow(
                                    "Quand je réponds correctement à une question : ",
                                    18),
                                _ruleRow(
                                    "    - si je n'ai pas encore le camembert, je le gagne !",
                                    16),
                                _ruleRow(
                                    "    - si j'avais déjà le camembert, il ne se passe rien.",
                                    16),
                              ]),
                        ),
                      ),
                      const SizedBox(height: 20),
                      Card(
                        color: Colors.red,
                        child: Padding(
                          padding: const EdgeInsets.all(8.0),
                          child: Column(
                              crossAxisAlignment: CrossAxisAlignment.start,
                              children: [
                                _ruleRow("Quand je me trompe : ", 18),
                                _ruleRow(
                                    "    - si j'ai le camembert, je le perds !",
                                    16),
                                _ruleRow(
                                    "    - si je n'ai pas encore le camembert, il ne se passe rien.",
                                    16),
                              ]),
                        ),
                      ),
                    ],
                  ),
                ),
              )),
    );
  }

  @override
  Widget build(BuildContext context) {
    return Container(
      decoration: BoxDecoration(
        image: DecorationImage(
            image: const AssetImage("lib/images/grey-wood.png"),
            fit: BoxFit.cover,
            colorFilter: ColorFilter.mode(
                Colors.grey.withOpacity(0.6), BlendMode.srcATop)),
      ),
      child: Column(
        mainAxisAlignment: MainAxisAlignment.spaceEvenly,
        children: [
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              Padding(
                padding: const EdgeInsets.only(top: 5, left: 40, bottom: 10),
                child: pie,
              ),
              const Spacer(),
              Padding(
                  padding: const EdgeInsets.only(right: 30),
                  child: dice.Dice(onTapDice, diceRollAnimation, diceDisabled)),
            ],
          ),
          Expanded(child: Center(
            child: LayoutBuilder(
              builder: (_, cts) {
                return Board(cts.biggest.shortestSide, onTapTile,
                    availableTiles, pawnTile);
              },
            ),
          )),
          Padding(
            padding: const EdgeInsets.only(top: 5, bottom: 2),
            child: Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                Padding(
                  padding:
                      const EdgeInsets.symmetric(horizontal: 15, vertical: 10),
                  child: FloatingActionButton(
                    foregroundColor: Colors.lightBlue,
                    backgroundColor: Colors.lightBlue.shade100,
                    onPressed: () => _showRules(context),
                    child: const Icon(
                      IconData(0xe33d, fontFamily: 'MaterialIcons'),
                      size: 40,
                    ),
                    tooltip: "Afficher la règle du jeu",
                  ),
                ),
                recapRow,
              ],
            ),
          ),
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
        Pie.asButton(onTap, 5, Categorie.values.map((e) => true).toList()),
        const Padding(
          padding: EdgeInsets.only(top: 6, bottom: 6),
          child: Text("Trivial poursuit"),
        ),
      ],
    );
  }
}
