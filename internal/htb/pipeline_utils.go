package htb

import "github.com/kathleenfrench/sneak/internal/entity"

func getPipelineMapKeys(pipelines entity.Pipelines) []string {
	names := []string{}
	for n := range pipelines {
		names = append(names, n)
	}

	return names
}
