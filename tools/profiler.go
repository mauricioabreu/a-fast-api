package tools

import (
	"github.com/mauricioabreu/a-fast-api/config"
	"github.com/pkg/profile"
)

func StartProfiling(cfg *config.Config) interface{ Stop() } {
	var profiler interface{ Stop() }
	switch cfg.ProfilerMode {
	case "cpu":
		profiler = profile.Start(
			profile.CPUProfile,
			profile.ProfilePath(cfg.ProfilerPath))
	case "mem":
		profiler = profile.Start(
			profile.MemProfile,
			profile.ProfilePath(cfg.ProfilerPath))
	}

	return profiler
}
