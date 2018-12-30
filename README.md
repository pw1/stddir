# Stddir


## Introduction
Stddir is a cross-platform libray for finding standard directories for configuration files, data
files, cached data and runtime files. The locations of these directories are different for different
operating systems. Supported operating systems are: Linux, Windows, MacOSX, FreeBSD, NetBSD and
OpenBSD.

Full documentation can be found on GoDoc: https://godoc.org/github.com/pw1/stddir


## Config directories
Config find directories where user-specific and system wide configuration is stored. Depending on
the operating system the array contains the following items:

### Linux
1. `${XDG_CONFIG_HOME}/<program>`
   * If `${XDG_CONFIG_HOME}` is undefined, then `${HOME}/.config/<program>`
   * For example: `/home/janedoe/.config/foobar`
2. `${XDG_CONFIG_DIRS}/<program>` (for each entry, XDG_CONFIG_DIRS may contain multiple items)
   * If `${XDG_CONFIG_DIRS}` is not defined, then `/etc/xdg/<program>`
   * For example: `/etc/xdg/foobar`
3. `/etc/<program>`
   * For example: `/etc/foobar`

### Windows
1. `%APPDATA%\<program>`
   * For example: `C:\Users\JaneDoe\AppData\Roaming\foobar`
2. `%LOCALAPPDATA%\<program>`
   * For example: `C:\Users\JaneDoe\AppData\Local\foobar`
3. `%ProgramData%\<program>`
   * For example: `C:\ProgramData\<program>`

### MacOSX
1. `~/Library/Application Support/<program>`
   * For example: `/Users/janedoe/Library/Application Support/foobar`
2. `/Library/Application Support/<program>`
   * For example: `/Library/Application Support/foobar`


## Cache directories
Cache finds directories where applications should cache information. Depending on the operating
system the array contains the following items:

### Linux:
1. `${XDG_CACHE_HOME}/<program>`
   * If `${XDG_CACHE_HOME}` is undefined, then `${HOME}/.cache/<program>`
   * For example: `/home/janedoe/.cache/foobar`
2. `/var/cache/<program>`
   For example: `/var/cache/foobar`

### Windows:
1. `%LOCALAPPDATA%\<program>\cache`
   * For example: `C:\Users\JaneDoe\AppData\Local\foobar\cache`
2. `%ProgramData%\<program>\cache`
   * For example: `C:\ProgramData\<program>\cache`

### MacOSX:
1. `~/Library/Caches/<program>`
   * For example: `/Users/janedoe/Library/Caches/foobar`
2. `/Library/Caches/<program>`
   * For example: `/Library/Caches/foobar`


## External resources
* http://standards.freedesktop.org/basedir-spec/basedir-spec-latest.html
* https://refspecs.linuxfoundation.org/fhs.shtml
* https://developer.apple.com/library/content/documentation/FileManagement/Conceptual/FileSystemProgrammingGuide/MacOSXDirectories/MacOSXDirectories.html
* https://en.wikipedia.org/wiki/Special_folder
* https://en.wikipedia.org/wiki/Directory_structure
