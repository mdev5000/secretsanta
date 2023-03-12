package config

import (
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

type configTestCase struct {
	env string
}

func eq[T any](ptr *T, expected T) func(t *testing.T) {
	return func(t *testing.T) {
		require.NotNil(t, ptr)
		require.Equal(t, expected, *ptr)
	}
}

func Test_keepsDefaults(t *testing.T) {
	c := Config{
		MongoURI: "default",
		Env:      "default",
	}
	require.NoError(t, LoadConfig(&c))
	require.Equal(t, "default", c.MongoURI)
	require.Equal(t, "default", c.Env)
}

func Test_canLoadConfig(t *testing.T) {
	c := Config{
		MongoURI: "default",
		Env:      "default",
	}

	cases := []struct {
		env    string
		value  string
		expect func(t *testing.T)
	}{
		{
			env:    "MONGO_URI",
			value:  "user@mongo",
			expect: eq(&c.MongoURI, "user@mongo"),
		},
		{
			env:    "ENV",
			value:  "development",
			expect: eq(&c.Env, "development"),
		},
	}

	existingEnv := map[string]string{}
	// Set up environment variables first
	for _, tc := range cases {
		existingEnv[tc.env] = os.Getenv(tc.env)
		require.NoError(t, os.Setenv(tc.env, tc.value))
	}
	t.Cleanup(func() {
		for env, value := range existingEnv {
			os.Setenv(env, value)
		}
	})

	err := LoadConfig(&c)
	require.NoError(t, err)
	for _, tc := range cases {
		t.Run(tc.env, func(t *testing.T) {
			tc.expect(t)
		})
	}
}
