package main

import (
	"html/template"
	"net/http"
	"site/test/inmemory"
	"site/internal/handler"
)

func HomePage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles(inmemory.IndexHtml())
		if err != nil {
			handler.Error(w, r, http.StatusBadGateway, err)
			return
		}
		data := struct{}{}
		tmpl.Execute(w, data)
	}
}

func main() {
	http.HandleFunc("/", HomePage())
	http.ListenAndServe(":8081", nil)
}

// func (s *Server) RegisterPage() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		tmpl, err := template.ParseFiles(inmemory.GetInstance().GetRegisterHtml())
// 		if err != nil {
// 			ioutils.Error(w, r, http.StatusBadGateway, err)
// 			return
// 		}
// 		data := struct{}{}
// 		tmpl.Execute(w, data)
// 	}
// }

// func (s *Server) LoginPage() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		tmpl, err := template.ParseFiles(inmemory.GetInstance().GetLoginHtml())
// 		if err != nil {
// 			ioutils.Error(w, r, http.StatusBadGateway, err)
// 			return
// 		}
// 		data := struct{}{}
// 		tmpl.Execute(w, data)
// 	}
// }

// func (s *Server) RatingsPage() http.HandlerFunc {
// 	type RatingsViewData struct {
// 		Users   *api.UserList
// 		Country string
// 		City    string
// 	}
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		country := r.URL.Query().Get("country")
// 		if country == "" {
// 			country = "not configured"
// 		}
// 		city := r.URL.Query().Get("city")
// 		if city == "" {
// 			city = "not configured"
// 		}
// 		users, err := s.store.Users().All(s.ctx, 1)
// 		if err != nil {
// 			ioutils.Error(w, r, http.StatusInternalServerError, err)
// 			return
// 		}
// 		data := RatingsViewData{
// 			Users:   users,
// 			Country: country,
// 			City:    city,
// 		}
// 		tmpl, err := template.ParseFiles(inmemory.GetInstance().GetRatingsHtml())
// 		if err != nil {
// 			ioutils.Error(w, r, http.StatusBadGateway, err)
// 			return
// 		}
// 		tmpl.Execute(w, data)
// 	}
// }
