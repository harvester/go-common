package version

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewHarvesterVersion(t *testing.T) {
	var testCases = []struct {
		name         string
		versionStr   string
		isDev        bool
		isPrerelease bool
		expectedErr  error
	}{
		{"formal release version string", "v1.2.2", false, false, nil},
		{"prerelease version string", "v1.2.2-rc1", false, true, nil},
		{"dev version string", "f024f49a", true, false, nil},
		{"empty version string", "", false, false, ErrInvalidVersion},
	}

	for _, tc := range testCases {
		hv, err := NewHarvesterVersion(tc.versionStr)
		assert.Equal(t, tc.expectedErr, err, tc.name)
		if tc.expectedErr == nil {
			assert.Equal(t, tc.isDev, hv.isDev, tc.name)
			assert.Equal(t, tc.isPrerelease, hv.isPrerelease, tc.name)
		}
	}
}

func TestHarvesterVersion_IsNewer(t *testing.T) {
	var testCases = []struct {
		name        string
		version1    string
		version2    string
		isNewer     bool
		expectedErr error
	}{
		{"same formal release versions", "v1.2.1", "v1.2.1", false, nil},
		{"same prerelease versions", "v1.2.2-rc1", "v1.2.2-rc1", false, nil},

		{"same dev versions", "aabbccdd", "aabbccdd", false, ErrIncomparableVersion},
		{"same dev versions (dirty build)", "aabbccdd-dirty", "aabbccdd-dirty", false, ErrIncomparableVersion},

		{"compare between formal release versions", "v1.2.2", "v1.2.1", true, nil},
		{"compare between formal release versions", "v1.2.1", "v1.2.2", false, nil},

		{"compare between prerelease versions", "v1.2.2-rc1", "v1.2.1-rc1", true, nil},
		{"compare between prerelease versions", "v1.2.1-rc1", "v1.2.2-rc1", false, nil},
		{"compare between prerelease versions", "v1.2.2-rc1", "v1.2.1-rc2", true, nil},
		{"compare between prerelease versions", "v1.2.1-rc2", "v1.2.2-rc1", false, nil},
		{"compare between prerelease versions", "v1.2.1-rc2", "v1.2.2-rc1", false, nil},
		{"compare between prerelease versions", "v1.2.1-rc3", "v1.2.1-rc2", true, nil},
		{"compare between prerelease versions", "v1.2.1-rc2", "v1.2.1-rc3", false, nil},

		{"compare between dev versions", "aabbccdd", "11223344", false, ErrIncomparableVersion},
		{"compare between dev versions (dirty build)", "aabbccdd-dirty", "11223344-dirty", false, ErrIncomparableVersion},

		{"compare between formal release and prerelease versions", "v1.2.2-rc1", "v1.2.1", true, nil},
		{"compare between formal release and prerelease versions", "v1.2.1", "v1.2.2-rc1", false, nil},
		{"compare between formal release and prerelease versions", "v1.2.2", "v1.2.2-rc1", true, nil},
		{"compare between formal release and prerelease versions", "v1.2.2-rc1", "v1.2.2", false, nil},
		{"compare between formal release and prerelease versions", "v1.2.2-rc1", "v1.3.0", false, nil},
		{"compare between formal release and prerelease versions", "v1.10.0", "v1.2.0-rc10", true, nil},
		{"compare between formal release and prerelease versions", "v1.10.0-rc1", "v1.2.0", true, nil},

		{"compare between formal release and dev version", "v1.2.1", "11223344", false, ErrIncomparableVersion},
		{"compare between formal release and dev version", "11223344", "v1.2.1", false, ErrIncomparableVersion},
		{"compare between prerelease and dev version", "v1.2.1-rc1", "11223344", false, ErrIncomparableVersion},
		{"compare between prerelease and dev version", "11223344", "v1.2.1-rc1", false, ErrIncomparableVersion},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			v1, err := NewHarvesterVersion(tc.version1)
			assert.Nil(t, err, tc.name)

			v2, err := NewHarvesterVersion(tc.version2)
			assert.Nil(t, err, tc.name)

			isNewer, err := v1.IsNewer(v2)
			if tc.expectedErr != nil {
				assert.Equal(t, tc.expectedErr, err, tc.name)
			} else {
				assert.Nil(t, err, tc.name)
			}
			assert.Equal(t, tc.isNewer, isNewer, tc.name)
		})
	}
}

