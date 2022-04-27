import 'package:eleve/exercices/fields.dart';
import 'package:eleve/exercices/number.dart';
import 'package:eleve/exercices/types.gen.dart';
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
                buildText([content], false, fontSize, inTable: true), 2),
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
  void disable() {
    super.disable();
    for (var row in _controllers) {
      for (var cell in row) {
        cell.disable();
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
    return _Table(
        ct.data.horizontalHeaders,
        ct.data.verticalHeaders,
        ct._controllers
            .map((row) => row
                .map((cell) => TableCell(
                        child: Padding(
                      padding: const EdgeInsets.symmetric(vertical: 4),
                      child: NumberField(
                        widget.color,
                        cell,
                        outlined: true,
                      ),
                    )))
                .toList())
            .toList());
  }
}
