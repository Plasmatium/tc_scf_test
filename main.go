package main

import (
	"fmt"
	"os"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	scf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/scf/v20180416"
)

var reqLimit = 20
var invoked = 0
var ch = make(chan struct{}, reqLimit)

func do(cpf *profile.ClientProfile, cred *common.Credential, params string) {
	status := "finished"
	defer func() {
		ch <- struct{}{}
		invoked++
		fmt.Println()
		fmt.Println(invoked, status)
		fmt.Println()
	}()

	c, _ := scf.NewClient(cred, "ap-shanghai", cpf)
	req := scf.NewInvokeRequest()

	err := req.FromJsonString(params)
	if err != nil {
		status = "failed"
		panic(err)
	}
	resp, err := c.Invoke(req)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		fmt.Printf("An API error has returned: %s\n", err)
		ch <- struct{}{}
	}
	if err != nil {
		status = "failed"
		panic(err)
	}
	fmt.Printf("%s", resp.ToJsonString())
}

func main() {

	credential := common.NewCredential(
		os.Getenv("TENCENTCLOUD_SECRET_ID"),
		os.Getenv("TENCENTCLOUD_SECRET_KEY"),
	)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "scf.tencentcloudapi.com"

	params := `{"FunctionName":"fromAWS"}`
	// err := request.FromJsonString(params)
	// if err != nil {
	// 	panic(err)
	// }
	// response, err := client.Invoke(request)
	// if _, ok := err.(*errors.TencentCloudSDKError); ok {
	// 	fmt.Printf("An API error has returned: %s", err)
	// 	return
	// }
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Printf("%s", response.ToJsonString())
	for i := 0; i < reqLimit; i++ {
		go do(cpf, credential, params)
	}
	for i := 0; i < reqLimit; i++ {
		<-ch
	}
	fmt.Println("all done")
}
