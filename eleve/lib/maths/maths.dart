import 'package:flutter/material.dart';
import 'package:flutter_math_fork/flutter_math.dart';
import 'package:json_annotation/json_annotation.dart';

part 'maths.g.dart';

@JsonSerializable(explicitToJson: true)
class Obj {
  final List<Part> l;

  Obj(this.l);

  factory Obj.fromJson(Map<String, dynamic> json) => _$ObjFromJson(json);

  Map<String, dynamic> toJson() => _$ObjToJson(this);
}

enum PartKind {
  text,
  inlineFormula,
  displayFormula,
  formulaAnswer,
  listAnswer,
}

// @JsonSerializable()
abstract class Part {
  factory Part.fromJson(Map<String, dynamic> json) {
    final type = json["__type__"] as PartKind;
    switch (type) {
      case PartKind.text:
        return TextBlock.fromJson(json);
      case PartKind.inlineFormula:
        return InlineFormula.fromJson(json);
      case PartKind.displayFormula:
        return DisplayFormula.fromJson(json);
      case PartKind.formulaAnswer:
        return FormulaAnswer.fromJson(json);
      case PartKind.listAnswer:
        return ListAnswser.fromJson(json);
    }
  }

  Map<String, dynamic> toJson();
}

@JsonSerializable()
class TextBlock implements Part {
  final String text;

  TextBlock(this.text);

  factory TextBlock.fromJson(Map<String, dynamic> json) =>
      _$TextBlockFromJson(json);

  @override
  Map<String, dynamic> toJson() {
    final out = _$TextBlockToJson(this);
    out["__type__"] = PartKind.text;
    return out;
  }
}

@JsonSerializable()
class InlineFormula implements Part {
  final String latex;

  InlineFormula(this.latex);

  factory InlineFormula.fromJson(Map<String, dynamic> json) =>
      _$InlineFormulaFromJson(json);

  @override
  Map<String, dynamic> toJson() {
    final out = _$InlineFormulaToJson(this);
    out["__type__"] = PartKind.inlineFormula;
    return out;
  }
}

@JsonSerializable()
class DisplayFormula implements Part {
  final String latex;
  DisplayFormula(this.latex);

  factory DisplayFormula.fromJson(Map<String, dynamic> json) =>
      _$DisplayFormulaFromJson(json);

  @override
  Map<String, dynamic> toJson() {
    final out = _$DisplayFormulaToJson(this);
    out["__type__"] = PartKind.displayFormula;
    return out;
  }
}

abstract class _Answer implements Part {}

@JsonSerializable()
class FormulaAnswer extends _Answer {
  FormulaAnswer();

  factory FormulaAnswer.fromJson(Map<String, dynamic> json) =>
      _$FormulaAnswerFromJson(json);

  @override
  Map<String, dynamic> toJson() {
    final out = _$FormulaAnswerToJson(this);
    out["__type__"] = PartKind.formulaAnswer;
    return out;
  }
}

@JsonSerializable()
class ListAnswser extends _Answer {
  ListAnswser();

  factory ListAnswser.fromJson(Map<String, dynamic> json) =>
      _$ListAnswserFromJson(json);

  @override
  Map<String, dynamic> toJson() {
    final out = _$ListAnswserToJson(this);
    out["__type__"] = PartKind.listAnswer;
    return out;
  }
}

class Exercice extends StatefulWidget {
  const Exercice({Key? key}) : super(key: key);

  @override
  _ExerciceState createState() => _ExerciceState();
}

class _ExerciceState extends State<Exercice> {
  @override
  Widget build(BuildContext context) {
    return Column(children: [
      const Text("Factoriser l'expression suivante :"),
      Math.tex(
        r"f(x) = u_{n+1} * 3 + \frac{2x - 5}{4}",
        textStyle: const TextStyle(fontSize: 20),
      ),
      const TextField(),
      TextButton(onPressed: () {}, child: const Text("Valider"))
    ]);
  }
}
