package apigw

import (
	"encoding/json"
	"fmt"
	coreapi "git.ecd.axway.org/apigov/apic_agents_sdk/pkg/api"
	"git.ecd.axway.org/apigov/apic_agents_sdk/pkg/apic"
	"git.ecd.axway.org/apigov/apic_agents_sdk/pkg/apic/apiserver/models/management/v1alpha1"

	//coreapi "git.ecd.axway.org/apigov/apic_agents_sdk/pkg/api"
	//"git.ecd.axway.org/apigov/apic_agents_sdk/pkg/apic/apiserver/models/management/v1alpha1"
	//corecfg "git.ecd.axway.org/apigov/apic_agents_sdk/pkg/config"
	"git.ecd.axway.org/apigov/apic_agents_sdk/pkg/transaction"
	coreerrors "git.ecd.axway.org/apigov/apic_agents_sdk/pkg/util/errors"
	//"git.ecd.axway.org/apigov/service-mesh-agent/pkg/apicauth"
	"github.com/rathnapandi/spring-traceability-agent/pkg/agent/config"
	"net/http"
	"strconv"
	"time"

	"regexp"

	"github.com/elastic/beats/v7/libbeat/logp"
	"github.com/elastic/beats/v7/libbeat/publisher"
	//"github.com/tidwall/gjson"
)

var debugf = logp.MakeDebug("apigw")

//EventProcessor -
type EventProcessor struct {
	tenantID      string
	deployment    string
	environment   string
	environmentID string
	teamID        string
	url           string
	authToken     string
	maxRetries    int
	//	v7Client                      *V7HTTPClient
	//	apiWatcher                    apimanager.APIWatcher
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

	//if agentConfig.Gateway.EnableAPICalls {
	//	ep.v7Client = NewV7Client(agentConfig.Gateway)
	//}

	//ep.prepareTransactionLogEventMap = map[string](func(transactionLeg Transaction, correlationID string) transaction.LogEvent){
	//	httpKey: ep.prepareHTTPTransactionLogEvent,
	//	//jmsKey:  ep.prepareJMSTransactionLogEvent,
	//}
	//tokenURL := agentConfig.Central.GetAuthConfig().GetTokenURL()
	//aud := agentConfig.Central.GetAuthConfig().GetAudience()
	//priKey := agentConfig.Central.GetAuthConfig().GetPrivateKey()
	//pubKey := agentConfig.Central.GetAuthConfig().GetPublicKey()
	//keyPwd := agentConfig.Central.GetAuthConfig().GetKeyPassword()
	//clientID := agentConfig.Central.GetAuthConfig().GetClientID()
	//authTimeout := agentConfig.Central.GetAuthConfig().GetTimeout()
	//platformTokenGetter := &platformTokenGetter{
	//	requester: apicauth.NewPlatformTokenGetter(priKey, pubKey, keyPwd, tokenURL, aud, clientID, authTimeout),
	//}
	//client := coreapi.NewClient(agentConfig.Central.GetTLSConfig(), agentConfig.Central.GetProxyURL())
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
	externalAPIID := getAPIServiceByExternalAPIID(springLogEntry)

	newEvents, err = p.processTransactions(newEvents, event, springLogEntry, externalAPIID)
	if err != nil {
		trxnErr := coreerrors.Wrap(ErrTrxnDataProcess, err.Error())
		logp.Error(trxnErr)
		return newEvents, trxnErr
	}

	return newEvents, nil
}

