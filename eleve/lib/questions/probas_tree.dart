import 'package:eleve/questions/dropdown.dart';
import 'package:eleve/questions/fields.dart';
import 'package:eleve/questions/number.dart';
import 'package:eleve/types/src_maths_questions_client.dart';
import 'package:flutter/foundation.dart';
import 'package:flutter/material.dart';

TreeShape shape(TreeNodeAnswer tree) {
  if (tree.children.isEmpty) {
    return [];
  }
  final levelWidth = tree.children.length;
  return [levelWidth] + shape(tree.children[0]);
}

class TreeController extends FieldController {
  final TreeFieldBlock data;

  TreeShape? selectedShape;
  // setup when the shape is chosen
  _NodeController? controllers;

  TreeController(this.data, void Function() onChange) : super(onChange);

  void setShapeIndex(int? index) {
    setShape(index == null ? null : data.shapeProposals[index]);
  }

  void setShape(TreeShape? shape) {
    selectedShape = shape;
    if (shape != null) {
      controllers = _NodeController.editableFromShape(
          onChange, shape, true, data.eventsProposals, isEnabled);
    }
  }

  @override
  bool hasValidData() {
    if (controllers == null) {
      return false;
    }
    return controllers!.hasValidData();
  }

  @override
  Answer getData() {
    return TreeAnswer(controllers!.getData());
  }

  @override
  void setData(Answer answer) {
    final tree = (answer as TreeAnswer).root;
    setShape(shape(tree));
    controllers!.setData(tree);
  }

  @override
  void setEnabled(bool b) {
    // also disable children controllers
    super.setEnabled(b);
    controllers?.setEnabled(b);
  }
}

class _NodeController {
  final bool isRoot;
  final DropDownController? valueController;
  final List<_NodeController> children; // empty for the leafs
  final List<NumberController>? edgesController; // same length as children

  _NodeController(
      this.isRoot, this.valueController, this.children, this.edgesController);

  int levels() {
    if (children.isEmpty) {
      return 0;
    }
    return 1 + children[0].levels();
  }

  factory _NodeController.staticFromShape(TreeShape shape, bool isRoot) {
    if (shape.isEmpty) {
      return _NodeController(false, null, [], []);
    }
    final children = List<_NodeController>.generate(shape[0],
        (index) => _NodeController.staticFromShape(shape.sublist(1), false));
    return _NodeController(isRoot, null, children, null);
  }

  factory _NodeController.editableFromShape(void Function() onChange,
      TreeShape shape, bool isRoot, List<TextOrMath> proposals, bool enabled) {
    final controller =
        DropDownController(onChange, proposals.map((e) => [e]).toList());

    if (shape.isEmpty) {
      return _NodeController(isRoot, controller, [], []);
    }

    final children = List<_NodeController>.generate(
        shape[0],
        (index) => _NodeController.editableFromShape(
            onChange, shape.sublist(1), false, proposals, enabled));

    final edgesControllers = List<NumberController>.generate(
        children.length, (index) => NumberController(onChange));
    final dd = DropDownController(onChange, proposals.map((e) => [e]).toList());
    return _NodeController(isRoot, dd, children, edgesControllers);
  }

  void setEnabled(bool b) {
    valueController?.setEnabled(b);
    edgesController?.forEach((c) => c.setEnabled(b));
    for (var c in children) {
      c.setEnabled(b);
    }
  }

  bool hasValidData() {
    if (!isRoot &&
        (valueController == null || !valueController!.hasValidData())) {
      return false;
    }
    if (edgesController == null ||
        !edgesController!.every((element) => element.hasValidData())) {
      return false;
    }
    return children.every((element) => element.hasValidData());
  }

  // only valid if hasValidData is true
  // by convention, the root is expected to has 0 as value
  TreeNodeAnswer getData() {
    final childrenAnswers = children.map((e) => e.getData()).toList();
    final edgesAnswers = edgesController!.map((e) => e.getNumber()).toList();
    return TreeNodeAnswer(
        childrenAnswers, edgesAnswers, isRoot ? 0 : valueController!.index!);
  }

