package pipeline

import (
	"strconv"
	"strings"
)

func (s *Service) IncreaseSemver(version string) string {

	// extract patch from semver string
	versionArray := strings.Split(version, ".")
	// convert patch to integer
	patchString := versionArray[len(versionArray)-1]

	// error is ignored since versioning is controlled by backend
	patchInteger, _ := strconv.Atoi(patchString)

	// increase patch by +1, convert it back to string and replace it in version array
	versionArray[len(versionArray)-1] = strconv.Itoa(patchInteger + 1)

	// join array back to semver string
	return strings.Join(versionArray, ".")
}
