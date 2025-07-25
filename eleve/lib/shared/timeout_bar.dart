import 'package:flutter/material.dart';

class TimeoutBar extends StatefulWidget {
  final Duration duration;
  final Color color;

  const TimeoutBar(this.duration, this.color, {Key? key}) : super(key: key);

  @override
  _TimeoutBarState createState() => _TimeoutBarState();
}

class _TimeoutBarState extends State<TimeoutBar> {
  int seconds = 0;

  double get value =>
      1 - seconds.toDouble() / widget.duration.inSeconds.toDouble();

  @override
  void didUpdateWidget(covariant TimeoutBar oldWidget) {
    seconds = 0;
    super.didUpdateWidget(oldWidget);
  }

  @override
  void initState() {
    seconds = -1;
    _updateClock();
    super.initState();
  }

  @override
  void dispose() {
    // cancel the next timer
    seconds = widget.duration.inSeconds;
    super.dispose();
  }

  void _updateClock() async {
    if (seconds >= widget.duration.inSeconds) {
      return;
    }
    setState(() {
      seconds += 1;
    });
    await Future.delayed(const Duration(seconds: 1), _updateClock);
  }

  @override
  Widget build(BuildContext context) {
    return LinearProgressIndicator(
      color: widget.color,
      value: value,
      minHeight: 10,
    );
  }
}
