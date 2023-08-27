package tools

import (
	"os"

	"github.com/pkg/profile"
)

func StartProfiling(profileMode string) interface{ Stop() } {
	var profiler interface{ Stop() }
	switch profileMode {
	case "cpu":
		profiler = profile.Start(
			profile.CPUProfile,
			profile.ProfilePath(os.Getenv("PROFILER_PATH")))
	case "mem":
		profiler = profile.Start(
			profile.MemProfile,
			profile.ProfilePath(os.Getenv("PROFILER_PATH")))
	}

	return profiler
}
