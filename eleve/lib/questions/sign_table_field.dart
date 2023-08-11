import 'package:eleve/questions/expression.dart';
import 'package:eleve/questions/fields.dart';
import 'package:eleve/questions/sign_table.dart';
import 'package:eleve/questions/table.dart';
import 'package:eleve/types/src_maths_questions_client.dart';
import 'package:flutter/material.dart';

class _FunctionSign {
  final List<SignSymbol> fxs; // with length [signsLength+1]
  final List<bool?>
      signs; // is positive ? or not set, with length [signsLength]

  _FunctionSign(int signsLength)
      : fxs = List<SignSymbol>.generate(
            signsLength + 1, (index) => SignSymbol.nothing),
        signs = List<bool?>.filled(signsLength, null);
}

/// [_STController] is the controller for one table (the length and number of functions are fixed)
class _STController {
  bool enabled = true;

  final List<ExpressionController> xs; // with length [signsLength+1]
  final List<_FunctionSign> functions;

  final void Function() onChange;

  _STController(
      FieldAPI api, int signsLength, int functionsLength, this.onChange)
      : xs = List<ExpressionController>.generate(
            signsLength + 1, (index) => ExpressionController(api, onChange)),
        functions = List<_FunctionSign>.generate(
            functionsLength, (index) => _FunctionSign(signsLength));

  void toggleSign(int row, int index) {
    final signs = functions[row].signs;
    if (signs[index] == null) {
      signs[index] = true;
    } else {
      signs[index] = !(signs[index]!);
    }
    onChange();
  }

  void onSymbolClick(int row, int index) {
    final fxs = functions[row].fxs;
    final newIndex = (fxs[index].index + 1) % SignSymbol.values.length;
    fxs[index] = SignSymbol.values[newIndex];
    onChange();
  }

  void setEnabled(bool enabled) {
    enabled = false;
    for (var ct in xs) {
      ct.setEnabled(enabled);
    }
  }

  bool hasValidData() {
    return xs.every((element) => element.hasValidData()) &&
        functions.every((row) => row.signs.every((element) => element != null));
  }

  Answer getData() {
    return SignTableAnswer(
      xs.map((e) => e.getExpression()).toList(),
      functions
          .map((e) => FunctionSign("", e.fxs, e.signs.map((e) => e!).toList()))
          .toList(),
    );
  }

  void setData(SignTableAnswer answer) {
    for (var i = 0; i < xs.length; i++) {
      xs[i].setExpression(answer.xs[i]);
    }
    for (var i = 0; i < functions.length; i++) {
      final inFunction = answer.functions[i];
      final fxs = functions[i].fxs;
      final signs = functions[i].signs;
      for (var j = 0; j < fxs.length; j++) {
        fxs[j] = inFunction.fxSymbols[j];
      }
      for (var j = 0; j < signs.length; j++) {
        signs[j] = inFunction.signs[j];
      }
    }
  }
}

class SignTableController extends FieldController {
  final FieldAPI api;
  final SignTableFieldBlock data;

  // setup when selecting a length
  _STController? _ct;

  int? get selectedSignsLength => _ct == null ? null : _ct!.xs.length - 1;

  SignTableController(this.api, this.data, void Function() onChange)
      : super(onChange);

  void setSignsLength(int? length) {
    _ct = length == null
        ? null
        : _STController(api, length, data.labels.length, onChange);
  }

  @override
  void setEnabled(bool enabled) {
    super.setEnabled(enabled);
    _ct?.setEnabled(enabled);
  }

  @override
  bool hasValidData() {
    return _ct != null && _ct!.hasValidData();
  }

  @override
  Answer getData() {
    return _ct!.getData();
  }

  @override
  void setData(Answer answer) {
    final ans = answer as SignTableAnswer;
    setSignsLength(ans.xs.length - 1);
    _ct!.setData(ans);
  }
}

class SignTableFieldW extends StatefulWidget {
  final Color color;
  final SignTableController controller;

  const SignTableFieldW(this.color, this.controller, {Key? key})
      : super(key: key);

  @override
  SignTableFieldWState createState() => SignTableFieldWState();
}

class SignTableFieldWState extends State<SignTableFieldW> {
  void _resetArrowLength() {
    setState(() {
      widget.controller._ct = null;
    });
  }

