import 'package:eleve/audio.dart';
import 'package:eleve/build_mode.dart';
import 'package:eleve/main_eleve_shared.dart';
import 'package:flutter/material.dart';

final bm = buildMode();
// final bm = BuildMode.dev;

void main() {
  final audio = Audio();
  // start with some defaults
  audio.setSongs([0, 1]);
  runApp(EleveApp(audio, bm));
}
