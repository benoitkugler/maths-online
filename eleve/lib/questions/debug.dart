// defines helper used when debugging

import 'package:eleve/questions/repere.gen.dart';
import 'package:eleve/questions/types.gen.dart';

TextOrMath T(String s) {
  return TextOrMath(s, false);
}

const bounds = RepereBounds(20, 20, Coord(4, 4));

const emptyFigure = Figure(Drawings({}, [], []), bounds, true);
