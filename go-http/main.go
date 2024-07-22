package main

import (
	"fmt"
	"log"
	"net/http"
)

const Logo = `
 '/¯¯¯/|¯¯¯|'         |\¯¯¯¯\'                   |\¯¯¯\ \¯¯¯\    |¯¯¯¯¯¯¯¯¯¯|    |¯¯¯¯¯¯¯¯¯¯|    |¯¯¯¯|\¯¯¯¯\  
 |    °| |___|'  /¯¯¯¯/\       \'                 | |     |_|     |'    |___       ___|    |___       ___|    |       | |       | °
 |    °| |    °|' '|       '| '|       |                '°\|            °|'    |     |      |     |    |     |      |     |    |       |/____/|  
 |    °|'\¯¯¯\' '|\       \/____/|'                °'|     |¯|     |'    |___|    °'|___|   |___|    °'|___|   |       |      °| |  
 |\___\|___|  '| '\____\      | |                 °|\___\¸\___\'         |___'|               |___'|         |____|____'|/'  
 | |     ||     |   \ '|       |___'|/                 °| |     | '|     '|         |      |'               |      |'         |       |           
'°\|__¸'||__¸'|    '\|____|     °                     \|__¸'| '|__¸¸|         |__¸¸'|'               |__¸¸'|'         |____|           

`

func main() {

	fmt.Print(Logo)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello world")
	})

	log.Fatal(http.ListenAndServe(":8080", nil))

}
