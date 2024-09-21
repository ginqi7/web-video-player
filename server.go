package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"io/fs"
	"log/slog"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
)

//go:embed templates static
var embedFS embed.FS

type Server struct {
	tmpl          *template.Template
	staticVersion string
}

func httpError(r *http.Request, w http.ResponseWriter, err error, code int) {
	http.Error(w, err.Error(), code)
	slog.Error("failed request",
		err,
		slog.String("url", r.URL.String()),
		slog.Int("code", code),
	)
}

// NormalizePath normalizes the request URL by removing the delimeter suffix.
func NormalizePath(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = strings.TrimRight(r.URL.Path, "/")
		h(w, r)
	}
}

// DisableFileListing disables file listing under directories. It can be used with the built-in http.FileServer.
func DisableFileListing(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") {
			http.NotFound(w, r)
			return
		}
		h.ServeHTTP(w, r)
	})
}

// Don't include sprig just for one function.
var templateFunctions = map[string]any{
	"defaultString": func(s string, def string) string {
		if s == "" {
			return def
		}
		return s
	},
}

type TemplateData struct {
	StaticVersion string
	Navigations   []Navigation
	Directories   []Directory
	Files         []File
	VideoURL      string
}

func (s *Server) PlayerHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	tmplData := TemplateData{
		StaticVersion: s.staticVersion,
		VideoURL:      fmt.Sprintf("/video/%s", query.Get("path")),
	}
	if err := s.tmpl.ExecuteTemplate(w, "player.html", tmplData); err != nil {
		httpError(r, w, err, http.StatusInternalServerError)
		return
	}
}

func (s *Server) VideoHandler(w http.ResponseWriter, r *http.Request) {
	// Open Audio File
	videoFile, err := openFile(r.URL.Path)
	if err != nil {
		http.Error(w, "Unable to open audio file", http.StatusInternalServerError)
		return
	}
	defer videoFile.Close()

	// Get File Info
	fileInfo, err := videoFile.Stat()
	if err != nil {
		http.Error(w, "Unable to get file info", http.StatusInternalServerError)
		return
	}
	rangeHeader := r.Header.Get("Range")

	if rangeHeader != "" {
		var start, end int64
		_, err := fmt.Sscanf(rangeHeader, "bytes=%d-%d", &start, &end)
		if err != nil {
			_, err = fmt.Sscanf(rangeHeader, "bytes=%d-", &start)
			end = 0
			if err != nil {
				http.Error(w, "Invalid Range", http.StatusRequestedRangeNotSatisfiable)
				return
			}
		}
		fileSize := fileInfo.Size()
		if start >= fileSize || end >= fileSize {
			http.Error(w, "Requested range not satisfiable", http.StatusRequestedRangeNotSatisfiable)
			return
		}

		if end < 0 || end == 0 {
			end = fileSize - 1
		}

		chunkSize := end - start + 1

		// Set Content-Range
		w.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, fileSize))
		w.Header().Set("Content-Length", strconv.FormatInt(chunkSize, 10))
		w.WriteHeader(http.StatusPartialContent)
		videoFile.Seek(start, 0)
		buffer := make([]byte, chunkSize)
		videoFile.Read(buffer)
		w.Write(buffer)
	} else {
		// Set Content-Length
		w.Header().Set("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))
		w.WriteHeader(http.StatusOK)
		http.ServeContent(w, r, fileInfo.Name(), fileInfo.ModTime(), videoFile)
	}

	// Set Content-Type
	w.Header().Set("Content-Type", "video/mp4")
	// Set Accept-Ranges, for Chrome setting `currentTime'
	// reference: https://segmentfault.com/q/1010000002908474
	w.Header().Set("Accept-Ranges", "bytes")
}

func (s *Server) ListingHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path)
	tmplData := TemplateData{
		StaticVersion: s.staticVersion,
		Navigations:   parseNavigation(r.URL.Path),
		Directories:   getDirectories(r.URL.Path),
		Files:         getFiles(r.URL.Path),
	}
	if err := s.tmpl.ExecuteTemplate(w, "listing.html", tmplData); err != nil {
		httpError(r, w, err, http.StatusInternalServerError)
		return
	}
}

func (s *Server) TranslateHandler(w http.ResponseWriter, r *http.Request) {
	urlQuery := r.URL.Query()
	jsonData, err := json.Marshal(query(urlQuery.Get("word")))
	if err != nil {
		fmt.Println(err)
	}
	w.Write(jsonData)
}

// StartServer starts HTTP server.
func StartServer(addr string) error {
	tmpl, err := template.New("").Funcs(templateFunctions).ParseFS(embedFS, "templates/*.html")

	if err != nil {
		return err
	}

	mux := http.NewServeMux()

	staticVersion := fmt.Sprintf("%x", rand.Uint64())
	staticFS, err := fs.Sub(embedFS, "static")
	if err != nil {
		return err
	}
	s := Server{
		tmpl:          tmpl,
		staticVersion: staticVersion,
	}
	staticPath := fmt.Sprintf("/static/%s/", staticVersion)
	mux.Handle(staticPath, DisableFileListing(http.StripPrefix(staticPath, http.FileServer(http.FS(staticFS)))))
	// mux.Handle("/", http.StripPrefix("/", (NormalizePath(s.PlayerHandler))))
	mux.Handle("/translate", http.StripPrefix("/translate", (NormalizePath(s.TranslateHandler))))
	// mux.Handle("/listing", http.StripPrefix("/listing", (NormalizePath(s.ListingHandler))))
	mux.Handle("/listing/", http.StripPrefix("/listing", (NormalizePath(s.ListingHandler))))
	mux.Handle("/play/", http.StripPrefix("/play", (NormalizePath(s.PlayerHandler))))
	mux.Handle("/video/", http.StripPrefix("/video/", (NormalizePath(s.VideoHandler))))
	mux.Handle("/", http.RedirectHandler("/listing/", http.StatusMovedPermanently))
	return http.ListenAndServe(addr, mux)
}
