import 'package:eleve/questions/expression.dart';
import 'package:eleve/questions/fields.dart';
import 'package:eleve/questions/sign_table.dart';
import 'package:eleve/questions/table.dart';
import 'package:eleve/questions/types.gen.dart';
import 'package:flutter/material.dart';

/// [_STController] is the controller for one table (the length is known)
class _STController {
  bool enabled = true;

  final List<ExpressionController> xs; // with length [signsLength+1]
  final List<SignSymbol> fxs; // with length [signsLength+1]
  final List<bool?>
      signs; // is positive ? or not set, with length [signsLength]

  final void Function() onChange;

  // zerosLength exclude edges
  _STController(FieldAPI api, int signsLength, this.onChange)
      : xs = List<ExpressionController>.generate(
            signsLength + 1, (index) => ExpressionController(api, onChange)),
        fxs = List<SignSymbol>.generate(
            signsLength + 1, (index) => SignSymbol.nothing),
        signs = List<bool?>.filled(signsLength, null);

  void toggleSign(int index) {
    if (signs[index] == null) {
      signs[index] = true;
    } else {
      signs[index] = !(signs[index]!);
    }
    onChange();
  }

  void onSymbolClick(int index) {
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
        signs.every((element) => element != null);
  }

  Answer getData() {
    return SignTableAnswer(
      xs.map((e) => e.getExpression()).toList(),
      fxs,
      signs.map((e) => e!).toList(),
    );
  }

  void setData(SignTableAnswer answer) {
    for (var i = 0; i < xs.length; i++) {
      xs[i].setExpression(answer.xs[i]);
    }
    for (var i = 0; i < fxs.length; i++) {
      fxs[i] = answer.fxSymbols[i];
    }
    for (var i = 0; i < signs.length; i++) {
      signs[i] = answer.signs[i];
    }
  }
}

class SignTableController extends FieldController {
  final FieldAPI api;
  final SignTableFieldBlock data;

  // setup when selecting a length
  _STController? ct;

  int? get selectedSignsLength => ct == null ? null : ct!.signs.length;

  SignTableController(this.api, this.data, void Function() onChange)
      : super(onChange);

  void setSignsLength(int? length) {
    ct = length == null ? null : _STController(api, length, onChange);
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
    final ans = answer as SignTableAnswer;
    setSignsLength(ans.signs.length);
    ct!.setData(ans);
  }
}

class SignTableField extends StatefulWidget {
  final Color color;
  final SignTableController controller;

  const SignTableField(this.color, this.controller, {Key? key})
      : super(key: key);

  @override
  SignTableFieldState createState() => SignTableFieldState();
}

class SignTableFieldState extends State<SignTableField> {
  void _resetArrowLength() {
    setState(() {
      widget.controller.ct = null;
    });
  }

  @override
  Widget build(BuildContext context) {
    final ct = widget.controller;
    return ct.selectedSignsLength == null
        ? Container(
            padding: const EdgeInsets.all(12),
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
                                    primary: widget.color),
                                child: Text(e.toString()),
                                onPressed: () => setState(() {
                                      ct.setSignsLength(e);
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
  final _STController controller;
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

      fxRow.add(
        SignSymbolButton(
            ct.fxs[i],
            ct.enabled
                ? () => setState(() {
                      ct.onSymbolClick(i);
                    })
                : null),
      );

      // sign columns
      if (i < ct.xs.length - 1) {
        final isUp = widget.controller.signs[i];
        xRow.add(const SizedBox());
        fxRow.add(isUp == null
            ? TableCell(
                verticalAlignment: TableCellVerticalAlignment.middle,
                child: InkWell(
                    borderRadius: BorderRadius.circular(10),
                    child: Container(
                      height: 20,
                      width: 80,
                      decoration: const BoxDecoration(
                        color: Colors.white,
                        shape: BoxShape.circle,
                      ),
                    ),
                    onTap: () => setState(() {
                          ct.toggleSign(i);
                        })),
              )
            : SignButton(isUp,
                onTap: ct.enabled
                    ? () => setState(() {
                          ct.toggleSign(i);
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
      BaseFunctionTable(widget.functionLabel, xRow, fxRow,
          headerColor: Colors.red.shade200)
    ]);
  }
}

class SignButton extends StatelessWidget {
  final bool isPositive;
  final void Function()? onTap;

  const SignButton(this.isPositive, {Key? key, this.onTap}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return TableCell(
        verticalAlignment: TableCellVerticalAlignment.middle,
        child: SizedBox(
          width: 80,
          child: TextButton(
            child: Text(
              isPositive ? "+" : "-",
              style: const TextStyle(fontSize: 18),
            ),
            onPressed: onTap,
          ),
        ));
  }
}

class SignSymbolButton extends StatelessWidget {
  final SignSymbol symbol;
  final void Function()? onTap;

  const SignSymbolButton(
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
        text = const Text("0");
        break;
      case SignSymbol.forbiddenValue:
        text = ForbiddenValueW(color: color);
        break;
    }
    return TableCell(
        verticalAlignment: TableCellVerticalAlignment.middle,
        child: OutlinedButton(
          child: text,
          onPressed: onTap,
        ));
  }
}
