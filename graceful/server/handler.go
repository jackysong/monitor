package server

import (
	"fmt"
	"mime"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/kisekivul/graceful/archiver"
)

type Service struct {
	Host   string `help:"Host interface"`
	Port   int    `help:"Listening port"`
	Config `type:"embedded"`
}

//handler configuration
type Config struct {
	// Auth       string `help:"Enable HTTP basic auth with the chosen username and password (must be in the form 'user:pass')"`
	Directory  string `type:"arg" help:"[directory] from which files will be served"`
	Slashing   bool   `help:"Disable automatic slash insertion when loading an index.html or directory"`
	Listable   bool   `help:"Disable directory listing"`
	Archivable bool   `help:"Disable directory archiving (download directories by appending .zip .tar .tar.gz - archives are streamed without buffering)"`
	PushState  bool   `help:"Enable PushState mode, causes missing directory paths to return the root index.html file, instead of a 404. Allows for sane usage of the HTML5 History API." short:"s"`
	Fallback   string `help:"Requests that yeild a 404, will instead proxy through to the provided path (swaps in the appropriate Host header)"`
}

//file service handler
type Handler struct {
	config Config
	served sync.Map
	proxy  *httputil.ReverseProxy
	host   string
	root   string
}

//create new Server
func NewHandler(c Config) (http.Handler, error) {
	h := &Handler{
		config: c,
		served: sync.Map{},
	}
	_, err := os.Stat(c.Directory)
	if c.Directory == "" || err != nil {
		return nil, fmt.Errorf("Missing directory: %s", c.Directory)
	}

	if c.PushState {
		h.root = filepath.Join(c.Directory, "index.html")
		if _, err := os.Stat(h.root); err != nil {
			return nil, fmt.Errorf("'%s' is required for pushstate", h.root)
		}
	}

	if c.Fallback != "" {
		u, err := url.Parse(c.Fallback)
		if err != nil {
			return nil, err
		}
		if !strings.HasPrefix(u.Scheme, "http") {
			return nil, fmt.Errorf("Invalid fallback protocol scheme")
		}
		h.host = u.Host
		h.proxy = httputil.NewSingleHostReverseProxy(u)
	}

	handler := http.Handler(h)

	//listen
	return handler, nil
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//requested target
	p := filepath.Join(h.config.Directory, r.URL.Path)
	//check file or dir
	var isDir, missing bool
	if info, err := os.Stat(p); err != nil {
		missing = true
	} else {
		isDir = info.IsDir()
	}

	if h.config.PushState && missing && filepath.Ext(p) == "" {
		//missing and pushstate and no ext
		p = h.root //change to request for the root
		isDir = false
		missing = false
	}

	if h.proxy != nil && (missing || isDir) {
		//fallback proxy enabled
		r.Host = h.host
		h.proxy.ServeHTTP(w, r)
		return
	}

	if h.config.Archivable && missing {
		//archive dir
		archive(w, p)
		return
	}

	if !isDir && missing {
		reply(w, 404, "Not found")
		return
	}

	//force trailing slash
	if isDir && !h.config.Slashing && !strings.HasSuffix(r.URL.Path, "/") {
		w.Header().Set("Location", r.URL.Path+"/")
		reply(w, 302, "Redirecting (must use slash for directories)")
		return
	}

	//directory list
	if isDir {
		if !h.config.Listable {
			reply(w, 403, "Listing not allowed")
			return
		}
		h.dirlist(w, r, p)
		return
	}

	//check file again
	info, err := os.Stat(p)
	if err != nil {
		reply(w, 404, "Not found")
		return
	}

	//stream file
	f, err := os.Open(p)
	if err != nil {
		reply(w, 500, err.Error())
		return
	}

	modtime := info.ModTime()
	//first time - dont use cache
	if served, ok := h.served.Load(p); !ok || !served.(bool) {
		h.served.Store(p, true)
		modtime = time.Now()
	}

	//http.ServeContent handles caching and range requests
	http.ServeContent(w, r, info.Name(), modtime, f)
}

func reply(w http.ResponseWriter, status int, msg string) {
	w.WriteHeader(status)
	if msg != "" {
		w.Write([]byte(msg))
	}
}

func archive(w http.ResponseWriter, path string) {
	var (
		err       error
		dir       string
		available bool
	)

	//check archivable
	ext := archiver.Extension(path)
	if ext != "" {
		var err error
		if dir, err = filepath.Abs(strings.TrimSuffix(path, ext)); err == nil {
			if info, err := os.Stat(dir); err == nil && info.IsDir() {
				available = true
			}
		}
	}

	if available {
		w.Header().Set("Content-Type", mime.TypeByExtension(ext))
		w.Header().Set("Content-Disposition", "attachment; filename="+filepath.Base(dir)+ext)
		w.WriteHeader(200)
		//archiver
		a, _ := archiver.NewWriter(w, ext, true)
		if err = a.AddDir(dir); err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		if err = a.Close(); err != nil {
			w.Write([]byte(err.Error()))
			return
		}
	}
}
