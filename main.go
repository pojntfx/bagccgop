package main

import (
	"fmt"
	"os/exec"
)

type Platform struct {
	GoOS             string
	GoArch           string
	DebianArch       string
	APTPackageSuffix string
	GCCArch          string
}

// Based on https://www.debian.org/ports/
// hppa, ia64, m68k and sh4 are not supported because they lack a gccgo package
var supportedPlatforms = []Platform{
	{
		"linux",
		"alpha",
		"alpha",
		"-alpha-linux-gnu",
		"alpha-linux-gnu",
	},
	{
		"linux",
		"ppc",
		"powerpc",
		"-powerpc-linux-gnu",
		"powerpc-linux-gnu",
	},
	{
		"linux",
		"ppc64",
		"ppc64",
		"-powerpc64-linux-gnu",
		"powerpc64-linux-gnu",
	},
	{
		"linux",
		"sparc64",
		"sparc64",
		"-sparc64-linux-gnu",
		"sparc64-linux-gnu",
	},
	{
		"linux",
		"riscv64",
		"riscv64",
		"-riscv64-linux-gnu",
		"riscv64-linux-gnu",
	},
	{
		"linux",
		"amd64",
		"amd64",
		"",
		"x86_64-linux-gnu",
	},
	{
		"linux",
		"arm64",
		"arm64",
		"-aarch64-linux-gnu",
		"aarch64-linux-gnu",
	},
	{
		"linux",
		"arm",
		"armel",
		"-arm-linux-gnueabi",
		"arm-linux-gnueabi",
	},
	{
		"linux",
		"arm",
		"armhf",
		"-arm-linux-gnueabihf",
		"arm-linux-gnueabihf",
	},
	{
		"linux",
		"386",
		"i386",
		"-i686-linux-gnu",
		"i686-linux-gnu",
	},
	{
		"linux",
		"mipsle",
		"mipsel",
		"-mipsel-linux-gnu",
		"mipsel-linux-gnu",
	},
	{
		"linux",
		"mips64le",
		"mips64el",
		"-mips64el-linux-gnuabi64",
		"mips64el-linux-gnuabi64",
	},
	{
		"linux",
		"powerpc64le",
		"ppc64el",
		"-powerpc64le-linux-gnu",
		"powerpc64le-linux-gnu",
	},
	{
		"linux",
		"s390x",
		"s390x",
		"-s390x-linux-gnu",
		"s390x-linux-gnu",
	},
}

func getSystemShell() []string {
	// Prefer Bash
	bash, err := exec.LookPath("bash")
	if err == nil {
		return []string{bash, "-c"}
	}

	// Fall back to POSIX shell
	return []string{"sh", "-c"}
}

func main() {
	// 	// Define usage
	// 	pflag.Usage = func() {
	// 		fmt.Printf(`Build for all gccgo-supported platforms by default, disable those which you don't want (bagop with CGo support).
	// Example usage: %s -b mybin -x '(linux/alpha|linux/ppc64el)' -j "$(nproc)" 'main.go'
	// Example usage (with plain flag): %s -b mybin -x '(linux/alpha|linux/ppc64el)' -j "$(nproc)" -p 'go build -o $DST main.go'
	// See https://github.com/pojntfx/bagccgop for more information.
	// Usage: %s [OPTION...] '<INPUT>'
	// `, os.Args[0], os.Args[0], os.Args[0])

	// 		pflag.PrintDefaults()
	// 	}

	// 	// Parse flags
	// 	binFlag := pflag.StringP("bin", "b", "mybin", "Prefix of resulting binary")
	// 	distFlag := pflag.StringP("dist", "d", "out", "Directory build into")
	// 	excludeFlag := pflag.StringP("exclude", "x", "", "Regex of platforms not to build for, i.e. (linux/alpha|linux/ppc64el)")
	// 	extraArgs := pflag.StringP("extra-args", "e", "", "Extra arguments to pass to the Go compiler")
	// 	jobsFlag := pflag.Int64P("jobs", "j", 1, "Maximum amount of parallel jobs")
	// 	goismsFlag := pflag.BoolP("goisms", "g", false, "Use Go's conventions (i.e. amd64) instead of uname's conventions (i.e. x86_64)")
	// 	plainFlag := pflag.BoolP("plain", "p", false, "Sets GOARCH, GOARCH, CC, GCCGO, GOFLAGS and DST and leaves the rest up to you (see example usage)")

	// 	pflag.Parse()
	for _, p := range supportedPlatforms {
		fmt.Print(p.DebianArch + " ")
	}

	fmt.Println()

	for _, p := range supportedPlatforms {
		fmt.Print("gccgo" + p.APTPackageSuffix + " gcc" + p.APTPackageSuffix + " ")
	}

	fmt.Println()
}
