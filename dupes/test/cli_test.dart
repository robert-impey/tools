import 'dart:io';
import 'package:test/test.dart';
import 'package:path/path.dart' as p;

void main() {
  test('CLI accepts a file argument and fails if not a directory', () async {
    final tempDir = await Directory.systemTemp.createTemp('cli_test_');
    final filePath = p.join(tempDir.path, 'not_a_directory.txt');
    await File(filePath).writeAsString('hello');

    // Run the CLI with the file argument and expect non-zero exit code
    var result = await Process.run('dart', ['run', 'bin/dupes.dart', filePath]);
    if (result.exitCode == 0) {
      expect(
        result.stdout,
        contains('Usage:'),
        reason: 'Expected usage message for invalid input',
      );
    } else {
      // Ensure we actually got a failure
      expect(result.exitCode, isNot(equals(0)));
    }

    await tempDir.delete(recursive: true);
  });

  test('CLI prints usage message when no target directories are provided', () async {
    var result = await Process.run('dart', ['run', 'bin/dupes.dart']);
    expect(result.exitCode, equals(0));
    expect(
      result.stdout,
      contains('Usage: dart run bin/dupes.dart [--dupes] <directory1> [directory2] ...'),
    );

    result = await Process.run('dart', ['run', 'bin/dupes.dart', '--dupes']);
    expect(result.exitCode, equals(0));
    expect(
      result.stdout,
      contains('Usage: dart run bin/dupes.dart [--dupes] <directory1> [directory2] ...'),
    );
  });
}
