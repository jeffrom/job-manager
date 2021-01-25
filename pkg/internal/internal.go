// Package internal contains code that we don't want to export outside of
// job-manager.
package internal

import "log"

func IgnoreError(err error) {
	if err != nil {
		log.Print("ignoring error:", err)
	}
}
