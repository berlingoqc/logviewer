package ty

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResolveString(t *testing.T) {

	t.Setenv("LOVE", "love")
	t.Setenv("BLIND", "visible")

	ms := MS{
		"test": "${LOVE}-is-${BLIND}",
	}

	resolvedMs := ms.ResolveVariables()

	assert.Equal(t, "love-is-visible", resolvedMs["test"], "failed to correctlty resolved varialbes")

}

func TestResolveNoEnv(t *testing.T) {

	ms := MS{
		"test": "${LOVE}-is-${BLIND}",
	}

	resolvedMs := ms.ResolveVariables()

	assert.Equal(t, "${LOVE}-is-${BLIND}", resolvedMs["test"], "failed to correctlty resolved varialbes")

}

func TestResolveStringDefault(t *testing.T) {

	t.Setenv("LOVE", "love")

	ms := MS{
		"test": "${LOVE}-is-${BLIND:-blind}",
	}

	resolvedMs := ms.ResolveVariables()

	assert.Equal(t, "love-is-blind", resolvedMs["test"], "failed to correctlty resolved varialbes")

}
