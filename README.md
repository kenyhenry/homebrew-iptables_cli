# iptables Cli
Iptables cli is a Cli interface for iptables.
Iptables must be a good firewall in Linux, but seem a bit complicated in use.
So let’s try to make it more simple using this new open source tool name iptables cli write in Go.

## Install
``` brew tap kenyhenry/iptables_cli ```
``` brew install iptables_cli ```

## Why ?
- Iptables give us a wide choice to manage rules, move rules up and down is damagable.
- iptables -L takes time
- There is a lot of options and no categories to manage
- Delete wrong rule with index can lead to strong damaged on your VPS

## Precaution
- -t option not implemented yet because there is no way to get tables using iptables -L
- Move and edit rule included adding a new rule then deleting the old one so be careful on which rule you are moving
- Sometimes Rendering widgets is not automatic so tape escape to rerender
- In some, Iptables cli is challenging iptables use and is a better approach to manage all your rules.

## Future
- It could be interesting to switch to nftable on iptables cli
- Create our own firewall using netfilter and map it to our custom iptables cli interface


Under mit license

All contributors is welcome here
