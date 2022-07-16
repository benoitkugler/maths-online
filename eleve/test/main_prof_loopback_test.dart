import 'dart:convert';

import 'package:eleve/loopback_types.gen.dart';
import 'package:eleve/questions/question.dart';
import 'package:eleve/questions/types.gen.dart';
import 'package:flutter_test/flutter_test.dart';

void main() {
  testWidgets('main prof loopback ...', (tester) async {
    final notif1 = ValidQuestionNotification(const QuestionAnswersIn({}));
    final json = jsonEncode(
        loopbackClientEventToJson(LoopbackQuestionValidIn(notif1.data)));
    final roundrip = jsonDecode(json) as Map<String, dynamic>;
    expect(roundrip.length, equals(2));
  });
}
