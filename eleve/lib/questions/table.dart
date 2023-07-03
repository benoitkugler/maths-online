import 'package:eleve/questions/fields.dart';
import 'package:eleve/questions/number.dart';
import 'package:eleve/types/src_maths_questions_client.dart';
import 'package:flutter/material.dart';

class _Cell extends StatelessWidget {
  static const fontSize = 14.0;

  final TextOrMath content;
  final bool isHeader;
  final TableCellVerticalAlignment align;

  const _Cell(this.content, this.isHeader,
      {Key? key, this.align = TableCellVerticalAlignment.middle})
      : super(key: key);

  @override
  Widget build(BuildContext context) {
    return TableCell(
        verticalAlignment: align,
        child: Container(
          padding: const EdgeInsets.symmetric(vertical: 4, horizontal: 8),
          color: isHeader ? Colors.grey : Colors.transparent,
          child: Center(
            child: TextRow(
                buildText([content], TextS(), fontSize, baselineMiddle: true),
                verticalPadding: 2),
          ),
        ));
  }
}

class TableW extends StatelessWidget {
  final TableBlock data;

  const TableW(this.data, {Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return _Table(
      data.horizontalHeaders,
      data.verticalHeaders,
      List<List<Widget>>.generate(
        data.values.length,
        (i) => List<Widget>.generate(
          data.values[i].length,
          (j) => _Cell(data.values[i][j], false),
        ),
      ),
    );
  }
}

// common part between static and editable table
class _Table extends StatelessWidget {
  final List<TextOrMath> horizontalHeaders;
  final List<TextOrMath> verticalHeaders;
  final List<List<Widget>> values;

  const _Table(this.horizontalHeaders, this.verticalHeaders, this.values,
      {Key? key})
      : super(key: key);

  @override
  Widget build(BuildContext context) {
    final firstRow = [
      if (verticalHeaders.isNotEmpty) const SizedBox(),
      ...horizontalHeaders.map((e) => _Cell(e, true)).toList(),
    ];
    return SingleChildScrollView(
      scrollDirection: Axis.horizontal,
      child: Table(
          defaultColumnWidth: const IntrinsicColumnWidth(),
          border: TableBorder.all(),
          children: [
            if (horizontalHeaders.isNotEmpty) TableRow(children: firstRow),
            ...List<TableRow>.generate(
                values.length,
                (i) => TableRow(children: [
                      if (verticalHeaders.isNotEmpty)
                        _Cell(
                          verticalHeaders[i],
                          true,
                          align: TableCellVerticalAlignment.fill,
                        ),
                      ...values[i],
                    ])),
          ]),
    );
  }
}

// editable table

class TableController extends FieldController {
  final TableFieldBlock data;
  final List<List<NumberController>> _controllers;

  TableController(this.data, void Function() onChange)
      : _controllers = List<List<NumberController>>.generate(
            data.verticalHeaders.length,
            (_) => List<NumberController>.generate(
                data.horizontalHeaders.length,
                (_) => NumberController(onChange))),
        super(onChange);

  @override
  void setEnabled(bool enabled) {
    super.setEnabled(enabled);
    for (var row in _controllers) {
      for (var cell in row) {
        cell.setEnabled(enabled);
      }
    }
  }

  @override
  bool hasValidData() {
    return _controllers
        .every((row) => row.every((cell) => cell.hasValidData()));
  }

  @override
  Answer getData() {
    return TableAnswer(_controllers
        .map((row) => row
            .map(
              (e) => e.getNumber(),
            )
            .toList())
        .toList());
  }

  @override
  void setData(Answer answer) {
    final rows = (answer as TableAnswer).rows;
    for (var i = 0; i < _controllers.length; i++) {
      final ctRow = _controllers[i];
      for (var j = 0; j < ctRow.length; j++) {
        ctRow[j].setNumber(rows[i][j]);
      }
    }
  }
}

class TableField extends StatefulWidget {
  final Color color;
  final TableController controller;

  const TableField(this.color, this.controller, {Key? key}) : super(key: key);

  @override
  _TableFieldState createState() => _TableFieldState();
}

class _TableFieldState extends State<TableField> {
  @override
  Widget build(BuildContext context) {
    final ct = widget.controller;
    final color = ct.hasError ? Colors.red : widget.color;
    return _Table(
        ct.data.horizontalHeaders,
        ct.data.verticalHeaders,
        ct._controllers
            .map((row) => row
                .map((cell) => TableCell(
                        child: Padding(
                      padding: const EdgeInsets.symmetric(vertical: 4),
                      child: NumberField(
                        color,
                        cell,
                        outlined: true,
                      ),
                    )))
                .toList())
            .toList());
  }
}

/// [BaseFunctionTable] serves as container for
/// variation and sign tables, editable or not.
class BaseFunctionTable extends StatelessWidget {
  final List<Widget> xRow;
  final List<MapEntry<String, List<Widget>>> fxRows; // (label, row) pairs
  final Color? headerColor;

  const BaseFunctionTable(this.xRow, this.fxRows, {Key? key, this.headerColor})
      : super(key: key);

  @override
  Widget build(BuildContext context) {
    final color = headerColor ?? Colors.grey.shade600;
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
            TableRow(decoration: BoxDecoration(color: color), children: [
              const MathTableCell(TableCellVerticalAlignment.middle, "x"),
              ...xRow
            ]),
            ...fxRows.map((e) => TableRow(children: [
                  MathTableCell(TableCellVerticalAlignment.middle, e.key,
                      width: 70),
                  ...e.value,
                ]))
          ]),
    );
  }
}
