package main

import (
	"fmt"
	"os/exec"
	"strings"
)

type IptablesCmd struct {
	Protocol        string
	Direction       string
	Port            string
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
	if cmd.Protocol != "" {
		args = append(args, "-p", cmd.Protocol)
	}
	if cmd.Direction != "" {
		args = append(args, "-i", cmd.Direction) // Assuming direction corresponds to interface input
	}
	if cmd.Port != "" {
		args = append(args, "--dport", cmd.Port) // Assuming port is destination port
	}
	if cmd.Module != "" {
		args = append(args, "-m", cmd.Module)
	}
	if cmd.ConnectionState != "" {
		args = append(args, "-m", "conntrack", "--ctstate", cmd.ConnectionState)
	}
	if cmd.Jump != "" {
		args = append(args, "-j", cmd.Jump)
	}
	return args
}

func iptablesCmd(option string) string {
	base := "iptables"
	full := exec.Command(base, option)
	out, err := full.Output()
	if err != nil {
		fmt.Println(err)
		return err.Error()
	}
	// fmt.Println(string(out))
	return string(out)
}

func IptablesAddRule(cmd IptablesCmd) string {
	option := " -p " + cmd.Protocol + " --" + cmd.Direction + cmd.Port + " -m " + cmd.Module + " --cstate " + cmd.ConnectionState + " -j " + cmd.Jump
	return iptablesCmd(option)
}

func IptablesAddChain(chainName string) string {
	option := "-N" + chainName
	return iptablesCmd(option)
}

func IptablesListChain() []string {
	option := "-L"
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

func IptablesList(chainName string) []string {
	option := "-L" + chainName
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

func IptablesGetRule(chainName string, ruleIndex int) string {
	option := "-L" + chainName
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

func IptablesMoveRule(chainName string, ruleIndex int, targetIndex int) string {
	option := "-L" + chainName
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

func IptablesEditRule(chainName string, ruleIndex int, cmd IptablesCmd) string {
	option := "-R" + chainName + string(ruleIndex)
	result := iptablesCmd(option)
	return result
}

func IptablesDeleteRule(chainName string, ruleIndex int) string {
	option := "-D" + chainName + string(ruleIndex)
	result := iptablesCmd(option)
	return result
}
