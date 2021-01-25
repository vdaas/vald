package decoder

import (
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook/conversion"
)

// Decoder represents a type alias of conversion.Decoder.
type Decoder = conversion.Decoder

// NewDecoder creates a Decoder given the runtime.Scheme.
// It will return an error when NewDecoder method failed.
func NewDecoder() (*Decoder, error) {
	return conversion.NewDecoder(runtime.NewScheme())
}
