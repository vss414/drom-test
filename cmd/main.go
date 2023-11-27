package main

import (
	"bufio"
	"github.com/vss414/drom-test/internal/parser"
	"log"
	"os"
	"strings"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	in := bufio.NewReader(os.Stdin)

	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for {
				var s string
				var err error
				if s, err = in.ReadString('\n'); err != nil {
					return
				}
				s = strings.TrimSpace(s)

				log.Printf("Start parsing report for %s", s)

				if car, err := parser.Parse(s); err == nil {
					if err := parser.Save(*car); err != nil {
						log.Printf("failed to find %s: %s", s, err)
					}
				} else {
					log.Printf("failed to find %s: %s", s, err)
				}

				log.Printf("Finish parsing report for %s", s)
			}
		}()
	}

	wg.Wait()
}
