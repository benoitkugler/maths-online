import 'package:eleve/questions/fields.dart';
import 'package:eleve/types/src_maths_expression_sets.dart';
import 'package:eleve/types/src_maths_questions_client.dart';
import 'package:flutter/material.dart';

// same as ListNode, but mutable
class _Node {
  List<_Node> args;
  SetOp op;
  Set leaf;
  _Node(this.args, this.op, this.leaf);

  factory _Node.empty() => _Node([], SetOp.sUnion, 0);
  factory _Node.leaf(Set s) => _Node([], SetOp.sLeaf, s);
  factory _Node.comp(_Node kid) => _Node([kid], SetOp.sComplement, 0);

  ListNode toListNode() =>
      ListNode(args.map((e) => e.toListNode()).toList(), op, leaf);

  factory _Node.fromListNode(ListNode n) =>
      _Node(n.args.map((e) => _Node.fromListNode(e)).toList(), n.op, n.leaf);

  _Node shallowCopy() => _Node(args, op, leaf);
  _Node deepCopy() => _Node(args.map((e) => e.deepCopy()).toList(), op, leaf);

  bool get uOrI => op == SetOp.sUnion || op == SetOp.sInter;

  // return true if the node is a sentinel empty value
  bool isEmpty() => (op == SetOp.sUnion && args.isEmpty);

  bool isValid() {
    switch (op) {
      case SetOp.sUnion:
        return args.isNotEmpty && args.every((element) => element.isValid());
      case SetOp.sInter:
        return args.isNotEmpty && args.every((element) => element.isValid());
      case SetOp.sComplement:
        return args.length == 1 && args.first.isValid();
      case SetOp.sLeaf:
        return true;
    }
  }

  // if the last element of args is empty, replace
  void addOrReplace(_Node kid) {
    if (args.isNotEmpty && args.last.isEmpty()) {
      args[args.length - 1] = kid;
    } else {
      args.add(kid);
    }
  }

  // replace all fields by other
  void replace(_Node other) {
    args = other.args;
    op = other.op;
    leaf = other.leaf;
  }
}

class _IndexedSet {
  final Set index;
  final String set;
  const _IndexedSet(this.index, this.set);
}

class SetController extends FieldController {
  final List<_IndexedSet> _sets; // shuffled

  final List<_Node> _history = []; // history API

  _Node _answer = _Node.empty();
  // pointer to one node of answer,
  // may be modified by user tap
  late _Node _cursor;

  SetController(super.onChange, List<String> sets)
      : _sets = List.generate(
            sets.length, (index) => _IndexedSet(index, sets[index])).toList() {
    // start with focus on the root
    _cursor = _answer;
    _sets.shuffle();
  }

  @override
  Answer getData() {
    return SetAnswer(_answer.toListNode());
  }

  @override
  bool hasValidData() {
    return _answer.isValid();
  }

  @override
  void setData(Answer answer) {
    _answer = _Node.fromListNode((answer as SetAnswer).root);
    _cursor = _answer; // reset the cursor to the root
  }

  void _saveHistory() {
    _history.add(_answer.deepCopy());
  }

  void _restoreHistory() {
    if (_history.isEmpty) return;
    _answer = _history.last;
    _cursor = _answer;
    _history.removeLast();
  }

  void _clear() {
    _saveHistory();
    _answer = _Node.empty();
    _cursor = _answer;
  }

  // add the given leaf to the current expression
  void _addLeaf(Set s) {
    _saveHistory();
    if (_cursor.uOrI) {
      // extend the union or intersection
      _cursor.addOrReplace(_Node.leaf(s));
    } else if (_cursor.op == SetOp.sComplement) {
      // add or replace
      _cursor.args = [_Node.leaf(s)];
    } else {
      // nothing to do
    }
  }

  void _addParenthesis() {
    _saveHistory();
    if (_cursor.uOrI) {
      if (_cursor.args.isNotEmpty && _cursor.args.last.isEmpty()) {
        // change the cursor
        _cursor = _cursor.args.last;
      }
    } else if (_cursor.op == SetOp.sComplement) {
    } else if (_cursor.op == SetOp.sLeaf) {}
  }

