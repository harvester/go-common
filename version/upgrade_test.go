package version

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHarvesterUpgradeVersion_IsUpgrade(t *testing.T) {
	var testCases = []struct {
		name           string
		currentVersion string
		upgradeVersion string
		expectedErr    error
	}{
		{"upgrade from a newer release version to an older one", "v1.3.0", "v1.2.1", ErrDowngrade},
		{"upgrade from an older release version to a newer one", "v1.2.2", "v1.3.0", nil},

		{"upgrade from a newer prerelease version to an older one", "v1.3.0-rc1", "v1.2.1-rc1", ErrDowngrade},
		{"upgrade from a newer prerelease version to an older one", "v1.2.1-rc2", "v1.2.1-rc1", ErrDowngrade},
		{"upgrade from an older prerelease version to a newer one", "v1.2.1-rc1", "v1.3.0-rc1", nil},
		{"upgrade from an older prerelease version to a newer one", "v1.2.1-rc1", "v1.2.1-rc2", nil},

		{"upgrade among two dev versions", "11223344", "aabbccdd", ErrIncomparableVersion},

		{"upgrade from a newer prerelease version to an older release version", "v1.2.2-rc2", "v1.2.1", ErrDowngrade},
		{"upgrade from an older prerelease version to a newer release version", "v1.2.1-rc2", "v1.2.2", nil},
		{"upgrade from a newer release version to an older prerelease version", "v1.2.2", "v1.2.2-rc2", ErrDowngrade},
		{"upgrade from a newer release version to an older prerelease version", "v1.2.2", "v1.2.1-rc2", ErrDowngrade},
		{"upgrade from an older release version to a newer prerelease version", "v1.2.1", "v1.2.2-rc1", nil},

		{"upgrade from a release version to a dev version", "v1.2.1", "aabbccdd", ErrIncomparableVersion},
		{"upgrade from a dev version to a release version", "aabbccdd", "v1.2.1", ErrIncomparableVersion},
		{"upgrade from a prerelease version to a dev version", "v1.2.2-rc1", "aabbccdd", ErrIncomparableVersion},
		{"upgrade from a dev version to a prerelease version", "aabbccdd", "v1.2.2-rc1", ErrIncomparableVersion},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cv, _ := NewHarvesterVersion(tc.currentVersion)
			uv, _ := NewHarvesterVersion(tc.upgradeVersion)

			huv := NewHarvesterUpgradeVersion(cv, uv, nil)
			actualErr := huv.IsUpgrade()
			if tc.expectedErr != nil {
				assert.Equal(t, tc.expectedErr, actualErr, tc.name)
			} else {
				assert.Nil(t, actualErr, tc.name)
			}
		})
	}
}

func TestHarvesterUpgradeVersion_IsUpgradable(t *testing.T) {
	var testCases = []struct {
		name                 string
		currentVersion       string
		minUpgradableVersion string
		expectedErr          error
	}{
		{"upgrade from a release version above the minimal requirement", "v1.2.1", "v1.2.0", nil},
		{"upgrade from a release version lower than the minimal requirement", "v1.2.1", "v1.2.2", ErrMinUpgradeRequirement},
		{"upgrade from the exact same release version of the minimal requirement", "v1.2.1", "v1.2.1", nil},

		{"upgrade from a prerelease version lower than the minimal requirement", "v1.2.1-rc1", "v1.2.1", ErrMinUpgradeRequirement},
		{"upgrade from a prerelease version above the minimal requirement (rc minUpgradableVersion)", "v1.2.1-rc2", "v1.2.1-rc1", nil},
		{"upgrade from a prerelease version lower than the minimal requirement (rc minUpgradableVersion)", "v1.2.1-rc1", "v1.2.1-rc2", ErrMinUpgradeRequirement},
		{"upgrade from a prerelease version lower than the minimal requirement (rc minUpgradableVersion)", "v1.2.0-rc1", "v1.2.1-rc2", ErrMinUpgradeRequirement},
		{"upgrade from the exact same prerelease version of the minimal requirement (rc minUpgradableVersion)", "v1.2.1-rc1", "v1.2.1-rc1", nil},

		{"upgrade from dev versions", "11223344", "v1.2.1", ErrIncomparableVersion},
		{"upgrade from dev versions", "aabbccdd", "v1.2.1", ErrIncomparableVersion},

		{"upgrade from a release version without any minimum requirement specified", "v1.2.1", "", nil},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cv, _ := NewHarvesterVersion(tc.currentVersion)
			mv, _ := NewHarvesterVersion(tc.minUpgradableVersion)

			huv := NewHarvesterUpgradeVersion(cv, nil, mv)
			actualErr := huv.IsUpgradable()
			if tc.expectedErr != nil {
				assert.Equal(t, tc.expectedErr, actualErr, tc.name)
			} else {
				assert.Nil(t, actualErr, tc.name)
			}
		})
	}
}

