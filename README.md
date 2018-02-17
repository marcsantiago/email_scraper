# [DEPRICATED]

A better implementation can be found here
https://github.com/marcsantiago/search_keyword

# email_scraper
Command Line Tool To Scrape Sites For Email Addresses


This is a very simple command line tool, which can be used to scrape hundreds (I've done 100k :-)) of websites very quickly.  The main limiting factor is the number of ports that can be open at the same time.  Often one has to increase the ulimit on one's machine to get around this issue.  For a very quick tutorial on how to do that, please check out this [article](https://viewsby.wordpress.com/2013/01/29/ubuntu-increase-number-of-open-files/), However you may not need to do this. This script has been tested on Mac OSX and Ubuntu 14.04.  I'm sure with some tweaking one can get it to work on Windows.  I leave that to anyone who forks the project.  The script can be built upon and perhaps given time I'll add additional functionality.  Of course for all you gophers, building the binary is probably the best way to go.

----------
# Default Use


    go run email_scraper.go -f  "path_to_file.txt"
   The file is expected to contain newline deliminated URLS. One URL per line.  The program saves a csv file to desktop by default.  That can be changed by providing an additional command.

    go run email_scraper.go -f  "path_to_file.txt" -s "save_file_path.csv"

Lastly, the speed at which the items are parsed can be increased or decreased by specifying the number of go workers you wish to spawn, the default is 15.

    go run email_scraper.go -f  "path_to_file.txt" -s "save_file_path.csv" -w 20

  You may have to play with that number to find you're sweet spot and a lot of it depends on your ulimit.

**Note:**
You can increase the speed of the script by exponentially if you increase Go's runtime, e.g. runtime.GOMAXPROCS(4) #4 cores. However, doing so often means that you have to reduce the number of workers because the ulimit can be hit very quickly causing socket/connectivity errors.
Also, on line 27 I supply a list of filters. That is because I am targeting very specific emails. By removing or modifying that section you can make the script more or less flexible.
