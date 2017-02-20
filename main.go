package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

type (
	Pizza struct {
		Rows, Columns, MinIng, MaxCells int
		Mushrooms, Tomatoes             int
		Raw                             string
		Arrangement                     [][]string
	}

	C [2]int

	Slice struct {
		Start, Stop C
	}
)

func main() {
	var (
		p   = &Pizza{}
		err error
	)
	p.Raw, err = ReadInput("./example.in")
	if err != nil {
		goto ERROR
	}
	err = p.ParseRaw()
	if err != nil {
		goto ERROR
	}
	p.SetTomatoes()
	p.SetMushrooms()
	p.SetArrangement()
	log.Printf("%+v\n", p)
	log.Println(p.isSlice(Slice{C{0, 0}, C{1, 2}}))
ERROR:
	log.Println(err)
	return
}

func ReadInput(filename string) (string, error) {
	buf := bytes.NewBuffer(nil)
	f, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	io.Copy(buf, f)
	f.Close()
	s := string(buf.Bytes())
	if len(s) < 1 {
		return s, errors.New("File is empty")
	}
	return strings.TrimSpace(s), nil
}

func (p *Pizza) ParseRaw() error {
	var err error
	p.Rows, err = strconv.Atoi(string(p.Raw[0]))
	if err != nil {
		err = fmt.Errorf(err.Error(), "Failed to parse number of rows")
		return err
	}
	p.Columns, err = strconv.Atoi(string(p.Raw[2]))
	if err != nil {
		err = fmt.Errorf(err.Error(), "Failed to parse number of columns")
		return err
	}
	p.MinIng, err = strconv.Atoi(string(p.Raw[4]))
	if err != nil {
		err = fmt.Errorf(err.Error(), "Failed to parse minimum ingredients per cell")
		return err
	}
	p.MaxCells, err = strconv.Atoi(string(p.Raw[6]))
	if err != nil {
		err = fmt.Errorf(err.Error(), "Failed to parse max number of cells per slice")
		return err
	}
	p.Arrangement = make([][]string, p.Rows)
	return nil
}

func (p *Pizza) SetTomatoes() {
	p.Tomatoes = strings.Count(p.Raw, "T")
}

func (p *Pizza) SetMushrooms() {
	p.Mushrooms = strings.Count(p.Raw, "M")
}

func (p *Pizza) SetArrangement() {
	ing := p.Raw[8:]
	tempRows := strings.Split(ing, "\n")
	for pos, row := range tempRows {
		p.Arrangement[pos] = strings.Split(row, "")
	}
}

func (p *Pizza) GetFirstSlices() Slice {
	var tempSlice Slice
	col, colVal := 0, ""
	for row, rowVal := range p.Arrangement {
		for col, colVal = range rowVal {
			tempSlice = Slice{C{row, 0}, C{row, col}}
			p.isSlice(tempSlice)
		}
	}
}

func (p *Pizza) isSlice(sl Slice) bool {
	var vals string
	for i, lastRow := 0, sl.Stop[0]; i <= lastRow; i++ {
		// log.Println(p.Arrangement[i][:sl.Stop[1]+1])
		vals = vals + strings.Join(p.Arrangement[i][:sl.Stop[1]+1], "")
	}
	// log.Println(vals)
	if strings.Count(vals, "T") >= p.MinIng && strings.Count(vals, "M") >= p.MinIng {
		return true
	}
	return false
}
