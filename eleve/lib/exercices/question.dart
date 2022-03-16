import 'package:eleve/exercices/expression.dart';
import 'package:eleve/exercices/fields.dart';
import 'package:eleve/exercices/number.dart';
import 'package:eleve/exercices/radio.dart';
import 'package:eleve/exercices/types.gen.dart';
import 'package:eleve/trivialpoursuit/events.gen.dart' as events;
import 'package:eleve/trivialpoursuit/pie.dart';
import 'package:eleve/trivialpoursuit/timeout_bar.dart';
import 'package:flutter/material.dart';
import 'package:flutter_math_fork/flutter_math.dart';

/// utility class used to layout the Block
class _ContentBuilder {
  final void Function(int) onFieldDone;
  final List<Block> _content;
  final Color _color;

  /// field controllers created by [initControllers]
  final Map<int, FieldController> _controllers;

  final List<Widget> rows = []; // final output

  List<InlineSpan> _currentRow = []; // current row
  static const _textStyle = TextStyle(fontSize: 18);

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
      }

      // TODO: handle more fields
    }
    return controllers;
  }

  // we use Wrap instead of Rows to avoid overflows
  void _flushCurrentRow() {
    if (_currentRow.isEmpty) {
      return;
    }

    rows.add(
      Padding(
        padding: const EdgeInsets.symmetric(vertical: 6),
        child: Text.rich(
          TextSpan(
              style: _textStyle.copyWith(height: 1.5), children: _currentRow),
        ),
      ),
    );
    _currentRow = [];
  }

  void _handleTextBlock(TextBlock element) {
    if (element.text.endsWith("\n")) {
      _currentRow.add(TextSpan(
        text: element.text.substring(0, element.text.length - 1),
      ));
      // add an explicit new line
      _flushCurrentRow();
    } else {
      _currentRow.add(TextSpan(
        text: element.text,
      ));
    }
  }

  // void _handleInlineFormulaBLock(String content) {
  //   Math math = Math.tex(
  //     content,
  //     mathStyle: MathStyle.text,
  //     textStyle: _textStyle,
  //   );
  //   List<Math> parts = math.texBreak().parts;
  //   List<InlineSpan> children = [];
  //   for (Math part in parts) {
  //     children.add(WidgetSpan(
  //       baseline: TextBaseline.alphabetic,
  //       alignment: PlaceholderAlignment.baseline,
  //       child: part,
  //     ));
  //     children.add(const TextSpan(text: ' '));
  //   }
  //   children.removeLast();

  //   _currentRow.add(Text.rich(TextSpan(
  //     children: children,
  //   )));
  // }

  WidgetSpan _inlineMath(String content) {
    return WidgetSpan(
        baseline: TextBaseline.alphabetic,
        alignment: PlaceholderAlignment.baseline,
        child: Math.tex(
          content,
          mathStyle: MathStyle.text,
          textScaleFactor: 1.15,
          textStyle: _textStyle.copyWith(fontSize: _textStyle.fontSize! - 1),
        ));
  }

  void _handleFormulaBlock(FormulaBlock element) {
    if (element.isInline) {
      _currentRow.add(_inlineMath(element.content));
    } else {
      // start a new row
      _flushCurrentRow();

      rows.add(Center(
          child: Math.tex(
        element.content,
        mathStyle: MathStyle.display,
        textStyle: _textStyle,
      )));
    }
  }

  void _handleNumberFieldBlock(NumberFieldBlock element) {
    final ct = _controllers[element.iD] as NumberController;
    _currentRow.add(WidgetSpan(
        child: NumberField(
            _color, ct.textController, () => onFieldDone(element.iD))));
  }

  void _handleExpressionFieldBlock(ExpressionFieldBlock element) {
    final ct = _controllers[element.iD] as ExpressionController;

    final field = WidgetSpan(
        child: ExpressionField(
            _color, ct.textController, () => onFieldDone(element.iD)));
    if (element.label.isNotEmpty) {
      // start a new line
      _flushCurrentRow();

      rows.add(
        Center(
          child: Text.rich(
            TextSpan(
              children: [
                _inlineMath(element.label),
                const TextSpan(text: " "),
                _inlineMath("="),
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

  /// populate [rows]
  void build() {
    for (var element in _content) {
      if (element is TextBlock) {
        _handleTextBlock(element);
      } else if (element is FormulaBlock) {
        _handleFormulaBlock(element);
      } else if (element is NumberFieldBlock) {
        _handleNumberFieldBlock(element);
      } else if (element is ExpressionFieldBlock) {
        _handleExpressionFieldBlock(element);
      } else if (element is RadioFieldBlock) {
        _handleRadioFieldBlock(element);
      } else {
        // TODO:
      }
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

class QuestionPage extends StatefulWidget {
  final Question question;
  final events.Categorie categorie;
  const QuestionPage(this.question, this.categorie, {Key? key})
      : super(key: key);

  @override
  State<QuestionPage> createState() => _QuestionPageState();
}

class _QuestionPageState extends State<QuestionPage> {
  late Map<int, FieldController> _controllers;

  @override
  void initState() {
    _controllers = _ContentBuilder.initControllers(widget.question.enonce, () {
      setState(() {});
    });
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
          color: widget.categorie.color.withOpacity(0.9),
          offset: const Offset(2, -2),
          blurRadius: 1.3)
    ];
    const spacing = SizedBox(height: 20.0);

    final builder = _ContentBuilder(_emitCheckSyntax, widget.question.enonce,
        _controllers, widget.categorie.color);
    builder.build();

    return Column(
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
        spacing,
        Expanded(
          child: ListView(
            shrinkWrap: true,
            children: builder.rows
                .map(
                  (e) => Padding(
                      padding: const EdgeInsets.symmetric(vertical: 6.0),
                      child: e),
                )
                .toList(),
          ),
        ),
        spacing,
        ElevatedButton(
          onPressed: areAnswersValid ? () => answers().dispatch(context) : null,
          style: ElevatedButton.styleFrom(primary: widget.categorie.color),
          child: const Text(
            "Valider",
            style: TextStyle(fontSize: 18),
          ),
        ),
        spacing,
        TimeoutBar(const Duration(seconds: 60), widget.categorie.color),
        spacing,
        const Text("", style: TextStyle(fontSize: 16)),
      ],
    );
  }
}