  void _addOp(SetOp op) {
    _saveHistory();
    if (_cursor.uOrI) {
      if (_cursor.op == op) {
        // extend the union or intersection
        _cursor.args.add(_Node.empty());
      } else if (op == SetOp.sComplement) {
        if (_cursor.args.isEmpty) {
          // add a new node with comp and pass the focus
          final newCursor = _Node.empty();
          _cursor.args.add(_Node.comp(newCursor));
          _cursor = newCursor;
        } else if (_cursor.args.last.isEmpty()) {
          // pass the focus
          final newCursor = _Node.empty();
          _cursor.args[_cursor.args.length - 1] = _Node.comp(newCursor);
          _cursor = newCursor;
        } else {
          // apply comp to the current cursor
          _cursor.replace(_Node.comp(_cursor.shallowCopy()));
        }
      } else {
        // replace by the given op applied to this and empty
        // and update the cursor
        final newCursor = _Node.empty();
        _cursor.replace(_Node([_cursor.shallowCopy(), newCursor], op, 0));
        _cursor = newCursor;
      }
    } else if (_cursor.op == SetOp.sComplement) {
      if (op == _cursor.op) {
        // remove the comp
        _cursor.replace(_cursor.args.isEmpty ? _Node.empty() : _cursor.args[0]);
      } else {
        _cursor.replace(_Node([_cursor.shallowCopy(), _Node.empty()], op, 0));
      }
    } else if (_cursor.op == SetOp.sLeaf && op == SetOp.sComplement) {
      _cursor.replace(_Node.comp(_cursor.shallowCopy()));
    }
  }
}

class SetFieldW extends StatefulWidget {
  final Color color;
  final SetController controller;

  const SetFieldW(this.color, this.controller, {super.key});

  @override
  State<SetFieldW> createState() => _SetFieldWState();
}

class _SetFieldWState extends State<SetFieldW> {
  Offset? dragGesturePosition = Offset.zero;

  @override
  Widget build(BuildContext context) {
    final ct = widget.controller;
    const operators = [SetOp.sUnion, SetOp.sInter, SetOp.sComplement];
    return Container(
      padding: const EdgeInsets.all(4.0),
      decoration: BoxDecoration(
        border: Border.all(color: widget.color),
        borderRadius: BorderRadius.circular(5),
      ),
      child: Column(
        mainAxisSize: MainAxisSize.min,
        crossAxisAlignment: CrossAxisAlignment.center,
        children: [
          Row(
            children: [
              IconButton(
                onPressed: () => setState(() => ct._clear()),
                icon: const Icon(Icons.clear),
                splashRadius: 20,
              ),
              Expanded(
                child: Center(
                  child: SingleChildScrollView(
                    scrollDirection: Axis.horizontal,
                    child: _NodeW(
                        false, widget.controller._answer, ct._sets, ct._cursor,
                        onTap: (cr) => setState(() => ct._cursor = cr)),
                  ),
                ),
              ),
              IconButton(
                onPressed: ct._history.isEmpty
                    ? null
                    : () => setState(() => ct._restoreHistory()),
                icon: const Icon(Icons.undo),
                splashRadius: 20,
              )
            ],
          ),
          Wrap(
              children: ct._sets
                  .map(
                    (si) => _Control(
                      () => setState(() => ct._addLeaf(si.index)),
                      si.set,
                    ),
                  )
                  .toList()),
          Row(mainAxisSize: MainAxisSize.min, children: [
            _Control(() => setState(() => ct._addParenthesis()), '()'),
            ...operators.map((op) =>
                _Control(() => setState(() => ct._addOp(op)), op.latex()))
          ]),
          // DEBUG only
          // ElevatedButton(
          //     onPressed: widget.controller.hasValidData()
          //         ? () {
          //             setState(() {
          //               widget.controller
          //                   .setData(SetAnswer(_Node.empty().toListNode()));
          //             });
          //           }
          //         : null,
          //     child: const Text("Valid"))
        ],
      ),
    );
  }
}

extension on SetOp {
  String latex() {
    switch (this) {
      case SetOp.sUnion:
        return "\\cup";
      case SetOp.sInter:
        return "\\cap";
      case SetOp.sComplement:
        return "\\overline{ \\quad }";
      case SetOp.sLeaf:
        return "";
    }
  }
}

