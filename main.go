package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
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
	Tree struct{
		Node Slice
		SubTrees []*Tree
	}
)

var (
	testPizza = &Pizza{}
	testErr error
)

func init() {
	args := os.Args
	if len(args) < 2 {
		panic(errors.New("You need to pass the full path of the input file"))
	}
	file := strings.TrimSpace(args[1])
	testPizza.Raw, testErr = ReadInput(file)
	if testErr != nil {
		fmt.Println(testErr)
		return
	}
	testErr = testPizza.ParseRaw()
	if testErr != nil {
		fmt.Println(testErr)
		return
	}
	testPizza.SetTomatoes()
	testPizza.SetMushrooms()
	testPizza.SetArrangement()
}

func main() {
	// fmt.Println(testPizza.GetViableSlices(C{1, 0}))
	fmt.Println(testPizza.traverse())
}

// ReadInput reads the input file
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

// ParseRaw takes the raw pizza input and parses it into meaningful information
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

// SetTomatoes counts the number of tomatoes in the pizza
// and sets the value under the Tomatoes property
func (p *Pizza) SetTomatoes() {
	p.Tomatoes = strings.Count(p.Raw, "T")
}

// SetMushrooms counts the number of mushrooms in the pizza
// and sets the value under the Mushroom property
func (p *Pizza) SetMushrooms() {
	p.Mushrooms = strings.Count(p.Raw, "M")
}

// SetArrangement parses the input file and maps the pizza into a two
// dimensional array set to its Arrangement property
func (p *Pizza) SetArrangement() {
	ing := p.Raw[8:]
	tempRows := strings.Split(ing, "\n")
	for pos, row := range tempRows {
		p.Arrangement[pos] = strings.Split(row, "")
	}
}

// GetViableSlices gets all viable slices given a starting point on the pizza point
func (p *Pizza) GetViableSlices(start C) []Slice {
	var (
		tempSlice   Slice
		validSlices []Slice
		col, row    int = start[1], start[0]
	)
	lenRows, rowLen, rowVal := 0, 0, []string{}
	for lenRows = len(p.Arrangement); row < lenRows; row++ {
		rowVal = p.Arrangement[row]
		rowLen = len(rowVal)
		if (((col - start[1]) + 1) * ((row - start[0]) + 1)) > p.MaxCells {
			break
		}
		for col = start[1]; col < rowLen; col++ {
			if (((col - start[1]) + 1) * ((row - start[0]) + 1)) > p.MaxCells {
				break
			}
			tempSlice = Slice{start, C{row, col}}
			if p.isSlice(tempSlice) {
				validSlices = append(validSlices, tempSlice)
			}
		}
		col = start[1]
	}
	return validSlices
}

// isSlices determines whether a slice is viable or not
func (p *Pizza) isSlice(sl Slice) bool {
	var vals string
	for i, lastRow := sl.Start[0], sl.Stop[0]; i <= lastRow; i++ {
		vals = vals + strings.Join(p.Arrangement[i][sl.Start[1]:sl.Stop[1]+1], "")
	}
	if strings.Count(vals, "T") >= p.MinIng && strings.Count(vals, "M") >= p.MinIng {
		return true
	}
	return false
}

// traverse returns a tree of all viable pizza slices
func (p *Pizza) traverse() *Tree{
	options := &Tree{}
	tempSlices := p.GetViableSlices(C{0, 0})
	for _, slice := range tempSlices{
		options.SubTrees = append(options.SubTrees, &Tree{Node: slice})
	}
	return options
}
