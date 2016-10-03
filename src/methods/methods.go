package methods

import (
	"bufio"
	"os"
	"sort"
	"strings"

	"../globals"
)

// CheckErr is a quick wrapper, so that I don't have to write numeruos if
// statement checks
func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}

// FindEmails uses a compound regex to find anything that might resemble an email
// in the raw html
func FindEmails(html string, filters []string) []string {
	emails := globals.EmailRE.FindAllString(html, -1)
	filteredEmails := []string{}
	for _, email := range emails {
		if StringInSlice(email, filters) {
			if !StringInSlice(email, filteredEmails) {
				filteredEmails = append(filteredEmails, email)
			}
		}
	}
	sort.Strings(filteredEmails)
	return filteredEmails
}

// FormatURL ensures that a protocol is attached to the url string passed in
// for more control specify your own protocol
func FormatURL(myURL string) string {
	if strings.Contains(myURL, "http://") || strings.Contains(myURL, "https://") {
		return myURL
	}
	return "http://" + myURL
}

// ReadLines opens the specified file and reads each line into a slice
func ReadLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.Trim(scanner.Text(), "")
		line = FormatURL(line)
		lines = append(lines, line)
	}
	sort.Strings(lines)
	return lines, scanner.Err()
}

// RemoveEqualSign Attempts to clean up some of the email data
func RemoveEqualSign(email string) string {
	emailParts := strings.Split(email, "=")
	if len(emailParts) == 1 {
		return emailParts[0]
	}
	return emailParts[1]
}

// StringInSlice quickly checks the list to avoid duplicates
func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if strings.Contains(a, b) {
			return true
		}
	}
	return false
}
