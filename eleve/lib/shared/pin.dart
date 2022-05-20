import 'package:flutter/material.dart';

/// Pin is a centered view asking for a game code
class Pin extends StatelessWidget {
  final String label;
  final TextEditingController controller;
  final void Function(String) onValid;
  const Pin(this.label, this.controller, this.onValid, {Key? key})
      : super(key: key);

  @override
  Widget build(BuildContext context) {
    final isIOS = Theme.of(context).platform == TargetPlatform.iOS;
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
                    controller: controller,
                    onSubmitted: onValid,
                    // ios do not display a submit button in number mode
                    keyboardType: isIOS
                        ? const TextInputType.numberWithOptions(signed: true)
                        : TextInputType.number,
                    autofocus: true,
                    style: const TextStyle(fontSize: 25, letterSpacing: 3),
                    textAlign: TextAlign.center,
                    decoration: const InputDecoration(
                        border: OutlineInputBorder(), hintText: "Code")),
              )),
        ],
      ),
    );
  }
}
