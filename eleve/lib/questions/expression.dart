import 'package:eleve/questions/fields.dart';
import 'package:eleve/questions/types.gen.dart';
import 'package:eleve/shared_gen.dart' as shared;
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';

class ExpressionController extends FieldController {
  final FieldAPI api;
  final TextEditingController textController;

  bool _isDirty = false;

  ExpressionController(this.api, void Function() onEditDone)
      : textController = TextEditingController(),
        super(onEditDone) {
    textController.addListener(() {
      _isDirty = true;
    });
  }

  void submit() {
    _isDirty = false;
    onChange();
  }

  Future<shared.CheckExpressionOut> _checkExpressionSyntax() async {
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
    print("setExpr $expr");
    textController.text = expr;
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

  // returns a float ratio between 0 and 1
  double get hintWidthRatio {
    // add some additional padding
    var clamped = (hintWidth + 3).clamp(5, 30);
    return clamped.toDouble() / 30.0;
  }

  const ExpressionField(this.color, this._controller,
      {Key? key, this.maxWidthFactor = 0.9, this.hintWidth = 30})
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
        !oldValue.text.endsWith(element + "("));
  }

  @override
  State<ExpressionField> createState() => _ExpressionFieldState();
}

class _ExpressionFieldState extends State<ExpressionField> {
  void _submit() async {
    if (widget._controller.getExpression().isEmpty) {
      // return early
      widget._controller.submit();
      return;
    }
    final rep = await widget._controller._checkExpressionSyntax();
    setState(() {
      widget._controller.setError(!rep.isValid);
    });

    if (!rep.isValid) {
      ScaffoldMessenger.of(context).showSnackBar(SnackBar(
        backgroundColor: Colors.red,
        content: Text.rich(TextSpan(children: [
          const TextSpan(text: "Syntaxe invalide: "),
          TextSpan(
              text: rep.reason,
              style: const TextStyle(fontWeight: FontWeight.bold)),
        ])),
      ));
    }

    widget._controller.submit();
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
          enabled: widget._controller.isEnabled,
          onSubmitted: (_) {
            _submit();
          },
          inputFormatters: [
            TextInputFormatter.withFunction((oldValue, newValue) {
              if (ExpressionField.isTypingFunc(oldValue, newValue)) {
                final sel = newValue.selection;
                return newValue.copyWith(
                    text: newValue.text + "()",
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
