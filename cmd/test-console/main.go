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
	"time"

	"github.com/avian-digital-forensics/timeline-investigator/cmd"
	"github.com/avian-digital-forensics/timeline-investigator/tests/client"
	"github.com/google/uuid"
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
					<p class="lead">{{ .User.Email }}</p>
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
	testSecret := flags.String("test-secret", "super-secret", "Secret for testing Timeline Investigator")
	if err := flags.Parse(args[1:]); err != nil {
		return err
	}

	// Parse the template
	tmpl, err := template.New("Test-console").Parse(tpl)
	if err != nil {
		return err
	}

	httpClient := client.New(*testURL)
	testService := client.NewTestService(httpClient, "")

	// Create the page-struct
	// for displaying data
	type page struct {
		Token string
		User  struct {
			ID    string
			Email string
		}
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

	var pages = make(map[string]*page)

	// Handle the requests @ "/test"
	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		var page page

		query := r.URL.Query()
		user := query.Get("user")
		if len(user) == 0 {
			page.Error = "Please specify username in url: /test?user={username}"
			tmpl.Execute(w, page)
			return
		}

		// Create a new test-user if it isn't a POST-request
		if r.Method != http.MethodPost {
			// First check if the user already exists in the cookie
			if cookie, err := r.Cookie(user); err == nil {
				if savedPage, ok := pages[cookie.Value]; ok {
					tmpl.Execute(w, savedPage)
					return
				}
			}

			// Generate ID for the new user
			uid := uuid.New().String()

			request := client.TestCreateUserRequest{
				Name:     user,
				ID:       uid,
				Email:    user + "@" + uid + ".com",
				Password: uuid.New().String(),
				Secret:   *testSecret,
			}

			created, err := testService.CreateUser(ctx, request)
			if err != nil {
				page.Error = fmt.Sprintf("Failed to create test-user: %s", err.Error())
				tmpl.Execute(w, page)
				return
			}

			// Set the data to the page
			page.User.Email = request.Email
			page.User.ID = uid
			page.Token = created.Token

			// Create a cookie for the user
			//r.AddCookie(&http.Cookie{Name: user, Value: uid, Expires: time.Now().AddDate(0, 0, 1)})
			http.SetCookie(w, &http.Cookie{Name: user, Value: uid, Expires: time.Now().AddDate(0, 0, 1)})

			// Save the page to the hashmap
			pages[uid] = &page

			// return the template
			tmpl.Execute(w, page)
			return
		}

		cookie, err := r.Cookie(user)
		if err != nil {
			page.Error = err.Error()
			tmpl.Execute(w, page)
			return
		}
		savedPage, ok := pages[cookie.Value]
		if !ok {
			page.Error = "No cookie found"
			tmpl.Execute(w, page)
			return
		}

		// Set the request-information to the page (store in-memory)
		savedPage.Request.Endpoint = r.FormValue("endpoint")
		savedPage.Request.Body = r.FormValue("body")

		// Send request to the testURL and return error if it failed
		req, err := http.NewRequest(http.MethodPost, *testURL+savedPage.Request.Endpoint, bytes.NewBuffer([]byte(savedPage.Request.Body)))
		if err != nil {
			savedPage.Error = fmt.Sprintf("Failed to create request: %s", err.Error())
			tmpl.Execute(w, savedPage)
			return
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", savedPage.Token)
		resp, err := httpClient.HTTPClient.Do(req)
		if err != nil {
			savedPage.Error = fmt.Sprintf("Failed to send request: %s", err.Error())
			tmpl.Execute(w, savedPage)
			return
		}
		defer resp.Body.Close()

		// Set the status-code to the page (store in-memory)
		savedPage.Response.StatusCode = resp.StatusCode

		// Decode the response-body and return the error if it failed
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			savedPage.Error = fmt.Sprintf("Failed to decode response: %s", err.Error())
			tmpl.Execute(w, savedPage)
			return
		}

		// Set the response-body to the savedPage (store in-memory)
		savedPage.Response.Body = string(body)
		// Return the savedPage
		tmpl.Execute(w, savedPage)
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
