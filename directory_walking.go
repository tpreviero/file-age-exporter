package main

import (
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"sync"
)

var FileCounters = struct {
	mu   sync.RWMutex
	data map[int]map[string]map[int]float64
}{
	data: make(map[int]map[string]map[int]float64),
}

func doWalk(configuration *Configuration) {
	FileCounters.mu.Lock()
	defer FileCounters.mu.Unlock()

	yearToMonthToWeekToCount := make(map[int]map[string]map[int]float64)
	countFileByYearMonthWeek := newWalkFunction(yearToMonthToWeekToCount)

	for _, directory := range configuration.Directories {
		err := filepath.Walk(directory, countFileByYearMonthWeek)
		if err != nil {
			log.Errorf("Error trying to walk directory %s: %v", directory, err)
		}
	}

	FileCounters.data = yearToMonthToWeekToCount
}

func newWalkFunction(yearToMonthToWeekToCount map[int]map[string]map[int]float64) filepath.WalkFunc {
	return memoizedWalkFunction(walkFunction(yearToMonthToWeekToCount))
}

func memoizedWalkFunction(walkFunc filepath.WalkFunc) filepath.WalkFunc {
	visited := make(map[string]bool)
	return func(path string, info os.FileInfo, err error) error {
		if visited[path] {
			return filepath.SkipDir
		}
		visited[path] = true
		return walkFunc(path, info, err)
	}
}

func walkFunction(yearToMonthToWeekToCount map[int]map[string]map[int]float64) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Errorf("Error trying to walk directory %s: %v", path, err)
			return filepath.SkipDir
		}
		for _, exclusion := range configuration.Exclusions {
			if exclusion.Match(path) {
				log.Infof("Excluding %s", path)
				return filepath.SkipDir
			}
		}

		if !info.IsDir() {
			year, week := info.ModTime().ISOWeek()
			month := info.ModTime().Month().String()

			if yearToMonthToWeekToCount[year] == nil {
				yearToMonthToWeekToCount[year] = make(map[string]map[int]float64)
			}
			if yearToMonthToWeekToCount[year][month] == nil {
				yearToMonthToWeekToCount[year][month] = make(map[int]float64)
			}

			yearToMonthToWeekToCount[year][month][week]++
		}
		log.Debugf("dir: %v: name: %s\n", info.IsDir(), path)
		return nil
	}
}
