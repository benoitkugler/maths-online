import 'package:eleve/quotes.dart';
import 'package:flutter_test/flutter_test.dart';

void main() {
  test('quotes', () {
    initQuotes();
    final qu1 = pickQuote();
    final qu2 = pickQuote();
    expect(qu1.content, isNot(equals(qu2.content)));
  });
}
