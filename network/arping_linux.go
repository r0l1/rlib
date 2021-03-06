package network

import (
    "net"
    "strings"
    "strconv"
    "github.com/r3boot/rlib/sys"
)

/*
 * Send count ARP Request packet(s) to ipaddr using arping. Return true if
 * the return code of arping is zero, false otherwise.
 */
func Arping(ipaddr net.IP, intf_name string, count int) (up bool, latency float64, err error) {
    if ipaddr == nil {
        return
    }

    arping, err := sys.BinaryPrefix("arping")
    if err != nil {
        return
    }

    stdout, _, err := sys.Run(arping, "-I", intf_name, "-c", strconv.Itoa(count), "-w", "3", ipaddr.String())

    var tot_latency float64 = 0

    if err == nil {
        up = true

        for _, line := range stdout {
            if ! strings.HasPrefix(line, "Unicast reply from") {
                continue
            }

            raw_latency := strings.Replace(strings.Split(line, " ")[6], "ms", "", -1)
            l, err := strconv.ParseFloat(raw_latency, 64)
            if err != nil {
                continue
            }

            tot_latency += l
        }
    }

    latency = (tot_latency / float64(count)) / 1000

    return
}
