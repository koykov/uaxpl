package uaxpl

import (
	"strings"
	"testing"
)

func TestTokenizer_Next(t *testing.T) {
	tests := []struct {
		name     string
		ua       string
		expected []string
	}{
		{
			name:     "empty string",
			ua:       "",
			expected: []string{},
		},
		{
			name:     "simple user agent",
			ua:       "Mozilla/5.0",
			expected: []string{"Mozilla", "5", "0"},
		},
		{
			name:     "full browser user agent",
			ua:       "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
			expected: []string{"Mozilla", "5", "0", "Windows", "NT", "10", "0", "Win64", "x64", "AppleWebKit", "537", "36"},
		},
		{
			name:     "bot with version",
			ua:       "Googlebot/2.1 (+http://www.google.com/bot.html)",
			expected: []string{"Googlebot", "2", "1", "http", "www", "google", "com", "bot", "html"},
		},
		{
			name:     "curl user agent",
			ua:       "curl/7.68.0",
			expected: []string{"curl", "7", "68", "0"},
		},
		{
			name:     "python requests",
			ua:       "python-requests/2.25.1",
			expected: []string{"python-requests", "2", "25", "1"},
		},
		{
			name:     "bot with underscores",
			ua:       "facebookexternalhit/1.1",
			expected: []string{"facebookexternalhit", "1", "1"},
		},
		{
			name:     "bot with multiple delimiters",
			ua:       "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)",
			expected: []string{"Mozilla", "5", "0", "compatible", "Googlebot", "2", "1", "http", "www", "google", "com", "bot", "html"},
		},
		{
			name:     "only delimiters",
			ua:       " / ( ) ; , . [ ] { } : = + * ",
			expected: []string{},
		},
		{
			name:     "tabs and newlines",
			ua:       "Mozilla\t5.0\nChrome\r91.0",
			expected: []string{"Mozilla", "5", "0", "Chrome", "91", "0"},
		},
		{
			name:     "braces and brackets",
			ua:       "{Mozilla}/5.0 [Chrome] (91.0)",
			expected: []string{"Mozilla", "5", "0", "Chrome", "91", "0"},
		},
		{
			name:     "colon and equals",
			ua:       "name:value=test",
			expected: []string{"name", "value", "test"},
		},
		{
			name:     "plus and star",
			ua:       "test+123*456",
			expected: []string{"test", "123", "456"},
		},
		{
			name:     "dot as delimiter",
			ua:       "Mozilla.5.0.Chrome.91",
			expected: []string{"Mozilla", "5", "0", "Chrome", "91"},
		},
		{
			name:     "real googlebot",
			ua:       "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)",
			expected: []string{"Mozilla", "5", "0", "compatible", "Googlebot", "2", "1", "http", "www", "google", "com", "bot", "html"},
		},
		{
			name:     "real bingbot",
			ua:       "Mozilla/5.0 (compatible; bingbot/2.0; +http://www.bing.com/bingbot.htm)",
			expected: []string{"Mozilla", "5", "0", "compatible", "bingbot", "2", "0", "http", "www", "bing", "com", "bingbot", "htm"},
		},
		{
			name:     "real yandexbot",
			ua:       "Mozilla/5.0 (compatible; YandexBot/3.0; +http://yandex.com/bots)",
			expected: []string{"Mozilla", "5", "0", "compatible", "YandexBot", "3", "0", "http", "yandex", "com", "bots"},
		},
		{
			name:     "real curl",
			ua:       "curl/7.68.0 (x86_64-pc-linux-gnu) libcurl/7.68.0 OpenSSL/1.1.1f zlib/1.2.11",
			expected: []string{"curl", "7", "68", "0", "x86_64-pc-linux-gnu", "libcurl", "7", "68", "0", "OpenSSL", "1", "1", "1f", "zlib", "1", "2", "11"},
		},
		{
			name:     "real python requests",
			ua:       "python-requests/2.25.1",
			expected: []string{"python-requests", "2", "25", "1"},
		},
		{
			name:     "complex with multiple slashes",
			ua:       "Mozilla/5.0 (Windows NT/10.0; Win64/x64) AppleWebKit/537.36 (KHTML/ like Gecko) Chrome/91.0.4472.124 Safari/537.36",
			expected: []string{"Mozilla", "5", "0", "Windows", "NT", "10", "0", "Win64", "x64", "AppleWebKit", "537", "36", "KHTML", "like", "Gecko", "Chrome", "91", "0", "4472", "124", "Safari", "537", "36"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var tkn Tokenizer
			var result []string

			for {
				token, eof := tkn.NextString(tt.ua)
				if eof {
					break
				}
				result = append(result, token)
			}

			if len(result) != len(tt.expected) {
				t.Errorf("got %d tokens, expected %d", len(result), len(tt.expected))
				t.Errorf("got: %v", result)
				t.Errorf("expected: %v", tt.expected)
				return
			}

			for i := range result {
				if result[i] != tt.expected[i] {
					t.Errorf("token %d: got %q, expected %q", i, result[i], tt.expected[i])
				}
			}
		})
	}
}

