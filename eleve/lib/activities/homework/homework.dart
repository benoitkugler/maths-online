import 'dart:convert';

import 'package:eleve/activities/homework/progression_bar.dart';
import 'package:eleve/activities/homework/sheet.dart';
import 'package:eleve/activities/homework/types.gen.dart';
import 'package:eleve/activities/homework/utils.dart';
import 'package:eleve/build_mode.dart';
import 'package:eleve/exercice/home.dart';
import 'package:eleve/questions/fields.dart';
import 'package:eleve/settings.dart';
import 'package:eleve/shared/errors.dart';
import 'package:eleve/shared_gen.dart';
import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;

typedef Sheets = List<SheetProgression>;

abstract class HomeworkAPI extends FieldAPI {
  Future<Sheets> loadSheets();
  Future<InstantiatedWork> loadWork(IdTask id);
  Future<StudentEvaluateTaskOut> evaluateExercice(
      IdTask idTask, EvaluateWorkIn ex);
}

class ServerHomeworkAPI implements HomeworkAPI {
  final BuildMode buildMode;
  final String studentID;

  const ServerHomeworkAPI(this.buildMode, this.studentID);

  @override
  Future<CheckExpressionOut> checkExpressionSyntax(String expression) {
    return ServerFieldAPI(buildMode).checkExpressionSyntax(expression);
  }

  @override
  Future<Sheets> loadSheets() async {
    const serverEndpoint = "/api/student/homework/sheets";
    final uri = Uri.parse(
        buildMode.serverURL(serverEndpoint, query: {studentIDKey: studentID}));
    final resp = await http.get(uri);
    return listSheetProgressionFromJson(jsonDecode(resp.body));
  }

  @override
  Future<InstantiatedWork> loadWork(IdTask idTask) async {
    const serverEndpoint = "/api/student/homework/task/instantiate";
    final uri = Uri.parse(buildMode.serverURL(serverEndpoint,
        query: {studentIDKey: studentID, "id": idTask.toString()}));
    final resp = await http.get(uri);
    return instantiatedWorkFromJson(jsonDecode(resp.body));
  }

  @override
  Future<StudentEvaluateTaskOut> evaluateExercice(
      IdTask idTask, EvaluateWorkIn ex) async {
    const serverEndpoint = "/api/student/homework/task/evaluate";
    final uri = Uri.parse(buildMode.serverURL(serverEndpoint));
    final resp = await http.post(uri,
        body: jsonEncode(studentEvaluateTaskInToJson(
            StudentEvaluateTaskIn(studentID, idTask, ex))),
        headers: {
          'Content-type': 'application/json',
        });
    return studentEvaluateTaskOutFromJson(checkServerError(resp.body));
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
                child: Padding(
              padding: const EdgeInsets.all(8.0),
              child: ErrorBar(
                  "Impossible de charger le travail à faire.", snapshot.error!),
            ));
          } else if (snapshot.hasData) {
            return _SheetList(widget.api, snapshot.data!);
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
      child:
          Column(mainAxisAlignment: MainAxisAlignment.center, children: const [
        CircularProgressIndicator(value: null),
        SizedBox(height: 20),
        Text("Chargement des fiches de travail..."),
      ]),
    );
  }
}

class _SheetList extends StatefulWidget {
  final HomeworkAPI api;
  final Sheets initialSheets;
  const _SheetList(this.api, this.initialSheets, {Key? key}) : super(key: key);

  @override
  State<_SheetList> createState() => _SheetListState();
}

class _SheetListState extends State<_SheetList> {
  late final Sheets sheets;

  @override
  void initState() {
    sheets = widget.initialSheets;
    super.initState();
  }

  bool updateMark(SheetMarkNotification notif) {
    final index =
        sheets.indexWhere((element) => element.sheet.id == notif.idSheet);
    final actual = sheets[index];
    notif.applyTo(actual.tasks);
    setState(() {
      sheets[index] = SheetProgression(actual.sheet, actual.tasks);
    });
    return true;
  }

  void onSelectSheet(SheetProgression sheet) {
    Navigator.of(context).push(MaterialPageRoute<void>(builder: (context) {
      return NotificationListener<SheetMarkNotification>(
          child: SheetW(widget.api, sheet), onNotification: updateMark);
    }));
  }

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
                onTap: () => onSelectSheet(e),
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
    final ma = sheetMark(sheet.tasks);
    final hasNotation = sheet.sheet.notation != Notation.noNotation;
    final started = sheet.tasks.where((ex) => ex.hasProgression).length;
    final completed = sheet.tasks
        .where((ex) => ex.hasProgression && ex.progression.isCompleted())
        .length;
    return Card(
      elevation: emphasize ? 3 : null,
      shadowColor: emphasize ? Colors.blueAccent : null,
      color: emphasize ? Colors.blueAccent.withOpacity(0.3) : null,
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
            if (sheet.tasks.isNotEmpty)
              Row(mainAxisSize: MainAxisSize.min, children: [
                Expanded(
                  child: Padding(
                    padding: const EdgeInsets.symmetric(vertical: 8.0),
                    child: ProgressionBar(
                        total: sheet.tasks.length,
                        completed: completed,
                        started: started),
                  ),
                ),
                if (hasNotation)
                  Padding(
                    padding: const EdgeInsets.only(left: 8.0),
                    child: Text("Note : ${ma.mark} / ${ma.bareme}"),
                  ),
              ])
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