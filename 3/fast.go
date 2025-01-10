package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

func FastSearch(out io.Writer) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	r := regexp.MustCompile("@")
	seenBrowsers := make(map[string]struct{})
	uniqueBrowsers := 0
	reAndroid := regexp.MustCompile("Android")
	reMSIE := regexp.MustCompile("MSIE")
	var foundUsers strings.Builder

	scanner := bufio.NewScanner(file)

	for i := 0; scanner.Scan(); i++ {
		line := scanner.Text()

		var user map[string]interface{}
		err := json.Unmarshal([]byte(line), &user)
		if err != nil {
			continue
		}

		isAndroid := false
		isMSIE := false

		browsers, ok := user["browsers"].([]interface{})
		if !ok {
			// log.Println("cant cast browsers")
			continue
		}

		for _, browserRaw := range browsers {
			browser, ok := browserRaw.(string)
			if !ok {
				// log.Println("cant cast browser to string")
				continue
			}

			if reAndroid.MatchString(browser) {
				isAndroid = true

				if _, seen := seenBrowsers[browser]; !seen {
					seenBrowsers[browser] = struct{}{}
					uniqueBrowsers++
				}
			}

			if reMSIE.MatchString(browser) {
				isMSIE = true

				if _, seen := seenBrowsers[browser]; !seen {
					seenBrowsers[browser] = struct{}{}
					uniqueBrowsers++
				}
			}
		}

		if !(isAndroid && isMSIE) {
			continue
		}

		// log.Println("Android and MSIE user:", user["name"], user["email"])
		email := r.ReplaceAllString(user["email"].(string), " [at] ")
		foundUsers.WriteString(fmt.Sprintf("[%d] %s <%s>\n", i, user["name"], email))
	}

	fmt.Fprintln(out, "found users:\n"+foundUsers.String())
	fmt.Fprintln(out, "Total unique browsers", len(seenBrowsers))
}
