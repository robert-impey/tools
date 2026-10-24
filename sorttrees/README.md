# sorttrees

`sorttrees` sorts space-indented text as a tree. Each level uses four leading
spaces, and siblings are sorted alphabetically. Blank lines are removed.

```text
dart run bin/sorttrees.dart TreeFile [--SeparateTopLevel] [--InPlace]
```

- `--SeparateTopLevel` adds a blank line between top-level trees.
- `--InPlace` replaces the input file after creating a `TreeFile.bak` backup.

Without `--InPlace`, sorted output is written to standard output.
