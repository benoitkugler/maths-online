import 'package:eleve/build_mode.dart';
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

/// [QuestionController] is the controller for a whole
/// question (internally, it is composed of multiple controllers, one for each field)
/// It should be created when displaying a question, then manipulated inside the [setState]
/// method of a [StatefulWidget].
class QuestionController {
  /// If [blockOnSubmit] is true, the validate button is disabled
  /// after one validation
  final bool blockOnSubmit;

  final Map<int, FieldController> _fields = {};
  bool _hasAnswered = false;

  bool get hasAnswered => _hasAnswered;

  void answer() {
    _hasAnswered = true;
    final enableFields = !blockOnSubmit || !_hasAnswered;
    if (!enableFields) {
      for (var element in _fields.values) {
        element.disable();
      }
    }
  }

  bool get enableValidate {
    final areAnswersValid =
        _fields.values.every((ct) => !ct.syntaxError && ct.hasValidData());
    return (!blockOnSubmit || !hasAnswered) && areAnswersValid;
  }

  QuestionAnswersIn answsers() {
    return QuestionAnswersIn(
      _fields.map((key, ct) => MapEntry(key, ct.getData())),
    );
  }

  void setAnswers(QuestionAnswersIn answers) {
    _fields.forEach((key, value) {
      _fields[key]!.setData(answers.data[key]!);
    });
  }

  /// Walks throught the question content and creates the field controllers,
  /// later used when building widgets.
  /// [onEditDone] is called when one field is updated by the user
  /// It should in return trigger a setState on the widget using the controller.
  /// [buildMode] is used to select the correct server endpoint
  /// when validating expression syntax
  QuestionController(List<Block> enonce, BuildMode buildMode,
      this.blockOnSubmit, void Function(int fieldID) onEditDone) {
    for (var block in enonce) {
      if (block is NumberFieldBlock) {
        _fields[block.iD] = NumberController(() => onEditDone(block.iD));
      } else if (block is ExpressionFieldBlock) {
        _fields[block.iD] =
            ExpressionController(buildMode, () => onEditDone(block.iD));
      } else if (block is RadioFieldBlock) {
        _fields[block.iD] =
            RadioController(() => onEditDone(block.iD), block.proposals);
      } else if (block is DropDownFieldBlock) {
        _fields[block.iD] =
            DropDownController(() => onEditDone(block.iD), block.proposals);
      } else if (block is OrderedListFieldBlock) {
        _fields[block.iD] =
            OrderedListController(() => onEditDone(block.iD), block);
      } else if (block is FigurePointFieldBlock) {
        _fields[block.iD] = FigurePointController(() => onEditDone(block.iD));
      } else if (block is FigureVectorFieldBlock) {
        _fields[block.iD] =
            FigureVectorController(block, () => onEditDone(block.iD));
      } else if (block is FigureVectorPairFieldBlock) {
        _fields[block.iD] = FigureVectorPairController(
            block.figure, () => onEditDone(block.iD));
      } else if (block is VariationTableFieldBlock) {
        _fields[block.iD] = VariationTableController(
            buildMode, block, () => onEditDone(block.iD));
      } else if (block is FunctionPointsFieldBlock) {
        _fields[block.iD] =
            FunctionPointsController(block, () => onEditDone(block.iD));
      } else if (block is TreeFieldBlock) {
        _fields[block.iD] = TreeController(block, () => onEditDone(block.iD));
      } else if (block is TableFieldBlock) {
        _fields[block.iD] = TableController(block, () => onEditDone(block.iD));
      }
    }
  }
}

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

/// utility class used to layout the blocks
class _ContentBuilder {
  final List<Block> _content;
  final Color _color;

  final Map<int, FieldController> _controllers;

  final List<Widget> rows = []; // final output
  final List<GlobalKey> zoomableWigets = [];

  List<InlineSpan> _currentRow = []; // current row
  bool lastIsText = false; // used to insert new line between to text block
  static const fontSize = 18.0;

  _ContentBuilder(this._content, this._controllers, this._color);

  void _flushCurrentRow() {
    if (_currentRow.isEmpty) {
      return;
    }

    rows.add(TextRow(_currentRow, verticalPadding: 6));
    _currentRow = [];
  }

  void _handleTextBlock(TextBlock element) {
    if (lastIsText) {
      _flushCurrentRow();
    }
    _currentRow.addAll(
        buildText(element.parts, TextS.fromTextBlock(element), fontSize));
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

    rows.add(Center(child: TableW(element)));
  }

  void _handleNumberFieldBlock(NumberFieldBlock element) {
    final ct = _controllers[element.iD] as NumberController;
    _currentRow.add(WidgetSpan(child: NumberField(_color, ct)));
  }

  void _handleExpressionFieldBlock(ExpressionFieldBlock element) {
    final ct = _controllers[element.iD] as ExpressionController;

    final field = WidgetSpan(child: ExpressionField(_color, ct));
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
      }

