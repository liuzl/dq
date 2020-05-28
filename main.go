package main

import (
	"flag"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/golang/glog"
	"zliu.org/filestore"
	"zliu.org/goutil/rest"
	"zliu.org/q"
)

var (
	addr   = flag.String("addr", ":9080", "bind address")
	data   = flag.String("data", "./data", "queue dir")
	fs     = flag.Bool("fs", true, "filestore flag")
	queue  *q.Queue
	fstore *filestore.FileStore
)

func EnqueueHandler(w http.ResponseWriter, r *http.Request) {
	glog.Infof("addr=%s  method=%s host=%s uri=%s",
		r.RemoteAddr, r.Method, r.Host, r.RequestURI)
	r.ParseForm()
	data := strings.TrimSpace(r.FormValue("data"))
	if data == "" {
		rest.MustEncode(w, rest.RestMessage{"error", "data is empty"})
		return
	}
	if err := queue.Enqueue(data); err != nil {
		rest.MustEncode(w, rest.RestMessage{"error", err.Error()})
		return
	}
	if *fs {
		fstore.WriteLine([]byte(data))
	}
	rest.MustEncode(w, rest.RestMessage{"ok", nil})
}

func DequeueHandler(w http.ResponseWriter, r *http.Request) {
	glog.Infof("addr=%s  method=%s host=%s uri=%s",
		r.RemoteAddr, r.Method, r.Host, r.RequestURI)
	r.ParseForm()
	peek := strings.ToLower(strings.TrimSpace(r.FormValue("peek")))
	var ret string
	var err error
	if peek == "true" {
		ret, err = queue.Peek()
	} else {
		_, ret, err = queue.Dequeue(-1)
	}
	if err != nil {
		rest.MustEncode(w, rest.RestMessage{"error", err.Error()})
		return
	}
	w.Write([]byte(ret))
}

func StatusHandler(w http.ResponseWriter, r *http.Request) {
	glog.Infof("addr=%s  method=%s host=%s uri=%s",
		r.RemoteAddr, r.Method, r.Host, r.RequestURI)
	rest.MustEncode(w, queue.Status())
}

func main() {
	flag.Parse()
	defer glog.Flush()
	var err error
	qdir := filepath.Join(*data, "queue")
	if queue, err = q.NewQueue(qdir); err != nil {
		glog.Fatal(err)
	}
	if *fs {
		fsdir := filepath.Join(*data, "fs")
		if fstore, err = filestore.NewFileStore(fsdir); err != nil {
			glog.Fatal(err)
		}
	}
	defer glog.Info("server exit")
	http.Handle("/dequeue/", rest.WithLog(DequeueHandler))
	http.Handle("/enqueue/", rest.WithLog(EnqueueHandler))
	http.Handle("/status/", rest.WithLog(StatusHandler))
	glog.Info("dq listen on", *addr)
	glog.Error(http.ListenAndServe(*addr, nil))
}
