package htb

import "github.com/kathleenfrench/sneak/internal/entity"

func getJobKeys(jobs map[string]*entity.Job) []string {
	names := []string{}
	for n := range jobs {
		names = append(names, n)
	}

	return names
}
