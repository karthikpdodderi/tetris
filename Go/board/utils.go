package board

import (
	"fmt"
	"math/rand"
	"strings"
)

func (data *boardData) print() {

	// not printing the first two rows, as its used to spawn new blocks
	rows := []string{}
	// for i := data.numStagingRows; i < data.height; i++ {
	for i := 0; i < data.height; i++ {
		row := data.arena[i]
		rows = append(rows, strings.Join(strings.Split(string(row), ""), " "))
	}
	fmt.Println(strings.Join(rows, "\n"))
}

func (data *boardData) clear() {

	// data.clearLines(data.height - data.numStagingRows)
	data.clearLines(data.height)
}
func (board *boardData) generateRandomBlock() {

	blockTypeList := []blockType{i_bl, o_bl, j_bl, l_bl, s_bl, z_bl, t_bl}
	//blockTypeList := []blockType{i_bl, i_bl, i_bl, i_bl, i_bl, i_bl, i_bl}
	blockType := blockTypeList[rand.Intn(7)]

	switch blockType {
	case i_bl:
		board.currentBlock = block{
			origin:      pos{rowNum: 0, colNum: (board.width / 2) - 2},
			blockType:   i_bl,
			relativePos: [4]pos{{rowNum: 2, colNum: 0}, {rowNum: 2, colNum: 1}, {rowNum: 2, colNum: 2}, {rowNum: 2, colNum: 3}},
		}

	case o_bl:
		board.currentBlock = block{
			origin:      pos{rowNum: 0, colNum: (board.width / 2) - 2},
			blockType:   o_bl,
			relativePos: [4]pos{{rowNum: 1, colNum: 1}, {rowNum: 1, colNum: 2}, {rowNum: 2, colNum: 1}, {rowNum: 2, colNum: 2}},
		}

	case j_bl:
		board.currentBlock = block{
			origin:      pos{rowNum: 0, colNum: (board.width / 2) - 2},
			blockType:   j_bl,
			relativePos: [4]pos{{rowNum: 0, colNum: 2}, {rowNum: 1, colNum: 0}, {rowNum: 1, colNum: 1}, {rowNum: 1, colNum: 2}},
		}

	case l_bl:
		board.currentBlock = block{
			origin:      pos{rowNum: 0, colNum: (board.width / 2) - 2},
			blockType:   l_bl,
			relativePos: [4]pos{{rowNum: 0, colNum: 0}, {rowNum: 1, colNum: 0}, {rowNum: 1, colNum: 1}, {rowNum: 1, colNum: 2}},
		}

	case s_bl:
		board.currentBlock = block{
			origin:      pos{rowNum: 0, colNum: (board.width / 2) - 2},
			blockType:   s_bl,
			relativePos: [4]pos{{rowNum: 0, colNum: 1}, {rowNum: 0, colNum: 2}, {rowNum: 1, colNum: 0}, {rowNum: 1, colNum: 1}},
		}

	case z_bl:
		board.currentBlock = block{
			origin:      pos{rowNum: 0, colNum: (board.width / 2) - 2},
			blockType:   z_bl,
			relativePos: [4]pos{{rowNum: 0, colNum: 0}, {rowNum: 0, colNum: 1}, {rowNum: 1, colNum: 1}, {rowNum: 1, colNum: 2}},
		}

	case t_bl:
		board.currentBlock = block{
			origin:      pos{rowNum: 0, colNum: (board.width / 2) - 2},
			blockType:   t_bl,
			relativePos: [4]pos{{rowNum: 0, colNum: 0}, {rowNum: 0, colNum: 1}, {rowNum: 0, colNum: 2}, {rowNum: 1, colNum: 1}},
		}

	default:
		panic(fmt.Sprintf("Invalid block %v ", blockType))
	}

}

func (board *boardData) clearCurentBlock() {
	absPos := getAbsoulutePosition(board.currentBlock.origin, board.currentBlock.relativePos)
	for i := 0; i < 4; i++ {
		board.arena[absPos[i].rowNum][absPos[i].colNum] = board.bgChar
	}
}

func (board *boardData) addCurentBlock() {
	absPos := getAbsoulutePosition(board.currentBlock.origin, board.currentBlock.relativePos)
	board.logger.Log(fmt.Sprintf("abs pos of curr block : %v ", absPos))
	for i := 0; i < 4; i++ {
		board.arena[absPos[i].rowNum][absPos[i].colNum] = board.blockChar
	}
}

func (board *boardData) clearCompletedRows(startingRowPos int) int {
	rowOffset := 0
	for i := startingRowPos; i >= 0; i-- {
		if board.isRowComplete(i) {
			rowOffset++
			continue
		}
		board.updateRow(i, rowOffset)
	}
	return rowOffset
}

