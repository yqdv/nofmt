package nofmt

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func redent(s string) string {
	return strings.TrimSpace(strings.Replace(s, "\n\t\t", "\n", -1)) + "\n"
}

var GOFMT_FORMAT_BIN = os.Getenv("HOME") + "/goroot/bin/gofmt.orig"
var GOFMT_IMPORTS_BIN = os.Getenv("HOME") + "/go/bin/goimports.orig"
var NOFMT_FORMAT_BIN = os.Getenv("HOME") + "/go/bin/gofmt-nofmt"
var NOFMT_IMPORTS_BIN = os.Getenv("HOME") + "/go/bin/goimports-nofmt"

func run_test(t *testing.T, input string, golden string) {
	tmpdir, err := os.MkdirTemp(os.TempDir(), "nofmt-tests"); if err != nil { panic(err) }

	inputFile := tmpdir + "/" + "infmt.go"
	nofmtFile := tmpdir + "/" + "nofmt.go"
	gofmtFile := tmpdir + "/" + "gofmt.go"
	err = os.WriteFile(inputFile, []byte(input), 0o644); if err != nil { panic(err) }
	err = os.WriteFile(nofmtFile, []byte(input), 0o644); if err != nil { panic(err) }
	err = os.WriteFile(gofmtFile, []byte(input), 0o644); if err != nil { panic(err) }

	// Test goimports-nofmt vs golden
	cmd := exec.Command(NOFMT_IMPORTS_BIN, "-w", nofmtFile)
	err = cmd.Run(); if err != nil { panic(err) }
	nofmtOutput, err := os.ReadFile(nofmtFile); if err != nil { panic(err) }

	if string(nofmtOutput) != golden {
		fmt.Println("===== NOFMT FAILURE =====")
		fmt.Println("----- input ----")
		fmt.Printf("%s", input)
		fmt.Println("----- output ----")
		fmt.Printf("%s", nofmtOutput)
		fmt.Println("----- expected ----")
		fmt.Printf("%s", golden)
		fmt.Println("----- end ----")
		t.FailNow()
	}

	// Test goimports-nofmt vs itself
	cmd = exec.Command(NOFMT_IMPORTS_BIN, "-w", nofmtFile)
	err = cmd.Run(); if err != nil { panic(err) }
	nofmtOutput2, err := os.ReadFile(nofmtFile); if err != nil { panic(err) }

	if string(nofmtOutput2) != string(nofmtOutput) {
		fmt.Println("===== NOFMT IDEMPOTENCE FAILURE =====")
		fmt.Println("----- input ----")
		fmt.Printf("%s", input)
		fmt.Println("----- output ----")
		fmt.Printf("%s", nofmtOutput2)
		fmt.Println("----- expected ----")
		fmt.Printf("%s", nofmtOutput)
		fmt.Println("----- end ----")
		t.FailNow()
	}

	// Test goimports-nofmt vs goimports
	cmd = exec.Command(GOFMT_IMPORTS_BIN, "-w", gofmtFile)
	err = cmd.Run(); if err != nil { panic(err) }
	gofmtOutput, err := os.ReadFile(gofmtFile); if err != nil { panic(err) }

	cmd = exec.Command(NOFMT_IMPORTS_BIN, "-w", gofmtFile)
	err = cmd.Run(); if err != nil { panic(err) }
	nofmtOutput, err = os.ReadFile(gofmtFile); if err != nil { panic(err) }

	if string(nofmtOutput) != string(gofmtOutput) {
		fmt.Println("===== GOFMT FAILURE =====")
		fmt.Println("----- input ----")
		fmt.Printf("%s", gofmtOutput)
		fmt.Println("----- output ----")
		fmt.Printf("%s", nofmtOutput)
		fmt.Println("----- expected ----")
		fmt.Printf("%s", gofmtOutput)
		fmt.Println("----- end ----")
		t.FailNow()
	}

	err = os.RemoveAll(tmpdir); if err != nil { panic(err) }
}

