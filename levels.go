package main

import (
	"regexp"
	"strconv"
)

type VolumeLevels struct {
	left  float64
	right float64
}

func makeVolumeLevels(left string, right string) VolumeLevels {
	left_, err := strconv.ParseFloat(left, 32)
	chk(err)

	right_, err := strconv.ParseFloat(right, 32)
	chk(err)

	levels := VolumeLevels{
		left:  left_,
		right: right_,
	}

	return levels
}

func (levels VolumeLevels) mean() float64 {
	return float64((levels.left + levels.right) / 2.0)
}

func makeVolumeLevelsParsingLine(line string, r *regexp.Regexp) *VolumeLevels { // VolumeLevels { (levels VolumeLevels)
	submatch := r.FindStringSubmatch(line)

	if len(submatch) > 2 {
		value := makeVolumeLevels(submatch[1], submatch[2])
		return &value
	}
	return nil
}
