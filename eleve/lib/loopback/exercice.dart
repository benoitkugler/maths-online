import 'package:eleve/exercice/exercice.dart';
import 'package:eleve/questions/fields.dart';
import 'package:eleve/types/src.dart';
import 'package:eleve/types/src_prof_editor.dart';

class LoopbackExerciceController {
  final ExerciceController controller;
  final LoopbackShowExercice data;

  LoopbackExerciceController(this.data, ServerFieldAPI api)
      : controller = ExerciceController(
            StudentWork(data.exercice, data.progression), null, api);
}
