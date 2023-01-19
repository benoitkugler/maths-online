import 'package:eleve/exercice/exercice.dart';
import 'package:eleve/questions/fields.dart';
import 'package:eleve/types/src.dart';
import 'package:eleve/types/src_prof_editor.dart';
import 'package:eleve/types/src_tasks.dart';

class LoopbackExerciceController {
  final ExerciceController controller;
  final LoopbackShowExercice data;

  LoopbackExerciceController(this.data, FieldAPI api)
      : controller = ExerciceController(
            StudentWork(
                data.exercice, data.progression.tryStartFirstQuestion()),
            data.progression.startQuestion,
            api);
}

extension _AdjustStart on ProgressionExt {
  // if positive, start at the given question, not in the summary
  int? get startQuestion => nextQuestion >= 0 ? nextQuestion : null;

  // if at summary, set nextQuestion to 0, so that is it enabled
  ProgressionExt tryStartFirstQuestion() {
    return ProgressionExt(questions, nextQuestion < 0 ? 0 : nextQuestion);
  }
}
