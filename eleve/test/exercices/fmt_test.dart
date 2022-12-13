import 'package:eleve/questions/number.dart';
import 'package:flutter_test/flutter_test.dart';

void main() {
  test('float formating ...', () async {
    var l = NumberController(() {});
    l.setNumber(6);
    expect(l.textController.text, equals("6"));

    l.setNumber(6.45);
    expect(l.textController.text, equals("6,45"));
  });
}
