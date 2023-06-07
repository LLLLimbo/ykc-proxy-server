package main

import (
	"bytes"
	hex2 "encoding/hex"
	"fmt"
	"strconv"
)

const (
	//start flag
	StartFlag = byte(0x68)

	//device -> platform
	Verification                = byte(0x01)
	Heartbeat                   = byte(0x03)
	BillingModelVerification    = byte(0x05)
	BillingModelRequest         = byte(0x08)
	OfflineDataReport           = byte(0x13)
	ChargingHandshake           = byte(0x15)
	Configuration               = byte(0x17)
	ChargingFinished            = byte(0x19)
	ErrorReport                 = byte(0x1b)
	BmsInterrupted              = byte(0x1d)
	ChargingPileInterrupted     = byte(0x21)
	ChargingMetrics             = byte(0x23)
	BmsInformation              = byte(0x25)
	ActiveChargingRequest       = byte(0x31)
	RemoteBootstrapResponse     = byte(0x33)
	RemoteShutdownResponse      = byte(0x35)
	TransactionRecord           = byte(0x3b)
	BalanceUpdateResponse       = byte(0x41)
	CardSynchronizationResponse = byte(0x43)
	CardClearingResponse        = byte(0x45)
	CardQueryingResponse        = byte(0x47)
	SetWorkingParamsResponse    = byte(0x51)
	NtpResponse                 = byte(0x55)
	SetBillingModelResponse     = byte(0x57)
	FloorLockDataUpload         = byte(0x61)
	Response                    = byte(0x63)
	RemoteRebootResponse        = byte(0x91)
	OtaResponse                 = byte(0x93)

	// platform -> device
	VerificationResponse             = byte(0x02)
	HeartbeatResponse                = byte(0x04)
	BillingModelVerificationResponse = byte(0x06)
	BillingModelResponse             = byte(0x0a)
	RealTimeDataRequest              = byte(0x12)
	ChargingRequestConfirmed         = byte(0x32)
	RemoteBootstrapRequest           = byte(0x34)
	RemoteShutdownRequest            = byte(0x36)
	TransactionRecordConfirmed       = byte(0x40)
	AccountBalanceRemoteUpdate       = byte(0x42)
	CardSynchronizationRequest       = byte(0x44)
	CardClearingRequest              = byte(0x46)
	CardQueryingRequest              = byte(0x48)
	SetWorkingParamsRequest          = byte(0x52)
	NtpRequest                       = byte(0x56)
	SetBillingModelRequest           = byte(0x58)
	UpDownFloorLock                  = byte(0x62)
	RemoteRebootRequest              = byte(0x92)
	OtaRequest                       = byte(0x94)
)

type Header struct {
	Length    int    `json:"length"`
	Seq       int    `json:"seq"`
	Encrypted bool   `json:"encrypted"`
	FrameId   string `json:"frameId"`
}

type VerificationMessage struct {
	Header          *Header `json:"header"`
	Id              string  `json:"Id"`
	ElcType         int     `json:"elcType"`
	Guns            int     `json:"guns"`
	ProtocolVersion int     `json:"protocolVersion"`
	SoftwareVersion string  `json:"softwareVersion"`
	Network         int     `json:"network"`
	Sim             string  `json:"sim"`
	Operator        int     `json:"operator"`
}

func PackVerificationMessage(buf []byte, hex []string, header *Header) *VerificationMessage {
	//Id
	id := ""
	for _, v := range hex[6:13] {
		id += v
	}

	//type
	elcType := int(buf[13])

	//gun number
	guns := int(buf[14])

	//protocol version
	protocolVersion := int(buf[15]) / 10

	//software version
	softwareVersionBytes, _ := hex2.DecodeString(MakeHexStringFromHexArray(hex[16:24]))
	softwareVersion := string(softwareVersionBytes)

	//network type
	network := int(buf[25])

	//sim
	var sim string
	for _, v := range hex[26:36] {
		sim += v
	}

	//operator
	operator := int(buf[36])

	msg := &VerificationMessage{
		Header:          header,
		Id:              id,
		ElcType:         elcType,
		Guns:            guns,
		ProtocolVersion: protocolVersion,
		SoftwareVersion: softwareVersion,
		Network:         network,
		Sim:             sim,
		Operator:        operator,
	}
	return msg
}

