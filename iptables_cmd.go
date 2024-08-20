package main

import (
	"os/exec"
	"regexp"
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

func getElement(str string, keyword string) string {
	start := strings.Index(str, keyword)
	if start != -1 {
		start += len(keyword)
		end := strings.Index(str[start:], " ")
		if end == -1 {
			end = len(str)
		} else {
			end += start
		}
		return str[start:end]
	}
	return ""
}

func ExtractAndGenerateCommands(line string, chaineName string) IptablesCmd {
	var cmd IptablesCmd

	parts := strings.Fields(line)

	protocol := parts[3]
	sourceIP := parts[7]
	destIP := parts[8]
	ifaceIn := parts[5]
	ifaceOut := parts[6]
	jump := parts[2]

	// TODO : handle -t table arg in iptables
	if protocol != "all" {
		cmd.Protocol = protocol
	}

	if sourceIP != "0.0.0.0/0" {
		cmd.Source = sourceIP
	}

	if destIP != "0.0.0.0/0" {
		cmd.Destination = destIP
	}

	if ifaceIn != "*" {
		cmd.InIface = ifaceIn
	}

	if ifaceOut != "*" {
		cmd.OutIface = ifaceOut
	}

	part := strings.Join(parts[9:], " ")

	// TODO : handle all modules
	if strings.HasPrefix(part, "icmptype") {
		icmp := strings.TrimPrefix(part, "icmptype ")
		cmd.Module = "icmp"
		cmd.ModuleArg = "--icmp-type " + icmp
	}

	if strings.HasPrefix(part, "dpt:") {
		port := getElement(part, "dpt:")
		cmd.DPort = port
	}
	if strings.HasPrefix(part, "spt:") {
		sport := getElement(part, "spt:")
		cmd.SPort = sport
	}
	if strings.HasPrefix(part, "ctstate ") {
		cstate := getElement(part, "ctstate ")
		cmd.ConnectionState = cstate
	}

	re := regexp.MustCompile(`"([^"]+)"`)
	if strings.HasPrefix(part, "LOG") {
		cmd.Jump = "LOG"
		var logPrefix string
		matches := re.FindAllStringSubmatch(part, -1)
		for _, match := range matches {
			logPrefix = match[1]
		}
		cmd.JumpArg = "--log-prefix " + "\"" + logPrefix + "\""
	} else {
		cmd.Jump = jump
	}

	return cmd
}

func ArraytToCmd(chain string, rules []string, base []int) IptablesCmd {
	return IptablesCmd{
		Chain:           chain,
		Table:           rules[0][base[0]:],
		Protocol:        rules[1][base[1]:],
		DPort:           rules[2][base[2]:],
		SPort:           rules[3][base[3]:],
		Source:          rules[4][base[4]:],
		Destination:     rules[5][base[5]:],
		Module:          rules[6][base[6]:],
		ModuleArg:       rules[7][base[7]:],
		ConnectionState: rules[8][base[8]:],
		Jump:            rules[9][base[9]:],
		JumpArg:         rules[10][base[10]:],
		InIface:         rules[11][base[11]:],
		OutIface:        rules[12][base[12]:],
	}
}

func generateIptablesArgs(cmd IptablesCmd) []string {
	var args []string
	if cmd.Chain != "" {
		args = append(args, "-A", cmd.Chain)
	}
	if cmd.Table != "" {
		args = append(args, "-t", cmd.Table)
	}
	if cmd.InIface != "" {
		args = append(args, "-i", cmd.InIface)
	}
	if cmd.OutIface != "" {
		args = append(args, "-o", cmd.OutIface)
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
		args = append(args, strings.Split(cmd.ModuleArg, " ")...)
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
	option := []string{"-L", chainName, "-n", "-v"}
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
