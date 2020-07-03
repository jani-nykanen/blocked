package main

import (
	"os"

	"github.com/jani-nykanen/blocked/src/core"
)

const (
	defaultSaveFilePath = "save.dat"
)

type completionInfo struct {
	states            []int32
	currentStage      int32
	endingPlayedState int32
	// This should go elsewhere, but since this data
	// should be loaded only once, let's put it here...
	sinfo *stageInfoContainer

	enterPressed bool // For this reason, RENAME THIS STRUCT
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

	bytes := make([]byte, cinfo.levelCount()+1)
	for i := range cinfo.states {

		bytes[i] = byte(cinfo.states[i])
	}
	bytes[len(bytes)-1] = byte(cinfo.endingPlayedState)

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

	bytes := make([]byte, cinfo.levelCount()+1)
	var n int
	n, err = file.Read(bytes)
	if err != nil {

		return err
	}

	for i, b := range bytes {

		if int32(i) < core.MinInt32(int32(n), cinfo.levelCount()) {

			cinfo.states[i] = int32(b)
		}
	}

	if int32(n) >= cinfo.levelCount()+1 {

		cinfo.endingPlayedState = int32(bytes[cinfo.levelCount()])
	}

	return nil
}

func (cinfo *completionInfo) checkIfNewEndingObtained() bool {

	min := int32(2)
	for _, s := range cinfo.states {

		if s <= cinfo.endingPlayedState {

			return false
		}

		if s < min {

			min = s
		}
	}

	cinfo.endingPlayedState = core.MaxInt32(cinfo.endingPlayedState, min)

	return true
}

func (cinfo *completionInfo) clear() {

	for i := range cinfo.states {

		cinfo.states[i] = 0
	}
	cinfo.endingPlayedState = 0

}

func newCompletionInfo() *completionInfo {

	cinfo := new(completionInfo)

	cinfo.currentStage = 1
	cinfo.sinfo = parseStageInfo("assets/maps")
	cinfo.states = make([]int32, len(cinfo.sinfo.entries))

	cinfo.endingPlayedState = 0

	cinfo.enterPressed = false

	return cinfo
}
