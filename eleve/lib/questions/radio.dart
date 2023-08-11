import 'package:eleve/questions/fields.dart';
import 'package:eleve/types/src_maths_questions_client.dart';
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

class RadioFieldW extends StatefulWidget {
  final Color _color;
  final RadioController _controller;

  const RadioFieldW(this._color, this._controller, {Key? key})
      : super(key: key);

  @override
  State<RadioFieldW> createState() => _RadioFieldWState();
}

class _RadioFieldWState extends State<RadioFieldW> {
  @override
  Widget build(BuildContext context) {
    final color = widget._controller.hasError ? Colors.red : widget._color;
    return Container(
      decoration: BoxDecoration(
        border: Border.all(color: color),
        borderRadius: BorderRadius.circular(5),
      ),
      child: Column(
          children: List<Widget>.generate(widget._controller.proposals.length,
              (index) {
        final prop = widget._controller.proposals[index];
        return RadioListTile<int>(
          activeColor: color,
          title: TextRow(buildText(prop, TextS(), 18), verticalPadding: 2),
          value: index,
          groupValue: widget._controller.index,
          onChanged: widget._controller.isEnabled
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
