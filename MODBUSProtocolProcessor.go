package modbusx

// PLEASE DO NOT MODIFY THIS FILE, ALL CHANGES WILL BE LOST.
// THE FILE IS REGENERATED AUTOMATICALLY, ALL CHANGES WILL BE LOST.

// 20201105-153128 {468bebe796308df98cc16adaa99d82353ace864a}

// Copyright (c) 2020 Leonid Titov, all rights reserved

import (
	"1axk.com/latitov/daob"
	"context"
	"fmt"
	"log"
	"time"
)

type pp_function_type (func(
	ctx_request context.Context,
	ch_o1_dmap chan *daob.DMAP,
	msg_dmap *daob.DMAP,
	MODBUS_RR_DMID int,
	req_frame *Frame,
	debug_mode bool,
) ([]byte, *Exception))

type msg_reg_PP_func struct {
}

// MODBUSProtocolProcessor()
func MODBUSProtocolProcessor(
	mictx context.Context,
	mictx_cancel func(),
	ch_i1 chan *Frame,
	ch_o1_dmap chan *daob.DMAP,
	MODBUS_RR_DMID int,
	debug_mode bool,
) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("modbusx.MODBUSProtocolProcessor: %v", err)
			log.Printf("modbusx.MODBUSProtocolProcessor critical error, about to cancel the Main Interconnect context")
			mictx_cancel()
		}
	}()
	pp_function := [256]pp_function_type{}

	// Add default functions.
	pp_function[1] = PP_ReadCoils
	pp_function[2] = PP_ReadDiscreteInputs
	pp_function[3] = PP_ReadRegisters
	pp_function[4] = PP_ReadRegisters
	pp_function[5] = PP_WriteSingleCoil
	pp_function[6] = PP_WriteHoldingRegister
	pp_function[15] = PP_WriteMultipleCoils
	pp_function[16] = PP_WriteHoldingRegisters

	log.Println("modbusx.MODBUSProtocolProcessor: Launching the main loop.")

	msg_dmap := &daob.DMAP{
		ChDown:    make(chan *daob.DMAP, 1), // make it bufferable!
		ReactMeta: map[string]interface{}{"modbusx": true},
	}

	for {
		select {
		case req_frame := <-ch_i1:

			// _this_ time to process the request
			ctx_request, cancel := context.WithTimeout(mictx, 100*time.Millisecond)

			resp_frame := req_frame.Copy()

			fN := req_frame.GetFunction()

			var exception *Exception
			var data []byte
			//
			if pp_function[fN] != nil {

				// call the MODBUS function handler
				data, exception = pp_function[fN](
					ctx_request,
					ch_o1_dmap,
					msg_dmap,
					MODBUS_RR_DMID,
					req_frame,
					debug_mode,
				)

				resp_frame.SetData(data)

			} else {
				exception = &IllegalFunction
			}

			if exception != &Success {
				resp_frame.SetException(exception)
			}

			select {
			case <-ctx_request.Done():
				return
			case req_frame.GetChDown() <- resp_frame:
			}

			cancel() // don't forget this!!
		case <-mictx.Done():
			return
		}
	}
}

// RegisterFunctionHandler()	override the default behavior for a given Modbus function in the Protocol Processor
func RegisterFunctionHandler(
	ch_PP chan msg_reg_PP_func,
	funcCode uint8,
	function pp_function_type,
) {
}

// PP_ReadCoils	function 1	bool	reads coils from internal memory.
func PP_ReadCoils(
	ctx_request context.Context,
	ch_o1_dmap chan *daob.DMAP,
	msg_dmap *daob.DMAP,
	MODBUS_RR_DMID int,
	req_frame *Frame,
	debug_mode bool,
) (
	[]byte,
	*Exception,
) {
	// ...
	return []byte{}, &IllegalFunction
}

// PP_ReadDiscreteInputs	function 2	bool	reads discrete inputs from internal memory.
func PP_ReadDiscreteInputs(
	ctx_request context.Context,
	ch_o1_dmap chan *daob.DMAP,
	msg_dmap *daob.DMAP,
	MODBUS_RR_DMID int,
	req_frame *Frame,
	debug_mode bool,
) (
	[]byte,
	*Exception,
) {
	// ...
	return []byte{}, &IllegalFunction
}

