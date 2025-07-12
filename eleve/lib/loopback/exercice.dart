import 'package:eleve/exercice/exercice.dart';
import 'package:eleve/loopback/loopback.dart';
import 'package:eleve/types/src.dart';
import 'package:eleve/types/src_prof_preview.dart';
import 'package:eleve/types/src_sql_homework.dart';
import 'package:eleve/types/src_tasks.dart';

class LoopbackExerciceController implements LoopbackController {
  final ExerciceController controller;
  final LoopbackShowExercice data;

  bool instantShowCorrection = false;

  LoopbackExerciceController(this.data)
      : controller = ExerciceController(
            StudentWork(
                data.exercice, data.progression.tryStartFirstQuestion()),
            QuestionRepeat.unlimited,
            0),
        instantShowCorrection = data.showCorrection;

  // if positive, start at the given question, not in the summary
  int? get initialQuestion =>
      data.progression.nextQuestion >= 0 ? data.progression.nextQuestion : null;
}

extension on ProgressionExt {
  // if at summary, set nextQuestion to 0, so that is it enabled
  ProgressionExt tryStartFirstQuestion() {
    return ProgressionExt(questions, nextQuestion < 0 ? 0 : nextQuestion);
  }
}