type VerificationResponseMessage struct {
	Header *Header `json:"header"`
	Id     string  `json:"id"`
	Result bool    `json:"result"`
}

func PackVerificationResponseMessage(msg *VerificationResponseMessage) []byte {
	var resp bytes.Buffer
	resp.Write([]byte{StartFlag, 0x0c})
	seqStr := fmt.Sprintf("%x", msg.Header.Seq)
	seq := ConvertIntSeqToReversedHexArr(seqStr)
	resp.Write(HexToBytes(MakeHexStringFromHexArray(seq)))
	encrypted := byte(0x00)
	if msg.Header.Encrypted {
		encrypted = byte(0x01)
	}
	resp.Write([]byte{encrypted})
	resp.Write([]byte{VerificationResponse})
	resp.Write(HexToBytes(msg.Id))
	result := byte(0x01)
	if msg.Result {
		result = byte(0x00)
	}
	resp.Write([]byte{result})
	resp.Write(ModbusCRC(resp.Bytes()[2:]))
	return resp.Bytes()
}

type BillingModelVerificationMessage struct {
	Header           *Header `json:"header"`
	Id               string  `json:"Id"`
	BillingModelCode string  `json:"billingModelCode"`
}

func PackBillingModelVerificationMessage(hex []string, header *Header) *BillingModelVerificationMessage {
	//Id
	id := ""
	for _, v := range hex[6:13] {
		id += v
	}

	//billing model code
	bmcode := hex[13] + hex[14]

	msg := &BillingModelVerificationMessage{
		Header:           header,
		Id:               id,
		BillingModelCode: bmcode,
	}
	return msg
}

type BillingModelVerificationResponseMessage struct {
	Header           *Header `json:"header"`
	Id               string  `json:"id"`
	BillingModelCode string  `json:"billingModelCode"`
	Result           bool    `json:"result"`
}

func PackBillingModelVerificationResponseMessage(msg *BillingModelVerificationResponseMessage) []byte {
	var resp bytes.Buffer
	resp.Write(HexToBytes("680e"))

	seqStr := fmt.Sprintf("%x", msg.Header.Seq)
	seq := ConvertIntSeqToReversedHexArr(seqStr)
	resp.Write(HexToBytes(MakeHexStringFromHexArray(seq)))

	encrypted := byte(0x00)
	if msg.Header.Encrypted {
		encrypted = byte(0x01)
	}
	resp.Write([]byte{encrypted})
	resp.Write(HexToBytes("06"))
	resp.Write(HexToBytes(msg.Id))
	resp.Write(HexToBytes(msg.BillingModelCode))

	result := "01"
	if msg.Result {
		result = "00"
	}
	resp.Write(HexToBytes(result))
	resp.Write(ModbusCRC(resp.Bytes()[2:]))
	return resp.Bytes()
}

type HeartbeatMessage struct {
	Header    *Header `json:"header"`
	Id        string  `json:"Id"`
	Gun       string  `json:"gun"`
	GunStatus int     `json:"gunStatus"`
}

func PackHeartbeatMessage(hex []string, header *Header) *HeartbeatMessage {
	//id
	id := ""
	for _, v := range hex[6:13] {
		id += v
	}
	//gun
	gun := hex[13]
	//gun status
	gs, _ := strconv.ParseInt(hex[14], 16, 64)

	msg := &HeartbeatMessage{
		Header:    header,
		Id:        id,
		Gun:       gun,
		GunStatus: int(gs),
	}
	return msg
}