      lastIsText = element is TextBlock;
    }

    // flush the current row
    _flushCurrentRow();
  }
}

/// CheckQuestionSyntaxeNotification is emitted when the player
/// has edited one field.
/// It should usually trigger a call to the backend to check early
/// if the syntax is correct, before doing the real correction of the answer.
class CheckQuestionSyntaxeNotification extends Notification {
  final QuestionSyntaxCheckIn data;
  CheckQuestionSyntaxeNotification(this.data);

  @override
  String toString() {
    return "CheckQuestionSyntaxeNotification($data)";
  }
}

/// ValidQuestionNotification is emitted when the player
/// validates his answer.
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

/// [QuestionW] is the widget used to display a question
/// The interactivity is handled internally; with the two
/// hooks [onCheckSyntax] and [onValid] provided as external API,
/// as well as the [syntaxFeedback] parameter
class QuestionW extends StatefulWidget {
  final BuildMode buildMode;

  final Question question;
  final Color color;
  final void Function(ValidQuestionNotification) onValid;

  /// [timeout] is the duration for the question
  /// It may be set to [null] to hide the timeout bar.
  final Duration? timeout;

  /// [footerQuote] is displayed at the bottom of the screen when the
  /// question has been answered
  final String footerQuote;

  /// If [blockOnSubmit] is true, the validate button is disabled
  /// after one validation
  final bool blockOnSubmit;

  /// If [answer] is provided, the question controllers and fields
  /// are filled using the answers given.
  final QuestionAnswersIn? answer;

  const QuestionW(this.buildMode, this.question, this.color, this.onValid,
      {Key? key,
      this.timeout = const Duration(seconds: 60),
      this.footerQuote = "",
      this.blockOnSubmit = true,
      this.answer})
      : super(key: key);

  @override
  State<QuestionW> createState() => _QuestionWState();
}

class _QuestionWState extends State<QuestionW> {
  late _ContentBuilder builder;
  late QuestionController controller;

  @override
  void initState() {
    _initController();
    _buildFields();

    super.initState();
  }

  @override
  void didUpdateWidget(QuestionW oldWidget) {
    _initController();
    _buildFields();

    super.didUpdateWidget(oldWidget);
  }

  void _initController() {
    controller = QuestionController(widget.question.enonce, widget.buildMode,
        widget.blockOnSubmit, _onEditDone);
    if (widget.answer != null) {
      controller.setAnswers(widget.answer!);
    }
  }

  void _buildFields() {
    builder = _ContentBuilder(
        widget.question.enonce, controller._fields, widget.color);
    builder.build();
  }

  void _onEditDone(int fieldID) async {
    final ct = controller._fields[fieldID]!;
    if (!ct.hasValidData()) {
      return;
    }
    setState(() {});
  }

  void _validate() {
    final answers = ValidQuestionNotification(controller.answsers());
    controller.answer();
    _buildFields();
    setState(() {});
    widget.onValid(answers);
  }

  @override
  Widget build(BuildContext context) {
    const spacing = SizedBox(height: 20.0);
    final timeout = widget.timeout;
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 5.0),
      child: Column(
        mainAxisAlignment: MainAxisAlignment.spaceEvenly,
        children: [
          Padding(
            padding: const EdgeInsets.symmetric(vertical: 8.0),
            child: ColoredTitle("Question", widget.color),
          ),
          Expanded(
              child: _OptionScrollList(
            builder,
            ElevatedButton(
              onPressed: controller.enableValidate ? _validate : null,
              style: ElevatedButton.styleFrom(primary: widget.color),
              child: Text(
                controller.blockOnSubmit && controller.hasAnswered
                    ? "En attente..."
                    : "Valider",
                style: const TextStyle(fontSize: 18),
              ),
            ),
          )),
          if (timeout != null) ...[
            spacing,
            TimeoutBar(timeout, widget.color),
            Padding(
              padding: const EdgeInsets.only(top: 6),
              child: AnimatedOpacity(
                duration: const Duration(seconds: 2),
                opacity: controller.hasAnswered ? 1 : 0,
                child: Text(controller.hasAnswered ? widget.footerQuote : "",
                    style: const TextStyle(
                        fontSize: 16, fontStyle: FontStyle.italic)),
              ),
            ),
          ]
        ],
      ),
    );
  }
}

class ColoredTitle extends StatelessWidget {
  final String title;
  final Color color;
  const ColoredTitle(this.title, this.color, {Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    final shadows = [
      Shadow(
          color: color.withOpacity(0.9),
          offset: const Offset(2, -2),
          blurRadius: 1.3)
    ];
    return DecoratedBox(
      decoration: BoxDecoration(
        border: Border.all(color: color),
        borderRadius: BorderRadius.circular(5),
      ),
      child: Padding(
        padding: const EdgeInsets.all(8.0),
        child: Text(
          title,
          style: TextStyle(fontSize: 22, shadows: shadows),
        ),
      ),
    );
  }
}
