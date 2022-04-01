import 'package:eleve/exercices/dropdown.dart';
import 'package:eleve/exercices/expression.dart';
import 'package:eleve/exercices/fields.dart';
import 'package:eleve/exercices/figure_point.dart';
import 'package:eleve/exercices/figure_vector.dart';
import 'package:eleve/exercices/figure_vector_pair.dart';
import 'package:eleve/exercices/function_graph.dart';
import 'package:eleve/exercices/function_points.dart';
import 'package:eleve/exercices/number.dart';
import 'package:eleve/exercices/ordered_list.dart';
import 'package:eleve/exercices/probas_tree.dart';
import 'package:eleve/exercices/radio.dart';
import 'package:eleve/exercices/repere.dart';
import 'package:eleve/exercices/sign_table.dart';
import 'package:eleve/exercices/table.dart';
import 'package:eleve/exercices/types.gen.dart';
import 'package:eleve/exercices/variation_table.dart';
import 'package:eleve/exercices/variation_table_field.dart';
import 'package:eleve/trivialpoursuit/timeout_bar.dart';
import 'package:flutter/material.dart';
import 'package:flutter_math_fork/flutter_math.dart';

WidgetSpan _inlineMath(String content, double fontSize) {
  return WidgetSpan(
    baseline: TextBaseline.alphabetic,
    alignment: PlaceholderAlignment.baseline,
    child: Math.tex(
      content,
      mathStyle: MathStyle.text,
      textScaleFactor: 1.15,
      textStyle: TextStyle(fontSize: fontSize - 1),
    ),
  );
}

/// utility class used to layout the Block
class _ContentBuilder {
  final void Function(int) onFieldDone;
  final List<Block> _content;
  final Color _color;

  /// field controllers created by [initControllers]
  final Map<int, FieldController> _controllers;

  final List<Widget> rows = []; // final output
  final List<GlobalKey> zoomableWigets = [];

  List<InlineSpan> _currentRow = []; // current row
  bool lastIsText = false; // used to insert new line between to text block
  static const fontSize = 18.0;

  _ContentBuilder(
      this.onFieldDone, this._content, this._controllers, this._color);

  /// walks throught the question content and creates field controllers,
  /// later used when building widgets
  static Map<int, FieldController> initControllers(
      List<Block> content, void Function() onChange) {
    final controllers = <int, FieldController>{};
    for (var block in content) {
      if (block is NumberFieldBlock) {
        controllers[block.iD] = NumberController(onChange);
      } else if (block is ExpressionFieldBlock) {
        controllers[block.iD] = ExpressionController(onChange);
      } else if (block is RadioFieldBlock) {
        controllers[block.iD] = RadioController(onChange, block.proposals);
      } else if (block is DropDownFieldBlock) {
        controllers[block.iD] = DropDownController(onChange, block.proposals);
      } else if (block is OrderedListFieldBlock) {
        controllers[block.iD] = OrderedListController(onChange, block);
      } else if (block is FigurePointFieldBlock) {
        controllers[block.iD] = FigurePointController(onChange);
      } else if (block is FigureVectorFieldBlock) {
        controllers[block.iD] = FigureVectorController(block, onChange);
      } else if (block is FigureVectorPairFieldBlock) {
        controllers[block.iD] =
            FigureVectorPairController(block.figure, onChange);
      } else if (block is VariationTableFieldBlock) {
        controllers[block.iD] =
            VariationTableController(block.length, onChange);
      } else if (block is FunctionPointsFieldBlock) {
        controllers[block.iD] = FunctionPointsController(block, onChange);
      } else if (block is TreeFieldBlock) {
        controllers[block.iD] = TreeController(block, onChange);
      } else if (block is TableFieldBlock) {
        controllers[block.iD] = TableController(block, onChange);
      }

      // TODO: handle more fields
    }
    return controllers;
  }

  void _flushCurrentRow() {
    if (_currentRow.isEmpty) {
      return;
    }

    rows.add(TextRow(_currentRow, 6));
    _currentRow = [];
  }

  void _handleTextBlock(TextBlock element) {
    if (lastIsText) {
      _flushCurrentRow();
    }
    _currentRow.addAll(buildText(element.parts, element.isHint, fontSize));
  }