type HeartbeatResponseMessage struct {
	Header   *Header `json:"header"`
	Id       string  `json:"Id"`
	Gun      string  `json:"gun"`
	Response int     `json:"response"`
}

func PackHeartbeatResponseMessage(msg *HeartbeatResponseMessage) []byte {
	var resp bytes.Buffer
	resp.Write(HexToBytes("680d"))

	seqStr := fmt.Sprintf("%x", msg.Header.Seq)
	seq := ConvertIntSeqToReversedHexArr(seqStr)
	resp.Write(HexToBytes(MakeHexStringFromHexArray(seq)))

	encrypted := "00"
	if msg.Header.Encrypted {
		encrypted = "01"
	}

	resp.Write(HexToBytes(encrypted))
	resp.Write(HexToBytes("04"))
	resp.Write(HexToBytes(msg.Id))
	resp.Write(HexToBytes(msg.Gun))
	resp.Write(HexToBytes("00"))
	resp.Write(ModbusCRC(resp.Bytes()[2:]))

	return resp.Bytes()
}

type RemoteBootstrapRequestMessage struct {
	Header       *Header `json:"header"`
	TradeSeq     string  `json:"tradeSeq"`
	Id           string  `json:"id"`
	GunId        string  `json:"gunId"`
	LogicCard    string  `json:"logicCard"`
	PhysicalCard string  `json:"physicalCard"`
	Balance      int     `json:"balance"`
}

func PackRemoteBootstrapRequestMessage(msg *RemoteBootstrapRequestMessage) []byte {
	var resp bytes.Buffer
	resp.Write(HexToBytes("6830"))

	seqStr := fmt.Sprintf("%x", GenerateSeq())
	seq := ConvertIntSeqToReversedHexArr(seqStr)
	resp.Write(HexToBytes(MakeHexStringFromHexArray(seq)))

	encrypted := byte(0x00)
	if msg.Header.Encrypted {
		encrypted = byte(0x01)
	}
	resp.Write([]byte{encrypted})
	resp.Write(HexToBytes("34"))
	resp.Write(HexToBytes(msg.TradeSeq))
	resp.Write(HexToBytes(msg.Id))
	resp.Write(HexToBytes(msg.GunId))
	resp.Write(PadArrayWithZeros(HexToBytes(msg.LogicCard), 8))
	resp.Write(PadArrayWithZeros(HexToBytes(msg.PhysicalCard), 8))

	balance := HexToBytes(fmt.Sprintf("%x", msg.Balance))
	balance = PadArrayWithZeros(balance, 4)
	resp.Write(balance)

	resp.Write(ModbusCRC(resp.Bytes()[2:]))
	return resp.Bytes()
}

type RemoteBootstrapResponseMessage struct {
	Header   *Header `json:"header"`
	TradeSeq string  `json:"tradeSeq"`
	Id       string  `json:"id"`
	GunId    string  `json:"gunId"`
	Result   bool    `json:"result"`
	Reason   int     `json:"reason"`
}

func PackRemoteBootstrapResponseMessage(hex []string, header *Header) *RemoteBootstrapResponseMessage {
	//trade sequence number
	tradeSeq := ""
	for _, v := range hex[6:22] {
		tradeSeq += v
	}

	//id
	id := ""
	for _, v := range hex[22:29] {
		id += v
	}

	//gun id
	gunId := hex[29]

	//result
	result := false
	if hex[30] == "01" {
		result = true
	}

	//fail reason
	reason, _ := strconv.ParseInt(hex[14], 16, 64)

	msg := &RemoteBootstrapResponseMessage{
		Header:   header,
		TradeSeq: tradeSeq,
		Id:       id,
		GunId:    gunId,
		Result:   result,
		Reason:   int(reason),
	}
	return msg
}

