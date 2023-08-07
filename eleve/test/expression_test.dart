import 'package:eleve/questions/expression.dart';
import 'package:flutter_test/flutter_test.dart';

void main() {
  test('split expression ...', () async {
    expect(
        detectNumerator("(2x + 5)"), equals(const SplitString("", "2x + 5")));
    expect(detectNumerator("((2x + 8) + 5)"),
        equals(const SplitString("", "(2x + 8) + 5")));
    expect(detectNumerator("(2x + 8)(1 + 5)"),
        equals(const SplitString("(2x + 8)", "1 + 5")));
    expect(detectNumerator(" -log(2x + 5)"),
        equals(const SplitString(" -", "log(2x + 5)")));
    expect(detectNumerator("a +(2x + 5)"),
        equals(const SplitString("a +", "2x + 5")));
    expect(detectNumerator("2x"), equals(const SplitString("", "2x")));
    expect(detectNumerator("x - 1"), equals(const SplitString("x - ", "1")));
    expect(detectNumerator("(x - 1)"), equals(const SplitString("", "x - 1")));
  });
}
