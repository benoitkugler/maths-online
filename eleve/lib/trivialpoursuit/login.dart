import 'package:eleve/build_mode.dart';
import 'package:eleve/trivialpoursuit/game.dart';
import 'package:flutter/material.dart';
import 'package:otp_text_field/otp_text_field.dart';
import 'package:otp_text_field/style.dart';

const pinLength = 4;

/// Loggin is an introduction screen to access
/// a TrivialPoursuit game
class TrivialPoursuitLoggin extends StatefulWidget {
  final BuildMode buildMode;
  const TrivialPoursuitLoggin(this.buildMode, {Key? key}) : super(key: key);

  @override
  _TrivialPoursuitLogginState createState() => _TrivialPoursuitLogginState();
}

class _TrivialPoursuitLogginState extends State<TrivialPoursuitLoggin> {
  OtpFieldController otpController = OtpFieldController();

  @override
  void initState() {
    Future.delayed(const Duration(milliseconds: 50), () {
      otpController.setFocus(0);
    });

    if (widget.buildMode == BuildMode.debug) {
      // skip loggin screen
      Future.delayed(
          const Duration(milliseconds: 50), () => _launchTrivialPoursuit(""));
    }

    super.initState();
  }

  @override
  Widget build(BuildContext context) {
    final otp = OTPTextField(
      controller: otpController,
      length: pinLength,
      width: MediaQuery.of(context).size.width,
      textFieldAlignment: MainAxisAlignment.center,
      otpFieldStyle: OtpFieldStyle(
          enabledBorderColor: Theme.of(context).colorScheme.secondary),
      fieldStyle: FieldStyle.box,
      fieldWidth: 45,
      onCompleted: _launchTrivialPoursuit,
      onChanged: (_) {},
    );

    return Scaffold(
      appBar: AppBar(
        title: const Text("Rejoindre une partie"),
      ),
      body: Card(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            const Padding(
              padding: EdgeInsets.symmetric(vertical: 20),
              child: Text(
                "Code de la partie",
                style: TextStyle(fontSize: 20),
              ),
            ),
            Padding(
                padding: const EdgeInsets.symmetric(vertical: 20), child: otp),
          ],
        ),
      ),
    );
  }

  void _launchTrivialPoursuit(String code) {
    final route = Navigator.of(context).push(MaterialPageRoute<void>(
      settings: const RouteSettings(name: "/board"),
      builder: (_) => Scaffold(
          appBar: AppBar(),
          body: TrivialPoursuitController(widget.buildMode, code)),
    ));

    route.then((value) {
      otpController.clear();
    });
  }
}