  @override
  Widget build(BuildContext context) {
    final ct = widget.controller;
    return ct.selectedSignsLength == null
        ? Container(
            padding: const EdgeInsets.all(12),
            decoration: BoxDecoration(
              border: Border.all(color: widget.color),
              borderRadius: BorderRadius.circular(5),
            ),
            child: Column(
              mainAxisSize: MainAxisSize.min,
              children: [
                const Text(
                  "Choisir le nombre de signes :",
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
                                    backgroundColor: widget.color),
                                child: Text(e.toString()),
                                onPressed: () => setState(() {
                                      ct.setSignsLength(e);
                                    })),
                          ),
                        )
                        .toList()),
              ],
            ),
          )
        : _OneTable(ct.hasError ? Colors.red : widget.color, ct._ct!,
            ct.data.labels, ct.isEnabled ? _resetArrowLength : null);
  }
}

class _OneTable extends StatefulWidget {
  final Color color;
  final _STController controller;
  final List<String> functionLabels;
  final void Function()? onBack;

  const _OneTable(this.color, this.controller, this.functionLabels, this.onBack,
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
    for (var i = 0; i < ct.xs.length; i++) {
      // alternate
      xRow.add(ExpressionCell(widget.color, widget.controller.xs[i],
          TableCellVerticalAlignment.middle));
      xRow.add(const SizedBox());
    }
    xRow.removeLast();

    final fxRows = List.generate(ct.functions.length, (i) {
      final List<Widget> fxRow = [];
      final fxs = ct.functions[i].fxs;
      final signs = ct.functions[i].signs;
      for (var j = 0; j < fxs.length; j++) {
        // alternate
        fxRow.add(
          _SignSymbolButton(
              fxs[j],
              ct.enabled
                  ? () => setState(() {
                        ct.onSymbolClick(i, j);
                      })
                  : null),
        );
        // sign columns
        if (j < fxs.length - 1) {
          final isUp = signs[j];

          fxRow.add(isUp == null
              ? TableCell(
                  verticalAlignment: TableCellVerticalAlignment.middle,
                  child: Container(
                      padding: const EdgeInsets.all(10),
                      width: 80,
                      child: Ink(
                        height: 30,
                        decoration: const ShapeDecoration(
                          color: Colors.lightBlue,
                          shape: CircleBorder(),
                        ),
                        child: IconButton(
                          splashRadius: 24,
                          padding: EdgeInsets.zero,
                          color: Colors.white,
                          icon: const Icon(Icons.question_mark),
                          onPressed: () => setState(() {
                            ct.toggleSign(i, j);
                          }),
                        ),
                      )))
              : _SignButton(isUp,
                  onTap: ct.enabled
                      ? () => setState(() {
                            ct.toggleSign(i, j);
                          })
                      : null));
        }
      }
      return MapEntry(widget.functionLabels[i], fxRow);
    }).toList();

    return Column(children: [
      Row(children: [
        IconButton(
            onPressed: widget.onBack,
            icon: const Icon(IconData(0xe092,
                fontFamily: 'MaterialIcons', matchTextDirection: true)))
      ]),
      BaseFunctionTable(xRow, fxRows, headerColor: Colors.red.shade200)
    ]);
  }
}

class _SignButton extends StatelessWidget {
  final bool isPositive;
  final void Function()? onTap;

  const _SignButton(this.isPositive, {Key? key, this.onTap}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return TableCell(
        verticalAlignment: TableCellVerticalAlignment.middle,
        child: SizedBox(
          width: 80,
          child: Padding(
            padding: const EdgeInsets.symmetric(horizontal: 4.0),
            child: OutlinedButton(
              onPressed: onTap,
              child: Text(
                isPositive ? "+" : "-",
                style:
                    const TextStyle(fontSize: 20, fontWeight: FontWeight.bold),
              ),
            ),
          ),
        ));
  }
}

class _SignSymbolButton extends StatelessWidget {
  final SignSymbol symbol;
  final void Function()? onTap;

  const _SignSymbolButton(
    this.symbol,
    this.onTap, {
    Key? key,
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    final color = Theme.of(context).colorScheme.primary;
    Widget text;
    switch (symbol) {
      case SignSymbol.nothing:
        text = const Text("");
        break;
      case SignSymbol.zero:
        text = const Text("0",
            style: TextStyle(fontSize: 18, fontWeight: FontWeight.bold));
        break;
      case SignSymbol.forbiddenValue:
        text = ForbiddenValueW(color: color);
        break;
    }
    return TableCell(
        verticalAlignment: TableCellVerticalAlignment.middle,
        child: OutlinedButton(
          onPressed: onTap,
          child: text,
        ));
  }
}
