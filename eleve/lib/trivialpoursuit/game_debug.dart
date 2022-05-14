import 'package:eleve/trivialpoursuit/events.gen.dart';

const updates = [
  StateUpdate(
      [
        // PlayerJoin(0),
        GameStart(),
        // PlayerTurn("", 0),
        // DiceThrow(2),
        // PossibleMoves("Katia", [1, 3, 7], 0),
        // Move([0, 1, 2, 3, 4, 5], 5),
        // ShowQuestion("test", 60, Categorie.orange),
        // PlayerAnswerResult("tetst", 0, false, Categorie.green),
        GameEnd({
          0: [24, 49]
        }, [
          0,
          1
        ], [
          "Pierre",
          "Benoit"
        ])
      ],
      GameState({
        0: PlayerStatus(
            "Player 2", QuestionReview([], []), [true, true, true, true, true]),
        1: PlayerStatus("Player 2", QuestionReview([], []),
            [false, false, false, false, false]),
        // 2: PlayerStatus("Player 2", QuestionReview([], []),
        //     [true, true, false, false, false]),
      }, 0, 0)),
  // StateUpdate(
  //     [
  //       // Move([0, 1, 2, 3, 4, 5], 5),
  //       // DiceThrow(2)
  //       // PlayerAnswerResult(0, true),
  //       // GameEnd([0], ["Pierre"])
  //     ],
  //     GameState({
  //       0: PlayerStatus("Player 1", QuestionReview([], []),
  //           [true, false, false, true, false]),
  //       1: PlayerStatus("Player 2", QuestionReview([], []),
  //           [true, false, false, false, false]),
  //       2: PlayerStatus("Player 3", QuestionReview([], []),
  //           [false, true, false, false, false]),
  //     }, 0, 0)),
  // StateUpdate(
  //     [GameTerminated()],
  //     GameState({
  //       0: PlayerStatus(
  //           "", QuestionReview([], []), [true, false, true, true, false])
  //     }, 0, 0))
];
