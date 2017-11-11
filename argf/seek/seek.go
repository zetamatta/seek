package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/mattn/go-colorable"
	"github.com/zetamatta/experimental/argf"
	"github.com/zetamatta/go-mbcs"
)

var ignoreCase = flag.Bool("i", false, "ignore case")

func main1() error {
	flag.Parse()
	args := flag.Args()

	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "Usage: %s [flags...] REGEXP Files...\n", os.Args[0])
		flag.PrintDefaults()
		return nil
	}
	var pattern string = args[0]
	if *ignoreCase {
		pattern = "(?i)" + pattern
	}

	rx, err := regexp.Compile(pattern)
	if err != nil {
		return err
	}
	out := colorable.NewColorableStdout()
	r := argf.NewFiles(args[1:])
	for r.Scan() {
		line := r.Bytes()

		var text string
		if utf8.Valid(line) {
			text = string(line)
		} else {
			var err error
			text, err = mbcs.AtoU(line)
			if err != nil {
				text = err.Error()
			}
		}
		text = strings.Replace(text, "\xEF\xBB\xBF", "", 1)

		m := rx.FindAllStringIndex(text, -1)
		if m != nil {
			fmt.Fprintf(out, "\x1B[35;1m%s:\x1B[32;1m%d\x1B[36;1m:\x1B[37;1m", r.Filename(), r.FNR())
			last := 0
			for i := 0; i < len(m); i++ {
				fmt.Fprintf(out, "%s\x1B[31;1m%s\x1B[37;1m",
					text[last:m[i][0]],
					text[m[i][0]:m[i][1]])
				last = m[i][1]
			}
			fmt.Fprintln(out, text[last:])
		}
	}
	return r.Err()
}

func main() {
	if err := main1(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	os.Exit(0)
}