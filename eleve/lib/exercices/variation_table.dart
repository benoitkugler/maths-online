import 'dart:math';

import 'package:eleve/exercices/fields.dart';
import 'package:eleve/exercices/types.gen.dart';
import 'package:flutter/material.dart';

class _WidgetPair {
  final Widget x;
  final Widget fx;
  _WidgetPair(this.x, this.fx);

  factory _WidgetPair.fromData(VariationColumn data) {
    if (data.isArrow) {
      return _WidgetPair(const SizedBox(), _Arrow(data.isUp));
    }

    return _WidgetPair(
      MathTableCell(TableCellVerticalAlignment.middle, data.x),
      MathTableCell(
        data.isUp
            ? TableCellVerticalAlignment.top
            : TableCellVerticalAlignment.bottom,
        data.y,
      ),
    );
  }
}

class VariationTable extends StatelessWidget {
  final VariationTableBlock data;

  const VariationTable(this.data, {Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    final xRow = <Widget>[];
    final fxRow = <Widget>[];
    for (var element in data.columns) {
      final pair = _WidgetPair.fromData(element);
      xRow.add(pair.x);
      fxRow.add(pair.fx);
    }

    return Table(
        border: const TableBorder(
          top: BorderSide(width: 1.5),
          left: BorderSide(width: 1.5),
          right: BorderSide(width: 1.5),
          bottom: BorderSide(width: 1.5),
          horizontalInside: BorderSide(),
        ),
        children: [
          TableRow(
              decoration: BoxDecoration(color: Colors.grey.shade600),
              children: [
                const MathTableCell(TableCellVerticalAlignment.middle, "x"),
                ...xRow
              ]),
          TableRow(children: [
            const MathTableCell(TableCellVerticalAlignment.middle, "f(x)"),
            ...fxRow,
          ])
        ]);
  }
}

class _Arrow extends StatelessWidget {
  final bool isUp;
  const _Arrow(this.isUp, {Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    const angle = pi / 4 + pi * 5 / 180;
    return TableCell(
      verticalAlignment: TableCellVerticalAlignment.middle,
      child: Padding(
        padding: const EdgeInsets.symmetric(vertical: 40.0),
        child: Transform.rotate(
            angle: isUp ? -angle : angle,
            child: Transform.scale(
                scaleX: 3,
                child:
                    const Icon(IconData(0xe09f, fontFamily: 'MaterialIcons')))),
      ),
    );
  }
}