type OfflineDataReportMessage struct {
	Header                  *Header `json:"header"`
	TradeSeq                string  `json:"tradeSeq"`
	Id                      string  `json:"id"`
	GunId                   string  `json:"gunId"`
	Status                  int     `json:"status"`
	Reset                   int     `json:"reset"`
	Plugged                 int     `json:"plugged"`
	Ov                      int     `json:"ov"`
	Oc                      int     `json:"oc"`
	LineTemp                int     `json:"lineTemp"`
	LineCode                string  `json:"lineCode"`
	Soc                     int     `json:"soc"`
	BpTopTemp               int     `json:"bpTopTemp"`
	AccumulatedChargingTime int     `json:"accumulatedChargingTime"`
	RemainingTime           int     `json:"remainingTime"`
	ChargingDegrees         int     `json:"chargingDegrees"`
	LossyChargingDegrees    int     `json:"lossyChargingDegrees"`
	ChargedAmount           int     `json:"chargedAmount"`
	HardwareFailure         int     `json:"hardwareFailure"`
}

func PackOfflineDataReportMessage(hex []string, header *Header) *OfflineDataReportMessage {
	//trade sequence number
	tradeSeq := ""
	for _, v := range hex[6:22] {
		tradeSeq += v
	}

	//id
	id := ""
	for _, v := range hex[22:29] {
		id += v
	}

	//gun id
	gunId := hex[29]

	//status
	status, _ := strconv.ParseInt(hex[30], 16, 64)

	//reset
	reset, _ := strconv.ParseInt(hex[31], 16, 64)

	//plugged
	plugged, _ := strconv.ParseInt(hex[32], 16, 64)

	//ov
	ov, _ := strconv.ParseInt(MakeHexStringFromHexArray(hex[33:35]), 16, 64)

	//oc
	oc, _ := strconv.ParseInt(MakeHexStringFromHexArray(hex[35:37]), 16, 64)

	//lineTemp
	lineTemp, _ := strconv.ParseInt(hex[37], 16, 64)

	//lineCode
	lineCode := MakeHexStringFromHexArray(hex[38:46])

	//soc
	soc, _ := strconv.ParseInt(hex[46], 16, 64)

	//bpTopTemp
	bpTopTemp, _ := strconv.ParseInt(hex[47], 16, 64)

	//accumulatedChargingTime
	accumulatedChargingTime, _ := strconv.ParseInt(MakeHexStringFromHexArray(hex[48:50]), 16, 64)

	//remainingTime
	remainingTime, _ := strconv.ParseInt(MakeHexStringFromHexArray(hex[50:52]), 16, 64)

	//chargingDegrees
	chargingDegrees, _ := strconv.ParseInt(MakeHexStringFromHexArray(hex[52:56]), 16, 64)

	//lossyChargingDegrees
	lossyChargingDegrees, _ := strconv.ParseInt(MakeHexStringFromHexArray(hex[56:60]), 16, 64)

	//chargedAmount
	chargedAmount, _ := strconv.ParseInt(MakeHexStringFromHexArray(hex[60:64]), 16, 64)

	//hardwareFailure
	hardwareFailure, _ := strconv.ParseInt(MakeHexStringFromHexArray(hex[64:66]), 16, 64)

	msg := &OfflineDataReportMessage{
		Header:                  header,
		TradeSeq:                tradeSeq,
		Id:                      id,
		GunId:                   gunId,
		Status:                  int(status),
		Reset:                   int(reset),
		Plugged:                 int(plugged),
		Ov:                      int(ov),
		Oc:                      int(oc),
		LineTemp:                int(lineTemp),
		LineCode:                lineCode,
		Soc:                     int(soc),
		BpTopTemp:               int(bpTopTemp),
		AccumulatedChargingTime: int(accumulatedChargingTime),
		RemainingTime:           int(remainingTime),
		ChargingDegrees:         int(chargingDegrees),
		LossyChargingDegrees:    int(lossyChargingDegrees),
		ChargedAmount:           int(chargedAmount),
		HardwareFailure:         int(hardwareFailure),
	}

	return msg
}