func TestTokenizer_Reset(t *testing.T) {
	var tkn Tokenizer
	ua := "Mozilla/5.0 Chrome/91.0"

	// First read
	token, eof := tkn.NextString(ua)
	if eof || token != "Mozilla" {
		t.Errorf("first token: got %q, expected Mozilla", token)
	}

	// Reset and read again
	tkn.Reset()
	token, eof = tkn.NextString(ua)
	if eof || token != "Mozilla" {
		t.Errorf("after reset: got %q, expected Mozilla", token)
	}
}

func TestTokenizer_ConsecutiveCalls(t *testing.T) {
	var tkn Tokenizer
	ua := "Mozilla/5.0 Chrome/91.0"

	expected := []string{"Mozilla", "5", "0", "Chrome", "91", "0"}
	var result []string

	for i := 0; i < 6; i++ {
		token, eof := tkn.NextString(ua)
		if eof {
			t.Errorf("unexpected EOF at position %d", i)
			break
		}
		result = append(result, token)
	}

	// Check EOF
	_, eof := tkn.NextString(ua)
	if !eof {
		t.Error("expected EOF, got more tokens")
	}

	if len(result) != len(expected) {
		t.Errorf("got %d tokens, expected %d", len(result), len(expected))
	}

	for i := range result {
		if result[i] != expected[i] {
			t.Errorf("token %d: got %q, expected %q", i, result[i], expected[i])
		}
	}
}

func TestTokenizer_EmptyTokens(t *testing.T) {
	var tkn Tokenizer
	ua := "   \t\n\r   "

	_, eof := tkn.NextString(ua)
	if !eof {
		t.Error("expected EOF for string with only delimiters")
	}
}

func TestTokenizer_NoDelimiters(t *testing.T) {
	var tkn Tokenizer
	ua := "SingleToken"

	token, eof := tkn.NextString(ua)
	if eof {
		t.Error("expected token, got EOF")
	}
	if token != "SingleToken" {
		t.Errorf("got %q, expected SingleToken", token)
	}

	_, eof = tkn.NextString(ua)
	if !eof {
		t.Error("expected EOF after single token")
	}
}

func TestTokenizer_IsDelimiter(t *testing.T) {
	delimiters := []byte{' ', '\t', '\n', '\r', '/', '(', ')', ';', ',', '.', '[', ']', '{', '}', ':', '=', '+', '*'}
	nonDelimiters := []byte{'a', 'b', 'c', 'A', 'B', 'C', '0', '1', '9', '@', '#', '$', '%', '^', '&', '~', '`', '|', '\\', '\''}

	for _, d := range delimiters {
		if !isDelimiter(d) {
			t.Errorf("isDelimiter(%q) returned false, expected true", d)
		}
	}

	for _, nd := range nonDelimiters {
		if isDelimiter(nd) {
			t.Errorf("isDelimiter(%q) returned true, expected false", nd)
		}
	}
}

func BenchmarkTokenizer_Next(b *testing.B) {
	benchmarks := []struct {
		name string
		ua   string
	}{
		{
			name: "simple",
			ua:   "Mozilla/5.0",
		},
		{
			name: "full_browser",
			ua:   "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
		},
		{
			name: "googlebot",
			ua:   "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)",
		},
		{
			name: "bingbot",
			ua:   "Mozilla/5.0 (compatible; bingbot/2.0; +http://www.bing.com/bingbot.htm)",
		},
		{
			name: "curl",
			ua:   "curl/7.68.0 (x86_64-pc-linux-gnu) libcurl/7.68.0 OpenSSL/1.1.1f zlib/1.2.11",
		},
		{
			name: "python_requests",
			ua:   "python-requests/2.25.1",
		},
		{
			name: "long_complex",
			ua:   "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36 Edg/91.0.864.59",
		},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			var tkn Tokenizer
			b.ResetTimer()
			b.ReportAllocs()

			for i := 0; i < b.N; i++ {
				tkn.Reset()
				for {
					_, eof := tkn.NextString(bm.ua)
					if eof {
						break
					}
				}
			}
		})
	}
}

