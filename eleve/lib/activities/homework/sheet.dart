import 'package:eleve/activities/homework/homework.dart';
import 'package:eleve/activities/homework/utils.dart';
import 'package:eleve/exercice/exercice.dart';
import 'package:eleve/exercice/home.dart';
import 'package:eleve/shared/errors.dart';
import 'package:eleve/shared/title.dart';
import 'package:eleve/types/src.dart';
import 'package:eleve/types/src_prof_homework.dart';
import 'package:eleve/types/src_sql_homework.dart';
import 'package:eleve/types/src_sql_tasks.dart';
import 'package:eleve/types/src_tasks.dart';
import 'package:flutter/material.dart';

MarkBareme sheetMark(List<TaskProgressionHeader> tasks) {
  int mark = 0;
  int bareme = 0;
  for (var element in tasks) {
    mark += element.mark;
    bareme += element.bareme;
  }
  return MarkBareme(mark, bareme);
}

class _ExerciceAPI implements ExerciceAPI {
  final HomeworkAPI api;
  final IdTask idTask;
  final IdTravail idTravail;

  StudentEvaluateTaskOut? lastState;

  _ExerciceAPI(this.api, this.idTask, this.idTravail);

  @override
  Future<EvaluateWorkOut> evaluate(EvaluateWorkIn params) async {
    final res = await api.evaluateExercice(idTask, idTravail, params);
    lastState = res;
    return res.ex;
  }

  @override
  Future<CheckExpressionOut> checkExpressionSyntax(String expression) {
    return api.checkExpressionSyntax(expression);
  }
}

class SheetMarkNotification extends Notification {
  final IdSheet idSheet;
  final IdTask idTask;
  final ProgressionExt? newProgression;
  final int newMark;

  const SheetMarkNotification(
      this.idSheet, this.idTask, this.newProgression, this.newMark);

  /// [applyTo] updates [tasks] in place
  void applyTo(List<TaskProgressionHeader> tasks) {
    final index = tasks.indexWhere((element) => element.id == idTask);
    final current = tasks[index];
    tasks[index] = TaskProgressionHeader(
        current.id,
        current.title,
        current.chapter,
        newProgression != null,
        newProgression ?? const ProgressionExt([], 0),
        newMark,
        current.bareme);
  }
}

class SheetW extends StatefulWidget {
  final HomeworkAPI api;
  final SheetProgression sheet;

  const SheetW(this.api, this.sheet, {Key? key}) : super(key: key);

  @override
  State<SheetW> createState() => _SheetWState();
}

class _SheetWState extends State<SheetW> {
  List<TaskProgressionHeader> tasks = [];

  @override
  void initState() {
    tasks = widget.sheet.tasks;
    super.initState();
  }

