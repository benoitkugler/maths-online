import 'dart:convert';

import 'package:eleve/activities/homework/progression_bar.dart';
import 'package:eleve/activities/homework/sheet.dart';
import 'package:eleve/activities/homework/utils.dart';
import 'package:eleve/build_mode.dart';
import 'package:eleve/exercice/home.dart';
import 'package:eleve/questions/fields.dart';
import 'package:eleve/settings.dart';
import 'package:eleve/shared/activity_start.dart';
import 'package:eleve/shared/errors.dart';
import 'package:eleve/types/src.dart';
import 'package:eleve/types/src_prof_homework.dart';
import 'package:eleve/types/src_sql_homework.dart';
import 'package:eleve/types/src_sql_tasks.dart';
import 'package:eleve/types/src_tasks.dart';
import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;

typedef Sheets = List<SheetProgression>;

extension on SheetProgression {
  Set<String> chapters() => tasks.map((e) => e.chapter).toSet();
}

abstract class HomeworkAPI extends FieldAPI {
  Future<Sheets> loadSheets(bool loadNonNoted);
  Future<InstantiatedWork> loadWork(IdTask id);
  Future<StudentEvaluateTaskOut> evaluateExercice(
      IdTask idTask, IdTravail idTravail, EvaluateWorkIn ex);
  Future<void> resetTask(IdTravail idTravail, IdTask idTask);
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
  Future<Sheets> loadSheets(loadNonNoted) async {
    final serverEndpoint = loadNonNoted
        ? "/api/student/homework/sheets/free"
        : "/api/student/homework/sheets";
    final uri = Uri.parse(
        buildMode.serverURL(serverEndpoint, query: {studentIDKey: studentID}));
    final resp = await http.get(uri);
    return listSheetProgressionFromJson(checkServerError(resp.body));
  }

  @override
  Future<InstantiatedWork> loadWork(IdTask idTask) async {
    const serverEndpoint = "/api/student/homework/task/instantiate";
    final uri = Uri.parse(buildMode.serverURL(serverEndpoint,
        query: {studentIDKey: studentID, "id": idTask.toString()}));
    final resp = await http.get(uri);
    return instantiatedWorkFromJson(checkServerError(resp.body));
  }

  @override
  Future<StudentEvaluateTaskOut> evaluateExercice(
      IdTask idTask, IdTravail idTravail, EvaluateWorkIn ex) async {
    const serverEndpoint = "/api/student/homework/task/evaluate";
    final uri = Uri.parse(buildMode.serverURL(serverEndpoint));
    final resp = await http.post(uri,
        body: jsonEncode(studentEvaluateTaskInToJson(
            StudentEvaluateTaskIn(studentID, idTask, ex, idTravail))),
        headers: {
          'Content-type': 'application/json',
        });
    return studentEvaluateTaskOutFromJson(checkServerError(resp.body));
  }

  @override
  Future<void> resetTask(IdTravail idTravail, IdTask idTask) async {
    const serverEndpoint = "/api/student/homework/task/reset";
    final uri = Uri.parse(buildMode.serverURL(serverEndpoint));
    final resp = await http.post(uri,
        body: jsonEncode(studentResetTaskInToJson(
            StudentResetTaskIn(studentID, idTravail, idTask))),
        headers: {
          'Content-type': 'application/json',
        });
    checkServerError(resp.body);
  }
}

class HomeworkDisabled extends StatelessWidget {
  const HomeworkDisabled({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
        appBar: AppBar(
          title: const Text("Travail à la maison"),
        ),
        body: const Center(
            child: Card(
          margin: EdgeInsets.all(20),
          child: Padding(
            padding: EdgeInsets.all(8.0),
            child: Text(
              "Cette activité n'est pas disponible, car tu n'es pas inscris sur une classe.",
              style: TextStyle(fontSize: 16),
            ),
          ),
        )));
  }
}

/// [HomeworkStart] is the entry point widget for the homework
/// activity, and diplays a home screen allowing the
/// user to choose between noted, limited vs free exercices.
///
/// See [HomeworkDisabled] is the student is not registred.
class HomeworkStart extends StatelessWidget {
  final HomeworkAPI api;

  const HomeworkStart(this.api, {super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text("Démarrer un exercice"),
      ),
      body: Column(mainAxisAlignment: MainAxisAlignment.spaceEvenly, children: [
        LaunchCard("Devoirs", "Je veux faire le travail prévu pour la classe.",
            const Icon(Icons.pending_actions), () {
          Navigator.of(context).push(MaterialPageRoute<void>(
              builder: (_) => Scaffold(body: _HomeworkW(api, false))));
        }),
        const Divider(thickness: 4),
        LaunchCard("Entrainement", "Je veux réviser ou m'entrainer librement.",
            const Icon(Icons.person), () {
          Navigator.of(context).push(MaterialPageRoute<void>(
              builder: (_) => Scaffold(body: _HomeworkW(api, true))));
        }),
      ]),
    );
  }
}

