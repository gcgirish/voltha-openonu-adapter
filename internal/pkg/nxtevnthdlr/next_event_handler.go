/*
 * Copyright 2021-present Open Networking Foundation
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

//Package nxtevnthdlr serializes Flow and TP events
package nxtevnthdlr

import (
	"context"
	"fmt"
	"github.com/looplab/fsm"
	"github.com/opencord/voltha-lib-go/v5/pkg/log"
	"github.com/opencord/voltha-openonu-adapter-go/internal/pkg/avcfg"
	cmn "github.com/opencord/voltha-openonu-adapter-go/internal/pkg/common"
	ic "github.com/opencord/voltha-protos/v4/go/inter_container"
	"sync"
)

const cBasePathTechProfileKVStore = "%s/technology_profiles"

const (
	// events of NexEventHandler FSM
	nextEventHandlerEventFlowReq         = "nextEventHandlerEventFlowReq"
	nextEventHandlerEventFlowResp        = "nextEventHandlerEventFlowResp"
	nextEventHandlerEventTpDeletePending = "nextEventHandlerEventTpDeletePending"
	nextEventHandlerEventTpAddPending    = "nextEventHandlerEventTpAddPending"
	nextEventHandlerEventTpDelReq        = "nextEventHandlerEventTpDelReq"
	nextEventHandlerEventTpAddReq        = "nextEventHandlerEventTpAddReq"
	nextEventHandlerEventTpResp          = "nextEventHandlerEventTpResp"
)
const (
	// states of NexEventHandler FSM
	nextEventHandlerStIdle           = "nextEventHandlerStIdle"
	nextEventHandlerStProcessingFlow = "nextEventHandlerStProcessingFlow"
	nextEventHandlerStProcessingTp   = "nextEventHandlerStProcessingTp"
	nextEventHandlerStWaitTpDelete   = "nextEventHandlerStWaitTpDelete"
	nextEventHandlerStWaitTpAdd      = "nextEventHandlerStWaitTpAdd"
)

type FlowCb struct {
	ctx     context.Context // Flow handler context
	addFlow bool            // if true flow to be added, else removed
	flow    *cmn.UniVlanFlowParams
	errChan *chan error // channel to report the Flow handling error
}

type TpCb struct {
	ctx           context.Context         // TP handler context
	addTpInstance bool                    // if true flow to be added, else removed
	msg           *ic.InterAdapterMessage // inter adapter TP message
	errChan       *chan error             // channel to report the TP handling error
}

type NextEventHandler struct {
	onuDeviceID    string
	pDeviceHandler cmn.IdeviceHandler
	pOnuTP         *avcfg.OnuUniTechProf
	pAniFsm        *avcfg.UniPonAniConfigFsm

	ponIntfID      uint32
	onuID          uint32
	uniID          uint32
	tpID           uint32
	tpInstancePath string

	flowCbQ     []FlowCb
	flowCbQLock sync.RWMutex
	tpCbQ       []TpCb
	tpCbQLock   sync.RWMutex
	flowCbChan  chan FlowCb
	tpCbChan    chan TpCb
	flowResp    chan error
	tpResp      chan error

	stopFlowMonitoringRoutine     chan bool
	stopTpMonitoringRoutine       chan bool
	isFlowMonitoringRoutineActive bool
	isTpMonitoringRoutineActive   bool

	nxtEventHdlrLock sync.RWMutex

	fsm     *fsm.FSM
	lockFsm sync.RWMutex
}

// NewNextEventHandler initializes the Next Event Handler module
func NewNextEventHandler(ctx context.Context, aDeviceHandler cmn.IdeviceHandler, aOnuTp *avcfg.OnuUniTechProf, aAniFsm *avcfg.UniPonAniConfigFsm,
	aPonIntfID uint32, aOnuID uint32, aUniID uint32, aTpID uint32, aTpInstancePath string) *NextEventHandler {
	nxtEvntHdlr := &NextEventHandler{}
	nxtEvntHdlr.pDeviceHandler = aDeviceHandler
	nxtEvntHdlr.onuDeviceID = nxtEvntHdlr.pDeviceHandler.GetDeviceID()
	nxtEvntHdlr.pOnuTP = aOnuTp
	nxtEvntHdlr.pAniFsm = aAniFsm
	nxtEvntHdlr.ponIntfID = aPonIntfID
	nxtEvntHdlr.onuID = aOnuID
	nxtEvntHdlr.uniID = aUniID
	nxtEvntHdlr.tpID = aTpID
	nxtEvntHdlr.tpInstancePath = aTpInstancePath

	nxtEvntHdlr.stopFlowMonitoringRoutine = make(chan bool)
	nxtEvntHdlr.stopTpMonitoringRoutine = make(chan bool)
	nxtEvntHdlr.flowResp = make(chan error)
	nxtEvntHdlr.tpResp = make(chan error)

	if nxtEvntHdlr.fsm = nxtEvntHdlr.initiateNextEventHandlerFsm(ctx); nxtEvntHdlr.fsm == nil {
		logger.Errorw(ctx, "failed to initialize fsm", log.Fields{"device-id": nxtEvntHdlr.onuDeviceID})
		return nil
	}
	go nxtEvntHdlr.monitorFlowMessages(ctx)
	go nxtEvntHdlr.monitorTpMessages(ctx)

	logger.Infow(ctx, "initializing next event handler success",
		log.Fields{"device-id": nxtEvntHdlr.onuDeviceID, "tpInstPath": nxtEvntHdlr.tpInstancePath})
	return nxtEvntHdlr
}

// Exported functions

func (neh *NextEventHandler) StopNextEventHandler(ctx context.Context) {
	logger.Debugw(ctx, "StopNextEventHandler - start",
		log.Fields{"device-id": neh.onuDeviceID, "pon": neh.ponIntfID, "onuID": neh.onuID, "uniID": neh.uniID, "tpID": neh.tpID})
	neh.nxtEventHdlrLock.Lock()
	defer neh.nxtEventHdlrLock.Unlock()
	if neh.isTpMonitoringRoutineActive {
		neh.stopTpMonitoringRoutine <- true
	}
	if neh.isFlowMonitoringRoutineActive {
		neh.stopFlowMonitoringRoutine <- true
	}
	logger.Debugw(ctx, "StopNextEventHandler - end",
		log.Fields{"device-id": neh.onuDeviceID, "pon": neh.ponIntfID, "onuID": neh.onuID, "uniID": neh.uniID, "tpID": neh.tpID})
}

func (neh *NextEventHandler) HandleTpReq(ctx context.Context, cb TpCb) {
	neh.nxtEventHdlrLock.RLock()
	defer neh.nxtEventHdlrLock.RUnlock()
	// Check that the monitoring routine is active
	// If active, push the event on the channel
	// Else immediately return error in the return channel in the Cb
	// The caller will wait in the return channel for acknowledgement of the request
	if neh.isTpMonitoringRoutineActive {
		neh.tpCbChan <- cb
	} else {
		err := fmt.Errorf("tp-monitoring-routine-not-active--device-id-%v-uni-id-%d-tp-id-%v", neh.onuDeviceID, neh.uniID, neh.tpID)
		*cb.errChan <- err
	}
}

func (neh *NextEventHandler) HandleFlowReq(ctx context.Context, cb FlowCb) {
	neh.nxtEventHdlrLock.RLock()
	defer neh.nxtEventHdlrLock.RUnlock()
	// Check that the monitoring routine is active
	// If active, push the event on the channel
	// Else immediately return error in the return channel in the Cb
	// The caller will wait in the return channel for acknowledgement of the request
	if neh.isFlowMonitoringRoutineActive {
		neh.flowCbChan <- cb
	} else {
		err := fmt.Errorf("flow-monitoring-routine-not-active--device-id-%v-uni-id-%d-tp-id-%v", neh.onuDeviceID, neh.uniID, neh.tpID)
		*cb.errChan <- err
	}
}

// unexported functions

// initiateNextEventHandlerFsm initializes the Next Event Handler FSM
func (neh *NextEventHandler) initiateNextEventHandlerFsm(ctx context.Context) *fsm.FSM {

	// Next Event Handler FSM
	return fsm.NewFSM(

		nextEventHandlerStIdle,
		fsm.Events{
			{Name: nextEventHandlerEventFlowReq, Src: []string{nextEventHandlerStIdle}, Dst: nextEventHandlerStProcessingFlow},
			{Name: nextEventHandlerEventFlowResp, Src: []string{nextEventHandlerStProcessingFlow}, Dst: nextEventHandlerStIdle},
			{Name: nextEventHandlerEventTpDeletePending, Src: []string{nextEventHandlerStProcessingFlow, nextEventHandlerStProcessingTp}, Dst: nextEventHandlerStWaitTpDelete},
			{Name: nextEventHandlerEventTpAddPending, Src: []string{nextEventHandlerStProcessingFlow}, Dst: nextEventHandlerStWaitTpAdd},
			{Name: nextEventHandlerEventTpDelReq, Src: []string{nextEventHandlerStWaitTpDelete}, Dst: nextEventHandlerStProcessingTp},
			{Name: nextEventHandlerEventTpAddReq, Src: []string{nextEventHandlerStWaitTpAdd}, Dst: nextEventHandlerStProcessingTp},
			{Name: nextEventHandlerEventTpResp, Src: []string{nextEventHandlerStProcessingTp}, Dst: nextEventHandlerStIdle},
		},
		fsm.Callbacks{
			"enter_state":                               func(e *fsm.Event) { neh.logFsmStateChange(ctx, e) },
			"enter_" + nextEventHandlerStIdle:           func(e *fsm.Event) { neh.nextEventHandlerFsmHandleIdle(ctx, e) },
			"enter_" + nextEventHandlerStProcessingFlow: func(e *fsm.Event) { neh.nextEventHandlerFsmHandleProcessingFlow(ctx, e) },
			"enter_" + nextEventHandlerStWaitTpAdd:      func(e *fsm.Event) { neh.nextEventHandlerFsmWaitTpAdd(ctx, e) },
			"enter_" + nextEventHandlerStWaitTpDelete:   func(e *fsm.Event) { neh.nextEventHandlerFsmWaitTpDelete(ctx, e) },
			"enter_" + nextEventHandlerStProcessingTp:   func(e *fsm.Event) { neh.nextEventHandlerFsmHandleProcessingTp(ctx, e) },
		},
	)
}

// FSM Handlers -- start

func (neh *NextEventHandler) nextEventHandlerFsmHandleIdle(ctx context.Context, e *fsm.Event) {

	neh.flowCbQLock.RLock()
	// If we have flows, process them
	if len(neh.flowCbQ) > 0 {
		neh.lockFsm.Lock()
		go func() {
			_ = neh.fsm.Event(nextEventHandlerEventFlowReq)
		}()
		neh.lockFsm.Unlock()
	}
	neh.flowCbQLock.RUnlock()

}

func (neh *NextEventHandler) nextEventHandlerFsmHandleProcessingFlow(ctx context.Context, e *fsm.Event) {
	// The expectation is the flow is ordered.
	// Handle batch of queued flows of each kind (add or remove).
	// Each flow is processed to completion before picking up next.

	var lFlowCb []FlowCb
	neh.flowCbQLock.Lock()
	// We bunch together flow types of a particular kind (add/remove) and process them together
	if len(neh.flowCbQ) > 0 { // We have some flows to process
		// Some initial checks before we start processing the flows.
		addFlow := neh.flowCbQ[0].addFlow
		if addFlow {
			//// TODO ////
			/*
				Pseudocode
				---------
				if !tpDoneForUniTp(uni-id,tp-id) {
					1. Trigger nextEventHandlerEventTpAddPending to move to
					2. Release flowCbQLock
					3. return
				}
			*/
		} else { // Remove flow
			//// TODO ////
			/*
				Pseudocode
				---------
				// If TP is not done remove send error response on the channel for these flows and remove them from queue
				if !tpDoneForUniTp(uni-id,tp-id) {
					1. Populate collocated flow removes from neh.flowCbQ
					2. Send error on the flowCb.errChan
					3. Remove these entries from neh.flowCbQ
					4. Trigger nextEventHandlerEventFlowResp on the Fsm to go back to nextEventHandlerStIdle state
					5. return
				}
			*/
		}

		// Populate all the flows of a given kind (add/rem) to a local queue to be processed in a bunch, but serially.
		for {
			if neh.flowCbQ[0].addFlow == addFlow {
				lFlowCb = append(lFlowCb, neh.flowCbQ[0])
				neh.flowCbQ = append(neh.flowCbQ[:0], neh.flowCbQ[1:]...)
			} else {
				break
			}
		}
	}
	neh.flowCbQLock.Unlock()

	// Start processing flows of the given kind (add/rem), one by one, blocking until response is received for each flow
	for _, v := range lFlowCb {
		logger.Debugw(ctx, "handling flow", log.Fields{"flow": v.flow})
		////// TODO /////
		// Trigger vlan FSM. Pass ** neh.flowResp ** channel to monitor for response
		// Block on neh.flowResp channel for response.
		// When the response is received, pass it to lFlowCb.errChan and then proceed to next flow
	}
}

