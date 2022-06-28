import 'package:eleve/build_mode.dart';
import 'package:eleve/exercice/home.dart';
import 'package:eleve/exercice/question.dart';
import 'package:eleve/shared_gen.dart';
import 'package:flutter/material.dart';

/// ExerciceW is the widget providing one exercice to
/// the student.
/// It is used in the editor loopback, and as the base for
/// at home training activity
class ExerciceW extends StatefulWidget {
  final BuildMode buildMode;

  /// [data] stores the server instantiated exercice with
  /// the initial progression state.
  final Exercice data;

  const ExerciceW(this.buildMode, this.data, {Key? key}) : super(key: key);

  @override
  State<ExerciceW> createState() => _ExerciceWState();
}

class _ExerciceWState extends State<ExerciceW> {
  void showQuestion(int index) {
    Navigator.of(context).push(MaterialPageRoute<void>(
        builder: (ct) => QuestionPage(
              widget.buildMode,
              widget.data.exercice.questions[index].question,
              (p0) {
                print(p0); // TODO:
              },
            )));
  }

  @override
  Widget build(BuildContext context) {
    return ExerciceHome(widget.data, showQuestion);
  }
}
