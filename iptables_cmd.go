package main

import (
	"fmt"
	"os/exec"
)

type IptablesCmd struct {
	protocol        string
	direction       string
	port            string
	module          string
	connectionState string
	jump            string
}

func iptablesCmd(option string) string {
	base := "iptables"
	full := exec.Command(base, option)
	out, err := full.Output()
	if err != nil {
		fmt.Println(err)
		return err.Error()
	}
	fmt.Println(string(out))
	return string(out)
}

func IptablesAddRule(cmd IptablesCmd) string {
	option := " -p " + cmd.protocol + " --" + cmd.direction + cmd.port + " -m " + cmd.module + " --cstate " + cmd.connectionState + " -j " + cmd.jump
	return iptablesCmd(option)
}

func IptablesList() string {
	option := " -L"
	return iptablesCmd(option)
}
