import 'package:eleve/activities/ceintures/animations.dart';
import 'package:eleve/shared/settings_shared.dart';
import 'package:eleve/types/src_sql_ceintures.dart';
import 'package:flutter/material.dart';

class CeinturesActivityIcon extends StatelessWidget {
  final void Function() onTap;

  const CeinturesActivityIcon(this.onTap, {Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Column(
      children: [
        RawMaterialButton(
          onPressed: onTap,
          child: Image.asset("assets/images/yellow-belt.png",
              width: 68, height: 60),
        ),
        const Padding(
          padding: EdgeInsets.only(top: 8, bottom: 6),
          child: Text("Ceintures de calcul"),
        ),
      ],
    );
  }
}

extension on Rank {
  Color get color {
    switch (this) {
      case Rank.startRank:
        return Colors.transparent;
      case Rank.blanche:
        return Colors.white;
      case Rank.jaune:
        return Colors.yellow;
      case Rank.orange:
        return Colors.orange;
      case Rank.verte:
        return Colors.lightGreen;
      case Rank.bleue:
        return Colors.blue;
      case Rank.rouge:
        return Colors.red;
      case Rank.marron:
        return Colors.brown;
      case Rank.noire:
        return Colors.black87;
    }
  }
}

class CeinturesStart extends StatefulWidget {
  final UserSettings settings;
  const CeinturesStart(this.settings, {super.key});

  @override
  State<CeinturesStart> createState() => _CeinturesStartState();
}

class _CeinturesStartState extends State<CeinturesStart> {
  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Center(
        child: UnlockAnimation(Rank.orange.color),
        // Image.asset(
        //   "assets/images/yellow-belt.png",
        //   width: 68,
        //   color: Rank.bleue.color,
        // ),
      ),
    );
  }
}
