import 'package:eleve/exercices/fields.dart';
import 'package:eleve/exercices/types.gen.dart';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';

class ExpressionController extends FieldController {
  final TextEditingController textController;

  ExpressionController(void Function() onChange)
      : textController = TextEditingController(),
        super(onChange) {
    textController.addListener(onChange);
  }

  @override
  bool hasValidData() {
    final content = textController.text.trim();
    return content.isNotEmpty;
  }

  @override
  Answer getData() {
    final content = textController.text.trim();
    return ExpressionAnswer(content);
  }
}

class ExpressionField extends StatelessWidget {
  final Color _color;
  final TextEditingController _controller;
  final void Function() onDone;

  const ExpressionField(this._color, this._controller, this.onDone, {Key? key})
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
  Widget build(BuildContext context) {
    return Container(
      width: MediaQuery.of(context).size.width * 0.4,
      padding: const EdgeInsets.symmetric(horizontal: 4),
      child: TextField(
        onSubmitted: (_) => onDone(),
        inputFormatters: [
          TextInputFormatter.withFunction((oldValue, newValue) {
            if (isTypingFunc(oldValue, newValue)) {
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
        controller: _controller,
        decoration: InputDecoration(
          isDense: true,
          contentPadding: const EdgeInsets.only(top: 10, bottom: 4),
          focusedBorder: UnderlineInputBorder(
            borderSide: BorderSide(
              color: _color,
            ),
          ),
        ),
        cursorColor: _color,
        style: TextStyle(color: Colors.yellow.shade100, letterSpacing: 1.5),
        textAlign: TextAlign.center,
        textAlignVertical: TextAlignVertical.center,
      ),
    );
  }
}
