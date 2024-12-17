package main

import (
	"fmt"
	"log"
	"main/board"
	"main/key_logger"
	"main/logger"
	"time"
)

func main() {

	logger, loggerCloser, err := logger.NewFileLogger(fmt.Sprintf("%v", time.Now().UnixMilli()), false)
	if err != nil {
		log.Printf("Error while bringing up a new logger. Error : %v \n", err)
		panic(fmt.Sprintf("error while bringing up a new logger. Error : %v ", err))
	}
	defer loggerCloser()

	leftMoveRune := 'j'
	rightMoveRune := 'l'
	downMoveRune := 'k'
	leftRotateRune := 'a'
	rightRotateRune := 'd'
	quitRune := 'q'
	pausePlayRune := 'p'
	dropRune := 'c'

	keyLogger := key_logger.NewKeyLogger(1*time.Millisecond, 100, leftMoveRune, rightMoveRune, downMoveRune, dropRune)
	keyLogger.Start()
	defer keyLogger.Stop()

	boardMover, boardPrinter, err := board.NewBoard(10, 20, '.', '#', 3, logger)
	if err != nil {
		log.Printf("Error while bringing up a new board. Error : %v \n", err)
		panic(fmt.Sprintf("error while bringing up a new board. Error : %v ", err))
	}

	boardPrinter.Start()
	defer boardPrinter.Stop()

	completeChan := make(chan bool)
	pausePlayChan := make(chan bool)
	isPaused := false

	go func() {
		for {
			select {
			case <-pausePlayChan:
				<-pausePlayChan
			default:
				time.Sleep(300 * time.Millisecond)
				state := boardMover.Down()
				boardPrinter.Refresh()
				if state == board.COMPLETE {
					completeChan <- true
				}
			}
		}
	}()

	go func() {
		for {

			keyPressed := keyLogger.Get()

			if isPaused && keyPressed != pausePlayRune && keyPressed != quitRune {
				continue
			}

			switch keyPressed {
			case quitRune:
				completeChan <- true

			case pausePlayRune:
				pausePlayChan <- true
				isPaused = !isPaused

			case leftMoveRune:
				boardMover.Left()

			case rightMoveRune:
				boardMover.Right()

			case leftRotateRune:
				boardMover.RotateLeft()

			case rightRotateRune:
				boardMover.RotateRight()

			case downMoveRune:
				state := boardMover.Down()
				if state == board.COMPLETE {
					completeChan <- true
				}

			case dropRune:
				state := boardMover.Drop()
				if state == board.COMPLETE {
					completeChan <- true
				}

			}

			boardPrinter.Refresh()
		}
	}()

	<-completeChan
}
