import 'package:eleve/questions/fields.dart';
import 'package:eleve/questions/number.dart';
import 'package:eleve/questions/types.gen.dart';
import 'package:flutter/material.dart';

class VectorController extends FieldController {
  final VectorFieldBlock data;

  final NumberController x;
  final NumberController y;

  VectorController(this.data, void Function() onChange)
      : x = NumberController(onChange),
        y = NumberController(onChange),
        super(onChange);

  @override
  bool hasValidData() {
    return x.hasValidData() && y.hasValidData();
  }

  @override
  Answer getData() {
    return VectorNumberAnswer(x.getNumber(), y.getNumber());
  }

  @override
  void setData(Answer answer) {
    final ans = answer as VectorNumberAnswer;
    x.setNumber(ans.x);
    y.setNumber(ans.y);
  }
}

class VectorField extends StatelessWidget {
  final Color color;

  final VectorController controller;

  const VectorField(this.color, this.controller, {Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    final color = controller.fieldError ? Colors.red : this.color;
    final x = NumberField(
      color,
      controller.x,
      onSubmitted: controller.onChange,
    );
    final y = NumberField(
      color,
      controller.y,
      onSubmitted: controller.onChange,
    );

    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 6.0),
      child: controller.data.displayColumn
          ? RichText(
              text: TextSpan(children: [
              const TextSpan(
                  text: "(",
                  style: TextStyle(fontSize: 60, fontWeight: FontWeight.w300)),
              WidgetSpan(
                  alignment: PlaceholderAlignment.bottom,
                  baseline: TextBaseline.alphabetic,
                  child: Transform.translate(
                    offset: const Offset(0, 6),
                    child: Column(
                      children: [x, const SizedBox(height: 5), y],
                    ),
                  )),
              const TextSpan(
                  text: ")",
                  style: TextStyle(fontSize: 60, fontWeight: FontWeight.w300)),
            ]))
          : Row(
              mainAxisSize: MainAxisSize.min,
              children: [
                const Text(
                  "(",
                  style: TextStyle(fontSize: 20),
                ),
                x,
                const Text(
                  ";",
                  style: TextStyle(fontSize: 20),
                ),
                y,
                const Text(
                  ")",
                  style: TextStyle(fontSize: 20),
                )
              ],
            ),
    );
  }
}
