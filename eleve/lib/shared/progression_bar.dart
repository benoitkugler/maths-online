import 'package:flutter/material.dart';

class ProgressionLayer {
  // in [0; 1]
  final double advance;
  final Color color;
  final bool stripped;
  const ProgressionLayer(this.advance, this.color, this.stripped);
}

class ProgressionBar extends StatelessWidget {
  final Color background;
  final List<ProgressionLayer> layers;

  const ProgressionBar(
      {required this.background, required this.layers, Key? key})
      : super(key: key);

  @override
  Widget build(BuildContext context) {
    return LayoutBuilder(builder: (context, constraints) {
      final width = constraints.maxWidth;
      return Stack(
        children: [
          _Layer(background, width),
          ...layers.map(
              (l) => _Layer(l.color, width * l.advance, stripped: l.stripped)),
        ],
      );
    });
  }
}

class _Layer extends StatelessWidget {
  final Color color;
  final double width;
  final bool stripped;
  const _Layer(this.color, this.width, {Key? key, this.stripped = false})
      : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Container(
      width: width,
      height: 12,
      decoration: BoxDecoration(
          gradient: stripped
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
