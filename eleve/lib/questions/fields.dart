import 'dart:convert';

import 'package:eleve/build_mode.dart';
import 'package:eleve/shared/hyperlink.dart';
import 'package:eleve/types/src.dart';
import 'package:eleve/types/src_maths_questions_client.dart';
import 'package:flutter/material.dart';
import 'package:flutter_math_fork/flutter_math.dart';
import 'package:http/http.dart' as http;

abstract class FieldController {
  bool _hasError = false;

  /// [hasError] may be set to true to indicate
  /// that the answer does not follow the correct syntax
  /// or is incorrect.
  bool get hasError => _hasError;

  /// [setError] sets the error state.
  /// It can be override for controllers which needs to
  /// update children.
  void setError(bool hasError) {
    _hasError = hasError;
  }

  bool _isEnabled = true;

  /// [isEnabled] is true if the field is actionnable.
  bool get isEnabled => _isEnabled;

  /// [setEnabled] disabled or enabled one field.
  /// If can be overriden for controllers which needs to
  /// update children.
  void setEnabled(bool enabled) {
    _isEnabled = enabled;
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

/// [FieldAPI] provides the server calls required
/// by the fields widgets.
abstract class FieldAPI {
  Future<CheckExpressionOut> checkExpressionSyntax(String expression);
}

/// [ServerFieldAPI] is the default implementation of [FieldAPI]
class ServerFieldAPI implements FieldAPI {
  final BuildMode buildMode;
  const ServerFieldAPI(this.buildMode);

  /// checkExpressionSyntaxCall implements [FieldAPI] with a server call
  @override
  Future<CheckExpressionOut> checkExpressionSyntax(String expression) async {
    final uri = Uri.parse(buildMode.serverURL("/api/check-expression"))
        .replace(queryParameters: {"expression": expression});

    final resp = await http.get(uri);
    return checkExpressionOutFromJson(jsonDecode(resp.body));
  }
}

Math textMath(String content, TextStyle style, {Key? key}) {
  style = style.copyWith(fontSize: (style.fontSize ?? 12) - 1);
  // remove bold since it also changes the font, with indesirable visual effect
  style = style.copyWith(fontWeight: FontWeight.normal);
  return Math.tex(
    content,
    key: key,
    mathStyle: MathStyle.text,
    textScaleFactor: 1.15,
    textStyle: style,
  );
}

// handle line breaking inside TextBlock by retuning a list of text chunk
Iterable<WidgetSpan> _inlineMath(
    String content, TextStyle style, PlaceholderAlignment aligment, Key? key) {
  final tex = textMath(content, style, key: key);
  final breakResult = tex.texBreak();
  return breakResult.parts.map((p) => WidgetSpan(
        baseline: TextBaseline.alphabetic,
        alignment: aligment,
        child: p,
      ));
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

List<InlineSpan> buildText(TextLine parts, TextS style, double fontSize,
    {bool baselineMiddle = false}) {
  final out = <InlineSpan>[];

  final ts = TextStyle(
      fontSize: style.smaller ? fontSize - 2 : fontSize,
      fontStyle: style.italic ? FontStyle.italic : FontStyle.normal,
      fontWeight: style.bold ? FontWeight.bold : FontWeight.normal);
  for (var part in parts) {
    if (part.isMath) {
      out.add(const TextSpan(text: " "));
      out.addAll(_inlineMath(
          part.text,
          ts,
          baselineMiddle
              ? PlaceholderAlignment.middle
              : PlaceholderAlignment.baseline,
          null));
      out.add(const TextSpan(text: " "));
    } else {
      // check for URLs
      out.addAll(parseURLs(part.text, ts));
    }
  }
  return out;
}

class TextRow extends StatelessWidget {
  final List<InlineSpan> content;
  final double verticalPadding;
  final double lineHeight;

  const TextRow(this.content,
      {Key? key, this.verticalPadding = 0, this.lineHeight = 1.5})
      : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: EdgeInsets.symmetric(vertical: verticalPadding),
      child: Text.rich(
        TextSpan(style: TextStyle(height: lineHeight), children: content),
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
  final double? width;
  const MathTableCell(this.align, this.mathContent, {Key? key, this.width})
      : super(key: key);

  @override
  Widget build(BuildContext context) {
    return TableCell(
      verticalAlignment: align,
      child: Container(
        alignment: Alignment.center,
        width: width,
        padding: const EdgeInsets.symmetric(horizontal: 6, vertical: 12),
        child: textMath(mathContent, const TextStyle(fontSize: fontSize - 1)),
      ),
    );
  }
}

/// SubmitOnLeave will call `submit` if it looses focus.
/// It should be used for keyboard entries, where the user
/// will not always validate.
class SubmitOnLeave extends StatelessWidget {
  final void Function() submit;
  final Widget child;
  const SubmitOnLeave({required this.submit, required this.child, Key? key})
      : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Focus(
      onFocusChange: (getFocus) {
        if (!getFocus) submit();
      },
      child: child,
    );
  }
}
