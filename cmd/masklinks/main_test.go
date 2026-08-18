package main

import (
	"bytes"
	"errors"
	"flag"
	"strings"
	"testing"
)

func mustParse(t *testing.T, args ...string) config {
	t.Helper()

	cfg, err := parseFlags("masklinks", args, new(bytes.Buffer))
	if err != nil {
		t.Fatalf("parseFlags(%q): %v", args, err)
	}

	return cfg
}

func TestRunMasksArguments(t *testing.T) {
	var stdout bytes.Buffer

	if err := run(mustParse(t, "мой сайт: http://example.com"), nil, &stdout); err != nil {
		t.Fatalf("run: %v", err)
	}

	const want = "мой сайт: http://***********\n"
	if got := stdout.String(); got != want {
		t.Errorf("получено %q, ожидалось %q", got, want)
	}
}

func TestRunJoinsArguments(t *testing.T) {
	var stdout bytes.Buffer

	if err := run(mustParse(t, "ссылка", "http://a.ru", "дальше"), nil, &stdout); err != nil {
		t.Fatalf("run: %v", err)
	}

	const want = "ссылка http://**** дальше\n"
	if got := stdout.String(); got != want {
		t.Errorf("получено %q, ожидалось %q", got, want)
	}
}

func TestRunReadsStdin(t *testing.T) {
	var stdout bytes.Buffer
	stdin := strings.NewReader("первая http://a.ru\nвторая https://b.ru")

	if err := run(mustParse(t), stdin, &stdout); err != nil {
		t.Fatalf("run: %v", err)
	}

	const want = "первая http://****\nвторая https://****"
	if got := stdout.String(); got != want {
		t.Errorf("получено %q, ожидалось %q", got, want)
	}
}

func TestRunCustomSchemes(t *testing.T) {
	var stdout bytes.Buffer

	cfg := mustParse(t, "-schemes", "ftp:// , sftp://", "архив на ftp://files.local и http://a.ru")
	if err := run(cfg, nil, &stdout); err != nil {
		t.Fatalf("run: %v", err)
	}

	const want = "архив на ftp://*********** и http://a.ru\n"
	if got := stdout.String(); got != want {
		t.Errorf("получено %q, ожидалось %q", got, want)
	}
}

func TestRunHandlesLineLongerThanScannerLimit(t *testing.T) {
	var stdout bytes.Buffer

	address := strings.Repeat("a", 128*1024)
	stdin := strings.NewReader("http://" + address + "\n")

	if err := run(mustParse(t), stdin, &stdout); err != nil {
		t.Fatalf("run: %v", err)
	}

	want := "http://" + strings.Repeat("*", len(address)) + "\n"
	if got := stdout.String(); got != want {
		t.Errorf("длинная строка обработана неверно: получено %d байт, ожидалось %d", len(got), len(want))
	}
}

func TestParseFlagsRejectsEmptySchemes(t *testing.T) {
	var stderr bytes.Buffer

	if _, err := parseFlags("masklinks", []string{"-schemes", " , ", "http://a.ru"}, &stderr); err == nil {
		t.Fatal("ожидалась ошибка, получено nil")
	}

	if stderr.Len() == 0 {
		t.Error("причина отказа должна попадать в stderr")
	}
}

func TestParseFlagsRejectsUnknownFlag(t *testing.T) {
	var stderr bytes.Buffer

	_, err := parseFlags("masklinks", []string{"-нет-такого-флага"}, &stderr)
	if err == nil {
		t.Fatal("ожидалась ошибка, получено nil")
	}

	if errors.Is(err, flag.ErrHelp) {
		t.Error("неизвестный флаг не должен выглядеть как запрос справки")
	}
}

func TestParseFlagsHelp(t *testing.T) {
	var stderr bytes.Buffer

	_, err := parseFlags("masklinks", []string{"-h"}, &stderr)
	if !errors.Is(err, flag.ErrHelp) {
		t.Fatalf("получено %v, ожидалось flag.ErrHelp", err)
	}

	if !strings.Contains(stderr.String(), "-schemes") {
		t.Error("справка должна перечислять флаги")
	}
}

func TestSplitSchemes(t *testing.T) {
	got, err := splitSchemes(" http:// ,, https:// ")
	if err != nil {
		t.Fatalf("splitSchemes: %v", err)
	}

	want := []string{"http://", "https://"}
	if len(got) != len(want) {
		t.Fatalf("получено %q, ожидалось %q", got, want)
	}

	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("получено %q, ожидалось %q", got, want)
		}
	}
}
