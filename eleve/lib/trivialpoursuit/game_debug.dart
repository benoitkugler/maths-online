import 'package:eleve/questions/types.gen.dart';
import 'package:eleve/trivialpoursuit/events.gen.dart';

/// used to collect shader warmup
const typicalUpdates = [
  StateUpdate(
      [],
      GameState({
        "0": PlayerStatus("Annonymous 065686", QuestionReview([], []),
            [true, true, false, true, false], false),
        "1": PlayerStatus("Annonymous 065686", QuestionReview([], []),
            [false, false, false, false, false], false),
      }, 0, "0")),
  StateUpdate(
      [
        PlayerJoin("0"),
        GameStart(),
        PlayerTurn("", "0"),
        DiceThrow(2),
        PossibleMoves("Katia", [1, 3, 7], "0"),
        Move([0, 1, 2, 3, 4, 5], 5),
        ShowQuestion(60, Categorie.orange, 0, Question("Test", [])),
        PlayerAnswerResults(Categorie.orange, {
          "0": PlayerAnswerResult(false, false),
          "1": PlayerAnswerResult(false, false),
        }),
        PlayerTurn("Ben", "0"),
        GameEnd({
          "0": [24, 49]
        }, [
          "0",
          "1"
        ], [
          "Pierre",
          "Benoit"
        ])
      ],
      GameState({
        "0": PlayerStatus("Annonymous 065686", QuestionReview([], []),
            [true, true, false, true, false], false),
        "1": PlayerStatus("Annonymous 065686", QuestionReview([], []),
            [false, false, false, false, false], false),
        "2": PlayerStatus("Annonymous 065686", QuestionReview([], []),
            [false, false, false, false, false], false),
        "3": PlayerStatus("Annonymous 065686", QuestionReview([], []),
            [false, false, false, false, false], false),
        "4": PlayerStatus("Annonymous 065686", QuestionReview([], []),
            [false, false, false, false, false], false),
      }, 0, "0")),
  StateUpdate(
      [
        Move([0, 1, 2, 3, 4, 5], 5),
        DiceThrow(2)
      ],
      GameState({
        "0": PlayerStatus("Player 1", QuestionReview([], []),
            [true, true, true, true, true], false),
        "1": PlayerStatus("Player 2", QuestionReview([], []),
            [true, false, false, false, false], false),
        "2": PlayerStatus("Player 3", QuestionReview([], []),
            [false, true, false, false, false], false),
        "3": PlayerStatus("Annonymous 065686", QuestionReview([], []),
            [false, false, false, false, false], false),
        "4": PlayerStatus("Annonymous 065686", QuestionReview([], []),
            [false, false, false, false, false], false),
      }, 0, "0"))
];

const updates = typicalUpdates;
// const updates = devUpdates;

const devUpdates = [
  StateUpdate(
      [
        PlayerJoin("0"),
        GameStart(),

        // PlayerTurn("", 0),
        // DiceThrow(2),
        // PossibleMoves("Katia", [1, 3, 7], 0),
        Move([0, 1, 2], 5),
        // ShowQuestion(60, Categorie.orange, 0, Question("Test", [])),
        // PlayerAnswerResults(
        //     Categorie.orange, {0: PlayerAnswerResult(false, false)}),
        // PlayerTurn("Ben", 0),
        // GameEnd({
        //   0: [24, 49]
        // }, [
        //   0,
        //   1
        // ], [
        //   "Pierre",
        //   "Benoit"
        // ])
      ],
      GameState({
        "0": PlayerStatus("Annonymous 065686", QuestionReview([], []),
            [true, true, false, true, false], false),
        "1": PlayerStatus("Annonymous 065686", QuestionReview([], []),
            [false, false, false, false, false], false),
        "2": PlayerStatus("Annonymous 065686", QuestionReview([], []),
            [false, false, false, false, false], false),
        "3": PlayerStatus("Annonymous 065686", QuestionReview([], []),
            [false, false, false, false, false], false),
        "4": PlayerStatus("Annonymous 065686", QuestionReview([], []),
            [false, false, false, false, false], false),
      }, 0, "0")),
];