func TestHarvesterUpgradeVersion_CheckUpgradeEligibility(t *testing.T) {
	var testCases = []struct {
		name                 string
		currentVersion       string
		upgradeVersion       string
		minUpgradableVersion string
		strictMode           bool
		expectedErr          error
	}{
		{"upgrade to same release version", "v1.2.1", "v1.2.1", "v1.1.2", true, nil},
		{"upgrade from an old release version same as the minimal requirement of a new release version", "v1.1.2", "v1.2.1", "v1.1.2", true, nil},
		{"upgrade from an old release version above the minimal requirement of a new release version", "v1.2.0", "v1.2.1", "v1.1.2", true, nil},
		{"upgrade from an old release version below the minimal requirement of a new release version", "v1.1.1", "v1.2.1", "v1.1.2", true, ErrMinUpgradeRequirement},
		{"upgrade from a new release version above the minimal requirement of an old release version", "v1.2.1", "v1.2.0", "v1.1.2", true, ErrDowngrade},
		{"upgrade from an old release version same the minimal requirement of a new prerelease version", "v1.1.2", "v1.2.1-rc1", "v1.1.2", true, nil},
		{"upgrade from an old release version above the minimal requirement of a new prerelease version", "v1.2.0", "v1.2.1-rc1", "v1.1.2", true, nil},
		{"upgrade from an old release version below the minimal requirement of a new prerelease version", "v1.1.1", "v1.2.1-rc1", "v1.1.2", true, ErrMinUpgradeRequirement},
		{"upgrade from a new release version above the minimal requirement of an old prerelease version", "v1.2.1", "v1.2.1-rc1", "v1.1.2", true, ErrDowngrade},
		{"upgrade from a release version to a dev version", "v1.2.1", "v1.2-ab12cd34", "", true, nil},

		{"upgrade to same prerelease version", "v1.2.2-rc1", "v1.2.2-rc1", "v1.2.1", true, nil},
		{"upgrade from an old prerelease version below the minimal requirement of a new release version", "v1.1.2-rc1", "v1.2.1", "v1.1.2", true, ErrMinUpgradeRequirement},
		{"upgrade from an old prerelease version above the minimal requirement of a new release version", "v1.2.0-rc1", "v1.2.1", "v1.1.2", true, ErrPrereleaseCrossVersionUpgrade},
		{"upgrade from an old prerelease version below the minimal requirement of a new release version", "v1.1.1-rc1", "v1.2.1", "v1.1.2", true, ErrMinUpgradeRequirement},
		{"upgrade from an old prerelease version above the minimal requirement of a new prerelease version", "v1.2.2-rc1", "v1.2.2-rc2", "v1.2.1", true, nil},
		{"upgrade from an old prerelease version below the minimal requirement of a new prerelease version", "v1.2.1-rc1", "v1.2.2-rc2", "v1.2.1", true, ErrMinUpgradeRequirement},
		{"upgrade from an old prerelease version above the minimal requirement of a new prerelease version", "v1.2.0-rc1", "v1.2.1-rc2", "v1.1.2", true, ErrPrereleaseCrossVersionUpgrade},
		{"upgrade from a prerelease version to a dev version", "v1.2.2-rc1", "v1.2-ab12cd34", "", true, nil},

		{"upgrade to same dev version", "v1.2-ab12cd34", "v1.2-ab12cd34", "", true, nil},
		{"upgrade among two dev versions", "v1.2-ab12cd34", "v1.2-1234567", "", true, nil},
		{"upgrade from a dev version to a release version (strict mode)", "v1.2-ab12cd34", "v1.2.1", "v1.1.2", true, ErrDevUpgrade},
		{"upgrade from a dev version to a release version (loose mode)", "v1.2-ab12cd34", "v1.2.1", "v1.1.2", false, nil},
		{"upgrade from a dev version to a prerelease version (strict mode)", "v1.2-ab12cd34", "v1.2.2-rc1", "v1.2.1", true, ErrDevUpgrade},
		{"upgrade from a dev version to a prerelease version (loose mode)", "v1.2-ab12cd34", "v1.2.2-rc1", "v1.2.1", false, nil},

		{"upgrade from an old release version to a new release version without any minimum requirement specified", "v1.2.1", "v1.3.0", "", true, nil},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cv, _ := NewHarvesterVersion(tc.currentVersion)
			uv, _ := NewHarvesterVersion(tc.upgradeVersion)
			mv, _ := NewHarvesterVersion(tc.minUpgradableVersion)

			huv := NewHarvesterUpgradeVersion(cv, uv, mv)
			actualErr := huv.CheckUpgradeEligibility(tc.strictMode)
			if tc.expectedErr != nil {
				assert.Equal(t, tc.expectedErr, actualErr, tc.name)
			} else {
				assert.Nil(t, actualErr, tc.name)
			}
		})
	}
}
