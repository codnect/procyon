package runtime

import "io"

// Banner interface represents a banner that can be printed to an io.Writer.
type Banner interface {
	// PrintBanner method prints the banner to the provided io.Writer.
	PrintBanner(writer io.Writer) error
}
