import 'dart:async';
import 'dart:convert';

import 'package:eleve/activities/trivialpoursuit/board.dart';
import 'package:eleve/activities/trivialpoursuit/dice.dart' as dice;
import 'package:eleve/activities/trivialpoursuit/game_debug.dart';
import 'package:eleve/activities/trivialpoursuit/game_end.dart';
import 'package:eleve/activities/trivialpoursuit/lobby.dart';
import 'package:eleve/activities/trivialpoursuit/pie.dart';
import 'package:eleve/activities/trivialpoursuit/question.dart';
import 'package:eleve/activities/trivialpoursuit/question_result.dart';
import 'package:eleve/activities/trivialpoursuit/success_recap.dart';
import 'package:eleve/build_mode.dart';
import 'package:eleve/shared/errors.dart';
import 'package:eleve/shared/settings_shared.dart';
import 'package:eleve/types/src_maths_questions_client.dart' hide Answer;
import 'package:eleve/types/src_trivial.dart';
import 'package:flutter/material.dart';
import 'package:web_socket_channel/web_socket_channel.dart';
import 'package:http/http.dart' as http;

class _LastQuestion {
  final ShowQuestion question;
  final Answers answer;
  const _LastQuestion(this.question, this.answer);
}

class GameAcces {
  final String code;
  final String studentPseudo;
  final String studentID;
  final String gameMeta;
  const GameAcces(this.code, this.studentID, this.studentPseudo, this.gameMeta);
}

const _infiniteDuration = Duration(days: 1000);

/// [GameTerminatedNotification] is emitted when a game is definitively
/// closed by the server.
class GameTerminatedNotification extends Notification {}

class TrivialPoursuitController extends StatefulWidget {
  final BuildMode buildMode;
  final GameAcces gameMeta;

  /// [isSelfLaunched] is [true] when the game was self launched by
  /// this user.
  final bool isSelfLaunched;

  static const gameMetaKey = "game-meta";

  Uri get apiURL => buildMode.websocketURL('/trivial/game/connect', query: {
        studentPseudoKey: gameMeta.studentPseudo,
        gameMetaKey: gameMeta.gameMeta
      });

  const TrivialPoursuitController(
      this.buildMode, this.gameMeta, this.isSelfLaunched,
      {Key? key})
      : super(key: key);

  @override
  _TrivialPoursuitControllerState createState() =>
      _TrivialPoursuitControllerState();
}

