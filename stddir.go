package stddir

import (
	"os"
	"strings"
)

// Flag represent an option that influences which directories are returned. For example, the
// ExcludeRoaming leaves out any roaming directories.
type Flag int

const (
	// ExcludeRoaming is a Flag that indicates that roaming directories must not be returned.
	ExcludeRoaming Flag = iota

	// ExcludeUser is a Flag that indicates that user-specific directories must not be returned.
	ExcludeUser

	// ExcludeSystem is a Flag that indicates that system wide directories must not be returned.
	ExcludeSystem
)

// flagCollection is a helper class to make it easier to work with a collection of flags.
type flagCollection []Flag

// Contain returns true if the flagCollection contains the specified flag. Returns false otherwise.
func (f flagCollection) Contain(flag Flag) bool {
	for _, value := range f {
		if value == flag {
			return true
		}
	}

	return false
}

// dirDef defines a directory. This may resolve into multiple directories.
type dirDef struct {
	Path    string // Path to the directory, may contain environment variables and "<program>"
	AltPath string // Alternative Path. This is used if if Path can't be resolved (missing env var).
	List    bool   // True if environment variables may contain a list of paths
	User    bool   // Is this a user-specific directory
	Roaming bool   // Is this a roaming user profile directory.
}

// Dir represent a single directory.
type Dir struct {
	Path    string // Absolute path to the directory
	User    bool   // True if this is a user-specific directory, false otherwise.
	Roaming bool   // True if this is a roaming user profile directory, false otherwise.
}

// Cache finds directories where applications should cache information. An array with one
// or more directories is returned. The array is sorted by the importance. The directory with the
// highest importance is the first item.
//
// By adding flags to the function call you can exclude certain directories.
//
// Depending on the operating system the array contains the following items (in this order):
//
// Linux:
//   1. ${XDG_CACHE_HOME}/<program>
//      If ${XDG_CACHE_HOME} is undefined, then ${HOME}/.cache/<program>
//      For example: /home/janedoe/.cache/foobar
//   2. /var/cache/<program>
//      For example: /var/cache/foobar
//
// Windows:
//   1. %LOCALAPPDATA%\<program>\cache
//      For example: C:\Users\JaneDoe\AppData\Local\foobar\cache
//   2. %ProgramData%\<program>\cache
//      For example: C:\ProgramData\<program>\cache
//
// MacOSX:
//   1. ~/Library/Caches/<program>
//      For example: /Users/janedoe/Library/Caches/foobar
//   2. /Library/Caches/<program>
//      For example: /Library/Caches/foobar
func Cache(program string, flags ...Flag) []Dir {
	return createDirList(program, cacheEntries, flags...)
}

// Config finds directories where user-specific and system wide configuration is stored. An array
// with one or more directories is returned. The array is sorted by the importance. The directory
// with the highest importance is the first item.
//
// By adding flags to the function call you can exclude certain directories.
//
// Depending on the operating system the array contains the following items (in this order):
//
// Linux:
//   1. ${XDG_CONFIG_HOME}/<program>
//      If ${XDG_CONFIG_HOME} is undefined, then ${HOME}/.config/<program>
//      For example: /home/janedoe/.config/foobar
//   2. ${XDG_CONFIG_DIRS}/<program> (for each entry, XDG_CONFIG_DIRS may contain multiple items)
//      If ${XDG_CONFIG_DIRS} is not defined, then /etc/xdg/<program>
//      For example: /etc/xdg/foobar
//   3. /etc/<program>
//      For example: /etc/foobar
//
// Windows:
//   1. %APPDATA%\<program>
//      For example: C:\Users\JaneDoe\AppData\Roaming\foobar
//   2. %LOCALAPPDATA%\<program>
//      For example: C:\Users\JaneDoe\AppData\Local\foobar
//   3. %ProgramData%\<program>
//      For example: C:\ProgramData\<program>
//
// MacOSX:
//   1. ~/Library/Application Support/<program>
//      For example: /Users/janedoe/Library/Application Support/foobar
//   2. /Library/Application Support/<program>
//      For example: /Library/Application Support/foobar
func Config(program string, flags ...Flag) []Dir {
	return createDirList(program, configEntries, flags...)
}

// Resolve a list of directory definitions. Returns the resolved directories. If a directory
// definition can't be resolved (e.g. because of a missing environment variable), then it is omitted
// from the returned directory list.
func createDirList(program string, entries []dirDef, flags ...Flag) []Dir {
	// Convert the []Flag into a flagCollection so that I can use the methods on this helper class.
	var flagCol flagCollection = flags

	dirs := []Dir{}
	for _, entry := range entries {
		// Skip over entries that are exluded by a flag.
		if (entry.Roaming && flagCol.Contain(ExcludeRoaming)) ||
			(entry.User && flagCol.Contain(ExcludeUser)) ||
			(!entry.User && flagCol.Contain(ExcludeSystem)) {
			continue
		}

		dirs = append(dirs, processDirDef(program, entry)...)
	}
	return dirs
}

// Resolve a single directory definition (dirDef). Returns list with zero of more items.
func processDirDef(program string, e dirDef) []Dir {
	path := e.Path
	path = strings.Replace(path, "<program>", program, -1)

	// Replace all environment variables (there can be multiple)
	loopCounter := 0
	for true {
		loopCounter++
		if loopCounter > 100 {
			// This is a safety measure against infinite loops
			return []Dir{}
		}

		i1, i2, varName := findEnvVar(path)
		if i1 < 0 {
			// There are no more environment variables
			break
		}

		varValue := os.Getenv(varName)

		// Handle the case where the environment variable is not defined
		if varValue == "" {
			if e.AltPath != "" {
				// We use the AltPath (create a new dirDef with AltPath as Path and process that)
				altDirDef := dirDef(e)
				altDirDef.Path = e.AltPath
				altDirDef.AltPath = ""
				return processDirDef(program, altDirDef)
			}

			// The environment variable is not defined and there is no AltPath, so we return nothing.
			return []Dir{}
		}

		// Handle the case where the environment variable contains a list of paths
		if e.List {
			parts := strings.Split(varValue, string(os.PathListSeparator))
			if len(parts) > 1 {
				dirs := []Dir{}
				for _, part := range parts {
					newDirDef := dirDef(e)
					newDirDef.Path = path[:i1] + part + path[(i2+1):]
					dirs = append(dirs, processDirDef(program, newDirDef)...)
				}
				return dirs
			}
		}

		path = path[:i1] + varValue + path[(i2+1):]
	}

	dir := Dir{
		Path:    path,
		User:    e.User,
		Roaming: e.Roaming,
	}

	return []Dir{dir}
}

// Finds the first environment variable in a path. The environment variable is denoted by
// ${var-name} (independent of operating system)
func findEnvVar(path string) (int, int, string) {
	i1 := strings.Index(path, "${")
	if i1 < 0 {
		return -1, -1, ""
	}

	iTmp := i1 + 2
	i2 := strings.Index(path[iTmp:], "}")
	if i2 < 0 {
		return -1, -1, ""
	}
	i2 += iTmp

	return i1, i2, path[(i1 + 2):(i2)]
}
