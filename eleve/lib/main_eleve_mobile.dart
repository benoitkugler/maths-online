import 'package:eleve/build_mode.dart';
import 'package:eleve/main_eleve_shared.dart';
import 'package:flutter/material.dart';
import 'package:upgrader/upgrader.dart';

// final bm = buildModeFromEnv();
final bm = BuildMode.debug;

void main() async {
  final audio = await loadAudioFromSettings();

  runApp(
    EleveApp(
      audio,
      bm,
      checkUprades: Upgrader(
        messages: _UpgraderMessages(),
        canDismissDialog: false,
        showIgnore: false,
        showLater: false,
        durationUntilAlertAgain: const Duration(seconds: 1),
      ),
    ),
  );
}

class _UpgraderMessages extends UpgraderMessages {
  @override
  String get buttonTitleUpdate => "OK"; // note that we loose localisation
}
