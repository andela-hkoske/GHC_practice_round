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
