import 'dart:convert';

import 'package:eleve/build_mode.dart';
import 'package:eleve/exercices/fields.dart';
import 'package:eleve/exercices/types.gen.dart';
import 'package:eleve/shared_gen.dart';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:http/http.dart' as http;

class ExpressionController extends FieldController {
  final BuildMode buildMode;
  final TextEditingController textController;

  bool _isDirty = false;

  ExpressionController(this.buildMode, void Function() onEditDone)
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

  Future<CheckExpressionOut> _checkExpressionSyntax() async {
    final uri = Uri.parse(buildMode.serverURL("/api/check-expression"))
        .replace(queryParameters: {"expression": getExpression()});

    final resp = await http.get(uri);
    return checkExpressionOutFromJson(jsonDecode(resp.body));
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
  final double widthFactor;

  const ExpressionField(this.color, this._controller,
      {Key? key, this.widthFactor = 0.4})
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
      return;
    }
    final rep = await widget._controller._checkExpressionSyntax();
    setState(() {
      widget._controller.syntaxError = !rep.isValid;
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
    final color = widget._controller.syntaxError ? Colors.red : widget.color;
    final textColor =
        widget._controller.syntaxError ? Colors.red : Colors.yellow.shade100;
    return Container(
      width: MediaQuery.of(context).size.width * widget.widthFactor,
      padding: const EdgeInsets.symmetric(horizontal: 4),
      child: Focus(
        onFocusChange: (getFocus) {
          if (!getFocus) _submit();
        },
        child: TextField(
          enabled: widget._controller.enabled,
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
