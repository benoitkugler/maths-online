import 'dart:math';

import 'package:eleve/questions/fields.dart';
import 'package:eleve/types/src.dart';
import 'package:eleve/types/src_maths_questions_client.dart';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';

class ExpressionController extends FieldController {
  final FieldAPI api;
  final TextEditingController textController;

  bool _isDirty = false;

  ExpressionController(this.api, void Function() onEditDone,
      {bool showFractionHelp = false})
      : textController = TextEditingController(),
        super(onEditDone) {
    if (showFractionHelp) {
      textController.text = "(  ) / (  )";
      textController.selection = const TextSelection.collapsed(offset: 2);
    }
    textController.addListener(() {
      _isDirty = true;
    });
  }

  void submit() {
    _isDirty = false;
    onChange();
  }

  Future<CheckExpressionOut> _checkExpressionSyntax() async {
    return api.checkExpressionSyntax(getExpression());
  }

  @override
  bool hasValidData() {
    final content = textController.text.trim();
    return !_isDirty && content.isNotEmpty;
  }

  String getExpression() {
    return textController.text.trim();
  }

  void setExpression(String expr) {
    textController.text = expr;
    // cleanup the _isDirty flag
    _isDirty = false;
  }

  @override
  Answer getData() {
    return ExpressionAnswer(getExpression());
  }

  @override
  void setData(Answer answer) {
    setExpression((answer as ExpressionAnswer).expression);
  }
}

class ExpressionField extends StatefulWidget {
  final Color color;
  final ExpressionController _controller;

  /// [maxWidthFactor] set the width of this widget
  /// to [maxWidthFactor] * size width, for a full width
  /// hint
  final double maxWidthFactor;

  /// [hintWidth] is a positive number enabling the field to
  /// have lower width than its normal full width.
  final int hintWidth;

  final bool autofocus;
  final void Function()? onSubmitted;

  // returns a float ratio between 0 and 1
  double get hintWidthRatio {
    // add some additional padding
    var clamped = (hintWidth + 3).clamp(5, 30);
    return clamped.toDouble() / 30.0;
  }

  const ExpressionField(this.color, this._controller,
      {Key? key,
      this.maxWidthFactor = 0.9,
      this.hintWidth = 30,
      this.autofocus = false,
      this.onSubmitted})
      : super(key: key);

  static bool isTypingFunc(
      TextEditingValue oldValue, TextEditingValue newValue) {
    /// To keep in sync with the server
    const functions = [
      "exp",
      "ln",
      "log",
      "sin",
      "cos",
      "tan",
      "asin",
      "arcsin",
      "acos",
      "arccos",
      "atan",
      "arctan",
      "abs",
      "sqrt",
    ];
    return functions.any((element) =>
        newValue.text.endsWith(element) &&
        !oldValue.text.endsWith("$element("));
  }

  @override
  State<ExpressionField> createState() => _ExpressionFieldState();
}

class _ExpressionFieldState extends State<ExpressionField> {
  void _showSyntaxError(String reason) {
    if (!mounted) return;
    ScaffoldMessenger.of(context).showSnackBar(SnackBar(
      backgroundColor: Colors.red,
      content: Text.rich(TextSpan(children: [
        const TextSpan(text: "Syntaxe invalide: "),
        TextSpan(
            text: reason, style: const TextStyle(fontWeight: FontWeight.bold)),
      ])),
    ));
  }

  void _submit() async {
    if (widget._controller.getExpression().isEmpty) {
      // return early
      widget._controller.submit();
      if (widget.onSubmitted != null) widget.onSubmitted!();
      return;
    }
    final rep = await widget._controller._checkExpressionSyntax();
    setState(() {
      widget._controller.setError(!rep.isValid);
    });

    if (!rep.isValid) {
      _showSyntaxError(rep.reason);
    }

    widget._controller.submit();
    if (widget.onSubmitted != null) widget.onSubmitted!();
  }

