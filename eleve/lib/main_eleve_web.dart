import 'package:web/web.dart';

import 'package:eleve/build_mode.dart';
import 'package:eleve/main_eleve_shared.dart';
import 'package:eleve/quotes.dart';
import 'package:eleve/shared/settings_web.dart';
import 'package:flutter/material.dart';

void main() async {
  initQuotes();

  // on the web, we enable dev mode with query param
  final uri = Uri.parse(window.location.href);
  final mode = uri.queryParameters["mode"];
  final bm = APISetting.fromString(mode ?? "");

  final handler = LocalStorageSettings();
  final audio = await loadAudioFromSettings(handler);

  runApp(EleveApp(audio, handler, bm));
}