class _Control extends StatelessWidget {
  final void Function() onTap;
  final String latex;
  const _Control(this.onTap, this.latex, {super.key});

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.all(2.0),
      child: ElevatedButton(
          onPressed: onTap,
          style: ElevatedButton.styleFrom(
            elevation: 4,
            backgroundColor: Colors.black54,
            visualDensity: VisualDensity.comfortable,
          ),
          child: textMath(latex, const TextStyle(fontSize: 16))),
    );
  }
}

// wrapper dispatching the concrete type
class _NodeW extends StatelessWidget {
  final bool showParenthesis;
  final _Node node;
  final List<_IndexedSet> sets;
  final _Node cursor;

  final void Function(_Node) onTap;

  const _NodeW(this.showParenthesis, this.node, this.sets, this.cursor,
      {required this.onTap, super.key});

  Widget kid() {
    switch (node.op) {
      case SetOp.sUnion:
        return _UniOrInt("\\cup", showParenthesis, node.args, sets, cursor,
            onTap: onTap);
      case SetOp.sInter:
        return _UniOrInt("\\cap", showParenthesis, node.args, sets, cursor,
            onTap: onTap);
      case SetOp.sComplement:
        return _ComplementW(
            node.args.isEmpty ? null : node.args.first, sets, cursor,
            onTap: onTap);
      case SetOp.sLeaf:
        return _LeafW(node.leaf, sets);
    }
  }

  @override
  Widget build(BuildContext context) {
    final isCursor = node == cursor;
    return InkWell(
      onTap: () => onTap(node),
      child: Padding(
        padding: node.op == SetOp.sLeaf
            ? const EdgeInsets.symmetric(vertical: 2)
            : const EdgeInsets.all(2),
        child: Container(
          decoration: BoxDecoration(
              borderRadius: const BorderRadius.all(Radius.circular(4)),
              color: isCursor
                  ? Colors.white.withOpacity(0.2)
                  : Colors.transparent),
          child: kid(),
        ),
      ),
    );
  }
}

const textStyle = TextStyle(fontSize: 20);

class _UniOrInt extends StatelessWidget {
  final String operator;
  final bool showParenthesis;
  final List<_Node> args;
  final List<_IndexedSet> sets;
  final _Node cursor;

  final void Function(_Node) onTap;

  const _UniOrInt(
      this.operator, this.showParenthesis, this.args, this.sets, this.cursor,
      {required this.onTap, super.key});

  @override
  Widget build(BuildContext context) {
    if (args.isEmpty) {
      // display an empty square
      return const SizedBox(height: 20, width: 10);
    }
    // add the n sign
    final l = <Widget>[];
    for (var i = 0; i < args.length; i++) {
      l.add(_NodeW(true, args[i], sets, cursor, onTap: onTap));
      if (i != args.length - 1) {
        l.add(textMath(operator, textStyle));
      }
    }
    if (args.length >= 2 && showParenthesis) {
      // add ()
      l.insert(0, textMath("(", textStyle));
      l.add(textMath(")", textStyle));
    }
    return Row(
        mainAxisSize: MainAxisSize.min,
        crossAxisAlignment: CrossAxisAlignment.baseline,
        textBaseline: TextBaseline.alphabetic,
        children: l);
  }
}

class _ComplementW extends StatelessWidget {
  final _Node? arg;
  final List<_IndexedSet> sets;
  final _Node cursor;

  final void Function(_Node) onTap;

  const _ComplementW(this.arg, this.sets, this.cursor,
      {required this.onTap, super.key});

  @override
  Widget build(BuildContext context) {
    final color = Theme.of(context).textTheme.bodyMedium?.color ?? Colors.black;
    return Container(
      decoration: BoxDecoration(border: Border(top: BorderSide(color: color))),
      child:
          arg == null ? null : _NodeW(false, arg!, sets, cursor, onTap: onTap),
    );
  }
}

class _LeafW extends StatelessWidget {
  final Set leaf;
  final List<_IndexedSet> sets;

  const _LeafW(this.leaf, this.sets, {super.key});

  @override
  Widget build(BuildContext context) {
    final set = sets.firstWhere((element) => element.index == leaf);
    return textMath(set.set, textStyle);
  }
}
