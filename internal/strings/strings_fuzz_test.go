// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2026 The Ebitengine Authors

package strings

import (
	"strings"
	"testing"
	"unsafe"
)

// FuzzCStringGoString fuzzes the CString -> GoString round trip.
//
// Whatever CString does - zero-copy for a string that already ends in NUL,
// copy otherwise - the bytes must survive untouched, and GoString must return
// exactly the prefix of the input up to its first NUL, which is the contract
// of a C string conversion.
//
// The string input needs no decoding: Go fuzzing supports string arguments
// directly, so the fuzzer mutates the value itself and coverage guidance
// applies to the real code (the zero-copy/copy branch and the NUL scan loop).
func FuzzCStringGoString(f *testing.F) {
	// Seed corpus: deterministic edge cases. These run on every plain
	// `go test`, so the fast deterministic checks stay in normal CI while
	// longer exploration happens under `go test -fuzz`.
	for _, s := range []string{
		"",
		"\x00",
		"a",
		"\x00\x00",
		"\x00a\x00",
		"hello",
		"hello\x00",
		"a\x00b",
		"日本語",
		"日本語\x00後",
		strings.Repeat("x", 1024),
		strings.Repeat("x", 1024) + "\x00",
	} {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, s string) {
		want := s
		if i := strings.IndexByte(s, 0); i >= 0 {
			want = s[:i]
		}
		if got := GoString(uintptr(unsafe.Pointer(CString(s)))); got != want {
			t.Fatalf("round trip of %q: got %q, want %q", s, got, want)
		}
	})
}
