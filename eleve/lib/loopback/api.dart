import 'dart:convert';

import 'package:eleve/build_mode.dart';
import 'package:eleve/shared/errors.dart';
import 'package:eleve/types/src_maths_questions.dart';
import 'package:eleve/types/src_maths_questions_client.dart';
import 'package:eleve/types/src_prof_editor.dart';
import 'package:eleve/types/src_tasks.dart';
import 'package:http/http.dart' as http;

class LoopbackAPI {
  final BuildMode buildMode;

  const LoopbackAPI(this.buildMode);

  Future<LoopbackEvaluateQuestionOut> evaluateQuestionAnswer(
      QuestionAnswersIn data, LoopbackShowQuestion origin) async {
    final uri =
        Uri.parse(buildMode.serverURL("/api/loopack/evaluate-question"));
    final params =
        LoopackEvaluateQuestionIn(origin.origin, AnswerP(origin.params, data));
    final resp = await http.post(uri,
        body: jsonEncode(loopackEvaluateQuestionInToJson(params)),
        headers: {
          'Content-type': 'application/json',
        });
    return loopbackEvaluateQuestionOutFromJson(checkServerError(resp.body));
  }

  Future<LoopbackShowQuestionAnswerOut> showQuestionAnswer(
      QuestionPage originPage, Params originParams) async {
    final uri = Uri.parse(buildMode.serverURL("/api/loopack/question-answer"));
    final params = LoopbackShowQuestionAnswerIn(originPage, originParams);
    final resp = await http.post(uri,
        body: jsonEncode(loopbackShowQuestionAnswerInToJson(params)),
        headers: {
          'Content-type': 'application/json',
        });
    return loopbackShowQuestionAnswerOutFromJson(checkServerError(resp.body));
  }
}
