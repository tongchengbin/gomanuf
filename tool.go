package gomanuf

import (
	"bufio"
	"fmt"
	"gopkg.in/yaml.v3"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type CredentialSubject struct {
	Behaviour       string   `yaml:"behaviour"`
	ID              string   `yaml:"id"`
	MacAddresses    []string `yaml:"macAddresses"`
	Manufacturer    string   `yaml:"manufacturer"`
	ManufacturerUri string   `yaml:"manufacturerUri"`
	Name            string   `yaml:"name"`
}

type DeviceAssertion struct {
	CredentialSubject CredentialSubject `yaml:"credentialSubject"`
}

func parseMacNet(macCIDR string) *MACNet {
	fields := strings.Split(macCIDR, "/")
	if len(fields) != 2 {
		return nil
	}
	mask, err := strconv.Atoi(fields[1])
	if err != nil {
		return nil
	}
	mac, err := net.ParseMAC(fields[0])
	if err != nil {
		return nil
	}
	return &MACNet{MAC: mac, Mask: mask}
}

func transManufactures() {
	file, err := os.Create("manufacture.gob")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	err = filepath.Walk("../ManySecured-D3DB/", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			println(err.Error())
			return err
		}
		// 如果是文件，并且扩展名是 .yaml 或 .yml，处理文件
		if !info.IsDir() && (strings.HasSuffix(info.Name(), ".yaml") || strings.HasSuffix(info.Name(), ".yml")) {
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			dev := &DeviceAssertion{}
			err = yaml.Unmarshal(content, dev)
			if err != nil {
				return err
			}
			for _, macCIDR := range dev.CredentialSubject.MacAddresses {
				macNet := parseMacNet(macCIDR)
				if macNet == nil {
					continue
				}
				if dev.CredentialSubject.Name == "" && dev.CredentialSubject.Manufacturer == "" {
					continue
				}
				line := macNet.String() + "\t" + dev.CredentialSubject.Name + "\t" + dev.CredentialSubject.Manufacturer + "\n"
				//println(line)
				_, err = writer.WriteString(line)
				if err != nil {
					println(err.Error())
					return err
				}

			}
		}

		return nil
	})
	writer.Flush()
}
