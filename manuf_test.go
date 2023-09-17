package gomanuf

import "testing"

func TestManufactureSource(t *testing.T) {

	f := Search("28:23:f5:a4:3f:98")
	if f != nil {
		println(f.mac.String(), f.Name, f.FactureName)
	}

}

func TestTrans(t *testing.T) {

	transManufactures()

}
