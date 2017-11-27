package pegmatch_test

import (
	"bytes"
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
	log.Printf("STRING: %s\n", operation)
	log.Printf("RESULT: %d\n", got.(int))
	log.Printf("TIME %s\n", time.Since(start))
	log.Println("---------------------------------------------------------------------------------")
}

func TestPegmatchMedium(t *testing.T) {
	operation := "3*(4+2)"
	start := time.Now()
	got, err := parser.ParseReader("", bytes.NewBufferString(operation))
	if err != nil || got.(int) != 18 {
		t.Error("failed")
	}
	log.Printf("STRING: %s\n", operation)
	log.Printf("RESULT: %d\n", got.(int))
	log.Printf("Time %s\n", time.Since(start))
	log.Println("---------------------------------------------------------------------------------")
}

func TestPegmatchHard(t *testing.T) {
	operation := "3*(4+2)/(2+4*4)+1"
	start := time.Now()
	got, err := parser.ParseReader("", bytes.NewBufferString(operation))
	if err != nil || got.(int) != 2 {
		t.Error("failed")
	}
	log.Printf("STRING: %s\n", operation)
	log.Printf("RESULT: %d\n", got.(int))
	log.Printf("Time %s\n", time.Since(start))
	log.Println("---------------------------------------------------------------------------------")
}
