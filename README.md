# magpie

**Snatches your shiny thoughts into a daily nest.**

Magpie is a small command-line tool for jotting down notes without leaving the terminal. 
Type `magpie` followed by your note and it's saved with the time you added it. 
Notes are grouped by day: each day gets its own nest, which is a plain text file that you can view, trim or delete with a single command.

## Usage

```
magpie [-h|--help] [-v|--version] <note>|pluck|nest|renest|nests|abandon
```

## Commands

| Command   | Description                                              |
|-----------|----------------------------------------------------------|
| `<note>`  | Save a note, thought or idea to today's nest file        |
| `pluck`   | Delete the last note from today's nest                   |
| `nest`    | Print notes from today's nest                            |
| `renest`  | Delete today's nest file and all its contents            |
| `nests`   | Print notes from all existing nest files                 |
| `abandon` | Delete all nest files and all their contents             |

Running `magpie` with no arguments prints the help message.

### `<note>` — save a note

Anything that isn't a command is saved as a note in today's nest. Quotes are optional.

```sh
$ magpie buy milk
Snatched:
09:14:02  buy milk

$ magpie "refactor the parser, it's getting messy"
Snatched:
11:47:39  refactor the parser, it's getting messy
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

Empty nests are skipped. If there are no notes at all: `No shiny notes in any nest yet.`

### `pluck` — remove the last note from today's nest

```sh
$ magpie pluck
Plucked:
11:47:39  refactor the parser, it's getting messy
```

If today's nest is empty: `No shiny notes in the nest to pluck.`

### `renest` — delete today's nest

Asks for confirmation first. Only `y`, `yes` or `Yes` proceeds.

```sh
$ magpie renest
This action will delete current nest file. Proceed? [y/N]: y
Renested.
```

Anything else aborts with `Aborted. Nest kept.`

### `abandon` — delete all nests

Asks for confirmation first. Only `y` or `yes` or `Yes` proceeds.

```sh
$ magpie abandon
This action will delete all nest files. Proceed? [y/N]: y
All nests abandoned.
```

Anything else aborts with `Aborted. All nests kept.`

## Options

| Option            | Description             |
|-------------------|-------------------------|
| `-h`, `--help`    | Print the help message  |
| `-v`, `--version` | Print the version       |

Output is coloured with ANSI escape codes; errors are printed in red to stderr.

## Where notes are stored

Nests are plain text files stored in the `HOME` directory and under `.magpie` directory, i.e. `~/.magpie/`.
Each day gets its own nest. Every nest file is saved in the following format: `YYYY-MM-DD-nest.txt`

## License
MIT. See [LICENSE](LICENSE)
