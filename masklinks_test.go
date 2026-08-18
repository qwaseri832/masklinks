package masklinks_test

import (
	"strings"
	"testing"

	"github.com/qwaseri832/masklinks"
)

func TestMask(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "пустая строка",
			in:   "",
			want: "",
		},
		{
			name: "текст без ссылок",
			in:   "Hello, its my page",
			want: "Hello, its my page",
		},
		{
			name: "ссылка в середине",
			in:   "Hello, its my page: http://localhost123.com See you",
			want: "Hello, its my page: http://**************** See you",
		},
		{
			name: "ссылка в конце строки",
			in:   "заходи на http://example.com",
			want: "заходи на http://***********",
		},
		{
			name: "https тоже маскируется",
			in:   "https://example.com готово",
			want: "https://*********** готово",
		},
		{
			name: "перевод строки завершает адрес",
			in:   "http://example.com\nПривет",
			want: "http://***********\nПривет",
		},
		{
			name: "табуляция завершает адрес",
			in:   "http://example.com\tдальше",
			want: "http://***********\tдальше",
		},
		{
			name: "две ссылки подряд",
			in:   "http://a.ru и https://b.ru",
			want: "http://**** и https://****",
		},
		{
			name: "схема без адреса",
			in:   "http:// пусто",
			want: "http:// пусто",
		},
		{
			name: "регистр схемы игнорируется",
			in:   "HTTP://Example.COM",
			want: "HTTP://***********",
		},
		{
			name: "кириллица в адресе маскируется посимвольно, а не побайтово",
			in:   "http://сайт.рф всё",
			want: "http://******* всё",
		},
		{
			name: "схема внутри слова тоже считается ссылкой",
			in:   "см.http://a.ru",
			want: "см.http://****",
		},
		{
			name: "https не разбивается на http + s",
			in:   "https://a.ru",
			want: "https://****",
		},
		{
			name: "İ перед схемой не сдвигает границу адреса",
			in:   "İ http://example.com",
			want: "İ http://***********",
		},
		{
			name: "знак Кельвина перед схемой не сдвигает границу адреса",
			in:   "Khttp://a.ru",
			want: "Khttp://****",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := masklinks.Mask(tt.in); got != tt.want {
				t.Errorf("Mask(%q)\n получено: %q\n ожидалось: %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestMaskSchemes(t *testing.T) {
	t.Run("пустой список схем оставляет текст как есть", func(t *testing.T) {
		const in = "http://example.com"
		if got := masklinks.MaskSchemes(in); got != in {
			t.Errorf("получено %q, ожидалось %q", got, in)
		}
	})

	t.Run("своя схема", func(t *testing.T) {
		got := masklinks.MaskSchemes("ftp://files.local готово", "ftp://")
		const want = "ftp://*********** готово"
		if got != want {
			t.Errorf("получено %q, ожидалось %q", got, want)
		}
	})

	t.Run("схема, которой нет в списке, не трогается", func(t *testing.T) {
		const in = "http://example.com"
		if got := masklinks.MaskSchemes(in, "ftp://"); got != in {
			t.Errorf("получено %q, ожидалось %q", got, in)
		}
	})
}

func TestMaskPreservesRuneCount(t *testing.T) {
	inputs := []string{
		"http://пример.рф и ещё http://a.b",
		"без ссылок вообще",
		"https://x\thttp://y\nhttps://z",
		"İK http://пример.рф",
	}
	for _, in := range inputs {
		got := masklinks.Mask(in)
		if a, b := len([]rune(got)), len([]rune(in)); a != b {
			t.Errorf("Mask(%q): рун стало %d, было %d", in, a, b)
		}
	}
}

func BenchmarkMask(b *testing.B) {
	text := strings.Repeat("текст http://example.com/path?q=1 ещё текст\n", 100)
	b.ReportAllocs()
	for b.Loop() {
		_ = masklinks.Mask(text)
	}
}