func (board *boardData) updateRow(rowNum int, rowOffset int) {
	for i := 0; i < board.width; i++ {
		board.arena[rowNum+rowOffset][i] = board.arena[rowNum][i]
	}
}

func (board *boardData) isRowComplete(rowNum int) bool {
	for i := 0; i < board.width; i++ {
		if board.arena[rowNum][i] != board.blockChar {
			return false
		}
	}
	return true
}

func (board *boardData) getDistanceToDrop(blockPositions [4]pos) int {
	lowestPositions := getLowestPositions(blockPositions)
	board.logger.Log(fmt.Sprintf("lowest positions : %v ", lowestPositions))
	lowestDistance := board.height
	for _, lowestPosition := range lowestPositions {
		currDist := board.getClosestPositionBelowDistance(lowestPosition)
		board.logger.Log(fmt.Sprintf("curr dist : %v", currDist))
		if lowestDistance > currDist{
			lowestDistance = currDist
		}
	}
	return lowestDistance
}

func (board *boardData) getClosestPositionBelowDistance(position pos) int {
	dist := 0
	for i := position.rowNum+1; i < board.height; i++ {
		if board.isPositionCollided(pos{rowNum: i, colNum: position.colNum}) {
			break
		}
		dist++
	}
	return dist
}

func getLowestPositions(blockPositions [4]pos) []pos {
	lowestPositions := []pos{}
	columnToLowestRow := map[int]int{}
	for _, blockPosition := range blockPositions {
		lowestRow, exists := columnToLowestRow[blockPosition.colNum]
		if !exists || blockPosition.rowNum > lowestRow {
			columnToLowestRow[blockPosition.colNum] = blockPosition.rowNum
		}
	}
	for col, row := range columnToLowestRow{
		lowestPositions = append(lowestPositions, pos{rowNum: row, colNum: col})
	}
	return lowestPositions
}

func getAbsoulutePosition(origin pos, relativePos [4]pos) [4]pos {
	absoultePositions := [4]pos{}
	for i := 0; i < 4; i++ {
		absoultePositions[i] = pos{
			rowNum: origin.rowNum + relativePos[i].rowNum,
			colNum: origin.colNum + relativePos[i].colNum,
		}
	}
	return absoultePositions
}

func (block *block) getRightRotateRelativePosition() [4]pos {
	newRelativePosition := [4]pos{}
	switch block.blockType {
	case i_bl, o_bl:
		for i := 0; i < 4; i++ {
			newRelativePosition[i] = fourRotateRight(block.relativePos[i])
		}
	case l_bl, j_bl, t_bl, s_bl, z_bl:
		for i := 0; i < 4; i++ {
			newRelativePosition[i] = threeRotateRight(block.relativePos[i])
		}
	default:
		panic(fmt.Sprintf("Invalid block %v ", block.blockType))
	}
	return newRelativePosition
}

func (block *block) getLeftRotateRelativePosition() [4]pos {
	newRelativePosition := [4]pos{}
	switch block.blockType {
	case i_bl, o_bl:
		for i := 0; i < 4; i++ {
			newRelativePosition[i] = fourRotateLeft(block.relativePos[i])
		}
	case l_bl, j_bl, t_bl, s_bl, z_bl:
		for i := 0; i < 4; i++ {
			newRelativePosition[i] = threeRotateLeft(block.relativePos[i])
		}
	default:
		panic(fmt.Sprintf("Invalid block %v ", block.blockType))
	}
	return newRelativePosition
}

func fourRotateRight(p pos) pos {
	return pos{
		rowNum: p.colNum,
		colNum: 3 - p.rowNum,
	}
}

func fourRotateLeft(p pos) pos {
	return pos{
		rowNum: 3 - p.colNum,
		colNum: p.rowNum,
	}
}

func threeRotateRight(p pos) pos {
	return pos{
		rowNum: p.colNum,
		colNum: 2 - p.rowNum,
	}
}

func threeRotateLeft(p pos) pos {
	return pos{
		rowNum: 2 - p.colNum,
		colNum: p.rowNum,
	}
}

func (board *boardData) isBlockCollided(b block) bool {
	absPositions := getAbsoulutePosition(b.origin, b.relativePos)
	for i := 0; i < 4; i++ {
		if board.isPositionCollided(absPositions[i]) {
			return true
		}
	}
	return false
}

func (board *boardData) isPositionCollided(position pos) bool {

	if position.rowNum > board.height-1 {
		return true
	}

	if position.colNum < 0 {
		return true
	}

	if position.colNum >= board.width {
		return true
	}

	if board.arena[position.rowNum][position.colNum] == board.blockChar {
		return true
	}

	return false
}

func getLowestPos(p [4]pos) pos {
	lowest := p[0]
	for i := 1; i < 4; i++ {
		if p[i].rowNum > lowest.rowNum {
			lowest = p[i]
		}
	}
	return lowest
}