type RemoteShutdownResponseMessage struct {
	Header *Header `json:"header"`
	Id     string  `json:"id"`
	GunId  string  `json:"gunId"`
	Result bool    `json:"result"`
	Reason int     `json:"reason"`
}

func PackRemoteShutdownResponseMessage(hex []string, header *Header) *RemoteShutdownResponseMessage {
	//id
	id := ""
	for _, v := range hex[6:13] {
		id += v
	}

	//gun id
	gunId := hex[13]

	//result
	result := false
	if hex[14] == "01" {
		result = true
	}

	//fail reason
	reason, _ := strconv.ParseInt(hex[15], 16, 64)

	msg := &RemoteShutdownResponseMessage{
		Header: header,
		Id:     id,
		GunId:  gunId,
		Result: result,
		Reason: int(reason),
	}
	return msg
}

type RemoteShutdownRequestMessage struct {
	Header *Header `json:"header"`
	Id     string  `json:"id"`
	GunId  string  `json:"gunId"`
}

func PackRemoteShutdownRequestMessage(msg *RemoteShutdownRequestMessage) []byte {
	var resp bytes.Buffer
	resp.Write(HexToBytes("680c"))
	seqStr := fmt.Sprintf("%x", GenerateSeq())
	seq := ConvertIntSeqToReversedHexArr(seqStr)
	resp.Write(HexToBytes(MakeHexStringFromHexArray(seq)))
	if msg.Header.Encrypted {
		resp.WriteByte(0x01)
	} else {
		resp.WriteByte(0x00)
	}
	resp.Write(HexToBytes("36"))
	resp.Write(HexToBytes(msg.Id))
	resp.Write(HexToBytes(msg.GunId))
	resp.Write(ModbusCRC(resp.Bytes()[2:]))
	return resp.Bytes()
}

type TransactionRecordMessage struct {
	Header                    *Header `json:"header"`
	TradeSeq                  string  `json:"tradeSeq"`
	Id                        string  `json:"id"`
	GunId                     string  `json:"gunId"`
	StartAt                   int64   `json:"startAt"`
	EndAt                     int64   `json:"endAt"`
	SharpUnitPrice            int64   `json:"sharpUnitPrice"`
	SharpElectricCharge       int64   `json:"sharpElectricCharge"`
	LossySharpElectricCharge  int64   `json:"lossySharpElectricCharge"`
	SharpPrice                int64   `json:"sharpPrice"`
	PeakUnitPrice             int64   `json:"peakUnitPrice"`
	PeakElectricCharge        int64   `json:"peakElectricCharge"`
	LossyPeakElectricCharge   int64   `json:"lossyPeakElectricCharge"`
	PeakPrice                 int64   `json:"peakPrice"`
	FlatUnitPrice             int64   `json:"flatUnitPrice"`
	FlatElectricCharge        int64   `json:"flatElectricCharge"`
	LossyFlatElectricCharge   int64   `json:"lossyFlatElectricCharge"`
	FlatPrice                 int64   `json:"flatPrice"`
	ValleyUnitPrice           int64   `json:"valleyUnitPrice"`
	ValleyElectricCharge      int64   `json:"valleyElectricCharge"`
	LossyValleyElectricCharge int64   `json:"lossyValleyElectricCharge"`
	ValleyPrice               int64   `json:"valleyPrice"`
	InitialMeterReading       int64   `json:"initialMeterReading"`
	FinalMeterReading         int64   `json:"finalMeterReading"`
	TotalElectricCharge       int64   `json:"totalElectricCharge"`
	LossyTotalElectricCharge  int64   `json:"lossyTotalElectricCharge"`
	ConsumptionAmount         int64   `json:"consumptionAmount"`
	Vin                       string  `json:"vin"`
	StartType                 int     `json:"startType"`
	TransactionDateTime       int64   `json:"transactionDateTime"`
	StopReason                int     `json:"stopReason"`
	PhysicalCardNumber        string  `json:"physicalCardNumber"`
}

