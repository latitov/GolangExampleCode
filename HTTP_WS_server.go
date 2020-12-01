package main

// (c) Leonid Titov, 2019, All rights reserved

import (
	"1axk.com/latitov/daob"
	"context"
	"crypto/tls"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/websocket"
	"log"
	"math"
	"net/http"
	"net/http/pprof"
	"os"
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type reg_ws_msg struct {
	ws      *websocket.Conn
	command int
}

func HTTP_WS_server(
	mictx context.Context,
	mictx_cancel func(),
	ch_o1_S1 chan *daob.DMAP,
	ch_CPUI_req chan *CPUI_req,
) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("HTTP_WS_server panicked: %v", err)
			mictx_cancel()
		}
	}()
	log.Printf("HTTP_WS_server: Starting.\n")

	certPem, err := FSByte(false, "/_R/resources/cert.pem")
	if err != nil {
		panic(fmt.Errorf("!!!-1: Can't load certificate file at \"certPem, err	:= FSByte(false, \"/_R/resources/cert.pem\")\":   %v", err))
	}
	keyPem, err := FSByte(false, "/_R/resources/key.pem")
	if err != nil {
		panic(fmt.Errorf("!!!-2: Can't load private key file at \"keyPem, err	:= FSByte(false, \"/_R/resources/key.pem\")\":   %v", err))
	}
	cert, err := tls.X509KeyPair(certPem, keyPem)
	if err != nil {
		panic(fmt.Errorf("!!!-3: X509KeyPair() failed at \"cert, err	:= tls.X509KeyPair(certPem, keyPem)\":   %v", err))
	}
	cfg := &tls.Config{Certificates: []tls.Certificate{cert}}

	srv := &http.Server{
		Addr: ":8100",

		Handler:      nil,
		TLSConfig:    cfg,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	}
	ctx_server, cancel := context.WithCancel(mictx)
	defer cancel()

	log.Println("HTTP_WS_server: Starting the connections registrar.")

	ch_reg_ws := make(chan *reg_ws_msg)
	go func() {
		map1 := make(map[*websocket.Conn]*websocket.Conn, 50)
		for {
			select {
			case msg := <-ch_reg_ws:
				switch msg.command {
				case 1:
					map1[msg.ws] = msg.ws
				case 0:
					delete(map1, msg.ws)
				}
			case <-ctx_server.Done():
				srv.Shutdown(ctx_server)
				for _, v := range map1 {
					v.Close()
				}
				return
			}
		}
	}()

	m := http.NewServeMux()
	srv.Handler = m

	m.HandleFunc("/ws/tv1",

		func(w http.ResponseWriter, r *http.Request) {
			var conn *websocket.Conn
			defer func() {
				if err := recover(); err != nil {
					log.Printf("WAUI tv1: WebSocket handler went away: failed OR CLOSED BY CLIENT: %v", err)
				}
			}()
			var chopfail bool
			_ = chopfail

			select {
			case <-time.After(500 * time.Millisecond):
				panic(fmt.Errorf("Failed to register a handler, \"ch_reg_ws <- &reg_ws_msg{command:1, ws:conn} @!! 200\""))
			case ch_reg_ws <- &reg_ws_msg{command: 1, ws: conn}:
			}
			defer func() {
				select {
				case <-time.After(500 * time.Millisecond):
					panic(fmt.Errorf("Failed to un-register a handler, \"ch_reg_ws <- &reg_ws_msg{command:0, ws:conn} @!! 200\""))
				case ch_reg_ws <- &reg_ws_msg{command: 0, ws: conn}:
				}
				conn.Close()
			}()
			conn, err := upgrader.Upgrade(w, r, nil)
			if err != nil {
				panic(fmt.Errorf("!!!-4: WebSocket upgrade failed at \"conn, err := upgrader.Upgrade(w, r, nil)\":   %v", err))
			}
			log.Printf("WAUI tv1: New active WebSocket connection, %s\n", conn.RemoteAddr())

			msg_dmap := &daob.DMAP{
				ChDown:    make(chan *daob.DMAP, 1),
				ReactMeta: map[string]interface{}{"WAUI": true},
			}

			msg := struct {
				Release_VERTIMESTAMP string
			}{
				"20201127-150757 {97f42b625fb4f1c86f4c5b1ddbe5298f83c01663}",
			}
			log.Printf("WAUI tv1: sending the initial context reset: %+v\n", msg)
			err = conn.WriteJSON(msg)
			if err != nil {
				panic(fmt.Errorf("!!!-5:  at \"err = conn.WriteJSON(msg)\":   %v", err))
			}

			for {

				req := TV1UI_req{}
				err = conn.ReadJSON(&req)
				if err != nil {
					panic(fmt.Errorf("!!!-6: Failed at \"err = conn.ReadJSON(&req)\":   %v", err))
				}

				ctx_request, cancel := context.WithTimeout(ctx_server, 1000*time.Millisecond)

				msg_dmap.Command = 255

				err := msg_dmap.Exchange(ch_o1_S1, ctx_request)
				if err != nil {
					panic(fmt.Errorf("!!!-7:  at \"err := msg_dmap. Exchange( ch_o1_S1, ctx_request)\":   %v", err))
				}

				switch req.RequestType {
				case 1:
					t := time.Now()
					msg := struct {
						RequestType   uint
						NoReset       uint
						Modbus_S_Time string
					}{
						NoReset:       1,
						RequestType:   req.RequestType,
						Modbus_S_Time: fmt.Sprintf("%v", t.Format("2006-01-02 15:04:05.000")),
					}
					conn.WriteJSON(msg)
				case 50:
					msg := struct {
						ProgramName string
						RequestType uint
						NoReset     uint
						Data        interface{}
						DMID        int
					}{
						ProgramName: req.Key,
						RequestType: req.RequestType,
						NoReset:     1,
					}

					msg_dmap.Command = 1
					msg_dmap.DMID = MNEMONIC_DMID_MAP["Programs"]
					msg_dmap.Key = msg.ProgramName

					err = msg_dmap.Exchange(ch_o1_S1, ctx_request)
					if err != nil {
						panic(fmt.Errorf("!!!-8:  at \"err = msg_dmap. Exchange(ch_o1_S1, ctx_request)\":   %v", err))
					}

					msg.DMID = msg_dmap.DMID
					msg.Data = msg_dmap.V

					conn.WriteJSON(msg)
				}

				cancel()
			}
		},
	)
	m.HandleFunc("/ws/controlplane",

		func(w http.ResponseWriter, r *http.Request) {
			var conn *websocket.Conn
			defer func() {
				if err := recover(); err != nil {
					log.Printf("WAUI controlplane: WebSocket handler went away: failed OR CLOSED BY CLIENT: %v", err)
				}
			}()
			var chopfail bool
			_ = chopfail

			select {
			case <-time.After(500 * time.Millisecond):
				panic(fmt.Errorf("Failed to register a handler, \"ch_reg_ws <- &reg_ws_msg{command:1, ws:conn} @!! 200\""))
			case ch_reg_ws <- &reg_ws_msg{command: 1, ws: conn}:
			}
			defer func() {
				select {
				case <-time.After(500 * time.Millisecond):
					panic(fmt.Errorf("Failed to un-register a handler, \"ch_reg_ws <- &reg_ws_msg{command:0, ws:conn} @!! 200\""))
				case ch_reg_ws <- &reg_ws_msg{command: 0, ws: conn}:
				}
				conn.Close()
			}()
			conn, err := upgrader.Upgrade(w, r, nil)
			if err != nil {
				panic(fmt.Errorf("!!!-9: WebSocket upgrade failed at \"conn, err := upgrader.Upgrade(w, r, nil)\":   %v", err))
			}
			log.Printf("WAUI controlplane: New active WebSocket connection, %s\n", conn.RemoteAddr())

			msg_dmap := &daob.DMAP{
				ChDown:    make(chan *daob.DMAP, 1),
				ReactMeta: map[string]interface{}{"WAUI": true},
			}

			msg := struct {
				Release_VERTIMESTAMP string
			}{
				"20201127-150757 {97f42b625fb4f1c86f4c5b1ddbe5298f83c01663}",
			}
			log.Printf("WAUI controlplane: sending the initial context reset: %+v\n", msg)
			err = conn.WriteJSON(msg)
			if err != nil {
				panic(fmt.Errorf("!!!-10:  at \"err = conn.WriteJSON(msg)\":   %v", err))
			}

			msg_dmap.Command = 1
			msg_dmap.DMID = 0
			msg_dmap.Key = "IO_MNEMONIC_MAP"

			err = msg_dmap.Exchange(ch_o1_S1, ctx_server)
			if err != nil {
				panic(fmt.Errorf("!!!-11:  at \"err = msg_dmap. Exchange(ch_o1_S1, ctx_server)\":   %v", err))
			}

			msg_dmap.Command = 5

			err = msg_dmap.Exchange(ch_o1_S1, ctx_server)
			if err != nil {
				panic(fmt.Errorf("!!!-12:  at \"err = msg_dmap. Exchange(ch_o1_S1, ctx_server)\":   %v", err))
			}

			var map1 map[string]int
			func() {
				defer func() {
					if errrec := recover(); errrec != nil {
						err = fmt.Errorf("try{} panicked at \"map1 = msg_dmap.V.(map[string]int)\": %v", errrec)
					}
				}()
				err = nil
				map1 = msg_dmap.V.(map[string]int)
			}()
			if err != nil {
				panic(fmt.Errorf("!!!-13: :   %v", err))
			}

			map2 := map[string]interface{}{}

			for {

				req := CPUI_req{}
				err = conn.ReadJSON(&req)
				if err != nil {
					panic(fmt.Errorf("!!!-14: Failed at \"err = conn.ReadJSON(&req)\":   %v", err))
				}

				ctx_request, cancel := context.WithTimeout(ctx_server, 2000*time.Millisecond)

				msg_dmap.Command = 255

				err := msg_dmap.Exchange(ch_o1_S1, ctx_request)
				if err != nil {
					panic(fmt.Errorf("!!!-15:  at \"err := msg_dmap. Exchange( ch_o1_S1, ctx_request)\":   %v", err))
				}

				chDown_103 := make(chan string, 1)
				go func() {
					for {
						select {
						case str := <-chDown_103:
							msg := struct {
								RequestType uint
								NoReset     uint
								Data        string
							}{
								RequestType: 999,
								NoReset:     1,
								Data:        fmt.Sprintf("CPUI Request 103 chDown diagnostic reply: %v", str),
							}
							conn.WriteJSON(msg)
							log.Print(msg.Data)
						}
					}
				}()

				switch req.RequestType {
				case 1:
					t := time.Now()
					msg := struct {
						RequestType   uint
						NoReset       uint
						Modbus_S_Time string
					}{
						NoReset:       1,
						RequestType:   req.RequestType,
						Modbus_S_Time: fmt.Sprintf("%v", t.Format("2006-01-02 15:04:05.000")),
					}
					conn.WriteJSON(msg)

				case 22:

					waui_msg := struct {
						RequestType uint
						Map         map[string]interface{}
					}{
						RequestType: 22,
						Map:         map2,
					}

					for key, dmid := range map1 {

						msg_dmap.Command = 1
						msg_dmap.DMID = dmid

						err = msg_dmap.Exchange(ch_o1_S1, ctx_request)
						if err == nil {

							err = func() (err2 error) {
								defer func() {
									if err := recover(); err != nil {

										err2 = err.(error)
									}
								}()

								var i2 interface{}
								if msg_dmap.Misc["WAUI_eness"] != nil {
									e := 0 - msg_dmap.Misc["WAUI_eness"].(int)
									switch v := msg_dmap.V.(type) {
									case uint8:
										i2 = float64(v) * math.Pow10(e)
									case uint16:
										i2 = float64(v) * math.Pow10(e)
									case uint32:
										i2 = float64(v) * math.Pow10(e)
									case uint64:
										i2 = float64(v) * math.Pow10(e)
									case int8:
										i2 = float64(v) * math.Pow10(e)
									case int16:
										i2 = float64(v) * math.Pow10(e)
									case int32:
										i2 = float64(v) * math.Pow10(e)
									case int64:
										i2 = float64(v) * math.Pow10(e)
									case float32:
										i2 = float64(v) * math.Pow10(e)
									case float64:
										i2 = float64(v) * math.Pow10(e)
									case int:
										i2 = float64(v) * math.Pow10(e)
									case uint:
										i2 = float64(v) * math.Pow10(e)
									default:
										panic(fmt.Errorf("Unsupported type in msg_dmap.V: %T", msg_dmap.V))
									}
								} else {
									i2 = msg_dmap.V
								}
								if msg_dmap.Misc["WAUI_format"] != nil {
									format := msg_dmap.Misc["WAUI_format"].(string)

									map2[key] = fmt.Sprintf(format, i2)
								} else {
									map2[key] = msg_dmap.V
								}
								return
							}()

							if err != nil {
								map2[key] = msg_dmap.V
							}

						} else {
							delete(map2, key)
						}
					}

					conn.WriteJSON(waui_msg)

				case 101:

					dmid := map1[req.Key]

					msg_dmap.Command = 1
					msg_dmap.DMID = dmid

					err = msg_dmap.Exchange(ch_o1_S1, ctx_request)
					if err != nil {
						log.Printf("WAUI controlplane failure, HTTP_WS_server.go:462: %v\n", err)
						break
					}

					v, ok := msg_dmap.V.(bool)
					if !ok {
						log.Printf("WAUI controlplane failure, HTTP_WS_server.go:465\n")
						break
					}

					if v {
						v = false
					} else {
						v = true
					}

					msg_dmap.Command = 2
					msg_dmap.DMID = dmid
					msg_dmap.V = v
					msg_dmap.Misc = nil

					err = msg_dmap.Exchange(ch_o1_S1, ctx_request)
					if err != nil {
						log.Printf("WAUI controlplane failure, HTTP_WS_server.go:476: %v\n", err)
						break
					}

				case 102:

					dmid := map1[req.Key]

					msg_dmap.Command = 1
					msg_dmap.DMID = dmid

					err = msg_dmap.Exchange(ch_o1_S1, ctx_request)
					if err != nil {
						log.Printf("WAUI controlplane failure, HTTP_WS_server.go:489: %v\n", err)
						break
					}

					err = msg_dmap.SetV_TypeAware(req.NewValue)
					if err != nil {
						msg := struct {
							RequestType uint
							NoReset     uint
							Data        string
						}{
							RequestType: 999,
							NoReset:     1,
							Data:        fmt.Sprintf("WAUI controlplane failure, HTTP_WS_server.go:500: msg_dmap.V is %T %v;; %v\n", msg_dmap.V, msg_dmap.V, err),
						}
						conn.WriteJSON(msg)

						log.Print(msg.Data)
						break
					}

					msg_dmap.Command = 2
					msg_dmap.DMID = dmid
					msg_dmap.Misc = nil

					err = msg_dmap.Exchange(ch_o1_S1, ctx_request)
					if err != nil {
						log.Printf("WAUI controlplane failure, HTTP_WS_server.go:514: %v\n", err)
						break
					}

				case 103:

					rq2 := req
					rq2.chDown = chDown_103

					select {
					case ch_CPUI_req <- &rq2:
					default:
					}

				case 50:
					msg := struct {
						ProgramName string
						RequestType uint
						NoReset     uint
						Data        interface{}
						DMID        int
					}{
						ProgramName: req.Key,
						RequestType: req.RequestType,
						NoReset:     1,
					}

					msg_dmap.Command = 1
					msg_dmap.DMID = MNEMONIC_DMID_MAP["Programs"]
					msg_dmap.Key = msg.ProgramName

					err = msg_dmap.Exchange(ch_o1_S1, ctx_request)
					if err != nil {
						panic(fmt.Errorf("!!!-18:  at \"err = msg_dmap. Exchange(ch_o1_S1, ctx_request)\":   %v", err))
					}

					msg.DMID = msg_dmap.DMID
					msg.Data = msg_dmap.V

					conn.WriteJSON(msg)

				}

				cancel()
			}
		},
	)

	m.HandleFunc("/debug/pprof/", pprof.Index)
	m.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	m.HandleFunc("/debug/pprof/profile", pprof.Profile)
	m.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	m.HandleFunc("/debug/pprof/trace", pprof.Trace)

	m.Handle("/", handlers.CombinedLoggingHandler(os.Stdout, http.FileServer(Dir(false, "/_R/resources/static/"))))

	log.Printf("HTTP_WS_server: Serving HTTP on port :8100\n")
	err = srv.ListenAndServeTLS("", "")

	panic(fmt.Errorf("HTTP_WS_server closed: %v", err))
}