func TestHarvesterVersion_IsEqual(t *testing.T) {
	var testCases = []struct {
		name        string
		version1    string
		version2    string
		isEqual     bool
		expectedErr error
	}{
		{"same formal release versions", "v1.2.1", "v1.2.1", true, nil},
		{"same prerelease versions", "v1.2.2-rc1", "v1.2.2-rc1", true, nil},

		{"same dev versions", "aabbccdd", "aabbccdd", true, nil},
		{"same dev versions (dirty build)", "aabbccdd-dirty", "aabbccdd-dirty", true, nil},

		{"compare between formal release versions", "v1.2.2", "v1.2.1", false, nil},
		{"compare between formal release versions", "v1.2.1", "v1.2.2", false, nil},

		{"compare between prerelease versions", "v1.2.2-rc1", "v1.2.1-rc1", false, nil},
		{"compare between prerelease versions", "v1.2.1-rc1", "v1.2.2-rc1", false, nil},
		{"compare between prerelease versions", "v1.2.2-rc1", "v1.2.1-rc2", false, nil},
		{"compare between prerelease versions", "v1.2.1-rc2", "v1.2.2-rc1", false, nil},
		{"compare between prerelease versions", "v1.2.1-rc2", "v1.2.2-rc1", false, nil},

		{"compare between dev versions", "aabbccdd", "11223344", false, ErrIncomparableVersion},
		{"compare between dev versions (dirty build)", "aabbccdd-dirty", "11223344-dirty", false, ErrIncomparableVersion},

		{"compare between formal release and prerelease versions", "v1.2.2-rc1", "v1.2.1", false, nil},
		{"compare between formal release and prerelease versions", "v1.2.1", "v1.2.2-rc1", false, nil},
		{"compare between formal release and prerelease versions", "v1.2.2", "v1.2.2-rc1", false, nil},
		{"compare between formal release and prerelease versions", "v1.2.2-rc1", "v1.2.2", false, nil},
		{"compare between formal release and prerelease versions", "v1.2.2-rc1", "v1.3.0", false, nil},
		{"compare between formal release and prerelease versions", "v1.10.0", "v1.2.0-rc10", false, nil},
		{"compare between formal release and prerelease versions", "v1.10.0-rc1", "v1.2.0", false, nil},

		{"compare between formal release and dev version", "v1.2.1", "11223344", false, ErrIncomparableVersion},
		{"compare between formal release and dev version", "11223344", "v1.2.1", false, ErrIncomparableVersion},
		{"compare between prerelease and dev version", "v1.2.1-rc1", "11223344", false, ErrIncomparableVersion},
		{"compare between prerelease and dev version", "11223344", "v1.2.1-rc1", false, ErrIncomparableVersion},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			v1, err := NewHarvesterVersion(tc.version1)
			assert.Nil(t, err, tc.name)

			v2, err := NewHarvesterVersion(tc.version2)
			assert.Nil(t, err, tc.name)

			isEqual, err := v1.IsEqual(v2)
			if tc.expectedErr != nil {
				assert.Equal(t, tc.expectedErr, err, tc.name)
			} else {
				assert.Nil(t, err, tc.name)
			}
			assert.Equal(t, tc.isEqual, isEqual, tc.name)
		})
	}
}

