# sokol-init :sparkles:

This is a command line tool that helps you start a
[Sokol](https://github.com/floooh/sokol) project with
[Zig](https://ziglang.org/) :zap:

`sokol-init` is used to create a directory containing code that
allows you to promptly get started coding with Sokol (via the
[sokol-zig](https://github.com/floooh/sokol-zig) bindings)

## Installation

(Requires you to have [Go](https://go.dev/) installed)

```sh
go install github.com/peterhellberg/sokol-init@latest
```

## Usage

(Requires you to have an up to date (_nightly_) version of
[Zig](https://ziglang.org/download/#release-master) installed

```sh
sokol-init myproject
cd myproject
zig build run
```

:seedling:
