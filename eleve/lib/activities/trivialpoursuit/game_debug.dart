import 'package:eleve/types/src_maths_questions_client.dart';
import 'package:eleve/types/src_trivial.dart';

/// used to collect shader warmup
const typicalUpdates = [
  StateUpdate(
      [],
      GameState({
        "0": PlayerStatus("Annonymous 065686", QuestionReview([], []),
            [true, true, false, true, false], false, 2),
        "1": PlayerStatus("Annonymous 065686", QuestionReview([], []),
            [false, false, false, false, false], false, 2),
      }, "0", 0)),
  StateUpdate(
      [
        PlayerJoin("0"),
        GameStart(),
        PlayerTurn("", "0"),
        DiceThrow(2),
        PossibleMoves("Katia", [1, 3, 7], "0"),
        Move([0, 1, 2, 3, 4, 5], 5),
        ShowQuestion(60, Categorie.orange, 0, Question([], [])),
        PlayerAnswerResults(Categorie.orange, {
          "0": PlayerAnswerResult(false, false),
          "1": PlayerAnswerResult(false, false),
        }, {}),
        PlayersStillInQuestionResult(["1"], ["Katia"]),
      ],
      GameState({
        "0": PlayerStatus("Annonymous 065686", QuestionReview([], []),
            [true, true, false, true, false], false, 0),
        "1": PlayerStatus("Annonymous 065686", QuestionReview([], []),
            [false, false, false, false, false], false, 1),
        "2": PlayerStatus("Annonymous 065686", QuestionReview([], []),
            [false, false, false, false, false], false, 2),
        "3": PlayerStatus("Annonymous 065686", QuestionReview([], []),
            [false, false, false, false, false], false, 3),
        "4": PlayerStatus("Annonymous 065686", QuestionReview([], []),
            [false, false, false, false, false], false, 4),
      }, "0", 0)),
  StateUpdate(
      [
        PlayerTurn("Ben", "0"),
        GameEnd({
          "0": [24, 49]
        }, [
          "0",
          "1"
        ], [
          "Pierre",
          "Benoit"
        ], {})
      ],
      GameState(
        {
          "0": PlayerStatus("Annonymous 065686", QuestionReview([], []),
              [true, true, false, true, false], false, 0),
          "1": PlayerStatus("Annonymous 065686", QuestionReview([], []),
              [false, false, false, false, false], false, 1),
          "2": PlayerStatus("Annonymous 065686", QuestionReview([], []),
              [false, false, false, false, false], false, 2),
          "3": PlayerStatus("Annonymous 065686", QuestionReview([], []),
              [false, false, false, false, false], false, 3),
          "4": PlayerStatus("Annonymous 065686", QuestionReview([], []),
              [false, false, false, false, false], false, 4),
        },
        "0",
        0,
      )),
];

const updates = typicalUpdates;
// const updates = devUpdates;

const devUpdates = [
  StateUpdate(
      [
        PlayerJoin("0"),
        GameStart(),
        // Move([0, 1, 2], 5),
      ],
      GameState({
        "0": PlayerStatus("Annonymous 065686", QuestionReview([], []),
            [true, true, false, true, false], false, 0),
        "1": PlayerStatus("Annonymous 065686", QuestionReview([], []),
            [false, false, false, false, false], false, 0),
        "2": PlayerStatus("Annonymous 065686", QuestionReview([], []),
            [false, false, false, false, false], false, 0),
      }, "0", 0)),
  StateUpdate(
      [
        PlayersStillInQuestionResult(["1", "2"], ["Bubeu", "Guigui"]),
      ],
      GameState({
        "0": PlayerStatus("Annonymous 065686", QuestionReview([], []),
            [true, true, false, true, false], false, 0),
        "1": PlayerStatus("Annonymous 065686", QuestionReview([], []),
            [false, false, false, false, false], false, 0),
        "2": PlayerStatus("Annonymous 065686", QuestionReview([], []),
            [false, false, false, false, false], false, 0),
      }, "0", 0)),
  StateUpdate(
      [
        PlayersStillInQuestionResult(["1"], ["Guigui"]),
      ],
      GameState({
        "0": PlayerStatus("Annonymous 065686", QuestionReview([], []),
            [true, true, false, true, false], false, 0),
        "1": PlayerStatus("Annonymous 065686", QuestionReview([], []),
            [false, false, false, false, false], false, 0),
        "2": PlayerStatus("Annonymous 065686", QuestionReview([], []),
            [false, false, false, false, false], false, 0),
      }, "0", 0)),
  // StateUpdate(
  //     [
  //       PlayerAnswerResults(Categorie.orange, {
  //         "0": PlayerAnswerResult(false, false),
  //         "1": PlayerAnswerResult(false, false),
  //         "2": PlayerAnswerResult(false, false),
  //       }),
  //     ],
  //     GameState({
  //       "0": PlayerStatus("Annonymous 065686", QuestionReview([], []),
  //           [true, true, false, true, false], false),
  //       "1": PlayerStatus("Annonymous 065686", QuestionReview([], []),
  //           [false, false, false, false, false], false),
  //       "2": PlayerStatus("Annonymous 065686", QuestionReview([], []),
  //           [false, false, false, false, false], false),
  //     }, "0",0 )),
];
