import 'package:eleve/build_mode.dart';
import 'package:eleve/homework/homework.dart';
import 'package:eleve/homework/types.gen.dart';
import 'package:eleve/main_shared.dart';
import 'package:flutter/material.dart' hide Flow;

void main() async {
  runApp(const _App());
}

const ex = Exercice(
    1, "Pythagore", "Facile !", Parameters([], []), Flow.sequencial, 1, false);
const exH = ExerciceProgressionHeader(
    ex, false, ProgressionExt(Progression(0, 0), [], -1));

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
        // body: Homework(BuildMode.debug, ""),
        appBar: AppBar(),
        body: SheetHome(BuildMode.debug, [
          SheetProgression(
              Sheet(1, 1, "Fonctions affines", Notation.successNotation, true,
                  DateTime.now()),
              [exH, exH]),
          SheetProgression(
              Sheet(1, 1, "Fonctions affines", Notation.successNotation, true,
                  DateTime.now()),
              []),
          SheetProgression(
              Sheet(1, 1, "Fonctions affines", Notation.successNotation, true,
                  DateTime.now()),
              [])
        ]),
      ),
    );
  }
}
