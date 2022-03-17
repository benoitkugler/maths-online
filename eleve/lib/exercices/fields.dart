import 'package:eleve/exercices/types.gen.dart';
import 'package:flutter/material.dart';
import 'package:flutter_math_fork/flutter_math.dart';

abstract class FieldController {
  /// returns true if the field is not empty and contains valid data
  bool hasValidData();

  /// returns the current answer
  Answer getData();
}

WidgetSpan _inlineMath(String content, double fontSize) {
  return WidgetSpan(
    baseline: TextBaseline.alphabetic,
    alignment: PlaceholderAlignment.baseline,
    child: Math.tex(
      content,
      mathStyle: MathStyle.text,
      textScaleFactor: 1.15,
      textStyle: TextStyle(fontSize: fontSize - 1),
    ),
  );
}

List<InlineSpan> buildText(List<TextOrMath> parts, double fontSize) {
  return parts
      .map((part) => part.isMath
          ? _inlineMath(part.text, fontSize)
          : TextSpan(text: part.text))
      .toList();
}

class TextRow extends StatelessWidget {
  final List<InlineSpan> content;
  final double verticalPadding;

  const TextRow(this.content, this.verticalPadding, {Key? key})
      : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: EdgeInsets.symmetric(vertical: verticalPadding),
      child: Text.rich(
        TextSpan(style: const TextStyle(height: 1.5), children: content),
      ),
    );
  }
}
