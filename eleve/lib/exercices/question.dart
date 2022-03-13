import 'package:eleve/exercices/types.gen.dart';
import 'package:eleve/trivialpoursuit/events.gen.dart';
import 'package:eleve/trivialpoursuit/pie.dart';
import 'package:eleve/trivialpoursuit/timeout_bar.dart';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_math_fork/flutter_math.dart';

/// utility class used to layout the ClientBlock
class _ContentBuilder {
  final List<ClientBlock> _content;
  final Color _color;

  final List<Widget> rows = []; // final output

  List<InlineSpan> _currentRow = []; // current row
  static const _textStyle = TextStyle(fontSize: 18);

  _ContentBuilder(this._content, this._color);

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

  void _handleTextBlock(ClientTextBlock element) {
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

  void _handleFormulaBlock(ClientFormulaBlock element) {
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

  void _handleNumberFieldBlock(ClientNumberFieldBlock element) {
    // TODO: stores controllers by ID
    final ct = TextEditingController();
    _currentRow.add(WidgetSpan(child: _NumberField(_color, ct)));
  }

  /// populate [rows]
  void build() {
    for (var element in _content) {
      if (element is ClientTextBlock) {
        _handleTextBlock(element);
      } else if (element is ClientFormulaBlock) {
        _handleFormulaBlock(element);
      } else if (element is ClientNumberFieldBlock) {
        _handleNumberFieldBlock(element);
      } else {
        // TODO:
      }
    }

    // flush the current row
    _flushCurrentRow();
  }
}

class Question extends StatelessWidget {
  final ClientQuestion question;
  final Categorie categorie;
  const Question(this.question, this.categorie, {Key? key}) : super(key: key);

  // group the text and inline math block on the same row
  List<Widget> _buildColumns() {
    final builder = _ContentBuilder(question.content, categorie.color);
    builder.build();
    return builder.rows;
  }

  @override
  Widget build(BuildContext context) {
    final shadows = [
      Shadow(
          color: categorie.color.withOpacity(0.9),
          offset: const Offset(2, -2),
          blurRadius: 1.3)
    ];
    const spacing = SizedBox(height: 20.0);
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
              question.title,
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
            children: _buildColumns()
                .map(
                  (e) => Padding(
                      padding: const EdgeInsets.symmetric(vertical: 6.0),
                      child: e),
                )
                .toList(),
          ),
        ),
        spacing,
        TextField(
          cursorColor: categorie.color,
          decoration: InputDecoration(
            focusedBorder: OutlineInputBorder(
                borderSide: BorderSide(
              color: categorie.color,
            )),
            border: const OutlineInputBorder(),
          ),
        ),
        spacing,
        ElevatedButton(
          onPressed: () => print("ok"),
          style: ElevatedButton.styleFrom(primary: categorie.color),
          child: const Text(
            "Valider",
            style: TextStyle(fontSize: 18),
          ),
        ),
        spacing,
        TimeoutBar(const Duration(seconds: 60), categorie.color),
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
