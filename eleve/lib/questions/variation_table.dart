import 'dart:math';

import 'package:eleve/questions/fields.dart';
import 'package:eleve/questions/types.gen.dart';
import 'package:flutter/material.dart';

class BaseVariationTable extends StatelessWidget {
  final String label;
  final List<Widget> xRow;
  final List<Widget> fxRow;

  const BaseVariationTable(this.label, this.xRow, this.fxRow, {Key? key})
      : super(key: key);

  @override
  Widget build(BuildContext context) {
    return SingleChildScrollView(
      scrollDirection: Axis.horizontal,
      child: Table(
          defaultColumnWidth: const IntrinsicColumnWidth(),
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
              MathTableCell(TableCellVerticalAlignment.middle, label,
                  width: 70),
              ...fxRow,
            ])
          ]),
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

    return BaseVariationTable(data.label, xRow, fxRow);
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