  void setData(TreeNodeAnswer answer) {
    if (!isRoot) {
      valueController!.setIndex(answer.value);
    }

    for (var i = 0; i < edgesController!.length; i++) {
      edgesController![i].setNumber(answer.probabilities[i]);
    }

    // recurse
    for (var i = 0; i < children.length; i++) {
      children[i].setData(answer.children[i]);
    }
  }
}

class TreeField extends StatefulWidget {
  final Color color;
  final TreeController controller;

  const TreeField(this.color, this.controller, {Key? key}) : super(key: key);

  @override
  _TreeFieldState createState() => _TreeFieldState();
}

class _TreeFieldState extends State<TreeField> {
  void _onSelectShape(int? shapeIndex) {
    setState(() {
      widget.controller.setShapeIndex(shapeIndex);
    });
  }

  void _showShapeSelection() async {
    final ct = widget.controller;
    final selected = await Navigator.of(context).push(
      MaterialPageRoute<int?>(
        builder: (context) => Scaffold(
          appBar: AppBar(),
          body: _ShapeSelection(
              widget.color, ct.data.shapeProposals, ct.selectedShape, (p0) {
            Navigator.of(context).pop(p0);
          }),
        ),
      ),
    );
    _onSelectShape(selected);
  }

  @override
  Widget build(BuildContext context) {
    final ct = widget.controller;
    return ct.selectedShape == null
        ? InkWell(
            onTap: ct.isEnabled ? _showShapeSelection : null,
            child: Container(
              padding: const EdgeInsets.all(12),
              decoration: BoxDecoration(
                border: Border.all(color: widget.color),
                borderRadius: BorderRadius.circular(5),
              ),
              child: const Text(
                "Sélectionner la forme de l'arbre...",
                style: TextStyle(fontStyle: FontStyle.italic, fontSize: 14),
              ),
            ),
          )
        : _OneTree(false, ct.hasError ? Colors.red : widget.color,
            _showShapeSelection, null, ct.controllers!);
  }
}

class _TreeView extends StatefulWidget {
  final Color color;
  final int nbLevels;
  final bool isSelected;
  final Positioned? backButton;
  final Widget root;
  final void Function()? onTap;

  const _TreeView(
      {super.key,
      required this.color,
      required this.nbLevels,
      required this.isSelected,
      required this.backButton,
      required this.root,
      required this.onTap});

  @override
  State<_TreeView> createState() => _TreeViewState();
}

class _TreeViewState extends State<_TreeView> {
  final ScrollController _controller = ScrollController();

  @override
  void dispose() {
    _controller.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    const boxPadding = 5.0;
    const levelHeightHint =
        _NodeLayout.edgesHeight + _NodeLayout.valueHeight + 2 * boxPadding + 2;

    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 10),
      child: Scrollbar(
        controller: _controller,
        radius: const Radius.circular(4),
        thumbVisibility: kIsWeb ? true : null,
        child: SingleChildScrollView(
          controller: _controller,
          scrollDirection: Axis.horizontal,
          child: InkWell(
            onTap: widget.onTap,
            child: SizedBox(
              height: widget.nbLevels * levelHeightHint,
              child: Stack(
                children: [
                  Container(
                      padding: const EdgeInsets.symmetric(vertical: boxPadding),
                      decoration: BoxDecoration(
                          border: Border.all(color: widget.color),
                          borderRadius: BorderRadius.circular(5),
                          color: widget.isSelected
                              ? Colors.white.withOpacity(0.3)
                              : null),
                      child: widget.root),
                  if (widget.backButton != null) widget.backButton!
                ],
              ),
            ),
          ),
        ),
      ),
    );
  }
}

class _OneTree extends StatelessWidget {
  final bool isSelected;
  final Color color;
  final void Function()? onBack;
  final void Function()? onTap;
  final _NodeController controller;

  const _OneTree(
      this.isSelected, this.color, this.onBack, this.onTap, this.controller,
      {Key? key})
      : super(key: key);

