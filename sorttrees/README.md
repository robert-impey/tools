# sorttrees

This Dart program was ported from the Python code in [tree-sorter](https://github.com/robert-impey/tree-sorter).

Program for sorting space indented text, treating blocks indented to the same level as groups. 
This is useful for grouping items in a "to do" list.

e.g.

    II Animals
      2 Dogs
      3 Cats
      1 Fish
    I Plants
      b Nettles
      c Oak
      a Grass
  
becomes

    I Plants
      a Grass
      b Nettles
      c Oak
    II Animals
      1 Fish
      2 Dogs
      3 Cats
      
Trees of grouped data can be of arbitrary depth.

## Usage

    $ sorttrees -h
    Sort an indented text tree.
    
    Usage: sorttrees TreeFile [--SeparateTopLevel] [--InPlace]
    
    Options:
      --SeparateTopLevel  Separate top level trees with a blank line.
      --InPlace           Write the output on top of the input file.
      -h, --help          Show this help message.

### Examples

The input file might contain a few trees in no particular order:

    $ cat fixtures/two-deep-with-gaps.txt
    1 one
        i one
        iii one
        ii one

    3 three
        iii three
        i three
        ii three

    2 two
        ii two
        iii two
        i two

The tool reads the contents of the text file and prints out the trees in order.

    $ sorttrees fixtures/two-deep-with-gaps.txt
    1 one
        i one
        ii one
        iii one
    2 two
        i two
        ii two
        iii two
    3 three
        i three
        ii three
        iii three

The --SeparateTopLevel puts a blank line between "top level" trees. This is useful for navigating between groups nimbly in vim.

    $ sorttrees fixtures/two-deep-with-gaps.txt --SeparateTopLevel
    1 one
        i one
        ii one
        iii one

    2 two
        i two
        ii two
        iii two

    3 three
        i three
        ii three
        iii three
        
The --InPlace flag writes the sorted trees on top of the input file. 

    $ sorttrees ~/docs/todo.txt --SeparateTopLevel --InPlace

In this case, the tool creates a back up file with the extension ".bak" first. Obviously, you've already got other back ups anyway. This is free software. I take no respobsibility for anything that happens to your files. See license.

## License

This software is released under the terms of the [MIT License](https://opensource.org/licenses/MIT).