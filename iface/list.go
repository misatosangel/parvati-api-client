package iface

import (
	"fmt"
	"github.com/misatosangel/traceroute"
	"io"
	"net"
	"strings"
)

type InterfaceList struct {
	List   []net.Interface
	Filter int
}

func NewList(want int) (*InterfaceList, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	return &InterfaceList{List: ifaces, Filter: want}, nil
}

func (self *InterfaceList) GetInterfaceNumber(name string) int {
	for _, iface := range self.List {
		if strings.EqualFold(iface.Name, name) {
			return iface.Index
		}
	}
	return 0
}

func (self *InterfaceList) GetPublicIP(ifaceNum, filter int) (*traceroute.IPAddrMap, error) {
	var found *traceroute.IPAddrMap
	total := 0
	if filter == 0 {
		filter = self.Filter
	}
	for _, iface := range self.List {
		if ifaceNum != 0 && ifaceNum != iface.Index {
			continue
		}
		addrs, err := traceroute.FilterInterfaceIPs(iface, filter)
		if err != nil {
			return nil, err
		}
		cnt := len(addrs)
		realAddrs := make([]traceroute.IPAddrMap, 0, cnt)
		for _, a := range addrs {
			if a.RemoteIP != nil {
				realAddrs = append(realAddrs, a)
			}
		}
		cnt = len(realAddrs)
		if ifaceNum != 0 {
			switch cnt {
			case 0:
				return nil, fmt.Errorf("Interface %d. (%s) did not contain any public IPs.\n", iface.Index, iface.Name)
			case 1:
				val := realAddrs[0]
				return &val, nil
			default:
				ips := make([]string, cnt, cnt)
				for i, a := range realAddrs {
					ips[i] = "'" + a.RemoteIP.String() + "'"
				}
				return nil, fmt.Errorf("Interface %d. (%s) contains multiple public IPs: %s.\n", iface.Index, iface.Name, strings.Join(ips, ", "))
			}
		}
		total += cnt
		if cnt == 0 {
			continue
		}
		found = &(realAddrs[0])
	}
	if ifaceNum != 0 {
		ifaces := ""
		for _, iface := range self.List {
			ifaces += " - " + InterfaceToPublicIPString(iface, filter, true)
		}
		if ifaces == "" {
			ifaces = "<none available>"
		}
		vStr := traceroute.FilterToString(filter)
		return nil, fmt.Errorf("No suitable %s IPs were found in interface %d. Known interfaces were:\n"+ifaces, vStr, ifaceNum)
	}

	if total == 0 {
		vStr := traceroute.FilterToString(filter)
		return nil, fmt.Errorf("No suitable %s IPs were found in any interface. Are you connected to the internet?\n", vStr)
	}
	if total == 1 {
		return found, nil
	}
	ifaces := ""
	for _, iface := range self.List {
		ifaces += " - " + InterfaceToPublicIPString(iface, filter, true)
	}
	if ifaces == "" { // should never happen unless we get a disconnect between above and here on recheck
		vStr := traceroute.FilterToString(filter)
		return nil, fmt.Errorf("No suitable %s IPs were found in any interface. Are you connected to the internet?\n", vStr)
	}
	vStr := traceroute.FilterToString(filter)
	return nil, fmt.Errorf("Found more than one %s IP. Please specify which interface you wish to use.\n"+ifaces, vStr)
}

func InterfaceToPublicIPString(iface net.Interface, filter int, terse bool) string {
	addrs, err := traceroute.FilterInterfaceIPs(iface, filter|traceroute.WANT_LIVE_IP)
	vStr := traceroute.FilterToString(filter)
	out := fmt.Sprintf("%d. (%s) ", iface.Index, iface.Name)
	if err != nil {
		return out + err.Error() + "\n"
	}
	cnt := len(addrs)
	switch cnt {
	case 0:
		if terse {
			return out + "<none>\n"
		}
		if vStr != "any" {
			vStr = "no " + vStr
		} else {
			vStr = "no"
		}
		return out + vStr + " IPs.\n"
	case 1:
		if terse {
			return out + addrs[0].IPNatStr() + "\n"
		}
		if vStr == "any" {
			vStr = " IP"
		} else {
			vStr += " IP"
		}
		return out + "has " + vStr + ": " + addrs[0].IPNatStr() + "\n"
	}
	ips := make([]string, cnt, cnt)
	for i, a := range addrs {
		ips[i] = "'" + a.IPNatStr() + "'"
	}
	if terse {
		return out + strings.Join(ips, ", ") + "\n"
	}

	if vStr == "any" {
		vStr = " IPs"
	} else {
		vStr += " IPs"
	}
	return out + "contains multiple " + vStr + ": " + strings.Join(ips, ", ") + ".\n"
}

