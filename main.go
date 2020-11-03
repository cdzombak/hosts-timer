package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/txn2/txeh"
)

const (
	usageExit = 1
	hostsErrExit = 2

	blockedAddrIPv4 = "127.0.0.1"
)

func usage() {
	fmt.Printf("Usage: %s [OPTIONS] DOMAIN\n", os.Args[0])
	fmt.Printf("Options:\n")
	flag.PrintDefaults()
	fmt.Printf("\nIssues:\n  https://github.com/cdzombak/hosts-timer/issues/new\n")
	fmt.Printf("\nAuthor: Chris Dzombak <https://www.dzombak.com>\n")
}

func main() {
	flag_disable := flag.Bool("disable", false, "Disable access to the domain. " +
		"(Cannot be used with -enable or -time.)")
	flag_install := flag.Bool("install", false, "alias for -disable")
	flag_enable := flag.Bool("enable", false, "(Re-)Enable access to the domain, indefinitely. " +
		"(Cannot be used with -disable or -time).")
	flag_uninstall := flag.Bool("uninstall", false, "alias for -enable")
	flag_time := flag.String("time", "", "Enable access to the domain for the given amount of time " +
		"(string like 1h5m30s). Access is disabled after that time, or when the process receives SIGINT/SIGTERM " +
		"(eg. Ctrl-C). Cannot be used with -enable or -disable.")
	flag.Usage = usage
	flag.Parse()
	if len(flag.Args()) != 1 {
		flag.Usage()
		os.Exit(usageExit)
	}

	// TODO(cdzombak): accept list of domains on cli
	domain := flag.Args()[0]
	domain = strings.ToLower(domain)
	if strings.HasPrefix(domain, "www.") {
		domain = domain[4:]
	}
	if domain == "" {
		flag.Usage()
		os.Exit(usageExit)
	}

	hosts, err := txeh.NewHostsDefault()
	if err != nil {
		fmt.Println(err)
		os.Exit(hostsErrExit)
	}

	if *flag_disable || *flag_install {
		if *flag_enable || *flag_uninstall {
			flag.Usage()
			os.Exit(usageExit)
		}
		disableDomain(hosts, domain)
	} else if *flag_enable || *flag_uninstall {
		if *flag_disable || *flag_install {
			flag.Usage()
			os.Exit(usageExit)
		}
		enableDomain(hosts, domain)
	} else if *flag_time != "" {
		if *flag_disable || *flag_install || *flag_enable || *flag_uninstall {
			flag.Usage()
			os.Exit(usageExit)
		}
		duration, err := time.ParseDuration(*flag_time)
		if err != nil {
			fmt.Println(err)
			os.Exit(usageExit)
		}
		timedEnable(hosts, domain, duration)
	} else {
		flag.Usage()
		os.Exit(usageExit)
	}
}

func disableDomain(hosts *txeh.Hosts, domain string) {
	hosts.AddHosts(blockedAddrIPv4, []string{domain, "www."+domain})
	if err := hosts.Save(); err != nil {
		fmt.Println(err)
		os.Exit(hostsErrExit)
	}
}

func enableDomain(hosts *txeh.Hosts, domain string) {
	hosts.RemoveHosts([]string{domain, "www."+domain})
	if err := hosts.Save(); err != nil {
		fmt.Println(err)
		os.Exit(hostsErrExit)
	}
}

func timedEnable(hosts *txeh.Hosts, domain string, duration time.Duration) {
	enableDomain(hosts, domain)

	finishUpMutex := sync.Mutex{}
	finishUp := func() {
		finishUpMutex.Lock()
		if err := hosts.Reload(); err != nil {
			fmt.Println(err)
			os.Exit(hostsErrExit)
		}
		disableDomain(hosts, domain)
		os.Exit(0)
	}

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		finishUp()
	}()

	time.Sleep(duration)
	finishUp()
}
