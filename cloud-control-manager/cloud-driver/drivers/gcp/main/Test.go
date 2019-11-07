// Proof of Concepts of CB-Spider.
// The CB-Spider is a sub-Framework of the Cloud-Barista Multi-Cloud Project.
// The CB-Spider Mission is to connect all the clouds with a single interface.
//
//      * Cloud-Barista: https://github.com/cloud-barista
//
// This is a Cloud Driver Example for PoC Test.
//
// by hyokyung.kim@innogrid.co.kr, 2019.07.

package main

import (
	"fmt"

	testconf "./conf"
	cblog "github.com/cloud-barista/cb-log"
	irs "github.com/cloud-barista/cb-spider/cloud-control-manager/cloud-driver/interfaces/resources"
	"github.com/davecgh/go-spew/spew"
	"github.com/sirupsen/logrus"
)

var cblogger *logrus.Logger

func init() {
	// cblog is a global variable.
	cblogger = cblog.GetLogger("GCP Resource Test")
	cblog.SetLevel("debug")
}

// Test VM
func handleVM() {
	cblogger.Info("Start VM Resource Test")

	ResourceHandler, err := testconf.GetResourceHandler("VM")
	if err != nil {
		panic(err)
	}

	handler := ResourceHandler.(irs.VMHandler)

	VmID := "vm01"

	for {
		fmt.Println("VM Management")
		fmt.Println("0. Quit")
		fmt.Println("1. VM Start")
		fmt.Println("2. VM Info")
		fmt.Println("3. Suspend VM")
		fmt.Println("4. Resume VM")
		fmt.Println("5. Reboot VM")
		fmt.Println("6. Terminate VM")

		fmt.Println("7. GetVMStatus VM")
		fmt.Println("8. ListVMStatus VM")
		fmt.Println("9. ListVM")

		var commandNum int
		inputCnt, err := fmt.Scan(&commandNum)
		if err != nil {
			panic(err)
		}

		if inputCnt == 1 {
			switch commandNum {
			case 0:
				return

			case 1:
				vmReqInfo := irs.VMReqInfo{
					VMName: "vm01",
					//ImageId:            config.Aws.ImageID,
					VirtualNetworkId: "cb-vnet",
					//NetworkInterfaceId: "eni-00befb6d8c3a87b24",
					PublicIPId: "publicip-vm01",
					//SecurityGroupIds: []string{"sg-0df1c209ea1915e4b"},
					//SecurityGroupIds: []string{config.Aws.SecurityGroupID},
					VMSpecId: "f1-micro",
					//KeyPairName:      config.Aws.KeyName,
				}

				vmInfo, err := handler.StartVM(vmReqInfo)
				if err != nil {
					//panic(err)
					cblogger.Error(err)
				} else {
					cblogger.Info("VM 생성 완료!!", vmInfo)
					spew.Dump(vmInfo)
				}

			case 2:
				vmInfo, err := handler.GetVM(VmID)
				if err != nil {
					cblogger.Errorf("[%s] VM 정보 조회 실패", VmID)
					cblogger.Error(err)
				} else {
					cblogger.Infof("[%s] VM 정보 조회 결과", VmID)
					cblogger.Info(vmInfo)
					spew.Dump(vmInfo)
				}

			case 3:
				cblogger.Info("Start Suspend VM ...")
				err := handler.SuspendVM(VmID)
				if err != nil {
					cblogger.Errorf("[%s] VM Suspend 실패", VmID)
					cblogger.Error(err)
				} else {
					cblogger.Infof("[%s] VM Suspend 성공", VmID)
				}

			case 4:
				cblogger.Info("Start Resume  VM ...")
				err := handler.ResumeVM(VmID)
				if err != nil {
					cblogger.Errorf("[%s] VM Resume 실패", VmID)
					cblogger.Error(err)
				} else {
					cblogger.Infof("[%s] VM Resume 성공", VmID)
				}

			case 5:
				cblogger.Info("Start Reboot  VM ...")
				err := handler.RebootVM(VmID)
				if err != nil {
					cblogger.Errorf("[%s] VM Reboot 실패", VmID)
					cblogger.Error(err)
				} else {
					cblogger.Infof("[%s] VM Reboot 성공", VmID)
				}

			case 6:
				cblogger.Info("Start Terminate  VM ...")
				err := handler.TerminateVM(VmID)
				if err != nil {
					cblogger.Errorf("[%s] VM Terminate 실패", VmID)
					cblogger.Error(err)
				} else {
					cblogger.Infof("[%s] VM Terminate 성공", VmID)
				}

			case 7:
				cblogger.Info("Start Get VM Status...")
				vmStatus, err := handler.GetVMStatus(VmID)
				if err != nil {
					cblogger.Errorf("[%s] VM Get Status 실패", VmID)
					cblogger.Error(err)
				} else {
					cblogger.Infof("[%s] VM Get Status 성공 : [%s]", VmID, vmStatus)
				}

			case 8:
				cblogger.Info("Start ListVMStatus ...")
				vmStatusInfos, err := handler.ListVMStatus()
				if err != nil {
					cblogger.Error("ListVMStatus 실패")
					cblogger.Error(err)
				} else {
					cblogger.Info("ListVMStatus 성공")
					cblogger.Info(vmStatusInfos)
					spew.Dump(vmStatusInfos)
				}

			case 9:
				cblogger.Info("Start ListVM ...")
				vmList, err := handler.ListVM()
				if err != nil {
					cblogger.Error("ListVM 실패")
					cblogger.Error(err)
				} else {
					cblogger.Info("ListVM 성공")
					cblogger.Info("=========== VM 목록 ================")
					cblogger.Info(vmList)
					spew.Dump(vmList)
				}

			}
		}
	}
}

func main() {
	cblogger.Info("GCP Resource Test")
	handleVM()
}