func (self *InterfaceList) GetFilteredList(filter int) ([]traceroute.IPAddrMap, error) {
	out := make([]traceroute.IPAddrMap, 0, 10)
	if filter == 0 {
		filter = self.Filter
	}
	for _, iface := range self.List {
		addrs, err := traceroute.FilterInterfaceIPs(iface, filter)
		if err != nil {
			return nil, err
		}
		if len(addrs) == 0 {
			continue
		}
		out = append(out, addrs...)
	}
	return out, nil
}

func (self *InterfaceList) Show(w io.Writer) {
	for _, iface := range self.List {
		addrs, err := traceroute.FilterInterfaceIPs(iface, self.Filter)
		if err == nil && len(addrs) == 0 {
			continue
		}
		fmt.Fprintf(w, "%d. %s", iface.Index, iface.Name)
		macStr := iface.HardwareAddr.String()
		if macStr != "" {
			fmt.Fprintf(w, " [MAC: %s]", macStr)
		}
		fmt.Fprint(w, "\n")
		if err != nil {
			fmt.Fprintf(w, "Could not find addreses: %s\n", err.Error())
			continue
		}
		for _, addr := range addrs {
			privIP := addr.LocalIP
			pubIP := addr.RemoteIP
			if addr.Error != nil {
				fmt.Fprintf(w, " - %s (%s)\n", privIP.String(), addr.Error.Error())
				continue
			}
			if addr.HasNAT() {
				var bestRoutes []traceroute.TraceRoute
				var routeErr error
				pubIpStr := pubIP.String()
				foundPublic := false
				lastLen := 0
				for hops := 1; hops > 0; hops++ {
					routes, err := traceroute.Trace(pubIpStr, "", iface.Name, hops, 500)
					if err != nil {
						routeErr = err
						break
					}
					curLen := len(routes)
					if curLen <= lastLen {
						break
					}
					lastLen = curLen
					bestRoutes = routes

					if pubIP.Equal(routes[curLen-1].IP) {
						foundPublic = true
						break
					}
					// last three are timeouts? give up
					if curLen > 3 && routes[curLen-1].IP == nil && routes[curLen-2].IP == nil && routes[curLen-3].IP == nil {
						break
					}
				}
				fmt.Fprintf(w, " - %s", privIP.String())
				gatewayIP, err := traceroute.FindGateway(pubIpStr, "", iface.Name, privIP)
				if err != nil {
					fmt.Fprintf(w, " [Gateway IP detection failed: '%s']", err.Error())
				} else {
					fmt.Fprintf(w, " [Router LAN IP: %s]", gatewayIP.String())
				}
				if bestRoutes != nil {
					for _, route := range bestRoutes {
						if route.Time == 0 {
							fmt.Fprintf(w, " --NAT [timed out]--> ???")
						} else {
							fmt.Fprintf(w, " --NAT [%0.3f ms]--> %s", route.Time, route.IP.String())
						}
					}
					if !foundPublic {
						fmt.Fprintf(w, " --???--> %s", pubIpStr)
					}
				} else {
					if routeErr == nil {
						fmt.Fprintf(w, " --NAT--> %s [No info from traceroute]", pubIpStr)
					} else {
						fmt.Fprintf(w, " --NAT--> %s [Error reading route: %s]", pubIpStr, routeErr.Error())
					}
				}
				fmt.Fprintf(w, "\n")
			} else {
				fmt.Fprintf(w, " - %s (no-NAT)\n", privIP.String())
			}
		}
		fmt.Fprint(w, "\n")
	}
}