  void _handleFormulaBlock(FormulaBlock element) {
    // start a new row
    _flushCurrentRow();

    rows.add(Center(
        child: Math.tex(
      element.formula,
      mathStyle: MathStyle.display,
      textStyle: const TextStyle(fontSize: fontSize),
    )));
  }

  void _handleVariationTableBlock(VariationTableBlock element) {
    // start a new row
    _flushCurrentRow();

    rows.add(Center(child: VariationTable(element)));
  }

  void _handleSignTableBlock(SignTableBlock element) {
    // start a new row
    _flushCurrentRow();

    rows.add(Center(child: SignTable(element)));
  }

  void _handleFigureBlock(FigureBlock element) {
    // start a new row
    _flushCurrentRow();

    rows.add(Center(child: StaticRepere(element.figure)));
  }

  void _handleFunctionGraphBlock(FunctionGraphBlock element) {
    // start a new row
    _flushCurrentRow();

    rows.add(Center(child: FunctionGraphW(element)));
  }

  void _handleTableBlock(TableBlock element) {
    // start a new row
    _flushCurrentRow();

    rows.add(Center(child: TableSimple(element)));
  }

  void _handleTableDoubleEntryBlock(TableDoubleEntryBlock element) {
    // start a new row
    _flushCurrentRow();

    rows.add(Center(child: TableDoubleEntry(element)));
  }

  void _handleNumberFieldBlock(NumberFieldBlock element) {
    final ct = _controllers[element.iD] as NumberController;
    _currentRow.add(WidgetSpan(child: NumberField(_color, ct)));
  }

  void _handleExpressionFieldBlock(ExpressionFieldBlock element) {
    final ct = _controllers[element.iD] as ExpressionController;

    final field = WidgetSpan(
        child: ExpressionField(_color, ct, () => onFieldDone(element.iD)));
    if (element.label.isNotEmpty) {
      // start a new line
      _flushCurrentRow();

      rows.add(
        Center(
          child: Text.rich(
            TextSpan(
              children: [
                _inlineMath(element.label, fontSize),
                const TextSpan(text: " "),
                field,
              ],
            ),
          ),
        ),
      );
    } else {
      // just add the field in the current row
      _currentRow.add(field);
    }
  }

  void _handleRadioFieldBlock(RadioFieldBlock element) {
    final ct = _controllers[element.iD] as RadioController;

    // start a new line
    _flushCurrentRow();

    rows.add(RadioField(_color, ct));
  }

  void _handleDropDownFieldBlock(DropDownFieldBlock element) {
    final ct = _controllers[element.iD] as DropDownController;

    // add inline
    _currentRow.add(WidgetSpan(child: DropDownField(_color, ct)));
  }

  void _handleOrderedListFieldBlock(OrderedListFieldBlock element) {
    final ct = _controllers[element.iD] as OrderedListController;

    // start a new line
    _flushCurrentRow();

    rows.add(OrderedListField(_color, ct));
  }

  void _handleFigurePointFieldBlock(FigurePointFieldBlock element) {
    final ct = _controllers[element.iD] as FigurePointController;

    // start a new line
    _flushCurrentRow();

    final key = GlobalKey();
    zoomableWigets.add(key);
    rows.add(Center(child: FigurePointField(element.figure, ct, key: key)));
  }

  void _handleFigureVectorFieldBlock(FigureVectorFieldBlock element) {
    final ct = _controllers[element.iD] as FigureVectorController;

    // start a new line
    _flushCurrentRow();

    final key = GlobalKey();
    zoomableWigets.add(key);
    rows.add(Center(child: FigureVectorField(ct, key: key)));
  }

  void _handleFigureVectorPairFieldBlock(FigureVectorPairFieldBlock element) {
    final ct = _controllers[element.iD] as FigureVectorPairController;

    // start a new line
    _flushCurrentRow();

    final key = GlobalKey();
    zoomableWigets.add(key);
    rows.add(Center(child: FigureVectorPairField(ct, key: key)));
  }

  void _handleVariationTableFieldBlock(VariationTableFieldBlock element) {
    final ct = _controllers[element.iD] as VariationTableController;

    // start a new line
    _flushCurrentRow();

    rows.add(Center(child: VariationTableField(_color, ct)));
  }

