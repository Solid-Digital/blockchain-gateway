package casbin

import (
	"regexp"
	"strings"
)

// KeyMatchDomain determines whether key1 matches the pattern of key2 (similar to RESTful path), key2 can contain a *.
// For example, "/foo/bar" matches "/foo/*", "/resource1" matches "/:resource"
func KeyMatchDomain(reqObj string, polObj string) bool {
	//(parseDomain(r.obj, p.obj) + "::" + r.sub, p.sub) &&
	domain := ParseDomain(reqObj, polObj)
	//domain := strings.Split(subj, "::")[0]
	polObj = strings.Replace(polObj, ":domain", domain, -1)
	return KeyMatch2(reqObj, polObj)
}

// KeyMatchDomainFunc is the wrapper for KeyMatchDomain.
func KeyMatchDomainFunc(args ...interface{}) (interface{}, error) {
	name1 := args[0].(string)
	name2 := args[1].(string)

	return bool(KeyMatchDomain(name1, name2)), nil
}

func ParseDomain(reqObj string, polObj string) string {
	if !KeyMatch2(reqObj, polObj) {
		return "*"
	}

	pp := strings.Split(polObj, "/")
	rr := strings.Split(reqObj, "/")

	for i, s := range pp {
		if s == ":domain" {
			return rr[i]
		}
	}

	return "*"
}

// KeyMatchDomainFunc is the wrapper for KeyMatchDomain.
func ParseDomainFunc(args ...interface{}) (interface{}, error) {
	name1 := args[0].(string)
	name2 := args[1].(string)
	return string(ParseDomain(name1, name2)), nil
}

// KeyMatch2 determines whether key1 matches the pattern of key2 (similar to RESTful path), key2 can contain a *.
// For example, "/foo/bar" matches "/foo/*", "/resource1" matches "/:resource"
func KeyMatch2(key1 string, key2 string) bool {
	key2 = strings.Replace(key2, "/*", "/.*", -1)

	re := regexp.MustCompile(`(.*):[^/]+(.*)`)
	for {
		if !strings.Contains(key2, "/:") {
			break
		}

		key2 = "^" + re.ReplaceAllString(key2, "$1[^/]+$2") + "$"
	}

	return RegexMatch(key1, key2)
}

// KeyMatch2Func is the wrapper for KeyMatch2.
func KeyMatch2Func(args ...interface{}) (interface{}, error) {
	name1 := args[0].(string)
	name2 := args[1].(string)

	return bool(KeyMatch2(name1, name2)), nil
}

// KeyMatch3 determines whether key1 matches the pattern of key2 (similar to RESTful path), key2 can contain a *.
// For example, "/foo/bar" matches "/foo/*", "/resource1" matches "/{resource}"
func KeyMatch3(key1 string, key2 string) bool {
	key2 = strings.Replace(key2, "/*", "/.*", -1)

	re := regexp.MustCompile(`(.*)\{[^/]+\}(.*)`)
	for {
		if !strings.Contains(key2, "/{") {
			break
		}

		key2 = re.ReplaceAllString(key2, "$1[^/]+$2")
	}

	return RegexMatch(key1, key2)
}

// KeyMatch3Func is the wrapper for KeyMatch3.
func KeyMatch3Func(args ...interface{}) (interface{}, error) {
	name1 := args[0].(string)
	name2 := args[1].(string)

	return bool(KeyMatch3(name1, name2)), nil
}

// RegexMatch determines whether key1 matches the pattern of key2 in regular expression.
func RegexMatch(key1 string, key2 string) bool {
	res, err := regexp.MatchString(key2, key1)
	if err != nil {
		panic(err)
	}
	return res
}

// RegexMatchFunc is the wrapper for RegexMatch.
func RegexMatchFunc(args ...interface{}) (interface{}, error) {
	name1 := args[0].(string)
	name2 := args[1].(string)

	return bool(RegexMatch(name1, name2)), nil
}
