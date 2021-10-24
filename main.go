package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"sync"

	"github.com/alessio/shellescape"
	"github.com/fatih/color"
	"github.com/spf13/pflag"
)

const (
	mountedPwd = "data"
)

func getChrootLocation(debianArch string) string {
	return path.Join("/var", "lib", "bagccgop", debianArch+"-chroot") // This is always a UNIX path hence `filepath` is not being used
}

func mountChroot(debianArch string) error {
	chrootLocation := getChrootLocation(debianArch)
	pwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("could not get home directory: %v", err)
	}

	for _, c := range [][]string{
		{"-o", "bind", "/dev", path.Join(chrootLocation, "dev")},
		{"-o", "bind", "/proc", path.Join(chrootLocation, "proc")},
		{"-o", "bind", pwd, path.Join(chrootLocation, mountedPwd)},
	} {
		cmd := exec.Command("mount", c...)

		// Capture stdout and stderr
		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		// Start the build
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("could not mount chroot: err=%v, stdout=%v, stderr=%v", err, stdout.String(), stderr.String())
		}
	}

	return nil
}

func execInChroot(debianArch string, cmds []string, env map[string]string) error {
	chrootLocation := getChrootLocation(debianArch)

	for _, c := range cmds {
		cmd := exec.Command("chroot", append([]string{chrootLocation, "/bin/bash", "-c"}, c)...)

		// Capture stdout and stderr
		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		// Set env vars
		cmd.Env = os.Environ()
		for key, value := range env {
			cmd.Env = append(cmd.Env, shellescape.Quote(key)+"="+value)
		}

		// Start the build
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("could not install packages: err=%v, stdout=%v, stderr=%v", err, stdout.String(), stderr.String())
		}
	}

	return nil
}

func getPkgNameForArch(pkg string, debianArch string) string {
	return pkg + ":" + debianArch
}

func getCC(gccArch string) string {
	return gccArch + "-gcc"
}

func getGCCGo(gccArch string) string {
	return gccArch + "-gccgo"
}

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
		"ppc64le",
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

