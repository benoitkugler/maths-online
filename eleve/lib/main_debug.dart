import 'package:eleve/main_shared.dart';
import 'package:eleve/questions/repere.dart';
import 'package:eleve/questions/repere.gen.dart';
import 'package:flutter/material.dart';

void main() async {
  runApp(const _App());
}

class _App extends StatelessWidget {
  const _App({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Isyro',
      theme: theme,
      debugShowCheckedModeBanner: false,
      localizationsDelegates: localizations,
      supportedLocales: locales,
      home: Scaffold(
        body: Center(
          child: const StaticRepere(Figure(
              Drawings({
                "A": LabeledPoint("", PosPoint(Coord(1, 1), LabelPos.bottom)),
                "B": LabeledPoint("", PosPoint(Coord(1, 4), LabelPos.bottom)),
                "C": LabeledPoint("", PosPoint(Coord(4, 4), LabelPos.bottom)),
                "D": LabeledPoint("", PosPoint(Coord(4, 10), LabelPos.bottom)),
              }, [], [], [
                Area("#FF00FF", ["A", "B", "D", "C"]),
                Area("#FF0000", ["A", "B", "C"]),
              ]),
              RepereBounds(20, 20, Coord(1, 1)),
              true,
              true)),
        ),
      ),
    );
  }
}
