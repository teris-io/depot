package server

import (
	"fmt"
	"github.com/teris-io/depot/config"
	"github.com/teris-io/log"
	"net/http"
)

func Start(args []string, options map[string]string) int {
	conf, err := config.Parse(nil)
	if err != nil {
		log.Error(err).Log("failed to load config")
		return 1
	}

	log.Level(log.InfoLevel).Log("starting depot server")

	addr := fmt.Sprintf("%s:%d", conf.Hostname, conf.Port)
	if err := http.ListenAndServe(addr, &handler{remotes: conf.Remotes}); err != nil {
		log.Error(err).Log("server failed")
		return 2
	}
	return 0
}

type handler struct {
	remotes []config.Remote
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Level(log.InfoLevel).Logf("%v", r)
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

/*
09 00:36:16.963882 INF &{GET /com/google/code/findbugs/annotations/3.0.1/annotations-3.0.1.pom HTTP/1.1 1 1 map[User-Agent:[Gradle/4.3 (Mac OS X;10.13.2;x86_64) (Oracle Corporation;1.8.0_152;25.152-b16)] Accept-Encoding:[gzip,deflate] Connection:[Keep-Alive]] {} <nil> 0 [] false localhost:9595 map[] map[] <nil> map[] 127.0.0.1:62209 /com/google/code/findbugs/annotations/3.0.1/annotations-3.0.1.pom <nil> <nil> <nil> 0xc42010e100}
09 00:36:53.497817 INF &{GET /com/google/code/findbugs/annotations/maven-metadata.xml HTTP/1.1 1 1 map[Connection:[Keep-Alive] User-Agent:[Gradle/4.3 (Mac OS X;10.13.2;x86_64) (Oracle Corporation;1.8.0_152;25.152-b16)] Accept-Encoding:[gzip,deflate]] {} <nil> 0 [] false localhost:9595 map[] map[] <nil> map[] 127.0.0.1:62212 /com/google/code/findbugs/annotations/maven-metadata.xml <nil> <nil> <nil> 0xc42010e280}

 */
