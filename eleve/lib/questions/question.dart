import 'package:eleve/questions/dropdown.dart';
import 'package:eleve/questions/expression.dart';
import 'package:eleve/questions/fields.dart';
import 'package:eleve/questions/function_graph.dart';
import 'package:eleve/questions/function_points.dart';
import 'package:eleve/questions/geometric_construction.dart';
import 'package:eleve/questions/image.dart';
import 'package:eleve/questions/number.dart';
import 'package:eleve/questions/ordered_list.dart';
import 'package:eleve/questions/probas_tree.dart';
import 'package:eleve/questions/proof.dart';
import 'package:eleve/questions/radio.dart';
import 'package:eleve/questions/repere.dart';
import 'package:eleve/questions/sets.dart';
import 'package:eleve/questions/sign_table.dart';
import 'package:eleve/questions/sign_table_field.dart';
import 'package:eleve/questions/table.dart';
import 'package:eleve/questions/variation_table.dart';
import 'package:eleve/questions/variation_table_field.dart';
import 'package:eleve/questions/vector.dart';
import 'package:eleve/quotes.dart';
import 'package:eleve/shared/timeout_bar.dart';
import 'package:eleve/shared/title.dart';
import 'package:eleve/shared/zommables.dart';
import 'package:eleve/types/src_maths_questions_client.dart';
import 'package:flutter/material.dart';
import 'package:flutter_math_fork/flutter_math.dart';

/// Walks throught the question content and creates the field controllers,
/// later used when building widgets.
Map<int, FieldController> _createFieldControllers(
    Question question, void Function() onChange) {
  final fields = <int, FieldController>{};
  for (var block in question.enonce) {
    if (block is NumberFieldBlock) {
      fields[block.iD] = NumberController(onChange);
    } else if (block is ExpressionFieldBlock) {
      fields[block.iD] = ExpressionController(onChange,
          showFractionHelp: block.showFractionHelp);
    } else if (block is RadioFieldBlock) {
      fields[block.iD] = RadioController(onChange, block.proposals);
    } else if (block is DropDownFieldBlock) {
      fields[block.iD] = DropDownController(onChange, block.proposals);
    } else if (block is OrderedListFieldBlock) {
      fields[block.iD] = OrderedListController(onChange, block);
    } else if (block is GeometricConstructionFieldBlock) {
      fields[block.iD] =
          GeometricConstructionController.fromBlock(onChange, block);
    } else if (block is VariationTableFieldBlock) {
      fields[block.iD] = VariationTableController(block, onChange);
    } else if (block is SignTableFieldBlock) {
      fields[block.iD] = SignTableController(block, onChange);
    } else if (block is FunctionPointsFieldBlock) {
      fields[block.iD] = FunctionPointsController(block, onChange);
    } else if (block is TreeFieldBlock) {
      fields[block.iD] = TreeController(block, onChange);
    } else if (block is TableFieldBlock) {
      fields[block.iD] = TableController(block, onChange);
    } else if (block is VectorFieldBlock) {
      fields[block.iD] = VectorController(block, onChange);
    } else if (block is ProofFieldBlock) {
      fields[block.iD] = ProofController(block, onChange);
    } else if (block is SetFieldBlock) {
      fields[block.iD] = SetController(onChange, block.sets);
    }
  }
  return fields;
}

typedef QuestionAnswers = Map<int, Answer>;
typedef QuestionFeedback = Map<int, bool>;

/// [QuestionController] defines the parameters
/// of a question widget.
/// It provides convenience methods to control the state
/// of each question-internal controllers.
class QuestionController extends ChangeNotifier {
  /// associated to the question ennonce
  Map<int, FieldController> fields = {};

  String buttonLabel = "Valider";
  bool buttonEnabled = false;

  /// [timeout] is the duration for the question,
  /// displaying a timeout bar. It does not trigger any
  /// action when done.
  /// It may be set to [null] to hide the timeout bar.
  Duration? timeout;

  /// If not null, [footerQuote] is displayed at the bottom of the screen.
  /// An empty [QuoteData] may be provided to occupy the space but hide
  /// the text.
  QuoteData? footerQuote;

  QuestionController();

  factory QuestionController.fromQuestion(Question question) {
    final out = QuestionController();
    out.fields = _createFieldControllers(question, out._onFieldChange);
    return out;
  }

  /// [_onFieldChange] is called when one field is updated
  /// by the user.
  /// It enables the primary button
  /// if all the fields are valid.
  void _onFieldChange() {
    buttonEnabled = fields.values.every((field) => field.hasValidData());
    notifyListeners();
  }

