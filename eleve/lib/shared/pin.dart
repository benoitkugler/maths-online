import 'package:flutter/material.dart';

/// Pin is a centered view asking for a game code
class Pin extends StatelessWidget {
  final String label;
  final void Function(String) onValid;
  const Pin(this.label, this.onValid, {Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Card(
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Center(
            child: Padding(
              padding: const EdgeInsets.symmetric(vertical: 20),
              child: Text(
                label,
                style: const TextStyle(fontSize: 20),
              ),
            ),
          ),
          Padding(
              padding: const EdgeInsets.symmetric(vertical: 20, horizontal: 10),
              child: SizedBox(
                width: 200,
                child: TextField(
                  onSubmitted: onValid,
                  keyboardType: TextInputType.number,
                  autofocus: true,
                  style: const TextStyle(fontSize: 25, letterSpacing: 3),
                  textAlign: TextAlign.center,
                  decoration:
                      const InputDecoration(border: OutlineInputBorder()),
                ),
              )),
        ],
      ),
    );
  }
}
