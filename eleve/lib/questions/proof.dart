import 'package:eleve/questions/drag_text.dart';
import 'package:eleve/questions/fields.dart';
import 'package:eleve/types/src_maths_questions_client.dart';
import 'package:flutter/material.dart';

class ProofController extends FieldController {
  final ProofFieldBlock enonce;

  _SequenceController controller; // root controller
  List<TextLine> currentProposals;

  ProofController(this.enonce, void Function() onChange)
      : controller = _SequenceController(enonce.shape.root, onChange),
        currentProposals = enonce.termProposals.toList(),
        super(onChange);

  @override
  void setEnabled(bool enabled) {
    controller.setEnabled(enabled);
    super.setEnabled(enabled);
  }

  @override
  void setError(bool hasError) {
    controller.setError(hasError);
    super.setError(hasError);
  }

  @override
  Answer getData() {
    return ProofAnswer(Proof(controller.getData() as Sequence));
  }

  @override
  void setData(Answer answer) {
    controller.setData((answer as ProofAnswer).proof.root);
    // by construction, the answer uses all the proposals,
    // so we don't bother with diffing
    currentProposals.clear();
  }

  @override
  bool hasValidData() {
    return _hasValidData(controller.getData());
  }
}

bool _hasValidData(Assertion a) {
  if (a is Statement) {
    return a.content.isNotEmpty;
  } else if (a is Equality) {
    return a.terms.every((element) => element.isNotEmpty);
  } else if (a is Node) {
    return a.op != Binary.invalid &&
        _hasValidData(a.left) &&
        _hasValidData(a.right);
  } else if (a is Sequence) {
    return a.parts.every(_hasValidData);
  } else {
    throw ("exhaustive type switch");
  }
}

// state object responsible for one node
abstract class _AssertionController {
  _AssertionController();

  bool _hasError = false;
  bool get hasError => _hasError;
  void setError(bool hasError) {
    _hasError = hasError;
  }

  bool _isEnabled = true;
  bool get isEnabled => _isEnabled;
  void setEnabled(bool b) {
    _isEnabled = b;
  }

  Assertion getData();
  void setData(Assertion data);

  factory _AssertionController.fromData(
      Assertion data, void Function() onChange) {
    if (data is Statement) {
      return _StatementController(data, onChange);
    } else if (data is Equality) {
      return _EqualityController(data, onChange);
    } else if (data is Node) {
      return _NodeController(data, onChange);
    } else if (data is Sequence) {
      return _SequenceController(data, onChange);
    } else {
      throw ("exhaustive type switch");
    }
  }
}

class _StatementController extends _AssertionController {
  Statement data;
  final void Function() onChange;

  _StatementController(this.data, this.onChange);

  @override
  Assertion getData() {
    return data;
  }

  @override
  void setData(Assertion data) {
    this.data = (data as Statement);
  }

  void setStatement(TextLine statement) {
    data = Statement(statement);
    onChange();
  }
}

class _EqualityController extends _AssertionController {
  Equality data;
  final void Function() onChange;
  _EqualityController(this.data, this.onChange);

  @override
  Assertion getData() {
    return data;
  }

  @override
  void setData(Assertion data) {
    this.data = (data as Equality);
  }

  void setTerm(TextLine term, int index) {
    data.terms[index] = term;
    onChange();
  }

  void setAvecDef(TextLine def) {
    data = Equality(data.terms, def, data.withDef);
    onChange();
  }
}

class _NodeController extends _AssertionController {
  Binary _op;
  _AssertionController left;
  _AssertionController right;
  final void Function() onChange;
  _NodeController(
    Node data,
    this.onChange,
  )   : _op = data.op,
        left = _AssertionController.fromData(data.left, onChange),
        right = _AssertionController.fromData(data.right, onChange);

  @override
  void setError(bool hasError) {
    left.setError(hasError);
    right.setError(hasError);
    super.setError(hasError);
  }

  @override
  void setEnabled(bool b) {
    left.setEnabled(b);
    right.setEnabled(b);
    super.setEnabled(b);
  }

  @override
  Assertion getData() {
    return Node(left.getData(), right.getData(), _op);
  }

