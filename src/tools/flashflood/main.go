package main
import (
//	"net/http"
//	"tools/flashflood/server"
)

func main() {
//	port := os.Getenv("PORT")
//	loggregatorURL := os.Getenv("LOGGREGATOR_URL")
//	uaaURL := os.Getenv("UAA_URL")
//	username := os.Getenv("USERNAME")
//	password := os.Getenv("PASSWORD")
//
//	vcapJSON := os.Getenv("VCAP_APPLICATION")
//	vcapData := make(map[string]interface{})
//	err := json.Unmarshal([]byte(vcapJSON), &vcapData)
//	if err != nil {
//		panic(err)
//	}
//
//	appID := vcapData["application_id"].(string)

//	http.ListenAndServe("localhost:1234", server.New())

	//for
		//noaa.Start
		//emitter.sendOne()
		//wait until receive the sendOne() msg

		//select
			//case: "/stop"

	// emitter.flood() //go

	// sent, received := startTest(100 * 2**n)
	// when we visit /start ...
		// emit logs at 100/s for 10s
		// count loss
		// emit logs at 200/s for 10s
		// count loss
		// emit at 400/s ...
		// 800/s ...

}
