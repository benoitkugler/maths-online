import 'package:eleve/types/src_maths_questions_client.dart';
import 'package:flutter/material.dart';

class ImageW extends StatelessWidget {
  final ImageBlock data;
  const ImageW(this.data, {super.key});

  @override
  Widget build(BuildContext context) {
    return Image.network(
      data.uRL,
      scale: 100 / data.scale.toDouble(), // scale works a inverse
      webHtmlElementStrategy: WebHtmlElementStrategy.fallback,
    );
  }
}
