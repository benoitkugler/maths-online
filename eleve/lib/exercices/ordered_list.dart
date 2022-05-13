import 'package:eleve/exercices/fields.dart';
import 'package:eleve/exercices/types.gen.dart';
import 'package:flutter/material.dart';

class _IndexedProposal {
  // as send by the server, it is also the
  // indice into the _references proposal list
  final int index;
  final TextLine text;
  _IndexedProposal(this.index, this.text);
}

class OrderedListController extends FieldController {
  final List<_IndexedProposal> _references = [];
  final String label;
  final int expectedLength;

  List<_IndexedProposal> answers = [];

  List<_IndexedProposal> proposals = [];

  OrderedListController(void Function() onChange, OrderedListFieldBlock field)
      : expectedLength = field.answerLength,
        label = field.label,
        super(onChange) {
    for (var i = 0; i < field.proposals.length; i++) {
      _references.add(_IndexedProposal(i, field.proposals[i]));
    }
    // start with all propositions
    proposals = _references;
  }

  /// insertAnswerAt inserts [symbol] right after [location]
  /// in the answers array
  /// is is also removed from its old location if it was already in
  /// the answers
  void insertAnswerAt(_IndexedProposal symbol, int location) {
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
  void swapWithAnswer(_IndexedProposal symbol, int answerIndex) {
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
  void removeAnswer(_IndexedProposal symbol) {
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

  @override
  void setData(Answer answer) {
    final ans = (answer as OrderedListAnswer).indices;
    answers = ans.map((e) => _IndexedProposal(e, _references[e].text)).toList();
    // set the proposals to the remaining items
    final used = ans.toSet();
    proposals =
        _references.where((element) => !used.contains(element.index)).toList();
  }
}

class _PositionnedItem {
  final _IndexedProposal item;
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
      bool dense, bool enabled, List<_IndexedProposal> list) {
    return List<_Symbol>.generate(
        list.length,
        (index) =>
            _Symbol(dense, enabled, _PositionnedItem(list[index], index)));
  }

  @override
  Widget build(BuildContext context) {
    final text = Padding(
      padding: EdgeInsets.symmetric(
        vertical: 6,
        horizontal: isAnswer ? 8 : 12,
      ),
      child: TextRow(buildText(symbol.item.text, TextS(), fontSize),
          lineHeight: 1),
    );
    return Draggable<_PositionnedItem>(
      maxSimultaneousDrags: enabled ? null : 0,
      data: symbol,
      feedback: Material(
        elevation: 8,
        borderRadius: BorderRadius.circular(10),
        child: Padding(
          padding: const EdgeInsets.all(8),
          child: TextRow(buildText(symbol.item.text, TextS(), fontSize),
              lineHeight: 1),
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
                  child: text,
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
  final List<_IndexedProposal> answers;
  final String label; // optional

  final void Function(_PositionnedItem, bool isStart) addAnswer;

  const _AnswerRow(
      this.color, this.enabled, this.answers, this.label, this.addAnswer,
      {Key? key})
      : super(key: key);

  @override
  Widget build(BuildContext context) {
    const minHeight = 80.0;
    final widgets = answers.isEmpty
        ? [
            Expanded(
              child: DragTarget<_PositionnedItem>(
                  builder: (context, candidateData, rejectedData) => Container(
                        constraints: const BoxConstraints(
                          minHeight: minHeight,
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
                    constraints: const BoxConstraints(
                      minHeight: minHeight,
                    ),
                    color: candidateData.isEmpty ? null : color,
                  ),
                ),
                onAccept: (_PositionnedItem symbol) {
                  addAnswer(symbol, true);
                },
              ),
            ),
            ConstrainedBox(
              constraints: BoxConstraints(
                  maxWidth: MediaQuery.of(context).size.width * 2 / 3),
              child: Wrap(
                alignment: WrapAlignment.center,
                children: _Symbol.fromList(true, enabled, answers),
              ),
            ),
            Expanded(
              child: DragTarget<_PositionnedItem>(
                builder: (context, candidateData, rejectedData) => Padding(
                  padding: const EdgeInsets.only(left: 8.0),
                  child: Container(
                    constraints: const BoxConstraints(
                      minHeight: minHeight,
                    ),
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
  final List<_IndexedProposal> props;

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
