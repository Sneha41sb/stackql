package prettyprint_test

import (
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stackql/stackql/pkg/prettyprint"
)

func newDefaultContext() prettyprint.Context {
	return prettyprint.NewPrettyPrintContext(false, 2, 0, "|", logrus.New())
}

func newPrettyContext() prettyprint.Context {
	return prettyprint.NewPrettyPrintContext(true, 2, 0, "|", logrus.New())
}

func TestNewPrettyPrintContext(t *testing.T) {
	ctx := prettyprint.NewPrettyPrintContext(true, 4, 2, ">", logrus.New())
	if ctx.PrettyPrint != true {
		t.Errorf("PrettyPrint = %v, want true", ctx.PrettyPrint)
	}
	if ctx.Indentation != 4 {
		t.Errorf("Indentation = %d, want 4", ctx.Indentation)
	}
	if ctx.BaseIndentation != 2 {
		t.Errorf("BaseIndentation = %d, want 2", ctx.BaseIndentation)
	}
	if ctx.Delimiter != ">" {
		t.Errorf("Delimiter = %q, want %q", ctx.Delimiter, ">")
	}
}

func TestNewPrettyPrinter(t *testing.T) {
	ctx := prettyprint.NewPrettyPrintContext(false, 2, 4, "|", logrus.New())
	pp := prettyprint.NewPrettyPrinter(ctx)
	if pp == nil {
		t.Fatal("NewPrettyPrinter returned nil")
	}
}

func TestRenderColumnName(t *testing.T) {
	cases := []struct {
		name         string
		base         int
		cn           string
		want         string
	}{
		{name: "zero base", base: 0, cn: "col", want: "col"},
		{name: "non-zero base", base: 2, cn: "col", want: "  col"},
		{name: "empty column", base: 0, cn: "", want: ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			pp := prettyprint.NewPrettyPrinter(prettyprint.NewPrettyPrintContext(false, 2, tc.base, "|", logrus.New()))
			if got := pp.RenderColumnName(tc.cn); got != tc.want {
				t.Errorf("RenderColumnName(%q) = %q, want %q", tc.cn, got, tc.want)
			}
		})
	}
}

func TestRenderTemplateVarAndDelimit(t *testing.T) {
	cases := []struct {
		name  string
		base  int
		delim string
		tv    string
		want  string
	}{
		{name: "zero base", base: 0, delim: "|", tv: "x", want: "|{{ .values.x }}|"},
		{name: "non-zero base", base: 2, delim: "|", tv: "x", want: "  |{{ .values.x }}|"},
		{name: "multi-char delim", base: 0, delim: "::", tv: "foo", want: "::{{ .values.foo }}::"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			pp := prettyprint.NewPrettyPrinter(prettyprint.NewPrettyPrintContext(false, 2, tc.base, tc.delim, logrus.New()))
			if got := pp.RenderTemplateVarAndDelimit(tc.tv); got != tc.want {
				t.Errorf("RenderTemplateVarAndDelimit(%q) = %q, want %q", tc.tv, got, tc.want)
			}
		})
	}
}

func TestRenderTemplateVarNoDelimit(t *testing.T) {
	cases := []struct {
		name  string
		base  int
		delim string
		tv    string
		want  string
	}{
		{name: "zero base", base: 0, delim: "|", tv: "x", want: " {{ .values.x }}"},
		{name: "non-zero base", base: 2, delim: "|", tv: "x", want: "   {{ .values.x }}"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			pp := prettyprint.NewPrettyPrinter(prettyprint.NewPrettyPrintContext(false, 2, tc.base, tc.delim, logrus.New()))
			if got := pp.RenderTemplateVarNoDelimit(tc.tv); got != tc.want {
				t.Errorf("RenderTemplateVarNoDelimit(%q) = %q, want %q", tc.tv, got, tc.want)
			}
		})
	}
}

func TestRenderTemplateVarPlaceholderNoDelimit(t *testing.T) {
	cases := []struct {
		name  string
		base  int
		delim string
		tv    string
		want  string
	}{
		{name: "zero base", base: 0, delim: "|", tv: "x", want: " x: << x >>"},
		{name: "non-zero base", base: 2, delim: "|", tv: "x", want: "   x: << x >>"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			pp := prettyprint.NewPrettyPrinter(prettyprint.NewPrettyPrintContext(false, 2, tc.base, tc.delim, logrus.New()))
			if got := pp.RenderTemplateVarPlaceholderNoDelimit(tc.tv); got != tc.want {
				t.Errorf("RenderTemplateVarPlaceholderNoDelimit(%q) = %q, want %q", tc.tv, got, tc.want)
			}
		})
	}
}

func TestRenderTemplateVarPlaceholderKeyNoDelimit(t *testing.T) {
	cases := []struct {
		name  string
		base  int
		tv    string
		want  string
	}{
		{name: "zero base", base: 0, tv: "x", want: " x"},
		{name: "non-zero base", base: 2, tv: "x", want: "   x"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			pp := prettyprint.NewPrettyPrinter(prettyprint.NewPrettyPrintContext(false, 2, tc.base, "|", logrus.New()))
			if got := pp.RenderTemplateVarPlaceholderKeyNoDelimit(tc.tv); got != tc.want {
				t.Errorf("RenderTemplateVarPlaceholderKeyNoDelimit(%q) = %q, want %q", tc.tv, got, tc.want)
			}
		})
	}
}

