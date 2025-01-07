import curses
import time
import random
from typing import List, Tuple
import sys
import math
from random import randint

class MonkeyTypeCLI:
    def __init__(self):
        self.words = [
            "the", "be", "to", "of", "and", "a", "in", "that", "have", "I",
            "it", "for", "not", "on", "with", "he", "as", "you", "do", "at",
            "this", "but", "his", "by", "from", "they", "we", "say", "her",
            "she", "or", "an", "will", "my", "one", "all", "would", "there",
            "their", "what", "so", "up", "out", "if", "about", "who", "get",
            "which", "go", "me", "when", "make", "can", "like", "time", "no",
            "just", "him", "know", "take", "people", "into", "year", "your",
            "good", "some", "could", "them", "see", "other", "than", "then",
            "now", "look", "only", "come", "its", "over", "think", "also",
            "back", "after", "use", "two", "how", "our", "work", "first",
            "well", "way", "even", "new", "want", "because", "any", "these",
            "give", "day", "most", "us"
        ]
        
    def generate_text(self, word_count: int = 23) -> str:
        return " ".join(random.choice(self.words) for _ in range(word_count))

    def init_colors(self):
        curses.start_color()
        # Force black background
        curses.use_default_colors()
        
        # Basic colors
        curses.init_pair(1, curses.COLOR_GREEN, -1)    # Correct
        curses.init_pair(2, curses.COLOR_RED, -1)      # Wrong
        curses.init_pair(3, curses.COLOR_WHITE, -1)    # Not typed yet
        curses.init_pair(4, curses.COLOR_YELLOW, -1)   # Current word
        curses.init_pair(5, curses.COLOR_CYAN, -1)     # Instructions/Stats
        
        # Confetti colors (6-11)
        for i, color in enumerate([
            curses.COLOR_RED,
            curses.COLOR_GREEN,
            curses.COLOR_YELLOW,
            curses.COLOR_BLUE,
            curses.COLOR_MAGENTA,
            curses.COLOR_CYAN
        ]):
            curses.init_pair(6 + i, color, -1)

    def calculate_stats(self, typed_chars: int, correct_chars: int, time_taken: float) -> Tuple[float, float]:
        wpm = (typed_chars / 5) / (time_taken / 60)
        accuracy = (correct_chars / typed_chars * 100) if typed_chars > 0 else 0
        return round(wpm, 2), round(accuracy, 2)

    def run(self, stdscr):
        self.init_colors()
        curses.curs_set(0)  # Hide cursor

        while True:
            stdscr.clear()
            height, width = stdscr.getmaxyx()
            
            # Generate test text
            test_text = self.generate_text()
            words = test_text.split()
            current_word_idx = 0
            current_input = ""
            typed_chars = 0
            correct_chars = 0
            started = False
            completed = False

            # Display title and instructions
            title = "🚀 MONKEY TYPE CLI 🚀"
            stdscr.addstr(0, (width - len(title)) // 2, title, curses.color_pair(5) | curses.A_BOLD)
            instructions = "Type the text below (ESC to quit, R to restart)"
            stdscr.addstr(1, (width - len(instructions)) // 2, instructions, curses.color_pair(5))
            
            # Add horizontal line
            stdscr.addstr(2, 0, "─" * width, curses.color_pair(5))

            # Display test text with padding
            padding = 2
            stdscr.addstr(4, padding, test_text, curses.color_pair(3))
            stdscr.refresh()

            # Add another horizontal line
            stdscr.addstr(6, 0, "─" * width, curses.color_pair(5))

            start_time = time.time()

            while not completed:
                # Display current progress
                stdscr.move(4, padding)
                x_pos = padding
                
                for i, word in enumerate(words):
                    if i < current_word_idx:
                        stdscr.addstr(word, curses.color_pair(1))
                        stdscr.addstr(" ", curses.color_pair(3))
                    elif i == current_word_idx:
                        stdscr.addstr(word, curses.color_pair(4) | curses.A_BOLD)
                        stdscr.addstr(" ", curses.color_pair(3))
                    else:
                        stdscr.addstr(word + " ", curses.color_pair(3))
                    x_pos += len(word) + 1

                # Show current input with a box around it
                input_y = 8
                stdscr.addstr(input_y, padding, "│ Your input: ", curses.color_pair(5))
                stdscr.addstr(" " * (width - padding - 14))
                stdscr.addstr(input_y, padding + 13, current_input)
                
                # Show stats if started
                if started:
                    current_time = time.time() - start_time
                    wpm, accuracy = self.calculate_stats(typed_chars, correct_chars, current_time)
                    stats = f"WPM: {wpm} │ Accuracy: {accuracy}% │ Time: {current_time:.1f}s"
                    stdscr.addstr(10, (width - len(stats)) // 2, stats, curses.color_pair(5) | curses.A_BOLD)

                stdscr.refresh()

                # Get input
                try:
                    ch = stdscr.getch()
                except KeyboardInterrupt:
                    return

                if ch == 27:  # ESC
                    return
                elif ch in (ord('R'), ord('r')):
                    return self.run(stdscr)
                elif ch in (curses.KEY_BACKSPACE, 127, 8):
                    if current_input:
                        current_input = current_input[:-1]
                        if typed_chars > 0:
                            typed_chars -= 1
                elif ch == ord(' '):
                    if not started:
                        started = True
                    if current_input:
                        # Check word
                        if current_input == words[current_word_idx]:
                            correct_chars += len(current_input)
                        typed_chars += len(current_input) + 1  # +1 for space
                        current_word_idx += 1
                        current_input = ""
                        
                        if current_word_idx >= len(words):
                            completed = True
                elif ch in range(32, 127):  # Printable characters
                    if not started:
                        started = True
                    current_input += chr(ch)
                    # Check if this is the last word and it matches
                    if current_word_idx == len(words) - 1 and current_input == words[-1]:
                        current_word_idx += 1
                        typed_chars += len(current_input)
                        correct_chars += len(current_input)
                        completed = True

            # Show completion screen (moved outside the input loop)
            total_time = time.time() - start_time
            final_wpm, final_accuracy = self.calculate_stats(typed_chars, correct_chars, total_time)
            
            # Clear the screen
            stdscr.clear()
            
            # Create confetti
            confetti = Confetti(height, width)
            
            # Animation loop
            animation_start = time.time()
            while time.time() - animation_start < 5:  # 5 seconds animation
                stdscr.clear()
                
                # Update and draw confetti
                confetti.update()
                confetti.draw(stdscr)
                
                # Draw results
                result_title = "🎉 TEST COMPLETED! 🎉"
                stdscr.addstr(height // 2 - 4, (width - len(result_title)) // 2, 
                             result_title, curses.color_pair(5) | curses.A_BOLD)
                
                stats = [
                    f"Final WPM: {final_wpm}",
                    f"Final Accuracy: {final_accuracy}%",
                    f"Total Time: {total_time:.1f}s"
                ]
                
                for i, stat in enumerate(stats):
                    stdscr.addstr(height // 2 - 2 + i, (width - len(stat)) // 2, 
                                 stat, curses.color_pair(5))
                
                prompt = "Press any key to try again, ESC to quit"
                stdscr.addstr(height // 2 + 2, (width - len(prompt)) // 2, 
                             prompt, curses.color_pair(5))
                
                stdscr.refresh()
                time.sleep(0.05)
                
                # Check for key press to exit early
                stdscr.nodelay(1)
                if stdscr.getch() == 27:  # ESC
                    return
                stdscr.nodelay(0)
            
            # Wait for final keypress after animation
            if stdscr.getch() == 27:  # ESC
                return

class Confetti:
    def __init__(self, height, width):
        self.height = height
        self.width = width
        self.particles = []
        self.chars = ['*', '•', '×', '+', '✦', '✶', '★']
        self.colors = [
            curses.COLOR_RED,
            curses.COLOR_GREEN,
            curses.COLOR_YELLOW,
            curses.COLOR_BLUE,
            curses.COLOR_MAGENTA,
            curses.COLOR_CYAN
        ]
        self.create_particles()

    def create_particles(self):
        for _ in range(50):  # Number of confetti particles
            self.particles.append({
                'x': randint(0, self.width-1),
                'y': randint(-self.height, 0),
                'char': self.chars[randint(0, len(self.chars)-1)],
                'color': self.colors[randint(0, len(self.colors)-1)],
                'speed': randint(1, 3) / 2
            })

    def update(self):
        for p in self.particles:
            p['y'] += p['speed']
            if p['y'] >= self.height:
                p['y'] = 0
                p['x'] = randint(0, self.width-1)

    def draw(self, stdscr):
        for p in self.particles:
            if 0 <= p['y'] < self.height and 0 <= p['x'] < self.width:
                try:
                    stdscr.addstr(int(p['y']), int(p['x']), p['char'],
                                curses.color_pair(6 + self.colors.index(p['color'])))
                except curses.error:
                    pass

if __name__ == "__main__":
    try:
        curses.wrapper(MonkeyTypeCLI().run)
    except KeyboardInterrupt:
        sys.exit(0)
