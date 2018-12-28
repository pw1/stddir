package stddir

var (
	cacheEntries = []dirDef{
		dirDef{Path: "${XDG_CACHE_HOME}/<program>", AltPath: "${HOME}/.cache/<program>", User: true},
		dirDef{Path: "/var/cache/<program>"},
	}

	configEntries = []dirDef{
		dirDef{Path: "${XDG_CONFIG_HOME}/<program>", AltPath: "${HOME}/.config/<program>", User: true},
		dirDef{Path: "${XDG_CONFIG_DIRS}/<program>", AltPath: "/etc/xdg/<program>", List: true},
		dirDef{Path: "/etc/<program>"},
	}

	roamingConfigEntries = []dirDef{}
)
