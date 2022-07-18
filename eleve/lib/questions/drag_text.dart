import 'package:eleve/questions/fields.dart';
import 'package:eleve/questions/types.gen.dart';
import 'package:flutter/material.dart';

class DragText<T extends Object> extends StatelessWidget {
  final T data; // the associated value
  final TextLine text; // the displayed text
  final bool enabled;
  final bool dense;

  const DragText(this.data, this.text,
      {required this.enabled, this.dense = false, Key? key})
      : super(key: key);

  static const _fontSize = 16.0;

  @override
  Widget build(BuildContext context) {
    return Draggable<T>(
      maxSimultaneousDrags: enabled ? null : 0,
      data: data,
      feedback: Material(
        elevation: 8,
        borderRadius: BorderRadius.circular(10),
        child: Padding(
          padding: const EdgeInsets.all(8),
          child: TextRow(buildText(text, TextS(), _fontSize), lineHeight: 1),
        ),
      ),
      child: Padding(
        padding: EdgeInsets.symmetric(horizontal: dense ? 0 : 6),
        child: Material(
            elevation: 8,
            borderRadius: BorderRadius.circular(dense ? 2 : 8),
            child: Padding(
              padding: EdgeInsets.symmetric(
                vertical: dense ? 4 : 6,
                horizontal: dense ? 8 : 12,
              ),
              child: TextRow(
                buildText(text, TextS(), _fontSize, baselineMiddle: true),
                lineHeight: 1,
                verticalPadding: 2,
              ),
            )),
      ),
    );
  }
}
