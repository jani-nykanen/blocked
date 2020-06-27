package main

import (
	"os"

	"github.com/jani-nykanen/ultimate-puzzle/src/core"
)

const (
	defaultSaveFilePath = "save.dat"
)

type completionInfo struct {
	states       []int32
	currentStage int32
	// This should go elsewhere, but since this data
	// should be loaded only once, let's put it here...
	sinfo *stageInfoContainer
}

func (cinfo *completionInfo) updateState(index int32, state int32) {

	if index < 1 || index > cinfo.levelCount() {
		return
	}

	cinfo.states[index-1] = core.MaxInt32(state, cinfo.states[index-1])
}

func (cinfo *completionInfo) getState(index int32) int32 {

	if index < 1 || index > cinfo.levelCount() {
		return 0
	}

	return cinfo.states[index-1]
}

func (cinfo *completionInfo) levelCount() int32 {

	return int32(len(cinfo.states))
}

func (cinfo *completionInfo) saveToFile(path string) error {

	file, err := os.Create(path)
	if err != nil {

		return err
	}

	bytes := make([]byte, cinfo.levelCount())
	for i := range cinfo.states {

		bytes[i] = byte(cinfo.states[i])
	}

	_, err = file.Write(bytes)
	if err != nil {

		file.Close()
		return err
	}

	file.Close()

	return nil
}

func (cinfo *completionInfo) readFromFile(path string) error {

	file, err := os.Open(path)
	if err != nil {

		return err
	}

	bytes := make([]byte, cinfo.levelCount())
	var n int
	n, err = file.Read(bytes)
	if err != nil {

		return err
	}

	for i, b := range bytes {

		if i < n {

			cinfo.states[i] = int32(b)
		}
	}

	return nil
}

func newCompletionInfo(count int32) *completionInfo {

	cinfo := new(completionInfo)

	cinfo.currentStage = 1
	cinfo.states = make([]int32, count)

	cinfo.sinfo = parseStageInfo("assets/maps")

	return cinfo
}
