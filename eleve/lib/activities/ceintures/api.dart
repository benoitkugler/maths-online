import 'dart:convert';

import 'package:eleve/build_mode.dart';
import 'package:eleve/shared/errors.dart';
import 'package:eleve/types/src_prof_ceintures.dart';
import 'package:http/http.dart' as http;

abstract class CeinturesAPI {
  Future<GetEvolutionOut> getEvolution(StudentTokens args);
  Future<CreateEvolutionOut> createEvolution(CreateEvolutionIn args);
  Future<SelectQuestionsOut> selectQuestions(SelectQuestionsIn args);
  Future<EvaluateAnswersOut> evaluateAnswers(EvaluateAnswersIn args);
}

class ServerCeinturesAPI implements CeinturesAPI {
  final BuildMode buildMode;
  const ServerCeinturesAPI(this.buildMode);

  @override
  Future<GetEvolutionOut> getEvolution(StudentTokens args) async {
    const serverEndpoint = "/api/student/ceintures";
    final uri = buildMode.serverURL(serverEndpoint);
    final resp = await http
        .post(uri, body: jsonEncode(studentTokensToJson(args)), headers: {
      'Content-type': 'application/json',
    });
    return getEvolutionOutFromJson(checkServerError(resp.body));
  }

  @override
  Future<CreateEvolutionOut> createEvolution(CreateEvolutionIn args) async {
    const serverEndpoint = "/api/student/ceintures";
    final uri = buildMode.serverURL(serverEndpoint);
    final resp = await http
        .put(uri, body: jsonEncode(createEvolutionInToJson(args)), headers: {
      'Content-type': 'application/json',
    });
    return createEvolutionOutFromJson(checkServerError(resp.body));
  }

  @override
  Future<SelectQuestionsOut> selectQuestions(SelectQuestionsIn args) async {
    const serverEndpoint = "/api/student/ceintures/stage";
    final uri = buildMode.serverURL(serverEndpoint);
    final resp = await http
        .post(uri, body: jsonEncode(selectQuestionsInToJson(args)), headers: {
      'Content-type': 'application/json',
    });
    return selectQuestionsOutFromJson(checkServerError(resp.body));
  }

  @override
  Future<EvaluateAnswersOut> evaluateAnswers(EvaluateAnswersIn args) async {
    const serverEndpoint = "/api/student/ceintures/stage";
    final uri = buildMode.serverURL(serverEndpoint);
    final resp = await http
        .put(uri, body: jsonEncode(evaluateAnswersInToJson(args)), headers: {
      'Content-type': 'application/json',
    });
    return evaluateAnswersOutFromJson(checkServerError(resp.body));
  }
}
