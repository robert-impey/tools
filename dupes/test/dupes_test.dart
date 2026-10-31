import 'dart:io';
import 'package:dupes/dupes.dart';
import 'package:test/test.dart';
import 'package:path/path.dart' as p;

void main() {
  late Directory tempDir;

  setUp(() async {
    tempDir = await Directory.systemTemp.createTemp('dupes_test_');
  });

  tearDown(() async {
    if (await tempDir.exists()) {
      await tempDir.delete(recursive: true);
    }
  });

  test('groupFilesByChecksum finds duplicates', () async {
    // Create some files
    final file1 = File(p.join(tempDir.path, 'file1.txt'));
    await file1.writeAsString('hello world');

    final file2 = File(p.join(tempDir.path, 'file2.txt'));
    await file2.writeAsString('hello world');

    final file3 = File(p.join(tempDir.path, 'file3.txt'));
    await file3.writeAsString('different content');

    final results = await groupFilesByChecksum([tempDir]);

    expect(results, hasLength(1));
    final dupes = results.values.first;
    expect(dupes, hasLength(2));
    expect(dupes, contains(file1.path));
    expect(dupes, contains(file2.path));
    expect(dupes, isNot(contains(file3.path)));
  });

  test('groupFilesByChecksum handles subdirectories', () async {
    // Create some files
    final file1 = File(p.join(tempDir.path, 'file1.txt'));
    await file1.writeAsString('hello world');

    final subDir = Directory(p.join(tempDir.path, 'sub'));
    await subDir.create();

    final file2 = File(p.join(subDir.path, 'file2.txt'));
    await file2.writeAsString('hello world');

    final results = await groupFilesByChecksum([tempDir]);

    expect(results, hasLength(1));
    final dupes = results.values.first;
    expect(dupes, hasLength(2));
    expect(dupes, contains(file1.path));
    expect(dupes, contains(file2.path));
  });

  test('groupFilesByChecksum ignores empty directories', () async {
    final subDir = Directory(p.join(tempDir.path, 'sub'));
    await subDir.create();

    final results = await groupFilesByChecksum([tempDir]);

    expect(results, isEmpty);
  });
}
