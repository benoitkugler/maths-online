import 'dart:convert';

import 'package:eleve/types/src_trivial.dart';

void main() async {
  final jsonMessage = jsonEncode(clientEventITFToJson(const Ping("info")));
  print(jsonMessage);
}
