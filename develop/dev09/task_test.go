package main

import (
	"sync"
	"testing"
)

func TestParseLinks(t *testing.T) {
	wg := &Wget{
		domain:  "https://test-domain.ru/",
		deep:    6,
		curDeep: 0,
		mu:      &sync.RWMutex{},
	}
	page := `<a href="https://test-domain.ru/home/"></a>
				   <a href="https://test-domain.ru/index"></a>
				   <a href="https://test-domain.ru/home/123"></a>
				   <a href="ftp://test-domain.ru/home/123"></a>
				   <a href="https://test-another-domain.ru/vk/"></a>
				   <a href="/another"></a>
				   <a href="../app/main.css"></a>
				   <a href="another"></a>`

	expectedRes := []string{
		"https://test-domain.ru/home/",
		"https://test-domain.ru/index",
		"https://test-domain.ru/home/123",
		"https://test-domain.ru/another",
		"https://test-domain.ru/app/main.css",
		"https://test-domain.ru/another",
	}
	res := wg.parseLinks(page)
	for i := 0; i != len(res); i++ {
		if expectedRes[i] != res[i] {
			t.Errorf("Error parseLinks expected: %s, but got: %s", expectedRes[i], res[i])
		}
	}
	if wg.curDeep != 1 {
		t.Errorf("Error parseLinks increment deep value expected: %v, but got: %v", 1, wg.curDeep)
	}
}

func TestGetSavePathAndName(t *testing.T) {
	wg := &Wget{
		domain:   "https://test-domain.ru/",
		savePath: "testDir",
	}
	tests := []struct {
		url  string
		dir  string
		file string
	}{
		{
			url:  "https://test-domain.ru/",
			file: "index.html",
			dir:  "",
		},
		{
			url:  "https://test-domain.ru/home/",
			file: "home.html",
			dir:  "",
		},
		{
			url:  "https://test-domain.ru/home/123",
			file: "123.html",
			dir:  "home/",
		},
		{
			url:  "https://test-domain.ru/home/dir/file",
			file: "file.html",
			dir:  "home/dir/",
		},
	}
	for _, tt := range tests {
		dir, file := wg.getSavePathAndName(tt.url)
		if dir != tt.dir && file != tt.file {
			t.Errorf("Error getSavePathAndName expected: %s - %s, but got: %s,  %s", tt.dir, tt.file, dir, file)
		}
	}
}

func TestIsJunkPage(t *testing.T) {
	wg := &Wget{}
	tests := map[string]bool{
		"https://test-domain.ru/":                         false,
		"https://test-domain.ru/home?page=1/":             true,
		"https://test-domain.ru/=":                        true,
		"https://test-domain.ru/home/dir/file&date=23-01": true,
	}
	for k, v := range tests {
		if wg.isJunkPage(k) != v {
			t.Errorf("Error isJunkPage expected: %v, but got: %v", v, wg.isJunkPage(k))
		}
	}
}
