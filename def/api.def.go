package def

// Service is the main-service
type Service interface {
	// Greet sends a polite greeting
	Greet(GreetRequest) GreetResponse
}

// GreetRequest is the request object for GreeterService.Greet.
type GreetRequest struct {
	// Namee of the person to greet
	//
	// example: "Simon"
	Name string
}

// GreetResponse is the response object containing a
// person's greeting.
type GreetResponse struct {
	// Greeting is a nice message welcoming somebody.
	//
	// example: "Hello Simon"
	Greeting string
}
