import 'package:eleve/exercices/fields.dart';
import 'package:eleve/exercices/types.gen.dart';
import 'package:flutter/material.dart';

class _WidgetPair {
  final Widget x;
  final Widget fx;
  _WidgetPair(this.x, this.fx);

  factory _WidgetPair.fromData(SignColumn data) {
    if (data.isSign) {
      return _WidgetPair(
          const SizedBox(),
          MathTableCell(
              TableCellVerticalAlignment.middle, data.isPositive ? "+" : "-"));
    }

    return _WidgetPair(
      MathTableCell(TableCellVerticalAlignment.middle, data.x),
      data.isYForbiddenValue
          ? const _ForbiddenValue()
          : MathTableCell(
              TableCellVerticalAlignment.middle, data.isPositive ? "0" : ""),
    );
  }
}

class SignTable extends StatelessWidget {
  final SignTableBlock data;

  const SignTable(this.data, {Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    final xRow = <Widget>[];
    final fxRow = <Widget>[];
    for (var element in data.columns) {
      final pair = _WidgetPair.fromData(element);
      xRow.add(pair.x);
      fxRow.add(pair.fx);
    }

    return SingleChildScrollView(
      scrollDirection: Axis.horizontal,
      child: Table(
          border: const TableBorder(
            top: BorderSide(width: 1.5),
            left: BorderSide(width: 1.5),
            right: BorderSide(width: 1.5),
            bottom: BorderSide(width: 1.5),
            horizontalInside: BorderSide(),
          ),
          defaultColumnWidth: const IntrinsicColumnWidth(),
          children: [
            TableRow(
                decoration: BoxDecoration(color: Colors.red.shade200),
                children: [
                  const MathTableCell(TableCellVerticalAlignment.middle, "x"),
                  ...xRow
                ]),
            TableRow(children: [
              MathTableCell(TableCellVerticalAlignment.middle, data.label),
              ...fxRow,
            ])
          ]),
    );
  }
}

class _ForbiddenValue extends StatelessWidget {
  static const height = 40.0;

  const _ForbiddenValue({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    final textColor =
        Theme.of(context).textTheme.bodyMedium?.color ?? Colors.black;
    return TableCell(
      verticalAlignment: TableCellVerticalAlignment.middle,
      child: LayoutBuilder(
        builder: (context, constraints) => Align(
          child: Container(
            width: 5,
            height: height,
            decoration: BoxDecoration(
                color: Colors.transparent,
                border: Border(
                  left: BorderSide(color: textColor),
                  right: BorderSide(color: textColor),
                )),
          ),
        ),
      ),
    );
  }
}
