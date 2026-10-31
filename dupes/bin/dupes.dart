import 'dart:io';
import 'package:dupes/dupes.dart';

void main(List<String> args) async {
  bool findDupes = false;
  List<Directory> targetDirs = [];

  for (var arg in args) {
    if (arg == '--dupes') {
      findDupes = true;
    } else {
      final dir = Directory(arg);
      if (await dir.exists()) {
        targetDirs.add(dir);
      } else {
        final file = File(arg);
        if (await file.exists()) {
          // If a file is provided, use its parent directory
          targetDirs.add(file.parent);
        }
      }
    }
  }

  if (targetDirs.isEmpty) {
    print(
      "Usage: dart run bin/dupes.dart [--dupes] <directory1> [directory2] ...",
    );
    return;
  }

  final Map<String, List<String>> duplicateChecksums =
      await groupFilesByChecksum(targetDirs);
  await printGroups(duplicateChecksums, findDupes);
}
