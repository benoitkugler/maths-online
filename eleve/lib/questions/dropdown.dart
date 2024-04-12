import 'package:eleve/questions/fields.dart';
import 'package:eleve/questions/radio.dart';
import 'package:flutter/material.dart';

typedef DropDownController = RadioController;

class DropDownFieldW extends StatefulWidget {
  final Color color;
  final DropDownController controller;

  const DropDownFieldW(this.color, this.controller, {Key? key})
      : super(key: key);

  @override
  _DropDownFieldWState createState() => _DropDownFieldWState();
}

class _DropDownFieldWState extends State<DropDownFieldW> {
  static const fontSize = 16.0;

  @override
  Widget build(BuildContext context) {
    final ct = widget.controller;
    return DropdownButton<int>(
      isDense: true,
      style: widget.controller.hasError
          ? TextStyle(color: Colors.red.shade200)
          : null,
      underline: widget.controller.hasError
          ? Container(height: 1.0, color: Colors.red
              // color: Colors.red,
              )
          : null,
      focusColor: widget.color,
      dropdownColor: widget.color,
      hint: const Text("Choisir"),
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
      onChanged: ct.isEnabled
          ? (v) => setState(() {
                ct.setIndex(v);
              })
          : null,
    );
  }
}
