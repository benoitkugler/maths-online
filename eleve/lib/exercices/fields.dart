import 'package:eleve/exercices/types.gen.dart';
import 'package:flutter/material.dart';
import 'package:flutter_math_fork/flutter_math.dart';

abstract class FieldController {
  /// [onChange] should be called when the state change
  /// to notify the question widget
  final void Function() onChange;

  FieldController(this.onChange);

  /// returns true if the field is not empty and contains valid data
  bool hasValidData();

  /// returns the current answer
  Answer getData();
}

Widget textMath(String content, double fontSize, {Key? key}) {
  return Math.tex(
    content,
    key: key,
    mathStyle: MathStyle.text,
    textScaleFactor: 1.15,
    textStyle: TextStyle(fontSize: fontSize - 1),
  );
}

WidgetSpan _inlineMath(String content, double fontSize, Key? key) {
  return WidgetSpan(
    baseline: TextBaseline.alphabetic,
    alignment: PlaceholderAlignment.baseline,
    child: textMath(content, fontSize, key: key),
  );
}

List<InlineSpan> buildText(List<TextOrMath> parts, double fontSize) {
  final out = <InlineSpan>[];
  for (var part in parts) {
    if (part.isMath) {
      out.add(const TextSpan(text: " "));
      out.add(_inlineMath(part.text, fontSize, null));
      out.add(const TextSpan(text: " "));
    } else {
      out.add(TextSpan(text: part.text, style: TextStyle(fontSize: fontSize)));
    }
  }
  return out;
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

/// [MathTableCell] is a [TableCell] containing
/// math text
class MathTableCell extends StatelessWidget {
  static const fontSize = 16.0;

  final TableCellVerticalAlignment align;
  final String mathContent;

  const MathTableCell(this.align, this.mathContent, {Key? key})
      : super(key: key);

  @override
  Widget build(BuildContext context) {
    return TableCell(
      verticalAlignment: align,
      child: Align(
        alignment: Alignment.center,
        child: Padding(
            padding: const EdgeInsets.symmetric(horizontal: 2, vertical: 12),
            child: textMath(mathContent, fontSize)),
      ),
    );
  }
}