  @override
  Widget build(BuildContext context) {
    final color = widget._controller.hasError ? Colors.red : widget.color;
    final textColor =
        widget._controller.hasError ? Colors.red : Colors.yellow.shade100;

    final width = MediaQuery.of(context).size.width *
        widget.maxWidthFactor *
        widget.hintWidthRatio;
    return Container(
      width: width,
      padding: const EdgeInsets.symmetric(horizontal: 4),
      child: SubmitOnLeave(
        submit: _submit,
        child: TextField(
          autofocus: widget.autofocus,
          enabled: widget._controller.isEnabled,
          onSubmitted: (_) {
            _submit();
          },
          inputFormatters: [
            TextInputFormatter.withFunction((oldValue, newValue) {
              if (ExpressionField.isTypingFunc(oldValue, newValue)) {
                final sel = newValue.selection;
                return newValue.copyWith(
                    text: "${newValue.text}()",
                    selection: sel.copyWith(
                        baseOffset: sel.baseOffset + 1,
                        extentOffset: sel.extentOffset + 1));
              }
              return newValue;
            })
          ],
          keyboardType: TextInputType.visiblePassword,
          controller: widget._controller.textController,
          decoration: InputDecoration(
            isDense: true,
            contentPadding: const EdgeInsets.only(top: 10, bottom: 4),
            focusedBorder: UnderlineInputBorder(
              borderSide: BorderSide(
                color: color,
              ),
            ),
            enabledBorder: UnderlineInputBorder(
              borderSide: BorderSide(
                color: color,
              ),
            ),
          ),
          cursorColor: color,
          style: TextStyle(color: textColor, letterSpacing: 1.5),
          textAlign: TextAlign.center,
          textAlignVertical: TextAlignVertical.center,
        ),
      ),
    );
  }
}

/// [ExpressionCell] wraps an [ExpressionField] in a [TableCell]
class ExpressionCell extends StatelessWidget {
  final Color color;
  final ExpressionController controller;
  final TableCellVerticalAlignment align;

  const ExpressionCell(this.color, this.controller, this.align, {Key? key})
      : super(key: key);

  @override
  Widget build(BuildContext context) {
    return TableCell(
      verticalAlignment: align,
      child: Padding(
        padding: const EdgeInsets.symmetric(vertical: 4.0),
        child: ExpressionField(
          color,
          controller,
          maxWidthFactor: 0.2,
        ),
      ),
    );
  }
}

abstract class _ExprPart {
  String toExpr();
}

class RegularExpr implements _ExprPart {
  final String expr;
  RegularExpr(this.expr);

  @override
  String toExpr() => expr;
}

class FractionExpr implements _ExprPart {
  final String num;
  final String den;
  FractionExpr(this.num, this.den);

  @override
  String toExpr() => "($num) / ($den)";
}

class PowerExpr implements _ExprPart {
  final String base;
  final String exp;

  PowerExpr(this.base, this.exp);

  @override
  String toExpr() => "$base^($exp)";
}

class _FocusNodes {
  final FocusNode first;
  final FocusNode second;
  _FocusNodes()
      : first = FocusNode(),
        second = FocusNode();

  void dispose() {
    first.dispose();
    second.dispose();
  }
}

class ExpressionController2 {
  List<_ExprPart> parts = [];
  List<_FocusNodes> focusNodes = [];

  ExpressionController2(List<_ExprPart> parts) {
    this.parts = parts.map((e) => e).toList();
    focusNodes = parts.map((e) => _FocusNodes()).toList();
  }

  String getExpression() {
    return parts.map((e) => e.toExpr()).join().trim();
  }

  void _insertAt(int index, _ExprPart part, {bool focus = true}) {
    parts.insert(index, part);
    focusNodes.insert(index, _FocusNodes());
    if (focus) focusNodes[index].first.requestFocus();
  }

  void _insertRegularAt(int index) {
    _insertAt(index, RegularExpr(""));
  }

  void _replaceAtByFraction(int index, String before, String after) {
    if (before.codeUnits.first == 0x200b) {
      before = before.substring(1);
    }
    final splitBefore = detectNumerator(before);
    final frac = FractionExpr(splitBefore.after, after);
    // check if there is content before the fraction
    if (splitBefore.before.trim().isNotEmpty) {
      _insertAt(index, RegularExpr(splitBefore.before.trim()), focus: true);
      _replaceAt(index + 1, frac);
    } else {
      _replaceAt(index, frac);
    }
  }

  void _replaceAt(int index, _ExprPart newPart) {
    parts[index] = newPart;
    focusNodes[index].dispose();
    focusNodes[index] = _FocusNodes();
    // focusNodes[index].second.requestFocus();
  }

  void _removeExponentAt(int index) {
    final base = (parts[index] as PowerExpr).base;
    parts[index] = RegularExpr(base);
    focusNodes[index].dispose();
    focusNodes[index] = _FocusNodes();
    focusNodes[index].first.requestFocus();
  }

  void _removeFractionAt(int index) {
    final num = (parts[index] as FractionExpr).num;
    parts[index] = RegularExpr(num);
    focusNodes[index].dispose();
    focusNodes[index] = _FocusNodes();
    focusNodes[index].first.requestFocus();
  }
}

class _RegularEdit extends StatelessWidget {
  final Color color;
  final Color textColor;
  final RegularExpr data;
  final _FocusNodes focus;

  final void Function(String before, String after) onAddFraction;
  final void Function(String before, String after) onAddPower;

  const _RegularEdit(this.color, this.textColor, this.data, this.focus,
      this.onAddFraction, this.onAddPower,
      {super.key});

