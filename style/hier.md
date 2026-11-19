<!--
  Copyright 2026 Samvel Khalatyan. All rights reserved.

  Use of this source code is governed by a BSD-style
  license that can be found in the LICENSE file.
-->

# NAME

*hier* - hier(7) style guide

# DESCRIPTION

## Folder naming: path as a sentence

A path reads left to right as a phrase. Each folder names the *kind* of thing
below it, never how many there are. It is name qualifier for the next item in
the path.

"the ‹folder› ‹name›" should read naturally.

- *pass*: `image/foo.iso` gives "the image foo.iso,"
- *fail*: `images/foo.iso` gives "the images foo.iso,".

### Rules

1. *Countable things*: use the singular noun.  
   `image/`, `disk/`, `script/`, `test/`, `doc/`, `user/`

2. *Mass nouns*: stay as they are, they have no plural to avoid.  
   `data/`, `config/`, `src/`, `media/`, `cache/`, `log/`

3. *Roles and stages*: use the noun for the role or stage.  
   `build/`, `dist/`, `tmp/`, `archive/`, `backup/`

4. *General to specific*: like English compound nouns. A modifier comes
   before the thing it modifies, so the path reads like a noun phrase:

   - `vm/image/debian.qcow2` reads "the vm image debian."
   - `db/backup/2026-11-19.sql.gz` reads "the db backup from 2026-11-19."

5. *Path holds the kind*: don't repeat it in the file name, the folder already
   says it. Use `image/debian.iso`, not `image/debian-image.iso`.

6. *Avoid plurals*: pick a singular equivalent:
   - `settings/` becomes `config/`
   - `news/` becomes `article/`
   - `assets/` becomes `asset/` (or split into `image/`, `font/`, and so on).

7. *One form per concept*: if you use `doc/`, never also use `documentation/`
   or `docs/`. Short forms like `src`, `lib`, `bin`, `doc` and `etc` are fine if
   used consistently.

8. *Names you don't control are exceptions*: leave tool-mandated or
   third-party names alone (`node_modules/`, `.git/`, a `tests/` that your test
   runner requires), and list them as known exceptions so nobody "fixes" them.

9. *Formatting*: lowercase, with hyphens if you really need several words
   (`vm-image/`). Prefer nesting (`vm/image/`) when the parts are a modifier
   plus a kind, because nesting keeps the sentence reading.

### Example

```
/srv
├── vm/
│   ├── image/        debian-13.qcow2, freebsd-15.qcow2
│   └── disk/         web01.raw
├── db/
│   └── backup/       2026-11-19.sql.gz
├── web/
│   ├── config/       nginx.conf
│   └── log/          access.log
└── tmp/
```

`/srv/vm/image/freebsd-15.qcow2` reads "under srv, the vm image freebsd-15.".
This sits comfortably next to FreeBSD's own `/usr/local/etc` and
`/usr/local/lib`.
