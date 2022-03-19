import 'package:eleve/exercices/fields.dart';
import 'package:eleve/exercices/types.gen.dart';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';

class NumberController extends FieldController {
  final TextEditingController textController;

  NumberController(void Function() onChange)
      : textController = TextEditingController() {
    textController.addListener(onChange);
  }

  String get text => textController.text.trim().replaceAll(RegExp(r","), ".");

  @override
  bool hasValidData() {
    if (text.isEmpty) {
      return false;
    }
    return double.tryParse(text) != null;
  }

  @override
  Answer getData() {
    return NumberAnswer(double.parse(text));
  }
}

class NumberField extends StatelessWidget {
  final Color _color;
  final TextEditingController _controller;
  final void Function() onDone;

  const NumberField(this._color, this._controller, this.onDone, {Key? key})
      : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 4),
      child: SizedBox(
        width: 100,
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
          style: TextStyle(color: Colors.yellow.shade100),
          textAlign: TextAlign.center,
          textAlignVertical: TextAlignVertical.center,
          keyboardType: const TextInputType.numberWithOptions(
              signed: true, decimal: true),
          inputFormatters: <TextInputFormatter>[
            FilteringTextInputFormatter.allow(RegExp(r'[0-9-,\.]')),
          ], // Only numbers, minus, dot and comma can be entered
        ),
      ),
    );
  }
}
