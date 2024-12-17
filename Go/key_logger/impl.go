package key_logger

import (
	"fmt"
	"time"

	"github.com/eiannone/keyboard"
)

type keyLoggerData struct {
	waitTime  time.Duration
	stopper   chan bool
	buffer    chan rune
	leftRune  rune
	rightRune rune
	downRune  rune
	dropRune  rune
}

func NewKeyLogger(waitTime time.Duration, loggerBufferLenght int, leftRune rune, rightRune rune, downRune rune, dropRune rune) KeyLogger {
	return &keyLoggerData{
		waitTime:  waitTime,
		stopper:   make(chan bool),
		buffer:    make(chan rune, loggerBufferLenght),
		leftRune:  leftRune,
		rightRune: rightRune,
		downRune:  downRune,
		dropRune: dropRune,
	}
}

func (data *keyLoggerData) Start() {
	err := keyboard.Open()
	if err != nil {
		panic(fmt.Sprintf("Error while opening the keyboard. Error : %v ", err))
	}
	go func() {
		for {
			select {
			case <-data.stopper:
				return
			default:
				char, key, err := keyboard.GetSingleKey()
				if err != nil {
					panic(fmt.Sprintf("Error while getting single key from keyboard. Error : %v ", err))
				}
				switch key {
				case keyboard.KeyArrowLeft:
					data.buffer <- data.leftRune
				case keyboard.KeyArrowRight:
					data.buffer <- data.rightRune
				case keyboard.KeyArrowDown:
					data.buffer <- data.downRune
				case keyboard.KeySpace:
					data.buffer <- data.dropRune
				default:
					data.buffer <- char
				}
			}
			time.Sleep(data.waitTime)
			// waiting for process, i.e the charecter pressed
			// result in stopping the channel
		}
	}()
}

func (data *keyLoggerData) Stop() {
	data.stopper <- true
}

func (data *keyLoggerData) Get() rune {
	val := <-data.buffer
	return val
}
