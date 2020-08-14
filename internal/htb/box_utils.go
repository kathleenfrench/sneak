package htb

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/kathleenfrench/sneak/internal/entity"
)

var osOptions = []string{
	"linux",
	"windows",
	"freeBSD",
	"openBSD",
	"other",
}

var difficulties = []string{
	"easy",
	"medium",
	"hard",
	"insane",
}

// completionColorizer returns an icon with the completion status of a box
func completionColorizer(completed bool) string {
	if completed {
		return color.HiGreenString("pwnd")
	}

	return color.HiYellowString("incomplete")
}

// difficultyColorizer colorizes based on difficulty
func difficultyColorizer(diff string) string {
	switch diff {
	case "easy":
		return color.GreenString("easy")
	case "medium":
		return color.YellowString("medium")
	case "hard":
		return color.HiRedString("hard")
	case "insane":
		return color.RedString("insane")
	}

	return ""
}

func constructBoxListing(box entity.Box) string {
	head := fmt.Sprintf(
		"%s - [%s][%s][%s]",
		box.Name,
		color.HiBlueString(box.OS),
		difficultyColorizer(box.Difficulty),
		completionColorizer(box.Completed),
	)

	return head
}

func makeGuiBoxMappings(boxes []entity.Box) (keys []string, mapping map[string]entity.Box) {
	mapping = make(map[string]entity.Box)

	for _, b := range boxes {
		name := constructBoxListing(b)
		keys = append(keys, name)
		mapping[name] = b
	}

	return keys, mapping
}
