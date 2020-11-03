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
	fmt.Printf("Usage: %s [OPTIONS] DOMAIN [DOMAIN] ...\n", os.Args[0])
	fmt.Printf("Options:\n")
	flag.PrintDefaults()
	fmt.Printf("\nIssues:\n  https://github.com/cdzombak/hosts-timer/issues/new\n")
	fmt.Printf("\nAuthor: Chris Dzombak <https://www.dzombak.com>\n")
}

func main() {
	flag_disable := flag.Bool("disable", false, "Disable access to the domain(s). " +
		"(Cannot be used with -enable or -time.)")
	flag_install := flag.Bool("install", false, "alias for -disable")
	flag_enable := flag.Bool("enable", false, "(Re-)Enable access to the domain(s), indefinitely. " +
		"(Cannot be used with -disable or -time).")
	flag_uninstall := flag.Bool("uninstall", false, "alias for -enable")
	flag_time := flag.String("time", "", "Enable access to the domain(s) for the given amount of time " +
		"(string like 1h5m30s). Access is disabled after that time, or when the process receives SIGINT/SIGTERM " +
		"(eg. Ctrl-C). Cannot be used with -enable or -disable.")
	flag.Usage = usage
	flag.Parse()

	var domains []string
	for _, d := range flag.Args() {
		d = strings.ToLower(d)
		if strings.HasPrefix(d, "www.") {
			d = d[4:]
		}
		if d != "" {
			domains = append(domains, d)
			domains = append(domains, "www." + d)
		}
	}

	if len(domains) == 0 {
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
		disableDomains(hosts, domains)
	} else if *flag_enable || *flag_uninstall {
		if *flag_disable || *flag_install {
			flag.Usage()
			os.Exit(usageExit)
		}
		enableDomains(hosts, domains)
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
		timedEnable(hosts, domains, duration)
	} else {
		flag.Usage()
		os.Exit(usageExit)
	}
}

func disableDomains(hosts *txeh.Hosts, domains []string) {
	hosts.AddHosts(blockedAddrIPv4, domains)
	if err := hosts.Save(); err != nil {
		fmt.Println(err)
		os.Exit(hostsErrExit)
	}
}

func enableDomains(hosts *txeh.Hosts, domains []string) {
	hosts.RemoveHosts(domains)
	if err := hosts.Save(); err != nil {
		fmt.Println(err)
		os.Exit(hostsErrExit)
	}
}

func timedEnable(hosts *txeh.Hosts, domains []string, duration time.Duration) {
	enableDomains(hosts, domains)

	finishUpMutex := sync.Mutex{}
	finishUp := func() {
		finishUpMutex.Lock()
		if err := hosts.Reload(); err != nil {
			fmt.Println(err)
			os.Exit(hostsErrExit)
		}
		disableDomains(hosts, domains)
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
