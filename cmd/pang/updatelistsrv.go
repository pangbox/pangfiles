package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/pangbox/pangfiles/crypto/pyxtea"
	"github.com/pangbox/pangfiles/encoding/litexml"
	"github.com/pangbox/pangfiles/updatelist"
)

type cacheentry struct {
	modTime time.Time
	fSize   int64
	fInfo   updatelist.FileInfo
}

type server struct {
	key   pyxtea.Key
	dir   string
	cache map[string]cacheentry
	mutex sync.RWMutex
}

func (s *server) calcEntry(wg *sync.WaitGroup, entry *updatelist.FileInfo, f os.FileInfo) {
	defer wg.Done()
	var err error

	name := f.Name()
	*entry, err = updatelist.MakeFileInfo(s.dir, "", f, f.Size())

	if err != nil {
		log.Printf("Error calculating entry for %s: %s", name, err)
		entry.Filename = name
	} else {
		log.Printf("Successfully calculated entry for %s", name)

		s.mutex.Lock()
		defer s.mutex.Unlock()

		s.cache[name] = cacheentry{
			modTime: f.ModTime(),
			fSize:   f.Size(),
			fInfo:   *entry,
		}
	}
}

func (s *server) updateList(rw io.Writer) error {
	start := time.Now()

	files, err := ioutil.ReadDir(s.dir)
	if err != nil {
		return err
	}

	doc := updatelist.Document{}
	doc.Info.Version = "1.0"
	doc.Info.Encoding = "euc-kr"
	doc.Info.Standalone = "yes"
	doc.PatchVer = "FakeVer"
	doc.PatchNum = 9999
	doc.UpdateListVer = "20090331"

	hit, miss := 0, 0

	var wg sync.WaitGroup
	doc.UpdateFiles.Files = make([]updatelist.FileInfo, 0, len(files))
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		name := f.Name()

		s.mutex.RLock()
		cache, ok := s.cache[name]
		s.mutex.RUnlock()

		if ok && cache.modTime == f.ModTime() && cache.fSize == f.Size() {
			// Cache hit
			hit++
			doc.UpdateFiles.Files = append(doc.UpdateFiles.Files, cache.fInfo)
			doc.UpdateFiles.Count++
		} else {
			// Cache miss, calculate concurrently.
			miss++
			doc.UpdateFiles.Files = append(doc.UpdateFiles.Files, updatelist.FileInfo{})
			doc.UpdateFiles.Count++
			entry := &doc.UpdateFiles.Files[len(doc.UpdateFiles.Files)-1]
			wg.Add(1)
			go s.calcEntry(&wg, entry, f)
		}
	}

	wg.Wait()

	data, err := litexml.Marshal(doc)
	if err != nil {
		return err
	}

	if err := pyxtea.EncipherStreamPadNull(s.key, bytes.NewReader(data), rw); err != nil {
		return err
	}

	log.Printf("Updatelist served in %s (cache hits: %d, misses: %d)", time.Since(start), hit, miss)
	return nil
}

func (s *server) extracontents(w io.Writer) error {
	if _, err := w.Write([]byte(`<?xml version="1.0" standalone="yes" ?><extracontents><themes><pangya_default src="pangya_default.xml" url="http://127.0.0.1:8080/S4_Patch/extracontents/default/"/></themes></extracontents>`)); err != nil {
		return err
	}
	return nil
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.Body.Close()
	log.Printf("%s %s", r.Method, r.URL)
	if strings.Contains(strings.ToLower(r.URL.Path), "updatelist") {
		if err := s.updateList(w); err != nil {
			log.Printf("Error writing updateList: %v", err)
		}
	} else if strings.Contains(strings.ToLower(r.URL.Path), "extracontents") {
		if err := s.extracontents(w); err != nil {
			log.Printf("Error writing extracontents: %v", err)
		}
	}
}
