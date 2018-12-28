package stddir

var (
	cacheEntries = []dirDef{
		dirDef{Path: "${HOME}/Library/Caches/<program>", User: true},
		dirDef{Path: "/Library/Caches/<program>"},
	}

	configEntries = []dirDef{
		dirDef{Path: "${HOME}/Library/Application Support/<program>", User: true},
		dirDef{Path: "/Library/Application Support/<program>"},
	}

	roamingConfigEntries = []dirDef{}
)
