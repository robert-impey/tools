import 'dart:io';

import 'package:sorttrees/randomtrees.dart';
import 'package:sorttrees/sorttrees.dart';

void main(List<String> args) {
  if (args.contains('--help') || args.contains('-h')) {
    stdout.writeln('''
Generate random trees and sort them.
Time the sorting part.
Print statistics on the times.

Usage: speedtest [--Depth <depth>] [--Items <items>] [--Length <length>] [--Alphabet <alphabet>] [--Iterations <iterations>]

Options:
  --Depth       The depth of the trees. (default: 3)
  --Items       The number of items for each node of the tree. (default: 10)
  --Length      The length of each item. (default: 8)
  --Alphabet    The alphabet of allowed characters.
  --Iterations  The number of iterations. (default: 100)
  -h, --help    Show this help message.
''');
    return;
  }

  var depth = 3;
  var items = 10;
  var length = 8;
  var alphabet = defaultAlphabet;
  var iterations = 100;

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
    } else if (arg == '--Iterations') {
      final val = nextValue();
      if (val == null || int.tryParse(val) == null) {
        _usageError('Expected integer after --Iterations');
        return;
      }
      iterations = int.parse(val);
    } else if (arg.startsWith('--Iterations=')) {
      final val = arg.substring('--Iterations='.length);
      final parsed = int.tryParse(val);
      if (parsed == null) {
        _usageError('Invalid integer for --Iterations: $val');
        return;
      }
      iterations = parsed;
    } else {
      _usageError('Unknown option: $arg');
      return;
    }
  }

  final timesInSeconds = <double>[];
  var remainingIterations = iterations;

  while (remainingIterations > 0) {
    remainingIterations--;

    final randomTreeLines = generateRandomTreeLines(
      depth: depth,
      items: items,
      length: length,
      alphabet: alphabet,
    );

    final stopwatch = Stopwatch()..start();

    sortLines(randomTreeLines);

    stopwatch.stop();

    timesInSeconds.add(stopwatch.elapsedMicroseconds / 1000000.0);
  }

  if (timesInSeconds.isEmpty) {
    stdout.writeln('Number of times: 0');
    return;
  }

  var minTime = timesInSeconds.first;
  var maxTime = timesInSeconds.first;
  var sumTime = 0.0;

  for (final t in timesInSeconds) {
    if (t < minTime) minTime = t;
    if (t > maxTime) maxTime = t;
    sumTime += t;
  }

  final meanTime = sumTime / timesInSeconds.length;

  stdout.writeln('Number of times: ${timesInSeconds.length}');
  stdout.writeln('Min time: $minTime s');
  stdout.writeln('Max time: $maxTime s');
  stdout.writeln('Mean time: $meanTime s');
}

void _usageError(String message) {
  stderr.writeln(message);
  stderr.writeln('Run "speedtest --help" for usage.');
  exitCode = 2;
}