func TestNoopNonflat(t *testing.T) {
	input := redent(`
		package main
		
		import "fmt"
		
		func main() {
			if true {
				fmt.Println("A")
			} else {
				fmt.Println("B")
			}
		}
	`)
	golden := redent(`
		package main
		
		import "fmt"
		
		func main() {
			if true {
				fmt.Println("A")
			} else {
				fmt.Println("B")
			}
		}
	`)
	run_test(t, input, golden)
}

func TestNoopFlat1(t *testing.T) {
	input := redent(`
		package main
		
		import "fmt"
		
		func main() {
			if true {
				fmt.Println("A")
			} else { fmt.Println("B") }
		}
	`)
	golden := redent(`
		package main
		
		import "fmt"
		
		func main() {
			if true {
				fmt.Println("A")
			} else { fmt.Println("B") }
		}
	`)
	run_test(t, input, golden)
}

func TestNoopFlat2(t *testing.T) {
	input := redent(`
		package main
		
		import "fmt"
		
		func main() {
			if true { fmt.Println("A")
			} else { fmt.Println("B") }
		}
	`)
	golden := redent(`
		package main
		
		import "fmt"
		
		func main() {
			if true { fmt.Println("A")
			} else { fmt.Println("B") }
		}
	`)
	run_test(t, input, golden)
}

func TestNoopFlat3(t *testing.T) {
	input := redent(`
		package main
		
		import "fmt"
		
		func main() {
			if true { fmt.Println("A") } else { fmt.Println("B") }
		}
	`)
	golden := redent(`
		package main
		
		import "fmt"
		
		func main() {
			if true { fmt.Println("A") } else { fmt.Println("B") }
		}
	`)
	run_test(t, input, golden)
}

func TestNoopFlatIf1(t *testing.T) {
	input := redent(`
		package main
		
		import "os"
		
		func main() {
			if _, err := os.ReadFile(""); err != nil { panic(err) }
		}
	`)
	golden := redent(`
		package main
		
		import "os"
		
		func main() {
			if _, err := os.ReadFile(""); err != nil { panic(err) }
		}
	`)
	run_test(t, input, golden)
}

func TestNoopFlatIf2(t *testing.T) {
	input := redent(`
		package main
		
		import "os"
		
		func main() {
			_, err := os.ReadFile(""); if err != nil { panic(err) }
		}
	`)
	golden := redent(`
		package main
		
		import "os"
		
		func main() {
			_, err := os.ReadFile(""); if err != nil { panic(err) }
		}
	`)
	run_test(t, input, golden)
}

func TestAddSpaceFlat1(t *testing.T) {
	input := redent(`
		package main
		
		import "fmt"
		
		func main() {
			if true {
				fmt.Println("A")
			} else {fmt.Println("B")}
		}
	`)
	golden := redent(`
		package main
		
		import "fmt"
		
		func main() {
			if true {
				fmt.Println("A")
			} else { fmt.Println("B") }
		}
	`)
	run_test(t, input, golden)
}

func TestAddSpaceFlat2(t *testing.T) {
	input := redent(`
		package main
		
		import "fmt"
		
		func main() {
			if true {fmt.Println("A")} else {fmt.Println("B")}
		}
	`)
	golden := redent(`
		package main
		
		import "fmt"
		
		func main() {
			if true { fmt.Println("A") } else { fmt.Println("B") }
		}
	`)
	run_test(t, input, golden)
}

func TestAddSpaceFlatIf1(t *testing.T) {
	input := redent(`
		package main
		
		import "os"
		
		func main() {
			if _,err:=os.ReadFile("");err!=nil{panic(err)}
		}
	`)
	golden := redent(`
		package main
		
		import "os"
		
		func main() {
			if _, err := os.ReadFile(""); err != nil { panic(err) }
		}
	`)
	run_test(t, input, golden)
}

func TestAddSpaceFlatIf2(t *testing.T) {
	input := redent(`
		package main
		
		import "os"
		
		func main() {
			_,err:=os.ReadFile("");if err!=nil{panic(err)}
		}
	`)
	golden := redent(`
		package main
		
		import "os"
		
		func main() {
			_, err := os.ReadFile(""); if err != nil { panic(err) }
		}
	`)
	run_test(t, input, golden)
}

