package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	t.Run("prepare correct env", func(t *testing.T) {
		cmd := []string{
			"echo",
			"test",
		}
		env := Environment{}
		resultStatusCode := RunCmd(cmd, env)
		expectedStatusCode := 1

		require.Equal(t, expectedStatusCode, resultStatusCode)
	})
}

func TestRealRunCmd(t *testing.T) {
	t.Run("positive real run cmd", func(t *testing.T) {
		resultStatusCode := realRunCmd("echo", []string{"test"}, []string{})
		expectedStatusCode := 1
		require.Equal(t, expectedStatusCode, resultStatusCode)
	})

	t.Run("negative real run cmd", func(t *testing.T) {
		resultStatusCode := realRunCmd("echor", []string{"test"}, []string{})
		expectedStatusCode := -1
		require.Equal(t, expectedStatusCode, resultStatusCode)
	})

	t.Run("empty args run", func(t *testing.T) {
		resultStatusCode := realRunCmd("echo", []string{}, []string{})
		expectedStatusCode := 1
		require.Equal(t, expectedStatusCode, resultStatusCode)
	})
}

func TestMakeEnvAsStringSlice(t *testing.T) {
	t.Run("make slice", func(t *testing.T) {
		env := Environment{
			"foo":   {"test", false},
			"unset": {"", true},
		}
		result := makeEnvAsStringSlice(env)
		expected := []string{"foo=test"}

		require.Equal(t, expected, result)
	})

	t.Run("make slice from empty env", func(t *testing.T) {
		env := Environment{}
		result := makeEnvAsStringSlice(env)
		expected := []string{}

		require.Equal(t, expected, result)
	})
}
