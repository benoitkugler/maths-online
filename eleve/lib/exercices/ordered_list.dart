import 'package:eleve/exercices/fields.dart';
import 'package:eleve/exercices/types.gen.dart';
import 'package:flutter/material.dart';

class IndexedText {
  // as send by the server, it is also the
  // indice into the reference proposal list
  final int index;
  final String text;
  IndexedText(this.index, this.text);
}

class OrderedListController extends FieldController {
  final List<IndexedText> _references = [];
  final String label;
  final int expectedLength;

  List<IndexedText> answers = [];

  List<IndexedText> proposals = [];

  OrderedListController(void Function() onChange, OrderedListFieldBlock field)
      : expectedLength = field.answerLength,
        label = field.label,
        super(onChange) {
    for (var i = 0; i < field.proposals.length; i++) {
      _references.add(IndexedText(i, field.proposals[i]));
    }
    // start with all propositions
    proposals = _references;
  }

  /// insertAnswerAt inserts [symbol] right after [location]
  /// in the answers array
  /// is is also removed from its old location if it was already in
  /// the answers
  void insertAnswerAt(IndexedText symbol, int location) {
    proposals.removeWhere((element) => element.index == symbol.index);
    final existing =
        answers.indexWhere((element) => element.index == symbol.index);
    if (existing != -1) {
      // remove the current symbol from answers
      answers.removeAt(existing);
      // adjust the new location
      if (location > existing) {
        location--;
      }
    }
    answers.insert(location, symbol);
    onChange();
  }

  /// swapWithAnswer adds [symbol] into the answers, removing the
  /// previous element at [answerIndex]
  void swapWithAnswer(IndexedText symbol, int answerIndex) {
    proposals.add(answers[answerIndex]);
    answers[answerIndex] = symbol;
    onChange();
  }

  /// swapBetweenAnswers swap the orders of symbols at indices
  /// [answerIndex1] and [answerIndex2]
  void swapBetweenAnswers(int answerIndex1, int answerIndex2) {
    final tmp = answers[answerIndex1];
    answers[answerIndex1] = answers[answerIndex2];
    answers[answerIndex2] = tmp;
    onChange();
  }

  /// remove [symbol] from the chosen answer and put it back
  /// in the proposals
  void removeAnswer(IndexedText symbol) {
    answers.removeWhere((element) => element.index == symbol.index);
    proposals.add(symbol);
    onChange();
  }

  @override
  bool hasValidData() {
    return answers.length == expectedLength;
  }

  @override
  Answer getData() {
    return OrderedListAnswer(answers.map((item) => item.index).toList());
  }
}

class _PositionnedItem {
  final IndexedText item;
  final int position;
  _PositionnedItem(this.item, this.position);
}

class OrderedListField extends StatefulWidget {
  final Color _color;
  final OrderedListController _controller;

  const OrderedListField(this._color, this._controller, {Key? key})
      : super(key: key);

  @override
  _OrderedListFieldState createState() => _OrderedListFieldState();
}

class _OrderedListFieldState extends State<OrderedListField> {
  @override
  Widget build(BuildContext context) {
    final ct = widget._controller;
    final props = ct.proposals;
    final answers = ct.answers;
    return Column(
      children: [
        _AnswerRow(
            widget._color,
            ct.enabled,
            answers,
            ct.label,
            (symbol, isStart) => setState(() {
                  // insert into answers
                  ct.insertAnswerAt(symbol.item, isStart ? 0 : answers.length);
                })),
        const SizedBox(height: 20),
        _PropsRow(widget._color, ct.enabled, props, (symbol) {
          setState(() {
            ct.removeAnswer(symbol.item);
          });
        })
      ],
    );
  }
}

class _Symbol extends StatelessWidget {
  final bool isAnswer;
  final bool enabled;
  final _PositionnedItem symbol;
  const _Symbol(this.isAnswer, this.enabled, this.symbol, {Key? key})
      : super(key: key);

  static const fontSize = 16.0;

  static List<_Symbol> fromList(
      bool dense, bool enabled, List<IndexedText> list) {
    return List<_Symbol>.generate(
        list.length,
        (index) =>
            _Symbol(dense, enabled, _PositionnedItem(list[index], index)));
  }

