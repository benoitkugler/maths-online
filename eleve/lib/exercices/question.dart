import 'package:eleve/exercices/types.gen.dart';
import 'package:eleve/trivialpoursuit/events.gen.dart' as events;
import 'package:eleve/trivialpoursuit/pie.dart';
import 'package:eleve/trivialpoursuit/timeout_bar.dart';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_math_fork/flutter_math.dart';

abstract class _FieldController {
  /// returns true if the field is not empty and contains valid data
  bool hasValidData();

  /// returns the current answer
  Answer getData();
}

class _NumberController extends _FieldController {
  final TextEditingController textController;

  _NumberController(void Function() onChange)
      : textController = TextEditingController() {
    textController.addListener(onChange);
  }

  @override
  bool hasValidData() {
    final content = textController.text.trim();
    if (content.isEmpty) {
      return false;
    }
    return double.tryParse(content) != null;
  }

  @override
  Answer getData() {
    final content = textController.text.trim();
    return NumberAnswer(double.parse(content));
  }
}

/// utility class used to layout the Block
class _ContentBuilder {
  final List<Block> _content;
  final Color _color;

  /// field controllers created by [initControllers]
  final Map<int, _FieldController> _controllers;

  final List<Widget> rows = []; // final output

  List<InlineSpan> _currentRow = []; // current row
  static const _textStyle = TextStyle(fontSize: 18);

  _ContentBuilder(this._content, this._controllers, this._color);

  /// walks throught the question content and creates field controllers,
  /// later used when building widgets
  static Map<int, _FieldController> initControllers(
      List<Block> content, void Function() onChange) {
    final controllers = <int, _FieldController>{};
    for (var block in content) {
      if (block is NumberFieldBlock) {
        controllers[block.iD] = _NumberController(onChange);
      } // TODO: handle more fields
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
    _currentRow.add(TextSpan(
      text: element.text,
    ));
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

  void _handleFormulaBlock(FormulaBlock element) {
    if (element.isInline) {
      _currentRow.add(WidgetSpan(
        baseline: TextBaseline.alphabetic,
        alignment: PlaceholderAlignment.baseline,
        child: Math.tex(
          element.content,
          mathStyle: MathStyle.text,
          textStyle: _textStyle,
        ),
      ));
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
    final ct = _controllers[element.iD] as _NumberController;
    _currentRow.add(WidgetSpan(child: _NumberField(_color, ct.textController)));
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
      } else {
        // TODO:
      }
    }

    // flush the current row
    _flushCurrentRow();
  }
}

/// ValidQuestionNotification is emitted when the player
/// validates his answer
class ValidQuestionNotification extends Notification {
  final Map<int, Answer> answers;
  ValidQuestionNotification(this.answers);

  @override
  String toString() {
    return "ValidQuestionNotification($answers)";
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
  late Map<int, _FieldController> _controllers;

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
        _controllers.map((key, ct) => MapEntry(key, ct.getData())));
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

    final builder = _ContentBuilder(
        widget.question.enonce, _controllers, widget.categorie.color);
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

class _NumberField extends StatelessWidget {
  final Color _color;
  final TextEditingController _controller;

  const _NumberField(this._color, this._controller, {Key? key})
      : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 4),
      child: SizedBox(
        width: 100,
        child: TextField(
          controller: _controller,
          decoration: InputDecoration(
            isDense: true,
            contentPadding: const EdgeInsets.only(top: 10, bottom: 4),
            focusedBorder: UnderlineInputBorder(
              borderSide: BorderSide(
                color: _color,
              ),
            ),
          ),
          cursorColor: _color,
          style: TextStyle(color: Colors.yellow.shade100),
          textAlign: TextAlign.center,
          textAlignVertical: TextAlignVertical.center,
          keyboardType: const TextInputType.numberWithOptions(
              signed: true, decimal: true),
          inputFormatters: <TextInputFormatter>[
            FilteringTextInputFormatter.allow(RegExp(r'[0-9.-]'))
          ], // Only numbers, minus and dot can be entered
        ),
      ),
    );
  }
}
