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

	// Define colors
	cyan := color.New(color.FgCyan).Add(color.Underline)
	yellow := color.New(color.FgYellow)
	red := color.New(color.FgRed)

	// Print the intro with colors
	cyan.Println("====================================")
	cyan.Println("        Welcome to MonkeyType CLI!")
	cyan.Println("====================================")
	yellow.Println("You will be shown random words to type.")
	red.Println("Press 'R' to refresh the test or CTRL+C to exit.")
	red.Println("Press Enter to start the typing test.")
	fmt.Scanln()
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
	wpm := (wordsTyped / secondsTaken) * 60

	// Calculate accuracy
	var correctChars, totalChars int
	for i := 0; i < len(typedText) && i < len(originalText); i++ {
		if typedText[i] == originalText[i] {
			correctChars++
		}
		totalChars++
	}

	accuracy := (float64(correctChars) / float64(totalChars)) * 100

	clearScreen()

	// Define color for results
	green := color.New(color.FgGreen) // Declare green color

	// Print results with colors
	green.Printf("\n====================================\n")
	green.Printf("Typing Speed: %.2f words per minute (WPM)\n", wpm)
	green.Printf("Accuracy: %.2f%%\n", accuracy)
	green.Printf("\n====================================\n")
}

// Function to start the typing test
func startTest() {
	// Generate a random sentence to type
	sentence := getRandomSentence()

	// Define colors
	blue := color.New(color.FgBlue)
	green := color.New(color.FgGreen) // Declare green color

	// Display sentence to type with colors
	blue.Println("\nYour sentence to type: ")
	green.Println(sentence)

	// Record the start time
	startTime := time.Now()

	// Get user input
	fmt.Println("\nStart typing:")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	typedText := scanner.Text()

	// Calculate typing speed and accuracy
	calculateSpeedAndAccuracy(startTime, typedText, sentence)
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
