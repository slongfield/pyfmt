#!/usr/bin/python
"""makedoc.py takes the README.md file and creates the doc.go file.

This makes keeping the Markdown and godoc in sync easy. It doesn't understand all aspects of
Markdown, or all aspects of Godoc, just the few that are used in this repo.

To update documentation, update the README.md, and then run this script to update the doc.go file.
Do not edit doc.go directly.
"""


def main():
    """Main."""
    doc = []

    for line in open("README.md"):
        if line.startswith("```"):
            continue
        if line.startswith("#"):
            line = line.lstrip("#").lstrip(" ")
        doc.append(line)

    with open("doc.go", 'w') as outp:
        outp.write("// Copyright header goes here.\n\n")
        outp.write("/*\n")
        for line in doc:
            outp.write(line)
        outp.write("*/\n")
        outp.write("package pyfmt\n")

if __name__ == '__main__':
    main()
