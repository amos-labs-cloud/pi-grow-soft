package measurement

import (
	"fmt"
	"github.com/influxdata/influxdb-client-go/v2"
	"testing"
	"time"
)

func TestMetrics(t *testing.T) {
	// You can generate a Token from the "Tokens Tab" in the UI
	const token = "F-QFQpmCL9UkR3qyoXnLkzWj03s6m4eCvYgDl1ePfHBf9ph7yxaSgQ6WN0i9giNgRTfONwVMK1f977r_g71oNQ=="
	const bucket = "pgs"
	const org = "amoslabs"

	client := influxdb2.NewClient("http://localhost:8086", token)
	// always close client at the end
	defer client.Close()

	// get non-blocking write client
	writeAPI := client.WriteAPI(org, bucket)

	// write line protocol
	for _, v := range []float32{1, 2, 3, 4, 5, 6, 7, 8} {
		writeAPI.WriteRecord(fmt.Sprintf("stat,unit=temperature avg=%f,max=%f", 23.0+v, 45.0+v))
		writeAPI.Flush()
		time.Sleep(time.Second * 10)
	}
	// Flush writes

}
