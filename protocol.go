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
	BillingModelRequest         = byte(0x09)
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

type BillingModelRequestMessage struct {
	Header *Header `json:"header"`
	Id     string  `json:"Id"`
}

func PackBillingModelRequestMessage(hex []string, header *Header) *BillingModelRequestMessage {
	//Id
	id := ""
	for _, v := range hex[6:13] {
		id += v
	}

	msg := &BillingModelRequestMessage{
		Header: header,
		Id:     id,
	}
	return msg
}

type BillingModelResponseMessage struct {
	Header           *Header `json:"header"`
	Id               string  `json:"id"`
	BillingModelCode string  `json:"billingModelCode"`
	SharpUnitPrice   int     `json:"sharpUnitPrice"`
	SharpServiceFee  int     `json:"sharpServiceFee"`
	PeakUnitPrice    int     `json:"peakUnitPrice"`
	PeakServiceFee   int     `json:"peakServiceFee"`
	FlatUnitPrice    int     `json:"flatUnitPrice"`
	FlatServiceFee   int     `json:"flatServiceFee"`
	ValleyUnitPrice  int     `json:"valleyUnitPrice"`
	ValleyServiceFee int     `json:"valleyServiceFee"`
	AccrualRatio     int     `json:"accrualRatio"`
	RateList         []int   `json:"rateList"`
}

func PackBillingModelResponseMessage(msg *BillingModelResponseMessage) []byte {
	var resp bytes.Buffer
	resp.Write([]byte{0x68, 0x5e})
	seqStr := fmt.Sprintf("%x", GenerateSeq())
	seq := ConvertIntSeqToReversedHexArr(seqStr)
	resp.Write(HexToBytes(MakeHexStringFromHexArray(seq)))
	if msg.Header.Encrypted {
		resp.WriteByte(0x01)
	} else {
		resp.WriteByte(0x00)
	}
	resp.Write([]byte{0x0a})
	resp.Write(HexToBytes(msg.Id))
	resp.Write(HexToBytes(msg.BillingModelCode))
	resp.Write(IntToBIN(msg.SharpUnitPrice, 4))
	resp.Write(IntToBIN(msg.SharpServiceFee, 4))
	resp.Write(IntToBIN(msg.PeakUnitPrice, 4))
	resp.Write(IntToBIN(msg.PeakServiceFee, 4))
	resp.Write(IntToBIN(msg.FlatUnitPrice, 4))
	resp.Write(IntToBIN(msg.FlatServiceFee, 4))
	resp.Write(IntToBIN(msg.ValleyUnitPrice, 4))
	resp.Write(IntToBIN(msg.ValleyServiceFee, 4))
	resp.Write([]byte(strconv.Itoa(msg.AccrualRatio)))

	for _, v := range msg.RateList {
		resp.Write(IntToBIN(v, 1))
	}

	resp.Write(ModbusCRC(resp.Bytes()[2:]))
	return resp.Bytes()
}

type BillingModelVerificationResponseMessage struct {
	Header           *Header `json:"header"`
	Id               string  `json:"id"`
	BillingModelCode string  `json:"billingModelCode"`
	Result           bool    `json:"result"`
}

func PackBillingModelVerificationResponseMessage(msg *BillingModelVerificationResponseMessage) []byte {
	var resp bytes.Buffer
	resp.Write([]byte{StartFlag, 0x0e})

	seqStr := fmt.Sprintf("%x", msg.Header.Seq)
	seq := ConvertIntSeqToReversedHexArr(seqStr)
	resp.Write(HexToBytes(MakeHexStringFromHexArray(seq)))

	encrypted := byte(0x00)
	if msg.Header.Encrypted {
		encrypted = byte(0x01)
	}
	resp.Write([]byte{encrypted})
	resp.Write([]byte{BillingModelVerificationResponse})
	resp.Write(HexToBytes(msg.Id))
	resp.Write(HexToBytes(msg.BillingModelCode))

	result := byte(0x01)
	if msg.Result {
		result = byte(0x00)
	}
	resp.Write([]byte{result})
	resp.Write(ModbusCRC(resp.Bytes()[2:]))
	return resp.Bytes()
}

