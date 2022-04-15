import 'dart:convert';

import 'package:eleve/exercices/question.dart';
import 'package:eleve/exercices/types.gen.dart';
import 'package:eleve/loopback_types.gen.dart';
import 'package:flutter_test/flutter_test.dart';

void main() {
  testWidgets('main prof loopback ...', (tester) async {
    final notif1 = ValidQuestionNotification(const QuestionAnswersIn({}));
    final json = jsonEncode({
      "Kind":
          loopbackClientDataKindToJson(LoopbackClientDataKind.validAnswerIn),
      "Data": questionAnswersInToJson(notif1.data)
    });

    final roundrip = jsonDecode(json) as Map<String, dynamic>;
    expect(roundrip.length, equals(2));

    final notif2 = CheckQuestionSyntaxeNotification(
        const QuestionSyntaxCheckIn(NumberAnswer(1), 0));
    final json2 = jsonEncode({
      "Kind":
          loopbackClientDataKindToJson(LoopbackClientDataKind.checkSyntaxIn),
      "Data": questionSyntaxCheckInToJson(notif2.data)
    });

    final roundrip2 = jsonDecode(json2) as Map<String, dynamic>;
    expect(roundrip2.length, equals(2));
  });
}
