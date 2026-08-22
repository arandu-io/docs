#!/usr/bin/env python3
"""Check that every internal reference in this repository resolves.

This repository is markdown and nothing else, so the links between the files
are the only thing a machine can verify here. Two kinds of reference exist,
and they are not the same check:

  links     -- real markdown links and images, `[text](target)`. A broken one
               is a defect and this script exits 1 for it.

  pointers  -- a file named in prose, `docs/07-pacotes.md` or `docs/31:398`.
               Most of these do not resolve, because the numbered documents
               moved into _archive/ and the pointers were written before the
               move. A pointer inside _archive/ is counted and does not fail
               the run: the archive is frozen, so correcting one would mean
               editing a file that no longer takes edits. A pointer anywhere
               else does fail, because there the file can be corrected.

Standard library only, on purpose: nothing in this project installs a package
manager to run a check.

Usage:
    python3 check-links.py            report both, exit 1 on a broken link
    python3 check-links.py --quiet    same, printing only what fails
"""

import collections
import os
import re
import subprocess
import sys

# A markdown link or image, with an optional title. Reference definitions
# ("[label]: target") are matched separately because they anchor at the line.
INLINE = re.compile(r'!?\[(?P<text>[^\]]*)\]\((?P<url>[^)\s]+)(?:\s+"[^"]*")?\)')
REFDEF = re.compile(r'^\s{0,3}\[(?P<text>[^\]]+)\]:\s*(?P<url>\S+)')
# Markdown allows html, and the README uses it for the logo and the badges.
HTML = re.compile(r'<(?:img|a)\s[^>]*?(?:src|href)\s*=\s*"(?P<url>[^"]+)"')
EXTERNAL = re.compile(r'^(https?:|mailto:|tel:|ftp:)')
CODE_SPAN = re.compile(r'`[^`]*`')
FENCE = re.compile(r'^\s*(```|~~~)')

# A file named in prose. The second form is the shorthand the numbered
# documents use for each other -- "docs/31:398" is document 31, line 398.
POINTER = re.compile(r'(?<![\w/-])((?:[\w.-]+/)*[\w.-]+\.md)(?![\w])')
SHORTHAND = re.compile(r'(?<![\w/-])docs/(\d{2})(?![\w.-])')

ARCHIVE = '_archive/'


def tracked_files():
    out = subprocess.run(['git', 'ls-files'], capture_output=True, text=True,
                         check=True).stdout
    return set(out.split())


def headings(path):
    """The anchors GitHub derives from the headings of one file."""
    slugs = set()
    for line in read(path):
        m = re.match(r'^#{1,6}\s+(.*?)\s*$', line)
        if not m:
            continue
        text = re.sub(r'`([^`]*)`', r'\1', m.group(1))
        text = re.sub(r'\[([^\]]*)\]\([^)]*\)', r'\1', text)
        slug = re.sub(r'[^\w\- ]', '', text.lower()).strip().replace(' ', '-')
        slugs.add(slug)
    return slugs


def read(path):
    with open(path, encoding='utf-8') as fh:
        return fh.read().split('\n')


def strip_code(line):
    """Blank out inline code so that Go generics do not read as a link."""
    return CODE_SPAN.sub(lambda m: ' ' * len(m.group(0)), line)


def scan(files, tracked):
    """Return (links checked, broken links, unresolved pointers)."""
    broken, pointers = [], []
    anchors = {}
    checked = 0
    for path in files:
        directory = os.path.dirname(path)
        fenced = False
        for number, raw in enumerate(read(path), 1):
            if FENCE.match(raw):
                fenced = not fenced
                continue
            if fenced:
                continue
            line = strip_code(raw)

            found = [m.group('url') for m in INLINE.finditer(line)]
            found += [m.group('url') for m in HTML.finditer(line)]
            definition = REFDEF.match(line)
            if definition:
                found.append(definition.group('url'))
            for url in found:
                if EXTERNAL.match(url) or url.startswith('#'):
                    continue
                target, _, fragment = url.partition('#')
                target = target.split('?')[0]
                if not target:
                    continue
                checked += 1
                resolved = os.path.normpath(os.path.join(directory, target))
                if resolved not in tracked:
                    broken.append((path, number, url, 'no such file'))
                elif fragment and resolved.endswith('.md'):
                    if resolved not in anchors:
                        anchors[resolved] = headings(resolved)
                    if fragment.lower() not in anchors[resolved]:
                        broken.append((path, number, url, 'no such heading'))

            for m in POINTER.finditer(raw):
                named = m.group(1)
                candidates = [os.path.normpath(os.path.join(directory, named)),
                              os.path.normpath(named)]
                if not any(c in tracked for c in candidates):
                    pointers.append((path, number, named, guess(named, tracked)))
            for m in SHORTHAND.finditer(raw):
                pointers.append((path, number, m.group(0),
                                 guess(m.group(1) + '-', tracked)))
    return checked, broken, pointers


def guess(named, tracked):
    """Where the named file went, when exactly one candidate carries the name."""
    base = os.path.basename(named)
    hits = [t for t in tracked if os.path.basename(t) == base] or \
           [t for t in tracked if os.path.basename(t).startswith(base)]
    return hits[0] if len(hits) == 1 else ''


def main():
    quiet = '--quiet' in sys.argv
    os.chdir(subprocess.run(['git', 'rev-parse', '--show-toplevel'],
                            capture_output=True, text=True,
                            check=True).stdout.strip())
    tracked = tracked_files()
    files = sorted(f for f in tracked if f.endswith('.md'))

    checked, broken, pointers = scan(files, tracked)

    outside = [p for p in pointers if not p[0].startswith(ARCHIVE)]
    inside = [p for p in pointers if p[0].startswith(ARCHIVE)]

    if not quiet:
        print(f'{len(files)} markdown files, {checked} internal links')
        print(f'{len(pointers)} pointers name a file that is not there: '
              f'{len(inside)} inside {ARCHIVE}, {len(outside)} outside it')
        by_file = collections.Counter(p[0] for p in pointers)
        for path, count in by_file.most_common(10):
            print(f'    {count:4}  {path}')
        if len(by_file) > 10:
            print(f'    and {len(by_file) - 10} more files')
        print()
        print(f'{ARCHIVE} is frozen, so the pointers written inside it stay as '
              f'they are.\nThey are counted here and they do not fail the run.')
        print()

    for path, number, name, moved in outside:
        where = f', moved to {moved}' if moved else ''
        print(f'{path}:{number}: pointer names {name}, which is not there{where}')

    for path, number, url, why in broken:
        print(f'{path}:{number}: {why}: {url}')

    if broken or outside:
        print(f'\n{len(broken)} broken link(s), '
              f'{len(outside)} unresolved pointer(s) outside {ARCHIVE}.')
        return 1
    if not quiet:
        print('every markdown link resolves.')
    return 0


if __name__ == '__main__':
    sys.exit(main())