  /// [answers] wraps all the fields answers,
  /// assuming all fields are valid.
  QuestionAnswersIn answers() {
    return QuestionAnswersIn(
        fields.map((key, ct) => MapEntry(key, ct.getData())));
  }

  QuestionFeedback feedback() {
    return fields.map((key, ct) => MapEntry(key, ct.hasError));
  }

  /// [setAnswers] updates the fields with the given answers,
  /// removing any existing feedback.
  void setAnswers(QuestionAnswers answers) {
    fields.forEach((key, field) => field.setData(answers[key]!));
    setFeedback(null);
  }

  /// If [feedback] is not null, [setFeedback] marks the fields with a false value
  /// as error, and disable all fields
  /// If is is null, it removes error indicator and enable them again.
  void setFeedback(QuestionFeedback? feedback) {
    fields.forEach((key, field) =>
        field.setError(feedback == null ? false : !(feedback[key] ?? false)));
    setFieldsEnabled(feedback == null);
  }

  /// [setFieldsEnabled] set the enabled property for all the question fields
  void setFieldsEnabled(bool enabled) {
    for (var field in fields.values) {
      field.setEnabled(enabled);
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

class _LaidoutFields {
  final List<Widget> rows; // final output
  final List<GlobalKey> zoomableKeys;
  _LaidoutFields(this.rows, this.zoomableKeys);
}

/// utility class used to layout the blocks
class _FieldsBuilder {
  final List<Block> _content;
  final Color _color;

  final Map<int, FieldController> fields;

  final List<Widget> rows = []; // final output
  final List<GlobalKey> zoomableKeys = [];

  List<InlineSpan> _currentRow = []; // current row
  bool lastIsText = false; // used to insert new line between to text block
  static const fontSize = 18.0;

  _FieldsBuilder(this._content, this.fields, this._color);

  static _LaidoutFields build(
      List<Block> content, Map<int, FieldController> fields, Color color) {
    final builder = _FieldsBuilder(content, fields, color);
    builder._build();
    return _LaidoutFields(builder.rows, builder.zoomableKeys);
  }

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
        child: SingleChildScrollView(
      scrollDirection: Axis.horizontal,
      child: Math.tex(
        element.formula,
        mathStyle: MathStyle.display,
        textStyle: const TextStyle(fontSize: fontSize),
      ),
    )));
  }

  void _handleVariationTableBlock(VariationTableBlock element) {
    // start a new row
    _flushCurrentRow();

    rows.add(Center(child: VariationTableW(element)));
  }

  void _handleSignTableBlock(SignTableBlock element) {
    // start a new row
    _flushCurrentRow();

    rows.add(Center(child: SignTableW(element)));
  }

  void _handleFigureBlock(FigureBlock element) {
    // start a new row
    _flushCurrentRow();

    rows.add(Center(child: StaticRepere(element.figure)));
  }

  void _handleFunctionsGraphBlock(FunctionsGraphBlock element) {
    // start a new row
    _flushCurrentRow();

    rows.add(Center(child: FunctionsGraphW(element)));
  }

  void _handleTableBlock(TableBlock element) {
    // start a new row
    _flushCurrentRow();

    rows.add(Center(child: TableW(element)));
  }

  void _handleImageBlock(ImageBlock element) {
    // start a new row
    _flushCurrentRow();

    rows.add(Center(child: ImageW(element)));
  }

  void _handleNumberFieldBlock(NumberFieldBlock element) {
    final ct = fields[element.iD] as NumberController;
    _currentRow.add(WidgetSpan(
        child: NumberFieldW(
      _color,
      ct,
      sizeHint: element.sizeHint,
      autofocus: true,
    )));
  }