func PackTransactionRecordMessage(raw []byte, hex []string, header *Header) *TransactionRecordMessage {
	//trade sequence number
	tradeSeq := ""
	for _, v := range hex[6:22] {
		tradeSeq += v
	}

	//id
	id := ""
	for _, v := range hex[22:29] {
		id += v
	}

	//gun id
	gunId := hex[29]

	//start time
	startAt := Cp56time2aToUnixMilliseconds(raw[30:37])

	//end time
	endAt := Cp56time2aToUnixMilliseconds(raw[37:44])

	//sharp unit price
	sharpUnitPrice, _ := strconv.ParseInt(MakeHexStringFromHexArray(hex[44:48]), 16, 64)

	//sharp electric charge
	sharpElectricCharge, _ := strconv.ParseInt(MakeHexStringFromHexArray(hex[48:52]), 16, 64)

	//lossy sharp electric charge
	lossySharpElectricCharge, _ := strconv.ParseInt(MakeHexStringFromHexArray(hex[52:56]), 16, 64)

	//sharp price
	sharpPrice, _ := strconv.ParseInt(MakeHexStringFromHexArray(hex[56:60]), 16, 64)

	//peak unit price
	peakUnitPrice, _ := strconv.ParseInt(MakeHexStringFromHexArray(hex[60:64]), 16, 64)

	//peak electric charge
	peakElectricCharge, _ := strconv.ParseInt(MakeHexStringFromHexArray(hex[64:68]), 16, 64)

	//lossy peak electric charge
	lossyPeakElectricCharge, _ := strconv.ParseInt(MakeHexStringFromHexArray(hex[68:72]), 16, 64)

	//peak price
	peakPrice, _ := strconv.ParseInt(MakeHexStringFromHexArray(hex[72:76]), 16, 64)

	//flat unit price
	flatUnitPrice, _ := strconv.ParseInt(MakeHexStringFromHexArray(hex[76:80]), 16, 64)

	//flat electric charge
	flatElectricCharge, _ := strconv.ParseInt(MakeHexStringFromHexArray(hex[80:84]), 16, 64)

	//lossy flat electric charge
	lossyFlatElectricCharge, _ := strconv.ParseInt(MakeHexStringFromHexArray(hex[84:88]), 16, 64)

	//flat price
	flatPrice, _ := strconv.ParseInt(MakeHexStringFromHexArray(hex[88:92]), 16, 64)

	//valley unit price
	valleyUnitPrice, _ := strconv.ParseInt(MakeHexStringFromHexArray(hex[92:96]), 16, 64)

	//valley electric charge
	valleyElectricCharge, _ := strconv.ParseInt(MakeHexStringFromHexArray(hex[96:100]), 16, 64)

	//lossy valley electric charge
	lossyValleyElectricCharge, _ := strconv.ParseInt(MakeHexStringFromHexArray(hex[100:104]), 16, 64)

	//valley price
	valleyPrice, _ := strconv.ParseInt(MakeHexStringFromHexArray(hex[104:108]), 16, 64)

	//initial meter reading
	initialMeterReading, _ := strconv.ParseInt(MakeHexStringFromHexArray(hex[108:113]), 16, 64)

	//final meter reading
	finalMeterReading, _ := strconv.ParseInt(MakeHexStringFromHexArray(hex[113:118]), 16, 64)

	//total electric charge
	totalElectricCharge, _ := strconv.ParseInt(MakeHexStringFromHexArray(hex[118:122]), 16, 64)

	//lossy total electric charge
	lossyTotalElectricCharge, _ := strconv.ParseInt(MakeHexStringFromHexArray(hex[122:126]), 16, 64)

	//consumption amount
	consumptionAmount, _ := strconv.ParseInt(MakeHexStringFromHexArray(hex[126:130]), 16, 64)

	//vin
	vin := MakeHexStringFromHexArray(hex[130:147])

	//start type
	startType, _ := strconv.ParseInt(hex[147], 16, 64)

	//transaction date time
	transactionDateTime := Cp56time2aToUnixMilliseconds(raw[148:155])

	//stop reason
	stopReason, _ := strconv.ParseInt(hex[155], 16, 64)

	//physical card number
	physicalCardNumber := MakeHexStringFromHexArray(hex[156:164])

	//fill all fields
	msg := &TransactionRecordMessage{
		TradeSeq:                  tradeSeq,
		Id:                        id,
		GunId:                     gunId,
		StartAt:                   startAt,
		EndAt:                     endAt,
		SharpUnitPrice:            sharpUnitPrice,
		SharpElectricCharge:       sharpElectricCharge,
		LossySharpElectricCharge:  lossySharpElectricCharge,
		SharpPrice:                sharpPrice,
		PeakUnitPrice:             peakUnitPrice,
		PeakElectricCharge:        peakElectricCharge,
		LossyPeakElectricCharge:   lossyPeakElectricCharge,
		PeakPrice:                 peakPrice,
		FlatUnitPrice:             flatUnitPrice,
		FlatElectricCharge:        flatElectricCharge,
		LossyFlatElectricCharge:   lossyFlatElectricCharge,
		FlatPrice:                 flatPrice,
		ValleyUnitPrice:           valleyUnitPrice,
		ValleyElectricCharge:      valleyElectricCharge,
		LossyValleyElectricCharge: lossyValleyElectricCharge,
		ValleyPrice:               valleyPrice,
		InitialMeterReading:       initialMeterReading,
		FinalMeterReading:         finalMeterReading,
		TotalElectricCharge:       totalElectricCharge,
		LossyTotalElectricCharge:  lossyTotalElectricCharge,
		ConsumptionAmount:         consumptionAmount,
		Vin:                       vin,
		StartType:                 int(startType),
		TransactionDateTime:       transactionDateTime,
		StopReason:                int(stopReason),
		PhysicalCardNumber:        physicalCardNumber,
	}
	return msg
}

