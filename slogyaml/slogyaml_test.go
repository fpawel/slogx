package slogyaml

import (
	"bytes"
	"log/slog"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestYAMLHandler_basic_test(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	h := NewYAMLHandler().WithWriter(buf).WithTimeLayout("")
	log := slog.New(h)
	log.WithGroup("group-name").
		Info("Hello, world!",
			"user", "alice",
			"id", 12,
			"slice", []string{"value5", "value6"},
			"something", map[string]any{
				"key1": "value1",
				"key2": 2,
				"key3": []string{"value3", "value4"},
			},
		)

	require.Equal(t, `- INFO  Hello, world!:
    - group-name:
        - user: alice
        - id: 12
        - slice:
            - value5
            - value6
        - something:
            key1: value1
            key2: 2
            key3:
                - value3
                - value4
`, buf.String())
}

func ExampleYAMLHandler_basic_test() {
	h := NewYAMLHandler().WithWriter(os.Stdout).WithTimeLayout("")
	log := slog.New(h)
	log.WithGroup("group-name").
		Info("Hello, world!",
			"user", "alice",
			"id", 12,
			"slice", []string{"value5", "value6"},
			"something", map[string]any{
				"key1": "value1",
				"key2": 2,
				"key3": []string{"value3", "value4"},
			},
		)
	// Output:
	// - INFO  Hello, world!:
	//     - group-name:
	//         - user: alice
	//         - id: 12
	//         - slice:
	//             - value5
	//             - value6
	//         - something:
	//             key1: value1
	//             key2: 2
	//             key3:
	//                 - value3
	//                 - value4
}
