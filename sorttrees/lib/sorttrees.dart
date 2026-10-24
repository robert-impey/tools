import 'dart:io';

const int _indentationWidth = 4;

class _Tree {
  _Tree(this.text);

  final String? text;
  final List<_Tree> children = [];
}

class _ParsedLine {
  const _ParsedLine(this.depth, this.text);

  final int depth;
  final String text;
}

Future<List<String>> readLines(File file) async {
  final contents = await file.readAsString();
  return contents
      .split('\n')
      .map((line) => line.trimRight())
      .where((line) => line.isNotEmpty)
      .toList();
}

String sortLines(List<String> lines, {bool separateTopLevel = false}) {
  final root = _Tree(null);
  final stack = <_Tree>[root];

  for (var index = 0; index < lines.length; index++) {
    final parsed = _parseLine(lines[index]);
    if (parsed.depth >= stack.length) {
      throw FormatException('Line ${index + 1} skips indentation depth.');
    }

    final parent = stack[parsed.depth];
    final child = _Tree(parsed.text);
    parent.children.add(child);
    stack
      ..removeRange(parsed.depth + 1, stack.length)
      ..add(child);
  }

  _finalise(root);
  return _writeTree(root, separateTopLevel: separateTopLevel);
}

bool areLinesSortedTree(List<String> lines) {
  final previousItemsByDepth = <String>[];
  var previousLineDepth = 0;

  for (var index = 0; index < lines.length; index++) {
    final parsed = _parseLine(lines[index]);
    if (parsed.depth < previousLineDepth) {
      if (parsed.depth >= previousItemsByDepth.length) {
        throw FormatException(
          'Line ${index + 1} has invalid indentation depth.',
        );
      }
      previousItemsByDepth.removeRange(
        parsed.depth + 1,
        previousItemsByDepth.length,
      );
    }

    if (parsed.depth == previousItemsByDepth.length) {
      previousItemsByDepth.add(parsed.text);
    } else if (parsed.depth < previousItemsByDepth.length) {
      if (parsed.text.compareTo(previousItemsByDepth[parsed.depth]) < 0) {
        return false;
      }
      previousItemsByDepth[parsed.depth] = parsed.text;
    } else {
      throw FormatException('Line ${index + 1} skips indentation depth.');
    }

    previousLineDepth = parsed.depth;
  }

  return true;
}

_ParsedLine _parseLine(String line) {
  var indentation = 0;
  while (indentation < line.length && line.codeUnitAt(indentation) == 0x20) {
    indentation++;
  }
  return _ParsedLine(
    indentation ~/ _indentationWidth,
    line.substring(indentation).trimRight(),
  );
}

void _finalise(_Tree node) {
  for (final child in node.children) {
    _finalise(child);
  }
  node.children.sort(_compareTrees);
}

int _compareTrees(_Tree left, _Tree right) {
  final textComparison = left.text!.compareTo(right.text!);
  if (textComparison != 0) {
    return textComparison;
  }
  return _treeString(left).compareTo(_treeString(right));
}

String _treeString(_Tree node) => _writeTree(node, isRoot: false);

String _writeTree(
  _Tree node, {
  String indentation = '',
  bool separateTopLevel = false,
  bool isRoot = true,
}) {
  final output = StringBuffer();
  var childIndentation = indentation;
  if (!isRoot) {
    output
      ..write(indentation)
      ..writeln(node.text);
    childIndentation += ' ' * _indentationWidth;
  }

  for (var index = 0; index < node.children.length; index++) {
    if (isRoot && separateTopLevel && index > 0) {
      output.writeln();
    }
    output.write(
      _writeTree(
        node.children[index],
        indentation: childIndentation,
        separateTopLevel: separateTopLevel,
        isRoot: false,
      ),
    );
  }
  return output.toString();
}
