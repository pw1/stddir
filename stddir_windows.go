package stddir

var (
	cacheEntries = []dirDef{
		dirDef{Path: "${LOCALAPPDATA}\\<program>\\cache", User: true},
		dirDef{Path: "${ProgramData}\\<program>\\cache"},
	}

	configEntries = []dirDef{
		dirDef{Path: "${APPDATA}\\<program>", User: true, Roaming: true},
		dirDef{Path: "${LOCALAPPDATA}\\<program>", User: true},
		dirDef{Path: "${ProgramData}\\<program>"},
	}
)
