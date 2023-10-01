package captureutil

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCapture(t *testing.T) {
	fmt.Println(GetCaptureCode())
}

func TestIsNeedCapture(t *testing.T) {
	assert.Equal(t, false, IsNeedCaptcha("1"))
}
