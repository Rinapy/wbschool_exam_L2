package cut

import (
	"database/sql"
	"flag"
	"log"
	"os"
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

var (
	f string
	d string
	s bool
)

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


func ParseFlag() (*Config, []line) {
	cfg := NewCfg()

	f := flag.String("f", "1", "индексы или интервал столбцов которые будут выведены")
	flag.StringVar(&cfg.d, "d", ":", "разделитель который используется для разделения столбцов")
	flag.StringVar(&cfg.ld, "ld", "\t", "разделитель который используется для разделения строк")
	flag.BoolVar(&cfg.s, "s", true, "выводит строки в котором есть хотя бы один разделитель")
	flag.Parse()
	err := cfg.parseF(*f)
	if err != nil {
		log.Fatal(err)
	}
	lineSlice, err := cfg.ParseData()
	if err != nil {
		log.Fatal(err)
	}
	return cfg, lineSlice
}

func (c *Config) ParseData() ([]line, error) {
	strLines := os.Args[1]
	if strLines == "" {
		return []line{}, &DataNotFound{}
	}
	lines := strings.Split(strLines, c.ld)
	data := make(lineSlice, len(lines))
	for i := 0; i != len(lines); i++ {
		data[i].text = lines[i]
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
	return nil
}

func (a *CutApp) Run()  {
	outData := make(lineSlice, 0)
	for lineIndex, l := range a.Line{
		rows := strings.Split(l.text, a.Cfg.d)
		for _, rowIndex := range a.Cfg.f{if lineIndex == rowIndex{
				outData[lineIndex] = rows[rowIndex]
		}
	}
}