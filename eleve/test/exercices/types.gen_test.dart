import 'dart:convert';

import 'package:eleve/exercices/types.gen.dart';
import 'package:test/test.dart';

void main() {
  test("JSON", () {
    const input = """
{"Title":"Calcul littéral","Content":[{"Data":{"Text":"Développer l’expression : "},"Kind":3},{"Data":{"Content":"\\\\left(x - 6\\\\right)\\\\left(4x - 3\\\\right)","IsInline":true},"Kind":2}]}
""";
    final question = clientQuestionFromJson(jsonDecode(input));

    expect(question.title, equals("Calcul littéral"));
    expect(question.content.length, equals(2));
    expect(question.content[0] is TextBlock, equals(true));
    expect(question.content[1] is FormulaBlock, equals(true));
  });
}
