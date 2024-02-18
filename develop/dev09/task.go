package main

import (
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"
)

/*
=== Утилита wget ===

Реализовать утилиту wget с возможностью скачивать сайты целиком

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

type Wget struct {
	domain      string
	savePath    string
	transport   *http.Transport
	urlArr      map[string]bool
	mu          *sync.RWMutex
	wg          *sync.WaitGroup
	rateLimiter <-chan time.Time
	deep        int
	curDeep     int
}

func NewWget(domain string, savePath string, rps int, deep int) *Wget {
	return &Wget{
		domain:   domain,
		savePath: savePath,
		transport: &http.Transport{
			MaxIdleConns:       5,
			IdleConnTimeout:    25 * time.Second,
			DisableCompression: true,
		},
		urlArr:      make(map[string]bool),
		mu:          &sync.RWMutex{},
		wg:          &sync.WaitGroup{},
		rateLimiter: time.Tick(time.Second / time.Duration(rps)),
		deep:        deep,
		curDeep:     1,
	}
}

func (w *Wget) AddURL(url string) {
	w.mu.Lock()
	w.urlArr[url] = false
	w.mu.Unlock()
}

func (w *Wget) SetInspected(url string) {
	w.mu.Lock()
	w.urlArr[url] = true
	w.mu.Unlock()
}

func (w *Wget) Run() {
	os.Mkdir(w.savePath, os.ModePerm)
	w.AddURL(w.domain)
	w.GetSite()
}

func (w *Wget) SavePage(url string, page []byte) (bool, error) {
	if w.IsJunkPage(url) {
		return true, nil
	}
	dirPath, fileName := w.GetSavePathAndName(url)
	if dirPath != "" {
		if err := os.MkdirAll(w.savePath+"/"+dirPath, os.ModePerm); err != nil {
			return false, err
		}
	}
	pathToFile := filepath.Join(w.savePath+"/"+dirPath, fileName)
	file, err := os.Create(pathToFile)
	if err != nil {
		return false, err
	}
	defer file.Close()

	_, err = file.Write(page)
	if err != nil {
		return false, err
	}
	return false, nil
}

func (w *Wget) IsJunkPage(url string) bool {
	if strings.Contains(url, "?") {
		return true
	}
	if strings.Contains(url, "=") {
		return true
	}
	if strings.Contains(url, "&") {
		return true
	}
	return false
}

func (w *Wget) GetURL(url string) ([]byte, error) {
	client := http.Client{
		Transport: w.transport,
		Timeout:   time.Second * 30,
	}
	response, err := client.Get(url)

	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	buff, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return buff, nil
}

func (w *Wget) ParseLinks(page string) []string {
	if w.deep == w.curDeep {
		return []string{}
	}
	w.mu.Lock()
	w.curDeep++
	w.mu.Unlock()
	invalidContains := []string{"https://", "http://", "ftp://", "www.", "mailto:"}

	r := regexp.MustCompile(`href="(.*?)"`)
	links := r.FindAllString(page, -1)
	var clearValidLinks []string
	var cleanLink string
	for _, v := range links {
		cleanLink = v[6 : len(v)-1]

		if strings.Contains(cleanLink, w.domain) {
			clearValidLinks = append(clearValidLinks, cleanLink)
			continue
		}
		valid := true
		for _, invalid := range invalidContains {
			if strings.Contains(cleanLink, invalid) {
				valid = false
				break
			}
		}
		if valid {
			cleanLink = strings.Replace(cleanLink, "../", "", -1)
			cleanLink = strings.Replace(cleanLink, "./", "", -1)
			cleanLink = strings.TrimLeft(cleanLink, "/")
			clearValidLinks = append(clearValidLinks, w.domain+"/"+cleanLink)
		}
	}
	return clearValidLinks
}

func (w *Wget) GetSavePathAndName(url string) (dir string, file string) {
	var dirtyPath string
	if strings.Contains(url, w.domain) {
		dirtyPath = strings.TrimPrefix(url, w.domain)
	}
	if strings.HasSuffix(dirtyPath, "/") {
		dirtyPath = strings.TrimSuffix(dirtyPath, "/")
	}
	dir, file = filepath.Split(dirtyPath)

	if file == "" {
		file = "index"
	}
	if filepath.Ext(file) == "" {
		file += ".html"
	}
	return
}

func (w *Wget) GetSite() {
	var pageLen int
	for {
		pageLen = len(w.urlArr)
		w.mu.RLock()
		for url, inspected := range w.urlArr {
			if !inspected {
				w.wg.Add(1)
				go w.ProcessPage(url)
			}
		}
		w.mu.RUnlock()
		w.wg.Wait()
		if pageLen == len(w.urlArr) {
			break
		}
	}
	w.transport.CloseIdleConnections()
}

func (w *Wget) ProcessPage(url string) {
	defer w.wg.Done()
	<-w.rateLimiter
	log.Println("Делаю запрос к ", url)
	page, err := w.GetURL(url)
	w.SetInspected(url)
	if err != nil {
		log.Printf("(Get URL): не удалось загрузить страницу %s error: %s\n", url, err.Error())
	}
	isJunk, err := w.SavePage(url, page)
	if err != nil {
		log.Printf("(Save page): не удалось сохранить страницу %s error: %s\n", url, err.Error())
	}
	if !isJunk {
		links := w.ParseLinks(string(page))
		for _, link := range links {
			w.AddURL(link)
		}
	}
}

func main() {
	var (
		domain   string
		savePath string
		rps      int
		deep     int
	)

	flag.StringVar(&domain, "url", "", "базовый url сайта в формате https://site.su")
	flag.StringVar(&savePath, "path", "site", "папка для сохранения")
	flag.IntVar(&deep, "deep", 1, "папка для сохранения")
	flag.IntVar(&rps, "rps", 10, "папка для сохранения")
	flag.Parse()
	wget := NewWget(domain, savePath, rps, deep)
	wget.Run()
}
