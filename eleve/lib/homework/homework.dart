import 'dart:convert';

import 'package:eleve/build_mode.dart';
import 'package:eleve/exercice/home.dart';
import 'package:eleve/homework/progression.dart';
import 'package:eleve/homework/sheet.dart';
import 'package:eleve/homework/types.gen.dart';
import 'package:eleve/homework/utils.dart';
import 'package:eleve/settings.dart';
import 'package:eleve/shared/errors.dart';
import 'package:eleve/shared_gen.dart';
import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;

typedef Sheets = List<SheetProgression>;

abstract class HomeworkAPI {
  Future<Sheets> loadSheets();
  Future<StudentExerciceInst> loadExercice(Exercice ex);
}

class ServerHomeworkAPI implements HomeworkAPI {
  static const serverEndpoint = "/api/student/homework/sheets";

  final BuildMode buildMode;
  final String studentID;

  const ServerHomeworkAPI(this.buildMode, this.studentID);

  @override
  Future<Sheets> loadSheets() async {
    final uri = Uri.parse(
        buildMode.serverURL(serverEndpoint, query: {studentIDKey: studentID}));
    final resp = await http.get(uri);
    return listSheetProgressionFromJson(jsonDecode(resp.body));
  }

  @override
  Future<StudentExerciceInst> loadExercice(Exercice ex) {
    // TODO: implement loadExercice
    throw UnimplementedError();
  }
}

/// Homework is the entry point for the homework
/// activity
class Homework extends StatefulWidget {
  final HomeworkAPI api;

  const Homework(this.api, {Key? key}) : super(key: key);

  @override
  State<Homework> createState() => _HomeworkState();
}

class _HomeworkState extends State<Homework> {
  late Future<Sheets> sheets;

  @override
  void initState() {
    sheets = widget.api.loadSheets();

    super.initState();
  }

  @override
  void didUpdateWidget(covariant Homework oldWidget) {
    sheets = widget.api.loadSheets();
    super.didUpdateWidget(oldWidget);
  }

  void onSelectSheet(SheetProgression sheet) {
    Navigator.of(context).push(MaterialPageRoute<void>(builder: (context) {
      return SheetW(widget.api, sheet);
    }));
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text("Travail à la maison"),
      ),
      body: FutureBuilder<Sheets>(
        future: sheets,
        builder: (context, snapshot) {
          if (snapshot.hasError) {
            return Center(
                child: ErrorBar(
                    "Impossible de charger les données.", snapshot.error!));
          } else if (snapshot.hasData) {
            return _SheetList(snapshot.data!, onSelectSheet);
          } else {
            return const _Loading();
          }
        },
      ),
    );
  }
}

class _Loading extends StatelessWidget {
  const _Loading({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Center(
      child: Column(children: const [
        CircularProgressIndicator(value: null),
        SizedBox(height: 20),
        Text("Chargement des fiches de travail..."),
      ]),
    );
  }
}

// shows the most interesting sheet, or nothing

class _SheetList extends StatelessWidget {
  final Sheets sheets;
  final void Function(SheetProgression) onSelect;
  const _SheetList(this.sheets, this.onSelect, {Key? key}) : super(key: key);

  // assume a one at a time sheet to do and emphasize it
  int selectMainSheetID() {
    // select the most recent one
    final shs = sheets.map((e) => e).toList();
    shs.sort(((a, b) => a.sheet.deadline.isAfter(b.sheet.deadline) ? 1 : -1));
    return sheets[0].sheet.id;
  }

  @override
  Widget build(BuildContext context) {
    if (sheets.isEmpty) {
      return const Center(
          child: Text("Aucun travail à la maison n'est planifié."));
    }
    final bestSheet = selectMainSheetID();
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 6.0, horizontal: 2),
      child: ListView(
        children: sheets
            .map((e) => InkWell(
                onTap: () => onSelect(e),
                child: _SheetSummary(
                  e,
                  emphasize: e.sheet.id == bestSheet,
                )))
            .toList(),
      ),
    );
  }
}

class _SheetSummary extends StatelessWidget {
  final SheetProgression sheet;
  final bool emphasize;
  const _SheetSummary(this.sheet, {this.emphasize = false, Key? key})
      : super(key: key);

  @override
  Widget build(BuildContext context) {
    final hasNotation = sheet.sheet.notation != Notation.noNotation;
    final started = sheet.exercices.where((ex) => ex.hasProgression).length;
    final completed =
        sheet.exercices.where((ex) => ex.progression.isCompleted()).length;
    return Card(
      color: emphasize ? Colors.blueAccent : null,
      child: Padding(
        padding: const EdgeInsets.all(12.0),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              sheet.sheet.title,
              style: const TextStyle(fontSize: 18),
            ),
            Padding(
              padding: const EdgeInsets.symmetric(vertical: 8.0),
              child: Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  Text(hasNotation ? "Travail noté" : "Travail non noté"),
                  if (hasNotation)
                    Text("Pour le ${formatTime(sheet.sheet.deadline)}")
                ],
              ),
            ),
            if (sheet.exercices.isNotEmpty)
              Padding(
                padding: const EdgeInsets.symmetric(vertical: 8.0),
                child: ProgressionBar(
                    total: sheet.exercices.length,
                    completed: completed,
                    started: started),
              )
          ],
        ),
      ),
    );
  }
}

class HomeworkActivityIcon extends StatelessWidget {
  final void Function() onTap;

  const HomeworkActivityIcon(this.onTap, {Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Column(
      children: [
        GestureDetector(
          onTap: onTap,
          child: Image.asset("lib/images/homework.png", width: 68),
        ),
        const Padding(
          padding: EdgeInsets.only(top: 8, bottom: 6),
          child: Text("Travail à la maison"),
        ),
      ],
    );
  }
}
