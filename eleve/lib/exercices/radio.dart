import 'package:eleve/exercices/fields.dart';
import 'package:eleve/exercices/types.gen.dart';
import 'package:flutter/material.dart';

class RadioController extends FieldController {
  List<ListFieldProposal> proposals;
  int? index;
  final void Function() onChange;

  RadioController(this.onChange, this.proposals);

  @override
  bool hasValidData() {
    return index != null;
  }

  @override
  Answer getData() {
    return RadioAnswer(index!);
  }

  void setIndex(int? index) {
    this.index = index;
    onChange();
  }
}

class RadioField extends StatefulWidget {
  final Color _color;
  final RadioController _controller;

  const RadioField(this._color, this._controller, {Key? key}) : super(key: key);

  @override
  State<RadioField> createState() => _RadioFieldState();
}

class _RadioFieldState extends State<RadioField> {
  @override
  Widget build(BuildContext context) {
    return Column(
        children:
            List<Widget>.generate(widget._controller.proposals.length, (index) {
      final prop = widget._controller.proposals[index];
      return RadioListTile<int>(
        title: TextRow(buildText(prop.content, 18), 2),
        value: index,
        groupValue: widget._controller.index,
        onChanged: (int? value) {
          setState(() {
            widget._controller.setIndex(value);
          });
        },
      );
    }));
  }
}