/// [_HomeworkW] is the entry point widget for the homework
/// activity.
class _HomeworkW extends StatefulWidget {
  final HomeworkAPI api;
  final bool isNonNoted;

  /// Creates a new [_HomeworkW] widget
  const _HomeworkW(this.api, this.isNonNoted, {Key? key}) : super(key: key);

  @override
  State<_HomeworkW> createState() => _HomeworkWState();
}

class _HomeworkWState extends State<_HomeworkW> {
  late Future<Sheets> sheets;

  @override
  void initState() {
    sheets = widget.api.loadSheets(widget.isNonNoted);
    super.initState();
  }

  @override
  void didUpdateWidget(covariant _HomeworkW oldWidget) {
    sheets = widget.api.loadSheets(widget.isNonNoted);
    super.didUpdateWidget(oldWidget);
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text(widget.isNonNoted
            ? "Travaux d'entrainement"
            : "Devoirs à la maison"),
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
            return _SheetList(widget.api, widget.isNonNoted, snapshot.data!);
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
    return const Center(
      child: Column(mainAxisAlignment: MainAxisAlignment.center, children: [
        CircularProgressIndicator(value: null),
        SizedBox(height: 20),
        Text("Chargement des fiches de travail..."),
      ]),
    );
  }
}

class _SheetList extends StatefulWidget {
  final HomeworkAPI api;
  final bool isNotNoted;
  final Sheets initialSheets;

  const _SheetList(this.api, this.isNotNoted, this.initialSheets, {Key? key})
      : super(key: key);

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
      sheets[index] =
          SheetProgression(actual.idTravail, actual.sheet, actual.tasks);
    });
    return true;
  }

  void onSelectSheet(SheetProgression sheet) {
    Navigator.of(context).push(MaterialPageRoute<void>(builder: (context) {
      return NotificationListener<SheetMarkNotification>(
          onNotification: updateMark, child: SheetW(widget.api, sheet));
    }));
  }

  @override
  Widget build(BuildContext context) {
    if (sheets.isEmpty) {
      return Center(
          child: Text(widget.isNotNoted
              ? "Aucun exercice n'est disponible en accès libre."
              : "Aucun travail à la maison n'est planifié."));
    }
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 6.0, horizontal: 4),
      child: widget.isNotNoted
          ? _FreeSheetList(sheets, onSelectSheet)
          : _NotedSheetList(sheets, onSelectSheet),
    );
  }
}

class _NotedSheetList extends StatelessWidget {
  final Sheets sheets;
  final void Function(SheetProgression) onTap;

  const _NotedSheetList(this.sheets, this.onTap, {super.key});

  // assume a one at a time sheet to do and emphasize it
  int? _selectMainSheetID() {
    // select the most recent one
    final now = DateTime.now();
    final shs = sheets
        .where((e) => e.sheet.noted && !e.sheet.deadline.isBefore(now))
        .toList();
    if (shs.isEmpty) return null;
    shs.sort(((a, b) => a.sheet.deadline.isAfter(b.sheet.deadline) ? 1 : -1));
    return shs[0].sheet.id;
  }

  @override
  Widget build(BuildContext context) {
    final bestSheet = _selectMainSheetID();
    final now = DateTime.now();
    return ListView(
        children: sheets
            .map((e) => InkWell(
                onTap: () => onTap(e),
                child: _SheetSummary(
                  e,
                  status: e.sheet.deadline.isBefore(now)
                      ? SheetStatus.expired
                      : (e.sheet.id == bestSheet
                          ? SheetStatus.suggested
                          : SheetStatus.normal),
                )))
            .toList());
  }
}

class _FreeSheetList extends StatefulWidget {
  final Sheets sheets;
  final void Function(SheetProgression) onTap;

  const _FreeSheetList(this.sheets, this.onTap, {super.key});

  List<MapEntry<String, List<SheetProgression>>> get _groups {
    final tmp = <String, List<SheetProgression>>{};
    for (var sheet in sheets) {
      final thisChapters = sheet.chapters();
      // if the sheet is empty, it has no chapters
      // still show it
      if (thisChapters.isEmpty) thisChapters.add("");

      for (var chapter in thisChapters) {
        final l = tmp.putIfAbsent(chapter, () => []);
        l.add(sheet);
      }
    }
    final out = tmp.entries.toList();
    out.sort((a, b) => a.key.compareTo(b.key));
    return out;
  }

  @override
  State<_FreeSheetList> createState() => __FreeSheetListState();
}

class __FreeSheetListState extends State<_FreeSheetList> {
  List<bool> _expanded = [];

