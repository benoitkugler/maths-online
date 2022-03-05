import 'dart:convert';

import 'package:eleve/trivialpoursuit/events.gen.dart';
import 'package:test/test.dart';

void main() {
  test(
    "load events JSON",
    () {
      const input = """
        [
 {
  "Events": [
   {
    "Data": {},
    "Kind": 3
   },
   {
    "Data": {
     "Player": 2,
     "PlayerName": "Haha"
    },
    "Kind": 7
   },
   {
    "Data": {
     "Face": 3
    },
    "Kind": 1
   },
   {
    "Data": {
     "CurrentPlayer": 2,
     "Tiles": [
      3,
      16
     ]
    },
    "Kind": 8
   },
   {
    "Data": {
     "Tile": 3
    },
    "Kind": 4
   },
   {
    "Data": {
     "Player": 1
    },
    "Kind": 6
   },
   {
    "Data": {
     "Question": "Super",
     "Categorie": 0
    },
    "Kind": 9
   },
   {
    "Data": {
     "Player": 0,
     "Success": true
    },
    "Kind": 5
   },
   {
    "Data": {
     "Player": 1,
     "Success": false
    },
    "Kind": 5
   },
   {
    "Data": {
     "Player": 2,
     "Success": true
    },
    "Kind": 5
   },
   {
    "Data": {
     "Player": 0,
     "PlayerName": ""
    },
    "Kind": 7
   },
   {
    "Data": {
     "Face": 3
    },
    "Kind": 1
   },
   {
    "Data": {
     "Tile": 4
    },
    "Kind": 4
   },
   {
    "Data": {
     "Question": "Super",
     "Categorie": 1
    },
    "Kind": 9
   },
   {
    "Data": {
     "Player": 0,
     "Success": false
    },
    "Kind": 5
   },
   {
    "Data": {
     "Player": 1,
     "Success": true
    },
    "Kind": 5
   },
   {
    "Data": {
     "Player": 2,
     "Success": true
    },
    "Kind": 5
   },
   {
    "Data": {
     "Player": 1,
     "PlayerName": ""
    },
    "Kind": 7
   }
  ],
  "State": {
   "Successes": null,
   "PawnTile": 0,
   "Player": 0
  }
 }
]
    """;
      final ev = listStateUpdateFromJson(jsonDecode(input));
      expect(ev.length, equals(1));
      expect(ev[0].events.length, equals(18));
      expect(ev[0].events[0] is GameStart, equals(true));
    },
  );

  test("load state JSON", () {
    const input = """
        {
  "Successes": {
    "0": [
    true,
    false,
    false,
    false,
    false
    ],
    "1": [
    false,
    true,
    true,
    false,
    false
    ]
  },
  "PawnTile": 2,
  "Player": 0
  }
  """;

    final state = gameStateFromJson(jsonDecode(input));
    expect(state.pawnTile, equals(2));
  });
}
