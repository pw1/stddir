package stddir

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
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

func TestProcessDirDef1(t *testing.T) {
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
	assert.False(t, dirs[0].Roaming)
}

func TestProcessDirDef2(t *testing.T) {
	// Clear the environment, so that we can control the variables
	env := os.Environ()
	os.Clearenv()
	defer restoreEnv(env)

	// There are no environment variables yet
	e := dirDef{
		Path: "${HOME}/test/.<program>",
		User: true,
	}
	dirs := processDirDef("foobar", e)
	assert.Equal(t, 0, len(dirs))
}

func TestProcessDirDef3(t *testing.T) {
	// Clear the environment, so that we can control the variables
	env := os.Environ()
	os.Clearenv()
	defer restoreEnv(env)

	os.Setenv("HOME", "/home/janedoe")
	e := dirDef{
		Path:    "${HOME}/test/.<program>",
		AltPath: "/somewhere/else/<program>",
		User:    true,
	}
	dirs := processDirDef("foobar", e)
	assert.Equal(t, 1, len(dirs))
	assert.Equal(t, "/home/janedoe/test/.foobar", dirs[0].Path)
	assert.True(t, dirs[0].User)
	assert.False(t, dirs[0].Roaming)
}

func TestProcessDirDef4(t *testing.T) {
	// Clear the environment, so that we can control the variables
	env := os.Environ()
	os.Clearenv()
	defer restoreEnv(env)

	os.Setenv("HOME", "/home/janedoe")
	e := dirDef{
		Path:    "${FOO}/somewhere/else/<program>",
		AltPath: "${HOME}/test/.<program>",
		User:    true,
	}
	dirs := processDirDef("foobar", e)
	assert.Equal(t, 1, len(dirs))
	assert.Equal(t, "/home/janedoe/test/.foobar", dirs[0].Path)
	assert.True(t, dirs[0].User)
	assert.False(t, dirs[0].Roaming)
}

func TestProcessDirDef5(t *testing.T) {
	// Clear the environment, so that we can control the variables
	env := os.Environ()
	os.Clearenv()
	defer restoreEnv(env)

	os.Setenv("HOME", "/home/janedoe")
	e := dirDef{
		Path: "${HOME}/test/${FOO}/.<program>",
		User: true,
	}
	dirs := processDirDef("foobar", e)
	assert.Equal(t, 0, len(dirs))
}

func TestProcessDirDef6(t *testing.T) {
	// Clear the environment, so that we can control the variables
	env := os.Environ()
	os.Clearenv()
	defer restoreEnv(env)

	os.Setenv("ONE", "/one"+string(os.PathListSeparator)+"/alpha")
	os.Setenv("TWO", "two"+string(os.PathListSeparator)+"beta")
	e := dirDef{
		Path: "${ONE}/${TWO}/<program>",
	}
	dirs := processDirDef("foobar", e)
	assert.Equal(t, 1, len(dirs))
	assert.Equal(t, "/one"+string(os.PathListSeparator)+"/alpha/two"+string(os.PathListSeparator)+"beta/foobar", dirs[0].Path)
	assert.False(t, dirs[0].User)
	assert.False(t, dirs[0].Roaming)
}

func TestProcessDirDef7(t *testing.T) {
	// Clear the environment, so that we can control the variables
	env := os.Environ()
	os.Clearenv()
	defer restoreEnv(env)

	os.Setenv("ONE", "/one"+string(os.PathListSeparator)+"/alpha")
	os.Setenv("TWO", "two"+string(os.PathListSeparator)+"beta")
	e := dirDef{
		Path: "${ONE}/${TWO}/<program>",
		List: true,
	}
	dirs := processDirDef("foobar", e)
	assert.Equal(t, 4, len(dirs))
	assert.Equal(t, "/one/two/foobar", dirs[0].Path)
	assert.False(t, dirs[0].User)
	assert.False(t, dirs[0].Roaming)
	assert.Equal(t, "/one/beta/foobar", dirs[1].Path)
	assert.False(t, dirs[1].User)
	assert.False(t, dirs[1].Roaming)
	assert.Equal(t, "/alpha/two/foobar", dirs[2].Path)
	assert.False(t, dirs[2].User)
	assert.False(t, dirs[2].Roaming)
	assert.Equal(t, "/alpha/beta/foobar", dirs[3].Path)
	assert.False(t, dirs[3].User)
	assert.False(t, dirs[3].Roaming)
}

func TestProcessDirDefRoamingTag(t *testing.T) {
	// Clear the environment, so that we can control the variables
	env := os.Environ()
	os.Clearenv()
	defer restoreEnv(env)

	// There are no environment variables yet
	e := dirDef{
		Path:    "/some/where/<program>",
		User:    true,
		Roaming: true,
	}
	dirs := processDirDef("foobar", e)
	assert.Equal(t, 1, len(dirs))
	assert.Equal(t, "/some/where/foobar", dirs[0].Path)
	assert.True(t, dirs[0].User)
	assert.True(t, dirs[0].Roaming)
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

func getStandardDirDefs() []dirDef {
	defs := []dirDef{
		dirDef{Path: "/the/place/one"},
		dirDef{Path: "/the/place/two", Roaming: true},
		dirDef{Path: "/the/place/three", User: true},
	}
	return defs
}

func TestCreateDirListExcludeRoaming(t *testing.T) {
	dirs := createDirList("foobar", getStandardDirDefs(), ExcludeRoaming)
	assert.Equal(t, 2, len(dirs))
	assert.Equal(t, "/the/place/one", dirs[0].Path)
	assert.Equal(t, "/the/place/three", dirs[1].Path)
}

func TestCreateDirListExcludeUser(t *testing.T) {
	dirs := createDirList("foobar", getStandardDirDefs(), ExcludeUser)
	assert.Equal(t, 2, len(dirs))
	assert.Equal(t, "/the/place/one", dirs[0].Path)
	assert.Equal(t, "/the/place/two", dirs[1].Path)
}

func TestCreateDirListExcludeSystem(t *testing.T) {
	dirs := createDirList("foobar", getStandardDirDefs(), ExcludeSystem)
	assert.Equal(t, 1, len(dirs))
	assert.Equal(t, "/the/place/three", dirs[0].Path)
}

func TestCreateDirListExcludeUserAndRoaming(t *testing.T) {
	dirs := createDirList("foobar", getStandardDirDefs(), ExcludeUser, ExcludeRoaming)
	assert.Equal(t, 1, len(dirs))
	assert.Equal(t, "/the/place/one", dirs[0].Path)
}
