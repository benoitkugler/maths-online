import 'package:eleve/questions/fields.dart';
import 'package:eleve/questions/proof.gen.dart';
import 'package:eleve/questions/types.gen.dart';
import 'package:flutter/material.dart';

class ProofController extends FieldController {
  final ProofFieldBlock enonce;

  Proof currentProof;

  ProofController(this.enonce, void Function() onChange)
      : currentProof = enonce.shape,
        super(onChange);

  @override
  Answer getData() {
    return ProofAnswer(currentProof);
  }

  @override
  void setData(Answer answer) {
    currentProof = (answer as ProofAnswer).proof;
  }

  @override
  bool hasValidData() {
    return _hasValidData(currentProof.root);
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

class ProofField extends StatelessWidget {
  final Color color;
  final ProofController controller;

  const ProofField(this.color, this.controller, {Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return _AssertionW(
      controller.currentProof.root,
      _AssertionParams(color, controller.fieldError, controller.enabled),
      isRoot: true,
    );
  }
}

class _AssertionParams {
  final Color color;
  final bool hasError;
  final bool isEnabled;
  const _AssertionParams(
    this.color,
    this.hasError,
    this.isEnabled,
  );
}

// widget for one assertion node in a proof
class _AssertionW extends StatelessWidget {
  final Assertion data;
  final _AssertionParams params;
  final bool isRoot;
  const _AssertionW(this.data, this.params, {Key? key, this.isRoot = false})
      : super(key: key);

  @override
  Widget build(BuildContext context) {
    final Widget child;
    final s = data;
    if (s is Statement) {
      child = _StatementW(s, params);
    } else if (s is Equality) {
      child = _EqualityW(s, params);
    } else if (s is Node) {
      child = _NodeW(s, params);
    } else if (s is Sequence) {
      child = _SequenceW(s, params);
    } else {
      throw ("exhaustive switch");
    }
    return Container(
      padding: const EdgeInsets.all(4),
      decoration: isRoot
          ? null
          : BoxDecoration(
              border: Border.all(width: 1.5, color: Colors.lightBlueAccent),
              borderRadius: const BorderRadius.all(Radius.circular(4))),
      child: child,
    );
  }
}

class _StatementW extends StatelessWidget {
  final Statement data;
  final _AssertionParams params;

  const _StatementW(this.data, this.params, {Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Text(data.content);
  }
}

class _EqualityW extends StatelessWidget {
  final Equality data;
  final _AssertionParams params;

  const _EqualityW(this.data, this.params, {Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    final List<Widget> children = [];
    for (var term in data.terms) {
      children.add(Text(term));
      children.add(const Padding(
        padding: EdgeInsets.symmetric(horizontal: 8.0),
        child: Text("="),
      ));
    }
    children.removeLast();
    return Wrap(
      children: children,
      alignment: WrapAlignment.center,
    );
  }
}

class _NodeW extends StatelessWidget {
  final Node data;
  final _AssertionParams params;

  const _NodeW(this.data, this.params, {Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    final left = _AssertionW(data.left, params);
    final right = _AssertionW(data.right, params);
    final op = _BinaryField(params.color, data.op, params.hasError,
        params.isEnabled ? print : null); // TODO:
    return Column(
      crossAxisAlignment: CrossAxisAlignment.stretch,
      children: [
        left,
        Padding(
          padding: const EdgeInsets.all(8.0),
          child: Center(
            widthFactor: 4,
            child: SizedBox(
              child: op,
              width: 100,
            ),
          ),
        ),
        right
      ],
    );
  }
}

class _SequenceW extends StatelessWidget {
  final Sequence data;
  final _AssertionParams params;

  const _SequenceW(this.data, this.params, {Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    final List<Widget> children = [];
    for (var part in data.parts) {
      children.add(_AssertionW(part, params));
      children.add(const Center(
          child: Padding(
        padding: EdgeInsets.all(8.0),
        child: Text("donc"),
      )));
    }
    children.removeLast();
    return Column(
      children: children,
      crossAxisAlignment: CrossAxisAlignment.stretch,
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
        return "or";
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
      hint: const Text("   "),
      value: value,
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
