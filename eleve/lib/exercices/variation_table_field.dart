import 'package:eleve/exercices/fields.dart';
import 'package:eleve/exercices/number.dart';
import 'package:eleve/exercices/types.gen.dart';
import 'package:eleve/exercices/variation_table.dart';
import 'package:flutter/material.dart';

class VariationTableController extends FieldController {
  final int length;
  final List<NumberController> xs; // with length [length]
  final List<NumberController> fxs; // with length [length]
  final List<bool?> arrows; // up, down or not set, with length [length-1]

  VariationTableController(this.length, void Function() onChange)
      : xs = List<NumberController>.generate(
            length, (index) => NumberController(onChange)),
        fxs = List<NumberController>.generate(
            length, (index) => NumberController(onChange)),
        arrows = List<bool?>.filled(length - 1, null),
        super(onChange);

  void toggleArrow(int index) {
    if (arrows[index] == null) {
      arrows[index] = true;
    } else {
      arrows[index] = !(arrows[index]!);
    }
    onChange();
  }

  TableCellVerticalAlignment numberAlignment(int index) {
    if (index == length - 1) {
      final arrow = arrows[index - 1];
      return arrow == null
          ? TableCellVerticalAlignment.middle
          : (arrow
              ? TableCellVerticalAlignment.top
              : TableCellVerticalAlignment.bottom);
    }
    final arrow = arrows[index];
    return arrow == null
        ? TableCellVerticalAlignment.middle
        : (arrow
            ? TableCellVerticalAlignment.bottom
            : TableCellVerticalAlignment.top);
  }

  @override
  void disable() {
    super.disable();
    for (var ct in xs) {
      ct.disable();
    }
    for (var ct in fxs) {
      ct.disable();
    }
  }

  @override
  bool hasValidData() {
    return xs.every((element) => element.hasValidData()) &&
        fxs.every((element) => element.hasValidData()) &&
        arrows.every((element) => element != null);
  }

  @override
  Answer getData() {
    return VariationTableAnswer(
      xs.map((e) => e.getNumber()).toList(),
      fxs.map((e) => e.getNumber()).toList(),
      arrows.map((e) => e!).toList(),
    );
  }
}

class _NumberCell extends StatelessWidget {
  final Color color;
  final NumberController controller;
  final TableCellVerticalAlignment align;

  const _NumberCell(this.color, this.controller, this.align, {Key? key})
      : super(key: key);

  @override
  Widget build(BuildContext context) {
    return TableCell(
      verticalAlignment: align,
      child: Padding(
        padding: const EdgeInsets.symmetric(vertical: 4.0),
        child: NumberField(
          color,
          controller,
          outlined: true,
        ),
      ),
    );
  }
}

class VariationTableField extends StatefulWidget {
  final Color color;
  final VariationTableController controller;

  const VariationTableField(this.color, this.controller, {Key? key})
      : super(key: key);

  @override
  _VariationTableFieldState createState() => _VariationTableFieldState();
}

class _VariationTableFieldState extends State<VariationTableField> {
  @override
  Widget build(BuildContext context) {
    final ct = widget.controller;
    final List<Widget> xRow = [];
    final List<Widget> fxRow = [];
    for (var i = 0; i < ct.length; i++) {
      // number column
      xRow.add(_NumberCell(widget.color, widget.controller.xs[i],
          TableCellVerticalAlignment.middle));
      fxRow.add(_NumberCell(
          widget.color, widget.controller.fxs[i], ct.numberAlignment(i)));

      // arrow column
      if (i < ct.length - 1) {
        final isUp = widget.controller.arrows[i];
        xRow.add(const SizedBox());
        fxRow.add(isUp == null
            ? TableCell(
                verticalAlignment: TableCellVerticalAlignment.middle,
                child: InkWell(
                    borderRadius: BorderRadius.circular(10),
                    child: Container(
                      height: 20,
                      decoration: const BoxDecoration(
                        color: Colors.white,
                        shape: BoxShape.circle,
                      ),
                    ),
                    onTap: () => setState(() {
                          ct.toggleArrow(i);
                        })),
              )
            : VariationArrow(isUp,
                onTap: ct.enabled
                    ? () => setState(() {
                          ct.toggleArrow(i);
                        })
                    : null));
      }
    }

    return BaseVariationTable(xRow, fxRow);
  }
}