func TestHarvesterVersion_IsOlder(t *testing.T) {
	var testCases = []struct {
		name        string
		version1    string
		version2    string
		isOlder     bool
		expectedErr error
	}{
		{"same formal release versions", "v1.2.1", "v1.2.1", false, nil},
		{"same prerelease versions", "v1.2.2-rc1", "v1.2.2-rc1", false, nil},

		{"same dev versions", "aabbccdd", "aabbccdd", false, ErrIncomparableVersion},
		{"same dev versions (dirty build)", "aabbccdd-dirty", "aabbccdd-dirty", false, ErrIncomparableVersion},

		{"compare between formal release versions", "v1.2.2", "v1.2.1", false, nil},
		{"compare between formal release versions", "v1.2.1", "v1.2.2", true, nil},

		{"compare between prerelease versions", "v1.2.2-rc1", "v1.2.1-rc1", false, nil},
		{"compare between prerelease versions", "v1.2.1-rc1", "v1.2.2-rc1", true, nil},
		{"compare between prerelease versions", "v1.2.2-rc1", "v1.2.1-rc2", false, nil},
		{"compare between prerelease versions", "v1.2.1-rc2", "v1.2.2-rc1", true, nil},
		{"compare between prerelease versions", "v1.2.1-rc2", "v1.2.2-rc1", true, nil},

		{"compare between dev versions", "aabbccdd", "11223344", false, ErrIncomparableVersion},
		{"compare between dev versions (dirty build)", "aabbccdd-dirty", "11223344-dirty", false, ErrIncomparableVersion},

		{"compare between formal release and prerelease versions", "v1.2.2-rc1", "v1.2.1", false, nil},
		{"compare between formal release and prerelease versions", "v1.2.1", "v1.2.2-rc1", true, nil},
		{"compare between formal release and prerelease versions", "v1.2.2", "v1.2.2-rc1", false, nil},
		{"compare between formal release and prerelease versions", "v1.2.2-rc1", "v1.2.2", true, nil},
		{"compare between formal release and prerelease versions", "v1.2.2-rc1", "v1.3.0", true, nil},
		{"compare between formal release and prerelease versions", "v1.10.0", "v1.2.0-rc10", false, nil},
		{"compare between formal release and prerelease versions", "v1.10.0-rc1", "v1.2.0", false, nil},

		{"compare between formal release and dev version", "v1.2.1", "11223344", false, ErrIncomparableVersion},
		{"compare between formal release and dev version", "11223344", "v1.2.1", false, ErrIncomparableVersion},
		{"compare between prerelease and dev version", "v1.2.1-rc1", "11223344", false, ErrIncomparableVersion},
		{"compare between prerelease and dev version", "11223344", "v1.2.1-rc1", false, ErrIncomparableVersion},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			v1, err := NewHarvesterVersion(tc.version1)
			assert.Nil(t, err, tc.name)

			v2, err := NewHarvesterVersion(tc.version2)
			assert.Nil(t, err, tc.name)

			isOlder, err := v1.IsOlder(v2)
			if tc.expectedErr != nil {
				assert.Equal(t, tc.expectedErr, err, tc.name)
			} else {
				assert.Nil(t, err, tc.name)
			}
			assert.Equal(t, tc.isOlder, isOlder, tc.name)
		})
	}
}

func TestHarvesterVersion_String(t *testing.T) {
	var testCases = []struct {
		name   string
		input  string
		output string
	}{
		{"formal release version", "v1.2.1", "v1.2.1"},
		{"prerelease version", "v1.2.2-rc1", "v1.2.2-rc1"},
		{"dev version", "aabbccdd", "aabbccdd"},
		{"dev version (dirty build)", "aabbccdd-dirty", "aabbccdd-dirty"},

		{"formal release version without prefix", "1.2.1", "v1.2.1"},
		{"prerelease version", "1.2.2-rc1", "v1.2.2-rc1"},
		{"dev version", "aabbccdd", "aabbccdd"},
		{"dev version (dirty build)", "aabbccdd-dirty", "aabbccdd-dirty"},
		{"dev version", "v1.2-11223344", "v1.2-11223344"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			v, err := NewHarvesterVersion(tc.input)
			assert.Nil(t, err, tc.name)

			actual := v.String()
			assert.Equal(t, tc.output, actual, tc.name)
		})
	}
}
