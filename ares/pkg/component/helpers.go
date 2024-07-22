package component

import (
	"fmt"
	"sort"

	"bitbucket.org/unchain/ares/pkg/ares"

	"github.com/unchainio/pkg/iferr"
	"github.com/unchainio/pkg/semver"
)

func versionFileName(componentType ares.ComponentType, name, version, orgName string) string {
	return fmt.Sprintf("%s.%s.%s.%s.so", name, componentType, version, orgName)
}

func versionFileID(componentType ares.ComponentType, name, version, orgName, fileName string) string {
	return fmt.Sprintf("%s/%s/%s/%s/%s.tar.gz", orgName, componentType, name, version, fileName)
}

func ActionVersionFileName(name, version, orgName string) string {
	return versionFileName(ares.ComponentTypeAction, name, version, orgName)
}

func ActionVersionFileID(name, version, orgName, fileName string) string {
	return versionFileID(ares.ComponentTypeAction, name, version, orgName, fileName)
}

func TriggerVersionFileName(name, version, orgName string) string {
	return versionFileName(ares.ComponentTypeTrigger, name, version, orgName)
}

func TriggerVersionFileID(name, version, orgName, fileName string) string {
	return versionFileID(ares.ComponentTypeTrigger, name, version, orgName, fileName)
}

func sortVersions(versions []string) {
	// Sort by semver
	sort.Slice(versions, func(i, j int) bool {
		iVer, err := semver.NewVersion(versions[i])
		if err != nil {
			iferr.Warn(err)
			return false
		}
		jVer, err := semver.NewVersion(versions[j])
		if err != nil {
			iferr.Warn(err)
			return false
		}

		return !iVer.LessThan(*jVer)
	})
}

// FIXME(e-nikolov) deduplicate due to https://github.com/volatiletech/sqlboiler/issues/457#issuecomment-532493095
func deduplicateVersions(versions []string) []string {
	// Fix issue with 'slice bounds out of range [:1] with capacity 0'
	if len(versions) == 0 {
		return versions
	}

	j := 0
	for i := 1; i < len(versions); i++ {
		if versions[j] == versions[i] {
			continue
		}
		j++
		versions[j] = versions[i]
	}

	return versions[:j+1]
}
