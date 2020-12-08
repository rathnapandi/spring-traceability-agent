package apigw

import (
	"encoding/json"
	"fmt"
)

// TransactionDetail - transaction detail from the API
type TransactionDetail struct {
	Leg           int         `json:"leg"`
	Timestamp     int64       `json:"timestamp"`
	Duration      int         `json:"duration"`
	CorrelationID string      `json:"correlationId"`
	ServiceName   string      `json:"serviceName"`
	Subject       string      `json:"subject"`
	Operation     string      `json:"operation"`
	Type          string      `json:"type"`
	FinalStatus   interface{} `json:"finalStatus"`
	Sslsubject    interface{} `json:"sslsubject"`
	Protocol      interface{}
}

// HTTPTransactionDetail - HTTP specific transaction details
type HTTPTransactionDetail struct {
	URI           string      `json:"uri"`
	Status        int         `json:"status"`
	Statustext    string      `json:"statustext"`
	Method        string      `json:"method"`
	Vhost         interface{} `json:"vhost"`
	WafStatus     int         `json:"wafStatus"`
	BytesSent     int         `json:"bytesSent"`
	BytesReceived int         `json:"bytesReceived"`
	RemoteName    string      `json:"remoteName"`
	RemoteAddr    string      `json:"remoteAddr"`
	LocalAddr     string      `json:"localAddr"`
	RemotePort    string      `json:"remotePort"`
	LocalPort     string      `json:"localPort"`
}

// Transaction - transaction info gathered from the api calls
type Transaction struct {
	Details  TransactionDetail   `json:"details"`
	Rheaders []map[string]string `json:"rheaders"`
	Sheaders []map[string]string `json:"sheaders"`
}

// UnmarshalJSON - custom unmarshal for Transaction
func (t *Transaction) UnmarshalJSON(b []byte) (err error) {
	// Create an intermittent type to unmarshal the base attributes
	type TransactionAlias Transaction
	var tranAlias TransactionAlias

	err = json.Unmarshal(b, &tranAlias)
	if err == nil {
		t.Details = tranAlias.Details
		t.Rheaders = tranAlias.Rheaders
		t.Sheaders = tranAlias.Sheaders
		return
	}
	err = nil

	// error was hit, unmarshall to a map[string]interface{} and check the rheaders/sheaders values
	jsonMap := make(map[string]interface{})
	json.Unmarshal(b, &jsonMap)

	tranAlias.Rheaders = make([]map[string]string, 0)
	if jsonMap["rheaders"] != nil {
		b2 := []byte(fmt.Sprintf("%v", jsonMap["sheaders"]))
		json.Unmarshal(b2, &tranAlias.Sheaders)
	}

	tranAlias.Sheaders = make([]map[string]string, 0)
	if jsonMap["sheaders"] != nil {
		b2 := []byte(fmt.Sprintf("%v", jsonMap["sheaders"]))
		json.Unmarshal(b2, &tranAlias.Sheaders)
	}

	t.Details = tranAlias.Details
	t.Rheaders = tranAlias.Rheaders
	t.Sheaders = tranAlias.Sheaders
	return
}

// ServiceContext - Service Context from the EventLogEntry
type ServiceContext struct {
	App      string `json:"app"`
	Client   string `json:"client"`
	Duration int    `json:"duration"`
	Method   string `json:"method"`
	Monitor  bool   `json:"monitor"`
	Org      string `json:"org"`
	Service  string `json:"service"`
	Status   string `json:"status"`
}

// EventLogEntry - log entry created when an API in v7 is used
//type EventLogEntry struct {
//	CorrelationID   string              `json:"correlationId"`
//	CustomMsgAtts   struct{}            `json:"customMsgAtts"`
//	Duration        int                 `json:"duration"`
//	Legs            []TransactionDetail `json:"legs"`
//	Path            string              `json:"path"`
//	Protocol        string              `json:"protocol"`
//	ProtocolSrc     string              `json:"protocolSrc"`
//	ServiceContexts []ServiceContext    `json:"serviceContexts"`
//	Status          string              `json:"status"`
//	Time            int                 `json:"time"`
//	Type            string              `json:"type"`
//}

type SpringLogEntry struct {
	Timestamp  string `json:"@timestamp"`
	Level      string
	LoggerName string `json:"logger_name"`
	TraceId    string `json:"@trace_id"`
	SpanId     string `json:"@span_id"`
	Group      string
	Name       string
	Version    string
	Message    string
}