//
//func (p *EventProcessor) checkEventDetails(newEvents []publisher.Event, event publisher.Event) ([]publisher.Event, error) {
//	debugMsg := "event received without CorrelationID, throwning it out"
//
//	// Check if this event has been processed
//	if _, ok := event.Content.Meta[condorKey]; ok {
//		if retries, _ := event.Content.Meta[retriesKey].(int); retries < p.maxRetries {
//			debugMsg = "event received without CorrelationID, it seems to have been processed returning it for retry"
//			event.Content.Meta[retriesKey] = retries + 1
//			newEvents = append(newEvents, event)
//		} else {
//			debugMsg = "event received without CorrelationID, it has reached retry limit"
//		}
//	}
//
//	debugf(debugMsg)
//
//	return newEvents, nil
//}
//
//func (p *EventProcessor) getTransactionDetails(event publisher.Event, eventLogEntry EventLogEntry, eventMsg string) ([]Transaction, error) {
//	var err error
//	var v7transactionDetails []Transaction
//	// var protocolDetails []byte
//	for index, leg := range eventLogEntry.Legs {
//		// Add the TransactionDetails from the log event to v7transactionDetails
//		newTransaction := Transaction{Details: leg}
//		newTransaction.Details.Protocol = []byte(gjson.Get(eventMsg, fmt.Sprintf("legs.%d", index)).String())
//		v7transactionDetails = append(v7transactionDetails, newTransaction)
//	}
//
//	// Using the event query the v7 API for transaction details
//	sourceInstance := p.getSourceInstance(event)
//	if p.v7Client != nil && sourceInstance != "" {
//		// Get the TransactionDetails from v7
//		v7transactionDetails, err = p.v7Client.getV7LinkedTransactions(sourceInstance, eventLogEntry.CorrelationID)
//		if err != nil {
//			logp.Error(coreerrors.Wrap(ErrTrxnDataGet, err.Error()))
//		}
//	}
//
//	// Cast the details from the log event without the v7 query
//	for index, v7transactionDetail := range v7transactionDetails {
//		switch v7transactionDetail.Details.Type {
//		case httpKey:
//			var httpDetail HTTPTransactionDetail
//			err = json.Unmarshal(v7transactionDetail.Details.Protocol.([]byte), &httpDetail)
//			if err != nil {
//				err = coreerrors.Wrap(ErrProtocolStructure, err.Error()).FormatError(strings.ToUpper(httpKey), v7transactionDetail.Details.Protocol)
//				logp.Error(err)
//			}
//			v7transactionDetails[index].Details.Protocol = httpDetail
//		case jmsKey:
//			var jmsDetail JMSTransactionDetail
//			err = json.Unmarshal(v7transactionDetail.Details.Protocol.([]byte), &jmsDetail)
//			if err != nil {
//				err = coreerrors.Wrap(ErrProtocolStructure, err.Error()).FormatError(strings.ToUpper(jmsKey), v7transactionDetail.Details.Protocol)
//				logp.Error(err)
//			}
//			v7transactionDetails[index].Details.Protocol = jmsDetail
//		}
//	}
//
//	return v7transactionDetails, err
//}
//
//func (p *EventProcessor) getSourceInstance(event publisher.Event) string {
//	fileState, ok := event.Content.Private.(fbFile.State)
//	if ok {
//		fileName := fileState.Fileinfo.Name()
//		fileSegs := eventFileRegEx.FindSubmatch([]byte(fileName))
//		if fileSegs != nil && len(fileSegs) > 1 {
//			return "instance-" + string(fileSegs[1])
//		}
//	}
//	return ""
//}
//
func (p *EventProcessor) processTransactions(newEvents []publisher.Event, origLogEvent publisher.Event, springLogEntry SpringLogEntry, externalAPIID string) ([]publisher.Event, error) {
	// Iterate over all transactions, creating Condor events for each
	//for _, leg := range v7transaction {
	//	debugf("Transaction Log Event: %+v", v7transaction)
	//	if _, ok := p.prepareTransactionLogEventMap[leg.Details.Type]; !ok {
	//		debugf("Can't handle transation log events of type %s", leg.Details.Type)
	//		continue
	//	}
	//	transLogEventForLeg, err := p.createCondorEvent(origLogEvent, p.prepareTransactionLogEventMap[leg.Details.Type](leg, eventLogEntry.CorrelationID))
	//	if err != nil {
	//		return newEvents, err // createCondorEvent raises an error code
	//	}
	//	if transLogEventForLeg != nil {
	//		newEvents = append(newEvents, *transLogEventForLeg)
	//	}
	//}

	datetime, err := time.Parse(time.RFC3339, springLogEntry.Timestamp)
	if err != nil {
		panic(err)
	}
	transactionDetail := TransactionDetail{

		Timestamp:   datetime.Unix(),
		ServiceName: springLogEntry.Name,
	}
	transSummaryLogEvent, err := p.createCondorEvent(origLogEvent, p.prepareTransactionLogSummary(transactionDetail, springLogEntry))
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
func (p *EventProcessor) prepareTransactionLogSummary(transactionLegDetail TransactionDetail, springLogEntry SpringLogEntry) transaction.LogEvent {
	// Cast the Protocol details appropriately, assuming the Summary is of type HTTP
	var httpDetails HTTPTransactionDetail
	if transactionLegDetail.Protocol != nil {
		httpDetails = transactionLegDetail.Protocol.(HTTPTransactionDetail)
	}

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
			//Proxy: &transaction.Proxy{
			//	Name:     pInfo.ProxyName,
			//	ID:       transaction.FormatProxyID(pInfo.ProxyID),
			//	Revision: 1,
			//},
		},
	}
	//if len(eventLogEntry.ServiceContexts) > 0 {
	//	appName := eventLogEntry.ServiceContexts[0].App
	//	transSum.TransactionSummary.Application = new(transaction.Application)
	//
	//	// Add the V7 Application ID, with prefix, and Name to the event
	//	//cachedApp := p.apiWatcher.GetCachedApplicationByName(appName)
	//	//if cachedApp != nil {
	//	//	transSum.TransactionSummary.Application.ID = transaction.FormatApplicationID(cachedApp.ID)
	//	//}
	//	transSum.TransactionSummary.Application.Name = appName
	//}
	return transSum
}

func getAPIServiceByExternalAPIID(springLogEntry SpringLogEntry) string {

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
		return externalAPIID
		// fmt.Println(attributes.["externalAPIID"])
	}
	return ""

}