func (neh *NextEventHandler) nextEventHandlerFsmWaitTpAdd(ctx context.Context, e *fsm.Event) {
	neh.lockFsm.Lock()
	defer neh.lockFsm.Unlock()
	neh.tpCbQLock.RLock()
	defer neh.tpCbQLock.RUnlock()
	// If we have a TP Add in the queue, process it. Trigger nextEventHandlerEventTpAddReq event on the FSM
	if len(neh.tpCbQ) > 0 && neh.tpCbQ[0].addTpInstance {
		go func() {
			_ = neh.fsm.Event(nextEventHandlerEventTpAddReq)
		}()
	}
	return
}

func (neh *NextEventHandler) nextEventHandlerFsmWaitTpDelete(ctx context.Context, e *fsm.Event) {
	neh.lockFsm.Lock()
	defer neh.lockFsm.Unlock()
	neh.tpCbQLock.RLock()
	defer neh.tpCbQLock.RUnlock()
	// If we have a TP Delete in the queue, process it. Trigger nextEventHandlerEventTpDeleteReq event on the FSM
	if len(neh.tpCbQ) > 0 && !neh.tpCbQ[0].addTpInstance {
		go func() {
			_ = neh.fsm.Event(nextEventHandlerEventTpDelReq)
		}()
	}
	return
}

