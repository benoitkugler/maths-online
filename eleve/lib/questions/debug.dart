// defines helper used when debugging

import 'dart:convert';

import 'package:eleve/types/src_maths_questions_client.dart';
import 'package:eleve/types/src_maths_repere.dart';
import 'package:eleve/types/src_sql_editor.dart';
import 'package:eleve/types/src_tasks.dart';

TextOrMath T(String s) {
  return TextOrMath(s, false);
}

const bounds = RepereBounds(20, 20, Coord(4, 4));

const emptyFigure = Figure(Drawings({}, [], [], [], []), bounds, true, true);

final questionList = [
  const InstantiatedQuestion(
      0, Question([NumberFieldBlock(0, 10)], []), DifficultyTag.diff1, []),
  InstantiatedQuestion(
      0,
      Question([
        const FigureVectorFieldBlock("test", emptyFigure, 0, true),
        const FigureVectorFieldBlock("test", emptyFigure, 0, true),
        TableFieldBlock([
          T("sdsd"),
          T("sdsd"),
        ], [
          T("sdsd"),
          T("sdsd"),
        ], 1)
      ], []),
      DifficultyTag.diffEmpty,
      []),
  const InstantiatedQuestion(
      0, Question([NumberFieldBlock(0, 10)], []), DifficultyTag.diff3, []),
];

final proofB = proofFieldBlockFromJson(jsonDecode("""
{
  "Shape": {
   "Root": {
    "Parts": [
     {
      "Data": {
       "Left": {
        "Data": {
         "Parts": [
          {
           "Data": {
            "Content": []
           },
           "Kind": "Statement"
          },
          {
           "Data": {
            "Content": []
           },
           "Kind": "Statement"
          }
         ]
        },
        "Kind": "Sequence"
       },
       "Right": {
        "Data": {
         "Parts": [
          {
           "Data": {
            "Content": []
           },
           "Kind": "Statement"
          },
          {
           "Data": {
            "Content": []
           },
           "Kind": "Statement"
          }
         ]
        },
        "Kind": "Sequence"
       },
       "Op": 0
      },
      "Kind": "Node"
     },
     {
      "Data": {
       "Terms": [
        null,
        null,
        null
       ]
      },
      "Kind": "Equality"
     },
     {
      "Data": {
       "Terms": [
        null,
        null
       ]
      },
      "Kind": "Equality"
     },
     {
      "Data": {
       "Content": []
      },
      "Kind": "Statement"
     }
    ]
   }
  },
  "TermProposals": [
   [
    {
     "Text": "2k''",
     "IsMath": true
    }
   ],
   [
    {
     "Text": "n est pair",
     "IsMath": false
    }
   ],
   [
    {
     "Text": "m+n",
     "IsMath": true
    }
   ],
   [
    {
     "Text": "m = 2k",
     "IsMath": true
    }
   ],
   [
    {
     "Text": "2(k+k')",
     "IsMath": true
    }
   ],
   [
    {
     "Text": "m+n",
     "IsMath": true
    }
   ],
   [
    {
     "Text": "m+n est pair",
     "IsMath": false
    }
   ],
   [
    {
     "Text": "n = 2k'",
     "IsMath": true
    }
   ],
   [
    {
     "Text": "2k+2k'",
     "IsMath": true
    }
   ],
   [
    {
     "Text": "m est pair",
     "IsMath": false
    }
   ]
  ],
  "ID": 0
 }
 """));
