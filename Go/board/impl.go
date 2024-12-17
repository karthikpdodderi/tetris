package board

import (
	"fmt"
	"log"
	"main/logger"
	"main/utils"
	"sync"
)

func NewBoard(width, height int, bgChar rune, blockChar rune, numStagingRows int, logger logger.Logger) (Mover, Printer, error) {

	if height <= 5 {
		log.Println("Height should not be less than or equal to 5")
		return nil, nil, fmt.Errorf("height should not be less than or equal to 5")
	}

	if width <= 8 {
		log.Println("Width should not be less than or equal to 8")
		return nil, nil, fmt.Errorf("width should not be less than or equal to 8")
	}

	height += numStagingRows

	data := boardData{}

	data.height = height
	data.width = width

	arena := [][]rune{}
	for i := 0; i < data.height; i++ {
		row := make([]rune, data.width)
		for j := 0; j < data.width; j++ {
			row[j] = bgChar
		}
		arena = append(arena, row)
	}
	data.arena = arena

	data.access = sync.Mutex{}

	data.bgChar = bgChar
	data.blockChar = blockChar
	data.currentBlock = block{}
	data.isBlockLocked = true
	data.numStagingRows = numStagingRows
	data.logger = logger

	data.refreshChan = make(chan bool, 1)
	data.stopChan = make(chan bool)

	data.clearLines = utils.ClerLines

	data.generateRandomBlock()

	return &data, &data, nil

}

func (data *boardData) Drop() State {

	data.access.Lock()
	defer data.access.Unlock()

	data.clearCurentBlock()
	defer data.addCurentBlock()

	absoluteCurrentPostition := getAbsoulutePosition(data.currentBlock.origin, data.currentBlock.relativePos)
	dropDistance := data.getDistanceToDrop(absoluteCurrentPostition)
	data.logger.Log(fmt.Sprintf("drop distance : %v ", dropDistance))

	if dropDistance == 0 {
		lowestPos := getLowestPos(absoluteCurrentPostition)

		if lowestPos.rowNum < data.numStagingRows {
			return COMPLETE
		}

		data.addCurentBlock()

		rowsCleared := data.clearCompletedRows(lowestPos.rowNum)
		data.score += rowsCleared

		data.generateRandomBlock()

		return CONTINUE
	}

	data.currentBlock.origin.rowNum += dropDistance

	return CONTINUE
}

func (data *boardData) Down() State {

	data.access.Lock()
	defer data.access.Unlock()

	data.clearCurentBlock()
	defer data.addCurentBlock()

	nextBlock := data.currentBlock
	nextBlock.origin.rowNum += 1

	if data.isBlockCollided(nextBlock) {
		lowestPos := getLowestPos(getAbsoulutePosition(data.currentBlock.origin, data.currentBlock.relativePos))

		if lowestPos.rowNum < data.numStagingRows {
			return COMPLETE
		}

		data.addCurentBlock()

		rowsCleared := data.clearCompletedRows(lowestPos.rowNum)
		data.score += rowsCleared

		data.generateRandomBlock()

		return CONTINUE
	}

	data.currentBlock = nextBlock

	return CONTINUE
}

func (data *boardData) Left() {

	data.access.Lock()
	defer data.access.Unlock()

	data.clearCurentBlock()
	defer data.addCurentBlock()

	nextBlock := data.currentBlock
	nextBlock.origin.colNum -= 1

	if !data.isBlockCollided(nextBlock) {
		data.currentBlock = nextBlock
	}
}

func (data *boardData) Right() {

	data.access.Lock()
	defer data.access.Unlock()

	data.clearCurentBlock()
	defer data.addCurentBlock()

	nextBlock := data.currentBlock
	nextBlock.origin.colNum += 1

	if !data.isBlockCollided(nextBlock) {
		data.currentBlock = nextBlock
	}

}

func (data *boardData) RotateLeft() {

	data.access.Lock()
	defer data.access.Unlock()

	data.clearCurentBlock()
	defer data.addCurentBlock()

	nextBlock := data.currentBlock
	nextBlock.relativePos = nextBlock.getLeftRotateRelativePosition()

	if !data.isBlockCollided(nextBlock) {
		data.currentBlock = nextBlock
	}

}

func (data *boardData) RotateRight() {

	data.access.Lock()
	defer data.access.Unlock()

	data.clearCurentBlock()
	defer data.addCurentBlock()

	nextBlock := data.currentBlock
	nextBlock.relativePos = nextBlock.getRightRotateRelativePosition()

	if !data.isBlockCollided(nextBlock) {
		data.currentBlock = nextBlock
	}

}

func (data *boardData) GetScore() int {
	return data.score
}

func (data *boardData) Start() {

	data.access.Lock()
	defer data.access.Unlock()

	data.print()

	go func() {
		defer func() {
			data.clear()
			fmt.Printf("Score :  %d \n", data.score)
			fmt.Println("Press any key to exit ... ")
		}()
		for {
			select {
			case <-data.refreshChan:
				data.clear()
				data.print()
			case <-data.stopChan:
				return
			}
		}
	}()

}

func (data *boardData) Refresh() {

	data.access.Lock()
	defer data.access.Unlock()

	data.refreshChan <- true
}

func (data *boardData) Stop() {

	data.access.Lock()
	defer data.access.Unlock()

	data.stopChan <- true
}
