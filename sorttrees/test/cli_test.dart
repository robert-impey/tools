import 'dart:io';

import 'package:test/test.dart';

void main() {
  late Directory directory;

  setUp(() async {
    directory = await Directory.systemTemp.createTemp('sorttrees-cli-test-');
  });

  tearDown(() async {
    await directory.delete(recursive: true);
  });

  test('prints sorted output with separated top-level trees', () async {
    final file = File('${directory.path}${Platform.pathSeparator}tree.txt');
    await file.writeAsString('b\n    z\n    a\n\na\n');

    final result = await _runCli([file.path, '--SeparateTopLevel']);

    expect(result.exitCode, 0);
    expect(result.stdout, 'a\n\nb\n    a\n    z\n ');
    expect(result.stderr, isEmpty);
  });

  test('sorts in place after creating a backup', () async {
    final file = File('${directory.path}${Platform.pathSeparator}tree.txt');
    const original = 'b\n    z\n    a\n';
    await file.writeAsString(original);

    final result = await _runCli([file.path, '--InPlace']);

    expect(result.exitCode, 0);
    expect(await file.readAsString(), 'b\n    a\n    z\n');
    expect(await File('${file.path}.bak').readAsString(), original);
  });

  test('does not create a backup for an already sorted file', () async {
    final file = File('${directory.path}${Platform.pathSeparator}tree.txt');
    await file.writeAsString('a\n    b\n');

    final result = await _runCli([file.path, '--InPlace']);

    expect(result.exitCode, 0);
    expect(result.stdout, 'In place sorting requested but already sorted.\n');
    expect(await File('${file.path}.bak').exists(), isFalse);
  });
}

Future<ProcessResult> _runCli(List<String> args) {
  return Process.run(Platform.resolvedExecutable, [
    'bin/sorttrees.dart',
    ...args,
  ], workingDirectory: Directory.current.path);
}
