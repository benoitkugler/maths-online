import 'dart:convert';

import 'package:eleve/types/src_trivial.dart';
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
    "Player": ""
   },
   "Kind": "PlayerJoin"
  },
  {
   "Data": {
    "ID": "",
    "Pseudo": ""
   },
   "Kind": "PlayerReconnected"
  },
  {
   "Data": {
    "PlayerPseudos": {
     "0": "Paul"
    },
    "Pseudo": "",
    "ID": "",
    "IsJoining": false
   },
   "Kind": "LobbyUpdate"
  },
  {
   "Data": {},
   "Kind": "GameStart"
  },
  {
   "Data": {
    "Player": "1"
   },
   "Kind": "PlayerLeft"
  },
  {
   "Data": {
    "PlayerName": "Haha",
    "Player": "2"
   },
   "Kind": "PlayerTurn"
  },
  {
   "Data": {
    "Face": 3
   },
   "Kind": "DiceThrow"
  },
  {
   "Data": {
    "Path": [
     0
    ],
    "Tile": 3
   },
   "Kind": "Move"
  },
  {
   "Data": {
    "PlayerName": "",
    "Tiles": [
     3,
     9
    ],
    "Player": "2"
   },
   "Kind": "PossibleMoves"
  },
  {
   "Data": {
    "TimeoutSeconds": 0,
    "Categorie": 0,
    "ID": 1,
    "Question": {
     "Enonce": [
      {
       "Data": {
        "ID": 0,
        "SizeHint": 0
       },
       "Kind": "NumberFieldBlock"
      }
     ]
    }
   },
   "Kind": "ShowQuestion"
  },
  {
   "Data": {
    "Categorie": 0,
    "Results": {
     "0": {
      "Success": true,
      "AskForMask": false
     },
     "1": {
      "Success": false,
      "AskForMask": false
     },
     "2": {
      "Success": true,
      "AskForMask": false
     }
    }
   },
   "Kind": "PlayerAnswerResults"
  },
  {
   "Data": {
    "QuestionDecrassageIds": {
     "0": [
      1
     ]
    },
    "Winners": [
     "2"
    ],
    "WinnerNames": [
     "Paul"
    ]
   },
   "Kind": "GameEnd"
  },
  {
   "Data": {},
   "Kind": "GameTerminated"
  }
 ],
 "State": {
  "Players": null,
  "PawnTile": 0,
  "PlayerTurn": ""
 }
}
    """;
      final ev = stateUpdateFromJson(jsonDecode(input));
      expect(ev.events.length, equals(13));
      expect(ev.events[0] is PlayerJoin, equals(true));
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
