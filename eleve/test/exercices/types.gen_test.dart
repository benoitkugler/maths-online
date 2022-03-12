import 'dart:convert';

import 'package:eleve/exercices/types.gen.dart';
import 'package:test/test.dart';

void main() {
  test("JSON", () {
    const input = """
{"Title":"Calcul littéral","Content":[{"Data":{"Content":"Développer l’expression : "},"Kind":2},{"Data":{"Content":"\$ \\\\left(x - 6\\\\right)\\\\left(4x - 3\\\\right) \$"},"Kind":2}]}    """;

    final question = clientQuestionFromJson(jsonDecode(input));

    expect(question.title, equals("Calcul littéral"));
    expect(question.content.length, equals(2));
    expect(question.content[0] is LaTeXBlock, equals(true));
    expect(question.content[0] is LaTeXBlock, equals(true));
  });
}
