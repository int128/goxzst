// Package log provides the logging feature.
package log

import "log"

var Printf = log.Printf

func init() {
	log.SetFlags(log.Lmicroseconds)
}
