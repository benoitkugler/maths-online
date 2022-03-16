import 'package:eleve/exercices/types.gen.dart';

abstract class FieldController {
  /// returns true if the field is not empty and contains valid data
  bool hasValidData();

  /// returns the current answer
  Answer getData();
}
