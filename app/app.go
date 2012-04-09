package app

import (
	"fmt"
	// "html/template"
	"net/http"
)

func init() {
	http.HandleFunc("/generate", generate)
}

func generate(w http.ResponseWriter, r *http.Request) {
	err := signTemplate.Execute(w, r.FormValue("content"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// var signTemplate = template.Must(template.New("sign").Parse(signTemplateHTML))

// const signTemplateHTML = `
// <html>
//   <body>
//     <p>You wrote:</p>
//     <pre>{{.}}</pre>
//   </body>
// </html>
// `
