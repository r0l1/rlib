package sys

type Sysctl struct {
    CmdSysctl   string
}

func (s Sysctl) Set (key string, value string) (err error) {
    _, _, err = Run(s.CmdSysctl, "-w", key, "=", value)
    return
}

func SysctlFactory () (s Sysctl, err error) {
    s = *new(Sysctl)

    sysctl, err := BinaryPrefix("sysctl")
    if err != nil {
        return
    }

    s.CmdSysctl = sysctl
    return
}
