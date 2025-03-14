import 'package:eleve/build_mode.dart';
import 'package:eleve/main_eleve_shared.dart';
import 'package:eleve/quotes.dart';
import 'package:eleve/shared/settings_mobile.dart';
import 'package:flutter/material.dart';
import 'package:upgrader/upgrader.dart';

final bm = buildModeFromEnv();
// const bm = BuildMode.dev;

void main() async {
  initQuotes();

  final settingsHandler = FileSettings();

  final audio = await loadAudioFromSettings(settingsHandler);

  runApp(
    EleveApp(
      audio,
      settingsHandler,
      bm,
      checkUprades: Upgrader(
        messages: _UpgraderMessages(),
        durationUntilAlertAgain: const Duration(seconds: 1),
      ),
    ),
  );
}

class _UpgraderMessages extends UpgraderMessages {
  @override
  String get buttonTitleUpdate => "OK"; // note that we loose localisation
}
