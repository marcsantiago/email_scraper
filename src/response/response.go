package response

// HTTPResponse is a very simple struct to store and pass the data in
// it can be build upon to add features and more data
type HTTPResponse struct {
	URL  string
	HTML string
}
