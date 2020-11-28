package main

import (
	"1axk.com/latitov/daob"
	"context"
	"fmt"
	"log"
	"os"
	"runtime"
	"sort"
	"time"
)

func S1(
	mictx context.Context,
	mictx_cancel func(),
	ch_i1 chan *daob.DMAP,
) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("S1 panicked: %v", err)
			mictx_cancel()
		}
	}()
	var chopfail bool
	_ = chopfail
	var err error

	var DM *daob.DataModelSuperrootType

	log.Println("S1: Creating the data model...")
	DM = daob.NewDM()

	log.Println("S1: Populating data model with data structures...")

	DM.R.MustAddNewAtSM("MODBUS1_RR", &daob.Node{})
	MNEMONIC_DMID_MAP["MODBUS1_RR"] = DM.R.SM["MODBUS1_RR"].DMID

	log.Println("\t\"IO_MNEMONIC_MAP\"")

	DM.R.MustAddNewAtSM("IO_MNEMONIC_MAP", &daob.Node{})

	sa6 := [6]uint16{1, 2, 3}

	DM.R.SM["IO_MNEMONIC_MAP"].MustAddNewAtSM(".io.1.d1", &daob.Node{
		V: sa6, Misc: map[string]interface{}{"Type": ""},
	})
	MNEMONIC_DMID_MAP[".io.1.d1"] = DM.R.SM["IO_MNEMONIC_MAP"].SM[".io.1.d1"].DMID

	DM.R.SM["IO_MNEMONIC_MAP"].MustAddNewAtSM(".io.1.d2", &daob.Node{
		V: sa6, Misc: map[string]interface{}{"Type": ""},
	})
	MNEMONIC_DMID_MAP[".io.1.d2"] = DM.R.SM["IO_MNEMONIC_MAP"].SM[".io.1.d2"].DMID

	DM.R.SM["IO_MNEMONIC_MAP"].MustAddNewAtSM(".io.2.d1", &daob.Node{
		V: [2]uint16{}, Misc: map[string]interface{}{"Type": ""},
	})

	DM.R.SM["IO_MNEMONIC_MAP"].MustAddNewAtSM(".io.2.d1.i1", &daob.Node{
		V: false, Misc: map[string]interface{}{"Type": "bool"},
	})
	MNEMONIC_DMID_MAP[".io.2.d1.i1"] = DM.R.SM["IO_MNEMONIC_MAP"].SM[".io.2.d1.i1"].DMID

	DM.R.SM["IO_MNEMONIC_MAP"].MustAddNewAtSM(".io.2.d1.i2", &daob.Node{
		V: false, Misc: map[string]interface{}{"Type": "bool"},
	})
	MNEMONIC_DMID_MAP[".io.2.d1.i2"] = DM.R.SM["IO_MNEMONIC_MAP"].SM[".io.2.d1.i2"].DMID

	DM.R.SM["IO_MNEMONIC_MAP"].MustAddNewAtSM(".io.2.d1.i3", &daob.Node{
		V: false, Misc: map[string]interface{}{"Type": "bool"},
	})
	MNEMONIC_DMID_MAP[".io.2.d1.i3"] = DM.R.SM["IO_MNEMONIC_MAP"].SM[".io.2.d1.i3"].DMID

	DM.R.SM["IO_MNEMONIC_MAP"].MustAddNewAtSM(".io.2.d1.i4", &daob.Node{
		V: false, Misc: map[string]interface{}{"Type": "bool"},
	})
	MNEMONIC_DMID_MAP[".io.2.d1.i4"] = DM.R.SM["IO_MNEMONIC_MAP"].SM[".io.2.d1.i4"].DMID

	DM.R.SM["IO_MNEMONIC_MAP"].MustAddNewAtSM(".io.2.d1.i5", &daob.Node{
		V: false, Misc: map[string]interface{}{"Type": "bool"},
	})
	MNEMONIC_DMID_MAP[".io.2.d1.i5"] = DM.R.SM["IO_MNEMONIC_MAP"].SM[".io.2.d1.i5"].DMID

	DM.R.SM["IO_MNEMONIC_MAP"].MustAddNewAtSM(".io.2.d1.i6", &daob.Node{
		V: false, Misc: map[string]interface{}{"Type": "bool"},
	})
	MNEMONIC_DMID_MAP[".io.2.d1.i6"] = DM.R.SM["IO_MNEMONIC_MAP"].SM[".io.2.d1.i6"].DMID

	DM.R.SM["IO_MNEMONIC_MAP"].MustAddNewAtSM(".io.2.d1.i7", &daob.Node{
		V: false, Misc: map[string]interface{}{"Type": "bool"},
	})
	MNEMONIC_DMID_MAP[".io.2.d1.i7"] = DM.R.SM["IO_MNEMONIC_MAP"].SM[".io.2.d1.i7"].DMID

	DM.R.SM["IO_MNEMONIC_MAP"].MustAddNewAtSM(".io.2.d1.i8", &daob.Node{
		V: false, Misc: map[string]interface{}{"Type": "bool"},
	})
	MNEMONIC_DMID_MAP[".io.2.d1.i8"] = DM.R.SM["IO_MNEMONIC_MAP"].SM[".io.2.d1.i8"].DMID

	DM.R.SM["IO_MNEMONIC_MAP"].MustAddNewAtSM(".io.2.d1.q1", &daob.Node{
		V: false, Misc: map[string]interface{}{"Type": "bool", "R": 1, "Bit": 0},
	})
	MNEMONIC_DMID_MAP[".io.2.d1.q1"] = DM.R.SM["IO_MNEMONIC_MAP"].SM[".io.2.d1.q1"].DMID

	DM.R.SM["IO_MNEMONIC_MAP"].MustAddNewAtSM(".io.2.d1.q2", &daob.Node{
		V: false, Misc: map[string]interface{}{"Type": "bool", "R": 1, "Bit": 1},
	})
	MNEMONIC_DMID_MAP[".io.2.d1.q2"] = DM.R.SM["IO_MNEMONIC_MAP"].SM[".io.2.d1.q2"].DMID

	DM.R.SM["IO_MNEMONIC_MAP"].MustAddNewAtSM(".io.2.d1.q3", &daob.Node{
		V: false, Misc: map[string]interface{}{"Type": "bool", "R": 1, "Bit": 2},
	})
	MNEMONIC_DMID_MAP[".io.2.d1.q3"] = DM.R.SM["IO_MNEMONIC_MAP"].SM[".io.2.d1.q3"].DMID

	DM.R.SM["IO_MNEMONIC_MAP"].MustAddNewAtSM(".io.2.d1.q4", &daob.Node{
		V: false, Misc: map[string]interface{}{"Type": "bool", "R": 1, "Bit": 3},
	})
	MNEMONIC_DMID_MAP[".io.2.d1.q4"] = DM.R.SM["IO_MNEMONIC_MAP"].SM[".io.2.d1.q4"].DMID

	DM.R.SM["IO_MNEMONIC_MAP"].MustAddNewAtSM(".io.2.d1.q5", &daob.Node{
		V: false, Misc: map[string]interface{}{"Type": "bool", "R": 1, "Bit": 4},
	})
	MNEMONIC_DMID_MAP[".io.2.d1.q5"] = DM.R.SM["IO_MNEMONIC_MAP"].SM[".io.2.d1.q5"].DMID

	DM.R.SM["IO_MNEMONIC_MAP"].MustAddNewAtSM(".io.2.d1.q6", &daob.Node{
		V: false, Misc: map[string]interface{}{"Type": "bool", "R": 1, "Bit": 5},
	})
	MNEMONIC_DMID_MAP[".io.2.d1.q6"] = DM.R.SM["IO_MNEMONIC_MAP"].SM[".io.2.d1.q6"].DMID

	DM.R.SM["IO_MNEMONIC_MAP"].MustAddNewAtSM(".io.2.d1.q7", &daob.Node{
		V: false, Misc: map[string]interface{}{"Type": "bool", "R": 1, "Bit": 6},
	})
	MNEMONIC_DMID_MAP[".io.2.d1.q7"] = DM.R.SM["IO_MNEMONIC_MAP"].SM[".io.2.d1.q7"].DMID

	DM.R.SM["IO_MNEMONIC_MAP"].MustAddNewAtSM(".io.2.d1.q8", &daob.Node{
		V: false, Misc: map[string]interface{}{"Type": "bool", "R": 1, "Bit": 7},
	})
	MNEMONIC_DMID_MAP[".io.2.d1.q8"] = DM.R.SM["IO_MNEMONIC_MAP"].SM[".io.2.d1.q8"].DMID

	DM.R.SM["IO_MNEMONIC_MAP"].MustAddNewAtSM(".io.2.d2", &daob.Node{
		V: [8]uint16{}, Misc: map[string]interface{}{"Type": ""},
	})
	MNEMONIC_DMID_MAP[".io.2.d2"] = DM.R.SM["IO_MNEMONIC_MAP"].SM[".io.2.d2"].DMID

	DM.R.SM["IO_MNEMONIC_MAP"].MustAddNewAtSM(".io.2.d12", &daob.Node{
		SA: make([]*daob.Node, 50),
	})
	MNEMONIC_DMID_MAP[".io.2.d12"] = DM.R.SM["IO_MNEMONIC_MAP"].SM[".io.2.d12"].DMID

	DM.R.SM["IO_MNEMONIC_MAP"].MustAddNewAtSM("runtime.is727.PrintAllStacksSignal", &daob.Node{
		V: false, Misc: map[string]interface{}{"Type": "bool"},
	})
	AUTO_ID_3 := false
	AUTO_ID_2 := &runtime.MemStats{}
	DM.R.SM["IO_MNEMONIC_MAP"].MustAddNewAtSM("runtime.MemStats.Alloc", &daob.Node{})
	DM.R.SM["IO_MNEMONIC_MAP"].MustAddNewAtSM("runtime.MemStats.NumGC", &daob.Node{})
	DM.R.SM["IO_MNEMONIC_MAP"].MustAddNewAtSM("runtime.MemStats.HeapObjects", &daob.Node{})
	DM.R.SM["IO_MNEMONIC_MAP"].MustAddNewAtSM("runtime.NumGoroutine", &daob.Node{})
	DM.R.SM["IO_MNEMONIC_MAP"].MustAddNewAtSM("S1_Heartbeat", &daob.Node{})

	log.Println("\tS1: Set Misc[Name] for all in the map...")
	for key, n1 := range DM.R.SM["IO_MNEMONIC_MAP"].SM {
		n1.Misc["Name"] = key
	}

	log.Println("\tS1: Link in the Ma to the MODBUS RR...")
	for _, n1 := range DM.R.SM["IO_MNEMONIC_MAP"].SM {
		r, ok := n1.Misc["MODBUS1_R"].(int)
		if ok && r != 0 {
			DM.R.SM["MODBUS1_RR"].
				MustLinkIn_SA(r, n1)
		}
	}

	DM.R.MustAddNewAtSM("Programs", &daob.Node{})
	DM.R.SM["Programs"].
		MustAddNewAtSM("PLC_PRG", &daob.Node{}).
		MustAddNewAtSM("Post_1", &daob.Node{}).
		MustAddNewAtSM("Post_2", &daob.Node{}).
		MustAddNewAtSM("Post_3", &daob.Node{}).
		MustAddNewAtSM("Post_4", &daob.Node{}).
		MustAddNewAtSM("Post_5", &daob.Node{})

	MNEMONIC_DMID_MAP["Programs"] = DM.R.SM["Programs"].DMID
	MNEMONIC_DMID_MAP["Programs__PLC_PRG"] = DM.R.SM["Programs"].SM["PLC_PRG"].DMID
	MNEMONIC_DMID_MAP["Programs__Post_1"] = DM.R.SM["Programs"].SM["Post_1"].DMID
	MNEMONIC_DMID_MAP["Programs__Post_2"] = DM.R.SM["Programs"].SM["Post_2"].DMID
	MNEMONIC_DMID_MAP["Programs__Post_3"] = DM.R.SM["Programs"].SM["Post_3"].DMID
	MNEMONIC_DMID_MAP["Programs__Post_4"] = DM.R.SM["Programs"].SM["Post_4"].DMID
	MNEMONIC_DMID_MAP["Programs__Post_5"] = DM.R.SM["Programs"].SM["Post_5"].DMID

	log.Println("\tS1: Set Misc[Name] for all in Ma...")
	for key, n1 := range DM.R.SM["Programs"].SM {
		n1.Misc["Name"] = key
	}

	log.Println("\t\"MODBUS1_RR\"")

	if TestLevel == 2 {
		log.Println("Test Level 2: MODBUS1_RR:")
		for i, n1 := range DM.R.SM["MODBUS1_RR"].SA {
			if !n1.Null {
				name, _ := n1.Misc["Name"].(string)
				dmid := n1.DMID
				log.Printf("\tR %v, %+v, %v\n", i, name, dmid)
			}
		}
		os.Exit(0)
	}

	log.Println("\tList of keys in Ma...")

	if TestLevel == 3 {
		all := []string{}

		log.Println("Test Level 3")
		log.Println("COMPLETE LIST OF NAMES in Ma:")
		for _, n1 := range DM.R.SM["IO_MNEMONIC_MAP"].SM {
			name, _ := n1.Misc["Name"].(string)
			all = append(all, name)
		}
		sort.Strings(all)

		log.Println("COMPLETE LIST OF NAMES, SORTED:")
		for _, el := range all {
			log.Printf("\t.Name = %v\n", el)
		}
		os.Exit(0)
	}
	if TestLevel == 4 {
		all := []string{}
		map1 := map[string]int{}

		for key, n1 := range DM.R.SM["IO_MNEMONIC_MAP"].SM {
			all = append(all, key)
			map1[key] = n1.DMID
		}
		sort.Strings(all)

		log.Println("Test Level 4")
		log.Println("COMPLETE LIST OF KEYS, SORTED:")
		for _, key := range all {
			log.Printf("\tMa[%v], DMID=%v\n", key, map1[key])
		}
		os.Exit(0)
	}

	if TestLevel == 5 {

		all := []string{}

		for key := range DM.R.SM["IO_MNEMONIC_MAP"].SM {
			all = append(all, key)
		}
		sort.Strings(all)

		log.Println("Test Level 5")
		log.Println("COMPLETE LIST OF NAMES/KEYS FOR WA UI, SORTED:\n\n")
		for _, el := range all {
			fmt.Printf("KEY>\t\"%v\",\n", el)
		}
		os.Exit(0)
	}

	ticker_1s_PA := time.NewTicker(time.Second)
	defer ticker_1s_PA.Stop()
	ticker_1s := time.NewTicker(time.Second)
	defer ticker_1s.Stop()

	log.Println("S1: Prepartions done; sent DATA_MODEL_READY signal; READY TO RUN")
	close(ch_DATA_MODEL_READY)

	for {
		select {
		case msg1 := <-ch_i1:

			err = msg1.Proc(DM, 50*time.Millisecond)
			if err != nil {
				log.Printf("S1: DMAP processor: %v\n", err)
			}

		case <-ticker_1s.C:

			AUTO_ID_1 := &DM.R.SM["IO_MNEMONIC_MAP"].SM["S1_Heartbeat"].V
			switch vt := (*AUTO_ID_1).(type) {
			case int:
				*AUTO_ID_1 = vt + 1
			default:
				*AUTO_ID_1 = 0
			}

			DM.R.SM["IO_MNEMONIC_MAP"].SM["runtime.NumGoroutine"].V = runtime.NumGoroutine()

			runtime.ReadMemStats(AUTO_ID_2)

			DM.R.SM["IO_MNEMONIC_MAP"].SM["runtime.MemStats.Alloc"].V = AUTO_ID_2.Alloc
			DM.R.SM["IO_MNEMONIC_MAP"].SM["runtime.MemStats.NumGC"].V = AUTO_ID_2.NumGC
			DM.R.SM["IO_MNEMONIC_MAP"].SM["runtime.MemStats.HeapObjects"].V = AUTO_ID_2.HeapObjects

			AUTO_ID_4 := DM.R.SM["IO_MNEMONIC_MAP"].SM["runtime.is727.PrintAllStacksSignal"].V.(bool)
			if AUTO_ID_4 == true && AUTO_ID_3 == false {
				buf := make([]byte, 1<<22)
				len := runtime.Stack(buf, true)
				log.Printf("==================================================\nS1: REQUEST TO PRINT THE STACK TRACE OF ALL GOROUTINES:\n\n%v\n\n", string(buf[:len]))
			}
			DM.R.SM["IO_MNEMONIC_MAP"].SM["runtime.is727.PrintAllStacksSignal"].V = AUTO_ID_4

		case <-ticker_1s_PA.C:

		case <-mictx.Done():
			return
		}
	}
}
