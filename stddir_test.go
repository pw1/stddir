package stddir

import (
	"github.com/stretchr/testify/assert"
	"os"
	"strings"
	"testing"
)

func restoreEnv(env []string) {
	// Restore the original environment
	for _, keyval := range env {
		parts := strings.SplitN(keyval, "=", 2)
		os.Setenv(parts[0], parts[1])
	}
}

func TestConfig(t *testing.T) {
	dirs := Config("foobar")
	assert.True(t, len(dirs) > 0)
}

func TestCache(t *testing.T) {
	dirs := Cache("foobar")
	assert.True(t, len(dirs) > 0)
}

func TestFindEnvVar(t *testing.T) {
	i1, i2, varName := findEnvVar("")
	assert.Equal(t, -1, i1)
	assert.Equal(t, -1, i2)
	assert.Equal(t, "", varName)

	i1, i2, varName = findEnvVar("${HOME}/test")
	assert.Equal(t, 0, i1)
	assert.Equal(t, 6, i2)
	assert.Equal(t, "HOME", varName)

	i1, i2, varName = findEnvVar("${HOME}")
	assert.Equal(t, 0, i1)
	assert.Equal(t, 6, i2)
	assert.Equal(t, "HOME", varName)

	i1, i2, varName = findEnvVar("/bla/${FOO}/test")
	assert.Equal(t, 5, i1)
	assert.Equal(t, 10, i2)
	assert.Equal(t, "FOO", varName)

	i1, i2, varName = findEnvVar("/bla/${FOO/test") // End } is missing
	assert.Equal(t, -1, i1)
	assert.Equal(t, -1, i2)
	assert.Equal(t, "", varName)

	i1, i2, varName = findEnvVar("/bla/$FOO}/test") // Start { is missing}
	assert.Equal(t, -1, i1)
	assert.Equal(t, -1, i2)
	assert.Equal(t, "", varName)

	i1, i2, varName = findEnvVar("/bla/${}/test")
	assert.Equal(t, 5, i1)
	assert.Equal(t, 7, i2)
	assert.Equal(t, "", varName)

	i1, i2, varName = findEnvVar("/bla/${FOO}/test/${BAR}/abc")
	assert.Equal(t, 5, i1)
	assert.Equal(t, 10, i2)
	assert.Equal(t, "FOO", varName)
}

func TestProcessDirDef(t *testing.T) {
	// Clear the environment, so that we can control the variables
	env := os.Environ()
	os.Clearenv()
	defer restoreEnv(env)

	// There are no environment variables yet
	e := dirDef{
		Path:    "${HOME}/test/.<program>",
		AltPath: "/somewhere/else/<program>",
		User:    true,
	}
	dirs := processDirDef("foobar", e)
	assert.Equal(t, 1, len(dirs))
	assert.Equal(t, "/somewhere/else/foobar", dirs[0].Path)
	assert.True(t, dirs[0].User)

	e = dirDef{
		Path: "${HOME}/test/.<program>",
		User: true,
	}
	dirs = processDirDef("foobar", e)
	assert.Equal(t, 0, len(dirs))

	os.Setenv("HOME", "/home/janedoe")
	e = dirDef{
		Path:    "${HOME}/test/.<program>",
		AltPath: "/somewhere/else/<program>",
		User:    true,
	}
	dirs = processDirDef("foobar", e)
	assert.Equal(t, 1, len(dirs))
	assert.Equal(t, "/home/janedoe/test/.foobar", dirs[0].Path)
	assert.True(t, dirs[0].User)

	os.Setenv("HOME", "/home/janedoe")
	e = dirDef{
		Path:    "${FOO}/somewhere/else/<program>",
		AltPath: "${HOME}/test/.<program>",
		User:    true,
	}
	dirs = processDirDef("foobar", e)
	assert.Equal(t, 1, len(dirs))
	assert.Equal(t, "/home/janedoe/test/.foobar", dirs[0].Path)
	assert.True(t, dirs[0].User)

	os.Setenv("HOME", "/home/janedoe")
	e = dirDef{
		Path: "${HOME}/test/${FOO}/.<program>",
		User: true,
	}
	dirs = processDirDef("foobar", e)
	assert.Equal(t, 0, len(dirs))

	os.Setenv("ONE", "/one:/alpha")
	os.Setenv("TWO", "two:beta")
	e = dirDef{
		Path: "${ONE}/${TWO}/<program>",
	}
	dirs = processDirDef("foobar", e)
	assert.Equal(t, 1, len(dirs))
	assert.Equal(t, "/one:/alpha/two:beta/foobar", dirs[0].Path)
	assert.False(t, dirs[0].User)

	os.Setenv("ONE", "/one:/alpha")
	os.Setenv("TWO", "two:beta")
	e = dirDef{
		Path: "${ONE}/${TWO}/<program>",
		List: true,
	}
	dirs = processDirDef("foobar", e)
	assert.Equal(t, 4, len(dirs))
	assert.Equal(t, "/one/two/foobar", dirs[0].Path)
	assert.False(t, dirs[0].User)
	assert.Equal(t, "/one/beta/foobar", dirs[1].Path)
	assert.False(t, dirs[1].User)
	assert.Equal(t, "/alpha/two/foobar", dirs[2].Path)
	assert.False(t, dirs[2].User)
	assert.Equal(t, "/alpha/beta/foobar", dirs[3].Path)
	assert.False(t, dirs[3].User)
}

func TestCreateDirList(t *testing.T) {
	// Clear the environment, so that we can control the variables
	env := os.Environ()
	os.Clearenv()
	defer restoreEnv(env)

	// There are no environment variables yet
	os.Setenv("HOME", "/home/janedoe")
	e1 := dirDef{
		Path:    "${HOME}/.<program>",
		AltPath: "/somewhere/else/<program>",
		User:    true,
	}
	dirs := createDirList("foobar", []dirDef{e1})
	assert.Equal(t, 1, len(dirs))
	assert.Equal(t, "/home/janedoe/.foobar", dirs[0].Path)
	assert.True(t, dirs[0].User)

	e2 := dirDef{
		Path:    "${FOO}/.<program>",
		AltPath: "/somewhere/else/<program>",
	}
	e3 := dirDef{
		Path: "${FOO}/.<program>",
		List: true,
	}
	dirs = createDirList("foobar", []dirDef{e1, e2, e3})
	assert.Equal(t, 2, len(dirs))
	assert.Equal(t, "/home/janedoe/.foobar", dirs[0].Path)
	assert.True(t, dirs[0].User)
	assert.Equal(t, "/somewhere/else/foobar", dirs[1].Path)
	assert.False(t, dirs[1].User)
}