func TestImportsNoop(t *testing.T) {
	input := redent(`
		package main
		
		import (
			"fmt"
			"os"
		)
		
		func main() {
			if true {
				fmt.Println("A")
			} else {
				fmt.Println("B")
				os.ReadFile("")
			}
		}
	`)
	golden := redent(`
		package main
		
		import (
			"fmt"
			"os"
		)
		
		func main() {
			if true {
				fmt.Println("A")
			} else {
				fmt.Println("B")
				os.ReadFile("")
			}
		}
	`)
	run_test(t, input, golden)
}

func TestImportsAdd(t *testing.T) {
	input := redent(`
		package main
		
		import "fmt"
		
		func main() {
			if true {
				fmt.Println("A")
			} else {
				fmt.Println("B")
				os.ReadFile("")
			}
		}
	`)
	golden := redent(`
		package main
		
		import (
			"fmt"
			"os"
		)
		
		func main() {
			if true {
				fmt.Println("A")
			} else {
				fmt.Println("B")
				os.ReadFile("")
			}
		}
	`)
	run_test(t, input, golden)
}

func TestImportsRemove(t *testing.T) {
	input := redent(`
		package main
		
		import (
			"fmt"
			"os"
		)
		
		func main() {
			if true {
				fmt.Println("A")
			} else {
				fmt.Println("B")
			}
		}
	`)
	golden := redent(`
		package main
		
		import (
			"fmt"
		)
		
		func main() {
			if true {
				fmt.Println("A")
			} else {
				fmt.Println("B")
			}
		}
	`)
	run_test(t, input, golden)
}

func TestImportsReorder(t *testing.T) {
	input := redent(`
		package main
		
		import (
			"os"
			"fmt"
		)
		
		func main() {
			if true {
				fmt.Println("A")
			} else {
				fmt.Println("B")
				os.ReadFile("")
			}
		}
	`)
	golden := redent(`
		package main
		
		import (
			"fmt"
			"os"
		)
		
		func main() {
			if true {
				fmt.Println("A")
			} else {
				fmt.Println("B")
				os.ReadFile("")
			}
		}
	`)
	run_test(t, input, golden)
}

func TestImportsReorderUnroll(t *testing.T) {
	input := redent(`
		package main
		
		import ("os"; "fmt")
		
		func main() {
			if true {
				fmt.Println("A")
			} else {
				fmt.Println("B")
				os.ReadFile("")
			}
		}
	`)
	golden := redent(`
		package main
		
		import (
			"fmt"
			"os"
		)
		
		func main() {
			if true {
				fmt.Println("A")
			} else {
				fmt.Println("B")
				os.ReadFile("")
			}
		}
	`)
	run_test(t, input, golden)
}

func TestImportsReorderCombine(t *testing.T) {
	input := redent(`
		package main
		
		import "os"
		import "fmt"
		
		func main() {
			if true {
				fmt.Println("A")
			} else {
				fmt.Println("B")
				os.ReadFile("")
			}
		}
	`)
	golden := redent(`
		package main
		
		import (
			"fmt"
			"os"
		)
		
		func main() {
			if true {
				fmt.Println("A")
			} else {
				fmt.Println("B")
				os.ReadFile("")
			}
		}
	`)
	run_test(t, input, golden)
}

func TestImportsNoopGap(t *testing.T) {
	input := redent(`
		package main
		
		import (
			"fmt"
		
			"os"
		)
		
		func main() {
			if true {
				fmt.Println("A")
			} else {
				fmt.Println("B")
				os.ReadFile("")
			}
		}
	`)
	golden := redent(`
		package main
		
		import (
			"fmt"
		
			"os"
		)
		
		func main() {
			if true {
				fmt.Println("A")
			} else {
				fmt.Println("B")
				os.ReadFile("")
			}
		}
	`)
	run_test(t, input, golden)
}

