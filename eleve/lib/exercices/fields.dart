import 'package:eleve/exercices/types.gen.dart';
import 'package:flutter/material.dart';
import 'package:flutter_math_fork/flutter_math.dart';

abstract class FieldController {
  /// [syntaxError] should be set to true to indicate
  /// that the answers doest not follow the correct syntax
  bool syntaxError = false;

  bool _enabled = true;
  bool get enabled => _enabled;

  void disable() {
    _enabled = false;
  }

  /// [onChange] should be called when the state change
  /// to notify the question widget
  final void Function() onChange;

  FieldController(this.onChange);

  /// returns true if the field is not empty and contains valid data
  bool hasValidData();

  /// returns the current answer
  Answer getData();

  /// [setData] set the controller data using the given answer
  void setData(Answer answer);
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

WidgetSpan _inlineMath(
    String content, double fontSize, PlaceholderAlignment aligment, Key? key) {
  return WidgetSpan(
    baseline: TextBaseline.alphabetic,
    alignment: aligment,
    child: textMath(content, fontSize, key: key),
  );
}

class TextS {
  final bool bold;
  final bool italic;
  final bool smaller;

  TextS({this.bold = false, this.italic = false, this.smaller = false});

  factory TextS.fromTextBlock(TextBlock block) {
    return TextS(
        bold: block.bold, italic: block.italic, smaller: block.smaller);
  }
}

List<InlineSpan> buildText(List<TextOrMath> parts, TextS style, double fontSize,
    {bool inTable = false}) {
  final out = <InlineSpan>[];

  final ts = TextStyle(
      fontSize: style.smaller ? fontSize - 2 : fontSize,
      fontStyle: style.italic ? FontStyle.italic : FontStyle.normal,
      fontWeight: style.bold ? FontWeight.bold : FontWeight.normal);
  for (var part in parts) {
    if (part.isMath) {
      out.add(const TextSpan(text: " "));
      out.add(_inlineMath(
          part.text,
          fontSize,
          inTable ? PlaceholderAlignment.middle : PlaceholderAlignment.baseline,
          null));
      out.add(const TextSpan(text: " "));
    } else {
      out.add(TextSpan(text: part.text, style: ts));
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
  static const fontSize = 14.0;

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
          child: textMath(mathContent, fontSize - 1),
        ),
      ),
    );
  }
}