  @override
  Widget build(BuildContext context) {
    return _TreeView(
        color: color,
        nbLevels: controller.levels(),
        isSelected: isSelected,
        onTap: onTap,
        backButton: onBack == null
            ? null
            : Positioned(
                top: 3,
                left: 3,
                child: FloatingActionButton(
                    mini: true,
                    onPressed: controller.valueController?.isEnabled == true
                        ? onBack
                        : null,
                    tooltip: "Changer de forme",
                    child: const Icon(
                      IconData(0xe092,
                          fontFamily: 'MaterialIcons',
                          matchTextDirection: true),
                    )),
              ),
        root: _NodeEditable(color, controller));
  }
}

class _ShapeSelection extends StatelessWidget {
  final Color color;
  final List<TreeShape> proposals;
  final TreeShape? selected;
  final void Function(int) onSelect;

  const _ShapeSelection(
      this.color, this.proposals, this.selected, this.onSelect,
      {Key? key})
      : super(key: key);

  @override
  Widget build(BuildContext context) {
    return SingleChildScrollView(
      child: Column(
        mainAxisSize: MainAxisSize.min,
        children: List<Widget>.generate(
          proposals.length,
          (index) {
            final shapeI = proposals[index];
            final isSelected = listEquals(selected, shapeI);
            return _OneTree(isSelected, color, null, () => onSelect(index),
                _NodeController.staticFromShape(shapeI, true));
          },
        ),
      ),
    );
  }
}

class _NodeEditable extends StatefulWidget {
  final Color color;

  final _NodeController data;
  const _NodeEditable(this.color, this.data, {Key? key}) : super(key: key);

  @override
  _NodeEditableState createState() => _NodeEditableState();
}

class _NodeEditableState extends State<_NodeEditable> {
  void editEdge(int index) {
    final cts = widget.data.edgesController;
    if (cts == null) {
      return;
    }
    showDialog<void>(
        context: context,
        builder: (context) => Dialog(
            insetPadding: const EdgeInsets.all(20),
            child: Center(
              heightFactor: 2,
              child: NumberField(widget.color, cts[index], autofocus: true,
                  onSubmitted: () {
                Navigator.of(context).maybePop();
                setState(() {});
              }),
            )));
  }

  @override
  Widget build(BuildContext context) {
    final isRoot = widget.data.isRoot;
    final valueCt = widget.data.valueController;
    final ecs = widget.data.edgesController;
    final edges = ecs == null
        ? widget.data.children.map((e) => "")
        : ecs.map((ct) => ct.hasValidData() ? "${ct.getNumber()}" : " ? ");
    return _NodeLayout(
        color: widget.color,
        isRoot: isRoot,
        edges: edges.toList(),
        onTapEdge: valueCt?.isEnabled == true ? editEdge : null,
        value: valueCt == null
            ? const Text("?")
            : DropDownField(widget.color, valueCt),
        children: widget.data.children
            .map((e) => _NodeEditable(widget.color, e))
            .toList());
  }
}

class _NodeLayout extends StatelessWidget {
  static const edgesHeight = 50.0;
  static const valueHeight = 30.0;

  final Color color;
  final bool isRoot;
  final void Function(int)? onTapEdge;

  final List<String> edges;
  final Widget value;
  final List<Widget> children;

  const _NodeLayout(
      {super.key,
      required this.color,
      required this.isRoot,
      required this.onTapEdge,
      required this.edges,
      required this.value,
      required this.children});

  @override
  Widget build(BuildContext context) {
    const marginX = 12.0;

    final painter = _EdgesPainter(color, edges,
        isRoot ? 0 : _NodeLayout.valueHeight, _NodeLayout.edgesHeight);
    final hasChildren = children.isNotEmpty;

    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 4.0),
      child: GestureDetector(
        onTapUp: onTapEdge == null
            ? null
            : (details) {
                final index = painter.onHit(details.localPosition);
                if (index != null) {
                  onTapEdge!(index);
                }
              },
        child: CustomPaint(
          painter: painter,
          child: Column(
            children: [
              if (!isRoot)
                Container(
                  height: _NodeLayout.valueHeight,
                  margin: const EdgeInsets.symmetric(horizontal: marginX),
                  padding: const EdgeInsets.symmetric(horizontal: 8.0),
                  decoration: BoxDecoration(
                      border: Border.all(color: color),
                      borderRadius: const BorderRadius.all(Radius.circular(6))),
                  child: Column(
                    mainAxisAlignment: MainAxisAlignment.center,
                    children: [
                      value,
                    ],
                  ),
                ),
              if (hasChildren) ...[
                // make room for edges, drawn by CustomPaint
                const SizedBox(height: _NodeLayout.edgesHeight),
                // children row
                IntrinsicWidth(
                  child: Row(
                      mainAxisAlignment: MainAxisAlignment.center,
                      mainAxisSize: MainAxisSize.min,
                      children:
                          children.map((e) => Expanded(child: e)).toList()),
                ),
              ],
            ],
          ),
        ),
      ),
    );
  }
}

