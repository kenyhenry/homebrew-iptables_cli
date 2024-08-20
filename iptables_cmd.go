package main

import (
	"os/exec"
	"strconv"
	"strings"
)

type IptablesCmd struct {
	Chain           string
	Table           string
	Protocol        string
	SPort           string
	DPort           string
	Source          string
	Destination     string
	Module          string
	ModuleArg       string
	ConnectionState string
	Jump            string
	JumpArg         string
	InIface         string
	OutIface        string
}

func ContainString(target string, substring []string) bool {
	for _, str := range substring {
		if strings.Contains(target, str) {
			return true
		}
	}
	return false
}

func generateIptablesArgs(cmd IptablesCmd) []string {
	var args []string
	if cmd.Chain != "" {
		args = append(args, "-A", cmd.Chain)
	}
	if cmd.Table != "" {
		args = append(args, "-t", cmd.Table)
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
	if cmd.DPort != "" {
		args = append(args, "--dport", cmd.DPort)
	}
	if cmd.SPort != "" {
		args = append(args, "--dport", cmd.SPort)
	}
	if cmd.Module != "" {
		args = append(args, "-m", cmd.Module)
	}
	if cmd.ModuleArg != "" {
		args = append(args, "-m", cmd.ModuleArg)
	}
	if cmd.ConnectionState != "" {
		args = append(args, "--ctstate", cmd.ConnectionState)
	}
	if cmd.Jump != "" {
		args = append(args, "-j", cmd.Jump)
	}
	if cmd.JumpArg != "" {
		args = append(args, strings.Split(cmd.JumpArg, " ")...)
	}
	if cmd.InIface != "" {
		args = append(args, "-i", cmd.InIface)
	}
	if cmd.OutIface != "" {
		args = append(args, "-o", cmd.OutIface)
	}
	return args
}

func iptablesCmd(option []string) (string, error) {
	base := "iptables"
	full := exec.Command(base, option...)
	// fmt.Println(full)
	out, err := full.CombinedOutput()
	return string(out), err
}

// Chain
func IptablesAddChain(chainName string) (string, error) {
	option := []string{"-N", chainName}
	return iptablesCmd(option)
}

func IptablesFlushChain(chainName string) (string, error) {
	option := []string{"-F", chainName}
	return iptablesCmd(option)
}

func IptablesDeleteChain(chainName string) (string, error) {
	option := []string{"-X", chainName}
	return iptablesCmd(option)
}

func IptablesRenameChain(chainName string, newChainName string) (string, error) {
	option := []string{"-E", chainName, newChainName}
	return iptablesCmd(option)
}

func IptablesListChain() ([]string, error) {
	option := []string{"-L"}
	result, err := iptablesCmd(option)
	ret := []string{}
	if len(result) > 0 {
		lines := strings.Split(result, "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "Chain") {
				ret = append(ret, line)
				// fmt.Printf(line + "\n")
			}
		}
	}
	return ret, err
}

// POLICY
func IptablesMapPolicy(chainName string, targetChainName string) (string, error) {
	option := []string{"-A", chainName, "-j", targetChainName}
	return iptablesCmd(option)
}

func IptablesSetPolicy(chainName string, policy string) (string, error) {
	option := []string{"-P", chainName, policy}
	return iptablesCmd(option)
}

// Rule
func IptablesList(chainName string) ([]string, string, error) {
	option := []string{"-L", chainName}
	substring := []string{"target", "prot", "opt", "source", "destination", "Chain"}
	result, err := iptablesCmd(option)
	var ret []string
	if err == nil {
		lines := strings.Split(result, "\n")
		for _, line := range lines {
			if !ContainString(line, substring) {
				ret = append(ret, line)
				// fmt.Printf(line + "\n")
			}
		}
	} else {
		ret = append(ret, result)
		ret = append(ret, err.Error())
	}
	return ret, result, err
}

func IptablesAddRule(cmd IptablesCmd) (string, error) {
	option := generateIptablesArgs(cmd)
	return iptablesCmd(option)
}

func IptablesInsertRule(pos int, cmd IptablesCmd) (string, error) {
	option := generateIptablesArgs(cmd)
	return iptablesCmd(option)
}

func IptablesGetRule(chainName string, ruleIndex int) (string, error) {
	option := []string{"-L", chainName}
	substring := []string{"target", "prot", "opt", "source", "destination", "Chain"}
	result, err := iptablesCmd(option)
	var rules []string
	lines := strings.Split(result, "\n")
	for _, line := range lines {
		if !ContainString(line, substring) {
			rules = append(rules, line)
		}
	}
	rule := ""
	if len(rules) > ruleIndex {
		rule = rules[ruleIndex]
	}
	return rule, err
}

func IptablesDeleteRule(chainName string, ruleIndex int) (string, error) {
	option := []string{"-D", chainName, strconv.Itoa(ruleIndex)}
	return iptablesCmd(option)
}
