# gdlint

A simple linter for gdscript

## Goal

Support a variety of linting rules based on the [gdscript style guide](https://docs.godotengine.org/en/stable/tutorials/scripting/gdscript/gdscript_styleguide.html) in a configurable manor.

## Usage

There are 2 main entrypoints, a simple cli call and an LSP over stdin/stdout. To use the CLI you can build the application using [just](https://github.com/casey/just).

```sh
just build
```

The application will be built to `./bin/gdlint` and can then be used as follows:

```sh
./bin/gdlint check ./fixtures/scripts/class_declarations.gd
```

The `check` command also support `--include` and `--exclude` flags, run `./bin/gdlint check` to see the help.

If you want to have your editor (for example Zed) use `gdlint` you can add the following override in `./.zed/settings.json`.

```json
{
  "lsp": {
    "gdscript": {
      "binary": {
        "path": "./bin/gdlint",
        "arguments": ["lsp"]
      }
    }
  },
  "languages": {
    "GDScript": {
      "language_servers": [
        "gdscript"
      ]
    }
  }
}
```

When you open a gdscript file (like `./fixtures/scripts/typed_assignments.gd`) you should see warnings pop up in the file.

# Architecture

TODO

# Roadmap

## Rules

- [ ] detect missing type declarations\
  - [x] variables / constants
  - [x] function arguments
  - [x] function return types
  - [ ] signals
- [ ] detect name casing
  - [x] functions
  - [x] variables / constants
  - [x] classes
  - [ ] signals
  - [ ] files
  - [ ] enums + members
- [ ] ordering (top level)
- [ ] whitespace
  - [x] variable assignments
  - [ ] conditions
  - [ ] math operations
- [ ] formatting
  - [ ] line length
  - [ ] blank lines
  - [ ] array, enum, dict elements alignment
  - [ ] quote style
  - [ ] numbers
    - [ ] underscore usage
    - [ ] float decimal
    - [ ] hex caps
  - [ ] single line ternary
- [ ] pattern detection
  - [ ] simplified loops with higher level method
  - [ ] unused locals/privates
  - [ ] unreachable code

## Fixes

- [ ] ordering
- [ ] whitespace
- [ ] formatting
- [ ] pattern detection

## Capabilities

- [x] built in rule detection
- [x] cli check command
- [ ] cli fix command
- [x] lsp didOpen/didChange
- [x] lsp publishDiagnostics
- [ ] lsp didSave
- [ ] lsp code actions (workspace/applyEdits)
- [ ] lsp textDocument/formatting
- [ ] lsp workspace/diagnostics
- [ ] lsp workspace/didChangeWatchedFiles
- [ ] configuration file
  - [ ] ignored rules
  - [ ] excluded/included files
  - [ ] severity level per rule
  - [ ] overridable settings
    - [ ] line length
    - [ ] default string quote
- [ ] inline disable/enable comments
