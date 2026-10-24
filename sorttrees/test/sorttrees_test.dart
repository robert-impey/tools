import 'dart:io';

import 'package:sorttrees/sorttrees.dart';
import 'package:test/test.dart';

void main() {
  group('sortLines', () {
    test('sorts nested trees and removes gaps', () {
      expect(
        sortLines(['2', '    2', '    1', '1', '    2', '    1']),
        '1\n    1\n    2\n2\n    1\n    2\n',
      );
    });

    test('separates top-level trees', () {
      expect(sortLines(['2', '1'], separateTopLevel: true), '1\n\n2\n');
    });

    test('supports an empty tree', () {
      expect(sortLines([]), isEmpty);
    });

    test('rejects skipped indentation depths', () {
      expect(
        () => sortLines(['root', '        grandchild']),
        throwsFormatException,
      );
    });
  });

  group('areLinesSortedTree', () {
    test('accepts empty and sorted trees', () {
      expect(areLinesSortedTree([]), isTrue);
      expect(areLinesSortedTree(['root', '    a', '    b']), isTrue);
    });

    test('rejects unsorted siblings', () {
      expect(areLinesSortedTree(['root', '    b', '    a']), isFalse);
    });

    test('checks ordering after dropping a depth', () {
      expect(
        areLinesSortedTree([
          'root',
          '    a',
          '        child',
          '    b',
          'next',
          '    z',
          '    c',
        ]),
        isFalse,
      );
    });
  });

  test('readLines removes blank lines and trailing whitespace', () async {
    final directory = await Directory.systemTemp.createTemp('sorttrees-test-');
    addTearDown(() => directory.delete(recursive: true));
    final file = File('${directory.path}${Platform.pathSeparator}tree.txt');
    await file.writeAsString('one  \r\n\r\n    two \t\r\n');

    expect(await readLines(file), ['one', '    two']);
  });
}
