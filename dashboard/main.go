package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
)

var dashboardHTML = `
<html>
	<h2>Welcome to the SSO Dashboard</h2>
	<a href="http://localhost:8082/?token={{.Token}}">Go to LeetCode Clone</a><br/>
	<a href="http://localhost:8083/?token={{.Token}}">Go to GFG Clone</a>
</html>`

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/dashboard", handleDashboard)

	fmt.Println("Dashboard running on http://localhost:8081")
	http.ListenAndServe(":8081", r)
}

func handleDashboard(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	tmpl := template.Must(template.New("dashboard").Parse(dashboardHTML))
	tmpl.Execute(w, struct{ Token string }{Token: token})
}
