package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strings"

	"./src/globals"
	"./src/methods"
	"./src/response"
)

func main() {
	var file = flag.String("f", "", "Enter new line deliminated text file")
	var saveFileName = flag.String("s", "contact_bot.csv", "Enter a file name ex. file.csv")
	var wokers = flag.Int("w", 15, "The number of working the run on the script, default is 15")
	flag.Parse()

	if *file != "" {
		globals.Filters = []string{"info", "ads", "sales", "sale", "info", "media", "mediarelations", "media_relations", "contact", "contacts", "contactus", "contact_us", "contact-us", "about_us", "general", "advertise", "support", "systems", "system"}
		globals.EmailRE = regexp.MustCompile(`([a-z0-9!#$%&'*+\/=?^_{|}~-]+(?:\.[a-z0-9!#$%&'*+\/=?^_{|}~-]+)*(@|\sat\s)(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?(\.|\sdot\s))+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?)`)
		replacer := strings.NewReplacer("\"", "",
			"//", "",
			" ", "",
			"'/mailto/", "",
		)
		seedUrls, err := methods.ReadLines(*file)
		methods.CheckErr(err)
		seedUrlsNum := len(seedUrls)
		usr, err := user.Current()
		methods.CheckErr(err)
		parentPath := filepath.Join(usr.HomeDir, "/Desktop/"+*saveFileName)

		file, err := os.Create(parentPath)
		methods.CheckErr(err)
		defer file.Close()

		writer := csv.NewWriter(file)
		defer writer.Flush()
		globals.Header = []string{"URL", "EMAILS"}

		err = writer.Write(globals.Header)
		methods.CheckErr(err)

		globals.Data = make(chan *response.HTTPResponse)
		// Create 10 workers. Adjust up or down as needed.
		for w := 0; w < *wokers; w++ {
			go func() {
				for {
					// Increment request count. Exit at end.
					globals.RequestMu.Lock()
					i := globals.RequestCount
					globals.RequestCount++
					globals.RequestMu.Unlock()
					if i >= seedUrlsNum {
						return
					}
					// Fetch the current URL.
					myURL := seedUrls[i]
					resp, err := http.Get(myURL)
					if err != nil {
						fmt.Println(myURL, err.Error(), i)
						globals.Data <- &response.HTTPResponse{myURL, err.Error()}
						continue
					}

					b, err := ioutil.ReadAll(resp.Body)
					resp.Body.Close()
					if err != nil {
						fmt.Println(myURL, err.Error(), i)
						globals.Data <- &response.HTTPResponse{myURL, err.Error()}
						continue
					}

					myHTML := string(b)
					globals.Data <- &response.HTTPResponse{myURL, myHTML}
				}
			}()
		}

		// Recieve expected number of results
		for i := 0; i < seedUrlsNum; i++ {
			result := <-globals.Data
			emails := methods.FindEmails(result.HTML, globals.Filters)
			// clean and filter the emails before writting to file
			// quickly join, replace, and break the items back up
			filteredEmails := strings.Split(replacer.Replace(strings.Join(emails, ",")), ",")
			var cleanEmails []string
			if filteredEmails[0] != "" {
				for _, e := range filteredEmails {
					parsedEmail := methods.RemoveEqualSign(e)
					if !strings.Contains(parsedEmail, "png") && !strings.Contains(parsedEmail, "jpg") && !strings.Contains(parsedEmail, "jpeg") {
						cleanEmails = append(cleanEmails, parsedEmail)
					}
				}
			}
			if len(cleanEmails) > 0 {
				fmt.Printf("%s, %s, %d\n", result.URL, cleanEmails, i)
				var row = []string{result.URL, strings.Join(cleanEmails, ",")}
				err := writer.Write(row)
				writer.Flush()
				methods.CheckErr(err)
			}
		}
	}

}
