package cut

import (
	"flag"
	"fmt"
	"log"
	"strconv"
	"strings"
)

/*
=== Утилита cut ===

Принимает STDIN, разбивает по разделителю (TAB) на колонки, выводит запрошенные

Поддержать флаги:
-f - "fields" - выбрать поля (колонки)
-d - "delimiter" - использовать другой разделитель
-s - "separated" - только строки с разделителем

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

type Config struct {
	f  []int
	d  string
	ld string
	s  bool
}

func NewCfg() *Config {
	return &Config{}
}

type CutApp struct {
	Cfg  *Config
	Line lineSlice
}

type line struct {
	text string
}

type lineSlice []line

func NewApp() (*CutApp, error) {
	cfg := NewCfg()
	var f string
	flag.StringVar(&f, "f", "1", "индексы или интервал столбцов которые будут выведены")
	flag.StringVar(&cfg.d, "d", ":", "разделитель который используется для разделения столбцов")
	flag.StringVar(&cfg.ld, "ld", " ", "разделитель который используется для разделения строк")
	flag.BoolVar(&cfg.s, "s", false, "выводит строки в котором есть хотя бы один разделитель")
	flag.Parse()
	err := cfg.parseF(f)
	if err != nil {
		log.Fatal(err)
	}
	strLines := flag.Arg(0)
	ls, err := cfg.parseData(strLines)
	if err != nil {
		log.Fatal(err)
	}

	return &CutApp{
		Cfg:  cfg,
		Line: ls,
	}, nil
}

func (c *Config) parseData(strLines string) (lineSlice, error) {
	if strLines == "" {
		return nil, &DataNotFound{}
	}
	lines := strings.Split(strLines, c.ld)
	data := make(lineSlice, len(lines))
	for i := 0; i != len(lines); i++ {
		if c.s {
			if strings.Contains(lines[i], c.d) {
				data[i] = line{lines[i]}
			}
		} else {
			data[i] = line{lines[i]}
		}
	}
	return data, nil
}

func (c *Config) parseF(f string) error {
	if strings.Contains(f, "-") {
		strRows := strings.Split(f, "-")
		startIndex, err := strconv.Atoi(strRows[0])
		if err != nil {
			return &IndexValueError{}
		}
		endIndex, err := strconv.Atoi(strRows[1])
		if err != nil {
			return &IndexValueError{}
		}
		if startIndex > endIndex || startIndex <= 0 {
			return &IndexValueError{}
		}
		c.f = make([]int, endIndex-startIndex+1)
		for i := range c.f {
			c.f[i] = startIndex + i - 1
		}
		return nil
	} else if strings.Contains(f, ",") {
		strRows := strings.Split(f, ",")
		c.f = make([]int, len(strRows))
		for i, v := range strRows {
			val, err := strconv.Atoi(v)
			if err != nil {
				return &IndexValueError{}
			}
			c.f[i] = val - 1
		}
		return nil
	}
	index, err := strconv.Atoi(f)
	if err != nil {
		return &IndexValueError{}
	}
	if index <= 0 {
		return &IndexValueError{}
	}
	c.f = make([]int, 1)
	c.f[0] = index - 1
	return nil
}

func (a *CutApp) Run() {
	outData := make(lineSlice, 0)
	for _, l := range a.Line {
		rows := strings.Split(l.text, a.Cfg.d)
		outLine := make([]string, 0)
		for _, needIndex := range a.Cfg.f {
			if needIndex < len(rows) {
				outLine = append(outLine, rows[needIndex])
			}
		}
		outData = append(outData, line{text: strings.Join(outLine, a.Cfg.d)})
	}
	for _, v := range outData {
		fmt.Println(v.text)
	}
}