  @override
  void setData(Assertion data) {
    final node = (data as Node);
    _op = node.op;
    left.setData(node.left);
    right.setData(node.right);
  }

  void setOp(Binary op) {
    _op = op;
    onChange();
  }
}

class _SequenceController extends _AssertionController {
  List<_AssertionController> parts;

  _SequenceController(Sequence data, void Function() onChange)
      : parts = data.parts
            .map((e) => _AssertionController.fromData(e, onChange))
            .toList();

  @override
  void setError(bool hasError) {
    for (var element in parts) {
      element.setError(hasError);
    }
    super.setError(hasError);
  }

  @override
  void setEnabled(bool b) {
    for (var element in parts) {
      element.setEnabled(b);
    }
    super.setEnabled(b);
  }

  @override
  Assertion getData() {
    return Sequence(parts.map((ct) => ct.getData()).toList());
  }

  @override
  void setData(Assertion data) {
    // we assume the shape does not change
    final seq = (data as Sequence).parts;
    for (var i = 0; i < parts.length; i++) {
      parts[i].setData(seq[i]);
    }
  }
}

class ProofField extends StatefulWidget {
  final Color color;
  final ProofController controller;

  const ProofField(this.color, this.controller, {Key? key}) : super(key: key);

  @override
  State<ProofField> createState() => _ProofFieldState();
}

class _ProofFieldState extends State<ProofField> {
  void _removeProposal(_TermUsed term) {
    // if the old value is empty, do not replace it
    setState(() {
      if (term.replaced.isEmpty) {
        widget.controller.currentProposals.removeAt(term.choosen.index);
      } else {
        widget.controller.currentProposals[term.choosen.index] = term.replaced;
      }
    });
  }

  @override
  Widget build(BuildContext context) {
    // two parts : top for the proof skeleton, bottom for a list
    // of proposals
    return Column(
      children: [
        NotificationListener<_TermUsed>(
          onNotification: (n) {
            _removeProposal(n);
            return true;
          },
          child: _AssertionW(
            widget.color,
            widget.controller.controller,
            isRoot: true,
          ),
        ),
        const SizedBox(height: 5),
        _Proposals(
            widget.controller.isEnabled, widget.controller.currentProposals)
      ],
    );
  }
}

class _Proposals extends StatelessWidget {
  final bool enabled;
  final List<TextLine> proposals;

  const _Proposals(this.enabled, this.proposals, {Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    final children = List<Widget>.generate(
        proposals.length,
        (i) => DragText(_TermFromProposals(proposals[i], i), proposals[i],
            enabled: enabled));
    return Wrap(
      runSpacing: 6,
      children: children,
    );
  }
}

abstract class _DragData {}

class _TermFromProposals implements _DragData {
  final TextLine term;
  final int index;
  _TermFromProposals(this.term, this.index);
}

class _TermFromProof implements _DragData {
  final TextLine term;
  final void Function(TextLine otherTerm) onSwap;

  _TermFromProof(this.term, this.onSwap);
}

// widget for one assertion node in a proof
class _AssertionW extends StatelessWidget {
  final Color color;
  final bool isRoot;
  final _AssertionController controller;

  const _AssertionW(this.color, this.controller,
      {Key? key, this.isRoot = false})
      : super(key: key);

  @override
  Widget build(BuildContext context) {
    final Widget child;
    final ct = controller;
    if (ct is _StatementController) {
      child = _StatementW(color, ct);
    } else if (ct is _EqualityController) {
      child = _EqualityW(color, ct);
    } else if (ct is _NodeController) {
      child = _NodeW(color, ct);
    } else if (ct is _SequenceController) {
      child = _SequenceW(color, ct);
    } else {
      throw ("exhaustive switch");
    }
    return Container(
      padding: const EdgeInsets.symmetric(vertical: 3, horizontal: 1),
      child: child,
    );
  }
}

class _TermUsed extends Notification {
  final _TermFromProposals choosen;
  final TextLine replaced;
  const _TermUsed(this.choosen, this.replaced);
}

class _TextDrop extends StatelessWidget {
  final bool enabled;
  final TextLine text;

