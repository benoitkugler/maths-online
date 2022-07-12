// defines helper used when debugging

import 'package:eleve/questions/repere.gen.dart';
import 'package:eleve/questions/types.gen.dart';
import 'package:eleve/shared_gen.dart';

TextOrMath T(String s) {
  return TextOrMath(s, false);
}

const bounds = RepereBounds(20, 20, Coord(4, 4));

const emptyFigure = Figure(Drawings({}, [], []), bounds, true, true);

final questionList = [
  const InstantiatedQuestion(0, Question("", [NumberFieldBlock(0)]), []),
  InstantiatedQuestion(
      0,
      Question("", [
        const FigureVectorFieldBlock("test", emptyFigure, 0, true),
        const FigureVectorFieldBlock("test", emptyFigure, 0, true),
        TableFieldBlock([
          T("sdsd"),
          T("sdsd"),
        ], [
          T("sdsd"),
          T("sdsd"),
        ], 1)
      ]),
      []),
  const InstantiatedQuestion(0, Question("", [NumberFieldBlock(0)]), []),
];
