import 'package:eleve/exercices/fields.dart';
import 'package:eleve/exercices/types.gen.dart';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';

class NumberController extends FieldController {
  final TextEditingController textController;

  NumberController(void Function() onChange)
      : textController = TextEditingController(),
        super(onChange) {
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

  double getNumber() => double.parse(text);

  @override
  Answer getData() {
    return NumberAnswer(getNumber());
  }

  void setNumber(double n) {
    textController.text = n.toString();
  }

  @override
  void setData(Answer answer) {
    setNumber((answer as NumberAnswer).value);
  }
}

class NumberField extends StatelessWidget {
  final Color _color;
  final NumberController _controller;
  final bool outlined;
  final bool autofocus;
  final void Function(String)? onSubmitted;

  const NumberField(this._color, this._controller,
      {Key? key,
      this.outlined = false,
      this.autofocus = false,
      this.onSubmitted})
      : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 4),
      child: SizedBox(
        width: 50,
        child: TextField(
          enabled: _controller.enabled,
          onSubmitted: onSubmitted,
          autofocus: autofocus,
          controller: _controller.textController,
          decoration: InputDecoration(
            isDense: true,
            contentPadding: const EdgeInsets.only(top: 10, bottom: 4),
            focusedBorder: outlined
                ? OutlineInputBorder(
                    borderSide: BorderSide(
                      color: _color,
                    ),
                  )
                : UnderlineInputBorder(
                    borderSide: BorderSide(
                      color: _color,
                    ),
                  ),
            border: outlined
                ? OutlineInputBorder(
                    borderSide: BorderSide(
                      color: _color,
                    ),
                  )
                : null,
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