  void _handleExpressionFieldBlock(ExpressionFieldBlock element) {
    final ct = fields[element.iD] as ExpressionController;

    final field = WidgetSpan(
        child: ExpressionFieldW(
      _color,
      ct,
      hintWidth: element.sizeHint,
    ));
    if (element.label.isNotEmpty || element.suffix.isNotEmpty) {
      // start a new line
      _flushCurrentRow();

      rows.add(
        Center(
          child: Text.rich(
            TextSpan(
              children: [
                if (element.label.isNotEmpty) ...[
                  _inlineMath(element.label, fontSize),
                  const TextSpan(text: " "),
                ],
                field,
                if (element.suffix.isNotEmpty) ...[
                  const TextSpan(text: " "),
                  _inlineMath(element.suffix, fontSize),
                ],
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
    final ct = fields[element.iD] as RadioController;

    // start a new line
    _flushCurrentRow();

    rows.add(RadioFieldW(_color, ct));
  }

  void _handleDropDownFieldBlock(DropDownFieldBlock element) {
    final ct = fields[element.iD] as DropDownController;

    // add inline
    _currentRow.add(WidgetSpan(child: DropDownFieldW(_color, ct)));
  }

  void _handleOrderedListFieldBlock(OrderedListFieldBlock element) {
    final ct = fields[element.iD] as OrderedListController;

    // start a new line
    _flushCurrentRow();

    rows.add(OrderedListFieldW(_color, ct));
  }

  void _handleGeometricConstructionFieldBlock(
      GeometricConstructionFieldBlock element) {
    final ct = fields[element.iD] as GeometricConstructionController;

    // start a new line
    _flushCurrentRow();

    final key = GlobalKey();
    final zoom = TransformationController();
    zoomableKeys.add(key);
    rows.add(Center(
        child: Zoomable(
            zoom, GeometricConstructionFieldW(element, ct, zoom), key)));
  }

  void _handleVariationTableFieldBlock(VariationTableFieldBlock element) {
    final ct = fields[element.iD] as VariationTableController;

    // start a new line
    _flushCurrentRow();

    rows.add(Center(child: VariationTableFieldW(_color, ct)));
  }

  void _handleSignTableFieldBlock(SignTableFieldBlock element) {
    final ct = fields[element.iD] as SignTableController;

    // start a new line
    _flushCurrentRow();

    rows.add(Center(child: SignTableFieldW(_color, ct)));
  }

  void _handleFunctionPointsFieldBlock(FunctionPointsFieldBlock element) {
    final ct = fields[element.iD] as FunctionPointsController;

    // start a new line
    _flushCurrentRow();

    final key = GlobalKey();
    final zoom = TransformationController();
    zoomableKeys.add(key);
    rows.add(Center(child: Zoomable(zoom, FunctionPointsW(ct, zoom), key)));
  }

  void _handleTreeBlock(TreeBlock element) {
    // start a new line
    _flushCurrentRow();

    rows.add(Center(child: TreeW(_color, element)));
  }

  void _handleTreeFieldBlock(TreeFieldBlock element) {
    final ct = fields[element.iD] as TreeController;

    // start a new line
    _flushCurrentRow();

    rows.add(Center(child: TreeFieldW(_color, ct)));
  }

  void _handleTableFieldBlock(TableFieldBlock element) {
    final ct = fields[element.iD] as TableController;

    // start a new line
    _flushCurrentRow();

    rows.add(Center(child: TableFieldW(_color, ct)));
  }

  void _handleVectorFieldBlock(VectorFieldBlock element) {
    final ct = fields[element.iD] as VectorController;
    _currentRow.add(WidgetSpan(
        child: VectorFieldW(_color, ct),
        alignment: element.displayColumn
            ? PlaceholderAlignment.middle
            : PlaceholderAlignment.bottom));
  }

  void _handleProofFieldBlock(ProofFieldBlock element) {
    final ct = fields[element.iD] as ProofController;

    // start a new line
    _flushCurrentRow();

    rows.add(ProofFieldW(_color, ct));
  }

  void _handleSetFieldBlock(SetFieldBlock element) {
    final ct = fields[element.iD] as SetController;

    // start a new line
    _flushCurrentRow();

    rows.add(SetFieldW(_color, ct));
  }

  /// populate [rows]
  void _build() {
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
      } else if (element is FunctionsGraphBlock) {
        _handleFunctionsGraphBlock(element);
      } else if (element is TableBlock) {
        _handleTableBlock(element);
      } else if (element is ImageBlock) {
        _handleImageBlock(element);

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
      } else if (element is GeometricConstructionFieldBlock) {
        _handleGeometricConstructionFieldBlock(element);
      } else if (element is VariationTableFieldBlock) {
        _handleVariationTableFieldBlock(element);
      } else if (element is SignTableFieldBlock) {
        _handleSignTableFieldBlock(element);
      } else if (element is FunctionPointsFieldBlock) {
        _handleFunctionPointsFieldBlock(element);
      } else if (element is TreeBlock) {
        _handleTreeBlock(element);
      } else if (element is TreeFieldBlock) {
        _handleTreeFieldBlock(element);
      } else if (element is TableFieldBlock) {
        _handleTableFieldBlock(element);
      } else if (element is VectorFieldBlock) {
        _handleVectorFieldBlock(element);
      } else if (element is ProofFieldBlock) {
        _handleProofFieldBlock(element);
      } else if (element is SetFieldBlock) {
        _handleSetFieldBlock(element);
      }

      lastIsText = element is TextBlock;
    }

    // flush the current row
    _flushCurrentRow();
  }
}

class _ListRows extends StatelessWidget {
  final List<Block> content;
  final Map<int, FieldController> fields;
  final Color color;

  final Widget? bottom;

  const _ListRows(this.content, this.fields, this.color, this.bottom);

  @override
  Widget build(BuildContext context) {
    final rows = _FieldsBuilder.build(content, fields, color);
    return ListWithZoomables([
      ...rows.rows.map(
        (e) => Padding(
            padding: const EdgeInsets.symmetric(vertical: 4.0), child: e),
      ),
      const SizedBox(height: 10.0),
      if (bottom != null) bottom!
    ], rows.zoomableKeys, shrinkWrap: true);
  }
}

/// [CorrectionView] displays the correction of a question
class CorrectionView extends StatelessWidget {
  final Enonce correction;
  final Color color;

  /// If not null, [footerQuote] is displayed at the bottom of the screen.
  final QuoteData? footerQuote;

  const CorrectionView(this.correction, this.color, this.footerQuote,
      {super.key});

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 5.0),
      child: Column(
        mainAxisAlignment: MainAxisAlignment.spaceEvenly,
        children: [
          Padding(
            padding: const EdgeInsets.symmetric(vertical: 8.0),
            child: ColoredTitle("Correction", color),
          ),
          Expanded(
              child: _ListRows(
            correction,
            const {},
            color,
            footerQuote == null
                ? null
                : AnimatedOpacity(
                    duration: const Duration(seconds: 2),
                    opacity: footerQuote!.isEmpty ? 0 : 1,
                    child: Quote(footerQuote!)),
          )),
        ],
      ),
    );
  }
}

/// [QuestionView] is the widget used to display a question.
class QuestionView extends StatefulWidget {
  final Question question;
  final QuestionController controller;

