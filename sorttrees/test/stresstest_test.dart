import 'dart:io';

import 'package:test/test.dart';

void main() {
  test('stresstest CLI runs and completes successfully after configured duration', () async {
    final result = await Process.run(
      Platform.resolvedExecutable,
      [
        'bin/stresstest.dart',
        '--Duration',
        '0.5',
        '--Depth',
        '2',
        '--Items',
        '2',
      ],
      workingDirectory: Directory.current.path,
    );

    expect(result.exitCode, 0);
    expect(
      result.stdout,
      contains('Stress test completed successfully'),
    );
    expect(result.stderr, isEmpty);
  });
}
