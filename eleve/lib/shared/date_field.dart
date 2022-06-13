import 'package:flutter/material.dart';
import 'package:flutter/services.dart';

class _NumberField extends StatelessWidget {
  final int length;
  final String hint;
  final bool autofocus;
  final void Function() onDone;
  final void Function() onBack;
  final FocusNode focusNode;
  final TextEditingController controller;

  const _NumberField(this.length, this.hint, this.onDone, this.onBack,
      this.focusNode, this.controller,
      {Key? key, this.autofocus = false})
      : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 5),
      child: SizedBox(
        height: 60,
        width: 5 + length * 15,
        child: TextField(
          controller: controller,
          focusNode: focusNode,
          autofocus: autofocus,
          decoration: InputDecoration(helperText: hint, counterText: ""),
          maxLength: length,
          maxLengthEnforcement: MaxLengthEnforcement.enforced,
          textAlign: TextAlign.center,
          style: const TextStyle(fontSize: 22),
          keyboardType: TextInputType.number,
          onChanged: (content) {
            if (content.length == length) {
              onDone();
            } else if (content.isEmpty) {
              onBack();
            }
          },
        ),
      ),
    );
  }
}

class DateField extends StatefulWidget {
  /// return a formated 2006-12-26 date
  final void Function(String date) onSubmit;

  const DateField(this.onSubmit, {Key? key}) : super(key: key);

  @override
  State<DateField> createState() => _DateFieldState();
}

class _DateFieldState extends State<DateField> {
  final focusNodes = [FocusNode(), FocusNode(), FocusNode()];
  final controllers = [
    TextEditingController(),
    TextEditingController(),
    TextEditingController()
  ]; // day, month, year

  void _onDone() {
    final day = controllers[0].text;
    final month = controllers[1].text;
    final year = controllers[2].text;
    widget.onSubmit(
        "${year.padLeft(4, '0')}-${month.padLeft(2, '0')}-${day.padLeft(2, '0')}");
  }

  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.center,
      children: [
        Row(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            _NumberField(
              2,
              "JJ",
              () => FocusScope.of(context).requestFocus(focusNodes[1]),
              () {},
              focusNodes[0],
              controllers[0],
              autofocus: true,
            ),
            _NumberField(
              2,
              "MM",
              () => FocusScope.of(context).requestFocus(focusNodes[2]),
              () => FocusScope.of(context).requestFocus(focusNodes[0]),
              focusNodes[1],
              controllers[1],
            ),
            _NumberField(
              4,
              "AAAA",
              _onDone,
              () => FocusScope.of(context).requestFocus(focusNodes[1]),
              focusNodes[2],
              controllers[2],
            ),
          ],
        ),
      ],
    );
  }
}
