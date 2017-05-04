package taxcloud_test

import (
	"fmt"
	"github.com/ottemo/foundation/app/actors/tax/taxcloud"
	"strconv"
	"time"
)

const (
	testAPILoginID = "8278E10"
	testAPIKey     = "F5AB1E9E-8F4F-406D-A04F-946DD9DAEF10"
)

var (
	testGateway = taxcloud.NewGateway(testAPILoginID, testAPIKey, nil)
)

type testErrorProcessorType struct{}

func (t *testErrorProcessorType) Process(uniqueCode string, err error) {
	fmt.Println("TEST ERROR PROCESSING", uniqueCode, err)
}

func getUniqueStr() string {
	return strconv.FormatInt(time.Now().UnixNano(), 10)
}
