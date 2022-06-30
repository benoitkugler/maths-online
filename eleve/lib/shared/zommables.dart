import 'package:flutter/material.dart';

/// [ListWithZoomables] is a work around a flutter limitation,
/// to enable zoomables widget inside lists
class ListWithZoomables extends StatefulWidget {
  final List<Widget> children;
  final List<GlobalKey<ZoomableState>> zoomableKeys;
  final bool shrinkWrap;

  const ListWithZoomables(this.children, this.zoomableKeys,
      {Key? key, this.shrinkWrap = false})
      : super(key: key);

  @override
  State<ListWithZoomables> createState() => _ListWithZoomablesState();
}

class _ListWithZoomablesState extends State<ListWithZoomables> {
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
    final isInX = position.dx > boxOffset.dx &&
        position.dx < boxOffset.dx + box.size.width;
    return isInY && isInX;
  }

  void _checkPan(Offset position) {
    final hasZommable = widget.zoomableKeys.any((key) {
      if (_isInWidget(position, key)) {
        // check if the zoomable widget is activated or not
        return key.currentState?.activated ?? false;
      }
      return false;
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

/// [Zoomable] wraps a widget into an InteractiveViewer,
/// which is activated with double tap
class Zoomable extends StatefulWidget {
  final TransformationController controller;
  final Widget child;

  const Zoomable(this.controller, this.child, GlobalKey<ZoomableState> key)
      : super(key: key);

  @override
  State<Zoomable> createState() => ZoomableState();
}

class ZoomableState extends State<Zoomable> {
  bool activated = false;

  void _switch() {
    setState(() {
      if (activated) {
        // also reset the zoom
        widget.controller.value = Matrix4.identity();
      }
      activated = !activated;
    });
  }

  @override
  Widget build(BuildContext context) {
    return GestureDetector(
      onTap: _switch,
      onDoubleTap: _switch,
      child: AbsorbPointer(
        absorbing: !activated,
        child: Container(
          decoration: BoxDecoration(
              borderRadius: const BorderRadius.all(Radius.circular(4)),
              border: Border.all(
                color: activated ? Colors.blueAccent : Colors.transparent,
                width: 2,
              )),
          child: InteractiveViewer(
            transformationController: widget.controller,
            child: widget.child,
            maxScale: 5,
            onInteractionEnd: (_) {
              if (widget.controller.value.getMaxScaleOnAxis() == 1) {
                setState(() {
                  activated = false;
                });
              }
            },
          ),
        ),
      ),
    );
  }
}
