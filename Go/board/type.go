package board

import (
	"main/logger"
	"sync"
)

type State int

var (
	CONTINUE State = 0
	COMPLETE State = 1
)

type pos struct {
	rowNum int
	colNum int
}

type block struct {
	origin      pos
	relativePos [4]pos
	blockType   blockType
}

type boardData struct {
	blockChar      rune
	bgChar         rune
	arena          [][]rune
	currentBlock   block
	isBlockLocked  bool
	access         sync.Mutex
	numStagingRows int
	score          int
	logger         logger.Logger
	refreshChan    chan bool
	stopChan       chan bool
	clearLines     func(int)
	height         int
	width          int
}

type blockType int

var (
	i_bl blockType = 0
	j_bl blockType = 1
	l_bl blockType = 2
	o_bl blockType = 3
	s_bl blockType = 4
	t_bl blockType = 5
	z_bl blockType = 6
)
