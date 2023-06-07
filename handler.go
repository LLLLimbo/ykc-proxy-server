package main

import log "github.com/sirupsen/logrus"

func ResponseToBillingModelVerification(req *BillingModelVerificationResponseMessage) error {
	c, err := GetClient(req.Id)
	if err != nil {
		return err
	}
	resp := PackBillingModelVerificationResponseMessage(req)
	_, _ = c.Write(resp)
	log.WithFields(log.Fields{
		"id":       req.Id,
		"response": BytesToHex(resp),
	}).Debug("[06] BillingModelVerificationResponse message sent")
	return nil
}

func ResponseToVerification(req *VerificationResponseMessage) error {
	c, err := GetClient(req.Id)
	if err != nil {
		return err
	}
	resp := PackVerificationResponseMessage(req)
	_, _ = c.Write(resp)
	log.WithFields(log.Fields{
		"id":       req.Id,
		"response": BytesToHex(resp),
	}).Debug("[02] VerificationResponse message sent")
	return nil
}

func ResponseToHeartbeat(req *HeartbeatResponseMessage) error {
	c, err := GetClient(req.Id)
	if err != nil {
		return err
	}
	resp := PackHeartbeatResponseMessage(req)
	_, _ = c.Write(resp)
	log.WithFields(log.Fields{
		"id":       req.Id,
		"response": BytesToHex(resp),
	}).Debug("[04] HeartbeatResponse message sent")
	return nil
}

func SendRemoteBootstrapRequest(req *RemoteBootstrapRequestMessage) error {
	c, err := GetClient(req.Id)
	if err != nil {
		return err
	}
	resp := PackRemoteBootstrapRequestMessage(req)
	_, _ = c.Write(resp)
	log.WithFields(log.Fields{
		"id":      req.Id,
		"request": BytesToHex(resp),
	}).Debug("[34] RemoteBootstrapRequest message sent")
	return nil
}

func SendRemoteShutdownRequest(req *RemoteShutdownRequestMessage) error {
	c, err := GetClient(req.Id)
	if err != nil {
		return err
	}
	resp := PackRemoteShutdownRequestMessage(req)
	_, _ = c.Write(resp)
	log.WithFields(log.Fields{
		"id":      req.Id,
		"request": BytesToHex(resp),
	}).Debug("[36] RemoteShutdownRequest message sent")
	return nil
}

func SendTransactionRecordConfirmed(req *TransactionRecordConfirmedMessage) error {
	c, err := GetClient(req.Id)
	if err != nil {
		return err
	}
	resp := PackTransactionRecordConfirmedMessage(req)
	_, _ = c.Write(resp)
	log.WithFields(log.Fields{
		"id":      req.Id,
		"request": BytesToHex(resp),
	}).Debug("[40] TransactionRecordConfirmed message sent")
	return nil
}

func SendRemoteRebootRequest(req *RemoteRebootRequestMessage) error {
	c, err := GetClient(req.Id)
	if err != nil {
		return err
	}
	resp := PackRemoteRebootRequestMessage(req)
	_, _ = c.Write(resp)
	log.WithFields(log.Fields{
		"id":      req.Id,
		"request": BytesToHex(resp),
	}).Debug("[92] RemoteRebootRequest message sent")
	return nil
}
