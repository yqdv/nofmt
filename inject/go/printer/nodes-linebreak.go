func (p *printer) linebreak(line, min int, ws whiteSpace, newSection bool) (nbreaks int) {
	n := nlimit(line - p.pos.Line)
	if n > 0 {
		p.print(ws)
		if newSection {
			p.print(formfeed)
			n--
			nbreaks = 2
		}
		nbreaks += n
		for ; n > 0; n-- {
			p.print(newline)
		}
	} else if min > 0 {
		if !newSection {
			p.print(token.SEMICOLON)
		}
		p.print(blank)
	}
	return
}