  void _handleFunctionPointsFieldBlock(FunctionPointsFieldBlock element) {
    final ct = _controllers[element.iD] as FunctionPointsController;

    // start a new line
    _flushCurrentRow();

    final key = GlobalKey();
    zoomableWigets.add(key);
    rows.add(Center(child: FunctionPoints(ct, key: key)));
  }

  void _handleTreeFieldBlock(TreeFieldBlock element) {
    final ct = _controllers[element.iD] as TreeController;

    // start a new line
    _flushCurrentRow();

    rows.add(Center(child: TreeField(_color, ct)));
  }

  void _handleTableFieldBlock(TableFieldBlock element) {
    final ct = _controllers[element.iD] as TableController;

    // start a new line
    _flushCurrentRow();

    rows.add(Center(child: TableField(_color, ct)));
  }

  /// populate [rows]
  void build() {
    for (var element in _content) {
      // plain widgets

      if (element is TextBlock) {
        _handleTextBlock(element);
      } else if (element is FormulaBlock) {
        _handleFormulaBlock(element);
      } else if (element is VariationTableBlock) {
        _handleVariationTableBlock(element);
      } else if (element is SignTableBlock) {
        _handleSignTableBlock(element);
      } else if (element is FigureBlock) {
        _handleFigureBlock(element);
      } else if (element is FunctionGraphBlock) {
        _handleFunctionGraphBlock(element);
      } else if (element is TableBlock) {
        _handleTableBlock(element);
      } else if (element is TableDoubleEntryBlock) {
        _handleTableDoubleEntryBlock(element);

        // editable widgets

      } else if (element is NumberFieldBlock) {
        _handleNumberFieldBlock(element);
      } else if (element is ExpressionFieldBlock) {
        _handleExpressionFieldBlock(element);
      } else if (element is RadioFieldBlock) {
        _handleRadioFieldBlock(element);
      } else if (element is DropDownFieldBlock) {
        _handleDropDownFieldBlock(element);
      } else if (element is OrderedListFieldBlock) {
        _handleOrderedListFieldBlock(element);
      } else if (element is FigurePointFieldBlock) {
        _handleFigurePointFieldBlock(element);
      } else if (element is FigureVectorFieldBlock) {
        _handleFigureVectorFieldBlock(element);
      } else if (element is FigureVectorPairFieldBlock) {
        _handleFigureVectorPairFieldBlock(element);
      } else if (element is VariationTableFieldBlock) {
        _handleVariationTableFieldBlock(element);
      } else if (element is FunctionPointsFieldBlock) {
        _handleFunctionPointsFieldBlock(element);
      } else if (element is TreeFieldBlock) {
        _handleTreeFieldBlock(element);
      } else if (element is TableFieldBlock) {
        _handleTableFieldBlock(element);
      } else {
        // TODO:
      }

      lastIsText = element is TextBlock;
    }

    // flush the current row
    _flushCurrentRow();
  }
}

/// CheckQuestionSyntaxeNotification is emitted when the player
/// has edited one field
class CheckQuestionSyntaxeNotification extends Notification {
  final int id;
  final Answer answer;
  CheckQuestionSyntaxeNotification(this.id, this.answer);

  @override
  String toString() {
    return "CheckQuestionSyntaxeNotification($id, $answer)";
  }
}

/// ValidQuestionNotification is emitted when the player
/// validates his answer
class ValidQuestionNotification extends Notification {
  final QuestionAnswersIn data;
  ValidQuestionNotification(this.data);

  @override
  String toString() {
    return "ValidQuestionNotification($data)";
  }
}

class _OptionScrollList extends StatefulWidget {
  final _ContentBuilder content;
  final Widget button;

  const _OptionScrollList(this.content, this.button, {Key? key})
      : super(key: key);

  @override
  _OptionScrollListState createState() => _OptionScrollListState();
}

class _OptionScrollListState extends State<_OptionScrollList> {
  bool _isPaningOverView = false;

  /// returns true if [position] is in the widget identified
  /// by [key]
  bool _isInWidget(Offset position, GlobalKey key) {
    // find your widget
    final box = key.currentContext?.findRenderObject();
    if (box is! RenderBox) {
      return false;
    }

    //get offset
    Offset boxOffset = box.localToGlobal(Offset.zero);

    // check if your pointerdown event is inside the widget (you could do the same for the width, in this case I just used the height)
    if (position.dy > boxOffset.dy &&
        position.dy < boxOffset.dy + box.size.height) {
      // check x dimension aswell
      if (position.dx > boxOffset.dx &&
          position.dx < boxOffset.dx + box.size.width) {
        return true;
      }
    }

    return false;
  }