class _TrivialPoursuitControllerState extends State<TrivialPoursuitController>
    with WidgetsBindingObserver {
  late WebSocketChannel channel;
  late Timer _keepAliveTimmer;
  final eventQueue = StreamController<StateUpdate>();

  // the id of the client player
  PlayerID playerID = "";

  bool hasGameStarted = false;

  LobbyUpdate lobby = const LobbyUpdate({}, "", "", false, {});

  GameState state = const GameState(
      {"": PlayerStatus("", QuestionReview([], []), [], false)}, 0, "");
  Set<int> highligthedTiles = {};

  /// null when no animation is displayed
  Stream<dice.Face>? diceRollAnimation;
  bool diceDisabled = true;

  /// empty until game end
  GameEnd? gameEnd;

  double pieGlowWidth = 10;

  /// the last question shown, used to go back when displaying
  /// the result page
  _LastQuestion? lastQuestion;

  @override
  void initState() {
    if (widget.apiURL.host.isEmpty) {
      // debug only
      Future.delayed(const Duration(milliseconds: 200), processEventsDebug);
    } else {
      /// API connection
      channel = WebSocketChannel.connect(widget.apiURL);
      channel.stream
          .listen(listen, onError: _onNetworkError, onDone: _onServerDone);

      /// websocket is close in case of inactivity
      /// prevent it by sending pings
      _keepAliveTimmer = Timer.periodic(const Duration(seconds: 50), (timer) {
        _sendEvent(const Ping("keeping alive"));
      });

      // start the main event loop
      _startLoop();
    }

    WidgetsBinding.instance.addObserver(this);

    super.initState();
  }

  void processEventsDebug() async {
    for (var update in updates) {
      await processEvents(update);

      await Future<void>.delayed(const Duration(seconds: 3));
    }
  }

  void listen(dynamic event) {
    try {
      final update = stateUpdateFromJson(jsonDecode(event as String));
      eventQueue.add(update);
    } catch (e) {
      _onNetworkError(e);
    }
  }

  void _startLoop() async {
    await for (final update in eventQueue.stream) {
      await processEvents(update);
    }
  }

  @override
  void dispose() {
    if (widget.apiURL.host.isNotEmpty) {
      channel.sink.close(1000, "Bye bye");
      _keepAliveTimmer.cancel();
      eventQueue.close();
    }
    diceRollAnimation = null;

    WidgetsBinding.instance.removeObserver(this);
    super.dispose();
  }

  void _onServerDone() {
    if (!mounted) {
      return;
    }
    ScaffoldMessenger.of(context).showSnackBar(SnackBar(
      duration: const Duration(seconds: 5),
      backgroundColor: Theme.of(context).colorScheme.error,
      content: const Text("Connection interrompue."),
    ));
    popRouteToHome();
  }

  void _onNetworkError(dynamic error) {
    if (!mounted) {
      return;
    }
    ScaffoldMessenger.of(context).showSnackBar(SnackBar(
      duration: const Duration(seconds: 5),
      backgroundColor: Theme.of(context).colorScheme.error,
      content: Text("Une erreur est survenue : $error"),
    ));
    popRouteToHome();
  }

  void _sendEvent(ClientEventITF event) {
    if (widget.apiURL.host.isNotEmpty) {
      channel.sink.add(jsonEncode(clientEventITFToJson(event)));
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
      lobby = event;
    });

    if (event.iD == playerID) {
      // do not notify our own connection
      return;
    }

    ScaffoldMessenger.of(context).showSnackBar(SnackBar(
      duration: const Duration(seconds: 2),
      backgroundColor: Theme.of(context).colorScheme.primary,
      content: event.isJoining
          ? Text("${event.pseudo} a rejoint la partie !")
          : Text("${event.pseudo} a quitté la partie."),
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
    if (state.playerTurn != playerID) {
      return;
    }
    if (!highligthedTiles.contains(tile)) {
      return;
    }

    _sendEvent(ClientMove([], tile));
  }

  void onTapDice() {
    if (state.playerTurn != playerID) {
      // ignore if this it not our turn
      return;
    }

    _sendEvent(const DiceClicked());
  }

  Future<void> _onPlayerTurn(PlayerTurn event) async {
    // reset the last question
    lastQuestion = null;

    // remove potential dialog
    Navigator.of(context).popUntil(ModalRoute.withName("/board"));

    final isOwnTurn = event.player == playerID;
    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(
        duration: _infiniteDuration, // hide in _onDiceThrow
        backgroundColor: Theme.of(context).colorScheme.secondary,
        content: isOwnTurn
            ? const Text("C'est à toi de lancer le dé !")
            : Text("Au tour de ${event.playerName}"),
        action: SnackBarAction(label: "OK", onPressed: _removeCurrentSnackbar),
      ),
    );

    setState(() {
      diceDisabled = !isOwnTurn;
    });
  }

  void _removeCurrentSnackbar() {
    if (!mounted) return;
    ScaffoldMessenger.of(context).hideCurrentSnackBar();
  }

  // triggers and wait for a dice roll
  // with the given value
  Future<void> _onDiceThrow(DiceThrow event) async {
    ScaffoldMessenger.of(context).hideCurrentSnackBar();

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
    if (!mounted) return;
    ScaffoldMessenger.of(context).hideCurrentSnackBar();

    setState(() {
      highligthedTiles.clear();
    });

    for (var tile in event.path) {
      setState(() {
        state = GameState(state.players, tile, state.playerTurn);
      });

      await Future<void>.delayed(const Duration(milliseconds: 800));
    }
  }

  Future<void> _onShowQuestion(ShowQuestion event) async {
    lastQuestion = _LastQuestion(event, {});

    Navigator.of(context).push(MaterialPageRoute<void>(
      settings: const RouteSettings(name: "/question"),
      builder: (context) => InGameQuestionRoute(
        event,
        (a) {
          // do not close the page now, it is handled when receiving result
          lastQuestion = _LastQuestion(event, a.data);
          _sendEvent(Answer(a));
        },
      ),
    ));
  }

  void _showLastQuestion() {
    if (lastQuestion == null) {
      return;
    }

    Navigator.of(context).push(MaterialPageRoute<void>(
      settings: const RouteSettings(name: "/last-question"),
      builder: (context) => LastQuestionRoute(
        lastQuestion!.question,
        () {
          Navigator.of(context).pop();
        },
        lastQuestion!.answer,
      ),
    ));
  }

  Future<void> _onPlayerAnswerResults(PlayerAnswerResults event) async {
    // close the additional routes (question or recap)
    // until the "main" board
    Navigator.of(context).popUntil(ModalRoute.withName("/board"));

    var wantNextTurn =
        await Navigator.of(context).push(MaterialPageRoute<WantNextTurn>(
      settings: const RouteSettings(name: "/answer"),
      builder: (context) => NotificationListener<WantNextTurnNotification>(
        onNotification: (notification) {
          Navigator.pop(context, notification.event);
          return true;
        },
        child:
            QuestionResult(playerID, event, state.players, _showLastQuestion),
      ),
    ));
    wantNextTurn ??= const WantNextTurn(false);
    _sendEvent(wantNextTurn);
  }

  Future<void> _onPlayersStillInQuestionResult(
      PlayersStillInQuestionResult event) async {
    // ignore the event if we are one of the waiting players
    if (event.players.contains(playerID)) {
      return;
    }

    // close a potential previous dialog
    Navigator.of(context).popUntil(ModalRoute.withName("/board"));

    showDialog<void>(
        context: context,
        barrierColor: Colors.transparent,
        builder: (context) => WaitForPlayersDialog(event.playerNames),
        barrierDismissible: false);
  }

  void _onPlayerReconnected(PlayerReconnected event) {
    ScaffoldMessenger.of(context).showSnackBar(SnackBar(
        duration: const Duration(seconds: 2),
        backgroundColor: Theme.of(context).colorScheme.primary,
        content: Text("${event.pseudo} s'est reconnecté(e) !")));

    // our own playerID is empty only if we are reconnecting
    if (playerID == "") {
      setState(() {
        playerID = event.iD;
      });
      _onGameStart();
    }
  }

  void _onGameEnd(GameEnd event) {
    // remove potential dialog
    Navigator.of(context).popUntil(ModalRoute.withName("/board"));

    setState(() {
      hasGameStarted = false;
      gameEnd = event;
    });

    GameTerminatedNotification().dispatch(context);
  }

  void popRouteToHome() {
    Navigator.of(context).popUntil((route) {
      if (route.settings.name == null || route.settings.name!.isEmpty) {
        return true;
      }
      return false;
    });
  }

  void _onGameTerminated() {
    ScaffoldMessenger.of(context).hideCurrentSnackBar();
    ScaffoldMessenger.of(context).showSnackBar(SnackBar(
        duration: const Duration(seconds: 5),
        backgroundColor: Theme.of(context).colorScheme.secondary,
        content: const Text("La partie a été interrompue par son créateur.")));

    popRouteToHome();

    GameTerminatedNotification().dispatch(context);
  }

  void _showSuccessRecap() {
    Navigator.of(context).push(
      MaterialPageRoute<void>(
          builder: (context) => SuccessRecapScaffold(state.players)),
    );
  }

  // process the given event
  Future<void> _processEvent(ServerEvent event) async {
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
    } else if (event is PlayerAnswerResults) {
      return _onPlayerAnswerResults(event);
    } else if (event is PlayersStillInQuestionResult) {
      return _onPlayersStillInQuestionResult(event);
    } else if (event is GameEnd) {
      _onGameEnd(event);
    } else if (event is GameTerminated) {
      _onGameTerminated();
    } else if (event is PlayerReconnected) {
      _onPlayerReconnected(event);
    } else {
      // exhaustive switch
      throw Exception("unexpected event type ${event.runtimeType}");
    }
  }

  Future<void> processEvents(StateUpdate update) async {
    if (update.events.any((element) => element is PlayerReconnected)) {
      // update the state before triggering game start,
      // since it is needed when reconnecting
      state = update.state;
    }

    if (widget.buildMode != BuildMode.production) {
      print("Processing ${update.events.map((e) => e.runtimeType)}...");
    }

    for (var event in update.events) {
      await _processEvent(event);
    }

    if (widget.buildMode != BuildMode.production) {
      print("Done.");
    }

    setState(() {
      state = update.state;
    });
  }

  // for self launched games, perform an http call to make the server start the game
  void _startGame() async {
    final uri = widget.buildMode
        .serverURL("/api/student/trivial/selfaccess/start", query: {
      "game-id": widget.gameMeta.code,
    });

    try {
      final resp = await http.get(uri);
      checkServerError(resp.body);
    } catch (e) {
      showError("Impossible de démarrer la partie.", e, context);
      return;
    }
  }

  Widget get _game {
    if (gameEnd != null) {
      return GameEndPannel(widget.buildMode, gameEnd!, state.players, playerID);
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
        : GameLobby(lobby.playerPseudos, lobby.playerRanks, playerID,
            widget.isSelfLaunched ? _startGame : null);
  }

  @override
  Widget build(BuildContext context) {
    return GestureDetector(
      // simplify developpement
      onDoubleTap: widget.apiURL.host.isEmpty ? processEventsDebug : null,
      child: AnimatedSwitcher(
        duration: const Duration(seconds: 3),
        child: _game,
      ),
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

  Future<bool> _confirmLeave(BuildContext context) async {
    final res = await showDialog<bool?>(
        context: context,
        builder: (BuildContext context) {
          return AlertDialog(
            title: const Text("Quitter la partie"),
            content: const Text("Es-tu sûr de vouloir quitter la partie ?"),
            actions: [
              TextButton(
                  child: const Text("Quitter"),
                  onPressed: () {
                    Navigator.pop(context, true);
                  })
            ],
          );
        });
    final leave = res ?? false;
    if (leave) {
      // cleanup potential snacks
      ScaffoldMessenger.of(context).clearSnackBars();
    }
    return leave;
  }

  @override
  Widget build(BuildContext context) {
    return WillPopScope(
      onWillPop: () => _confirmLeave(context),
      child: Container(
        padding: const EdgeInsets.symmetric(horizontal: 5, vertical: 4),
        decoration: BoxDecoration(
          image: DecorationImage(
              image: const AssetImage("assets/images/grey-wood.png"),
              fit: BoxFit.cover,
              colorFilter: ColorFilter.mode(
                  const Color.fromARGB(255, 57, 115, 119).withOpacity(0.6),
                  BlendMode.srcATop)),
        ),
        child: Column(
          mainAxisAlignment: MainAxisAlignment.spaceEvenly,
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            Row(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                pie,
                Expanded(child: SizedBox(height: 95, child: recapRow))
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
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                Padding(
                  padding:
                      const EdgeInsets.symmetric(horizontal: 15, vertical: 10),
                  child: FloatingActionButton(
                    foregroundColor: const Color.fromARGB(255, 60, 209, 255),
                    backgroundColor: const Color.fromARGB(255, 110, 171, 182),
                    onPressed: () => _showRules(context),
                    tooltip: "Afficher la règle du jeu",
                    child: const Icon(
                      IconData(0xe33d, fontFamily: 'MaterialIcons'),
                      size: 40,
                    ),
                  ),
                ),
                Padding(
                    padding: const EdgeInsets.only(right: 10),
                    child:
                        dice.Dice(onTapDice, diceRollAnimation, diceDisabled)),
              ],
            ),
          ],
        ),
      ),
    );
  }
}

class TrivialActivityIcon extends StatelessWidget {
  final void Function() onTap;

  const TrivialActivityIcon(this.onTap, {Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Column(
      children: [
        Pie.asButton(onTap, 5, Categorie.values.map((e) => true).toList()),
        const Padding(
          padding: EdgeInsets.only(top: 6, bottom: 6),
          child: Text("Isy'Triv"),
        ),
      ],
    );
  }
}