func BenchmarkTokenizer_Next_CompareWithSplit(b *testing.B) {
	ua := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"

	b.Run("tokenizer", func(b *testing.B) {
		var tkn Tokenizer
		b.ResetTimer()
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			tkn.Reset()
			for {
				_, eof := tkn.NextString(ua)
				if eof {
					break
				}
			}
		}
	})

	b.Run("strings.Split", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = strings.FieldsFunc(ua, func(r rune) bool {
				return r == ' ' || r == '/' || r == '_' || r == '-' ||
					r == '(' || r == ')' || r == ';' || r == ',' ||
					r == '.' || r == '[' || r == ']' || r == '{' ||
					r == '}' || r == ':' || r == '=' || r == '+' ||
					r == '*' || r == '\t' || r == '\n' || r == '\r'
			})
		}
	})

	b.Run("strings.SplitN", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = strings.Split(ua, " ")
		}
	})
}

func BenchmarkTokenizer_Next_RealWorld(b *testing.B) {
	userAgents := []string{
		"Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)",
		"Mozilla/5.0 (compatible; bingbot/2.0; +http://www.bing.com/bingbot.htm)",
		"Mozilla/5.0 (compatible; YandexBot/3.0; +http://yandex.com/bots)",
		"curl/7.68.0 (x86_64-pc-linux-gnu) libcurl/7.68.0 OpenSSL/1.1.1f zlib/1.2.11",
		"python-requests/2.25.1",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.114 Safari/537.36",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.114 Safari/537.36",
		"facebookexternalhit/1.1 (+http://www.facebook.com/externalhit_uatext.php)",
		"Twitterbot/1.0",
		"Mozilla/5.0 (compatible; DotBot/1.1; http://www.opensiteexplorer.org/dotbot, help@moz.com)",
		"Mozilla/5.0 (compatible; AhrefsBot/7.0; +http://ahrefs.com/robot/)",
		"Mozilla/5.0 (compatible; MJ12bot/v1.4.8; http://mj12bot.com/)",
		"Mozilla/5.0 (compatible; SemrushBot/7~bl; +http://www.semrush.com/bot.html)",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 14_6 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.0 Mobile/15E148 Safari/604.1",
	}

	b.Run("tokenizer", func(b *testing.B) {
		var tkn Tokenizer
		b.ResetTimer()
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			ua := userAgents[i%len(userAgents)]
			tkn.Reset()
			for {
				_, eof := tkn.NextString(ua)
				if eof {
					break
				}
			}
		}
	})

	b.Run("strings.FieldsFunc", func(b *testing.B) {
		delimFunc := func(r rune) bool {
			return r == ' ' || r == '/' || r == '_' || r == '-' ||
				r == '(' || r == ')' || r == ';' || r == ',' ||
				r == '.' || r == '[' || r == ']' || r == '{' ||
				r == '}' || r == ':' || r == '=' || r == '+' ||
				r == '*' || r == '\t' || r == '\n' || r == '\r'
		}

		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			ua := userAgents[i%len(userAgents)]
			_ = strings.FieldsFunc(ua, delimFunc)
		}
	})
}

func BenchmarkTokenizer_Next_LongUA(b *testing.B) {
	ua := strings.Repeat("Mozilla/5.0 (Windows NT 10.0; Win64; x64) ", 10)

	b.Run("tokenizer", func(b *testing.B) {
		var tkn Tokenizer
		b.ResetTimer()
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			tkn.Reset()
			for {
				_, eof := tkn.NextString(ua)
				if eof {
					break
				}
			}
		}
	})

	b.Run("strings.FieldsFunc", func(b *testing.B) {
		delimFunc := func(r rune) bool {
			return r == ' ' || r == '/' || r == '_' || r == '-' ||
				r == '(' || r == ')' || r == ';' || r == ',' ||
				r == '.' || r == '[' || r == ']' || r == '{' ||
				r == '}' || r == ':' || r == '=' || r == '+' ||
				r == '*' || r == '\t' || r == '\n' || r == '\r'
		}

		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = strings.FieldsFunc(ua, delimFunc)
		}
	})
}
