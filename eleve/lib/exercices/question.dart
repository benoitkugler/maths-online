import 'package:eleve/exercices/types.gen.dart';
import 'package:eleve/trivialpoursuit/timeout_bar.dart';
import 'package:flutter/material.dart';
import 'package:flutter_math_fork/flutter_math.dart';

class Question extends StatelessWidget {
  final ClientQuestion question;
  const Question(this.question, {Key? key}) : super(key: key);

  // we use Wrap instead of Rows to avoid overflows
  Wrap _buildRow(List<Widget> widgets) {
    return Wrap(
      children: widgets,
      spacing: 5,
      runSpacing: 5,
    );
  }

  // group the text and inline math block on the same row
  List<Widget> _buildColumns() {
    const textStyle = TextStyle(fontSize: 18);
    final List<Widget> out = [];
    List<Widget> currentRow = [];
    question.content.forEach((element) {
      if (element is TextBlock) {
        currentRow.add(Text(element.text, style: textStyle));
      } else if (element is FormulaBlock) {
        if (element.isInline) {
          currentRow.add(Math.tex(
            element.content,
            mathStyle: MathStyle.text,
            textStyle: textStyle,
          ));
        } else {
          // start a new row
          if (currentRow.isNotEmpty) {
            // close the current if needed
            out.add(_buildRow(currentRow));
            currentRow = [];
          }

          out.add(Center(
              child: Math.tex(
            element.content,
            mathStyle: MathStyle.display,
            textStyle: textStyle,
          )));
        }
      } else {
        // TODO:
      }
    });

    // flush the current row
    if (currentRow.isNotEmpty) {
      // close the current if needed
      out.add(_buildRow(currentRow));
    }

    return out;
  }

  @override
  Widget build(BuildContext context) {
    const spacing = SizedBox(height: 20.0);
    return Column(
      mainAxisAlignment: MainAxisAlignment.spaceEvenly,
      children: [
        const Text(
          "Thème test",
          style: TextStyle(fontSize: 20),
        ),
        Padding(
          padding: const EdgeInsets.symmetric(vertical: 12.0),
          child: Text(
            question.title,
            style: const TextStyle(fontSize: 20),
            textAlign: TextAlign.left,
          ),
        ),
        spacing,
        Expanded(
          child: ListView(
            shrinkWrap: true,
            children: _buildColumns()
                .map(
                  (e) => Padding(
                      padding: const EdgeInsets.symmetric(vertical: 6.0),
                      child: e),
                )
                .toList(),
          ),
        ),
        spacing,
        const TextField(
          cursorColor: Colors.green,
          decoration: InputDecoration(
            focusedBorder: OutlineInputBorder(
                borderSide: BorderSide(
              color: Colors.green,
            )),
            border: OutlineInputBorder(),
          ),
        ),
        spacing,
        ElevatedButton(
          onPressed: () => print("ok"),
          style: ElevatedButton.styleFrom(primary: Colors.green),
          child: const Text(
            "Valider",
            style: TextStyle(fontSize: 18),
          ),
        ),
        spacing,
        const TimeoutBar(Duration(seconds: 60), Colors.green),
        spacing,
        const Text("", style: TextStyle(fontSize: 16)),
      ],
    );
  }
}
