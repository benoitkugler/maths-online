import 'dart:convert';

import 'package:eleve/activities/homework/sheet.dart';
import 'package:eleve/activities/homework/utils.dart';
import 'package:eleve/build_mode.dart';
import 'package:eleve/exercice/exercice.dart';
import 'package:eleve/shared/activity_start.dart';
import 'package:eleve/shared/animated_logo.dart';
import 'package:eleve/shared/errors.dart';
import 'package:eleve/shared/progression_bar.dart';
import 'package:eleve/shared/settings_shared.dart';
import 'package:eleve/types/src_prof_homework.dart';
import 'package:eleve/types/src_sql_homework.dart';
import 'package:eleve/types/src_sql_tasks.dart';
import 'package:eleve/types/src_sql_teacher.dart';
import 'package:eleve/types/src_tasks.dart';
import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;

typedef Sheets = List<SheetProgression>;

extension on SheetProgression {
  Set<String> chapters() => tasks.map((e) => e.chapter).toSet();
}

abstract class HomeworkAPI {
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
  Future<Sheets> loadSheets(loadNonNoted) async {
    final serverEndpoint = loadNonNoted
        ? "/api/student/homework/sheets/free"
        : "/api/student/homework/sheets";
    final uri =
        buildMode.serverURL(serverEndpoint, query: {studentIDKey: studentID});
    final resp = await http.get(uri);
    return listSheetProgressionFromJson(checkServerError(resp.body));
  }

  @override
  Future<InstantiatedWork> loadWork(IdTask idTask) async {
    const serverEndpoint = "/api/student/homework/task/instantiate";
    final uri = buildMode.serverURL(serverEndpoint,
        query: {studentIDKey: studentID, "id": idTask.toString()});
    final resp = await http.get(uri);
    return instantiatedWorkFromJson(checkServerError(resp.body));
  }

  @override
  Future<StudentEvaluateTaskOut> evaluateExercice(
      IdTask idTask, IdTravail idTravail, EvaluateWorkIn ex) async {
    const serverEndpoint = "/api/student/homework/task/evaluate";
    final uri = buildMode.serverURL(serverEndpoint);
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
    final uri = buildMode.serverURL(serverEndpoint);
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
              "Cette activité n'est pas disponible car tu n'es pas inscrit dans une classe.",
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
  final bool isNotNoted;

  /// Creates a new [_HomeworkW] widget
  const _HomeworkW(this.api, this.isNotNoted);

  @override
  State<_HomeworkW> createState() => _HomeworkWState();
}

class _HomeworkWState extends State<_HomeworkW> {
  late Future<Sheets> sheets;

  @override
  void initState() {
    sheets = widget.api.loadSheets(widget.isNotNoted);
    super.initState();
  }

  @override
  void didUpdateWidget(covariant _HomeworkW oldWidget) {
    sheets = widget.api.loadSheets(widget.isNotNoted);
    super.didUpdateWidget(oldWidget);
  }

  @override
  Widget build(BuildContext context) {
    return FutureBuilder<Sheets>(
        future: sheets,
        builder: (context, snapshot) {
          if (!snapshot.hasError && snapshot.hasData) {
            return _SheetListView(
                widget.api, widget.isNotNoted, snapshot.data!);
          }
          return Scaffold(
            appBar: AppBar(
                title: Text(widget.isNotNoted
                    ? "Travaux d'entrainement"
                    : "Devoirs à la maison")),
            body: snapshot.hasError
                ? Center(
                    child: Padding(
                    padding: const EdgeInsets.all(8.0),
                    child: ErrorBar("Impossible de charger le travail à faire.",
                        snapshot.error!),
                  ))
                : const _Loading(),
          );
        });
  }
}

class _Loading extends StatelessWidget {
  const _Loading();

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

class _SheetListView extends StatefulWidget {
  final HomeworkAPI api;
  final bool isNotNoted;
  final Sheets initialSheets;

  const _SheetListView(this.api, this.isNotNoted, this.initialSheets);

  @override
  State<_SheetListView> createState() => _SheetListViewState();
}

class _SheetListViewState extends State<_SheetListView> {
  late final Sheets sheets;

  @override
  void initState() {
    sheets = widget.initialSheets;
    super.initState();
  }

  double get averageMark {
    double total = 0;
    int nbSheets = 0;
    for (var sheet in sheets) {
      if (sheet.sheet.ignoreForMark) continue; // the student has a dispense
      final ma = sheetMark(sheet.tasks);
      nbSheets += 1;
      total += 20.0 * ma.mark.toDouble() / ma.bareme; // convert to 20
    }
    if (nbSheets == 0) return 0;
    return total / nbSheets;
  }

