package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"unsafe"

	"github.com/maxzender/jv/colorwriter"
	"github.com/maxzender/jv/jsonfmt"
	"github.com/maxzender/jv/jsontree"
	"github.com/maxzender/jv/terminal"
	termbox "github.com/nsf/termbox-go"
	"github.com/tidwall/gjson"
)

var (
	colorMap = map[jsonfmt.TokenType]termbox.Attribute{
		jsonfmt.DelimiterType: termbox.ColorDefault,
		jsonfmt.BoolType:      termbox.ColorRed,
		jsonfmt.StringType:    termbox.ColorGreen,
		jsonfmt.NumberType:    termbox.ColorYellow,
		jsonfmt.NullType:      termbox.ColorMagenta,
		jsonfmt.KeyType:       termbox.ColorBlue,
	}
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [-o] <query> [file]\n", os.Args[0])
	flag.PrintDefaults()
}

func main() {
	var showHelp, output bool
	flag.BoolVar(&showHelp, "h", false, "print usage")
	flag.BoolVar(&output, "o", false, "pretty output")
	flag.BoolVar(&showHelp, "help", false, "print usage")

	flag.Usage = usage
	flag.Parse()
	if showHelp {
		flag.Usage()
		os.Exit(0)
	}

	reader := os.Stdin
	var err error
	query := ""
	args := flag.Args()
	if len(args) > 0 {
		query = args[0]
	}

	if len(args) > 1 {
		reader, err = os.Open(args[1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}
		defer reader.Close()
	}

	content, err := ioutil.ReadAll(reader)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	if query != "" && query != "." {
		value := gjson.Get(unsafeBytesToString(content), query)
		content = unsafeFastStringToReadOnlyBytes(value.String())
	}

	if output {
		var buf bytes.Buffer
		err = json.Indent(&buf, content, "", "  ")
		if err != nil {
			fmt.Fprintf(os.Stdout, "%s", unsafeBytesToString(content))
		} else {
			fmt.Fprintf(os.Stdout, "%s", unsafeBytesToString(buf.Bytes()))
		}
		return
	}

	os.Exit(run(content))
}

func unsafeBytesToString(b []byte) string {
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sh := reflect.StringHeader{bh.Data, bh.Len}
	return *(*string)(unsafe.Pointer(&sh))
}

func unsafeFastStringToReadOnlyBytes(s string) []byte {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := reflect.SliceHeader{sh.Data, sh.Len, sh.Len}
	return *(*[]byte)(unsafe.Pointer(&bh))
}

func genTree(content []byte) (*jsontree.JsonTree, error) {
	writer := colorwriter.New(colorMap, termbox.ColorDefault)
	formatter := jsonfmt.New(content, writer)
	if err := formatter.Format(); err != nil {
		return nil, err
	}
	formattedJson := writer.Lines

	tree := jsontree.New(formattedJson)
	return tree, nil
}

func run(content []byte) int {
	raw := content
	tree, err := genTree(content)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return 1
	}
	root := tree

	term, err := terminal.New(tree)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return 1
	}

	defer term.Close()
	for {
		tree, err = genTree(content)
		if err != nil {
			tree = root
		}
		term.Change(tree)

		for {
			term.Render()
			e := term.Poll()
			if e.Key == termbox.KeyCtrlC {
				return 0
			}
			query := handleKeypress(term, tree, e)
			if query {
				value := gjson.Get(unsafeBytesToString(raw), string(term.Query))
				content = unsafeFastStringToReadOnlyBytes(value.String())
				break
			}
		}
	}
}

func handleKeypress(t *terminal.Terminal, j *jsontree.JsonTree, e termbox.Event) bool {
	query := false
	if e.Ch == 0 {
		switch e.Key {
		case termbox.KeyBackspace, termbox.KeyBackspace2:
			t.DelQuery()
		case termbox.KeyArrowUp:
			t.MoveCursor(0, -1)
		case termbox.KeyArrowDown:
			t.MoveCursor(0, +1)
		case termbox.KeyArrowLeft:
			j.ToggleLine(t.CursorY + t.OffsetY)
		case termbox.KeyArrowRight:
			j.ToggleLine(t.CursorY + t.OffsetY)
		case termbox.KeyEnter:
			query = true
		case termbox.KeySpace:
			j.ToggleLine(t.CursorY + t.OffsetY)
		case termbox.KeyTab:
			j.ToggleLine(t.CursorY + t.OffsetY)
		case termbox.KeyCtrlU:
			query = t.ClearQuery()
		case termbox.KeyCtrlW:
			query = t.DelWordQuery()
		}
	} else {
		t.AddQuery(e.Ch)
	}
	return query
}