type HeartbeatMessage struct {
	Header    *Header `json:"header"`
	Id        string  `json:"id"`
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
	reason, _ := strconv.ParseInt(hex[31], 16, 64)

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

func PackOfflineDataReportMessage(hex []string, raw []byte, header *Header) *OfflineDataReportMessage {
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
	status := BINToInt([]byte{raw[30]})

	//reset
	reset := BINToInt([]byte{raw[31]})

	//plugged
	plugged := BINToInt([]byte{raw[32]})

	//ov
	ov := BINToInt(raw[33:35])

	//oc
	oc := BINToInt(raw[35:37])

	//lineTemp
	lineTemp := BINToInt([]byte{raw[37]})

	//lineCode
	lineCode := BINToInt(raw[38:46])

	//soc
	soc := BINToInt([]byte{raw[46]})

	//bpTopTemp
	bpTopTemp := BINToInt([]byte{raw[47]})

	//accumulatedChargingTime
	accumulatedChargingTime := BINToInt(raw[48:50])

	//remainingTime
	remainingTime := BINToInt(raw[50:52])

	//chargingDegrees
	chargingDegrees := BINToInt(raw[52:56])

	//lossyChargingDegrees
	lossyChargingDegrees := BINToInt(raw[56:60])

	//chargedAmount
	chargedAmount := BINToInt(raw[60:64])

	//hardwareFailure
	hardwareFailure := BINToInt(raw[64:66])

	msg := &OfflineDataReportMessage{
		Header:                  header,
		TradeSeq:                tradeSeq,
		Id:                      id,
		GunId:                   gunId,
		Status:                  status,
		Reset:                   reset,
		Plugged:                 plugged,
		Ov:                      ov,
		Oc:                      oc,
		LineTemp:                lineTemp,
		LineCode:                strconv.Itoa(lineCode),
		Soc:                     soc,
		BpTopTemp:               bpTopTemp,
		AccumulatedChargingTime: accumulatedChargingTime,
		RemainingTime:           remainingTime,
		ChargingDegrees:         chargingDegrees,
		LossyChargingDegrees:    lossyChargingDegrees,
		ChargedAmount:           chargedAmount,
		HardwareFailure:         hardwareFailure,
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
	sharpUnitPrice := BINToInt(raw[44:48])

	//sharp electric charge
	sharpElectricCharge := BINToInt(raw[48:52])

	//lossy sharp electric charge
	lossySharpElectricCharge := BINToInt(raw[52:56])

	//sharp price
	sharpPrice := BINToInt(raw[56:60])

	//peak unit price
	peakUnitPrice := BINToInt(raw[60:64])

	//peak electric charge
	peakElectricCharge := BINToInt(raw[64:68])

	//lossy peak electric charge
	lossyPeakElectricCharge := BINToInt(raw[68:72])

	//peak price
	peakPrice := BINToInt(raw[72:76])

	//flat unit price
	flatUnitPrice := BINToInt(raw[76:80])

	//flat electric charge
	flatElectricCharge := BINToInt(raw[80:84])

	//lossy flat electric charge
	lossyFlatElectricCharge := BINToInt(raw[84:88])

	//flat price
	flatPrice := BINToInt(raw[88:92])

	//valley unit price
	valleyUnitPrice := BINToInt(raw[92:96])

	//valley electric charge
	valleyElectricCharge := BINToInt(raw[96:100])

	//lossy valley electric charge
	lossyValleyElectricCharge := BINToInt(raw[100:104])

	//valley price
	valleyPrice := BINToInt(raw[104:108])

	//initial meter reading
	initialMeterReading := BINToInt(raw[108:113])

	//final meter reading
	finalMeterReading := BINToInt(raw[113:118])

	//total electric charge
	totalElectricCharge := BINToInt(raw[118:122])

	//lossy total electric charge
	lossyTotalElectricCharge := BINToInt(raw[122:126])

	//consumption amount
	consumptionAmount := BINToInt(raw[126:130])

	//vin
	vin := MakeHexStringFromHexArray(hex[130:147])

	//start type
	startType := BINToInt([]byte{raw[147]})

	//transaction date time
	transactionDateTime := Cp56time2aToUnixMilliseconds(raw[148:155])

	//stop reason
	stopReason := BINToInt([]byte{raw[155]})

	//physical card number
	physicalCardNumber := MakeHexStringFromHexArray(hex[156:164])

	//fill all fields
	msg := &TransactionRecordMessage{
		Header:                    header,
		TradeSeq:                  tradeSeq,
		Id:                        id,
		GunId:                     gunId,
		StartAt:                   startAt,
		EndAt:                     endAt,
		SharpUnitPrice:            int64(sharpUnitPrice),
		SharpElectricCharge:       int64(sharpElectricCharge),
		LossySharpElectricCharge:  int64(lossySharpElectricCharge),
		SharpPrice:                int64(sharpPrice),
		PeakUnitPrice:             int64(peakUnitPrice),
		PeakElectricCharge:        int64(peakElectricCharge),
		LossyPeakElectricCharge:   int64(lossyPeakElectricCharge),
		PeakPrice:                 int64(peakPrice),
		FlatUnitPrice:             int64(flatUnitPrice),
		FlatElectricCharge:        int64(flatElectricCharge),
		LossyFlatElectricCharge:   int64(lossyFlatElectricCharge),
		FlatPrice:                 int64(flatPrice),
		ValleyUnitPrice:           int64(valleyUnitPrice),
		ValleyElectricCharge:      int64(valleyElectricCharge),
		LossyValleyElectricCharge: int64(lossyValleyElectricCharge),
		ValleyPrice:               int64(valleyPrice),
		InitialMeterReading:       int64(initialMeterReading),
		FinalMeterReading:         int64(finalMeterReading),
		TotalElectricCharge:       int64(totalElectricCharge),
		LossyTotalElectricCharge:  int64(lossyTotalElectricCharge),
		ConsumptionAmount:         int64(consumptionAmount),
		Vin:                       vin,
		StartType:                 startType,
		TransactionDateTime:       transactionDateTime,
		StopReason:                stopReason,
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

type RemoteRebootResponseMessage struct {
	Header *Header `json:"header"`
	Id     string  `json:"id"`
	Result int     `json:"result"`
}

func PackRemoteRebootResponseMessage(hex []string, header *Header) *RemoteRebootResponseMessage {
	//id
	id := ""
	for _, v := range hex[6:13] {
		id += v
	}

	//result 0-fail 1-success
	result := 1
	if hex[13] == "00" {
		result = 0
	}

	msg := &RemoteRebootResponseMessage{
		Header: header,
		Id:     id,
		Result: result,
	}
	return msg
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

type SetBillingModelRequestMessage struct {
	Header           *Header `json:"header"`
	Id               string  `json:"id"`
	BillingModelCode string  `json:"billingModelCode"`
	SharpUnitPrice   int     `json:"sharpUnitPrice"`
	SharpServiceFee  int     `json:"sharpServiceFee"`
	PeakUnitPrice    int     `json:"peakUnitPrice"`
	PeakServiceFee   int     `json:"peakServiceFee"`
	FlatUnitPrice    int     `json:"flatUnitPrice"`
	FlatServiceFee   int     `json:"flatServiceFee"`
	ValleyUnitPrice  int     `json:"valleyUnitPrice"`
	ValleyServiceFee int     `json:"valleyServiceFee"`
	AccrualRatio     int     `json:"accrualRatio"`
	RateList         []int   `json:"rateList"`
}

func PackSetBillingModelRequestMessage(msg *SetBillingModelRequestMessage) []byte {
	var resp bytes.Buffer
	resp.Write([]byte{0x68, 0x5e})
	seqStr := fmt.Sprintf("%x", GenerateSeq())
	seq := ConvertIntSeqToReversedHexArr(seqStr)
	resp.Write(HexToBytes(MakeHexStringFromHexArray(seq)))
	if msg.Header.Encrypted {
		resp.WriteByte(0x01)
	} else {
		resp.WriteByte(0x00)
	}
	resp.Write([]byte{0x0a})
	resp.Write(HexToBytes(msg.Id))
	resp.Write(HexToBytes(msg.BillingModelCode))
	resp.Write(IntToBIN(msg.SharpUnitPrice, 4))
	resp.Write(IntToBIN(msg.SharpServiceFee, 4))
	resp.Write(IntToBIN(msg.PeakUnitPrice, 4))
	resp.Write(IntToBIN(msg.PeakServiceFee, 4))
	resp.Write(IntToBIN(msg.FlatUnitPrice, 4))
	resp.Write(IntToBIN(msg.FlatServiceFee, 4))
	resp.Write(IntToBIN(msg.ValleyUnitPrice, 4))
	resp.Write(IntToBIN(msg.ValleyServiceFee, 4))
	resp.Write([]byte(strconv.Itoa(msg.AccrualRatio)))

	for _, v := range msg.RateList {
		resp.Write(IntToBIN(v, 1))
	}

	resp.Write(ModbusCRC(resp.Bytes()[2:]))
	return resp.Bytes()
}

type SetBillingModelResponseMessage struct {
	Header *Header `json:"header"`
	Id     string  `json:"id"`
	Result int     `json:"result"`
}

func PackSetBillingModelResponseMessage(hex []string, header *Header) *SetBillingModelResponseMessage {
	//id
	id := ""
	for _, v := range hex[6:13] {
		id += v
	}

	//result 0-fail 1-success
	result := 1
	if hex[13] == "00" {
		result = 0
	}

	msg := &SetBillingModelResponseMessage{
		Header: header,
		Id:     id,
		Result: result,
	}
	return msg
}

type ChargingFinishedMessage struct {
	Header                           *Header `json:"header"`
	TradeSeq                         string  `json:"tradeSeq"`
	Id                               string  `json:"id"`
	GunId                            string  `json:"gunId"`
	BmsSoc                           int     `json:"bmsSoc"`
	BmsBatteryPackLowestVoltage      int     `json:"bmsBatteryPackLowestVoltage"`
	BmsBatteryPackHighestVoltage     int     `json:"bmsBatteryPackHighestVoltage"`
	BmsBatteryPackLowestTemperature  int     `json:"bmsBatteryPackLowestTemperature"`
	BmsBatteryPackHighestTemperature int     `json:"bmsBatteryPackHighestTemperature"`
	CumulativeChargingDuration       int     `json:"cumulativeChargingDuration"`
	OutputPower                      int     `json:"outputPower"`
	ChargingUnitId                   int     `json:"chargingUnitId"`
}

func PackChargingFinishedMessage(hex []string, header *Header) *ChargingFinishedMessage {
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

	//soc
	soc, _ := strconv.ParseInt(hex[30], 16, 64)

	bmsBatteryPackLowestVoltage, _ := strconv.ParseInt(MakeHexStringFromHexArray(hex[31:33]), 16, 64)
	bmsBatteryPackHighestVoltage, _ := strconv.ParseInt(MakeHexStringFromHexArray(hex[33:35]), 16, 64)
	bmsBatteryPackLowestTemperature, _ := strconv.ParseInt(hex[35], 16, 64)
	bmsBatteryPackHighestTemperature, _ := strconv.ParseInt(hex[36], 16, 64)
	cumulativeChargingDuration, _ := strconv.ParseInt(MakeHexStringFromHexArray(hex[37:39]), 16, 64)
	outputPower, _ := strconv.ParseInt(MakeHexStringFromHexArray(hex[39:41]), 16, 64)
	chargingUnitId, _ := strconv.ParseInt(MakeHexStringFromHexArray(hex[41:45]), 16, 64)

	msg := &ChargingFinishedMessage{
		Header:                           header,
		TradeSeq:                         tradeSeq,
		Id:                               id,
		GunId:                            gunId,
		BmsSoc:                           int(soc),
		BmsBatteryPackLowestVoltage:      int(bmsBatteryPackLowestVoltage),
		BmsBatteryPackHighestVoltage:     int(bmsBatteryPackHighestVoltage),
		BmsBatteryPackLowestTemperature:  int(bmsBatteryPackLowestTemperature),
		BmsBatteryPackHighestTemperature: int(bmsBatteryPackHighestTemperature),
		CumulativeChargingDuration:       int(cumulativeChargingDuration),
		OutputPower:                      int(outputPower),
		ChargingUnitId:                   int(chargingUnitId),
	}
	return msg
}
