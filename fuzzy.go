// Package fuzzy provides methods for fuzzy string matching relying on an
// ensemble of different distance metrics.
package fuzzy

import (
	"fmt"
	"io"
	"math"
	"os"
	"os/exec"
	"strings"

	"github.com/antzucaro/matchr"
)

// Match returns the best match and a score between 0 and 1 for s in the corpus cs.
func Match(s string, cs []string) (string, float64) {
	w := make(map[string]float64)

	for _, c := range cs {
		dl := 1 - norm(float64(matchr.DamerauLevenshtein(s, c)), s, c)
		jw := matchr.JaroWinkler(s, c, false)
		sw := norm(matchr.SmithWaterman(s, c), s, c)

		w[c] = (dl + jw + sw) / 3
	}

	m := ""
	max := 0.0
	for _, c := range cs {
		v, _ := w[c]
		if v >= max {
			max = v
			m = c
		}
	}

	return m, max
}

func norm(s float64, a string, b string) float64 {
	return s / math.Max(float64(len(a)), float64(len(b)))
}

// Returns the match from the corpus cs selected by fzf.
func FzfMatch(cs []string) (string, error) {
	if !available("fzf") {
		return "", fmt.Errorf("fzf is not available")
	}
	return withFilter("fzf", func(in io.WriteCloser) {
		for _, c := range cs {
			fmt.Fprintln(in, c)
		}
	})
}

func available(name string) bool {
	cmd := exec.Command(shell(), "-c", "command -v "+name)
	err := cmd.Run()
	return err == nil
}

func withFilter(command string, input func(in io.WriteCloser)) (string, error) {
	cmd := exec.Command(shell(), "-c", command)
	cmd.Stderr = os.Stderr
	in, err := cmd.StdinPipe()
	if err != nil {
		return "", err
	}
	go func() {
		input(in)
		in.Close()
	}()
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func shell() string {
	sh := os.Getenv("SHELL")
	if len(sh) == 0 {
		sh = "sh"
	}
	return sh
}
