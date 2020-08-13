package zstd

import "github.com/klauspost/compress/zstd"

// WithEncoderLevel calls zstd.WithEncoderLevel.
func WithEncoderLevel(level int) EOption {
	return zstd.WithEncoderLevel(
		zstd.EncoderLevelFromZstd(level),
	)
}