func TestPrintTemplatedJSONString(t *testing.T) {
	// plain string → baseIndentationNoDelimit (delim len=1 → 1 leading space)
	pp := prettyprint.NewPrettyPrinter(newDefaultContext())
	got, err := pp.PrintTemplatedJSON("hello")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := " hello"
	if got != want {
		t.Errorf("PrintTemplatedJSON(\"hello\") = %q, want %q", got, want)
	}
}

func TestPrintTemplatedJSONStringQuoted(t *testing.T) {
	// quoted string → trim quotes → baseIndentationAndDelimit
	pp := prettyprint.NewPrettyPrinter(newDefaultContext())
	got, err := pp.PrintTemplatedJSON(`"value"`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "|value|"
	if got != want {
		t.Errorf("PrintTemplatedJSON(`\"value\"`) = %q, want %q", got, want)
	}
}

func TestPrintTemplatedJSONMap(t *testing.T) {
	pp := prettyprint.NewPrettyPrinter(newDefaultContext())
	body := map[string]interface{}{"a": "x", "b": "y"}
	got, err := pp.PrintTemplatedJSON(body)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "|{ \"a\": x, \"b\": y }|"
	if got != want {
		t.Errorf("PrintTemplatedJSON(map) = %q, want %q", got, want)
	}
}

func TestPrintTemplatedJSONMapPretty(t *testing.T) {
	pp := prettyprint.NewPrettyPrinter(newPrettyContext())
	body := map[string]interface{}{"a": "x"}
	got, err := pp.PrintTemplatedJSON(body)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "|{\n  \"a\": x\n }|"
	if got != want {
		t.Errorf("PrintTemplatedJSON(pretty map) = %q, want %q", got, want)
	}
}

func TestPrintTemplatedJSONArray(t *testing.T) {
	pp := prettyprint.NewPrettyPrinter(newDefaultContext())
	body := []interface{}{"x", "y"}
	got, err := pp.PrintTemplatedJSON(body)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "|[ x, y ]|"
	if got != want {
		t.Errorf("PrintTemplatedJSON(arr) = %q, want %q", got, want)
	}
}

func TestPrintTemplatedJSONArrayPretty(t *testing.T) {
	pp := prettyprint.NewPrettyPrinter(newPrettyContext())
	body := []interface{}{"x"}
	got, err := pp.PrintTemplatedJSON(body)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "|[\n  x\n ]|"
	if got != want {
		t.Errorf("PrintTemplatedJSON(pretty arr) = %q, want %q", got, want)
	}
}

func TestPrintTemplatedJSONUnsupportedType(t *testing.T) {
	pp := prettyprint.NewPrettyPrinter(newDefaultContext())
	_, err := pp.PrintTemplatedJSON(42)
	if err == nil {
		t.Error("PrintTemplatedJSON(42) expected error, got nil")
	}
	if !strings.Contains(err.Error(), "int") {
		t.Errorf("error %q does not mention type 'int'", err.Error())
	}
}

func TestPrintPlaceholderJSONString(t *testing.T) {
	pp := prettyprint.NewPrettyPrinter(newDefaultContext())
	got, err := pp.PrintPlaceholderJSON("hello")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := " hello"
	if got != want {
		t.Errorf("PrintPlaceholderJSON(\"hello\") = %q, want %q", got, want)
	}
}

func TestPrintPlaceholderJSONMapPretty(t *testing.T) {
	pp := prettyprint.NewPrettyPrinter(newPrettyContext())
	body := map[string]interface{}{"a": "x"}
	got, err := pp.PrintPlaceholderJSON(body)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := " {\n   a: x\n }"
	if got != want {
		t.Errorf("PrintPlaceholderJSON(pretty map) = %q, want %q", got, want)
	}
}

func TestPrintPlaceholderJSONArray(t *testing.T) {
	pp := prettyprint.NewPrettyPrinter(newDefaultContext())
	body := []interface{}{"x", "y"}
	got, err := pp.PrintPlaceholderJSON(body)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := " [ x, y ]"
	if got != want {
		t.Errorf("PrintPlaceholderJSON(arr) = %q, want %q", got, want)
	}
}

func TestPrintPlaceholderJSONArrayPretty(t *testing.T) {
	pp := prettyprint.NewPrettyPrinter(newPrettyContext())
	body := []interface{}{"x"}
	got, err := pp.PrintPlaceholderJSON(body)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := " [\n  x\n ]"
	if got != want {
		t.Errorf("PrintPlaceholderJSON(pretty arr) = %q, want %q", got, want)
	}
}

func TestPrintPlaceholderJSONUnsupportedType(t *testing.T) {
	pp := prettyprint.NewPrettyPrinter(newDefaultContext())
	_, err := pp.PrintPlaceholderJSON(42)
	if err == nil {
		t.Error("PrintPlaceholderJSON(42) expected error, got nil")
	}
	if !strings.Contains(err.Error(), "int") {
		t.Errorf("error %q does not mention type 'int'", err.Error())
	}
}

func TestPrintPlaceholderJSONNestedMap(t *testing.T) {
	pp := prettyprint.NewPrettyPrinter(newPrettyContext())
	body := map[string]interface{}{
		"a": map[string]interface{}{"b": "x"},
	}
	got, err := pp.PrintPlaceholderJSON(body)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := " {\n   a: {\n       b: x\n   }\n }"
	if got != want {
		t.Errorf("PrintPlaceholderJSON(nested map) = %q, want %q", got, want)
	}
}
