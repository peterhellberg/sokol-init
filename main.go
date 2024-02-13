package main

import (
	"bytes"
	"embed"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

//go:embed all:content
var content embed.FS

type config struct {
	dir      string
	withMath bool
}

func main() {
	if err := run(os.Args, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run(args []string, stdout io.Writer) error {
	var cfg config

	flags := flag.NewFlagSet(args[0], flag.ExitOnError)

	flags.BoolVar(&cfg.withMath, "with-math", false, "Include src/math.zig")

	flags.Usage = func() {
		format := "Usage: %s [OPTION]... DIRECTORY\n\nOptions:\n"

		fmt.Fprintf(flags.Output(), format, os.Args[0])

		flags.PrintDefaults()
	}

	if err := flags.Parse(args[1:]); err != nil {
		return err
	}

	rest := flags.Args()

	// Require a directory name
	if len(rest) < 1 {
		return fmt.Errorf("no name given as the first argument")
	}

	cfg.dir = rest[0]

	// Make sure that dir does not already exist
	if _, err := os.Stat(cfg.dir); !os.IsNotExist(err) {
		return fmt.Errorf("%q already exists", cfg.dir)
	}

	// Create the dir and dir/src
	if err := os.MkdirAll(cfg.dir+"/src", os.ModePerm); err != nil {
		return err
	}

	// Enter the new directory
	if err := os.Chdir(cfg.dir); err != nil {
		return err
	}

	entries, err := content.ReadDir("content")
	if err != nil {
		return err
	}

	for _, e := range entries {
		if !e.IsDir() {
			if err := writeFile(cfg, e.Name(), replacer); err != nil {
				return err
			}

		} else {
			if e.Name() == "src" {
				srcEntries, err := content.ReadDir("content/src")
				if err != nil {
					return err
				}

				for _, e := range srcEntries {
					if !e.IsDir() {
						if e.Name() == "math.zig" && !cfg.withMath {
							continue
						}

						if err := writeFile(cfg, "src/"+e.Name(), replacer); err != nil {
							return err
						}
					}
				}
			}
		}
	}

	return nil
}

func writeFile(cfg config, name string, dataFuncs ...dataFunc) error {
	data, err := content.ReadFile("content/" + name)
	if err != nil {
		return fmt.Errorf("writeFile: %w", err)
	}

	for i := range dataFuncs {
		data = dataFuncs[i](cfg, name, data)
	}

	return os.WriteFile(name, data, 0644)
}

type dataFunc func(config, string, []byte) []byte

func replacer(cfg config, name string, data []byte) []byte {
	title := filepath.Base(cfg.dir)

	if len(title) < 2 {
		return data
	}

	switch name {
	case "build.zig", "build.zig.zon", "README.md", "src/main.zig":
		return replaceOne(data, "sokol-init", title)
	default:
		return data
	}
}

func replaceOne(data []byte, old, new string) []byte {
	return bytes.Replace(data, []byte(old), []byte(new), 1)
}
