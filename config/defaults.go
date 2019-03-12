package config

import "runtime"

func DefaultConfig() *Config {
	return &Config{
		DebugMode: false,
		ColourScheme: ColourScheme{
			Cursor:       strToColourNoErr("#e8dfd6"),
			Foreground:   strToColourNoErr("#e8dfd6"),
			Background:   strToColourNoErr("#021b21"),
			Black:        strToColourNoErr("#000000"),
			Red:          strToColourNoErr("#800000"),
			Green:        strToColourNoErr("#008000"),
			Yellow:       strToColourNoErr("#808000"),
			Blue:         strToColourNoErr("#000080"),
			Magenta:      strToColourNoErr("#800080"),
			Cyan:         strToColourNoErr("#008080"),
			LightGrey:    strToColourNoErr("#f2f2f2"),
			DarkGrey:     strToColourNoErr("#808080"),
			LightRed:     strToColourNoErr("#ff0000"),
			LightGreen:   strToColourNoErr("#00ff00"),
			LightYellow:  strToColourNoErr("#ffff00"),
			LightBlue:    strToColourNoErr("#0000ff"),
			LightMagenta: strToColourNoErr("#ff00ff"),
			LightCyan:    strToColourNoErr("#00ffff"),
			White:        strToColourNoErr("#ffffff"),
			Selection:    strToColourNoErr("#333366"),
		},
		KeyMapping: KeyMappingConfig(map[string]string{
			string(ActionCopy):        addMod("c"),
			string(ActionPaste):       addMod("v"),
			string(ActionSearch):      addMod("g"),
			string(ActionToggleDebug): addMod("d"),
			string(ActionToggleSlomo): addMod(";"),
			string(ActionReportBug):   addMod("r"),
			string(ActionBufferClear): addMod("k"),
		}),
		SearchURL:             "https://www.google.com/search?q=$QUERY",
		MaxLines:              1000,
		CopyAndPasteWithMouse: true,
	}
}

func addMod(keys string) string {
	standardMod := "ctrl + shift + "
	if runtime.GOOS == "darwin" {
		standardMod = "super + "
	}
	return standardMod + keys
}
