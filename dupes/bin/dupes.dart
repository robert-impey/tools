import 'dart:io';
import 'package:crypto/crypto.dart';

/// Groups files found within specified directories that are duplicates based on their content checksum.
Future<Map<String, List<String>>> groupFilesByChecksum(List<Directory> dirs) async {
  final Map<int, List<String>> fileSizes = {};

  // 1. Traverse directories and group paths by size
  for (var dir in dirs) {
    if (await dir.exists()) {
      await _walkDir(dir, fileSizes);
    }
  }

  final Map<String, List<String>> duplicateChecksums = {};

  // 2. Process sizes to find duplicates
  final sortedEntries = fileSizes.entries.toList()
    ..sort((a, b) => a.key.compareTo(b.key));

  for (var entry in sortedEntries) {
    final List<String> paths = entry.value;

    if (paths.length > 1) {
      // Group paths by checksum
      final Map<String, List<String>> checksums = {};
      for (var path in paths) {
        try {
          final file = File(path);
          final data = await file.readAsBytes();

          // Calculate SHA256 hash
          final digest = sha256.convert(data);
          final checksum = digest.toString();

          checksums.putIfAbsent(checksum, () => []).add(path);
        } catch (e) {
          print('Warning: Could not process file $path: $e');
        }
      }

      // Collect checksum groups with more than one path
      for (var entry in checksums.entries) {
        if (entry.value.length > 1) {
          duplicateChecksums[entry.key] = entry.value;
        }
      }
    }
  }

  return duplicateChecksums;
}

Future<void> printGroups(Map<String, List<String>> duplicateChecksums, bool findDupes) async {
  if (duplicateChecksums.isEmpty) {
    print('No duplicate files found.');
    return;
  }

  final sortedChecksums = duplicateChecksums.entries.toList()
    ..sort((a, b) => a.key.compareTo(b.key));

  for (var entry in sortedChecksums) {
    final String checksum = entry.key;
    final List<String> paths = entry.value.toList()..sort();

    if (findDupes) {
      for (var path in paths.skip(1)) {
        print('$path');
      }
    } else {
      print('\n--- Checksum: $checksum ---');
      for (var path in paths) {
        print('    $path');
      }
    }
  }
}

/// Recursively walks a directory, adding file paths to the size map.
Future<void> _walkDir(Directory dir, Map<int, List<String>> fileSizes) async {
  try {
    final List<FileSystemEntity> entities = await dir.list(recursive: false).toList();

    for (final entity in entities) {
      if (entity is Directory) {
        await _walkDir(entity, fileSizes);
      } else if (entity is File) {
        try {
          final fileSize = await entity.length();
          fileSizes.putIfAbsent(fileSize, () => []).add(entity.path);
        } catch (e) {
          print('Warning: Could not get size for ${entity.path}: $e');
        }
      }
    }
  } catch (e) {
    print('Warning: Could not list directory ${dir.path}: $e');
  }
}

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
    print("Usage: dart run bin/dupes.dart [--dupes] <directory1> [directory2] ...");
    return;
  }

  final Map<String, List<String>> duplicateChecksums = await groupFilesByChecksum(targetDirs);
  await printGroups(duplicateChecksums, findDupes);
}
