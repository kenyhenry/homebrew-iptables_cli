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

func Iptables_cmd(cmd IptablesCmd) string {
	// TODO us iptables instead
	base := "iptables"
	// TODO sperate option by place in command iptables
	option := " -p " + cmd.protocol + " --" + cmd.direction + cmd.port + " -m " + cmd.module + " --cstate " + cmd.connectionState + " -j " + cmd.jump
	full := exec.Command(base, option)
	out, err := full.Output()
	if err != nil {
		fmt.Println(err)
		return err.Error()
	}

	fmt.Println(string(out))
	return string(out)
}
