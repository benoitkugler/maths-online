import 'dart:convert';

import 'package:eleve/questions/types.gen.dart';
import 'package:test/test.dart';

void main() {
  test("JSON", () {
    const input =
        "{\"Title\":\"Calcul littéral\",\"Enonce\":[{\"Data\":{\"Text\":\"Développer l’expression : \"},\"Kind\":4},{\"Data\":{\"Content\":\"\\\\left(x - 6\\\\right)\\\\left(4x - 3\\\\right)\",\"IsInline\":true},\"Kind\":0}]}";
    final question = questionFromJson(jsonDecode(input));

    expect(question.title, equals("Calcul littéral"));
    expect(question.enonce.length, equals(2));
    expect(question.enonce[0] is TextBlock, equals(true));
    expect(question.enonce[1] is FormulaBlock, equals(true));
  });
}