func TestNoopNestedFunc(t *testing.T) {
	input := redent(`
		package main
		
		func main() {
			if err := c.aaaa(&bbbb, cccc, "dddddddddd", eeeeeeee, fffffffff,
				func(ggggggg, hhhhhhh interface{}) (bool, error) {
					return iiiiiiiiiii(jjjjjjjjjjj.(string), kkkkkkkkkk.(string))
				}, c.lllllllllllll, c.mmmmmmmmmmmmmm); err != nil {
				return err
			}
		}
	`)
	golden := redent(`
		package main
		
		func main() {
			if err := c.aaaa(&bbbb, cccc, "dddddddddd", eeeeeeee, fffffffff,
				func(ggggggg, hhhhhhh interface{}) (bool, error) {
					return iiiiiiiiiii(jjjjjjjjjjj.(string), kkkkkkkkkk.(string))
				}, c.lllllllllllll, c.mmmmmmmmmmmmmm); err != nil {
				return err
			}
		}
	`)
	run_test(t, input, golden)
}

func TestJaggedIndent(t *testing.T) {
	input := redent(`
		package main
		
		import "os/exec"
		
			func expandTabs(src string) string {
			    expandCmd := exec.Command("expand", "-t", "4")
			 expandStdinPipe, err := expandCmd.StdinPipe(); if err!=nil {panic(err)}
			 expandStdinPipe.Write([]byte(src))
			  expandStdinPipe.Close()
			       expandStdout, err:=expandCmd.Output();if err != nil { panic(err)}
			    return string(expandStdout)
		}
	`)
	golden := redent(`
		package main
		
		import "os/exec"
		
		func expandTabs(src string) string {
			expandCmd := exec.Command("expand", "-t", "4")
			expandStdinPipe, err := expandCmd.StdinPipe(); if err != nil { panic(err) }
			expandStdinPipe.Write([]byte(src))
			expandStdinPipe.Close()
			expandStdout, err := expandCmd.Output(); if err != nil { panic(err) }
			return string(expandStdout)
		}
	`)
	run_test(t, input, golden)
}

func TestIotaNoop(t *testing.T) {
	input := redent(`
		package main
		
		const (
			x = iota
			y
			z
		)
		
		func main() {
			for x == y { }
		}
	`)
	golden := redent(`
		package main
		
		const (
			x = iota
			y
			z
		)
		
		func main() {
			for x == y { }
		}
	`)
	run_test(t, input, golden)
}

func TestIotaIndent(t *testing.T) {
	input := redent(`
		package main
		
		const (
		x = iota
		y
		z
		)
		
		func main() {
			for x == y { }
		}
	`)
	golden := redent(`
		package main
		
		const (
			x = iota
			y
			z
		)
		
		func main() {
			for x == y { }
		}
	`)
	run_test(t, input, golden)
}

func TestVarIndent(t *testing.T) {
	input := redent(`
		package main
		
		var (
		x string
						y int
		z uint64
		)
		
		func main() {
			for x == y { }
		}
	`)
	golden := redent(`
		package main
		
		var (
			x string
			y int
			z uint64
		)
		
		func main() {
			for x == y { }
		}
	`)
	run_test(t, input, golden)
}

func TestAddSpacesEmptyBlock1(t *testing.T) {
	input := redent(`
		package main
		
		func main() {
			for x==y{}
		}
	`)
	golden := redent(`
		package main
		
		func main() {
			for x == y { }
		}
	`)
	run_test(t, input, golden)
}

func TestAddSpacesEmptyBlock2(t *testing.T) {
	input := redent(`
		package main
		
		func main() {
			for x==y{
		}
		}
	`)
	golden := redent(`
		package main
		
		func main() {
			for x == y {
			}
		}
	`)
	run_test(t, input, golden)
}

func TestAddSpacesEmptyBlock3(t *testing.T) {
	input := redent(`
		package main
		
		func main() {
			for x==y  {
		                 }
		}
	`)
	golden := redent(`
		package main
		
		func main() {
			for x == y {
			}
		}
	`)
	run_test(t, input, golden)
}

func TestIndentComment1(t *testing.T) {
	input := redent(`
		package main
		
		func main() {
		// indent this
		// ...
		// ...
			x := 1
		}
	`)
	golden := redent(`
		package main
		
		func main() {
			// indent this
			// ...
			// ...
			x := 1
		}
	`)
	run_test(t, input, golden)
}

func TestIndentComment2(t *testing.T) {
	input := redent(`
		package main
		
		func main() {
			x := 1
		
		// indent this
		// ...
		// ...
		}
	`)
	golden := redent(`
		package main
		
		func main() {
			x := 1
		
			// indent this
			// ...
			// ...
		}
	`)
	run_test(t, input, golden)
}