  void _checkPan(Offset position) {
    final hasZommable =
        widget.content.zoomableWigets.any((key) => _isInWidget(position, key));
    setState(() {
      _isPaningOverView = hasZommable;
    });
  }

  @override
  Widget build(BuildContext context) {
    return Listener(
      onPointerUp: (ev) {
        // restore the scroll posibility
        setState(() {
          _isPaningOverView = false;
        });
      },
      onPointerDown: (ev) {
        _checkPan(ev.position);
      },
      child: ListView(
          // if dragging over your widget, disable scroll, otherwise allow scrolling
          physics: _isPaningOverView
              ? const NeverScrollableScrollPhysics()
              : const ScrollPhysics(),
          shrinkWrap: true,
          children: [
            ...widget.content.rows
                .map(
                  (e) => Padding(
                      padding: const EdgeInsets.symmetric(vertical: 6.0),
                      child: e),
                )
                .toList(),
            const SizedBox(height: 10.0),
            widget.button
          ]),
    );
  }
}

class QuestionPage extends StatefulWidget {
  final Question question;
  final Color color;
  const QuestionPage(this.question, this.color, {Key? key}) : super(key: key);

  static Widget withEvents(
    void Function(CheckQuestionSyntaxeNotification) onCheckSyntax,
    void Function(ValidQuestionNotification) onValid,
    Question question,
    Color color,
  ) {
    return NotificationListener<CheckQuestionSyntaxeNotification>(
      onNotification: (v) {
        onCheckSyntax(v);
        return true;
      },
      child: NotificationListener<ValidQuestionNotification>(
        onNotification: (v) {
          onValid(v);
          return true;
        },
        child: QuestionPage(question, color),
      ),
    );
  }

  @override
  State<QuestionPage> createState() => _QuestionPageState();
}

class _QuestionPageState extends State<QuestionPage> {
  late Map<int, FieldController> _controllers;

  _ContentBuilder? builder;

  @override
  void initState() {
    _controllers = _ContentBuilder.initControllers(widget.question.enonce, () {
      setState(() {});
    });

    builder = _ContentBuilder(
        _emitCheckSyntax, widget.question.enonce, _controllers, widget.color);
    builder!.build();

    super.initState();
  }

  bool get areAnswersValid =>
      _controllers.values.every((ct) => ct.hasValidData());

  ValidQuestionNotification answers() {
    return ValidQuestionNotification(
      QuestionAnswersIn(
        _controllers.map((key, ct) => MapEntry(key, ct.getData())),
      ),
    );
  }

  void _emitCheckSyntax(int id) {
    final ct = _controllers[id]!;
    if (!ct.hasValidData()) {
      return;
    }
    CheckQuestionSyntaxeNotification(id, ct.getData()).dispatch(context);
  }

  @override
  Widget build(BuildContext context) {
    final shadows = [
      Shadow(
          color: widget.color.withOpacity(0.9),
          offset: const Offset(2, -2),
          blurRadius: 1.3)
    ];
    const spacing = SizedBox(height: 20.0);

    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 5.0),
      child: Column(
        mainAxisAlignment: MainAxisAlignment.spaceEvenly,
        children: [
          Text(
            "Thème test",
            style: TextStyle(fontSize: 22, shadows: shadows),
          ),
          Padding(
            padding: const EdgeInsets.symmetric(vertical: 12.0),
            child: Align(
              alignment: Alignment.centerLeft,
              child: Text(
                widget.question.title,
                style: TextStyle(
                  shadows: shadows,
                  fontSize: 20,
                ),
              ),
            ),
          ),
          if (builder != null)
            Expanded(
                child: _OptionScrollList(
              builder!,
              ElevatedButton(
                onPressed:
                    areAnswersValid ? () => answers().dispatch(context) : null,
                style: ElevatedButton.styleFrom(primary: widget.color),
                child: const Text(
                  "Valider",
                  style: TextStyle(fontSize: 18),
                ),
              ),
            )),
          spacing,
          TimeoutBar(const Duration(seconds: 60), widget.color),
          spacing,
          const Text("", style: TextStyle(fontSize: 16)),
        ],
      ),
    );
  }
}
