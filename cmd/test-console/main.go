package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/avian-digital-forensics/timeline-investigator/cmd"
)

const tpl = `
<!doctype html>
<html lang="en">
	<head>
		<!-- Required meta tags -->
		<meta charset="utf-8">
		<meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">

		<!-- Bootstrap CSS -->
		<link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/4.0.0/css/bootstrap.min.css" integrity="sha384-Gn5384xqQ1aoWXA+058RXPxPg6fy4IWvTNh0E263XmFcJlSAwiGgFAW/dAiS6JXm" crossorigin="anonymous">

		<title>Timeline Investigator</title>
	</head>
	<body>
		<div class="container">
			<div class="row">
				<div class="col-lg-12 text-center">
					<h1 class="mt-5">Timeline Investigator</h1>
					<p class="lead">Test-console</p>
				</div>
			</div>
		</div>

		<div class="container">
			<div class="row">
				<div class="col">
					<h2>Request</h2>
					<form method="POST">
						<div class="mb-3">
							<label for="endpoint" class="form-label">Endpoint</label>
							<input type="text" class="form-control" id="endpoint" name="endpoint" value="{{ .Request.Endpoint }}">
						</div>
						<div class="mb-3">
							<label for="body" class="form-label">Body</label>
							<textarea class="form-control" id="body" name="body" rows="10">{{ .Request.Body }}</textarea>
						</div>
						{{if .Error}}
						<div>
							<p style="color:red">{{.Error}}</p>
						</div>
						{{end}}
						<div class="col-auto">
							<button type="submit" class="btn btn-primary mb-3">Send request</button>
						</div>
					</form>
				</div>
				<div class="col">
					<h2>Response</h2>
					<div class="mb-3">
						Status code <h3>{{ .Response.StatusCode }}</h3>
					</div>
					<div class="mb-3">
						<label for="body" class="form-label">Body</label>
						<div class="card">
							<div class="card-body">
								{{ .Response.Body }}
							</div>
						</div>
					</div>
				</div>
			</div>
		</div>
	</body>
</html>`

func main() {
	ctx := cmd.ContextWithSignal(context.Background())
	if err := run(ctx, os.Args, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, args []string, stdout io.Writer) error {
	// Get cli arguments
	flags := flag.NewFlagSet(args[0], flag.ExitOnError)
	address := flags.String("address", ":8080", "address to listen on")
	testURL := flags.String("ti-url", "http://localhost:8000/api/", "URL for testing Timeline Investigator")
	if err := flags.Parse(args[1:]); err != nil {
		return err
	}

	// Parse the template
	tmpl, err := template.New("Test-console").Parse(tpl)
	if err != nil {
		return err
	}

	// Create the page-struct
	// for displaying data
	var page struct {
		Request struct {
			Endpoint string
			Body     string
		}
		Response struct {
			StatusCode int
			Body       string
		}
		Error string
	}

	// Set an example-endpoint for placeholder
	page.Request.Endpoint = "CaseService.New"

	// Handle the requests @ "/test"
	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		// Return the template immidieatly if it isn't a POST-request
		if r.Method != http.MethodPost {
			tmpl.Execute(w, page)
			return
		}

		// Set the request-information to the page (store in-memory)
		page.Request.Endpoint = r.FormValue("endpoint")
		page.Request.Body = r.FormValue("body")

		// Send request to the testURL and return error if it failed
		resp, err := http.Post(*testURL+page.Request.Endpoint, "application/json", bytes.NewBuffer([]byte(page.Request.Body)))
		if err != nil {
			page.Error = fmt.Sprintf("Failed to send request: %s", err.Error())
			tmpl.Execute(w, page)
			return
		}
		defer resp.Body.Close()

		// Set the status-code to the page (store in-memory)
		page.Response.StatusCode = resp.StatusCode

		// Decode the response-body and return the error if it failed
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			page.Error = fmt.Sprintf("Failed to decode response: %s", err.Error())
			tmpl.Execute(w, page)
			return
		}

		// Set the response-body to the page (store in-memory)
		page.Response.Body = string(body)
		// Return the page
		tmpl.Execute(w, page)
	})

	// endpoint for for gke-healthchecks
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("ok"))
	})

	log.Printf("Testing against URL: %s", *testURL)
	log.Printf("Listening @ %s", *address)
	return http.ListenAndServe(*address, nil)
}
