import 'package:eleve/exercices/fields.dart';
import 'package:eleve/exercices/types.gen.dart';
import 'package:flutter/material.dart';

class RadioController extends FieldController {
  List<TextLine> proposals;
  int? index;

  RadioController(void Function() onChange, this.proposals) : super(onChange);

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

  @override
  void setData(Answer answer) {
    setIndex((answer as RadioAnswer).index);
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
    return Container(
      decoration: BoxDecoration(
        border: Border.all(color: widget._color),
        borderRadius: BorderRadius.circular(5),
      ),
      child: Column(
          children: List<Widget>.generate(widget._controller.proposals.length,
              (index) {
        final prop = widget._controller.proposals[index];
        return RadioListTile<int>(
          activeColor: widget._color,
          title: TextRow(buildText(prop, TextS(), 18), verticalPadding: 2),
          value: index,
          groupValue: widget._controller.index,
          onChanged: widget._controller.enabled
              ? (int? value) {
                  setState(() {
                    widget._controller.setIndex(value);
                  });
                }
              : null,
        );
      })),
    );
  }
}
