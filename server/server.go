package server

import (
	"crypto/tls"
	"fmt"
	"github.com/teris-io/depot/config"
	"github.com/teris-io/log"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"github.com/teris-io/bus"
	"errors"
)

func Start(args []string, options map[string]string) int {
	conf, err := config.Parse(nil)
	if err != nil {
		log.Error(err).Log("failed to load config")
		return 1
	}

	log.Level(log.InfoLevel).Log("starting depot server")

	addr := fmt.Sprintf("%s:%d", conf.Hostname, conf.Port)
	if err := http.ListenAndServe(addr, &handler{dataDir: conf.DataDir, remotes: conf.Remotes}); err != nil {
		log.Error(err).Log("server failed")
		return 2
	}
	return 0
}

type handler struct {
	dataDir string
	remotes []config.Remote
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var err error
	for _, rmt := range h.remotes {
		u := &unit{name: rmt.Name, url: rmt.Url, headers: rmt.Headers, dataDir: h.dataDir}
		if err = u.handle(w, r); err == nil {
			w.WriteHeader(200)
			return
		}
	}
	w.WriteHeader(404)
	w.Write([]byte(err.Error()))
}

type unit struct {
	name    string
	dataDir string
	url     string
	headers map[string]string
}

func (u *unit) handle(w http.ResponseWriter, r *http.Request) (err error) {
	artifact := path.Join(u.dataDir, u.name, r.URL.Path)
	if r.Method == "HEAD" {
		if _, err := os.Stat(artifact); os.IsExist(err) {
			err = nil
		} else {
			err = u.ping(r.URL.Path)
		}
		return
	}
	if _, err := os.Stat(artifact); os.IsNotExist(err) {
		if err = u.download(r.URL.Path, artifact); err != nil {
			return
		}
	}
	return u.deliver(w, artifact)
}

func (u *unit) ping(resourceUrl string) (err error) {
	var resp *http.Response
	if resp, err = u.client().Head(u.url + resourceUrl); err == nil {
		if resp.StatusCode >= 400 {
			err = errors.New("not found")
		}
	}
	return
}

func (u *unit) download(resourceUrl, artifact string) (err error) {
	var resp *http.Response
	if resp, err = u.client().Get(u.url + resourceUrl); err == nil {
		if resp.StatusCode < 400 {
			if err = os.MkdirAll(filepath.Dir(artifact), os.ModePerm); err == nil {
				var out *os.File
				if out, err = os.Create(artifact); err == nil {
					if _, err = io.Copy(out, resp.Body); err == nil {
						log.Level(log.InfoLevel).Logf("Downloaded %s", artifact)
					}
					out.Close()
					resp.Body.Close()
				}
			}
		} else {
			err = errors.New("not found")
		}
	}
	return
}

func (u *unit) deliver(w http.ResponseWriter, artifact string) (err error) {
	if _, err = os.Stat(artifact); err == nil {
		var in *os.File
		if in, err = os.Open(artifact); err == nil {
			defer in.Close()
			_, err = io.Copy(w, in)
		}
	}
	return
}

func (u *unit) client() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
}

/*
09 00:36:16.963882 INF &{GET /com/google/code/findbugs/annotations/3.0.1/annotations-3.0.1.pom HTTP/1.1 1 1 map[User-Agent:[Gradle/4.3 (Mac OS X;10.13.2;x86_64) (Oracle Corporation;1.8.0_152;25.152-b16)] Accept-Encoding:[gzip,deflate] Connection:[Keep-Alive]] {} <nil> 0 [] false localhost:9595 map[] map[] <nil> map[] 127.0.0.1:62209 /com/google/code/findbugs/annotations/3.0.1/annotations-3.0.1.pom <nil> <nil> <nil> 0xc42010e100}
09 00:36:53.497817 INF &{GET /com/google/code/findbugs/annotations/maven-metadata.xml HTTP/1.1 1 1 map[Connection:[Keep-Alive] User-Agent:[Gradle/4.3 (Mac OS X;10.13.2;x86_64) (Oracle Corporation;1.8.0_152;25.152-b16)] Accept-Encoding:[gzip,deflate]] {} <nil> 0 [] false localhost:9595 map[] map[] <nil> map[] 127.0.0.1:62212 /com/google/code/findbugs/annotations/maven-metadata.xml <nil> <nil> <nil> 0xc42010e280}

*/
