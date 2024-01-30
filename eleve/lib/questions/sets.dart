import 'package:eleve/questions/drag_text.dart';
import 'package:eleve/questions/fields.dart';
import 'package:eleve/types/src_maths_expression_sets.dart';
import 'package:eleve/types/src_maths_questions_client.dart';
import 'package:flutter/material.dart';

class SetController extends FieldController {
  final List<String> sets;
  SetController(super.onChange, this.sets);

  ListNode answer = ListNode([], SetOp.sLeaf, 0);

  @override
  Answer getData() {
    return SetAnswer(answer);
  }

  @override
  bool hasValidData() {
    return answer.isValid();
  }

  @override
  void setData(Answer answer) {
    this.answer = (answer as SetAnswer).root;
  }
}

extension on ListNode {
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
}

class SetFieldW extends StatefulWidget {
  final Color color;
  final SetController controller;

  const SetFieldW(this.color, this.controller, {super.key});

  @override
  State<SetFieldW> createState() => _SetFieldWState();
}

class _SetFieldWState extends State<SetFieldW> {
  @override
  Widget build(BuildContext context) {
    final sets = widget.controller.sets;
    return Column(
      mainAxisSize: MainAxisSize.min,
      children: [
        _NodeW(widget.controller.answer, sets),
        // TODO: controls
        Row(
            children: List.generate(
                sets.length,
                (index) => DragText(index, [TextOrMath(sets[index], true)],
                    enabled: true)).toList()),
      ],
    );
  }
}

// wrapper dispatching the concrete type
class _NodeW extends StatelessWidget {
  final ListNode node;
  final List<String> sets;

  const _NodeW(this.node, this.sets, {super.key});

  Widget kid() {
    switch (node.op) {
      case SetOp.sUnion:
        return _UnionW(node.args, sets);
      case SetOp.sInter:
        return _InterW(node.args, sets);
      case SetOp.sComplement:
        return _ComplementW(node.args.isEmpty ? null : node.args.first, sets);
      case SetOp.sLeaf:
        return _LeafW(node.leaf, sets);
    }
  }

  @override
  Widget build(BuildContext context) {
    return Padding(padding: const EdgeInsets.all(4), child: kid());
  }
}

const textStyle = TextStyle(fontSize: 16);

class _UnionW extends StatelessWidget {
  final List<ListNode> args;
  final List<String> sets;

  const _UnionW(this.args, this.sets, {super.key});

  @override
  Widget build(BuildContext context) {
    // add the U sign
    final l = <Widget>[];
    for (var i = 0; i < args.length; i++) {
      l.add(_NodeW(args[i], sets));
      if (i != args.length - 1) {
        l.add(textMath("\\cup", textStyle));
      }
    }
    if (args.length >= 2) {
      // add ()
      l.insert(0, const Text("("));
      l.add(const Text(")"));
    }
    return Row(children: l);
  }
}

class _InterW extends StatelessWidget {
  final List<ListNode> args;
  final List<String> sets;

  const _InterW(this.args, this.sets, {super.key});

  @override
  Widget build(BuildContext context) {
    // add the n sign
    final l = <Widget>[];
    for (var i = 0; i < args.length; i++) {
      l.add(_NodeW(args[i], sets));
      if (i != args.length - 1) {
        l.add(textMath("\\cap", textStyle));
      }
    }
    if (args.length >= 2) {
      // add ()
      l.insert(0, const Text("("));
      l.add(const Text(")"));
    }
    return Row(children: l);
  }
}

class _ComplementW extends StatelessWidget {
  final ListNode? arg;
  final List<String> sets;

  const _ComplementW(this.arg, this.sets, {super.key});

  @override
  Widget build(BuildContext context) {
    final color = Theme.of(context).textTheme.bodyMedium?.color ?? Colors.black;
    return Container(
      decoration: BoxDecoration(border: Border(top: BorderSide(color: color))),
      // padding: const EdgeInsets.only(top: 2),
      child: arg == null ? null : _NodeW(arg!, sets),
    );
  }
}

class _LeafW extends StatelessWidget {
  final Set leaf;
  final List<String> sets;

  const _LeafW(this.leaf, this.sets, {super.key});

  @override
  Widget build(BuildContext context) {
    return textMath(sets[leaf], textStyle);
  }
}
