package main

import (
	"testing"

	"github.com/pkg/profile"
	"github.com/stretchr/testify/require"
)

func TestProfilingOptions(t *testing.T) {
	const (
		profDir = "/tmp"
	)

	testCases := []struct {
		name     string
		profMode string
		profDir  string
		options  []func(*profile.Profile)
		err      error
	}{
		{
			name: "emtpy profile dir",
			err:  errEmptyProfDir,
		},
		{
			name:     "block",
			profDir:  profDir,
			profMode: "block",
			options: []func(*profile.Profile){
				profile.BlockProfile,
				profile.ProfilePath(profDir),
			},
		},
		{
			name:     "cpu",
			profDir:  profDir,
			profMode: "cpu",
			options: []func(*profile.Profile){
				profile.CPUProfile,
				profile.ProfilePath(profDir),
			},
		},
		{
			name:     "mem",
			profDir:  profDir,
			profMode: "mem",
			options: []func(*profile.Profile){
				profile.MemProfile,
				profile.ProfilePath(profDir),
			},
		},
		{
			name:     "mutex",
			profDir:  profDir,
			profMode: "mutex",
			options: []func(*profile.Profile){
				profile.MutexProfile,
				profile.ProfilePath(profDir),
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			opts, err := getProfOpts(tc.profMode, tc.profDir)
			if tc.err != nil {
				require.Error(t, err)
				require.Equal(t, tc.err, err)
				return
			}
			require.NoError(t, err)
			require.Len(t, opts, len(tc.options))
			for i, opt := range opts {
				require.IsType(t, tc.options[i], opt)
			}
		})
	}
}
