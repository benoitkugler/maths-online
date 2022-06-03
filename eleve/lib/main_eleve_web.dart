import 'dart:js' as js;

import 'package:eleve/build_mode.dart';
import 'package:eleve/main_eleve_shared.dart';
import 'package:flutter/material.dart';

void main() async {
  // on the web, we enable dev mode with query param
  final uri = Uri.parse(js.context['location']['href'] as String);
  final mode = uri.queryParameters["mode"];
  final bm = APISetting.fromString(mode ?? "");

  final audio = await loadAudioFromSettings();

  runApp(EleveApp(audio, bm));
}
