import 'package:eleve/build_mode.dart';
import 'package:eleve/main_eleve_shared.dart';
import 'package:flutter/material.dart';

final bm = buildMode();
// final bm = BuildMode.dev;

void main() async {
  final audio = await loadAudioFromSettings();

  runApp(EleveApp(audio, bm));
}
