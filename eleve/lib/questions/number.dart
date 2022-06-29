import 'package:eleve/questions/fields.dart';
import 'package:eleve/questions/types.gen.dart';
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
    final s = n.truncateToDouble() == n ? n.toInt().toString() : n.toString();
    textController.text = s;
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
  final void Function()? onSubmitted;

  const NumberField(this._color, this._controller,
      {Key? key,
      this.outlined = false,
      this.autofocus = false,
      this.onSubmitted})
      : super(key: key);

  Color get color => _controller.fieldError ? Colors.red : _color;

  @override
  Widget build(BuildContext context) {
    final border = outlined
        ? OutlineInputBorder(
            borderSide: BorderSide(
              color: color,
            ),
          )
        : UnderlineInputBorder(
            borderSide: BorderSide(
              color: color,
            ),
          );
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 4),
      child: SizedBox(
        width: 80,
        child: SubmitOnLeave(
          submit: onSubmitted ?? () {},
          child: TextField(
            enabled: _controller.enabled,
            onSubmitted: onSubmitted != null ? (_) => onSubmitted!() : null,
            autofocus: autofocus,
            controller: _controller.textController,
            decoration: InputDecoration(
              isDense: true,
              contentPadding: const EdgeInsets.only(top: 10, bottom: 2),
              disabledBorder: border,
              focusedBorder: border,
              enabledBorder: border,
              border: border,
            ),

            cursorColor: color,
            style: TextStyle(
                color: _controller.fieldError
                    ? Colors.red.shade200
                    : Colors.white),
            textAlign: TextAlign.center,
            textAlignVertical: TextAlignVertical.center,
            keyboardType: const TextInputType.numberWithOptions(
                signed: true, decimal: true),
            inputFormatters: <TextInputFormatter>[
              FilteringTextInputFormatter.allow(RegExp(r'[0-9-,\.]')),
            ], // Only numbers, minus, dot and comma can be entered
          ),
        ),
      ),
    );
  }
}
