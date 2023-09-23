import 'package:eleve/questions/fields.dart';
import 'package:eleve/types/src_maths_questions_client.dart';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';

class ExpressionController extends FieldController {
  final TextEditingController textController;

  bool _isDirty = false;

  ExpressionController(void Function() onEditDone,
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

class ExpressionFieldW extends StatefulWidget {
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
    var clamped = hintWidth.clamp(10, 27) + 3;
    return clamped.toDouble() / 30.0;
  }

  const ExpressionFieldW(this.color, this._controller,
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
      "sgn",
      "isZero",
      "isPrime",
    ];
    return functions.any((element) =>
        newValue.text.endsWith(element) &&
        !oldValue.text.endsWith("$element("));
  }

  @override
  State<ExpressionFieldW> createState() => _ExpressionFieldWState();
}

class _ExpressionFieldWState extends State<ExpressionFieldW> {
  void _submit() async {
    if (widget._controller.getExpression().isEmpty) {
      // return early
      widget._controller.submit();
      if (widget.onSubmitted != null) widget.onSubmitted!();
      return;
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
          // submitting trigger a focus change : do not call _submit twice
          // onSubmitted:

          inputFormatters: [
            TextInputFormatter.withFunction((oldValue, newValue) {
              if (ExpressionFieldW.isTypingFunc(oldValue, newValue)) {
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

/// [ExpressionCell] wraps an [ExpressionFieldW] in a [TableCell]
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
        child: ExpressionFieldW(
          color,
          controller,
          maxWidthFactor: 0.2,
        ),
      ),
    );
  }
}
