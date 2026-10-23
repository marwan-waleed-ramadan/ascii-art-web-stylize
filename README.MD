# ASCII Art Web Stylize

## Description
Go ASCII art web server with three banners and a responsive CSS interface.

## Authors
Marwan Waleed Ramadan

## Usage
Run `go run .` from this directory, then open http://localhost:8080.
Requires Go. Only standard library packages are used.

## Implementation details: algorithm
GET / displays templates/index.html. POST /ascii-art validates the text and banner,
loads eight rows for each printable ASCII character, joins character rows side by
side, and renders the result. static/style.css handles styling and small screens.
A small JavaScript function updates the character counter.

Study order: main -> home -> generate -> makeArt -> showPage.
main starts the server. home displays the form. generate reads and validates it.
makeArt builds the drawing. showPage displays the HTML and HTTP status.
Status codes: 200 success, 400 invalid input, 404 missing resources, 500 other errors.

Imports match the earlier ascii-art-web project: fmt, html/template, net/http,
os, strings. Uses fmt.Println, os.IsNotExist and strings.Builder instead of
log, errors and bytes. The organized CSS is retained.