class _EdgesPainter extends CustomPainter {
  final Color color;
  final List<String> edges;
  final double startY;
  final double height;

  _EdgesPainter(this.color, this.edges, this.startY, this.height);

  List<Offset> _middles = []; // cached during paint

  @override
  void paint(Canvas canvas, Size size) {
    final N = edges.length;
    final childWidth = size.width / N;
    final start = Offset(size.width / 2, startY);
    final ends = List<Offset>.generate(
        N, (i) => Offset((i + 0.5) * childWidth, startY + height));

    _middles = List<Offset>.generate(N, (i) {
      final end = ends[i];
      return Offset((start.dx + end.dx) / 2, (start.dy + end.dy) / 2);
    });

    for (var end in ends) {
      canvas.drawLine(
          start,
          end,
          Paint()
            ..color = color
            ..strokeWidth = 2);
    }

    for (var i = 0; i < N; i++) {
      final text = edges[i];
      final painter = TextPainter(
        text: TextSpan(
            text: text,
            style: const TextStyle(
                color: Colors.black, fontWeight: FontWeight.bold)),
        textDirection: TextDirection.ltr,
      );
      painter.layout();

      canvas.drawRRect(
          RRect.fromRectAndRadius(
              Rect.fromCenter(
                  center: _middles[i],
                  width: painter.width + 15,
                  height: painter.height + 10),
              const Radius.circular(10)),
          Paint()
            ..style = PaintingStyle.fill
            ..color = Colors.white);
      painter.paint(
          canvas, _middles[i] - Offset(painter.width / 2, painter.height / 2));
    }
  }

  @override
  bool shouldRepaint(covariant _EdgesPainter oldDelegate) {
    return oldDelegate != this;
  }

  @override
  bool? hitTest(Offset position) {
    return onHit(position) == null ? false : true;
  }

  int? onHit(Offset position) {
    const fieldRadius = 20.0;
    for (var i = 0; i < _middles.length; i++) {
      final m = position - _middles[i];
      if (m.distanceSquared <= fieldRadius * fieldRadius) {
        return i;
      }
    }
    return null;
  }
}

extension on TreeNodeAnswer {
  int levels() {
    if (children.isEmpty) {
      return 0;
    }
    return 1 + children[0].levels();
  }
}

class Tree extends StatelessWidget {
  final Color color;
  final TreeBlock data;

  const Tree(this.color, this.data, {super.key});

  @override
  Widget build(BuildContext context) {
    return _TreeView(
        onTap: null,
        color: color,
        nbLevels: data.root.levels(),
        isSelected: false,
        backButton: null,
        root: _NodeStatic(data.eventsProposals, color, true, data.root));
  }
}

class _NodeStatic extends StatelessWidget {
  final List<TextOrMath> eventProposals;

  final Color color;
  final bool isRoot;
  final TreeNodeAnswer node;

  const _NodeStatic(this.eventProposals, this.color, this.isRoot, this.node,
      {super.key});

  @override
  Widget build(BuildContext context) {
    return _NodeLayout(
        color: color,
        isRoot: isRoot,
        onTapEdge: null,
        edges: node.probabilities.map((e) => "$e").toList(),
        value: isRoot
            ? const SizedBox()
            : TextRow(
                buildText([eventProposals[node.value]], TextS(), 14),
                lineHeight: 1.2,
              ),
        children: node.children
            .map((e) => _NodeStatic(eventProposals, color, false, e))
            .toList());
  }
}
