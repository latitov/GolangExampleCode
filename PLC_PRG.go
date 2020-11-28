package main

import (
	"1axk.com/latitov/daob"
	"1axk.com/latitov/hjson"
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"time"
)

type PLC_PRG struct {
	Heartbeat int
	re_json   *regexp.Regexp
	re_hjson  *regexp.Regexp
}

func (me *PLC_PRG) Run(
	mictx context.Context,
	mictx_cancel func(),
	ch_o1_S1 chan *daob.DMAP,
	ch_CPUI_req chan *CPUI_req,
) {
	var err error
	_ = err

	me.Heartbeat++

	tick1 := time.NewTicker(100 * time.Millisecond)

	ctx, ctx_cancel := context.WithTimeout(mictx, 200*time.Millisecond)

	msg_dmap := &daob.DMAP{
		ChDown:    make(chan *daob.DMAP, 1),
		ReactMeta: map[string]interface{}{"PLC_PRG": true},
	}

	Post_1 := &Post_PRG{}

	func() {
		defer func() {
			if errrec := recover(); errrec != nil {
				err = fmt.Errorf("v=.io. panicked at \"Programs__Post_1 = Post_1\": %v", errrec)
			}
		}()
		err = nil
		msg_dmap.Command = 2
		msg_dmap.DMID = MNEMONIC_DMID_MAP["Programs__Post_1"]
		msg_dmap.V = Post_1
		msg_dmap.Misc = nil
		err = msg_dmap.Exchange(ch_o1_S1, ctx)
		if err != nil {
			panic(err)
		}
	}()

	Post_2 := &Post_PRG{}

	func() {
		defer func() {
			if errrec := recover(); errrec != nil {
				err = fmt.Errorf("v=.io. panicked at \"Programs__Post_2 = Post_2\": %v", errrec)
			}
		}()
		err = nil
		msg_dmap.Command = 2
		msg_dmap.DMID = MNEMONIC_DMID_MAP["Programs__Post_2"]
		msg_dmap.V = Post_2
		msg_dmap.Misc = nil
		err = msg_dmap.Exchange(ch_o1_S1, ctx)
		if err != nil {
			panic(err)
		}
	}()

	Post_3 := &Post_PRG{}

	func() {
		defer func() {
			if errrec := recover(); errrec != nil {
				err = fmt.Errorf("v=.io. panicked at \"Programs__Post_3 = Post_3\": %v", errrec)
			}
		}()
		err = nil
		msg_dmap.Command = 2
		msg_dmap.DMID = MNEMONIC_DMID_MAP["Programs__Post_3"]
		msg_dmap.V = Post_3
		msg_dmap.Misc = nil
		err = msg_dmap.Exchange(ch_o1_S1, ctx)
		if err != nil {
			panic(err)
		}
	}()

	Post_4 := &Post_PRG{}

	func() {
		defer func() {
			if errrec := recover(); errrec != nil {
				err = fmt.Errorf("v=.io. panicked at \"Programs__Post_4 = Post_4\": %v", errrec)
			}
		}()
		err = nil
		msg_dmap.Command = 2
		msg_dmap.DMID = MNEMONIC_DMID_MAP["Programs__Post_4"]
		msg_dmap.V = Post_4
		msg_dmap.Misc = nil
		err = msg_dmap.Exchange(ch_o1_S1, ctx)
		if err != nil {
			panic(err)
		}
	}()

	Post_5 := &Post_PRG{}

	func() {
		defer func() {
			if errrec := recover(); errrec != nil {
				err = fmt.Errorf("v=.io. panicked at \"Programs__Post_5 = Post_5\": %v", errrec)
			}
		}()
		err = nil
		msg_dmap.Command = 2
		msg_dmap.DMID = MNEMONIC_DMID_MAP["Programs__Post_5"]
		msg_dmap.V = Post_5
		msg_dmap.Misc = nil
		err = msg_dmap.Exchange(ch_o1_S1, ctx)
		if err != nil {
			panic(err)
		}
	}()

	for {
		select {
		case <-tick1.C:

			Post_1.Run(
				ctx, ctx_cancel,
				ch_o1_S1,
			)

			Post_2.Run(
				ctx, ctx_cancel,
				ch_o1_S1,
			)

			Post_3.Run(
				ctx, ctx_cancel,
				ch_o1_S1,
			)

			Post_4.Run(
				ctx, ctx_cancel,
				ch_o1_S1,
			)

			Post_5.Run(
				ctx, ctx_cancel,
				ch_o1_S1,
			)

		case cpuireq := <-ch_CPUI_req:

			var prg interface{}
			switch cpuireq.Key {
			case "PLC_PRG":
				prg = me
			case "Post_1":
				prg = Post_1
			case "Post_2":
				prg = Post_2
			case "Post_3":
				prg = Post_3
			case "Post_4":
				prg = Post_4
			case "Post_5":
				prg = Post_5
			default:
				errmsg := fmt.Sprintf("Unknown program '%v'", cpuireq.Key)
				select {
				case cpuireq.chDown <- errmsg:
				default:
				}
			}

			if me.re_json == nil || me.re_hjson == nil {
				me.re_json = regexp.MustCompile("^(?mi)json:")
				me.re_hjson = regexp.MustCompile("^(?mi)hjson:")
			}
			if me.re_json.MatchString(cpuireq.NewValue) {
				cpuireq.NewValue = me.re_json.ReplaceAllString(cpuireq.NewValue, "")
				err = json.Unmarshal([]byte(cpuireq.NewValue), prg)
				if err != nil {
					errmsg := fmt.Sprintf("json.Unmarshal(): %v;; the string was: '%v'", err, cpuireq.NewValue)
					select {
					case cpuireq.chDown <- errmsg:
					default:
					}
				} else {
					errmsg := fmt.Sprintf("<span style=\"font-size: 80%%; color: #080;\">json.Unmarshal(): PERFORMED OK, program='%v' (%T);; the string was: '%v'</span>", cpuireq.Key, prg, cpuireq.NewValue)
					select {
					case cpuireq.chDown <- errmsg:
					default:
					}
				}
			} else if me.re_hjson.MatchString(cpuireq.NewValue) {
				cpuireq.NewValue = me.re_hjson.ReplaceAllString(cpuireq.NewValue, "")
				err = hjson.Unmarshal([]byte(cpuireq.NewValue), prg)
				if err != nil {
					errmsg := fmt.Sprintf("hjson.Unmarshal(): %v;; the string was: '%v'", err, cpuireq.NewValue)
					select {
					case cpuireq.chDown <- errmsg:
					default:
					}
				} else {
					errmsg := fmt.Sprintf("<span style=\"font-size: 80%%; color: #080;\">json.Unmarshal(): PERFORMED OK, program='%v' (%T);; the string was: '%v'</span>", cpuireq.Key, prg, cpuireq.NewValue)
					select {
					case cpuireq.chDown <- errmsg:
					default:
					}
				}
			} else {
				errmsg := fmt.Sprintf("Syntax unknown; use 'json:' or 'hjson:' prefix;; the string was: '%v'", cpuireq.NewValue)
				select {
				case cpuireq.chDown <- errmsg:
				default:
				}
			}
		}
	}
}
