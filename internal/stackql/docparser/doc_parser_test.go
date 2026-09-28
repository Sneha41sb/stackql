package docparser_test

import (
	"testing"

	"github.com/stackql/stackql/internal/stackql/docparser"
)

func TestTranslateServiceKeyGenericProviderToIql(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{name: "no colon unchanged", in: "google", want: "google"},
		{name: "single colon replaced", in: "google:oauth", want: "google__oauth"},
		{name: "multiple colons replaced", in: "google:foo:bar", want: "google__foo__bar"},
		{name: "empty string", in: "", want: ""},
		{name: "leading colon", in: ":foo", want: "__foo"},
		{name: "trailing colon", in: "foo:", want: "foo__"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := docparser.TranslateServiceKeyGenericProviderToIql(tc.in)
			if got != tc.want {
				t.Errorf("TranslateServiceKeyGenericProviderToIql(%q) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}

func TestTranslateServiceKeyIqlToGenericProvider(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{name: "no double underscore unchanged", in: "google", want: "google"},
		{name: "single double-underscore replaced", in: "google__oauth", want: "google:oauth"},
		{name: "multiple double-underscores replaced", in: "google__foo__bar", want: "google:foo:bar"},
		{name: "empty string", in: "", want: ""},
		{name: "leading double underscore", in: "__foo", want: ":foo"},
		{name: "trailing double underscore", in: "foo__", want: "foo:"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := docparser.TranslateServiceKeyIqlToGenericProvider(tc.in)
			if got != tc.want {
				t.Errorf("TranslateServiceKeyIqlToGenericProvider(%q) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}

func TestTranslateServiceKeyRoundTrip(t *testing.T) {
	original := "google:foo:bar"
	generic := docparser.TranslateServiceKeyGenericProviderToIql(original)
	restored := docparser.TranslateServiceKeyIqlToGenericProvider(generic)
	if restored != original {
		t.Errorf("round-trip failed: %q -> %q -> %q", original, generic, restored)
	}
}

func TestTranslateServiceKeyDelimiters(t *testing.T) {
	// Verify the delimiter constants behave correctly through the public API:
	// a generic-provider key containing a single ":" maps to a stackql key with "__",
	// and vice versa. This confirms the constants used internally.
	generic := docparser.TranslateServiceKeyGenericProviderToIql("a:b")
	if generic != "a__b" {
		t.Errorf("generic->iql = %q, want %q", generic, "a__b")
	}
	iql := docparser.TranslateServiceKeyIqlToGenericProvider("a__b")
	if iql != "a:b" {
		t.Errorf("iql->generic = %q, want %q", iql, "a:b")
	}
}