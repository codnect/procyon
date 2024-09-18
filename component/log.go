package component

import "codnect.io/logy"

// log is a global variable that holds the logger instance.
// It uses the Get function from the logy to retrieve the logger.
var (
	log = logy.Get()
)