  /// [onButtonClick] is called when clicking on
  /// the primary button.
  final void Function() onButtonClick;

  final Color color;

  /// [title] is the title of the question
  final String title;

  final Widget? leadingButton;

  const QuestionView(
      this.question, this.controller, this.onButtonClick, this.color,
      {super.key, this.title = "Question", this.leadingButton});

  @override
  State<QuestionView> createState() => _QuestionViewState();
}

class _QuestionViewState extends State<QuestionView> {
  @override
  void initState() {
    _initController();

    super.initState();
  }

  @override
  void didUpdateWidget(QuestionView oldWidget) {
    _initController();

    super.didUpdateWidget(oldWidget);
  }

  @override
  void dispose() {
    widget.controller.removeListener(_onControllerChange);
    super.dispose();
  }

  void _initController() {
    widget.controller.removeListener(_onControllerChange);
    widget.controller.addListener(_onControllerChange);
  }

  void _onControllerChange() {
    if (mounted) setState(() {});
  }

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 5.0),
      child: Column(
        mainAxisAlignment: MainAxisAlignment.spaceEvenly,
        children: [
          Padding(
            padding: const EdgeInsets.symmetric(vertical: 8.0),
            child: ColoredTitle(widget.title, widget.color),
          ),
          if (widget.controller.timeout != null) ...[
            TimeoutBar(widget.controller.timeout!, widget.color),
            const SizedBox(height: 10.0),
          ],
          Expanded(
              child: _ListRows(
                  widget.question.enonce,
                  widget.controller.fields,
                  widget.color,
                  Row(
                    children: [
                      if (widget.leadingButton != null) widget.leadingButton!,
                      Expanded(
                        child: OutlinedButton(
                          onPressed: widget.controller.buttonEnabled
                              ? widget.onButtonClick
                              : null,
                          style: OutlinedButton.styleFrom(
                              foregroundColor: widget.color),
                          child: Text(
                            widget.controller.buttonLabel,
                            style: const TextStyle(fontSize: 18),
                          ),
                        ),
                      ),
                    ],
                  ))),
          if (widget.controller.footerQuote != null)
            Padding(
              padding: const EdgeInsets.only(top: 6),
              child: AnimatedOpacity(
                  duration: const Duration(seconds: 2),
                  opacity: widget.controller.footerQuote!.isEmpty ? 0 : 1,
                  child: Quote(widget.controller.footerQuote!)),
            ),
        ],
      ),
    );
  }
}
