import 'package:flutter/material.dart';

class TimeoutBar extends StatelessWidget {
  final TimeoutBarController controller;
  final Color color;

  const TimeoutBar(this.controller, this.color, {super.key});

  @override
  Widget build(BuildContext context) {
    return LinearProgressIndicator(
      color: color,
      value: controller.value,
      minHeight: 10,
      borderRadius: BorderRadius.circular(4),
    );
  }
}

class TimeoutBarController {
  int totalDuration;
  int seconds = 0;

  TimeoutBarController(this.totalDuration);

  // in [0;1]
  double get value => 1 - seconds.toDouble() / totalDuration.toDouble();
}
