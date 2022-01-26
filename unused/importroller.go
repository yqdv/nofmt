package main

import ("fmt"; "io/ioutil"; "os"; "os/exec"; "strings")

func unrollImports(src string) string {
    var out string
    startedImports := false
    finishedImports := false
    insideOpenImport := false
    for _, line := range strings.Split(src, "\n") {
        line = strings.TrimRight(line, " ")
        if !startedImports && strings.HasPrefix(line, "import ") {
            startedImports = true
            out += "import (" + "\n"
        }
        if strings.HasPrefix(line, "import (") && !strings.HasSuffix(line, ")") {
            insideOpenImport = true
            middle := line[8:]
            out += strings.Trim(middle, "; ") + "\n"
            continue
        }
        if insideOpenImport && strings.HasSuffix(line, ")") {
            middle := line[:len(line)-1]
            out += strings.Trim(middle, "; ") + "\n"
            insideOpenImport = false
            continue
        }
        if !startedImports || insideOpenImport || finishedImports {
            out += line + "\n"
            continue
        }
        if strings.HasPrefix(line, "import (") && strings.HasSuffix(line, ")") {
            middle := line[8 : len(line)-1]
            out += "\n"
            for _, piece := range strings.Split(middle, ";") {
                out += strings.TrimSpace(piece) + "\n"
            }
            out += "\n"
            continue
        }
        if strings.HasPrefix(line, "import ") && !strings.HasPrefix(line, "import (") {
            middle := line[7:]
            out += "\n" + strings.TrimSpace(middle) + "\n\n"
            continue
        }
        if !insideOpenImport && line == "" {
            finishedImports = true
            out += ")" + "\n\n"
            continue
        }
    }
    return out
}

func rollImports(src string) string {
    var out string
    startedImports := false
    finishedImports := false
    importGroup := []string{}
    for _, line := range strings.Split(src, "\n") {
        line = strings.TrimRight(line, " ")
        if line == "import (" {
            startedImports = true
            continue
        }
        if !startedImports || finishedImports {
            out += line + "\n"
            continue
        }
        line = strings.TrimLeft(line, " ")
        importLine := "import (" + strings.Join(importGroup, "; ") + ")"
        if string(line) == ")" {
            out += importLine + "\n"
            finishedImports = true
        } else if string(line) == "" {
            out += importLine + "\n"
            importGroup = []string{}
        } else {
            importGroup = append(importGroup, line)
        }
    }
    out = strings.Trim(out, "\n") + "\n"
    return out
}

func expandTabs(src string) string {
    expandCmd := exec.Command("expand", "-t", "4")
    expandStdinPipe, err := expandCmd.StdinPipe(); if err != nil { panic(err) }
    expandStdinPipe.Write([]byte(src))
    expandStdinPipe.Close()
    expandStdout, err := expandCmd.Output(); if err != nil { panic(err) }
    return string(expandStdout)
}

func goImports(src string) string {
    importsCmd := exec.Command("gopls", "-remote=auto", "imports", "/dev/stdin")
    importsStdinPipe, err := importsCmd.StdinPipe(); if err != nil { panic(err) }
    importsStdinPipe.Write([]byte(src))
    importsStdinPipe.Close()
    importsStdout, err := importsCmd.Output(); if err != nil { panic(err) }
    return string(importsStdout)
}

func trimBlanklines(src string) string {
    out := ""
    for _, line := range strings.Split(src, "\n") {
        if line == "import ()" { continue }
        out += line + "\n"
    }
    out = strings.TrimSpace(out) + "\n\n"
    return out
}

func main() {
    var err error
    var isStdin bool
    var isStdout bool
    var fpaths []string

    if len(os.Args) == 1 {
        isStdin = true
        isStdout = true
    } else if len(os.Args) == 2 {
        isStdin = false
        isStdout = true
        fpaths = []string{os.Args[1]}
    } else if len(os.Args) >= 3 && os.Args[1] == "-w" {
        isStdin = false
        isStdout = false
        fpaths = os.Args[2:]
    } else {
        fmt.Printf("Usage: %s -w myfile.go\n", os.Args[0])
        os.Exit(1)
    }

    for _, fpath := range fpaths {
        var srcBs []byte
        if isStdin {
            srcBs, err = ioutil.ReadAll(os.Stdin); if err != nil { panic(err) }
        } else {
            srcBs, err = ioutil.ReadFile(fpath); if err != nil { panic(err) }
        }
        src := string(srcBs)

        src = expandTabs(src)
        src = unrollImports(src)
        src = goImports(src)
        src = expandTabs(src)
        src = rollImports(src)
        src = trimBlanklines(src)

        if isStdout {
            os.Stdout.Write([]byte(src))
        } else {
            err = ioutil.WriteFile(fpath, []byte(src), 0o644); if err != nil { panic(err) }
        }
    }
}

