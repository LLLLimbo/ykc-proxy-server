package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net"
)

func VerificationRouter(opt *Options, buf []byte, hex []string, header *Header, conn net.Conn) {
	msg := PackVerificationMessage(buf, hex, header)

	log.WithFields(log.Fields{
		"id":               msg.Id,
		"elc_type":         msg.ElcType,
		"guns":             msg.Guns,
		"protocol_version": msg.ProtocolVersion,
		"software_version": msg.SoftwareVersion,
		"network":          msg.Network,
		"sim":              msg.Sim,
		"operator":         msg.Operator,
	}).Debug("[01] Verification message")
	StoreClient(msg.Id, conn)

	//auto response
	if opt.AutoVerification {
		m := &VerificationResponseMessage{
			Header: &Header{
				Seq:       0,
				Encrypted: false,
			},
			Id:     msg.Id,
			Result: true,
		}
		_ = ResponseToVerification(m)
		return
	}

	//forward
	if opt.MessageForwarder != nil {
		//convert msg to json string bytes
		b, _ := json.Marshal(msg)
		_ = opt.MessageForwarder.Publish("01", b)
	}

}

func VerificationResponseRouter(c *gin.Context) {
	var req VerificationResponseMessage
	if c.ShouldBind(&req) == nil {
		err := ResponseToVerification(&req)
		if err != nil {
			c.JSON(500, gin.H{"message": err})
			return
		}
	}
	c.JSON(200, gin.H{"message": "done"})
}

func HeartbeatRouter(opt *Options, hex []string, header *Header, conn net.Conn) {
	msg := PackHeartbeatMessage(hex, header)
	log.WithFields(log.Fields{
		"id":         msg.Id,
		"gun":        msg.Gun,
		"gun_status": msg.GunStatus,
	}).Debug("[03] Heartbeat message")

	//auto response
	if opt.AutoHeartbeatResponse {
		_ = ResponseToHeartbeat(&HeartbeatResponseMessage{
			Header:   header,
			Id:       msg.Id,
			Gun:      msg.Gun,
			Response: 0,
		})
		return
	}

	//forward
	if opt.MessageForwarder != nil {
		//convert msg to json string bytes
		b, _ := json.Marshal(msg)
		_ = opt.MessageForwarder.Publish("03", b)
	}
}

func BillingModelVerificationRouter(opt *Options, hex []string, header *Header, conn net.Conn) {
	msg := PackBillingModelVerificationMessage(hex, header)
	log.WithFields(log.Fields{
		"id":                 msg.Id,
		"billing_model_code": msg.BillingModelCode,
	}).Debug("[05] BillingModelRequest message")

	//auto response
	if opt.AutoBillingModelVerify {
		m := &BillingModelVerificationResponseMessage{
			Header: &Header{
				Seq:       0,
				Encrypted: false,
			},
			Id:               msg.Id,
			BillingModelCode: msg.BillingModelCode,
			Result:           true,
		}
		_ = ResponseToBillingModelVerification(m)
		return
	}

	//forward
	if opt.MessageForwarder != nil {
		//convert msg to json string bytes
		b, _ := json.Marshal(msg)
		_ = opt.MessageForwarder.Publish("05", b)
	}
}

func BillingModelVerificationResponseRouter(c *gin.Context) {
	var req BillingModelVerificationResponseMessage
	if c.ShouldBind(&req) == nil {
		err := ResponseToBillingModelVerification(&req)
		if err != nil {
			c.JSON(500, gin.H{"message": err})
			return
		}
	}
	c.JSON(200, gin.H{"message": "done"})
}

func RemoteBootstrapRequestRouter(c *gin.Context) {
	var req RemoteBootstrapRequestMessage
	if c.ShouldBind(&req) == nil {
		err := SendRemoteBootstrapRequest(&req)
		if err != nil {
			c.JSON(500, gin.H{"message": err})
			return
		}
	}
	c.JSON(200, gin.H{"message": "done"})
}

