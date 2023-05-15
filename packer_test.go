package main

import "testing"

func TestPackVerificationMessage(t *testing.T) {
	hexInput := []string{"00", "00", "00", "00", "00", "00", "01", "02", "03", "04", "05", "06", "07", "01", "02", "14", "31", "32", "33", "34", "35", "36", "37", "38", "00", "01", "41", "42", "43", "44", "45", "46", "47", "48", "49", "50", "02"}
	header := &Header{
		Length:    37,
		Seq:       1,
		Encrypted: false,
		FrameId:   "01",
	}
	expected := &VerificationMessage{
		Header:          header,
		Id:              "01020304050607",
		ElcType:         "AC",
		Guns:            2,
		ProtocolVersion: 2,
		SoftwareVersion: "12345678",
		Network:         "LAN",
		Sim:             "41424344454647484950",
		Operator:        "CT",
	}
	result := PackVerificationMessage(hexInput, header)
	if *result != *expected {
		t.Errorf("Expected %+v, got %+v", expected, result)
	}
}
