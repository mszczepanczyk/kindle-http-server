package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strconv"
	"strings"

	"github.com/coreos/go-iptables/iptables"
)

func AddIPTablesRule(table string, chain string, rulespec ...string) error {
	ipt, err := iptables.New()
	if err != nil {
		return err
	}

	log.Printf("Adding iptables rule %s %s %v", table, chain, rulespec)
	err = ipt.AppendUnique(table, chain, rulespec...)
	if err != nil {
		return err
	}

	return nil
}

func ExecCommands(commands [][]string) (string, error) {
	stdout := ""
	for _, command := range commands {
		cmdStdout, err := exec.Command(command[0], command[1:]...).Output()
		stdout += string(cmdStdout)
		if err != nil {
			return stdout, err
		}
	}
	return stdout, nil
}

func Swipe(x1 int, y1 int, x2 int, y2 int) error {
	_, err := ExecCommands([][]string{
		{"evemu-event", "/dev/input/event1", "--type", "EV_KEY", "--code", "BTN_TOUCH", "--value", "1"},
		{"evemu-event", "/dev/input/event1", "--type", "EV_ABS", "--code", "ABS_MT_TRACKING_ID", "--value", "0"},
		{"evemu-event", "/dev/input/event1", "--type", "EV_ABS", "--code", "ABS_MT_POSITION_X", "--value", strconv.Itoa(x1)},
		{"evemu-event", "/dev/input/event1", "--type", "EV_ABS", "--code", "ABS_MT_POSITION_Y", "--value", strconv.Itoa(y1)},
		{"evemu-event", "/dev/input/event1", "--type", "EV_SYN", "--code", "SYN_REPORT", "--value", "0"},
		{"evemu-event", "/dev/input/event1", "--type", "EV_ABS", "--code", "ABS_MT_POSITION_X", "--value", strconv.Itoa(x2)},
		{"evemu-event", "/dev/input/event1", "--type", "EV_ABS", "--code", "ABS_MT_POSITION_Y", "--value", strconv.Itoa(y2)},
		{"evemu-event", "/dev/input/event1", "--type", "EV_SYN", "--code", "SYN_REPORT", "--value", "0"},
		{"evemu-event", "/dev/input/event1", "--type", "EV_ABS", "--code", "ABS_MT_TRACKING_ID", "--value", "-1"},
		{"evemu-event", "/dev/input/event1", "--type", "EV_SYN", "--code", "SYN_REPORT", "--value", "0"},
		{"evemu-event", "/dev/input/event1", "--type", "EV_KEY", "--code", "BTN_TOUCH", "--value", "0"},
		{"evemu-event", "/dev/input/event1", "--type", "EV_SYN", "--code", "SYN_REPORT", "--value", "0"},
	})
	return err
}

func main() {
	const port = 80
	address := fmt.Sprintf("0.0.0.0:%d", port)

	err := AddIPTablesRule("filter", "INPUT", "-i", "wlan0", "-p", "tcp", "-m", "tcp", "--dport", strconv.Itoa(port), "-j", "ACCEPT")
	if err != nil {
		log.Println("Error adding iptables rule:", err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		json.NewEncoder(w).Encode(2137)
	})

	http.HandleFunc("/battery-level", func(w http.ResponseWriter, req *http.Request) {
		stdout, err := exec.Command("lipc-get-prop", "-i", "com.lab126.powerd", "battLevel").Output()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(-1)
			return
		}

		battery_level, err := strconv.Atoi(strings.TrimSpace(string(stdout)))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(-1)
		} else {
			json.NewEncoder(w).Encode(battery_level)
		}
	})

	http.HandleFunc("/toggle-power-button", func(w http.ResponseWriter, req *http.Request) {
		_, err := exec.Command("lipc-set-prop", "-i", "com.lab126.powerd", "powerButton", "1").Output()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	})

	http.HandleFunc("/lipc-get-prop/{publisher}/{property}", func(w http.ResponseWriter, req *http.Request) {
		out, err := exec.Command("lipc-get-prop", req.PathValue("publisher"), req.PathValue("property")).CombinedOutput()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(strings.TrimSpace(string(out)))
		} else {
			json.NewEncoder(w).Encode(strings.TrimSpace(string(out)))
		}
	})

	http.HandleFunc("/swipe-left", func(w http.ResponseWriter, req *http.Request) {
		err := Swipe(300, 700, 500, 700)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}
		json.NewEncoder(w).Encode("OK")
	})

	http.HandleFunc("/swipe-right", func(w http.ResponseWriter, req *http.Request) {
		err := Swipe(500, 700, 300, 700)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}
		json.NewEncoder(w).Encode("OK")
	})

	log.Println("Starting http server on", address)

	err = http.ListenAndServe(address, nil)
	if err != nil {
		log.Fatal(err)
	}
}
