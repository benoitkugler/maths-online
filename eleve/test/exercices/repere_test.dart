import 'package:eleve/exercices/repere.dart';
import 'package:flutter/widgets.dart';
import 'package:flutter_test/flutter_test.dart';

void main() {
  testWidgets('repere ...', (tester) async {
    var l = parseSubscript("Normal text", const TextStyle(fontSize: 16));
    expect(l.length, equals(1));
    expect((l[0] as TextSpan).text, equals("Normal text"));

    l = parseSubscript("Normal text with C_f", const TextStyle(fontSize: 16));
    expect(l.length, equals(2));
    expect((l[0] as TextSpan).text, equals("Normal text with C"));
    expect(l[1].runtimeType, equals(WidgetSpan));

    l = parseSubscript("Normal text with C_fg", const TextStyle(fontSize: 16));
    expect(l.length, equals(2));
    expect((l[0] as TextSpan).text, equals("Normal text with C"));
    expect(l[1].runtimeType, equals(WidgetSpan));

    l = parseSubscript(
        "Normal text with C_fg then other", const TextStyle(fontSize: 16));
    expect(l.length, equals(3));
    expect((l[0] as TextSpan).text, equals("Normal text with C"));
    expect(l[1].runtimeType, equals(WidgetSpan));
    expect((l[2] as TextSpan).text, equals(" then other"));
  });
}
