// defines helper used when debugging

import 'dart:convert';
import 'dart:math';

import 'package:eleve/loopback/question.dart';
import 'package:eleve/types/src_maths_functiongrapher.dart';
import 'package:eleve/types/src_maths_questions_client.dart';
import 'package:eleve/types/src_maths_repere.dart';
import 'package:eleve/types/src_sql_editor.dart';
import 'package:eleve/types/src_tasks.dart';
import 'package:flutter/material.dart';

TextOrMath T(String s) {
  return TextOrMath(s, false);
}

const bounds = RepereBounds(20, 20, Coord(4, 4));

const emptyFigure = Figure(Drawings({}, [], [], [], []), bounds, true, true);

String toHex(int v) {
  return v.toRadixString(16).padLeft(2, '0');
}

FunctionGraph function({String? label, Color? color}) {
  color = color ?? randColor();
  return FunctionGraph(
      FunctionDecoration(label ?? "C_f",
          "${toHex(color.red)}${toHex(color.green)}${toHex(color.blue)}"),
      List.generate(100, (index) {
        final x1 = 0 + (index / 100) * 15;
        final x2 = x1 + (1 / 100) * 15;
        final y = x1;
        return BezierCurve(
            Coord(x1, y), Coord(x1, y + Random().nextDouble()), Coord(x2, y));
      }));
}

final questionList = [
  const InstantiatedQuestion(
      0, Question([NumberFieldBlock(0, 10)], []), DifficultyTag.diff1, []),
  InstantiatedQuestion(
      0,
      Question([
        const GeometricConstructionFieldBlock(
            0, GFVector("test", true), FigureBlock(emptyFigure)),
        const GeometricConstructionFieldBlock(
            0, GFVector("test", false), FigureBlock(emptyFigure)),
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
