package connect

import (
	"errors"
	cblog "github.com/cloud-barista/cb-log"
	mkrs "github.com/cloud-barista/cb-spider/cloud-control-manager/cloud-driver/drivers/mock/resources"
	idrv "github.com/cloud-barista/cb-spider/cloud-control-manager/cloud-driver/interfaces"
	irs "github.com/cloud-barista/cb-spider/cloud-control-manager/cloud-driver/interfaces/resources"
	"github.com/sirupsen/logrus"
)

var cblogger *logrus.Logger

func init() {
	// cblog is a global variable.
	cblogger = cblog.GetLogger("CB-SPIDER")
}

type MockConnection struct {
	Region   idrv.RegionInfo
	MockName string
}

func (cloudConn *MockConnection) CreateImageHandler() (irs.ImageHandler, error) {
	cblogger.Info("Mock Driver: called CreateImageHandler()!")
	handler := mkrs.MockImageHandler{cloudConn.MockName}
	return &handler, nil
}

func (cloudConn *MockConnection) CreateVMHandler() (irs.VMHandler, error) {
	cblogger.Info("Mock Driver: called CreateVMHandler()!")
	handler := mkrs.MockVMHandler{cloudConn.Region, cloudConn.MockName}
	return &handler, nil
}

func (cloudConn *MockConnection) CreateVPCHandler() (irs.VPCHandler, error) {
	cblogger.Info("Mock Driver: called CreateVPCHandler()!")
	handler := mkrs.MockVPCHandler{cloudConn.MockName}
	return &handler, nil
}

func (cloudConn MockConnection) CreateSecurityHandler() (irs.SecurityHandler, error) {
	cblogger.Info("Mock Driver: called CreateSecurityHandler()!")
	handler := mkrs.MockSecurityHandler{cloudConn.MockName}
	return &handler, nil
}

func (cloudConn *MockConnection) CreateKeyPairHandler() (irs.KeyPairHandler, error) {
	cblogger.Info("Mock Driver: called CreateKeyPairHandler()!")
	handler := mkrs.MockKeyPairHandler{cloudConn.MockName}
	return &handler, nil
}

func (cloudConn *MockConnection) CreateVMSpecHandler() (irs.VMSpecHandler, error) {
	cblogger.Info("Mock Driver: called CreateVMSpecHandler()!")
	handler := mkrs.MockVMSpecHandler{cloudConn.MockName}
	return &handler, nil
}

func (cloudConn *MockConnection) CreateNLBHandler() (irs.NLBHandler, error) {
	cblogger.Info("Mock Driver: called CreateNLBHandler()!")
	handler := mkrs.MockNLBHandler{cloudConn.MockName}
	return &handler, nil
}

func (cloudConn *MockConnection) IsConnected() (bool, error) {
	cblogger.Info("Mock Driver: called IsConnected()!")
	if cloudConn == nil {
		return false, nil
	}

	return true, nil
}

func (cloudConn *MockConnection) Close() error {
	cblogger.Info("Mock Driver: called Close()!")
	return nil
}

func (cloudConn *MockConnection) CreateDiskHandler() (irs.DiskHandler, error) {
	cblogger.Info("Mock Driver: called CreateDiskHandler()!")
	handler := mkrs.MockDiskHandler{cloudConn.MockName}
	return &handler, nil
}

func (cloudConn *MockConnection) CreateClusterHandler() (irs.ClusterHandler, error) {
	cblogger.Info("Mock Driver: called CreateClusterHandler()!")
	handler := mkrs.MockClusterHandler{cloudConn.MockName}
	return &handler, nil
}

func (cloudConn *MockConnection) CreateMyImageHandler() (irs.MyImageHandler, error) {
	cblogger.Info("Mock Driver: called CreateMyImageHandler()!")
	handler := mkrs.MockMyImageHandler{cloudConn.MockName}
	return &handler, nil
}

func (cloudConn *MockConnection) CreateAnyCallHandler() (irs.AnyCallHandler, error) {
	cblogger.Info("Mock Driver: called CreateAnyCallHandler()!")
	handler := mkrs.MockAnyCallHandler{cloudConn.MockName}
	return &handler, nil
}

// CreateRegionZoneHandler implements connect.CloudConnection.
func (cloudConn *MockConnection) CreateRegionZoneHandler() (irs.RegionZoneHandler, error) {
	cblogger.Info("Mock Driver: called CreateRegionZoneHandler()!")
	handler := mkrs.MockRegionZoneHandler{Region: cloudConn.Region, MockName: cloudConn.MockName}
	return &handler, nil
}

func (cloudConn *MockConnection) CreatePriceInfoHandler() (irs.PriceInfoHandler, error) {
	cblogger.Info("Mock Driver: called CreatePriceInfoHandler()!")
	handler := mkrs.MockPriceInfoHandler{Region: cloudConn.Region, MockName: cloudConn.MockName}
	return &handler, nil
}

func (cloudConn *MockConnection) CreateTagHandler() (irs.TagHandler, error) {
	cblogger.Info("Mock Driver: called CreateTagHandler()!")
	handler := mkrs.MockTagHandler{MockName: cloudConn.MockName}
	return &handler, nil
}

func (cloudConn *MockConnection) CreateMonitoringHandler() (irs.MonitoringHandler, error) {
	return nil, errors.New("Mock Driver: not implemented")
}
