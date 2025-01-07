package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/nsf/termbox-go"
)

const timeLimit = 60 // Time limit in seconds

var sampleTexts = []string{
	"Go is an open-source programming language that makes it easy to build simple, reliable, and efficient software.",
	"Concurrency is not parallelism.",
	"Don't communicate by sharing memory; share memory by communicating.",
}

func generateText() string {
	return sampleTexts[rand.Intn(len(sampleTexts))]
}

func drawText(x, y int, text string, fg, bg termbox.Attribute) {
	for i, ch := range text {
		termbox.SetCell(x+i, y, ch, fg, bg)
	}
}

func drawScreen(input string, words []string, wordIdx int, elapsed int, completed bool) {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	// Title
	drawText(2, 1, "🚀 MONKEY TYPE CLI 🚀", termbox.ColorGreen|termbox.AttrBold, termbox.ColorDefault)

	// Time left
	timeLeft := timeLimit - elapsed
	timeText := fmt.Sprintf("Time Left: %ds", timeLeft)
	drawText(2, 3, timeText, termbox.ColorYellow, termbox.ColorDefault)

	// Sample text
	sample := strings.Join(words, " ")
	drawText(2, 5, sample, termbox.ColorWhite, termbox.ColorDefault)

	// Current input
	if completed {
		drawText(2, 7, "Test Completed! Press Cmd+R to Restart or ESC to Quit.", termbox.ColorGreen, termbox.ColorDefault)
	} else {
		inputText := fmt.Sprintf("Your Input: %s", input)
		drawText(2, 7, inputText, termbox.ColorCyan, termbox.ColorDefault)
	}

	termbox.Flush()
}

func main() {
	rand.Seed(time.Now().UnixNano())

	// Initialize termbox
	if err := termbox.Init(); err != nil {
		fmt.Println("Failed to initialize termbox:", err)
		return
	}
	defer termbox.Close()

	for {
		// Generate sample text
		words := strings.Split(generateText(), " ")
		wordIdx := 0
		input := ""
		completed := false
		start := time.Now()

		for {
			elapsed := int(time.Since(start).Seconds())

			// If time limit exceeded
			if elapsed >= timeLimit {
				completed = true
			}

			drawScreen(input, words, wordIdx, elapsed, completed)

			if completed {
				break
			}

			event := termbox.PollEvent()
			switch event.Type {
			case termbox.EventKey:
				switch event.Key {
				case termbox.KeyEsc:
					return // Exit program
				case termbox.KeyCtrlR: // Restart program (Cmd+R is Ctrl+R in termbox)
					main()
				case termbox.KeyBackspace, termbox.KeyBackspace2:
					if len(input) > 0 {
						input = input[:len(input)-1]
					}
				case termbox.KeyEnter, termbox.KeySpace:
					if strings.TrimSpace(input) == words[wordIdx] {
						wordIdx++
						input = ""
						if wordIdx >= len(words) {
							completed = true
						}
					} else {
						input += " " // Add space visually
					}
				default:
					if event.Ch != 0 {
						input += string(event.Ch)
					}
				}
			}
		}
	}
}
