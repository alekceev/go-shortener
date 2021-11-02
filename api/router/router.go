package router

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"

	"github.com/alekceev/go-shortener/api/handler"
	"github.com/alekceev/go-shortener/api/openapi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type Router struct {
	*chi.Mux
	hs *handler.Handlers
}

func NewRouter(hs *handler.Handlers) *Router {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	ret := &Router{
		hs: hs,
	}

	r.Get("/", ret.GetMainPage)
	r.Get("/openapi", ret.GetOpenAPI)
	r.Get("/static/{filename}", func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RequestURI)
		http.StripPrefix("/static", http.FileServer(http.Dir("./web/static"))).ServeHTTP(w, r)
	})

	r.Mount("/", openapi.Handler(ret))
	swg, err := openapi.GetSwagger()
	if err != nil {
		log.Fatal("swagger fail")
	}

	r.Get("/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		enc := json.NewEncoder(w)
		_ = enc.Encode(swg)
	})

	ret.Mux = r
	return ret
}

type Url handler.URL

func (Url) Bind(r *http.Request) error {
	return nil
}

func (Url) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type ResponseURL struct {
	ShortURL string `json:"short_url"`
	StatsURL string `json:"stats_url"`
}

func (ResponseURL) Bind(r *http.Request) error {
	return nil
}

func (ResponseURL) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (rt *Router) CreateShortURL(w http.ResponseWriter, r *http.Request) {
	ru := Url{}
	if err := render.Bind(r, &ru); err != nil {
		_ = render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	u, err := rt.hs.Create(r.Context(), handler.URL(ru))
	if err != nil {
		_ = render.Render(w, r, ErrRender(err))
		return
	}

	// render.Render(w, r, Url(u))
	responseURL := &ResponseURL{
		ShortURL: "/" + u.ShortURL,
		StatsURL: "/" + u.ShortURL + "/stats",
	}

	_ = render.Render(w, r, responseURL)
}

func (rt *Router) RedirectURL(w http.ResponseWriter, r *http.Request, shortURL string) {
	url, err := rt.hs.GetURL(r.Context(), shortURL)
	if err != nil {
		log.Println(err)
		// _ = render.Render(w, r, ErrNotFound)
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	http.Redirect(w, r, url.URL, http.StatusSeeOther)
}

func (rt *Router) GetStats(w http.ResponseWriter, r *http.Request, shortURL string) {
	url, err := rt.hs.GetStats(r.Context(), shortURL)
	if err != nil {
		log.Println(err)
		_ = render.Render(w, r, ErrNotFound)
		// http.Error(w, "not found", http.StatusNotFound)
		return
	}

	_ = render.Render(w, r, Url(url))
}

func (rt *Router) GetMainPage(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("./web/templates/index.html")
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", 500)
		return
	}

	if err = ts.Execute(w, nil); err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", 500)
		return
	}
}

func (rt *Router) GetOpenAPI(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("./web/templates/openapi.html")
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", 500)
		return
	}

	if err = ts.Execute(w, nil); err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", 500)
		return
	}
}
