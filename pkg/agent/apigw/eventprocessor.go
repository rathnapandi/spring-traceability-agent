package apigw

import (
	"encoding/json"
	"fmt"
	coreapi "git.ecd.axway.org/apigov/apic_agents_sdk/pkg/api"
	"git.ecd.axway.org/apigov/apic_agents_sdk/pkg/apic"
	"git.ecd.axway.org/apigov/apic_agents_sdk/pkg/apic/apiserver/models/management/v1alpha1"
	"git.ecd.axway.org/apigov/apic_agents_sdk/pkg/transaction"
	coreerrors "git.ecd.axway.org/apigov/apic_agents_sdk/pkg/util/errors"
	"github.com/elastic/beats/v7/libbeat/logp"
	"github.com/elastic/beats/v7/libbeat/publisher"
	"github.com/rathnapandi/spring-traceability-agent/pkg/agent/config"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

var debugf = logp.MakeDebug("apigw")

//EventProcessor -
type EventProcessor struct {
	tenantID                      string
	deployment                    string
	environment                   string
	environmentID                 string
	teamID                        string
	url                           string
	authToken                     string
	maxRetries                    int
	eventGenerator                transaction.EventGenerator
	prepareTransactionLogEventMap map[string](func(transactionLeg Transaction, correlationID string) transaction.LogEvent)
}

var eventFileRegEx = regexp.MustCompile(`group-[0-9]*_instance-([0-9]*).*`)

const (
	condorKey  = "condor"
	retriesKey = "retries"
	httpKey    = "http"
	//jmsKey     = "jms"
)

var client apic.Client
var environmentURL string

// New - return a new EventProcessor
func New(agentConfig config.AgentConfig, maxRetries int, apicClient apic.Client) *EventProcessor {
	client = apicClient
	ep := &EventProcessor{
		tenantID:      agentConfig.Central.GetTenantID(),
		deployment:    agentConfig.Central.GetAPICDeployment(),
		environment:   agentConfig.Central.GetEnvironmentName(),
		environmentID: agentConfig.Central.GetEnvironmentID(),
		teamID:        agentConfig.Central.GetTeamID(),
		//	apiWatcher:     apimanager.GetWatcher(),
		maxRetries:     maxRetries,
		eventGenerator: transaction.NewEventGenerator(),
	}
	//environment = agentConfig.Central.GetEnvironmentName()
	environmentURL = agentConfig.Central.GetEnvironmentURL()
	debugf("Event Processor Created with EnvironmentID: %s", ep.environmentID)
	return ep
}

// Process - Process the log file, waiting for events
func (p *EventProcessor) Process(events []publisher.Event) []publisher.Event {

	newEvents := make([]publisher.Event, 0)
	for _, event := range events {
		newEvents, _ = p.ProcessEvent(newEvents, event)
	}
	return newEvents
}

//ProcessEvent - Process an event from the log file
func (p *EventProcessor) ProcessEvent(newEvents []publisher.Event, event publisher.Event) ([]publisher.Event, error) {
	// Get the message from the log file
	eventMsgFieldVal, err := event.Content.Fields.GetValue("message")
	if err != nil {
		return newEvents, coreerrors.Wrap(ErrEventNoMsg, err.Error()).FormatError(event)
	}
	eventMsg, ok := eventMsgFieldVal.(string)
	if !ok {
		return newEvents, nil
	}

	// Unmarshal the message into a Log Entry Event
	var springLogEntry SpringLogEntry
	err = json.Unmarshal([]byte(eventMsg), &springLogEntry)
	if err != nil {
		msgErr := coreerrors.Wrap(ErrEventMsgStructure, err.Error()).FormatError(eventMsg)
		logp.Error(msgErr)
		return newEvents, msgErr
	}

	fmt.Printf("*** JSON : %v", springLogEntry)
	externalAPIID, name := getAPIServiceByExternalAPIID(springLogEntry)

	newEvents, err = p.processTransactions(newEvents, event, springLogEntry, externalAPIID, name)
	if err != nil {
		trxnErr := coreerrors.Wrap(ErrTrxnDataProcess, err.Error())
		logp.Error(trxnErr)
		return newEvents, trxnErr
	}

	return newEvents, nil
}

func (p *EventProcessor) processTransactions(newEvents []publisher.Event, origLogEvent publisher.Event, springLogEntry SpringLogEntry, externalAPIID, name string) ([]publisher.Event, error) {

	datetime, err := time.Parse(time.RFC3339, springLogEntry.Timestamp)
	if err != nil {
		panic(err)
	}
	transactionDetail := TransactionDetail{

		Timestamp:   datetime.Unix(),
		ServiceName: springLogEntry.Name,
	}
	transSummaryLogEvent, err := p.createCondorEvent(origLogEvent, p.prepareTransactionLogSummary(transactionDetail, springLogEntry, externalAPIID, name))
	if err != nil {
		return newEvents, err // createCondorEvent raises an error code
	}
	if transSummaryLogEvent != nil {
		debugf("Summary from Log Events: %+v", transactionDetail)
		newEvents = append(newEvents, *transSummaryLogEvent)
	}

	return newEvents, nil
}

func (p *EventProcessor) createCondorEvent(originalLogEvent publisher.Event, transactionEvent transaction.LogEvent) (*publisher.Event, error) {
	// Create the beat event then wrap that in the publisher event for Condor

	// Add a Retry count to Meta
	if originalLogEvent.Content.Meta == nil {
		originalLogEvent.Content.Meta = make(map[string]interface{}, 0)
	}
	originalLogEvent.Content.Meta[retriesKey] = 0
	originalLogEvent.Content.Meta[condorKey] = true

	beatEvent, err := p.eventGenerator.CreateEvent(transactionEvent, originalLogEvent.Content.Timestamp, originalLogEvent.Content.Meta, originalLogEvent.Content.Fields, originalLogEvent.Content.Private)

	if err != nil {
		return nil, coreerrors.Wrap(ErrCreateCondorEvent, err.Error())
	}
	event := publisher.Event{
		Content: beatEvent,
		Flags:   originalLogEvent.Flags,
	}
	return &event, nil
}

//
//func prepareHeaders(transactionHeaders []map[string]string) (string, string) {
//	// Add all headers to a single header object
//	headers := make(map[string]string, 0)
//	userAgent := ""
//
//	// skip processing the headers
//	if transactionHeaders != nil {
//		for _, transactionHeader := range transactionHeaders {
//			for key, value := range transactionHeader {
//				headers[key] = value
//				if key == "User-Agent" {
//					userAgent = value
//				}
//			}
//		}
//	}
//
//	// Marshall the headers to JSON
//	headerBytes, err := json.Marshal(headers)
//	if err != nil {
//		logp.Error(coreerrors.Wrap(ErrTrxnHeaders, err.Error()))
//	}
//
//	return string(headerBytes), userAgent
//}
//
//func (p *EventProcessor) prepareHTTPTransactionLogEvent(transactionLeg Transaction, correlationID string) transaction.LogEvent {
//	// Cast the Protocol details appropriately
//	var httpDetails HTTPTransactionDetail
//	if transactionLeg.Details.Protocol != nil {
//		httpDetails = transactionLeg.Details.Protocol.(HTTPTransactionDetail)
//	}
//
//	// Create the LogEvent for the Transaction Leg
//	var requestHeader, responseHeader, direction, userAgent, parentID string
//	var remotePort, localPort int
//	var remoteName, remoteAddr, localAddr string
//	if transactionLeg.Details.Leg == 0 {
//		direction = "inbound"
//
//		requestHeader, userAgent = prepareHeaders(transactionLeg.Rheaders)
//		responseHeader, _ = prepareHeaders(transactionLeg.Sheaders)
//
//		remotePort, _ = strconv.Atoi(httpDetails.RemotePort)
//		localPort, _ = strconv.Atoi(httpDetails.LocalPort)
//
//		remoteName = httpDetails.RemoteName
//		remoteAddr = httpDetails.RemoteAddr
//		localAddr = httpDetails.LocalAddr
//	} else {
//		parentID = "leg" + strconv.Itoa(transactionLeg.Details.Leg-1)
//		direction = "outbound"
//
//		requestHeader, userAgent = prepareHeaders(transactionLeg.Sheaders)
//		responseHeader, _ = prepareHeaders(transactionLeg.Rheaders)
//
//		remoteAddr = httpDetails.LocalAddr
//		localAddr = httpDetails.RemoteAddr
//
//		remotePort, _ = strconv.Atoi(httpDetails.LocalPort)
//		localPort, _ = strconv.Atoi(httpDetails.RemotePort)
//	}
//
//	transStatus := "Fail"
//	if httpDetails.Status >= http.StatusOK && httpDetails.Status < http.StatusBadRequest {
//		transStatus = "Pass"
//	}
//
//	transLogEventLeg := transaction.LogEvent{
//		Version:           "1.0",
//		Stamp:             transactionLeg.Details.Timestamp,
//		TransactionID:     correlationID,
//		TenantID:          p.tenantID,
//		TrcbltPartitionID: p.tenantID,
//		APICDeployment:    p.deployment,
//		EnvironmentID:     p.environmentID,
//		Environment:       p.environment,
//		Type:              "transactionEvent",
//		TransactionEvent: &transaction.Event{
//			ID:          "leg" + strconv.Itoa(transactionLeg.Details.Leg),
//			ParentID:    parentID,
//			Source:      remoteAddr,
//			Destination: localAddr,
//			Duration:    transactionLeg.Details.Duration,
//			Direction:   direction,
//			Status:      transStatus,
//			Protocol: &transaction.Protocol{
//				Type:            httpKey,
//				URI:             httpDetails.URI,
//				Args:            "",
//				Method:          httpDetails.Method,
//				Status:          httpDetails.Status,
//				StatusText:      http.StatusText(httpDetails.Status),
//				UserAgent:       userAgent,
//				Host:            httpDetails.LocalAddr,
//				Version:         "",
//				RequestHeaders:  requestHeader,
//				ResponseHeaders: responseHeader,
//				BytesSent:       httpDetails.BytesSent,
//				BytesReceived:   httpDetails.BytesReceived,
//				RemoteName:      remoteName,
//				RemoteAddr:      remoteAddr,
//				RemotePort:      remotePort,
//				LocalAddr:       localAddr,
//				LocalPort:       localPort,
//			},
//		},
//	}
//	return transLogEventLeg
//}
//
//
func (p *EventProcessor) prepareTransactionLogSummary(transactionLegDetail TransactionDetail, springLogEntry SpringLogEntry, externalAPIID, name string) transaction.LogEvent {
	// Cast the Protocol details appropriately, assuming the Summary is of type HTTP
	var httpDetails HTTPTransactionDetail
	if transactionLegDetail.Protocol != nil {
		httpDetails = transactionLegDetail.Protocol.(HTTPTransactionDetail)
	}

	httpDetails.Status = http.StatusOK
	httpDetails.Method = http.MethodGet
	httpDetails.URI = "/healthcheck"
	httpDetails.LocalAddr = "127.0.0.1"
	transactionLegDetail.Duration = 0

	// Create the LogEvent for the Transaction Summary
	transSummaryStatus := "Unknown"
	if httpDetails.Status >= http.StatusOK && httpDetails.Status < http.StatusBadRequest {
		transSummaryStatus = "Success"
	} else if httpDetails.Status >= http.StatusBadRequest && httpDetails.Status < http.StatusInternalServerError {
		transSummaryStatus = "Failure"
	} else if httpDetails.Status >= http.StatusInternalServerError && httpDetails.Status < http.StatusNetworkAuthenticationRequired {
		transSummaryStatus = "Exception"
	}
	status := strconv.Itoa(httpDetails.Status)
	transSum := transaction.LogEvent{
		Version:           "1.0",
		Stamp:             transactionLegDetail.Timestamp,
		TransactionID:     springLogEntry.TraceId,
		TenantID:          p.tenantID,
		TrcbltPartitionID: p.tenantID,
		APICDeployment:    p.deployment,
		EnvironmentID:     p.environmentID,
		Environment:       p.environment,
		Type:              "transactionSummary",
		TransactionSummary: &transaction.Summary{
			Status:       transSummaryStatus,
			StatusDetail: status,
			Duration:     transactionLegDetail.Duration,
			Team: &transaction.Team{
				ID: p.teamID,
			},
			EntryPoint: &transaction.EntryPoint{
				Type:   httpKey,
				Method: httpDetails.Method,
				Path:   httpDetails.URI,
				Host:   httpDetails.LocalAddr,
			},
			Proxy: &transaction.Proxy{
				//Name:     pInfo.ProxyName,
				//ID:       transaction.FormatProxyID(pInfo.ProxyID),
				Name:     name,
				ID:       "remoteApiId_" + externalAPIID,
				Revision: 1,
			},
		},
	}
	jsonData, _ := json.Marshal(transSum)
	fmt.Println(string(jsonData))
	return transSum
}

func getAPIServiceByExternalAPIID(springLogEntry SpringLogEntry) (string, string) {

	query := map[string]string{
		"query": "attributes.name" + "==\"" + springLogEntry.Name + "\" and attributes.version==\"" + springLogEntry.Version + "\"",
	}

	resp, err := client.ExecuteAPI(coreapi.GET, environmentURL+"/apiservices", query, nil)
	if err != nil {
		panic(err)
	}
	apiServices := make([]v1alpha1.APIService, 0)
	json.Unmarshal(resp, &apiServices)
	if len(apiServices) > 0 {
		attributes := &apiServices[0].Attributes
		externalAPIID := (*attributes)["externalAPIID"]
		name := (*attributes)["name"]
		return externalAPIID, name
	}
	return "", ""

}