  @override
  Widget build(BuildContext context) {
    return _PartExprField(
      color,
      textColor,
      data.toExpr(),
      focus.first,
      true,
      onInsertFraction: onAddFraction,
      onInsertPower: onAddPower,
    );
  }
}

class _PartExprField extends StatefulWidget {
  final Color color;
  final Color textColor;
  final String data;
  final FocusNode focus;
  final bool alignCenter;

  final void Function()? onClear;
  final void Function(String before, String after)? onInsertPower;
  final void Function(String before, String after)? onInsertFraction;

  const _PartExprField(
      this.color, this.textColor, this.data, this.focus, this.alignCenter,
      {this.onClear, this.onInsertPower, this.onInsertFraction, super.key});

  @override
  State<_PartExprField> createState() => _PartExprFieldState();
}

class _PartExprFieldState extends State<_PartExprField> {
  late TextEditingController controller;

  @override
  void initState() {
    controller = TextEditingController(text: "\u200b${widget.data}");
    super.initState();
  }

  @override
  void didUpdateWidget(covariant _PartExprField oldWidget) {
    controller = TextEditingController(text: "\u200b${widget.data}");
    super.didUpdateWidget(oldWidget);
  }

  @override
  void dispose() {
    controller.dispose();
    super.dispose();
  }

  void onChanged() {
    final sel = controller.selection;
    if (controller.text.trim().isEmpty && widget.onClear != null) {
      widget.onClear!();
    }
    if (sel.isCollapsed) {
      final before = sel.textBefore(controller.text);
      final after = sel.textAfter(controller.text);
      final lastChar = before.characters.last;
      if (lastChar == "/" && widget.onInsertFraction != null) {
        widget.onInsertFraction!(before.substring(0, before.length - 1), after);
      } else if (lastChar == "^" && widget.onInsertPower != null) {
        widget.onInsertPower!(before.substring(0, before.length - 1), after);
      }
    }
    setState(() {}); // rebuild length
  }

  @override
  Widget build(BuildContext context) {
    const cursorWidth = 2;
    final style = TextStyle(color: widget.textColor, letterSpacing: 1.2);
    // TextField merges given textStyle with text style from current theme
    // Do the same to get final TextStyle
    final themeData = Theme.of(context);
    final effectiveStyle = themeData.textTheme.titleMedium!.merge(style);

    // Use TextPainter to calculate the width of our text
    TextSpan ts = TextSpan(style: effectiveStyle, text: controller.text);
    TextPainter tp = TextPainter(text: ts, textDirection: TextDirection.ltr);
    tp.layout();

    return SizedBox(
      width: max(4, tp.width + cursorWidth),
      child: TextField(
        focusNode: widget.focus,
        onChanged: (_) => onChanged(),
        controller: controller,
        keyboardType: TextInputType.visiblePassword,
        decoration: InputDecoration(
            isDense: true,
            contentPadding: const EdgeInsets.symmetric(vertical: 4),
            focusedBorder: UnderlineInputBorder(
              borderSide: BorderSide(
                color: widget.color,
              ),
            ),
            enabledBorder: InputBorder.none
            // enabledBorder: UnderlineInputBorder(
            //   borderSide: BorderSide(
            //     color: color,
            //   ),
            // ),
            ),
        cursorColor: widget.color,
        style: style,
        textAlign: widget.alignCenter ? TextAlign.center : TextAlign.start,
        textAlignVertical: TextAlignVertical.center,
      ),
    );
  }
}

class _FractionEdit extends StatelessWidget {
  final Color color;
  final Color textColor;
  final FractionExpr data;
  final _FocusNodes focus;

  final void Function() onClearDen;

  const _FractionEdit(
      this.color, this.textColor, this.data, this.focus, this.onClearDen,
      {super.key});

  @override
  Widget build(BuildContext context) {
    final num = _PartExprField(color, textColor, data.num, focus.first, true);
    final den = _PartExprField(color, textColor, data.den, focus.second, true,
        onClear: onClearDen);
    return Padding(
        padding: const EdgeInsets.symmetric(horizontal: 4),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Container(
              decoration: const BoxDecoration(
                  border: Border(bottom: BorderSide(color: Colors.white))),
              padding: const EdgeInsets.only(bottom: 2),
              child: num,
            ),
            Padding(
              padding: const EdgeInsets.only(top: 2),
              child: den,
            ),
          ],
        ));
  }
}

class _PowerEdit extends StatelessWidget {
  final Color color;
  final Color textColor;
  final PowerExpr data;
  final _FocusNodes focus;

  final bool showRightSpace;

  final void Function() onRightSpaceTap;
  final void Function() onClearExponent;

  const _PowerEdit(this.color, this.textColor, this.data, this.focus,
      this.showRightSpace, this.onRightSpaceTap, this.onClearExponent,
      {super.key});