func (neh *NextEventHandler) nextEventHandlerFsmHandleProcessingTp(ctx context.Context, e *fsm.Event) {
	var lTpCb TpCb
	neh.tpCbQLock.RLock()
	if len(neh.tpCbQ) > 0 {
		lTpCb = neh.tpCbQ[0]
	} else {
		// nothing to handler. Should not ideally land here
		logger.Warnw(ctx, "no tp messages to process", log.Fields{"device-id": neh.onuDeviceID, "pon": neh.ponIntfID, "onuID": neh.onuID, "uniID": neh.uniID, "tpID": neh.tpID})
		go func() {
			_ = neh.fsm.Event(nextEventHandlerEventTpResp)
		}()
		neh.tpCbQLock.RUnlock()
		return
	}
	neh.tpCbQLock.RUnlock()

	if lTpCb.msg.Header.Type == ic.InterAdapterMessageType_TECH_PROFILE_DOWNLOAD_REQUEST {
		// Remove the TP Cb entry from the queue to be processed
		neh.tpCbQ = append(neh.tpCbQ[:0], neh.tpCbQ[1:]...)
		/// TODO ///
		/*
			Pseudocode
			---------
				1. Trigger vlan FSM. Pass ** neh.tpResp ** channel to monitor for response
				2. Block on neh.tpResp channel for response.
				3. When the response is received, pass it to lTpCb.errChan and then trigger event nextEventHandlerEventTpResp on the FSM
		*/
	} else if lTpCb.msg.Header.Type == ic.InterAdapterMessageType_DELETE_GEM_PORT_REQUEST || lTpCb.msg.Header.Type == ic.InterAdapterMessageType_DELETE_TCONT_REQUEST {

	}
}

