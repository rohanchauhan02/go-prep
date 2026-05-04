#!/usr/bin/env python3
"""Fix Go files where \\n inside string/rune literals became real newlines."""
import os, re, glob

base = "/Users/rohanchauhan/go_interview_prep"

def fix_file(path):
    with open(path) as f:
        content = f.read()

    lines = content.split('\n')
    fixed = []
    in_string = False
    quote_char = None
    result_lines = []

    # Strategy: join lines that are "broken" inside a string literal
    # We detect this by tracking whether we're inside an unclosed " or `
    i = 0
    while i < len(lines):
        line = lines[i]
        # Count unescaped quotes to detect unclosed strings
        j = 0
        in_dq = False
        while j < len(line):
            c = line[j]
            if c == '\\' and in_dq:
                j += 2
                continue
            if c == '"':
                in_dq = not in_dq
            j += 1

        if in_dq:
            # Line has an unclosed double-quote — merge with next line
            merged = line
            while in_dq and i + 1 < len(lines):
                i += 1
                next_line = lines[i].strip()
                merged = merged + '\\n' + next_line
                # Recheck
                j = 0
                in_dq2 = False
                for c2 in merged:
                    if c2 == '"':
                        in_dq2 = not in_dq2
                in_dq = in_dq2
            result_lines.append(merged)
        else:
            result_lines.append(line)
        i += 1

    fixed_content = '\n'.join(result_lines)
    with open(path, 'w') as f:
        f.write(fixed_content)

# Find all main.go files and fix them
files = glob.glob(os.path.join(base, "q*/main.go"))
for path in sorted(files):
    fix_file(path)
    folder = os.path.basename(os.path.dirname(path))
    print(f"  fixed {folder}/main.go")

print("\nDone fixing string literals.")