  // called when receiving a drop
  final void Function(_DragData origin) onDrop;
  // called when the drag started from this is accepted
  final void Function(TextLine otherTerm) onSwap;

  const _TextDrop(this.enabled, this.text, this.onDrop, this.onSwap, {Key? key})
      : super(key: key);

  @override
  Widget build(BuildContext context) {
    return DragTarget<_DragData>(
        builder: (context, candidateData, rejectedData) {
          return Container(
            padding: EdgeInsets.symmetric(
                vertical: text.isEmpty ? 6 : 2, horizontal: 2),
            decoration: BoxDecoration(
              color: candidateData.isNotEmpty ? Colors.blue : null,
              border:
                  text.isEmpty ? Border.all(color: Colors.blueAccent) : null,
              borderRadius: const BorderRadius.all(Radius.circular(4)),
            ),
            constraints: const BoxConstraints(minWidth: 40),
            child: text.isEmpty
                ? const Padding(
                    padding: EdgeInsets.symmetric(horizontal: 6.0),
                    child: Text(
                      "Glisser...",
                      style: TextStyle(fontStyle: FontStyle.italic),
                    ),
                  )
                : DragText(
                    _TermFromProof(text, onSwap),
                    text,
                    enabled: enabled,
                    dense: true,
                  ),
          );
        },
        onAccept: (term) => onDrop(term));
  }
}

class _StatementW extends StatefulWidget {
  final Color color;
  final _StatementController controller;

  const _StatementW(this.color, this.controller, {Key? key}) : super(key: key);

  @override
  State<_StatementW> createState() => _StatementWState();
}

class _StatementWState extends State<_StatementW> {
  @override
  Widget build(BuildContext context) {
    final text = widget.controller.data.content;
    final content = _TextDrop(widget.controller.isEnabled, text, (origin) {
      if (origin is _TermFromProposals) {
        setState(() {
          widget.controller.setStatement(origin.term);
        });
        _TermUsed(origin, text).dispatch(context);
      } else if (origin is _TermFromProof) {
        setState(() {
          widget.controller.setStatement(origin.term);
        });
        origin.onSwap(text);
      }
    }, (otherTerm) {
      setState(() {
        widget.controller.setStatement(otherTerm);
      });
    });
    return text.isEmpty
        ? content // fill the width
        : Align(
            alignment: Alignment.centerLeft,
            child: content,
          );
  }
}

class _EqualityW extends StatefulWidget {
  final Color color;
  final _EqualityController controller;

  const _EqualityW(this.color, this.controller, {Key? key}) : super(key: key);

  @override
  State<_EqualityW> createState() => _EqualityWState();
}

class _EqualityWState extends State<_EqualityW> {
  @override
  Widget build(BuildContext context) {
    final terms = widget.controller.data.terms;

    final List<Widget> children = [];

    for (var i = 0; i < terms.length; i++) {
      children.add(_TextDrop(widget.controller.isEnabled, terms[i], (notif) {
        final currentText = terms[i];
        if (notif is _TermFromProposals) {
          setState(() {
            widget.controller.setTerm(notif.term, i);
          });
          _TermUsed(notif, currentText).dispatch(context);
        } else if (notif is _TermFromProof) {
          setState(() {
            widget.controller.setTerm(notif.term, i);
          });
          notif.onSwap(currentText);
        }
      }, (otherTerm) {
        setState(() {
          widget.controller.setTerm(otherTerm, i);
        });
      }));

      children.add(const Padding(
        padding: EdgeInsets.symmetric(horizontal: 8.0),
        child: Text("="),
      ));
    }
    children.removeLast();

    final data = widget.controller.data;
    if (data.withDef) {
      children.add(const Padding(
        padding: EdgeInsets.symmetric(horizontal: 8.0),
        child: Text(" avec"),
      ));
      children.add(_TextDrop(widget.controller.isEnabled, data.def, (notif) {
        final currentText = data.def;
        if (notif is _TermFromProposals) {
          setState(() {
            widget.controller.setAvecDef(notif.term);
          });
          _TermUsed(notif, currentText).dispatch(context);
        } else if (notif is _TermFromProof) {
          setState(() {
            widget.controller.setAvecDef(notif.term);
          });
          notif.onSwap(currentText);
        }
      }, (otherTerm) {
        setState(() {
          widget.controller.setAvecDef(otherTerm);
        });
      }));
    }

    return Wrap(
      alignment: WrapAlignment.center,
      crossAxisAlignment: WrapCrossAlignment.center,
      runSpacing: 4,
      children: children,
    );
  }
}

class _WithBorder extends StatelessWidget {
  final Widget child;
  final Color color;