func main() {
	// Define usage
	pflag.Usage = func() {
		fmt.Printf(`Build for all gccgo-supported platforms by default, disable those which you don't want (bagop with CGo support).
	Example usage: %s -b mybin -x '(linux/alpha|linux/ppc64el)' -j "$(nproc)" 'main.go'
	Example usage (with plain flag): %s -b mybin -x '(linux/alpha|linux/ppc64el)' -j "$(nproc)" -p 'go build -o $DST main.go'
	See https://github.com/pojntfx/bagccgop for more information.
	Usage: %s [OPTION...] '<INPUT>'
	`, os.Args[0], os.Args[0], os.Args[0])

		pflag.PrintDefaults()
	}

	// Parse flags
	binFlag := pflag.StringP("bin", "b", "mybin", "Prefix of resulting binary")
	distFlag := pflag.StringP("dist", "d", "out", "Directory build into")
	excludeFlag := pflag.StringP("exclude", "x", "", "Regex of platforms not to build for, i.e. (linux/alpha|linux/ppc64el)")
	extraArgs := pflag.StringP("extra-args", "e", "", "Extra arguments to pass to the Go compiler")
	jobsFlag := pflag.Int64P("jobs", "j", 1, "Maximum amount of parallel jobs")
	goismsFlag := pflag.BoolP("goisms", "g", false, "Use Go's conventions (i.e. amd64) instead of uname's conventions (i.e. x86_64)")
	plainFlag := pflag.BoolP("plain", "p", false, "Sets GOARCH, GOARCH, CC, GCCGO, GOFLAGS and DST and leaves the rest up to you (see example usage)")

	prepareCommandFlag := pflag.StringP("prepare", "r", "", "Command to run before running the main command; will have only CC and GCCGO set (i.e. for code generation)")
	hostPackagesFlag := pflag.StringSliceP("hostPackages", "s", []string{}, "Comma-seperated list of Debian packages to install for the host architecture")
	packagesFlag := pflag.StringSliceP("packages", "a", []string{}, "Comma-seperated list of Debian packages to install for the selected architectures")
	manualPackagesFlag := pflag.StringSliceP("manualPackages", "m", []string{}, "Comma-seperated list of Debian packages to manually install for the selected architectures (i.e. those which would break the dependency graph)")

	pflag.Parse()

	// Validate arguments
	if pflag.NArg() == 0 {
		help := `command needs an argument: 'INPUT'`

		fmt.Println(help)

		pflag.Usage()

		fmt.Println(help)

		os.Exit(2)
	}

	// Interpret arguments
	input := pflag.Args()[0]

	// Limits the max. amount of concurrent builds
	// See https://play.golang.org/p/othihEtsOBZ
	var wg = sync.WaitGroup{}
	guard := make(chan struct{}, *jobsFlag)

	for _, lplatform := range supportedPlatforms {
		guard <- struct{}{}
		wg.Add(1)

		go func(platform Platform) {
			defer func() {
				wg.Done()

				<-guard
			}()

			// Construct the filename
			output := filepath.Join(*distFlag, *binFlag+"."+platform.GoOS+"-")

			// Add the arch identifier
			archIdentifier := platform.DebianArch
			if *goismsFlag {
				archIdentifier = platform.GoArch
			}
			output += archIdentifier

			// Check if current platform should be skipped
			skip := false
			if *excludeFlag != "" {
				iskip, err := regexp.MatchString(*excludeFlag, platform.GoOS+"/"+platform.GoArch)
				if err != nil {
					log.Fatal("could not match check if platform should be blocked based on regex:", err)
				}

				skip = iskip
			}

			// Skip the platform if it matches the exclude regex
			if skip {
				log.Printf("%v %v/%v (platform matched the provided regex)", color.New(color.FgYellow).SprintFunc()("skipping"), color.New(color.FgCyan).SprintFunc()(platform.GoOS), color.New(color.FgMagenta).SprintFunc()(platform.GoArch))

				return
			}

			// Continue if platform is enabled
			log.Printf("%v %v/%v (%v)", color.New(color.FgGreen).SprintFunc()("building"), color.New(color.FgCyan).SprintFunc()(platform.GoOS), color.New(color.FgMagenta).SprintFunc()(platform.GoArch), output)

			// Mount chroot
			if err := mountChroot(platform.DebianArch); err != nil {
				log.Fatalf("could not mount chroot for platform %v/%v: err=%v", platform.GoOS, platform.GoArch, err)
			}

			// Install host packages
			for _, pkg := range *hostPackagesFlag {
				if err := execInChroot(
					platform.DebianArch,
					[]string{`apt install -y ` + shellescape.Quote(pkg)},
					nil,
				); err != nil {
					log.Fatalf("could not install host packages for platform %v/%v: err=%v", platform.GoOS, platform.GoArch, err)
				}
			}

			// Install packages
			for _, rawPkg := range *packagesFlag {
				pkg := getPkgNameForArch(rawPkg, platform.DebianArch)

				if err := execInChroot(
					platform.DebianArch,
					[]string{`apt install -y ` + shellescape.Quote(pkg)},
					nil,
				); err != nil {
					log.Fatalf("could not install packages for platform %v/%v: err=%v", platform.GoOS, platform.GoArch, err)
				}
			}

			// Install manual packages
			for _, rawPkg := range *manualPackagesFlag {
				pkg := getPkgNameForArch(rawPkg, platform.DebianArch)

				if err := execInChroot(
					platform.DebianArch,
					[]string{
						`mkdir -p /tmp/bagccgop-packages/` + shellescape.Quote(pkg),
						`cd /tmp/bagccgop-packages/` + shellescape.Quote(pkg) + ` && apt download ` + shellescape.Quote(pkg),
						`dpkg -i --force-all /tmp/bagccgop-packages/` + shellescape.Quote(pkg) + `/*.deb`,
					},
					nil,
				); err != nil {
					log.Fatalf("could not manually install packages for platform %v/%v: err=%v", platform.GoOS, platform.GoArch, err)
				}
			}

			// Run prepare command
			if *prepareCommandFlag != "" {
				if err := execInChroot(
					platform.DebianArch,
					[]string{"cd " + mountedPwd + " && " + *prepareCommandFlag},
					map[string]string{
						"CC":    getCC(platform.GCCArch),
						"GCCGO": getGCCGo(platform.GCCArch),
					},
				); err != nil {
					log.Fatalf("could not run prepare command for platform %v/%v: err=%v", platform.GoOS, platform.GoArch, err)
				}
			}

			// Construct build command
			buildLine := "go build -o " + output + " " + input
			if *extraArgs != "" {
				buildLine = "go build -o " + output + " " + *extraArgs + " " + input
			}

			// If the plain flag is set, use the custom command
			if *plainFlag {
				buildLine = input
			}

			// Set env vars
			buildEnv := map[string]string{
				"CC":          getCC(platform.GCCArch),
				"GCCGO":       getGCCGo(platform.GCCArch),
				"CGO_ENABLED": "1",
				"GOOS":        platform.GoOS,
				"GOARCH":      platform.GoArch,
				"GOFLAGS":     "-compiler=gccgo " + os.Getenv("GOFLAGS"),
			}

			// If the plain flag is set, also set DST
			if *plainFlag {
				buildEnv["DST"] = shellescape.Quote(output)
			}

			// Start the build
			if err := execInChroot(
				platform.DebianArch,
				[]string{"cd " + mountedPwd + " && " + buildLine},
				buildEnv,
			); err != nil {
				log.Fatalf("could not build for platform %v/%v: err=%v", platform.GoOS, platform.GoArch, err)
			}
		}(lplatform)
	}

	wg.Wait()
}
