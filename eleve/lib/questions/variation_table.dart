import 'dart:math';

import 'package:eleve/questions/fields.dart';
import 'package:eleve/questions/table.dart';
import 'package:eleve/types/src_maths_questions_client.dart';
import 'package:flutter/material.dart';

class VariationTable extends StatelessWidget {
  final VariationTableBlock data;

  const VariationTable(this.data, {Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    final xRow = <Widget>[];
    final fxRow = <Widget>[];

    for (var i = 0; i < data.columns.length; i++) {
      final numberCol = data.columns[i];
      xRow.add(MathTableCell(TableCellVerticalAlignment.middle, numberCol.x));
      fxRow.add(MathTableCell(
        numberCol.isUp
            ? TableCellVerticalAlignment.top
            : TableCellVerticalAlignment.bottom,
        numberCol.y,
      ));

      if (i < data.columns.length - 1) {
        xRow.add(const SizedBox());
        fxRow.add(VariationArrow(data.arrows[i]));
      }
    }

    return BaseFunctionTable(data.label, xRow, fxRow);
  }
}

class VariationArrow extends StatelessWidget {
  final bool isUp;
  final void Function()? onTap;

  const VariationArrow(this.isUp, {Key? key, this.onTap}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    const angle = pi / 4 + pi * 5 / 180;
    return TableCell(
      verticalAlignment: TableCellVerticalAlignment.middle,
      child: InkWell(
        onTap: onTap,
        child: Padding(
          padding: const EdgeInsets.symmetric(vertical: 40.0, horizontal: 5),
          child: Transform.rotate(
              angle: isUp ? -angle : angle,
              child: Transform.scale(
                  scaleX: 3,
                  child: const Icon(
                    IconData(0xe09f, fontFamily: 'MaterialIcons'),
                  ))),
        ),
      ),
    );
  }
}
