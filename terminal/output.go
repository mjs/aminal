package terminal

import (
	"time"
)

// Wish list here: http://invisible-island.net/xterm/ctlseqs/ctlseqs.html

type TerminalCharSet int

// single rune handler
type runeHandler func(terminal *Terminal) error

type escapeSequenceHandler func(pty chan rune, terminal *Terminal) error

var runeMap = map[rune]runeHandler{
	0x05: enqHandler,
	0x07: bellHandler,
	0x08: backspaceHandler,
	0x09: tabHandler,
	0x0a: newLineHandler,
	0x0b: newLineHandler,
	0x0c: newLineHandler,
	0x0d: carriageReturnHandler,
	0x0e: shiftOutHandler,
	0x0f: shiftInHandler,
}

func newLineHandler(terminal *Terminal) error {
	terminal.ActiveBuffer().NewLine()
	terminal.NotifyDirty()
	return nil
}

func tabHandler(terminal *Terminal) error {
	terminal.ActiveBuffer().Tab()
	terminal.NotifyDirty()
	return nil
}

func carriageReturnHandler(terminal *Terminal) error {
	terminal.ActiveBuffer().CarriageReturn()
	terminal.NotifyDirty()
	return nil
}

func backspaceHandler(terminal *Terminal) error {
	terminal.ActiveBuffer().Backspace()
	terminal.NotifyDirty()
	return nil
}

func bellHandler(terminal *Terminal) error {
	// @todo ring bell - flash red or some shit?
	return nil
}

func enqHandler(terminal *Terminal) error {
	terminal.logger.Errorf("Received ENQ!")
	return nil
}

func shiftOutHandler(terminal *Terminal) error {
	terminal.logger.Debugf("Received shift out")
	terminal.terminalState.CurrentCharset = 1
	return nil
}

func shiftInHandler(terminal *Terminal) error {
	terminal.logger.Debugf("Received shift in")
	terminal.terminalState.CurrentCharset = 0
	return nil
}

func (terminal *Terminal) processRuneLocked(b rune) {
	terminal.Lock()
	defer terminal.Unlock()

	terminal.processRune(b)
}

func (terminal *Terminal) processRune(b rune) {
	defer terminal.NotifyDirty()

	if handler, ok := runeMap[b]; ok {
		if err := handler(terminal); err != nil {
			terminal.logger.Errorf("Error handling control code: %s", err)
		}
		return
	}
	//terminal.logger.Debugf("Received character 0x%X: %q", b, string(b))
	terminal.ActiveBuffer().Write(terminal.translateRune(b))
}

func (terminal *Terminal) translateRune(b rune) rune {
	table := terminal.terminalState.Charsets[terminal.terminalState.CurrentCharset]
	if table == nil {
		return b
	}
	chr, ok := (*table)[b]
	if ok {
		return chr
	}
	return b
}

func (terminal *Terminal) processInput(pty chan rune) {

	// https://en.wikipedia.org/wiki/ANSI_escape_code

	var b rune

	//	debug := ""

	for {

		if terminal.config.Slomo {
			time.Sleep(time.Millisecond * 100)
		}

		b = <-pty

		// debug += fmt.Sprintf("0x%x ", b)

		if b == 0x1b {
			// terminal.logger.Debug(debug)
			// debug = ""
			//terminal.logger.Debugf("Handling escape sequence: 0x%x", b)
			if err := ansiHandler(pty, terminal); err != nil {
				terminal.logger.Errorf("Error handling escape sequence: %s", err)
			}
			terminal.NotifyDirty()
			continue
		}

		terminal.processRuneLocked(b)
	}
}
