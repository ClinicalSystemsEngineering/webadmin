package webadmin

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"
)

type webpage struct {
	Title   string
	Heading string
	Body    []string
	Nav     []string
}

//used in the nav element to easily have a bottom navigation area on each page
var webpageurls = []string{"home", "status", "page"}

var queuesize = 0 //the size of the processed message channel

var parsedmsgs = make(chan string, 10000) //message processing channel for xml2tap conversions

var timeoutDuration = 5 * time.Second //read / write timeout duration

//HomePage not yet really implemented
func HomePage(w http.ResponseWriter, req *http.Request) {

	homepage := webpage{Title: "XML2TAP Homepage", Heading: "List of Commands:", Nav: webpageurls}

	tpl, err := template.ParseFiles("index.gohtml")
	if err != nil {
		log.Printf("error parsing index template: %v", err)
	}
	err = tpl.ExecuteTemplate(w, "index.gohtml", homepage)
	if err != nil {
		log.Printf("error executing template index: %v", err)
	}
}

//StatusPage displays the size of the queue for monitoring purposes
//queue values above 100 are considered an error if the queue stays at 100 or above for a prolonged period
func StatusPage(w http.ResponseWriter, req *http.Request) {
	//get the latest queue value
	queuemonitor()

	var queuestatus []string

	//determine if queue size is in error state or not currently hardcoded to 100
	if queuesize <= 100 {
		queuestatus = append(queuestatus, "OK: Current Queue Size:"+strconv.Itoa(queuesize))
	} else {
		queuestatus = append(queuestatus, "ERROR: Current Queue Size:"+strconv.Itoa(queuesize))
	}

	statuspage := webpage{Title: "XML2TAP Statuspage", Heading: "Current Queue Status:", Body: queuestatus, Nav: webpageurls}

	tpl, err := template.ParseFiles("status.gohtml")
	if err != nil {
		log.Printf("error parsing status template: %v", err)
	}
	err = tpl.ExecuteTemplate(w, "status.gohtml", statuspage)
	if err != nil {
		log.Printf("error executing status template: %v", err)
	}

}

//SendPage request page for pin;message to add to the queue for processing
func SendPage(w http.ResponseWriter, req *http.Request) {

	if req.Method == "GET" {
		sendpage := webpage{Title: "XML2TAP Sendpage", Heading: "Input pin and message below:", Nav: webpageurls}

		tpl, err := template.ParseFiles("sendpage.gohtml")
		if err != nil {
			log.Printf("error parsing sendpage template: %v", err)
		}
		err = tpl.ExecuteTemplate(w, "sendpage.gohtml", sendpage)
		if err != nil {
			log.Printf("error executing send template: %v", err)
		}
	} else {
		// put pin and message into the processing queue
		req.ParseForm()
		pin := req.PostFormValue("pin")
		msg := req.PostFormValue("message")
		if pin != "" && msg != "" {
			parsedmsgs <- pin + ";" + msg
			sendpage := webpage{Title: "XML2TAP Sendpage", Heading: "Message submitted. Input pin and message below:", Nav: webpageurls}

			tpl, err := template.ParseFiles("sendpage.gohtml")
			if err != nil {
				log.Printf("error parsing sendpage template: %v", err)
			}
			err = tpl.ExecuteTemplate(w, "sendpage.gohtml", sendpage)
			if err != nil {
				log.Printf("error executing send template: %v", err)
			}
		} else {
			sendpage := webpage{Title: "XML2TAP Sendpage", Heading: "Error with Submission Input Try Again. Input pin and message below:", Nav: webpageurls}

			tpl, err := template.ParseFiles("sendpage.gohtml")
			if err != nil {
				log.Printf("error parsing sendpage template: %v", err)
			}
			err = tpl.ExecuteTemplate(w, "sendpage.gohtml", sendpage)
			if err != nil {
				log.Printf("error executing send template: %v", err)
			}
		}

	}

}

//Webserver is a simplified admin inteface placeholder
func Webserver(portnum string) {
	http.HandleFunc("/", HomePage)
	http.HandleFunc("/home", HomePage)
	http.HandleFunc("/status", StatusPage)
	http.HandleFunc("/page", SendPage)
	http.Handle("/favicon.ico", http.NotFoundHandler())
	for {
		log.Println(http.ListenAndServe(":"+portnum, nil))
	}
}

func queuemonitor() {

	queuesize = len(parsedmsgs)

}