type TransactionRecordConfirmedMessage struct {
	Header   *Header `json:"header"`
	Id       string  `json:"id"`
	TradeSeq string  `json:"tradeSeq"`
	Result   int     `json:"result"`
}

func PackTransactionRecordConfirmedMessage(msg *TransactionRecordConfirmedMessage) []byte {
	var resp bytes.Buffer
	resp.Write([]byte{0x68, 0x15})
	seqStr := fmt.Sprintf("%x", GenerateSeq())
	seq := ConvertIntSeqToReversedHexArr(seqStr)
	resp.Write(HexToBytes(MakeHexStringFromHexArray(seq)))
	if msg.Header.Encrypted {
		resp.WriteByte(0x01)
	} else {
		resp.WriteByte(0x00)
	}
	resp.Write([]byte{0x40})
	resp.Write(HexToBytes(msg.TradeSeq))
	resp.Write([]byte(strconv.Itoa(msg.Result)))
	resp.Write(ModbusCRC(resp.Bytes()[2:]))
	return resp.Bytes()
}

type RemoteRebootRequestMessage struct {
	Header  *Header `json:"header"`
	Id      string  `json:"id"`
	Control int     `json:"control"`
}

func PackRemoteRebootRequestMessage(msg *RemoteRebootRequestMessage) []byte {
	var resp bytes.Buffer
	resp.Write([]byte{0x68, 0x0c})
	seqStr := fmt.Sprintf("%x", GenerateSeq())
	seq := ConvertIntSeqToReversedHexArr(seqStr)
	resp.Write(HexToBytes(MakeHexStringFromHexArray(seq)))
	if msg.Header.Encrypted {
		resp.WriteByte(0x01)
	} else {
		resp.WriteByte(0x00)
	}
	resp.Write([]byte{0x92})
	resp.Write(HexToBytes(msg.Id))
	resp.Write([]byte(strconv.Itoa(msg.Control)))
	resp.Write(ModbusCRC(resp.Bytes()[2:]))
	return resp.Bytes()
}
