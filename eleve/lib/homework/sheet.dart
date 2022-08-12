import 'package:eleve/exercice/exercice.dart';
import 'package:eleve/exercice/home.dart';
import 'package:eleve/homework/homework.dart';
import 'package:eleve/homework/types.gen.dart';
import 'package:eleve/homework/utils.dart';
import 'package:eleve/shared/title.dart';
import 'package:eleve/shared_gen.dart';
import 'package:flutter/material.dart';

class MarkBareme {
  final int mark;
  final int bareme;
  MarkBareme(this.mark, this.bareme);
}

MarkBareme taskMark(List<TaskProgressionHeader> tasks) {
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

  StudentEvaluateExerciceOut? lastState;

  _ExerciceAPI(this.api, this.idTask);

  @override
  Future<EvaluateExerciceOut> evaluate(EvaluateExerciceIn params) async {
    final res = await api.evaluateExercice(idTask, params);
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
  final ProgressionExt newProgression;
  final int newMark;

  const SheetMarkNotification(
      this.idSheet, this.idTask, this.newProgression, this.newMark);

  /// [updateTasks] updates [tasks] in place
  void updateTasks(List<TaskProgressionHeader> tasks) {
    final index = tasks.indexWhere((element) => element.id == idTask);
    final current = tasks[index];
    tasks[index] = TaskProgressionHeader(current.id, current.idExercice,
        current.titleExercice, true, newProgression, newMark, current.bareme);
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

  void _startExercice(TaskProgressionHeader task, BuildContext context) async {
    showDialog<void>(
        barrierDismissible: false,
        context: context,
        builder: (context) => Dialog(
              child: Padding(
                padding:
                    const EdgeInsets.symmetric(vertical: 16.0, horizontal: 12),
                child: Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: const [
                    Text(
                      "Chargement de l'exercice...",
                      style: TextStyle(fontSize: 16),
                    ),
                    CircularProgressIndicator(),
                  ],
                ),
              ),
            ));
    final instantiatedExercice = await widget.api.loadExercice(task.idExercice);
    Navigator.of(context).pop(); // remove the dialog

    final studentEx =
        StudentExerciceInst(instantiatedExercice, task.progression);
    final controller = _ExerciceAPI(widget.api, task.id);
    // actually launch the exercice
    await Navigator.of(context).push(MaterialPageRoute<void>(
        builder: (context) => ExerciceW(controller, studentEx)));

    if (controller.lastState != null) {
      final state = controller.lastState!;
      final notif = SheetMarkNotification(
          widget.sheet.sheet.id, task.id, state.ex.progression, state.mark);

      // locally update the task mark
      setState(() {
        notif.updateTasks(tasks);
      });

      // inform the top level sheet widget of the modification
      notif.dispatch(context);
    }
  }

  @override
  Widget build(BuildContext context) {
    final hasNotation = widget.sheet.sheet.notation != Notation.noNotation;
    return Scaffold(
      appBar: AppBar(title: const Text("Fiche de travail")),
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
            if (hasNotation)
              Padding(
                padding: const EdgeInsets.symmetric(vertical: 8.0),
                child: Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    const Text("Travail noté", style: TextStyle(fontSize: 18)),
                    Text("pour le ${formatTime(widget.sheet.sheet.deadline)}",
                        style: const TextStyle(fontSize: 18)),
                  ],
                ),
              ),
            Expanded(
                child: _TaskList(widget.sheet.tasks, hasNotation,
                    (ex) => _startExercice(ex, context))),
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

  const _TaskList(this.tasks, this.hasNotation, this.onStart, {Key? key})
      : super(key: key);

  @override
  Widget build(BuildContext context) {
    final children = tasks
        .map((ex) => ListTile(
              onTap: () => onStart(ex),
              leading: getCompletion(ex).icon,
              title: Text(ex.titleExercice),
              trailing: hasNotation ? Text("${ex.mark} / ${ex.bareme}") : null,
            ))
        .toList();
    if (hasNotation) {
      // add the total score of the sheet
      final total = taskMark(tasks);
      children.add(ListTile(
        title: const Text("Total"),
        trailing: Text("${total.mark} / ${total.bareme}"),
      ));
    }
    return ListView(children: children);
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
        return const Icon(IconData(0xf587,
            fontFamily: 'MaterialIcons', matchTextDirection: true));
      case ExerciceCompletion.started:
        return const Icon(
            IconData(0xf587,
                fontFamily: 'MaterialIcons', matchTextDirection: true),
            color: Colors.orange);
      case ExerciceCompletion.completed:
        return const Icon(IconData(0xe156, fontFamily: 'MaterialIcons'),
            color: Colors.green);
    }
  }
}
