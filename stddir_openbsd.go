package stddir

var (
	cacheEntries = []dirDef{
		dirDef{Path: "${HOME}/.<program>/cache", User: true},
		dirDef{Path: "/var/cache/<program>"},
	}

	configEntries = []dirDef{
		dirDef{Path: "${HOME}/.<program>", User: true},
		dirDef{Path: "/etc/<program>"},
	}
)
