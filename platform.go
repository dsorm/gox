package main

import (
	"fmt"
	"log"
	"strings"

	version "github.com/hashicorp/go-version"
)

// Platform is a combination of OS/arch that can be built against.
type Platform struct {
	OS   string
	Arch string

	// Default, if true, will be included as a default build target
	// if no OS/arch is specified. We try to only set as a default popular
	// targets or targets that are generally useful. For example, Android
	// is not a default because it is quite rare that you're cross-compiling
	// something to Android AND something like Linux.
	Default bool
}

func (p *Platform) String() string {
	return fmt.Sprintf("%s/%s", p.OS, p.Arch)
}

func removeElements(from []Platform, elements []Platform) []Platform {
	// make sure we don't change the passed underlying array
	np := make([]Platform, len(from), len(from))
	copy(np, from)

	for _, toRemove := range elements {
		for k, v := range np {
			if v.Arch == toRemove.Arch && v.OS == toRemove.OS {
				np = append(np[:k], np[k+1:]...)
				break
			}
		}
	}

	return np
}

var (
	Platforms_1_0 = []Platform{
		{"darwin", "386", true},
		{"darwin", "amd64", true},
		{"linux", "386", true},
		{"linux", "amd64", true},
		{"linux", "arm", true},
		{"freebsd", "386", true},
		{"freebsd", "amd64", true},
		{"openbsd", "386", true},
		{"openbsd", "amd64", true},
		{"windows", "386", true},
		{"windows", "amd64", true},
	}

	Platforms_1_1 = append(Platforms_1_0, []Platform{
		{"freebsd", "arm", true},
		{"netbsd", "386", true},
		{"netbsd", "amd64", true},
		{"netbsd", "arm", true},
		{"plan9", "386", false},
	}...)

	Platforms_1_3 = append(Platforms_1_1, []Platform{
		{"dragonfly", "386", false},
		{"dragonfly", "amd64", false},
		{"nacl", "amd64", false},
		{"nacl", "amd64p32", false},
		{"nacl", "arm", false},
		{"solaris", "amd64", false},
	}...)

	Platforms_1_4 = append(Platforms_1_3, []Platform{
		{"android", "arm", false},
		{"plan9", "amd64", false},
	}...)

	Platforms_1_5 = append(Platforms_1_4, []Platform{
		{"darwin", "arm", false},
		{"darwin", "arm64", true},
		{"linux", "arm64", true},
		{"linux", "ppc64", false},
		{"linux", "ppc64le", false},
	}...)

	Platforms_1_6 = append(Platforms_1_5, []Platform{
		{"android", "386", false},
		{"linux", "mips64", false},
		{"linux", "mips64le", false},
	}...)

	Platforms_1_7 = append(Platforms_1_5, []Platform{
		// While not fully supported s390x is generally useful
		{"linux", "s390x", true},
		{"plan9", "arm", false},
		// Add the 1.6 Platforms, but reflect full support for mips64 and mips64le
		{"android", "386", false},
		{"linux", "mips64", true},
		{"linux", "mips64le", true},
	}...)

	Platforms_1_8 = append(Platforms_1_7, []Platform{
		{"linux", "mips", true},
		{"linux", "mipsle", true},
	}...)

	// no new platforms in 1.9
	Platforms_1_9 = Platforms_1_8

	// no new platforms in 1.10
	Platforms_1_10 = Platforms_1_9

	Platforms_1_11 = append(Platforms_1_10, []Platform{
		{"js", "wasm", true},
	}...)

	// no new platforms in 1.12
	Platforms_1_12 = Platforms_1_11

	// no new platforms in 1.12
	Platforms_1_13 = Platforms_1_12

	Platforms_1_14 = removeElements(Platforms_1_13, []Platform{
		// Native Client was removed in 1.14 (https://golang.org/doc/go1.14#nacl)
		{"nacl", "amd64", false},
		{"nacl", "amd64p32", false},
		{"nacl", "arm", false},
	})

	Platforms_1_15 = append(
		removeElements(Platforms_1_14, []Platform{
			// darwin/386, darwin/arm is unsupported from Go 1.15 (https://golang.org/doc/go1.15#darwin)
			{"darwin", "386", true},
			{"darwin", "arm", false},
		}),
		[]Platform{
			{"linux", "riscv64", true},
		}...)

	Platforms_1_16 = append(Platforms_1_15,
		Platform{"ios", "amd64", false}, // iOS simulator on macOS devices with x86 CPU
		Platform{"ios", "arm64", false}, // regular iOS devices
	)

	Platforms_1_17 = append(Platforms_1_16,
		Platform{"windows", "arm64", true},
	)

	PlatformsLatest = Platforms_1_17
)

// SupportedPlatforms returns the full list of supported platforms for
// the version of Go that is
func SupportedPlatforms(v string) []Platform {
	// Use latest if we get an unexpected version string
	if !strings.HasPrefix(v, "go") {
		return PlatformsLatest
	}
	// go-version only cares about version numbers
	v = v[2:]

	current, err := version.NewVersion(v)
	if err != nil {
		log.Printf("Unable to parse current go version: %s\n%s", v, err.Error())

		// Default to latest
		return PlatformsLatest
	}

	var platforms = []struct {
		constraint string
		plat       []Platform
	}{
		{"<= 1.0", Platforms_1_0},
		{">= 1.1, < 1.3", Platforms_1_1},
		{">= 1.3, < 1.4", Platforms_1_3},
		{">= 1.4, < 1.5", Platforms_1_4},
		{">= 1.5, < 1.6", Platforms_1_5},
		{">= 1.6, < 1.7", Platforms_1_6},
		{">= 1.7, < 1.8", Platforms_1_7},
		{">= 1.8, < 1.9", Platforms_1_8},
		{">= 1.9, < 1.10", Platforms_1_9},
		{">=1.10, < 1.11", Platforms_1_10},
		{">=1.11, < 1.12", Platforms_1_11},
		{">=1.12, < 1.13", Platforms_1_12},
		{">=1.13, < 1.14", Platforms_1_13},
		{">=1.14, < 1.15", Platforms_1_14},
		{">=1.15, < 1.16", Platforms_1_15},
		{">=1.16, < 1.17", Platforms_1_16},
		{">=1.17, < 1.18", Platforms_1_17},
	}

	for _, p := range platforms {
		constraints, err := version.NewConstraint(p.constraint)
		if err != nil {
			panic(err)
		}
		if constraints.Check(current) {
			return p.plat
		}
	}

	// Assume latest
	return PlatformsLatest
}
