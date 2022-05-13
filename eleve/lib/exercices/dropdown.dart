import 'package:eleve/exercices/fields.dart';
import 'package:eleve/exercices/radio.dart';
import 'package:flutter/material.dart';

typedef DropDownController = RadioController;

class DropDownField extends StatefulWidget {
  final Color color;
  final DropDownController controller;

  const DropDownField(this.color, this.controller, {Key? key})
      : super(key: key);

  @override
  _DropDownFieldState createState() => _DropDownFieldState();
}

class _DropDownFieldState extends State<DropDownField> {
  static const fontSize = 16.0;

  @override
  Widget build(BuildContext context) {
    final ct = widget.controller;
    return DropdownButton<int>(
      isDense: true,
      focusColor: widget.color,
      dropdownColor: widget.color,
      hint: const Text("   "),
      value: ct.index,
      iconSize: 0,
      alignment: Alignment.center,
      // we use selectedItemBuilder since Math.tex do not handle
      // keys in a way that permit reusing the widgets in items
      selectedItemBuilder: (_) => List.generate(
        ct.proposals.length,
        (index) => Padding(
          padding: const EdgeInsets.symmetric(horizontal: 5.0),
          child: TextRow(buildText(ct.proposals[index], TextS(), fontSize)),
        ),
      ),
      items: List.generate(
          ct.proposals.length,
          (index) => DropdownMenuItem<int>(
                value: index,
                child: Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 3),
                  child: TextRow(
                      buildText(ct.proposals[index], TextS(), fontSize),
                      verticalPadding: 1),
                ),
              )),
      onChanged: ct.enabled
          ? (v) => setState(() {
                ct.setIndex(v);
              })
          : null,
    );
  }
}
