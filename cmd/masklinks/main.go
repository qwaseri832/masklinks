package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/qwaseri832/masklinks"
)

func main() {
	cfg, err := parseFlags(os.Args[0], os.Args[1:], os.Stderr)
	if err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return
		}
		os.Exit(2)
	}

	if err := run(cfg, os.Stdin, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "не удалось замаскировать текст:", err)
		os.Exit(1)
	}
}

type config struct {
	schemes []string
	text    []string
}

func parseFlags(prog string, args []string, stderr io.Writer) (config, error) {
	fs := flag.NewFlagSet(prog, flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.Usage = func() {
		fmt.Fprintf(stderr, "masklinks скрывает адреса ссылок, оставляя видимой только схему.\n\n")
		fmt.Fprintf(stderr, "Использование:\n  %s [флаги] [текст...]\n\n", prog)
		fmt.Fprintf(stderr, "Без текста в аргументах читается stdin.\n\nФлаги:\n")
		fs.PrintDefaults()
	}

	schemes := fs.String("schemes", strings.Join(masklinks.DefaultSchemes(), ","),
		"схемы через запятую, адреса после которых маскируются")

	if err := fs.Parse(args); err != nil {
		return config{}, err
	}

	list, err := splitSchemes(*schemes)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return config{}, err
	}

	return config{schemes: list, text: fs.Args()}, nil
}

func run(cfg config, stdin io.Reader, stdout io.Writer) error {
	out := bufio.NewWriter(stdout)

	if len(cfg.text) > 0 {
		masked := masklinks.MaskSchemes(strings.Join(cfg.text, " "), cfg.schemes...)
		if _, err := fmt.Fprintln(out, masked); err != nil {
			return err
		}
		return out.Flush()
	}

	if err := maskStream(bufio.NewReader(stdin), out, cfg.schemes); err != nil {
		return err
	}

	return out.Flush()
}

func maskStream(in *bufio.Reader, out io.Writer, schemes []string) error {
	for {
		line, readErr := in.ReadString('\n')

		if line != "" {
			if _, err := io.WriteString(out, masklinks.MaskSchemes(line, schemes...)); err != nil {
				return err
			}
		}

		if readErr != nil {
			if errors.Is(readErr, io.EOF) {
				return nil
			}
			return fmt.Errorf("чтение ввода: %w", readErr)
		}
	}
}

func splitSchemes(s string) ([]string, error) {
	var schemes []string
	for _, part := range strings.Split(s, ",") {
		if part = strings.TrimSpace(part); part != "" {
			schemes = append(schemes, part)
		}
	}

	if len(schemes) == 0 {
		return nil, errors.New("не задано ни одной схемы")
	}

	return schemes, nil
}