  const _WithBorder({required this.child, required this.color, Key? key})
      : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.symmetric(vertical: 3, horizontal: 3),
      decoration: BoxDecoration(
          border: Border.all(width: 1, color: color),
          borderRadius: const BorderRadius.all(Radius.circular(4))),
      child: child,
    );
  }
}

class _NodeW extends StatefulWidget {
  final Color color;
  final _NodeController controller;

  const _NodeW(this.color, this.controller, {Key? key}) : super(key: key);

  @override
  State<_NodeW> createState() => _NodeWState();
}

class _NodeWState extends State<_NodeW> {
  void _onOpChanged(Binary op) {
    setState(() {
      widget.controller.setOp(op);
    });
  }

  @override
  Widget build(BuildContext context) {
    final ct = widget.controller;
    return _WithBorder(
      color: widget.color,
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.stretch,
        children: [
          _AssertionW(widget.color, ct.left),
          Row(children: [
            Expanded(child: Container(height: 1, color: widget.color)),
            _WithBorder(
              color: widget.color,
              child: Padding(
                padding: const EdgeInsets.all(8.0),
                child: _BinaryField(widget.color, ct._op, ct.hasError,
                    ct.isEnabled ? _onOpChanged : null),
              ),
            ),
            Expanded(child: Container(height: 1, color: widget.color)),
          ]),
          _AssertionW(widget.color, ct.right),
        ],
      ),
    );
  }
}

class _SequenceW extends StatelessWidget {
  final Color color;
  final _SequenceController controller;

  const _SequenceW(this.color, this.controller, {Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    final List<Widget> children = [];
    for (var part in controller.parts) {
      children.add(_AssertionW(color, part));
      children.add(const Padding(
        padding: EdgeInsets.symmetric(vertical: 8.0, horizontal: 2),
        child: Text("donc"),
      ));
    }
    children.removeLast();
    return Column(
      crossAxisAlignment: CrossAxisAlignment.stretch,
      children: children,
    );
  }
}

extension _BinaryLabels on Binary {
  String label() {
    switch (this) {
      case Binary.invalid:
        return "";
      case Binary.and:
        return "et";
      case Binary.or:
        return "ou";
    }
  }
}

class _BinaryField extends StatelessWidget {
  final Color color;
  final Binary value;
  final bool hasError;
  final void Function(Binary)? onChanged;

  const _BinaryField(this.color, this.value, this.hasError, this.onChanged,
      {Key? key})
      : super(key: key);

  @override
  Widget build(BuildContext context) {
    return DropdownButton<Binary>(
      isDense: true,
      style: hasError ? TextStyle(color: Colors.red.shade200) : null,
      underline: hasError
          ? Container(height: 1.0, color: Colors.red
              // color: Colors.red,
              )
          : null,
      focusColor: color,
      dropdownColor: color,
      hint: const Text("Choisir..."),
      value: value != Binary.invalid ? value : null,
      iconSize: 0,
      alignment: Alignment.center,
      // we use selectedItemBuilder since Math.tex do not handle
      // keys in a way that permit reusing the widgets in items
      selectedItemBuilder: (_) => List.generate(
        Binary.values.length,
        (index) => Padding(
          padding: const EdgeInsets.symmetric(horizontal: 5.0),
          child: Text(value.label()),
        ),
      ),
      items: List.generate(
        Binary.values.length,
        (index) => DropdownMenuItem<Binary>(
          value: Binary.values[index],
          child: Padding(
            padding: const EdgeInsets.symmetric(horizontal: 3),
            child: Text(Binary.values[index].label()),
          ),
        ),
      ),
      onChanged:
          onChanged != null ? (b) => onChanged!(b ?? Binary.invalid) : null,
    );
  }
}