  @override
  void initState() {
    _initExpanded();
    super.initState();
  }

  @override
  void didUpdateWidget(covariant _FreeSheetList oldWidget) {
    _initExpanded();
    super.didUpdateWidget(oldWidget);
  }

  void _initExpanded() {
    _expanded = List.filled(widget._groups.length, false);
    if (_expanded.isNotEmpty) _expanded[0] = true;
  }

  void _expand(int index) {
    final isExpanded = _expanded[index];
    setState(() {
      // reset other panels
      for (var i = 0; i < _expanded.length; i++) {
        _expanded[i] = false;
      }
      _expanded[index] = !isExpanded;
    });
  }

  @override
  Widget build(BuildContext context) {
    return SingleChildScrollView(
      child: ExpansionPanelList(
          expansionCallback: (panelIndex, isExpanded) => _expand(panelIndex),
          dividerColor: Colors.lightBlue.withOpacity(0.8),
          expandedHeaderPadding: const EdgeInsets.all(6),
          children: List.generate(widget._groups.length, (index) {
            final e = widget._groups[index];
            return ExpansionPanel(
                backgroundColor: Colors.transparent,
                canTapOnHeader: true,
                isExpanded: _expanded[index],
                headerBuilder: (context, isExpanded) => Align(
                    alignment: Alignment.centerLeft,
                    child: Padding(
                        padding: const EdgeInsets.symmetric(horizontal: 8.0),
                        child: Row(
                          children: [
                            Flexible(
                              child: Text(e.key.isEmpty ? "Non classé" : e.key,
                                  style:
                                      Theme.of(context).textTheme.titleMedium),
                            ),
                            Chip(
                              label: Text("${e.value.length}"),
                              visualDensity: const VisualDensity(vertical: -2),
                            ),
                          ],
                        ))),
                body: Column(
                  mainAxisSize: MainAxisSize.min,
                  children: e.value
                      .map((e) => InkWell(
                            onTap: () => widget.onTap(e),
                            child: _SheetSummary(
                              e,
                              status: SheetStatus.normal,
                            ),
                          ))
                      .toList(),
                ));
          })),
    );
  }
}

enum SheetStatus { normal, suggested, expired }

const sheetSuggestedColor = Colors.blueAccent;
final sheetExpiredColor = Colors.red.shade300;

class _SheetSummary extends StatelessWidget {
  final SheetProgression sheet;
  final SheetStatus status;
  const _SheetSummary(this.sheet, {this.status = SheetStatus.normal, Key? key})
      : super(key: key);

  @override
  Widget build(BuildContext context) {
    final ma = sheetMark(sheet.tasks);
    final hasNotation = sheet.sheet.noted;
    final started = sheet.tasks.where((ex) => ex.hasProgression).length;
    final completed = sheet.tasks
        .where((ex) => ex.hasProgression && ex.progression.isCompleted())
        .length;
    final highlight = status == SheetStatus.suggested;
    final isExpired = status == SheetStatus.expired;
    return Card(
      shape: highlight
          ? RoundedRectangleBorder(
              side: const BorderSide(color: sheetSuggestedColor, width: 2.0),
              borderRadius: BorderRadius.circular(4.0))
          : null,
      elevation: highlight ? 3 : null,
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
              child: hasNotation
                  ? Row(
                      mainAxisAlignment: MainAxisAlignment.spaceBetween,
                      children: [
                          const Text("Travail noté"),
                          DeadlineCard(
                              isExpired: isExpired,
                              deadline: sheet.sheet.deadline)
                        ])
                  : const SizedBox(),
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
                Padding(
                  padding: const EdgeInsets.only(left: 8.0),
                  child: Text(
                      "${hasNotation ? 'Note' : 'Score'} : ${ma.mark} / ${ma.bareme}"),
                ),
              ])
          ],
        ),
      ),
    );
  }
}

class DeadlineCard extends StatelessWidget {
  const DeadlineCard({
    super.key,
    required this.isExpired,
    required this.deadline,
  });

  final bool isExpired;
  final DateTime deadline;

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.all(4),
      decoration: BoxDecoration(
          color: isExpired ? sheetExpiredColor : sheetSuggestedColor,
          borderRadius: const BorderRadius.all(Radius.circular(4))),
      child: RichText(
          text: TextSpan(children: [
        const TextSpan(text: "A rendre avant le\n"),
        TextSpan(
            text: formatTime(deadline),
            style: const TextStyle(fontWeight: FontWeight.bold)),
      ])),
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
        RawMaterialButton(
          onPressed: onTap,
          child: Image.asset("assets/images/homework.png", width: 68),
        ),
        const Padding(
          padding: EdgeInsets.only(top: 8, bottom: 6),
          child: Text("Travail à la maison"),
        ),
      ],
    );
  }
}
