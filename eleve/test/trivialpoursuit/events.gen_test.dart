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
   "Data": {},
   "Kind": 6
  },
  {
   "Data": {
    "PlayerName": "Haha",
    "Player": 2
   },
   "Kind": 9
  },
  {
   "Data": {
    "Face": 3
   },
   "Kind": 4
  },
  {
   "Data": {
    "PlayerName": "",
    "Tiles": [
     3,
     9
    ],
    "Player": 2
   },
   "Kind": 3
  },
  {
   "Data": {
    "Path": null,
    "Tile": 3
   },
   "Kind": 1
  },
  {
   "Data": {
    "Player": 1
   },
   "Kind": 8
  },
  {
   "Data": {
    "TimeoutSeconds": 0,
    "Categorie": 0,
    "ID": 1,
    "Question": {
     "Title": "",
     "Enonce": []
    }
   },
   "Kind": 10
  },
  {
   "Data": {
    "Player": 0,
    "Success": true,
    "Categorie": 0,
    "AskForMask": false
   },
   "Kind": 7
  },
  {
   "Data": {
    "Player": 1,
    "Success": false,
    "Categorie": 0,
    "AskForMask": false
   },
   "Kind": 7
  },
  {
   "Data": {
    "Player": 2,
    "Success": true,
    "Categorie": 0,
    "AskForMask": false
   },
   "Kind": 7
  },
  {
   "Data": {
    "PlayerName": "",
    "Player": 0
   },
   "Kind": 9
  },
  {
   "Data": {
    "Face": 3
   },
   "Kind": 4
  },
  {
   "Data": {
    "Path": null,
    "Tile": 4
   },
   "Kind": 1
  },
  {
   "Data": {
    "TimeoutSeconds": 0,
    "Categorie": 1,
    "ID": 2,
    "Question": {
     "Title": "",
     "Enonce": []
    }
   },
   "Kind": 10
  },
  {
   "Data": {
    "Player": 0,
    "Success": false,
    "Categorie": 0,
    "AskForMask": false
   },
   "Kind": 7
  },
  {
   "Data": {
    "Player": 1,
    "Success": true,
    "Categorie": 0,
    "AskForMask": false
   },
   "Kind": 7
  },
  {
   "Data": {
    "Player": 2,
    "Success": true,
    "Categorie": 0,
    "AskForMask": false
   },
   "Kind": 7
  },
  {
   "Data": {
    "PlayerName": "",
    "Player": 1
   },
   "Kind": 9
  }
 ],
 "State": {
  "Players": null,
  "PawnTile": 0,
  "Player": 0
 }
}
    """;
      final ev = stateUpdateFromJson(jsonDecode(input));
      expect(ev.events.length, equals(18));
      expect(ev.events[0] is GameStart, equals(true));
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
