import 'dart:io';

import 'package:sorttrees/sorttrees.dart';

Future<void> main(List<String> args) async {
  if (args.contains('--help') || args.contains('-h')) {
    stdout.writeln('''
Sort an indented text tree.

Usage: sorttrees TreeFile [--SeparateTopLevel] [--InPlace]

Options:
  --SeparateTopLevel  Separate top level trees with a blank line.
  --InPlace           Write the output on top of the input file.
  -h, --help          Show this help message.
''');
    return;
  }

  final options = <String>{};
  final positional = <String>[];
  for (final arg in args) {
    if (arg == '--SeparateTopLevel' || arg == '--InPlace') {
      options.add(arg);
    } else if (arg.startsWith('-')) {
      stderr.writeln('Unknown option: $arg');
      stderr.writeln('Run "sorttrees --help" for usage.');
      exitCode = 2;
      return;
    } else {
      positional.add(arg);
    }
  }

  if (positional.length != 1) {
    stderr.writeln('Expected exactly one tree file.');
    stderr.writeln('Run "sorttrees --help" for usage.');
    exitCode = 2;
    return;
  }

  try {
    final file = File(positional.single);
    final lines = await readLines(file);
    final inPlace = options.contains('--InPlace');
    if (inPlace && areLinesSortedTree(lines)) {
      stdout.writeln('In place sorting requested but already sorted.');
      return;
    }

    if (inPlace) {
      await file.copy('${file.path}.bak');
    }

    final output = sortLines(
      lines,
      separateTopLevel: options.contains('--SeparateTopLevel'),
    );
    if (inPlace) {
      await file.writeAsString(output);
    } else {
      stdout.write('$output ');
    }
  } on FileSystemException catch (error) {
    stderr.writeln(error.message);
    exitCode = 1;
  } on FormatException catch (error) {
    stderr.writeln(error.message);
    exitCode = 1;
  }
}
