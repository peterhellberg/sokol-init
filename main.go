package main

import (
	"bufio"
	"bytes"
	"embed"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

//go:embed all:content
var content embed.FS

type config struct {
	dir      string
	withMath bool
	zon      ZON
}

type ZON struct {
	name        string
	fingerprint string
}

func main() {
	if err := run(os.Args, os.Stderr); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func parse(args []string, stderr io.Writer) (config, error) {
	var cfg config

	flags := flag.NewFlagSet(args[0], flag.ExitOnError)

	flags.SetOutput(stderr)

	flags.BoolVar(&cfg.withMath, "with-math", false, "Include src/math.zig")

	flags.Usage = func() {
		format := "Usage: %s [OPTION]... DIRECTORY\n\nOptions:\n"

		fmt.Fprintf(flags.Output(), format, os.Args[0])

		flags.PrintDefaults()
	}

	if err := flags.Parse(args[1:]); err != nil {
		return cfg, err
	}

	rest := flags.Args()

	// Require a directory name
	if len(rest) < 1 {
		return cfg, fmt.Errorf("no name given as the first argument")
	}

	cfg.dir = rest[0]

	zon, err := initZON(cfg.dir)
	if err != nil {
		return cfg, err
	}

	cfg.zon = zon

	return cfg, nil
}

func run(args []string, stderr io.Writer) error {
	cfg, err := parse(args, stderr)
	if err != nil {
		return err
	}

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

	return os.WriteFile(name, data, 0o644)
}

type dataFunc func(config, string, []byte) []byte

func replacer(cfg config, name string, data []byte) []byte {
	title := filepath.Base(cfg.dir)

	if len(title) < 2 {
		return data
	}

	switch name {
	case "build.zig", "README.md", "src/main.zig":
		return replaceOne(data, "sokol-init", title)
	case "build.zig.zon":
		data = replaceOne(data, ".sokol_init", cfg.zon.name)
		data = replaceOne(data, "0xaf14342a89e8115e", cfg.zon.fingerprint)
		return data
	default:
		return data
	}
}

func replaceOne(data []byte, old, new string) []byte {
	return bytes.Replace(data, []byte(old), []byte(new), 1)
}

func initZON(dir string) (ZON, error) {
	tmp, err := os.MkdirTemp("", "sokol-init-")
	if err != nil {
		return ZON{}, err
	}
	defer os.RemoveAll(tmp)

	cwd, err := os.Getwd()
	if err != nil {
		return ZON{}, err
	}
	defer os.Chdir(cwd)

	tmpDir := filepath.Join(tmp, dir)

	if err := os.Mkdir(tmpDir, 0o755); err != nil {
		return ZON{}, err
	}

	if err := os.Chdir(tmpDir); err != nil {
		return ZON{}, err
	}

	cmd := exec.Command("zig", "init")

	if err := cmd.Run(); err != nil {
		return ZON{}, err
	}

	zonPath := filepath.Join(tmpDir, "build.zig.zon")

	return extractZON(zonPath)
}

func extractZON(zonPath string) (ZON, error) {
	var zon ZON

	f, err := os.Open(zonPath)
	if err != nil {
		return zon, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())

		if prefix := ".name = "; strings.Contains(text, prefix) {
			zon.name = strings.TrimSuffix(strings.TrimPrefix(text, prefix), ",")
		}

		if prefix := ".fingerprint = "; strings.Contains(text, prefix) {
			fingerprint, _, _ := strings.Cut(strings.TrimPrefix(text, prefix), ",")
			zon.fingerprint = fingerprint
		}
	}

	return zon, nil
}
