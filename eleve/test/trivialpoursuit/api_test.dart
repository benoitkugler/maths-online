import 'dart:convert';

import 'package:eleve/trivialpoursuit/events.gen.dart';

void main() async {
  final jsonMessage = jsonEncode(clientEventITFToJson(const Ping("info")));
  print(jsonMessage);
}
