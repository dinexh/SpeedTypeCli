package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"math/rand"
	"strings"
	"time"

	"github.com/fatih/color"
)

const (
	testDurationSeconds = 30 // Test duration in seconds
	uiWidth            = 60 // Width of the UI elements
)

// Function to show clear screen
func clearScreen() {
	cmd := exec.Command("clear") // For Unix-based systems (Linux/macOS)
	if err := cmd.Run(); err != nil {
		fmt.Println("Error clearing screen:", err)
	}
}

// Function to show intro screen
func showIntro() {
	clearScreen()

	cyan := color.New(color.FgCyan).Add(color.Bold)
	yellow := color.New(color.FgYellow)
	red := color.New(color.FgRed)

	fmt.Println(strings.Repeat("=", uiWidth))
	cyan.Println(centerText("🐒 MonkeyType CLI 🐒", uiWidth))
	fmt.Println(strings.Repeat("=", uiWidth))
	
	yellow.Println("\n" + centerText("Test Duration: 30 seconds", uiWidth))
	yellow.Println(centerText("Type as many words as you can!", uiWidth))
	fmt.Println()
	red.Println(centerText("Press 'R' to refresh | CTRL+C to exit", uiWidth))
	red.Println(centerText("Press Enter to start...", uiWidth))
	fmt.Scanln()
}

// Add this helper function for centering text
func centerText(text string, width int) string {
	padding := (width - len(text)) / 2
	return strings.Repeat(" ", padding) + text
}

// Function to generate a random sentence using random words
func getRandomSentence() string {
	words := []string{
		"quick", "fox", "brown", "jumped", "over", "the", "lazy", "dog", "speed", "typing", "is", "fun", "practice", "accuracy", "keyboard", "fingers", "hands", "code", "challenge", "efficiency", "debug", "optimizing", "programming", "language", "testing", "output", "keyboard", "compile", "software", "developer", "project", "platform", "task", "assignment", "focus", "goal", "habit", "practice", "endurance",
	}

	rand.Seed(time.Now().UnixNano())
	numWords := rand.Intn(10) + 10 // Create a sentence with 15 to 30 words

	var sentence []string
	for i := 0; i < numWords; i++ {
		sentence = append(sentence, words[rand.Intn(len(words))])
	}

	return strings.Join(sentence, " ")
}

// Function to calculate typing speed and accuracy
func calculateSpeedAndAccuracy(start time.Time, typedText, originalText string) {
	duration := time.Since(start)
	wordsTyped := float64(len(strings.Fields(typedText)))
	secondsTaken := duration.Seconds()
	
	// Handle the case where no input was provided (time's up)
	if typedText == "" {
		wordsTyped = 0
		secondsTaken = float64(testDurationSeconds)
	}
	
	wpm := (wordsTyped / secondsTaken) * 60

	var correctChars, totalChars int
	if typedText != "" {
		for i := 0; i < len(typedText) && i < len(originalText); i++ {
			if typedText[i] == originalText[i] {
				correctChars++
			}
			totalChars++
		}
	}

	accuracy := 0.0
	if totalChars > 0 {
		accuracy = (float64(correctChars) / float64(totalChars)) * 100
	}

	clearScreen()
	green := color.New(color.FgGreen).Add(color.Bold)
	yellow := color.New(color.FgYellow)

	fmt.Println(strings.Repeat("=", uiWidth))
	green.Println(centerText("🎯 Results 🎯", uiWidth))
	fmt.Println(strings.Repeat("=", uiWidth))
	yellow.Printf("\n%s\n", centerText(fmt.Sprintf("Speed: %.0f WPM", wpm), uiWidth))
	yellow.Printf("%s\n", centerText(fmt.Sprintf("Accuracy: %.1f%%", accuracy), uiWidth))
	fmt.Println(strings.Repeat("=", uiWidth))
}

// Add this function to manage cursor position
func moveCursorToInput() {
	// Move cursor to input line (6 lines from bottom)
	fmt.Print("\033[6A")  // Move up 6 lines
	fmt.Print("\033[2K")  // Clear the line
}

// Update the displayTimer function
func displayTimer(remainingSeconds int) {
	cyan := color.New(color.FgCyan).Add(color.Bold)
	
	// Save cursor position
	currentPos := "\033[s"
	// Restore cursor position
	restorePos := "\033[u"
	
	fmt.Print(currentPos)
	fmt.Print("\033[H") // Move to top
	fmt.Print("\033[2K") // Clear line
	timerBar := fmt.Sprintf("Time: %02d seconds", remainingSeconds)
	cyan.Println(centerText(timerBar, uiWidth))
	fmt.Print(restorePos)
}

// Function to start the typing test
func startTest() {
	clearScreen()
	sentence := getRandomSentence()

	blue := color.New(color.FgBlue)
	green := color.New(color.FgGreen)

	// Create channels for timer and input
	timerDone := make(chan bool)
	inputDone := make(chan string)
	
	fmt.Println(strings.Repeat("=", uiWidth))
	blue.Println("\nType this:")
	green.Println(sentence)
	fmt.Println("\nYour input:")
	fmt.Print("> ") // Add a prompt
	
	// Start the timer in a goroutine
	go func() {
		for i := testDurationSeconds; i >= 0; i-- {
			displayTimer(i)
			if i == 0 {
				timerDone <- true
				return
			}
			time.Sleep(time.Second)
		}
	}()

	// Start getting input in a goroutine
	go func() {
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		inputDone <- input
	}()

	// Wait for either timer completion or user input
	startTime := time.Now()
	select {
	case <-timerDone:
		moveCursorToInput()
		fmt.Println("\n\nTime's up!")
		time.Sleep(time.Second)
		calculateSpeedAndAccuracy(startTime, "", sentence)
	case typedText := <-inputDone:
		calculateSpeedAndAccuracy(startTime, typedText, sentence)
	}
}

// Handle key press to refresh
func waitForRefresh() {
	var ch byte
	for {
		// Define color for prompt
		red := color.New(color.FgRed) // Declare red color

		// Prompt user with color
		red.Print("Press 'R' to refresh the test or CTRL+C to exit: ")
		fmt.Scanf("%c\n", &ch)
		if ch == 'r' || ch == 'R' {
			clearScreen()
			startTest()
		}
	}
}

func main() {
	showIntro()
	startTest()
	waitForRefresh()
}
