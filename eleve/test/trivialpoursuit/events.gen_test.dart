import 'dart:convert';

import 'package:eleve/trivialpoursuit/events.gen.dart';
import 'package:test/test.dart';

void main() {
  test(
    "load events JSON",
    () {
      const input = """
      {
    "Events": [
      {
      "Data": {
        "Player": 2
      },
      "Kind": 3
      },
      {
      "Data": {
        "Face": 3
      },
      "Kind": 0
      },
      {
      "Data": {
        "Tiles": [
        3,
        16
        ]
      },
      "Kind": 4
      },
      {
      "Data": {
        "Tile": 3
      },
      "Kind": 1
      },
      {
      "Data": {
        "Question": "Super",
        "Categorie": 0
      },
      "Kind": 5
      },
      {
      "Data": {
        "Player": 0,
        "Success": true
      },
      "Kind": 2
      },
      {
      "Data": {
        "Player": 1,
        "Success": false
      },
      "Kind": 2
      },
      {
      "Data": {
        "Player": 2,
        "Success": true
      },
      "Kind": 2
      },
      {
      "Data": {
        "Player": 0
      },
      "Kind": 3
      },
      {
      "Data": {
        "Face": 3
      },
      "Kind": 0
      },
      {
      "Data": {
        "Tile": 4
      },
      "Kind": 1
      },
      {
      "Data": {
        "Question": "Super",
        "Categorie": 1
      },
      "Kind": 5
      },
      {
      "Data": {
        "Player": 0,
        "Success": false
      },
      "Kind": 2
      },
      {
      "Data": {
        "Player": 1,
        "Success": true
      },
      "Kind": 2
      },
      {
      "Data": {
        "Player": 2,
        "Success": true
      },
      "Kind": 2
      },
      {
      "Data": {
        "Player": 1
      },
      "Kind": 3
      }
    ],
    "Start": 0
    }
    """;
      final ev = eventRangeFromJson(jsonDecode(input));
      expect(ev.start, equals(0));
      expect(ev.events.length, equals(16));
      expect(ev.events[0] is PlayerTurn, equals(true));
    },
  );

  test("load state JSON", () {
    const input = """
      {
  "Question": {
    "Question": "",
    "Categorie": 0
  },
  "Successes": [
    [
    true,
    false,
    false,
    false,
    false
    ],
    [
    false,
    true,
    true,
    false,
    false
    ]
  ],
  "PawnTile": 2,
  "Player": 0,
  "Dice": 0
  }
  """;

    final state = gameStateFromJson(jsonDecode(input));
    expect(state.pawnTile, equals(2));
  });
}
