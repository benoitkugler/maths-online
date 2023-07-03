import 'package:eleve/questions/fields.dart';
import 'package:eleve/questions/table.dart';
import 'package:eleve/types/src_maths_questions_client.dart';
import 'package:flutter/material.dart';

Widget _forSymbol(SignSymbol symbol) {
  switch (symbol) {
    case SignSymbol.forbiddenValue:
      return const _ForbiddenValue();
    case SignSymbol.zero:
      return const MathTableCell(TableCellVerticalAlignment.middle, "0");
    case SignSymbol.nothing:
      return const SizedBox();
  }
}

class SignTable extends StatelessWidget {
  final SignTableBlock data;

  const SignTable(this.data, {Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    final xRow = <Widget>[];
    // alternate value and void
    for (var x in data.xs) {
      xRow.add(MathTableCell(TableCellVerticalAlignment.middle, x));
      xRow.add(const SizedBox());
    }
    // remove trailing void
    xRow.removeLast();

    final fxRows = data.functions.map((function) {
      final fxRow = <Widget>[];
      // alternate symbols and signs
      for (var i = 0; i < function.fxSymbols.length; i++) {
        fxRow.add(_forSymbol(function.fxSymbols[i]));
        if (i != function.fxSymbols.length - 1) {
          final signIsPositive = function.signs[i];
          fxRow.add(MathTableCell(
              TableCellVerticalAlignment.middle, signIsPositive ? "+" : "-"));
        }
      }
      return MapEntry(function.label, fxRow);
    }).toList();

    return BaseFunctionTable(xRow, fxRows, headerColor: Colors.red.shade200);
  }
}

class ForbiddenValueW extends StatelessWidget {
  static const height = 40.0;

  /// override theme text color
  final Color? color;

  const ForbiddenValueW({Key? key, this.color}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    final textColor =
        color ?? Theme.of(context).textTheme.bodyMedium?.color ?? Colors.black;
    return Align(
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
    );
  }
}

class _ForbiddenValue extends StatelessWidget {
  const _ForbiddenValue({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return const TableCell(
      verticalAlignment: TableCellVerticalAlignment.middle,
      child: ForbiddenValueW(),
    );
  }
}
