package main

import (
	"html/template"
	"log"
	"net/http"

	"gopl.io/ch4/github"
)

func main() {
	searchFilter := []string{"repo:golang/go", "is:open", "json", "decoder"}
	list := template.Must(template.New("list").Parse(`
		<table>
			<thead>
				<th>Issue</th>
				<th>Reporter</th>
			</thead>
			<tbody>
				{{range .Items}}
				<tr>
					<td><a href="{{.HTMLURL}}">{{.Number}}</a></td>
					<td><a href="{{.User.HTMLURL}}">{{.User.Login}}</a></td>
				</tr>
				{{end}}
			</tbody>
		</table>
		`))
	result, err := github.SearchIssues(searchFilter)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err = list.Execute(w, result)
		if err != nil {
			log.Fatal(err)
		}
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
