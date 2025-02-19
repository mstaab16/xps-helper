package main

import (
	"embed"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

//go:embed transformed_data.csv
var dataCSV embed.FS

func tableComponentBuilder(data []DataRow, queryParams url.Values) (templ.Component, error) {
	var filtered_data []DataRow
	search := queryParams.Get("search")
	width_str := queryParams.Get("width")
	if isAlpha(search) {
		filtered_data = filterByElement(data, search)
	} else {
		energy, err := strconv.ParseFloat(strings.TrimSpace(search), 64)
		width, err := strconv.ParseFloat(strings.TrimSpace(width_str), 64)
		if err != nil {
			return dataTable([]DataRow{}), nil
		}
		filtered_data = filterByEnergy(data, energy-width, energy+width)
	}
	// err := hello(name, age).Render(context.Background(), w)
	// Render the templ component
	component := dataTable(filtered_data)
	return component, nil
}

func xpsTableHandler(data []DataRow) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		queryParams := r.URL.Query()
		component, err := tableComponentBuilder(data, queryParams)
		if err != nil {
			http.Error(w, "Failed to render template 1", http.StatusInternalServerError)
			return
		}
		err = component.Render(r.Context(), w)
		if err != nil {
			http.Error(w, "Failed to render template 2", http.StatusInternalServerError)
		}
	}
}

func main() {
	data, err := loadCSV()
	if err != nil {
		fmt.Println(err)
		return
	}

	// http.Handle("/", templ.Handler(component))

	fmt.Println("Listening on :3000")
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Post("/", xpsTableHandler(data).ServeHTTP)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		// if r has header HX-Request, then render the table component only
		queryParams := r.URL.Query()
		tableComponent, err := tableComponentBuilder(data, queryParams)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Failed to render template 3", http.StatusInternalServerError)
			return
		}
		if r.Header.Get("HX-Request") == "true" {
			tableComponent.Render(r.Context(), w)
		} else {
			index(tableComponent, r).Render(r.Context(), w)
		}
	})
	http.ListenAndServe(":3000", r)
}
