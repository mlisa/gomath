package pegmatch_test

import (
	"bytes"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/mlisa/gomath/parser"
)

func TestPegmatchSimple(t *testing.T) {
	operation := "4+2"
	start := time.Now()
	got, err := parser.ParseReader("", bytes.NewBufferString(operation))
	if err != nil || got.(int) != 6 {
		t.Error("failed")
	}
	log.Printf("Write %s", time.Since(start))
	fmt.Println(got)
}

func TestPegmatchMedium(t *testing.T) {
	operation := "3*(4+2)"
	start := time.Now()
	got, err := parser.ParseReader("", bytes.NewBufferString(operation))
	if err != nil || got.(int) != 18 {
		t.Error("failed")
	}
	log.Printf("Write %s", time.Since(start))
	fmt.Println(got)
}

func TestPegmatchHard(t *testing.T) {
	operation := "3*(4+2)/(2+4*4)+1"
	start := time.Now()
	got, err := parser.ParseReader("", bytes.NewBufferString(operation))
	if err != nil || got.(int) != 2 {
		t.Error("failed")
	}
	log.Printf("Write %s", time.Since(start))
	fmt.Println(got)
}
