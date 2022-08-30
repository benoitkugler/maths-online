import 'dart:ui';

import 'package:flutter/material.dart';

class PawnImage extends StatelessWidget {
  final Offset center;
  final double size;
  const PawnImage(this.center, this.size, {Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return AnimatedPositioned(
      duration: const Duration(milliseconds: 800),
      top: center.dy - size / 2,
      left: center.dx - size / 2,
      width: size,
      height: size,
      curve: Curves.easeIn,
      child: CustomPaint(
        size: Size(size, size),
        painter: _PawnPainter(),
      ),
    );
  }
}

class _PawnPainter extends CustomPainter {
  static const color = Colors.redAccent;

  // produced from SVG
  static final path = Path()
    ..moveTo(450, 1279.4)
    ..cubicTo(414.3, 1276.1, 380.4, 1264, 350.5, 1244.2)
    ..cubicTo(312.1, 1218.6, 291, 1188.4, 281.9, 1146)
    ..cubicTo(279.9, 1137, 279.6, 1132.4, 279.6, 1113.5)
    ..cubicTo(279.5, 1090.1, 280.3, 1083.2, 285.6, 1062.6)
    ..cubicTo(290.5, 1043.4, 302.3, 1015.5, 314.4, 994.3)
    ..lineTo(317.9, 988.2)
    ..lineTo(309.2, 984.5)
    ..cubicTo(261.5, 964, 222.8, 933.5, 199.8, 898)
    ..cubicTo(185.4, 875.8, 175.9, 850.9, 170.8, 821.2)
    ..cubicTo(168.8, 810.1, 168.5, 804.9, 168.6, 781)
    ..cubicTo(168.6, 757.1, 168.9, 751.8, 170.8, 740.5)
    ..cubicTo(179.6, 689.6, 202, 645.9, 233.7, 617.8)
    ..cubicTo(247.1, 605.9, 269.8, 592.7, 285.8, 587.6)
    ..cubicTo(288.6, 586.6, 291, 585.7, 291, 585.5)
    ..cubicTo(291, 585.3, 282.3, 580.7, 271.8, 575.3)
    ..cubicTo(189.5, 533, 124.1, 477.8, 82.3, 415.5)
    ..cubicTo(42.2, 355.7, 16.6, 283.3, 6, 200)
    ..cubicTo(1.1, 161, 0.6, 149.7, 0.6, 85)
    ..cubicTo(0.5, 51.2, 0.9, 18.2, 1.4, 11.7)
    ..lineTo(2.2, 0)
    ..lineTo(467, 0)
    ..lineTo(931.8, 0)
    ..lineTo(932.6, 11.8)
    ..cubicTo(934, 31.4, 933.7, 141.3, 932.2, 161)
    ..cubicTo(927.2, 228.7, 915.6, 287.4, 898.8, 330)
    ..cubicTo(877.6, 383.7, 847.4, 428.8, 802.6, 473.6)
    ..cubicTo(762.8, 513.4, 717.6, 546.2, 660.8, 576.7)
    ..cubicTo(648.5, 583.4, 645.2, 585.6, 646.6, 586.1)
    ..cubicTo(667.3, 593.9, 674.3, 598.1, 693.9, 614.4)
    ..cubicTo(708.9, 626.8, 716.2, 634.3, 726.5, 648.1)
    ..cubicTo(744, 671.5, 755.2, 697.6, 761.4, 729)
    ..cubicTo(767.9, 762, 767.3, 805.1, 760, 833)
    ..cubicTo(752.8, 860.7, 738.1, 893.1, 724.6, 911.1)
    ..cubicTo(702.2, 940.9, 668.1, 965.9, 624.8, 984.5)
    ..lineTo(616.1, 988.2)
    ..lineTo(619.1, 993.4)
    ..cubicTo(642.5, 1033.3, 654.3, 1077.2, 652.7, 1118)
    ..cubicTo(651.4, 1149, 643.6, 1173.2, 626.3, 1199.5)
    ..cubicTo(608.1, 1227.1, 586.5, 1246.7, 559.5, 1260)
    ..cubicTo(532.2, 1273.5, 509, 1278.8, 474, 1279.4)
    ..cubicTo(463.3, 1279.6, 452.5, 1279.6, 450, 1279.4)
    ..close();

  @override
  void paint(Canvas canvas, Size size) {
    final pathBounds = path.getBounds();

    // map pathBounds to size
    final factor = 1 * size.width / pathBounds.width;
    canvas.scale(factor, -factor);
    canvas.translate(0, -1100);

    canvas.drawPath(
        path,
        Paint()
          ..style = PaintingStyle.stroke
          ..color = Colors.white
          ..strokeWidth = 200
          ..imageFilter = ImageFilter.blur(sigmaX: 100, sigmaY: 100));
    canvas.drawPath(
        path,
        Paint()
          ..style = PaintingStyle.fill
          ..color = color);
  }

  @override
  bool shouldRepaint(covariant CustomPainter oldDelegate) {
    return false;
  }
}