// PP_ReadRegisters	function 3,4	uint16	reads holding registers from internal memory.
func PP_ReadRegisters(
	ctx_request context.Context,
	ch_o1_dmap chan *daob.DMAP,
	msg_dmap *daob.DMAP,
	MODBUS_RR_DMID int,
	req_frame *Frame,
	debug_mode bool,
) (
	data []byte,
	exception *Exception,
) {
	var register, numRegs, beforeRegister int
	defer func() {
		if err := recover(); err != nil {
			if debug_mode {
				log.Printf("modbusx.MODBUSProtocolProcessor, function 3,4: R=%v-%v: %v", register, beforeRegister-1, err)
			}
			exception = &IllegalDataAddress
			return
		}
	}()
	register, numRegs, beforeRegister = registerAddressAndNumber(req_frame)

	values := make([]uint16, numRegs)

	for i, j := register, 0; i < beforeRegister; i++ {

		msg_dmap.Command = 1
		msg_dmap.DMID = MODBUS_RR_DMID
		msg_dmap.Index1 = i + 1 // THIS BECAUSE Index1 counts from 1 by DMAP spec, not because of MODBUS numbering!!

		err := msg_dmap.Exchange(ch_o1_dmap, ctx_request)
		if err != nil {
			panic(fmt.Errorf("!!!-1:  at \"err := msg_dmap. Exchange(ch_o1_dmap, ctx_request)\":   %v", err))
		}

		switch v := msg_dmap.V.(type) {
		case uint8:
			values[j] = uint16(v) // synonym for byte
		case uint16:
			values[j] = uint16(v)
		case uint32:
			values[j] = uint16(v)
		case uint64:
			values[j] = uint16(v)
		case int8:
			values[j] = uint16(v)
		case int16:
			values[j] = uint16(v)
		case int32:
			values[j] = uint16(v)
		case int64:
			values[j] = uint16(v)
		case float32:
			values[j] = uint16(v)
		case float64:
			values[j] = uint16(v)
		case int:
			values[j] = uint16(v)
		case uint:
			values[j] = uint16(v)
		default:
			panic(fmt.Errorf("Unsupported type in msg_dmap.V: %T", msg_dmap.V))
		}
		j++
	}

	data = append(data, byte(numRegs*2))
	data = append(data, Uint16ToBytes(values)...)

	exception = &Success
	return
}

// PP_WriteSingleCoil	function 5	bool	write a coil to internal memory.
func PP_WriteSingleCoil(
	ctx_request context.Context,
	ch_o1_dmap chan *daob.DMAP,
	msg_dmap *daob.DMAP,
	MODBUS_RR_DMID int,
	req_frame *Frame,
	debug_mode bool,
) (
	[]byte,
	*Exception,
) {
	// ...
	return []byte{}, &IllegalFunction
}

// PP_WriteHoldingRegister	function 6	uint16	write a holding register to internal memory.
func PP_WriteHoldingRegister(
	ctx_request context.Context,
	ch_o1_dmap chan *daob.DMAP,
	msg_dmap *daob.DMAP,
	MODBUS_RR_DMID int,
	req_frame *Frame,
	debug_mode bool,
) (
	data []byte,
	exception *Exception,
) {
	var (
		register int
		value    uint16
	)
	defer func() {
		if err := recover(); err != nil {
			if debug_mode {
				log.Printf("modbusx.MODBUSProtocolProcessor, function 6: R=%v: %v", register, err)
			}
			exception = &IllegalDataAddress
			return
		}
	}()
	register, value = registerAddressAndValue(req_frame)

	msg_dmap.Command = 1
	msg_dmap.DMID = MODBUS_RR_DMID
	msg_dmap.Index1 = register + 1 // THIS BECAUSE Index1 counts from 1 by DMAP spec, not because of MODBUS numbering!!

	err := msg_dmap.Exchange(ch_o1_dmap, ctx_request)
	if err != nil {
		panic(fmt.Errorf("!!!-3:  at \"err := msg_dmap. Exchange(ch_o1_dmap, ctx_request)\":   %v", err))
	}

	err = msg_dmap.SetV_TypeAware(value)
	if err != nil {
		panic(fmt.Errorf("!!!-4:  at \"err = msg_dmap. SetV_TypeAware( value)\":   %v", err))
	}

	msg_dmap.Command = 2
	msg_dmap.DMID = MODBUS_RR_DMID
	msg_dmap.Index1 = register + 1 // THIS BECAUSE Index1 counts from 1 by DMAP spec, not because of MODBUS numbering!!
	msg_dmap.Misc = nil            // ATTENTION; if err != nil { panic(fmt.Errorf("!!!-5: MAKE SURE TO SET HERE NIL IF YOU DON'T WANNA CHANGE IT!!! at \"msg_dmap.Misc	= nil	// ATTENTION\":   %v", err)) }

	err = msg_dmap.Exchange(ch_o1_dmap, ctx_request)
	if err != nil {
		panic(fmt.Errorf("!!!-6:  at \"err = msg_dmap. Exchange(ch_o1_dmap, ctx_request)\":   %v", err))
	}

	data = append(data, req_frame.GetData()[0:4]...)

	exception = &Success
	return
}
