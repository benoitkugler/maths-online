import 'package:eleve/exercices/fields.dart';
import 'package:eleve/exercices/types.gen.dart';
import 'package:flutter/material.dart';

class ExpressionController extends FieldController {
  final TextEditingController textController;

  ExpressionController(void Function() onChange)
      : textController = TextEditingController() {
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

  @override
  Widget build(BuildContext context) {
    return Container(
      width: 200,
      padding: const EdgeInsets.symmetric(horizontal: 4),
      child: TextField(
        onSubmitted: (_) => onDone(),
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
