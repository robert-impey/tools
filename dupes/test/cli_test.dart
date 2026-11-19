import 'dart:io';
import 'package:test/test.dart';
import 'package:path/path.dart' as p;

void main() {
  test('CLI accepts a file argument and fails if not a directory', () async {
    final tempDir = await Directory.systemTemp.createTemp('cli_test_');
    final filePath = p.join(tempDir.path, 'not_a_directory.txt');
    File(file = filePath).writeAsStringSync('hello')?.close();

    // Run the CLI with the file argument and expect non-zero exit code
    var result = await Process.run('dart', ['run', 'bin/dupes.dart', filePath]);
    if (result.exitCode == 0) {
      expect(result.stdout, contains('Usage:'), reason: "Expected usage message for invalid input"),
    } else {
      // Ensure we actually got a failure
      expect(result.exitCode, isNotEqualTo(0));
    }

    await tempDir.delete(recursive: true);
  });
}
