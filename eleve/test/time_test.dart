import 'package:eleve/activities/homework/utils.dart';
import 'package:flutter_test/flutter_test.dart';

void main() {
  test('format time', () {
    final dt = DateTime.parse("2023-09-16T18:00:00Z");
    expect(dt.isUtc, true);
    expect(formatTime(dt), "Samedi 16 Sept. 2023, 20h");
  });
}
