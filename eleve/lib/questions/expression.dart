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

class ExpressionController2 {
  List<_ExprPart> parts = [];

  ExpressionController2(this.parts);

  String getExpression() {
    return parts.map((e) => e.toExpr()).join().trim();
  }
}

class _RegularStatic extends StatelessWidget {
  final RegularExpr data;
  const _RegularStatic(this.data, {super.key});

  @override
  Widget build(BuildContext context) {
    return Text(data.toExpr(), style: TextStyle(letterSpacing: 1.5));
  }
}

class _RegularEdit extends StatelessWidget {
  final Color color;
  final Color textColor;
  final RegularExpr data;

  const _RegularEdit(this.color, this.textColor, this.data, {super.key});

  @override
  Widget build(BuildContext context) {
    return _PartExprField(color, textColor, data.toExpr());
  }
}

class _PartExprField extends StatefulWidget {
  final Color color;
  final Color textColor;
  final String data;

  const _PartExprField(this.color, this.textColor, this.data, {super.key});

  @override
  State<_PartExprField> createState() => _PartExprFieldState();
}

class _PartExprFieldState extends State<_PartExprField> {
  late TextEditingController controller;

  @override
  void initState() {
    controller = TextEditingController(text: widget.data);
    super.initState();
  }

  @override
  void didUpdateWidget(covariant _PartExprField oldWidget) {
    controller = TextEditingController(text: widget.data);
    super.didUpdateWidget(oldWidget);
  }

  @override
  Widget build(BuildContext context) {
    return SizedBox(
      width: 4 + controller.text.length * 10,
      child: TextField(
        onChanged: (_) => setState(() {}),
        controller: controller,
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
        style: TextStyle(color: widget.textColor, letterSpacing: 1.5),
        textAlign: TextAlign.center,
        textAlignVertical: TextAlignVertical.center,
      ),
    );
  }
}

class _FracLayout extends StatelessWidget {
  final Widget num;
  final Widget den;

  const _FracLayout(this.num, this.den, {super.key});

  @override
  Widget build(BuildContext context) {
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

class _FractionStatic extends StatelessWidget {
  final FractionExpr data;

  const _FractionStatic(this.data, {super.key});

  @override
  Widget build(BuildContext context) {
    return _FracLayout(Text(data.num, style: TextStyle(letterSpacing: 1.5)),
        Text(data.den, style: TextStyle(letterSpacing: 1.5)));
  }
}

class _FractionEdit extends StatelessWidget {
  final Color color;
  final Color textColor;
  final FractionExpr data;

  const _FractionEdit(this.color, this.textColor, this.data, {super.key});

  @override
  Widget build(BuildContext context) {
    return _FracLayout(_PartExprField(color, textColor, data.num),
        _PartExprField(color, textColor, data.den));
  }
}

class _PowerLayout extends StatelessWidget {
  final Widget base;
  final Widget exp;
  const _PowerLayout(this.base, this.exp, {super.key});

  @override
  Widget build(BuildContext context) {
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
      ],
    );
  }
}

class _PowerStatic extends StatelessWidget {
  final PowerExpr data;

  const _PowerStatic(this.data, {super.key});

  @override
  Widget build(BuildContext context) {
    return _PowerLayout(Text(data.base, style: TextStyle(letterSpacing: 1.5)),
        Text(data.exp, style: TextStyle(letterSpacing: 1.5)));
  }
}

class _PowerEdit extends StatelessWidget {
  final Color color;
  final Color textColor;
  final PowerExpr data;

  const _PowerEdit(this.color, this.textColor, this.data, {super.key});

  @override
  Widget build(BuildContext context) {
    return _PowerLayout(_PartExprField(color, textColor, data.base),
        _PartExprField(color, textColor, data.exp));
  }
}

class ExpressionField2 extends StatelessWidget {
  final ExpressionController2 data;

  const ExpressionField2(this.data, {super.key});

  @override
  Widget build(BuildContext context) {
    final color = Colors.green;
    final textColor = Colors.yellow.shade100;
    return SingleChildScrollView(
      scrollDirection: Axis.horizontal,
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.center,
        children: data.parts.map((e) {
          if (e is RegularExpr) {
            return _RegularEdit(color, textColor, e);
          } else if (e is PowerExpr) {
            return _PowerEdit(color, textColor, e);
          } else if (e is FractionExpr) {
            return _FractionEdit(color, textColor, e);
          }
          throw "not reachable";
        }).toList(),
      ),
    );
  }
}