func TestNoopSwitch(t *testing.T) {
	input := redent(`
		package main
		
		func main() {
			switch x {
			case 2: e = f
			case 1: { a = b }
			case 2: {
					c = d
				}
			default: { g = h }
			}
		}
	`)
	golden := redent(`
		package main
		
		func main() {
			switch x {
			case 2: e = f
			case 1: { a = b }
			case 2: {
					c = d
				}
			default: { g = h }
			}
		}
	`)
	run_test(t, input, golden)
}

func TestNoopTwoVar(t *testing.T) {
	input := redent(`
		package main

		func main() {
			x := 1; y := 2
		}
	`)
	golden := redent(`
		package main

		func main() {
			x := 1; y := 2
		}
	`)
	run_test(t, input, golden)
}

func TestNoopInitVar(t *testing.T) {
	input := redent(`
		package main

		func main() {
			var foo int; if bar == baz { foo = qux }
		}
	`)
	golden := redent(`
		package main

		func main() {
			var foo int; if bar == baz { foo = qux }
		}
	`)
	run_test(t, input, golden)
}

func TestNoopLabel(t *testing.T) {
	input := redent(`
		package main

		func main() {
			for a {
				for i := range b {
					if c {
						goto done
					} else if d {
						break
					}
					a = b[i]
				}
			done:
			}
		}
	`)
	golden := redent(`
		package main

		func main() {
			for a {
				for i := range b {
					if c {
						goto done
					} else if d {
						break
					}
					a = b[i]
				}
			done:
			}
		}
	`)
	run_test(t, input, golden)
}

func TestCloseBraceLabel(t *testing.T) {
	input := redent(`
		package main

		func main() {
			for a {
				for i := range b {
					if c {
						goto done
					} else if d {
						break
					}
					a = b[i]
				}
			done: }
		}
	`)
	golden := redent(`
		package main

		func main() {
			for a {
				for i := range b {
					if c {
						goto done
					} else if d {
						break
					}
					a = b[i]
				}
			done:
			}
		}
	`)
	run_test(t, input, golden)
}

func TestLoopLabelNoop(t *testing.T) {
	input := redent(`
		package main
		
		func main() {
		level0start:
			for y, row := range rows {
			level1start:
				for x, data := range row {
				level2start:
					if data == something {
						continue level1
					}
					row[x] = stuff
				level2end:
				}
			level1end:
			}
		level0end:
		}
	`)
	golden := redent(`
		package main
		
		func main() {
		level0start:
			for y, row := range rows {
			level1start:
				for x, data := range row {
				level2start:
					if data == something {
						continue level1
					}
					row[x] = stuff
				level2end:
				}
			level1end:
			}
		level0end:
		}
	`)
	run_test(t, input, golden)
}

func TestLoopLabelIndent(t *testing.T) {
	input := redent(`
		package main
		
		func main() {
		level0start:
			for y, row := range rows {
		level1start:
				for x, data := range row {
		level2start:
					if data == something {
						continue level1
					}
					row[x] = stuff
		level2end:
				}
		level1end:
			}
		level0end:
		}
	`)
	golden := redent(`
		package main
		
		func main() {
		level0start:
			for y, row := range rows {
			level1start:
				for x, data := range row {
				level2start:
					if data == something {
						continue level1
					}
					row[x] = stuff
				level2end:
				}
			level1end:
			}
		level0end:
		}
	`)
	run_test(t, input, golden)
}

func TestLoopLabelDedent(t *testing.T) {
	input := redent(`
		package main
		
		func main() {
						level0start:
			for y, row := range rows {
						level1start:
				for x, data := range row {
						level2start:
					if data == something {
						continue level1
					}
					row[x] = stuff
						level2end:
				}
						level1end:
			}
						level0end:
		}
	`)
	golden := redent(`
		package main
		
		func main() {
		level0start:
			for y, row := range rows {
			level1start:
				for x, data := range row {
				level2start:
					if data == something {
						continue level1
					}
					row[x] = stuff
				level2end:
				}
			level1end:
			}
		level0end:
		}
	`)
	run_test(t, input, golden)
}
