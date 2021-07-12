package main

import (
	"fmt"
	"net/http"
	"time"
)
// main go routine
func main() {
	links := []string{
		"https://google.com",
		"https://facebook.com",
		"https://stackoverflow.com",
		"https://golang.org",
		"https://amazon.com",
	}
	c := make(chan string) // creating a channel
	for _, link := range links {
		go checkLink(link,c) // 'go' keyword tells to run this function in another (child) go routine
	}

	// infinite loop // this program will check the status of above links forever
	// Option 1: floods the channel
	/*for {
		//fmt.Println(<- c)
		go checkLink(<- c,c) // receiving data from every channel // blocking call
	}*/

	// Option 2: better function readability; throttles the channel
	/*for l := range c { // range works on channel too // here, length of channel will vary based on number of messages in channel
		// time.Sleep(2 * time.Second) // wrong place to add sleep statement // this will block main routine every 2 secs whereas, main routine should be always listening for (and thus, clearing) incoming messages from the channel // this pause will just pile up unread messages inside the channel
		go checkLink(l,c) // receiving data from every channel // blocking call
	}*/

	// Option 3: better function readability; main goroutine and child goroutine modify same variable causing irregularities
	/*for l := range c { // range works on channel too // here, length of channel will vary based on number of messages in channel
		// time.Sleep(2 * time.Second) // wrong place to add sleep statement // this will block main routine every 2 secs whereas, main routine should be always listening for incoming messages from the channel // this pause will just pile up unread messages inside the channel
		// a function literal // anonymous function
		go func() {
			time.Sleep(2 * time.Second) // correct place to add sleep statement
			checkLink(l,c) // here main goroutine and child goroutine (anonymous function) reference the same variable 'l' thus giving wrong results as they both try to modify the same memory location
		}() // '()' here invokes this function
	}*/

	// Option 4: better function readability; correct way
	for l := range c { // range works on channel too // here, length of channel will vary based on number of messages in channel
		// time.Sleep(2 * time.Second) // wrong place to add sleep statement // this will block main routine every 2 secs whereas, main routine should be always listening for incoming messages from the channel // this pause will just pile up unread messages inside the channel
		// a function literal // anonymous function
		go func(link string) {
			time.Sleep(2 * time.Second) // correct place to add sleep statement
			checkLink(link,c) // here child goroutine (anonymous function) references the copy of the variable 'l' i.e., 'link' thus giving correct results
		}(l) // '()' here invokes this function
	}
}

func checkLink(link string, c chan string) {
	// time.Sleep(2 * time.Second) // not a best place to add sleep statement // this will block child go routine for 2 secs while it should be probing the url instantly (even if this will work though)
	_, err := http.Get(link) // a blocking function call; control of program is meanwhile passed to main go routine
	if err != nil {
		fmt.Println("Error: ", err)
		fmt.Printf("%v might be down", link)
		c <- link // sending data into channel // link is of string type
		return
	}
	fmt.Printf("%v is up!\n", link)
	c <- link // sending data into channel
}
