import 'dart:io';

import 'package:sorttrees/randomtrees.dart';
import 'package:test/test.dart';

void main() {
  test('generateRandomTreeLines generates correct structure', () {
    final lines = generateRandomTreeLines(depth: 2, items: 3, length: 5);
    // depth 2 with 3 items each:
    // root level: 3 items
    // each item has 3 sub items: 3 * 3 = 9
    // total = 12 lines
    expect(lines.length, 12);
    for (final line in lines) {
      expect(line.trim().length, 5);
    }
  });

  test('speedtest CLI runs and outputs expected stats', () async {
    final result = await Process.run(
      Platform.resolvedExecutable,
      ['bin/speedtest.dart', '--Iterations', '3', '--Depth', '2', '--Items', '2'],
      workingDirectory: Directory.current.path,
    );

    expect(result.exitCode, 0);
    expect(result.stdout, contains('Number of times: 3'));
    expect(result.stdout, contains('Min time:'));
    expect(result.stdout, contains('Max time:'));
    expect(result.stdout, contains('Mean time:'));
    expect(result.stderr, isEmpty);
  });
}