  @override
  Widget build(BuildContext context) {
    final base =
        _PartExprField(color, textColor, data.base, focus.first, false);
    final exp = _PartExprField(color, textColor, data.exp, focus.second, true,
        onClear: onClearExponent);
    return Row(
      mainAxisSize: MainAxisSize.min,
      crossAxisAlignment: CrossAxisAlignment.end,
      children: [
        Column(
          mainAxisSize: MainAxisSize.min,
          children: [base, const SizedBox(height: 8)],
        ),
        Column(
          mainAxisSize: MainAxisSize.min,
          children: [exp, const SizedBox(height: 16)],
        ),
        if (showRightSpace)
          GestureDetector(
              onTap: onRightSpaceTap,
              child: Container(
                color: Colors.yellow,
                width: 15,
                height: 40,
              )),
      ],
    );
  }
}

class ExpressionField2 extends StatefulWidget {
  final List<_ExprPart> data;

  const ExpressionField2(this.data, {super.key});

  @override
  State<ExpressionField2> createState() => _ExpressionField2State();
}

class _ExpressionField2State extends State<ExpressionField2> {
  ExpressionController2 controller = ExpressionController2([]);

  @override
  void initState() {
    controller = ExpressionController2(widget.data);
    super.initState();
  }

  @override
  void didUpdateWidget(covariant ExpressionField2 oldWidget) {
    controller = ExpressionController2(widget.data);
    super.didUpdateWidget(oldWidget);
  }

  void _fixFractionFocus(int i) {
    setState(() {
      controller.focusNodes[i].second.requestFocus();
    });
  }

  @override
  Widget build(BuildContext context) {
    final color = Colors.green;
    final textColor = Colors.yellow.shade100;
    final children = <Widget>[];
    for (var i = 0; i < controller.parts.length; i++) {
      final e = controller.parts[i];
      final focus = controller.focusNodes[i];
      if (e is RegularExpr) {
        children.add(_RegularEdit(
          color,
          textColor,
          e,
          focus,
          (b, a) => setState(() {
            controller._replaceAtByFraction(i, b, a);
            Future.delayed(const Duration(milliseconds: 20),
                () => _fixFractionFocus(i + 1));
          }),
          (b, a) => setState(() => controller._replaceAt(i, PowerExpr(b, a))),
        ));
      } else if (e is PowerExpr) {
        final showSpace = i == controller.parts.length - 1 ||
            (controller.parts[i + 1] is FractionExpr);
        children.add(_PowerEdit(
          color,
          textColor,
          e,
          focus,
          showSpace,
          () => setState(() => controller._insertRegularAt(i + 1)),
          () => setState(() => controller._removeExponentAt(i)),
        ));
      } else if (e is FractionExpr) {
        children.add(_FractionEdit(
          color,
          textColor,
          e,
          focus,
          () => setState(() => controller._removeFractionAt(i)),
        ));
      } else {
        throw "not reachable";
      }
    }

    return SingleChildScrollView(
      scrollDirection: Axis.horizontal,
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.center,
        children: children,
      ),
    );
  }
}

class SplitString {
  final String before;
  final String after;
  const SplitString(this.before, this.after);

  @override
  bool operator ==(Object other) =>
      (other is SplitString) && before == other.before && after == other.after;

  @override
  int get hashCode => before.hashCode + after.hashCode;
}

final a = "a".codeUnits.first;
final z = "z".codeUnits.first;
final A = "A".codeUnits.first;
final Z = "Z".codeUnits.first;
final c0 = "0".codeUnits.first;
final c9 = "9".codeUnits.first;

bool _isLetterOrDigit(Characters range) {
  final r = range.last.codeUnits.first;
  return (a <= r && r <= z || A <= r && r <= Z || c0 <= r && r <= c9);
}

SplitString detectNumerator(String expr) {
  expr = expr.trimRight();
  // detect ()
  var stack = 0;
  var mayHaveEdgeParenthesis = false;
  final iter = expr.characters.iteratorAtEnd;
  while (iter.moveBack()) {
    if (iter.currentCharacters.last == ')') {
      stack += 1;
    } else if (iter.currentCharacters.last == '(') {
      stack -= 1;
    }
    if (stack == 0) {
      mayHaveEdgeParenthesis = true;
      // add letters or digit until space
      while (iter.moveBack()) {
        if (!_isLetterOrDigit(iter.currentCharacters)) break;
      }
      break;
    }
  }

  var numerator = iter.stringAfter;
  if (mayHaveEdgeParenthesis &&
      numerator.characters.first == "(" &&
      numerator.characters.last == ")") {
    numerator = numerator.substring(1, numerator.length - 1);
  }
  return SplitString(iter.stringBefore + iter.current, numerator);
}
