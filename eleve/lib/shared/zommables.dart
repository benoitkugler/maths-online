import 'package:flutter/material.dart';

/// [ListWithZoomables] is a work around a flutter limitation,
/// to enable zoomables widget inside lists.
class ListWithZoomables extends StatefulWidget {
  final List<Widget> children;
  final List<GlobalKey> zoomableKeys;
  final bool shrinkWrap;

  const ListWithZoomables(this.children, this.zoomableKeys,
      {Key? key, this.shrinkWrap = false})
      : super(key: key);

  @override
  State<ListWithZoomables> createState() => _ListWithZoomablesState();
}

class _ListWithZoomablesState extends State<ListWithZoomables> {
  /// this ensure that scrolling on the edge of the widget works
  static const scrollArea = 30;

  bool enableScroll = true;

  /// returns true if [position] is in the widget identified
  /// by [key]
  static bool _isInWidget(Offset position, GlobalKey key) {
    // find your widget
    final box = key.currentContext?.findRenderObject();
    if (box is! RenderBox) {
      return false;
    }

    // get offset
    Offset boxOffset = box.localToGlobal(Offset.zero);

    // check if your pointerdown event is inside the widget
    final isInY = position.dy > boxOffset.dy &&
        position.dy < boxOffset.dy + box.size.height;
    final isInX = position.dx > boxOffset.dx + scrollArea &&
        position.dx < boxOffset.dx + box.size.width - scrollArea;
    return isInY && isInX;
  }

  void _checkPan(Offset position) {
    final hasZommable = widget.zoomableKeys.any((key) {
      return _isInWidget(position, key);
    });
    setState(() {
      enableScroll = !hasZommable;
    });
  }

  @override
  Widget build(BuildContext context) {
    return Listener(
      onPointerUp: (ev) {
        // restore the scroll posibility
        setState(() {
          enableScroll = true;
        });
      },
      onPointerDown: (ev) {
        _checkPan(ev.position);
      },
      child: ListView(
        // if dragging over your widget, disable scroll, otherwise allow scrolling
        physics: enableScroll
            ? const ScrollPhysics()
            : const NeverScrollableScrollPhysics(),
        shrinkWrap: widget.shrinkWrap,
        children: widget.children,
      ),
    );
  }
}

/// [Zoomable] makes a widget zoomable,
/// also adding horizontal margin to allow vertical scrolling.
class Zoomable extends StatelessWidget {
  final TransformationController controller;
  final Widget child;
  final GlobalKey innerKey;

  const Zoomable(this.controller, this.child, this.innerKey, {Key? key})
      : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Padding(
        padding: const EdgeInsets.symmetric(horizontal: 10),
        child: InteractiveViewer(
          transformationController: controller,
          maxScale: 5,
          key: innerKey,
          child: child,
        ));
  }
}
