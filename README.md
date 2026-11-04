# magpie

Collects your shiny thoughts into a daily nest.

magpie is a small command-line tool for jotting down notes without leaving the terminal. 
Type `magpie` followed by your note and it's saved with the time you added it. 
Notes are grouped by day: each day gets its own nest, a plain text file, which you can view, trim or delete with a single command.

## Usage

```
magpie [-h|--help] [-v|--version] <note>|pluck|nest|renest|nests|abandon
```

## Commands

| Command   | Description                                              |
|-----------|----------------------------------------------------------|
| `<note>`  | Save a note, thought or idea to today's nest file        |
| `pluck`   | Delete the last note from today's nest                  |
| `nest`    | Print notes from today's nest            |
| `renest`  | Delete today's nest file and all its contents |
| `nests`   | Print notes from all existing nest files                 |
| `abandon` | Delete all nest files and all their contents                    |

### `<note>` — save a note

Anything that isn't a command is saved as a note in today's nest. Quotes are optional.

```sh
$ magpie buy milk
$ magpie "refactor the parser, it's getting messy"
```

> Tip: quote notes with shell characters like `?`, `!`, `*` or `'`, e.g. `magpie "why is it slow?"`.
>
> A note that is exactly a command name (e.g. `magpie nest`) runs the command instead of being saved.

### `nest` — show today's notes

```sh
$ magpie nest
09:14:02  buy milk
11:47:39  refactor the parser, it's getting messy
```

If there are no notes yet: `No shiny notes in the nest yet.`

### `nests` — show notes from all days

```sh
$ magpie nests
======== 2026-11-03 ========
16:20:11  call the dentist

======== 2026-11-04 ========
09:14:02  buy milk
11:47:39  refactor the parser, it's getting messy
```

Empty nests are skipped.

### `pluck` — remove the last note from today's nest

```sh
$ magpie pluck
$ magpie nest
09:14:02  buy milk
```

### `renest` — delete today's nest

Asks for confirmation first. Only `y` or `yes` proceeds.

```sh
$ magpie renest
This action will delete current nest file. Proceed? [y/N]: y
Renested.
```

### `abandon` — delete all nests

Asks for confirmation first. Only `y` or `yes` proceeds.

```sh
$ magpie abandon
This action will delete all nest files. Proceed? [y/N]: y
All nests abandoned.
```

## Options

| Option            | Description             |
|-------------------|-------------------------|
| `-h`, `--help`    | Print the help message  |
| `-v`, `--version` | Print the version       |

## Where notes are stored

Nests are plain text files in `~/.magpie/`, one per day:

```
~/.magpie/2026-11-04-nest.txt
```

## License
MIT. See [LICENSE](LICENSE)
