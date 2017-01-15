package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func chooseMirrorCountry() (string, error) {
	fmt.Print("Enter a mirror host country (United States): ")
	return readAnswer(validMirrorCountries)
}

func chooseZoneInfo() (string, error) {
	fmt.Print("Enter your timezone (US/Pacific): ")
	return readLine()
}

func chooseLocale() (string, error) {
	fmt.Print("Enter a locale (en_US.UTF-8 UTF-8): ")
	return readAnswer(validLocales)
}

func chooseKeyboardLayout() (string, error) {
	fmt.Print("Enter a keyboard layout (us): ")
	return readLine()
}

func chooseHostname() (string, error) {
	fmt.Print("Enter a hostname: ")
	return readLine()
}

func readAnswer(valid []string) (string, error) {
	var err error
	var ans string

	for {
		ans, err = readLine()
		if err != nil {
			log.Print(err)
			fmt.Println("Try again.")
			continue
		}

		if v := findValid(ans, valid); v != "" {
			ans = v
			break
		}

		fmt.Printf("Invalid answer: %q\n", ans)
		fmt.Println("Try again.")
	}

	return ans, nil
}

func readLine() (string, error) {
	r := bufio.NewReader(os.Stdin)

	s, err := r.ReadString('\n')
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(s), nil
}

func findValid(ans string, valid []string) string {
	ans = strings.ToLower(ans)
	for _, v := range valid {
		if ans == strings.ToLower(v) {
			// Correct case.
			return v
		}
	}

	return ""
}
