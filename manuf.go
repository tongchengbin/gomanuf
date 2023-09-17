package gomanuf

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"path"
	"runtime"
	"strings"
)

type MACNet struct {
	MAC  net.HardwareAddr
	Mask int
}

func (mn *MACNet) Contains(macToCheck net.HardwareAddr) bool {
	maskBitCount := mn.Mask
	maskBytes := make(net.HardwareAddr, 6)
	for i := 0; i < maskBitCount/8; i++ {
		maskBytes[i] = 0xFF
	}
	maskBitsRemainder := maskBitCount % 8
	if maskBitsRemainder > 0 {
		maskBytes[maskBitCount/8] = byte(0xFF << (8 - maskBitsRemainder))
	}

	for i := 0; i < 6; i++ {
		if (macToCheck[i] & maskBytes[i]) != (mn.MAC[i] & maskBytes[i]) {
			if mn.MAC.String() == "28:23:f5:00:00:00/24" {
				println(">>>")
			}
			return false
		}
	}

	return true

}

type Manufacture struct {
	mac         *MACNet
	Name        string
	FactureName string
}

var manufacture []*Manufacture

func (mn *MACNet) String() string {
	return fmt.Sprintf("%s/%d", mn.MAC.String(), mn.Mask)
}

func init() {
	_, file, _, _ := runtime.Caller(0)
	f := path.Join(path.Dir(file), "manufacture.gob")
	fo, err := os.Open(f)
	if err != nil {
		return
	}
	defer func(fo *os.File) {
		_ = fo.Close()
	}(fo)
	scanner := bufio.NewScanner(fo)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Split(line, "\t")
		if len(fields) != 3 {
			continue
		}
		mac := fields[0]
		name := fields[1]
		factureName := fields[2]
		manufacture = append(manufacture, &Manufacture{Name: name, FactureName: factureName, mac: parseMacNet(mac)})
	}
}

func Search(mac string) *Manufacture {
	hw, err := net.ParseMAC(mac)
	if err != nil {
		println(err.Error())
		return nil
	}
	for _, manuf := range manufacture {
		if manuf.mac.Contains(hw) {
			return manuf
		}
	}
	return nil
}