// FSM Handlers -- end

// logFsmStateChange logs FSM state changes
func (neh *NextEventHandler) logFsmStateChange(ctx context.Context, e *fsm.Event) {
	logger.Debugw(ctx, "FSM state change", log.Fields{"device-id": neh.onuDeviceID,
		"event name": e.Event, "src state": e.Src, "dst state": e.Dst})
}

func (neh *NextEventHandler) monitorTpMessages(ctx context.Context) {
	neh.nxtEventHdlrLock.Lock()
	neh.isTpMonitoringRoutineActive = true
	neh.nxtEventHdlrLock.Unlock()
	for {
		select {
		case tpCb := <-neh.tpCbChan:
			neh.tpCbQLock.Lock()
			neh.tpCbQ = append(neh.tpCbQ, tpCb)
			neh.tpCbQLock.Unlock()

			neh.lockFsm.Lock()
			if tpCb.addTpInstance {
				// Push TP Add event to FSM only if it can handle it, else it will picked up later anyway
				if neh.fsm.Can(nextEventHandlerEventTpAddReq) {
					go func() {
						_ = neh.fsm.Event(nextEventHandlerEventTpAddReq)
					}()
				}
			} else {
				// Push TP Delete event to FSM only if it can handle it, else it will picked up later anyway
				if neh.fsm.Can(nextEventHandlerEventTpDelReq) {
					go func() {
						_ = neh.fsm.Event(nextEventHandlerEventTpDelReq)
					}()
				}
			}
			neh.lockFsm.Unlock()
		case <-neh.stopTpMonitoringRoutine:
			neh.isTpMonitoringRoutineActive = false // access to this should already protected by caller to stop the routine
			neh.tpCbQLock.Lock()
			neh.tpCbQ = make([]TpCb, 0) // flush tpCbQ
			neh.tpCbQLock.Unlock()
			return
		}
	}
}

func (neh *NextEventHandler) monitorFlowMessages(ctx context.Context) {
	neh.nxtEventHdlrLock.Lock()
	neh.isFlowMonitoringRoutineActive = true
	neh.nxtEventHdlrLock.Unlock()
	for {
		select {
		case flowCb := <-neh.flowCbChan:
			neh.flowCbQLock.Lock()
			neh.flowCbQ = append(neh.flowCbQ, flowCb)
			neh.flowCbQLock.Unlock()

			neh.lockFsm.Lock()
			// Push Flow event to FSM only if it can handle it, else it will picked up later anyway
			if neh.fsm.Can(nextEventHandlerEventFlowReq) {
				go func() {
					_ = neh.fsm.Event(nextEventHandlerEventFlowReq)
				}()
			}
			neh.lockFsm.Unlock()
		case <-neh.stopFlowMonitoringRoutine:
			neh.isFlowMonitoringRoutineActive = false //access to this should already protected by caller to stop the routine
			neh.flowCbQLock.Lock()
			neh.flowCbQ = make([]FlowCb, 0) // flush flowCbQ
			neh.flowCbQLock.Unlock()
			return
		}
	}
}
