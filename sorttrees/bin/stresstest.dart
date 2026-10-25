import 'dart:io';

import 'package:sorttrees/randomtrees.dart';
import 'package:sorttrees/sorttrees.dart';

void main(List<String> args) {
  if (args.contains('--help') || args.contains('-h')) {
    stdout.writeln('''
Generate random tree lines and attempt to sort them.
Test that the output lines are in the correct order.
Repeat until failure, or until the timeout is reached.

Usage: stresstest [--Depth <depth>] [--Items <items>] [--Length <length>] [--Alphabet <alphabet>] [--Duration <seconds>]

Options:
  --Depth       The depth of the trees. (default: 3)
  --Items       The number of items for each node of the tree. (default: 10)
  --Length      The length of each item. (default: 8)
  --Alphabet    The alphabet of allowed characters.
  --Duration    The duration in seconds to run before finishing successfully. (default: 10)
  -h, --help    Show this help message.
''');
    return;
  }

  var depth = 3;
  var items = 10;
  var length = 8;
  var alphabet = defaultAlphabet;
  var durationInSeconds = 10.0;

  for (var i = 0; i < args.length; i++) {
    final arg = args[i];
    String? nextValue() {
      if (i + 1 < args.length && !args[i + 1].startsWith('-')) {
        return args[++i];
      }
      return null;
    }

    if (arg == '--Depth') {
      final val = nextValue();
      if (val == null || int.tryParse(val) == null) {
        _usageError('Expected integer after --Depth');
        return;
      }
      depth = int.parse(val);
    } else if (arg.startsWith('--Depth=')) {
      final val = arg.substring('--Depth='.length);
      final parsed = int.tryParse(val);
      if (parsed == null) {
        _usageError('Invalid integer for --Depth: $val');
        return;
      }
      depth = parsed;
    } else if (arg == '--Items') {
      final val = nextValue();
      if (val == null || int.tryParse(val) == null) {
        _usageError('Expected integer after --Items');
        return;
      }
      items = int.parse(val);
    } else if (arg.startsWith('--Items=')) {
      final val = arg.substring('--Items='.length);
      final parsed = int.tryParse(val);
      if (parsed == null) {
        _usageError('Invalid integer for --Items: $val');
        return;
      }
      items = parsed;
    } else if (arg == '--Length') {
      final val = nextValue();
      if (val == null || int.tryParse(val) == null) {
        _usageError('Expected integer after --Length');
        return;
      }
      length = int.parse(val);
    } else if (arg.startsWith('--Length=')) {
      final val = arg.substring('--Length='.length);
      final parsed = int.tryParse(val);
      if (parsed == null) {
        _usageError('Invalid integer for --Length: $val');
        return;
      }
      length = parsed;
    } else if (arg == '--Alphabet') {
      final val = nextValue();
      if (val == null) {
        _usageError('Expected string after --Alphabet');
        return;
      }
      alphabet = val;
    } else if (arg.startsWith('--Alphabet=')) {
      alphabet = arg.substring('--Alphabet='.length);
    } else if (arg == '--Duration') {
      final val = nextValue();
      if (val == null || double.tryParse(val) == null) {
        _usageError('Expected number after --Duration');
        return;
      }
      durationInSeconds = double.parse(val);
    } else if (arg.startsWith('--Duration=')) {
      final val = arg.substring('--Duration='.length);
      final parsed = double.tryParse(val);
      if (parsed == null) {
        _usageError('Invalid number for --Duration: $val');
        return;
      }
      durationInSeconds = parsed;
    } else {
      _usageError('Unknown option: $arg');
      return;
    }
  }

  final stopwatch = Stopwatch()..start();
  final targetMicroseconds = (durationInSeconds * 1000000).round();
  var iterationsCompleted = 0;

  while (stopwatch.elapsedMicroseconds < targetMicroseconds) {
    final randomTreeLines = generateRandomTreeLines(
      depth: depth,
      items: items,
      length: length,
      alphabet: alphabet,
    );

    final sortedOutput = sortLines(randomTreeLines);
    final sortedLines = sortedOutput
        .split('\n')
        .where((line) => line.isNotEmpty)
        .toList();

    if (!areLinesSortedTree(sortedLines)) {
      stderr.writeln('A tree was not sorted correctly');
      stderr.writeln('Input');
      for (final line in randomTreeLines) {
        stderr.writeln(line);
      }
      stderr.writeln('Output');
      for (final line in sortedLines) {
        stderr.writeln(line);
      }
      exitCode = 1;
      return;
    }

    iterationsCompleted++;
  }

  stdout.writeln(
    'Stress test completed successfully ($iterationsCompleted iterations in ${stopwatch.elapsed.inMilliseconds / 1000.0}s).',
  );
}

void _usageError(String message) {
  stderr.writeln(message);
  stderr.writeln('Run "stresstest --help" for usage.');
  exitCode = 2;
}
