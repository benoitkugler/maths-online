import 'package:flutter/gestures.dart';
import 'package:flutter/material.dart';
import 'package:url_launcher/url_launcher.dart';

TextSpan _hyperlink(String text, Uri url, TextStyle style) {
  return TextSpan(
    style: style.copyWith(
        color: Colors.lightBlue.shade200, decoration: TextDecoration.underline),
    text: text,
    recognizer: TapGestureRecognizer()
      ..onTap = () async {
        if (await canLaunchUrl(url)) {
          await launchUrl(url, mode: LaunchMode.externalApplication);
        }
      },
  );
}

/// [hyperlink] return a [Text] widget displaying [text] and
/// pointing to [urlText]
Text hyperlink(String text, String urlText, {TextStyle? style}) {
  final uri = Uri.tryParse(urlText);
  if (uri == null) {
    // invalid url
    return Text(text, style: style);
  }

  return Text.rich(_hyperlink(text, uri, style ?? const TextStyle()));
}

List<TextSpan> parseURLs(String text, TextStyle style) {
  final out = <TextSpan>[];
  final re = RegExp(r"https:\/\/(\S+)");
  var currentIndex = 0;
  for (var match in re.allMatches(text)) {
    if (match.start > currentIndex) {
      // add normal text
      out.add(TextSpan(
          text: text.substring(currentIndex, match.start), style: style));
    }
    final urlText = match.group(0)!;
    // handle the potential url
    final uri = Uri.tryParse(urlText);
    if (uri == null) {
      // invalid url
      out.add(TextSpan(text: urlText, style: style));
    } else {
      out.add(_hyperlink(match.group(1)!, uri, style));
    }
    currentIndex = match.end;
  }
  // handle the remaining non url part
  if (currentIndex < text.length) {
    out.add(TextSpan(text: text.substring(currentIndex), style: style));
  }
  return out;
}