  void _startExercice(TaskProgressionHeader task) async {
    showDialog<void>(
        barrierDismissible: false,
        context: context,
        builder: (context) => const Dialog(
              child: Padding(
                padding: EdgeInsets.symmetric(vertical: 16.0, horizontal: 12),
                child: Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    Text(
                      "Chargement de l'exercice...",
                      style: TextStyle(fontSize: 16),
                    ),
                    CircularProgressIndicator(),
                  ],
                ),
              ),
            ));
    final instantiatedExercice = await widget.api.loadWork(task.id);
    Navigator.of(context).pop(); // remove the dialog

    final studentEx = StudentWork(instantiatedExercice, task.progression);
    final exeAPI = _ExerciceAPI(widget.api, task.id, widget.sheet.idTravail);
    final exController = ExerciceController(studentEx, null, exeAPI);
    // actually launch the exercice

    // TODO: for now we always show a correction (when available)
    // we might want a setting to let the teacher choose to display it or not
    await Navigator.of(context).push(MaterialPageRoute<void>(
        builder: (context) => ExerciceW(
              exeAPI,
              exController,
              showCorrectionButtonOnFail: true,
            )));

    if (exeAPI.lastState != null) {
      final state = exeAPI.lastState!;
      final notif = SheetMarkNotification(
          widget.sheet.sheet.id, task.id, state.ex.progression, state.mark);

      // locally update the task mark
      setState(() {
        notif.applyTo(tasks);
      });

      // inform the top level sheet widget of the modification
      notif.dispatch(context);
    }
  }

  void _resetTask(TaskProgressionHeader task) async {
    try {
      await widget.api.resetTask(widget.sheet.idTravail, task.id);
    } catch (e) {
      showError("Impossible de recommencer la taĉhe.", e, context);
      return;
    }
    final notif =
        SheetMarkNotification(widget.sheet.sheet.id, task.id, null, 0);
    // locally update the task mark
    setState(() {
      notif.applyTo(tasks);
    });
    // inform the top level sheet widget of the modification
    notif.dispatch(context);
  }

  @override
  Widget build(BuildContext context) {
    final hasNotation = widget.sheet.sheet.noted;
    final isExpired =
        hasNotation && widget.sheet.sheet.deadline.isBefore(DateTime.now());
    return Scaffold(
      appBar: AppBar(title: const Text("Feuille de travail")),
      body: Padding(
        padding: const EdgeInsets.all(8.0),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Padding(
              padding: const EdgeInsets.all(8.0),
              child: Center(
                  child: ColoredTitle(
                      widget.sheet.sheet.title, Colors.blueAccent)),
            ),
            if (hasNotation) ...[
              Padding(
                padding: const EdgeInsets.symmetric(vertical: 8.0),
                child: Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    const Text("Travail noté", style: TextStyle(fontSize: 18)),
                    RichText(
                        text: TextSpan(
                            style: const TextStyle(fontSize: 16),
                            children: [
                          const TextSpan(text: "A rendre avant le\n"),
                          TextSpan(
                              text: formatTime(widget.sheet.sheet.deadline),
                              style:
                                  const TextStyle(fontWeight: FontWeight.bold)),
                        ]))
                  ],
                ),
              ),
              if (isExpired)
                Card(
                  color: SheetStatus.expired.color,
                  child: const Padding(
                    padding: EdgeInsets.all(8.0),
                    child: Text(
                        "La progression et les notes de cette fiche sont verrouillées, car sa date de rendu est dépassée."),
                  ),
                )
            ],
            Expanded(
                child: _TaskList(
              widget.sheet.tasks,
              hasNotation,
              (ex) => _startExercice(ex),
              _resetTask,
            )),
          ],
        ),
      ),
    );
  }
}

class _TaskList extends StatelessWidget {
  final List<TaskProgressionHeader> tasks;
  final bool hasNotation;
  final void Function(TaskProgressionHeader) onStart;
  final void Function(TaskProgressionHeader) onReset;

  const _TaskList(this.tasks, this.hasNotation, this.onStart, this.onReset,
      {Key? key})
      : super(key: key);

  @override
  Widget build(BuildContext context) {
    final children = tasks
        .map((task) => Padding(
            padding: const EdgeInsets.all(4.0),
            child: Container(
              decoration: BoxDecoration(
                borderRadius: BorderRadius.circular(4.0),
                border: Border.all(color: Colors.lightBlue),
              ),
              child: Column(
                mainAxisSize: MainAxisSize.min,
                children: [
                  ListTile(
                    dense: true,
                    onTap: () => onStart(task),
                    leading: getCompletion(task).icon,
                    title: Text(task.title),
                    subtitle: task.chapter.isEmpty ? null : Text(task.chapter),
                    trailing: Text("${task.mark} / ${task.bareme}"),
                  ),
                  if (!hasNotation &&
                      getCompletion(task) == ExerciceCompletion.completed)
                    Row(
                      mainAxisAlignment: MainAxisAlignment.end,
                      children: [
                        Padding(
                          padding: const EdgeInsets.symmetric(horizontal: 8.0),
                          child: ElevatedButton.icon(
                              onPressed: () => onReset(task),
                              icon: const Icon(Icons.refresh),
                              label: const Text("Recommencer")),
                        )
                      ],
                    )
                ],
              ),
            )))
        .toList();

    final total = sheetMark(tasks);
    return ListView(children: [
      ...children,
      if (hasNotation)
        // add the total score of the sheet
        ListTile(
          title: const Text("Total"),
          trailing: Text("${total.mark} / ${total.bareme}"),
        )
    ]);
  }
}

enum ExerciceCompletion { notStarted, started, completed }

ExerciceCompletion getCompletion(TaskProgressionHeader ex) {
  if (!ex.hasProgression) {
    return ExerciceCompletion.notStarted;
  }
  if (ex.progression.isCompleted()) {
    return ExerciceCompletion.completed;
  }
  return ExerciceCompletion.started;
}

extension IconE on ExerciceCompletion {
  Icon get icon {
    switch (this) {
      case ExerciceCompletion.notStarted:
        return const Icon(assignementIcon);
      case ExerciceCompletion.started:
        return const Icon(
            IconData(0xf587,
                fontFamily: 'MaterialIcons', matchTextDirection: true),
            color: Colors.orange);
      case ExerciceCompletion.completed:
        return const Icon(completedIcon, color: Colors.green);
    }
  }
}