  @override
  Widget build(BuildContext context) {
    final text = Padding(
      padding: EdgeInsets.symmetric(
        vertical: 8,
        horizontal: isAnswer ? 8 : 12,
      ),
      child: textMath(symbol.item.text, fontSize),
    );
    return Draggable<_PositionnedItem>(
      maxSimultaneousDrags: enabled ? null : 0,
      data: symbol,
      feedback: Material(
        elevation: 8,
        borderRadius: BorderRadius.circular(10),
        child: Padding(
          padding: const EdgeInsets.all(8),
          child: textMath(symbol.item.text, fontSize),
        ),
      ),
      child: Padding(
        padding: EdgeInsets.symmetric(horizontal: isAnswer ? 0 : 6),
        child: Material(
          elevation: 8,
          borderRadius: BorderRadius.circular(isAnswer ? 1 : 10),
          child: isAnswer
              ? SizedBox(
                  height: 35,
                  child: Center(child: text),
                )
              : text,
        ),
      ),
    );
  }
}

class _AnswerRow extends StatelessWidget {
  final Color color;
  final bool enabled;
  final List<IndexedText> answers;
  final String label; // optional

  final void Function(_PositionnedItem, bool isStart) addAnswer;

  const _AnswerRow(
      this.color, this.enabled, this.answers, this.label, this.addAnswer,
      {Key? key})
      : super(key: key);

  @override
  Widget build(BuildContext context) {
    final widgets = answers.isEmpty
        ? [
            Expanded(
              child: DragTarget<_PositionnedItem>(
                  builder: (context, candidateData, rejectedData) => Container(
                        constraints: const BoxConstraints(
                          minHeight: 40,
                        ),
                        decoration: BoxDecoration(
                          color: candidateData.isEmpty ? null : color,
                        ),
                        child: const Center(
                          child: Text("Glisser les symboles...",
                              style: TextStyle(fontStyle: FontStyle.italic)),
                        ),
                      ),
                  onAccept: (_PositionnedItem symbol) {
                    addAnswer(symbol, true);
                  }),
            )
          ]
        : [
            Expanded(
              child: DragTarget<_PositionnedItem>(
                builder: (context, candidateData, rejectedData) => Padding(
                  padding: const EdgeInsets.only(right: 8.0),
                  child: Container(
                    constraints: const BoxConstraints(minHeight: 40),
                    color: candidateData.isEmpty ? null : color,
                  ),
                ),
                onAccept: (_PositionnedItem symbol) {
                  addAnswer(symbol, true);
                },
              ),
            ),
            ..._Symbol.fromList(true, enabled, answers),
            Expanded(
              child: DragTarget<_PositionnedItem>(
                builder: (context, candidateData, rejectedData) => Padding(
                  padding: const EdgeInsets.only(left: 8.0),
                  child: Container(
                    constraints: const BoxConstraints(minHeight: 40),
                    color: candidateData.isEmpty ? null : color,
                  ),
                ),
                onAccept: (_PositionnedItem symbol) {
                  addAnswer(symbol, false);
                },
              ),
            ),
          ];
    return Container(
        decoration: BoxDecoration(border: Border.all(color: color)),
        child: Row(children: [
          if (label.isNotEmpty)
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 8.0),
              child: textMath(label, _Symbol.fontSize),
            ),
          ...widgets,
        ]));
  }
}

class _PropsRow extends StatelessWidget {
  final Color color;
  final bool enabled;
  final List<IndexedText> props;

  final void Function(_PositionnedItem) removeAnswer;

  const _PropsRow(this.color, this.enabled, this.props, this.removeAnswer,
      {Key? key})
      : super(key: key);

  bool _isProposition(_PositionnedItem? candidate) {
    if (candidate == null) {
      return false;
    }
    return props
            .indexWhere((element) => element.index == candidate.item.index) !=
        -1;
  }

  @override
  Widget build(BuildContext context) {
    return DragTarget<_PositionnedItem>(
      builder: (context, candidateData, rejectedData) {
        bool accept =
            candidateData.isNotEmpty && !_isProposition(candidateData.first);
        return Container(
          padding: const EdgeInsets.all(4),
          decoration: BoxDecoration(
              border: Border.all(color: accept ? color : Colors.transparent)),
          child: Wrap(
            alignment: WrapAlignment.center,
            crossAxisAlignment: WrapCrossAlignment.center,
            children: _Symbol.fromList(false, enabled, props),
          ),
        );
      },
      onWillAccept: (_PositionnedItem? symbol) => !_isProposition(symbol),
      onAccept: removeAnswer,
    );
  }
}
