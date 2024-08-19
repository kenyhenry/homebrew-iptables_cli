package main

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

type IptablesCmd struct {
	Chain           string
	Protocol        string
	Port            string
	Source          string
	Destination     string
	Module          string
	ConnectionState string
	Jump            string
}

// TODO : on sending every iptables cmd print return in a msgBox
// msgBox := MsgBox()

func containString(target string, substring []string) bool {
	for _, str := range substring {
		if strings.Contains(target, str) {
			return true
		}
	}
	return false
}

func generateIptablesArgs(cmd IptablesCmd) []string {
	var args []string

	// Check each field and append corresponding argument if not empty
	if cmd.Chain != "" {
		args = append(args, "-A", cmd.Chain)
	}
	if cmd.Protocol != "" {
		args = append(args, "-p", cmd.Protocol)
	}
	if cmd.Source != "" {
		args = append(args, "-s", cmd.Source)
	}
	if cmd.Destination != "" {
		args = append(args, "-d", cmd.Destination)
	}
	if cmd.Port != "" {
		args = append(args, "--dport", cmd.Port) // Assuming port is destination port
	}
	if cmd.Module != "" {
		args = append(args, "-m", cmd.Module)
	}
	if cmd.ConnectionState != "" {
		args = append(args, "--ctstate", cmd.ConnectionState)
	}
	if cmd.Jump != "" {
		args = append(args, "-j", cmd.Jump)
	}
	return args
}

func iptablesCmd(option []string) string {
	base := "iptables"
	full := exec.Command(base, option...)
	fmt.Println(full)
	out, err := full.CombinedOutput()
	return string(out) + err.Error()
}

// Chain
func IptablesAddChain(chainName string) string {
	option := []string{"-N", chainName}
	return iptablesCmd(option)
}

func IptablesFlushChain(chainName string) string {
	option := []string{"-F", chainName}
	return iptablesCmd(option)
}

func IptablesDeleteChain(chainName string) string {
	option := []string{"-X", chainName}
	return iptablesCmd(option)
}

func IptablesMapPolicy(chainName string, targetChainName string) string {
	option := []string{"-A", chainName, "-j", targetChainName}
	return iptablesCmd(option)
}

func IptablesSetPolicy(chainName string, policy string) string {
	option := []string{"-P", chainName, policy}
	return iptablesCmd(option)
}

func IptablesListChain() []string {
	option := []string{"-L"}
	result := iptablesCmd(option)
	var ret []string
	lines := strings.Split(result, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "Chain") {
			ret = append(ret, line)
			fmt.Printf(line + "\n")
		}
	}
	return ret
}

// Rule
func IptablesList(chainName string) []string {
	option := []string{"-L", chainName}
	substring := []string{"target", "prot", "opt", "source", "destination", "Chain"}
	result := iptablesCmd(option)
	var ret []string
	lines := strings.Split(result, "\n")
	for _, line := range lines {
		if !containString(line, substring) {
			ret = append(ret, line)
			fmt.Printf(line + "\n")
		}
	}
	return ret
}

func IptablesAddRule(cmd IptablesCmd) string {
	option := generateIptablesArgs(cmd)
	return iptablesCmd(option)
}

func IptablesGetRule(chainName string, ruleIndex int) string {
	option := []string{"-L", chainName}
	substring := []string{"target", "prot", "opt", "source", "destination", "Chain"}
	result := iptablesCmd(option)
	var rules []string
	lines := strings.Split(result, "\n")
	for _, line := range lines {
		if !containString(line, substring) {
			rules = append(rules, line)
		}
	}
	rule := ""
	if len(rules) > ruleIndex {
		rule = rules[ruleIndex]
	}
	return rule
}

func IptablesDeleteRule(chainName string, ruleIndex int) string {
	option := []string{"-D", chainName, strconv.Itoa(ruleIndex)}
	return iptablesCmd(option)
}
