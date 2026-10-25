import 'dart:math';

const String defaultAlphabet =
    'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';

String generateRandomItem(
  int length, {
  String alphabet = defaultAlphabet,
  Random? random,
}) {
  final rng = random ?? Random();
  final buffer = StringBuffer();
  for (var i = 0; i < length; i++) {
    final index = rng.nextInt(alphabet.length);
    buffer.write(alphabet[index]);
  }
  return buffer.toString();
}

List<String> generateRandomTreeLines({
  int depth = 3,
  int items = 10,
  int length = 8,
  String alphabet = defaultAlphabet,
  String currentIndentation = '',
  Random? random,
}) {
  final lines = <String>[];
  if (depth > 0) {
    final rng = random ?? Random();
    var remainingItems = items;
    while (remainingItems > 0) {
      lines.add(
        '$currentIndentation${generateRandomItem(length, alphabet: alphabet, random: rng)}',
      );
      remainingItems--;
      final subLines = generateRandomTreeLines(
        depth: depth - 1,
        items: items,
        length: length,
        alphabet: alphabet,
        currentIndentation: '$currentIndentation    ',
        random: rng,
      );
      lines.addAll(subLines);
    }
  }
  return lines;
}
