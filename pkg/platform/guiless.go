//go:build !windows && !macos
// +build !windows,!macos

package platform

import (
	`aurora/pkg/server`
)

func Start(s *server.Server) {
	s.Start()
}
