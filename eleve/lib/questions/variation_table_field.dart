import 'package:eleve/questions/expression.dart';
import 'package:eleve/questions/fields.dart';
import 'package:eleve/questions/table.dart';
import 'package:eleve/questions/types.gen.dart';
import 'package:eleve/questions/variation_table.dart';
import 'package:flutter/material.dart';

/// [_VTController] is the controller for one table (the length is known)
class _VTController {
  bool enabled = true;

  final List<ExpressionController> xs; // with length [length]
  final List<ExpressionController> fxs; // with length [length]
  final List<bool?> arrows; // up, down or not set, with length [length-1]

  final void Function() onChange;

  _VTController(FieldAPI api, int arrowLength, this.onChange)
      : xs = List<ExpressionController>.generate(
            arrowLength + 1, (index) => ExpressionController(api, onChange)),
        fxs = List<ExpressionController>.generate(
            arrowLength + 1, (index) => ExpressionController(api, onChange)),
        arrows = List<bool?>.filled(arrowLength, null);

  void toggleArrow(int index) {
    if (arrows[index] == null) {
      arrows[index] = true;
    } else {
      arrows[index] = !(arrows[index]!);
    }
    onChange();
  }

  TableCellVerticalAlignment numberAlignment(int index) {
    if (index == xs.length - 1) {
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

  void setEnabled(bool enabled) {
    for (var ct in xs) {
      ct.setEnabled(enabled);
    }
    for (var ct in fxs) {
      ct.setEnabled(enabled);
    }
  }

  bool hasValidData() {
    return xs.every((element) => element.hasValidData()) &&
        fxs.every((element) => element.hasValidData()) &&
        arrows.every((element) => element != null);
  }

  Answer getData() {
    return VariationTableAnswer(
      xs.map((e) => e.getExpression()).toList(),
      fxs.map((e) => e.getExpression()).toList(),
      arrows.map((e) => e!).toList(),
    );
  }

  void setData(VariationTableAnswer answer) {
    for (var i = 0; i < xs.length; i++) {
      xs[i].setExpression(answer.xs[i]);
    }
    for (var i = 0; i < fxs.length; i++) {
      fxs[i].setExpression(answer.fxs[i]);
    }
    for (var i = 0; i < arrows.length; i++) {
      arrows[i] = answer.arrows[i];
    }
  }
}

class VariationTableController extends FieldController {
  final FieldAPI api;
  final VariationTableFieldBlock data;

  // setup when selecting a length
  _VTController? ct;

  int? get selectedArrowLength => ct == null ? null : ct!.arrows.length;

  VariationTableController(this.api, this.data, void Function() onChange)
      : super(onChange);

  void setArrowLength(int? length) {
    ct = length == null ? null : _VTController(api, length, onChange);
  }

  @override
  void setEnabled(bool enabled) {
    super.setEnabled(enabled);
    ct?.setEnabled(enabled);
  }

  @override
  bool hasValidData() {
    return ct != null && ct!.hasValidData();
  }

  @override
  Answer getData() {
    return ct!.getData();
  }

  @override
  void setData(Answer answer) {
    final ans = answer as VariationTableAnswer;
    setArrowLength(ans.arrows.length);
    ct!.setData(ans);
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
  void _resetArrowLength() {
    setState(() {
      widget.controller.ct = null;
    });
  }

  @override
  Widget build(BuildContext context) {
    final ct = widget.controller;
    return ct.selectedArrowLength == null
        ? Container(
            padding: const EdgeInsets.all(12),
            child: Column(
              mainAxisSize: MainAxisSize.min,
              children: [
                const Text(
                  "Choisir le nombre de flèches :",
                  style: TextStyle(fontStyle: FontStyle.italic, fontSize: 14),
                ),
                Row(
                    mainAxisAlignment: MainAxisAlignment.center,
                    children: ct.data.lengthProposals
                        .map(
                          (e) => Padding(
                            padding: const EdgeInsets.all(8.0),
                            child: ElevatedButton(
                                style: ElevatedButton.styleFrom(
                                    primary: widget.color),
                                child: Text(e.toString()),
                                onPressed: () => setState(() {
                                      ct.setArrowLength(e);
                                    })),
                          ),
                        )
                        .toList()),
              ],
            ),
            decoration: BoxDecoration(
              border: Border.all(color: widget.color),
              borderRadius: BorderRadius.circular(5),
            ),
          )
        : _OneTable(ct.hasError ? Colors.red : widget.color, ct.ct!,
            ct.data.label, ct.isEnabled ? _resetArrowLength : null);
  }
}

class _OneTable extends StatefulWidget {
  final Color color;
  final _VTController controller;
  final String functionLabel;
  final void Function()? onBack;

  const _OneTable(this.color, this.controller, this.functionLabel, this.onBack,
      {Key? key})
      : super(key: key);

  @override
  __OneTableState createState() => __OneTableState();
}

class __OneTableState extends State<_OneTable> {
  @override
  Widget build(BuildContext context) {
    final ct = widget.controller;
    final List<Widget> xRow = [];
    final List<Widget> fxRow = [];
    for (var i = 0; i < ct.xs.length; i++) {
      // number column
      xRow.add(ExpressionCell(widget.color, widget.controller.xs[i],
          TableCellVerticalAlignment.middle));
      fxRow.add(ExpressionCell(
          widget.color, widget.controller.fxs[i], ct.numberAlignment(i)));

      // arrow column
      if (i < ct.xs.length - 1) {
        final isUp = widget.controller.arrows[i];
        xRow.add(const SizedBox());
        fxRow.add(isUp == null
            ? TableCell(
                verticalAlignment: TableCellVerticalAlignment.middle,
                child: InkWell(
                    borderRadius: BorderRadius.circular(10),
                    child: Container(
                      height: 20,
                      width: 30,
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

    return Column(children: [
      Row(children: [
        IconButton(
            onPressed: widget.onBack,
            icon: const Icon(IconData(0xe092,
                fontFamily: 'MaterialIcons', matchTextDirection: true)))
      ]),
      BaseFunctionTable(widget.functionLabel, xRow, fxRow)
    ]);
  }
}
