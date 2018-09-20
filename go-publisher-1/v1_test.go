package kubeless

import (
	"testing"

	"github.com/kubeless/kubeless/pkg/functions"
)

func TestGenerateMessage(t *testing.T) {
	s := generateMessage(functions.Event{})
	if s == "" {
		t.Fatal("generateMessage returned the empty string, expected json.")
	} else {
		t.Log(s)
	}
}
