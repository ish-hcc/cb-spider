// Tencent Driver of CB-Spider.
// The CB-Spider is a sub-Framework of the Cloud-Barista Multi-Cloud Project.
// The CB-Spider Mission is to connect all the clouds with a single interface.
//
//   - Cloud-Barista: https://github.com/cloud-barista
//
// This is Tencent Driver.
//
// by CB-Spider Team, 2022.09.

package connect

import (
	cblog "github.com/cloud-barista/cb-log"
	trs "github.com/cloud-barista/cb-spider/cloud-control-manager/cloud-driver/drivers/tencent/resources"
	idrv "github.com/cloud-barista/cb-spider/cloud-control-manager/cloud-driver/interfaces"
	irs "github.com/cloud-barista/cb-spider/cloud-control-manager/cloud-driver/interfaces/resources"
	"github.com/sirupsen/logrus"

	cbs "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cbs/v20170312"
	clb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb/v20180317"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
	tag "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tag/v20180813"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"

	tke "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/tke/v20180525"

	"errors"
)

type TencentCloudConnection struct {
	CredentialInfo   idrv.CredentialInfo
	Region           idrv.RegionInfo
	VNetworkClient   *vpc.Client
	NLBClient        *clb.Client
	VMClient         *cvm.Client
	KeyPairClient    *cvm.Client
	ImageClient      *cvm.Client
	SecurityClient   *vpc.Client
	VmSpecClient     *cvm.Client
	DiskClient       *cbs.Client
	MyImageClient    *cvm.Client
	RegionZoneClient *cvm.Client
	TagClient        *tag.Client
	ClusterClient    *tke.Client
}

var cblogger *logrus.Logger

func init() {
	// cblog is a global variable.
	cblogger = cblog.GetLogger("CB-SPIDER TencentCloudConnection")
}

func (cloudConn *TencentCloudConnection) CreateKeyPairHandler() (irs.KeyPairHandler, error) {
	cblogger.Info("Start CreateKeyPairHandler()")

	keyPairHandler := trs.TencentKeyPairHandler{Region: cloudConn.Region, Client: cloudConn.KeyPairClient}

	return &keyPairHandler, nil
}

func (cloudConn *TencentCloudConnection) CreateVMHandler() (irs.VMHandler, error) {
	cblogger.Info("Start CreateVMHandler()")

	vmHandler := trs.TencentVMHandler{Region: cloudConn.Region, Client: cloudConn.VMClient, DiskClient: cloudConn.DiskClient, VPCClient: cloudConn.VNetworkClient}
	return &vmHandler, nil
}

func (cloudConn *TencentCloudConnection) IsConnected() (bool, error) {
	return true, nil
}
func (cloudConn *TencentCloudConnection) Close() error {
	return nil
}

func (cloudConn *TencentCloudConnection) CreateVPCHandler() (irs.VPCHandler, error) {
	cblogger.Info("Start")
	handler := trs.TencentVPCHandler{Region: cloudConn.Region, Client: cloudConn.VNetworkClient}

	return &handler, nil
}

func (cloudConn *TencentCloudConnection) CreateNLBHandler() (irs.NLBHandler, error) {
	cblogger.Info("Start")
	handler := trs.TencentNLBHandler{Region: cloudConn.Region, Client: cloudConn.NLBClient, VpcClient: cloudConn.VNetworkClient}

	return &handler, nil
}

func (cloudConn *TencentCloudConnection) CreateImageHandler() (irs.ImageHandler, error) {
	cblogger.Info("Start")
	handler := trs.TencentImageHandler{Region: cloudConn.Region, Client: cloudConn.ImageClient}

	return &handler, nil
}

func (cloudConn *TencentCloudConnection) CreateSecurityHandler() (irs.SecurityHandler, error) {
	cblogger.Info("Start")
	handler := trs.TencentSecurityHandler{Region: cloudConn.Region, Client: cloudConn.SecurityClient, TagClient: cloudConn.TagClient}

	return &handler, nil
}

func (cloudConn *TencentCloudConnection) CreateVMSpecHandler() (irs.VMSpecHandler, error) {
	cblogger.Info("Start")
	handler := trs.TencentVmSpecHandler{Region: cloudConn.Region, Client: cloudConn.VmSpecClient}
	return &handler, nil
}

func (cloudConn *TencentCloudConnection) CreateDiskHandler() (irs.DiskHandler, error) {

	cblogger.Info("Start")
	handler := trs.TencentDiskHandler{Region: cloudConn.Region, Client: cloudConn.DiskClient}

	return &handler, nil
}

func (cloudConn *TencentCloudConnection) CreateMyImageHandler() (irs.MyImageHandler, error) {
	cblogger.Info("Start")
	handler := trs.TencentMyImageHandler{Region: cloudConn.Region, Client: cloudConn.MyImageClient, CbsClient: cloudConn.DiskClient}

	return &handler, nil
}

func (cloudConn *TencentCloudConnection) CreateClusterHandler() (irs.ClusterHandler, error) {
	// temp
	// getEnv & Setting
	clusterHandler := trs.TencentClusterHandler{RegionInfo: cloudConn.Region, CredentialInfo: cloudConn.CredentialInfo}

	return &clusterHandler, nil

}

func (cloudConn *TencentCloudConnection) CreateAnyCallHandler() (irs.AnyCallHandler, error) {
	return nil, errors.New("Tencent Driver: not implemented")
}

func (cloudConn *TencentCloudConnection) CreateRegionZoneHandler() (irs.RegionZoneHandler, error) {
	handler := trs.TencentRegionZoneHandler{Region: cloudConn.Region, Client: cloudConn.RegionZoneClient}
	return &handler, nil
}

func (cloudConn *TencentCloudConnection) CreatePriceInfoHandler() (irs.PriceInfoHandler, error) {
	handler := trs.TencentPriceInfoHandler{Region: cloudConn.Region, Client: cloudConn.VMClient}
	return &handler, nil
}

func (cloudConn *TencentCloudConnection) CreateTagHandler() (irs.TagHandler, error) {
	handler := trs.TencentTagHandler{
		Region:    cloudConn.Region,
		TagClient: cloudConn.TagClient,
		// below client is for validate resources
		VNetworkClient: cloudConn.VNetworkClient,
		VMClient:       cloudConn.VMClient,
		NLBClient:      cloudConn.NLBClient,
		DiskClient:     cloudConn.DiskClient,
		ClusterClient:  cloudConn.ClusterClient,
	}
	return &handler, nil
}

func (cloudConn *TencentCloudConnection) CreateMonitoringHandler() (irs.MonitoringHandler, error) {
	return nil, errors.New("Tencent Driver: not implemented")
}