func RemoteBootstrapResponseRouter(opt *Options, hex []string, header *Header) {
	msg := PackRemoteBootstrapResponseMessage(hex, header)
	log.WithFields(log.Fields{
		"id":                    msg.Id,
		"trade_sequence_number": msg.TradeSeq,
		"gun_id":                msg.GunId,
		"result":                msg.Result,
		"reason":                msg.Reason,
	}).Debug("[33] RemoteBootstrapResponse message")

	//forward
	if opt.MessageForwarder != nil {
		//convert msg to json string bytes
		b, _ := json.Marshal(msg)
		_ = opt.MessageForwarder.Publish("33", b)
	}
}

func OfflineDataReportMessageRouter(opt *Options, hex []string, header *Header) {
	msg := PackOfflineDataReportMessage(hex, header)
	log.WithFields(log.Fields{
		"id":                               msg.Id,
		"trade_sequence_number":            msg.TradeSeq,
		"gun_id":                           msg.GunId,
		"status":                           msg.Status,
		"reset":                            msg.Reset,
		"plugged":                          msg.Plugged,
		"output_voltage":                   msg.Ov,
		"output_current":                   msg.Oc,
		"gun_line_temperature":             msg.LineTemp,
		"gun_line_encoding":                msg.LineCode,
		"battery_pack_highest_temperature": msg.BpTopTemp,
		"accumulated_charging_time":        msg.AccumulatedChargingTime,
		"remaining_time":                   msg.RemainingTime,
		"charging_degrees":                 msg.ChargingDegrees,
		"lossy_charging_degrees":           msg.LossyChargingDegrees,
		"charged_amount":                   msg.ChargedAmount,
		"hardware_failure":                 msg.HardwareFailure,
	}).Debug("[13] OfflineDataReport message")

	//forward
	if opt.MessageForwarder != nil {
		//convert msg to json string bytes
		b, _ := json.Marshal(msg)
		_ = opt.MessageForwarder.Publish("13", b)
	}
}

func RemoteShutdownResponseRouter(opt *Options, hex []string, header *Header) {
	msg := PackRemoteShutdownResponseMessage(hex, header)
	log.WithFields(log.Fields{
		"id":     msg.Id,
		"gun_id": msg.GunId,
		"result": msg.Result,
		"reason": msg.Reason,
	}).Debug("[35] RemoteShutdownResponse message")

	//forward
	if opt.MessageForwarder != nil {
		//convert msg to json string bytes
		b, _ := json.Marshal(msg)
		_ = opt.MessageForwarder.Publish("35", b)
	}
}

func RemoteShutdownRequestRouter(c *gin.Context) {
	var req RemoteShutdownRequestMessage
	if c.ShouldBind(&req) == nil {
		err := SendRemoteShutdownRequest(&req)
		if err != nil {
			c.JSON(500, gin.H{"message": err})
			return
		}
	}
	c.JSON(200, gin.H{"message": "done"})
}

func TransactionRecordConfirmedRouter(c *gin.Context) {
	var req TransactionRecordConfirmedMessage
	if c.ShouldBind(&req) == nil {
		err := SendTransactionRecordConfirmed(&req)
		if err != nil {
			c.JSON(500, gin.H{"message": err})
			return
		}
	}
	c.JSON(200, gin.H{"message": "done"})
}

func TransactionRecordMessageRouter(opt *Options, raw []byte, hex []string, header *Header) {
	msg := PackTransactionRecordMessage(raw, hex, header)
	msgJson, _ := json.Marshal(msg)
	log.WithFields(log.Fields{
		"msg": string(msgJson),
	}).Debug("[3b] TransactionRecord message")

	if opt.AutoTransactionRecordConfirm {
		m := &TransactionRecordConfirmedMessage{
			Header: &Header{
				Seq:       0,
				Encrypted: false,
			},
			Id:       msg.Id,
			TradeSeq: msg.TradeSeq,
			Result:   0,
		}
		_ = SendTransactionRecordConfirmed(m)
		return
	}

	//forward
	if opt.MessageForwarder != nil {
		//convert msg to json string bytes
		b, _ := json.Marshal(msg)
		_ = opt.MessageForwarder.Publish("3b", b)
	}
}
