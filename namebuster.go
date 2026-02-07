package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
)

var usage = `
  Usage: namebuster <text|url|file>

  Example (names): namebuster "John Doe" > usernames.txt
  Example (single): namebuster "admin" > usernames.txt
  Example (url):    namebuster https://example.com > usernames.txt
  Example (file):   namebuster employees.txt > usernames.txt
`

// Namebuster acts as the primary method for the CLI tool.
func Namebuster(input string) []string {
	var parsedNames []string
	var result []string
	var names []string

	// 1. Determine Input Type and Extract Content
	if isValidUrl(input) {
		content := fetchSiteContent(input)
		// For URLs, we still scrape for patterns because HTML is messy
		names = findNamesInText(content)
	} else if _, err := os.Stat(input); err == nil {
		// FIXED: Read file line-by-line instead of Regex scraping
		fileNames, err := readLines(input)
		if err == nil {
			names = fileNames
		} else {
			fmt.Printf("[!] Error reading file: %v\n", err)
		}
	} else {
		// Treat the input directly as the name
		names = []string{input}
	}

	if len(names) == 0 {
		fmt.Println("[!] No names found/extracted.")
		return result
	}

	// 2. Generate Usernames
	for _, name := range names {
		// Clean up the name
		name = strings.TrimSpace(name)
		if name == "" { 
			continue 
		}

		// Avoid processing the exact same string twice
		if contains(parsedNames, name) {
			continue
		}

		parsedNames = append(parsedNames, name)
		result = append(result, generateUsernames(name)...)
	}

	return result
}

// --- Helper Functions ---

func isValidUrl(toTest string) bool {
	return strings.HasPrefix(toTest, "http://") || strings.HasPrefix(toTest, "https://")
}

func fetchSiteContent(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("[!] Could not fetch URL: %v\n", err)
		return ""
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ""
	}
	return string(body)
}

// readLines reads a whole file into memory and returns a slice of its lines.
func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			lines = append(lines, line)
		}
	}
	return lines, scanner.Err()
}

func findNamesInText(text string) []string {
	// A simple heuristic for WEBSITES only: Look for "Capitalized Capitalized"
	r := regexp.MustCompile(`[A-Z][a-z]+\s[A-Z][a-z]+`)
	return r.FindAllString(text, -1)
}

// --- Core Logic ---

func contains(items []string, target string) bool {
	for _, item := range items {
		if item == target {
			return true
		}
	}
	return false
}

func combineNames(left []string, right []string) []string {
	left = append(left, addSeparators(left)...)
	return stringProduct(left, right)
}

func addSeparators(nameList []string) []string {
	var result []string
	for _, element := range nameList {
		result = append(result, element+".")
		result = append(result, element+"_")
		result = append(result, element+"-")
	}
	return result
}

func stringProduct(left []string, right []string) []string {
	var result []string
	for _, elementL := range left {
		for _, elementR := range right {
			result = append(result, elementL+elementR)
		}
	}
	return result
}

func generateUsernames(name string) []string {
	var result []string

	// Clean extra whitespace
	name = strings.TrimSpace(name)
	splitNames := strings.Fields(name) 

	if len(splitNames) == 0 {
		return result
	}

	// Handle Single Word Input (e.g., "bob")
	if len(splitNames) == 1 {
		word := splitNames[0]
		return []string{
			strings.ToLower(word),
			strings.ToUpper(word),
			strings.Title(strings.ToLower(word)),
		}
	}

	// Handle Full Name Input (First Last)
	firstName := splitNames[0]
	lastName := splitNames[1]

	// Common first name variations
	firstNames := []string{
		strings.ToLower(firstName),
		strings.Title(strings.ToLower(firstName)),
		strings.ToUpper(firstName),
	}

	// Common last name variations
	lastNames := []string{
		strings.ToLower(lastName),
		strings.Title(strings.ToLower(lastName)),
		strings.ToUpper(lastName),
	}

	// Add basic names
	result = append(result, firstNames...)
	result = append(result, lastNames...)
	
	// Prepare First Initial variations
	fInitials := []string{strings.ToLower(string(firstName[0])), strings.ToUpper(string(firstName[0]))}
	
	// Prepare Last Initial variations
	lInitials := []string{strings.ToLower(string(lastName[0])), strings.ToUpper(string(lastName[0]))}

	// 1 -- Full first + Full last
	result = append(result, combineNames(firstNames, lastNames)...)
	
	// 2 -- First Initial + Full Last
	result = append(result, combineNames(fInitials, lastNames)...)

	// 3 -- Full First + Last Initial
	result = append(result, combineNames(firstNames, lInitials)...)

	// 4 -- Full Last + First Initial
	result = append(result, combineNames(lastNames, fInitials)...)
	
	// 5 -- Full Last + Full First
	result = append(result, combineNames(lastNames, firstNames)...)

	return result
}

func main() {
	if len(os.Args) < 2 {
		fmt.Print(usage)
		os.Exit(1)
	}

	// Handle potential multi-word args passed without quotes
	input := strings.Join(os.Args[1:], " ")

	result := Namebuster(input)

	if len(result) > 0 {
		for _, value := range result {
			fmt.Println(value)
		}
	} else {
		fmt.Println("No usernames generated.")
	}
}