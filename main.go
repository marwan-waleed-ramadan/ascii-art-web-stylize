package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"
)

// Data sent to the HTML page.
type PageData struct {
	Text   string
	Banner string
	Result string
	Error  string
}

func main() {
	// Allow the browser to read the CSS file.
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", home)
	http.HandleFunc("/ascii-art", generate)

	fmt.Println("Open http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Server error:", err)
	}
}

// Render the home page.
func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Bad Request", 400)
		return
	}
	showPage(w, PageData{Banner: "standard"}, 200)
}

// Read the form, validate the data, and generate the result.
func generate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Bad Request", 400)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad Request", 400)
		return
	}

	data := PageData{
		Text:   r.PostFormValue("text"),
		Banner: r.PostFormValue("banner"),
	}
	if data.Banner == "" {
		data.Banner = "standard"
	}

	// Validate the text and banner before reading the file.
	if strings.TrimSpace(data.Text) == "" {
		data.Error = "Please enter some text."
	} else if len([]rune(data.Text)) > 200 {
		data.Error = "Maximum 200 characters."
	} else if data.Banner != "standard" && data.Banner != "shadow" && data.Banner != "thinkertoy" {
		data.Error = "Invalid banner."
	}

	text := strings.ReplaceAll(data.Text, "\r\n", "\n")
	text = strings.ReplaceAll(text, `\n`, "\n")
	for _, character := range text {
		if character != '\n' && (character < 32 || character > 126) {
			data.Error = "Use English letters, numbers and ASCII symbols only."
			break
		}
	}
	if data.Error != "" {
		showPage(w, data, 400)
		return
	}

	// Read the selected banner file.
	font, err := os.ReadFile(data.Banner + ".txt")
	if err != nil {
		status := 500
		data.Error = "Could not read the banner."
		if os.IsNotExist(err) {
			status = 404
			data.Error = "Banner file not found."
		}
		showPage(w, data, status)
		return
	}

	data.Result, err = makeArt(text, string(font))
	if err != nil {
		data.Error = "The banner file is incomplete."
		showPage(w, data, 500)
		return
	}
	showPage(w, data, 200)
}

// Each character has 8 drawing rows, followed by a separator line.
func makeArt(text string, font string) (string, error) {
	font = strings.ReplaceAll(font, "\r\n", "\n")
	fontLines := strings.Split(font, "\n")
	textLines := strings.Split(text, "\n")
	result := ""

	for _, line := range textLines {
		if line == "" {
			result += "\n"
			continue
		}
		for row := 0; row < 8; row++ {
			for _, character := range line {
				// The first character is a space with code 32. The first line in the file is a separator.
				start := 1 + (int(character)-32)*9
				if start < 1 || start+8 > len(fontLines) {
					return "", fmt.Errorf("invalid banner data")
				}
				result += fontLines[start+row]
			}
			result += "\n"
		}
	}
	return result, nil
}

// Place the data in HTML and send the page with the status code.
func showPage(w http.ResponseWriter, data PageData, status int) {
	page, err := template.ParseFiles("templates/index.html")
	if err != nil {
		if os.IsNotExist(err) {
			http.Error(w, "Page not found", 404)
		} else {
			http.Error(w, "Internal Server Error", 500)
		}
		return
	}

	// Render the page before sending it so we can return 500 if an error occurs.
	var output strings.Builder
	if err := page.Execute(&output, data); err != nil {
		http.Error(w, "Internal Server Error", 500)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	fmt.Fprint(w, output.String())
}
