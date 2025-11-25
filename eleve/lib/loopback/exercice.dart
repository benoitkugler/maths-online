import 'package:eleve/activities/homework/exercice.dart';
import 'package:eleve/loopback/loopback.dart';
import 'package:eleve/types/src_prof_preview.dart';
import 'package:eleve/types/src_sql_homework.dart';

class LoopbackExerciceController implements LoopbackController {
  final ExerciceController controller;
  final LoopbackShowExercice data;

  bool instantShowCorrection = false;

  LoopbackExerciceController(this.data)
    : controller = ExerciceController(
        data.exercice,
        data.progression,
        QuestionRepeat.unlimited,
        0,
      ),
      instantShowCorrection = data.showCorrection;

  // if positive, start at the given question, not in the summary
  int? get initialQuestion => data.nextQuestion >= 0 ? data.nextQuestion : null;
}
