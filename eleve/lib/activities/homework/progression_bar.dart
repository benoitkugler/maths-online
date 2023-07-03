import 'package:flutter/material.dart';

class ProgressionBar extends StatelessWidget {
  final int total;
  final int completed;
  final int started;

  const ProgressionBar(
      {required this.total,
      required this.completed,
      required this.started,
      Key? key})
      : super(key: key);

  @override
  Widget build(BuildContext context) {
    return LayoutBuilder(builder: (context, constraints) {
      final width = constraints.maxWidth;
      final startedWidth = width * started.toDouble() / total;
      final completedWidth = width * completed.toDouble() / total;
      return Stack(
        children: [
          _Layer(Colors.grey, width),
          _Layer(
            Colors.yellow.shade200,
            startedWidth,
            stripe: true,
          ),
          _Layer(Colors.lightGreenAccent, completedWidth),
        ],
      );
    });
  }
}

class _Layer extends StatelessWidget {
  final Color color;
  final double width;
  final bool stripe;
  const _Layer(this.color, this.width, {Key? key, this.stripe = false})
      : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Container(
      width: width,
      height: 12,
      decoration: BoxDecoration(
          gradient: stripe
              ? LinearGradient(
                  begin: Alignment.topLeft,
                  end: const Alignment(-0.95, 0),
                  stops: const [0.0, 0.5, 0.5, 1],
                  colors: [
                    color,
                    color,
                    Colors.transparent,
                    Colors.transparent,
                  ],
                  tileMode: TileMode.repeated,
                )
              : null,
          color: color,
          border: Border.all(width: 2, color: Colors.transparent),
          borderRadius: const BorderRadius.all(Radius.circular(4))),
    );
  }
}