  void _showNotedSheetAverage() {
    showDialog<void>(
        context: context,
        builder: (context) => AlertDialog(
            backgroundColor: Theme.of(context).scaffoldBackgroundColor,
            elevation: 8,
            shadowColor: Colors.white,
            title: const Text("Moyenne"),
            content: Column(
              mainAxisSize: MainAxisSize.min,
              children: [
                SizedBox(
                    width: 80,
                    height: 80,
                    child: AnimatedLogo(averageMark / 20)),
                const SizedBox(height: 20),
                Text("${averageMark.toStringAsFixed(1)} / 20",
                    style: Theme.of(context).textTheme.headlineSmall)
              ],
            )));
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

  void onSelectSheet(SheetProgression sheet) async {
    // show activity start only for math
    if (sheet.sheet.matiere == MatiereTag.mathematiques) {
      final onDone = await Navigator.of(context).push(MaterialPageRoute<bool>(
          builder: (context) =>
              MathActivityStart(() => Navigator.of(context).pop(true))));
      if (onDone == null) return;
      if (!mounted) return;
    }

    Navigator.of(context).push(MaterialPageRoute<void>(builder: (context) {
      return NotificationListener<SheetMarkNotification>(
          onNotification: updateMark, child: SheetW(widget.api, sheet));
    }));
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
        appBar: AppBar(
            title: Text(widget.isNotNoted
                ? "Travaux d'entrainement"
                : "Devoirs à la maison"),
            actions: widget.isNotNoted
                ? []
                : [
                    IconButton(
                      onPressed: _showNotedSheetAverage,
                      icon: const Icon(Icons.pie_chart),
                    )
                  ]),
        body: sheets.isEmpty
            ? Center(
                child: Text(widget.isNotNoted
                    ? "Aucun exercice n'est disponible en accès libre."
                    : "Aucun travail à la maison n'est planifié."))
            : Padding(
                padding:
                    const EdgeInsets.symmetric(vertical: 6.0, horizontal: 4),
                child: widget.isNotNoted
                    ? _FreeSheetList(sheets, onSelectSheet)
                    : _NotedSheetList(sheets, onSelectSheet),
              ));
  }
}

class _NotedSheetList extends StatelessWidget {
  final Sheets sheets;
  final void Function(SheetProgression) onTap;

  const _NotedSheetList(this.sheets, this.onTap);

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

  const _FreeSheetList(this.sheets, this.onTap);

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

  static double progressionRatio(List<SheetProgression> group) {
    int totalMark = 0;
    int totalBareme = 0;
    for (var sheet in group) {
      final ma = sheetMark(sheet.tasks);
      totalMark += ma.mark;
      totalBareme += ma.bareme;
    }
    return totalMark.toDouble() / totalBareme;
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
          dividerColor: Colors.lightBlue.withValues(alpha: 0.8),
          expandedHeaderPadding: const EdgeInsets.all(6),
          children: List.generate(widget._groups.length, (index) {
            final group = widget._groups[index];
            return ExpansionPanel(
                backgroundColor: Colors.transparent,
                canTapOnHeader: true,
                isExpanded: _expanded[index],
                headerBuilder: (context, isExpanded) => Align(
                    alignment: Alignment.centerLeft,
                    child: Padding(
                        padding: const EdgeInsets.symmetric(horizontal: 8.0),
                        child: Row(
                          mainAxisAlignment: MainAxisAlignment.spaceBetween,
                          children: [
                            Flexible(
                              child: Text(
                                  group.key.isEmpty ? "Non classé" : group.key,
                                  style:
                                      Theme.of(context).textTheme.titleMedium),
                            ),
                            const SizedBox(width: 4),
                            isExpanded
                                ? SizedBox(
                                    width: 40,
                                    height: 40,
                                    child: AnimatedLogo(
                                        _FreeSheetList.progressionRatio(
                                            group.value)))
                                : Chip(
                                    label: Text("${group.value.length}"),
                                    visualDensity:
                                        const VisualDensity(vertical: -2),
                                  ),
                          ],
                        ))),
                body: Column(
                  mainAxisSize: MainAxisSize.min,
                  children: group.value
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
  const _SheetSummary(this.sheet, {this.status = SheetStatus.normal});

  @override
  Widget build(BuildContext context) {
    final ma = sheetMark(sheet.tasks);

    final hasNotation = sheet.sheet.noted;
    final total = sheet.tasks.length;
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
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                Text(
                  sheet.sheet.title,
                  style: const TextStyle(fontSize: 18),
                ),
                Text(
                  matiereTagLabel(sheet.sheet.matiere),
                  style: const TextStyle(fontSize: 10),
                )
              ],
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
                    child: ProgressionBar(background: Colors.grey, layers: [
                      ProgressionLayer(started.toDouble() / total,
                          Colors.yellow.shade200, true),
                      ProgressionLayer(completed.toDouble() / total,
                          Colors.lightGreenAccent, false),
                    ]),
                  ),
                ),
                Padding(
                  padding: const EdgeInsets.only(left: 8.0),
                  child: sheet.sheet.ignoreForMark
                      ? const Text("Note ignorée")
                      : Text(
                          "${hasNotation ? 'Note' : 'Score'} : ${(20 * ma.mark / ma.bareme).toStringAsFixed(1)} / 20"),
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

  const HomeworkActivityIcon(this.onTap, {super.key});

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
