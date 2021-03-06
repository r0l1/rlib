package network

import (
    "errors"
    "io/ioutil"
    "net"
    "strconv"
    "github.com/r3boot/rlib/sys"
)

/*
 * Open /sys/class/net/<interface>/carrier and determine link status. Return
 * true if the content equals "1" (0x31), false otherwise. If the carrier file
 * cannot be read, return an error.
 */
func (l Link) HasCarrier () (result bool, err error) {
    carrier_file := "/sys/class/net/" + l.Interface.Name + "/carrier"

    content, err := ioutil.ReadFile(carrier_file)
    if err != nil {
        return
    }

    result = content[0] == LINK_UP
    return
}

func (l Link) SetLinkStatus (link_status byte) (err error) {
    var status string
    if link_status == LINK_UP {
        status = "up"
    } else if link_status == LINK_DOWN {
        status = "down"
    } else {
        err = errors.New("Unknown link status: " + strconv.Itoa(int(link_status)))
        return
    }

    _, _, err = sys.Run("/sbin/ip", "link", "set", l.Interface.Name, status)
    return
}

/*
 * Look in /sys/class/net/<interface>/type to see if this interface is
 * a loopback interface. Return if it is. Afterwards, look in
 * /sys/class/net/<interface>/device/class and check the pci class of the
 * device. If * this is "20000", it's an ethernet nic, if it's "28000", it's
 * a wireless nic. All other pci classes get flagged unknown.
 */
func (link *Link) GetType () (intf_type byte, err error) {
    var sys_file string

    l := *link
    if ! l.HasLink() {
        if err = l.SetLinkStatus(LINK_UP); err != nil {
            err = errors.New("SetLinkStatus failed: " + err.Error())
            return
        }
    }

    sys_file = "/sys/class/net/" + l.Interface.Name + "/type"
    content, err := ioutil.ReadFile(sys_file)
    if err != nil {
        return
    }

    value := string(content[0:3])
    if value == LINK_LOOPBACK {
        intf_type = INTF_TYPE_LOOPBACK
        return
    }

    sys_file = "/sys/class/net/" + l.Interface.Name + "/tun_flags"
    if sys.FileExists(sys_file) {
        content, err = ioutil.ReadFile(sys_file)
        if err != nil {
            return
        }

        value = string(content[2:6])
        if value == LINK_TAP {
            intf_type = INTF_TYPE_TAP
            return
        }
    }

    sys_file = "/sys/class/net/" + l.Interface.Name + "/device/class"
    if sys.FileExists(sys_file) {
        if content, err = ioutil.ReadFile(sys_file); err != nil {
            return
        }

        value = string(content[0:8])
        if value == LINK_WIRELESS {
            intf_type = INTF_TYPE_WIRELESS
            return
        } else if value == LINK_ETHERNET {
            intf_type = INTF_TYPE_ETHERNET
            return
        }
    }

    sys_file = "/sys/class/net/" + l.Interface.Name + "/bonding/mode"
    if sys.FileExists(sys_file) {
        if content, err = ioutil.ReadFile(sys_file); err != nil {
            return
        }

        value = string(content[0:13])
        if value == LINK_BONDING {
            intf_type = INTF_TYPE_BONDING
            return
        }
    }

    err = errors.New("Unknown interface type")

    return
}

func (l Link) GetMTU () (mtu int, err error) {
    mtu_file := "/sys/class/net/" + l.Interface.Name + "/mtu"

    content, err := ioutil.ReadFile(mtu_file)
    if err != nil { return }

    mtu, err  = strconv.Atoi(string(content[0:3]))

    return
}

func (l Link) SetMTU (mtu int) (err error) {
    cur_mtu, err := l.GetMTU()
    if err != nil {
        return
    }

    if cur_mtu != mtu {
        mtu_file := "/sys/class/net/" + l.Interface.Name + "/mtu"
        value := []byte(strconv.Itoa(mtu))
        err = ioutil.WriteFile(mtu_file, value, 0755)
    }

    return
}

func LinkFactory (intf net.Interface) (l Link, err error) {
    l = *new(Link)

    ifconfig, err := sys.BinaryPrefix("ip")
    if err != nil {
        return
    }

    l.Interface = intf
    l.CmdIfconfig = ifconfig

    return
}
