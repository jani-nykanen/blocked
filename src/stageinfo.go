package main

import (
	"strconv"

	"github.com/jani-nykanen/blocked/src/core"
)

type stageInfoEntry struct {
	name       string
	difficulty int32
}

type stageInfoContainer struct {
	entries []stageInfoEntry
}

func (sinfo *stageInfoContainer) getStageInfo(index int32) stageInfoEntry {

	if index < 0 || index >= int32(len(sinfo.entries)) {

		return stageInfoEntry{name: "null", difficulty: 1}
	}

	return sinfo.entries[index]
}

func parseStageInfo(folder string) *stageInfoContainer {

	sinfo := new(stageInfoContainer)

	sinfo.entries = make([]stageInfoEntry, 0)

	// This is a very rough way to do this, but... it works
	// TODO: ...but I still should rework it!
	var tmap *core.Tilemap
	var err error
	for i := 0; ; i++ {

		tmap, err = core.ParseTMX(folder + "/" + strconv.Itoa(i+1) + ".tmx")
		if err != nil {

			break
		}

		sinfo.entries = append(sinfo.entries,
			stageInfoEntry{
				name:       tmap.GetProperty("name", "null"),
				difficulty: tmap.GetNumericProperty("difficulty", 1)})
	}

	return sinfo
}
