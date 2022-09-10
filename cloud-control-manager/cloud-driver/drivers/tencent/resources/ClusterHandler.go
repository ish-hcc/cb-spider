// Proof of Concepts of CB-Spider.
// The CB-Spider is a sub-Framework of the Cloud-Barista Multi-Cloud Project.
// The CB-Spider Mission is to connect all the clouds with a single interface.
//
//      * Cloud-Barista: https://github.com/cloud-barista
//
// This is a Cloud Driver Example for PoC Test.
//
// by ETRI, 2022.08.

package resources

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	call "github.com/cloud-barista/cb-spider/cloud-control-manager/cloud-driver/call-log"
	"github.com/cloud-barista/cb-spider/cloud-control-manager/cloud-driver/drivers/tencent/utils/tencent"
	idrv "github.com/cloud-barista/cb-spider/cloud-control-manager/cloud-driver/interfaces"
	irs "github.com/cloud-barista/cb-spider/cloud-control-manager/cloud-driver/interfaces/resources"

	"github.com/jeremywohl/flatten"
	"github.com/sirupsen/logrus"
	tke "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/tke/v20180525"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
)

// tempCalllogger
// 공통로거 만들기 이전까지 사용
var once sync.Once
var tempCalllogger *logrus.Logger

func init() {
	once.Do(func() {
		tempCalllogger = call.GetLogger("HISCALL")
	})
}

type TencentClusterHandler struct {
	RegionInfo     idrv.RegionInfo
	CredentialInfo idrv.CredentialInfo
}

func (clusterHandler *TencentClusterHandler) CreateCluster(clusterReqInfo irs.ClusterInfo) (irs.ClusterInfo, error) {
	cblogger.Info("Tencent Cloud Driver: called CreateCluster()")
	callLogInfo := getCallLogScheme(clusterHandler.RegionInfo.Region, call.CLUSTER, "CreateCluster()", "CreateCluster()")

	// 클러스터 생성 요청 변환
	request, err := getCreateClusterRequest(clusterReqInfo)
	if err != nil {
		cblogger.Error(err)
		return irs.ClusterInfo{}, err
	}

	start := call.Start()
	res, err := tencent.CreateCluster(clusterHandler.CredentialInfo.ClientId, clusterHandler.CredentialInfo.ClientSecret, clusterHandler.RegionInfo.Region, request)
	loggingInfo(callLogInfo, start)
	if err != nil {
		cblogger.Error(err)
		loggingError(callLogInfo, err)
		return irs.ClusterInfo{}, err
	}
	println(res.ToJsonString())

	// var response_json_obj map[string]interface{}
	// json.Unmarshal([]byte(response_json_str), &response_json_obj)
	// cluster_id := response_json_obj["cluster_id"].(string)
	cluster_info, err := getClusterInfo(clusterHandler.CredentialInfo.ClientId, clusterHandler.CredentialInfo.ClientSecret, clusterHandler.RegionInfo.Region, *res.Response.ClusterId)
	if err != nil {
		return irs.ClusterInfo{}, err
	}

	// // 리턴할 ClusterInfo 만들기
	// // 일단은 단순하게 만들어서 반환한다.
	// // 추후에 정보 추가 필요

	// // NodeGroup 생성 정보가 있는경우 생성을 시도한다.
	// // 문제는 Cluster 생성이 완료되어야 NodeGroup 생성이 가능하다.
	// // Cluster 생성이 완료되려면 최소 10분 이상 걸린다.
	// // 성공할때까지 반복하면서 생성을 시도해야 하는가?
	for _, node_group := range clusterReqInfo.NodeGroupList {
		res, err := clusterHandler.AddNodeGroup(clusterReqInfo.IId, node_group)
		if err != nil {
			cblogger.Error(err)
			return irs.ClusterInfo{}, err
		}
		printFlattenJSON(res)
	}

	return *cluster_info, nil
}

func (clusterHandler *TencentClusterHandler) ListCluster() ([]*irs.ClusterInfo, error) {
	cblogger.Info("Tencent Cloud Driver: called ListCluster()")
	callLogInfo := getCallLogScheme(clusterHandler.RegionInfo.Region, call.CLUSTER, "ListCluster()", "ListCluster()")

	start := call.Start()
	res, err := tencent.GetClusters(clusterHandler.CredentialInfo.ClientId, clusterHandler.CredentialInfo.ClientSecret, clusterHandler.RegionInfo.Region)
	loggingInfo(callLogInfo, start)
	if err != nil {
		return nil, err
	}

	cluster_info_list := make([]*irs.ClusterInfo, *res.Response.TotalCount)
	for i, cluster := range res.Response.Clusters {
		cluster_info_list[i], err = getClusterInfo(clusterHandler.CredentialInfo.ClientId, clusterHandler.CredentialInfo.ClientSecret, clusterHandler.RegionInfo.Region, *cluster.ClusterId)
		if err != nil {
			return nil, err
		}
	}

	return cluster_info_list, nil
}

func (clusterHandler *TencentClusterHandler) GetCluster(clusterIID irs.IID) (irs.ClusterInfo, error) {
	cblogger.Info("Tencent Cloud Driver: called GetCluster()")
	callLogInfo := getCallLogScheme(clusterHandler.RegionInfo.Region, call.CLUSTER, clusterIID.NameId, "GetCluster()")

	start := call.Start()
	cluster_info, err := getClusterInfo(clusterHandler.CredentialInfo.ClientId, clusterHandler.CredentialInfo.ClientSecret, clusterHandler.RegionInfo.Region, clusterIID.SystemId)
	loggingInfo(callLogInfo, start)
	if err != nil {
		return irs.ClusterInfo{}, err
	}

	return *cluster_info, nil
}

func (clusterHandler *TencentClusterHandler) DeleteCluster(clusterIID irs.IID) (bool, error) {
	cblogger.Info("Tencent Cloud Driver: called DeleteCluster()")
	callLogInfo := getCallLogScheme(clusterHandler.RegionInfo.Region, call.CLUSTER, clusterIID.NameId, "DeleteCluster()")

	start := call.Start()
	res, err := tencent.DeleteCluster(clusterHandler.CredentialInfo.ClientId, clusterHandler.CredentialInfo.ClientSecret, clusterHandler.RegionInfo.Region, clusterIID.SystemId)
	loggingInfo(callLogInfo, start)
	if err != nil {
		return false, err
	}
	println(res)

	return true, nil
}

func (clusterHandler *TencentClusterHandler) AddNodeGroup(clusterIID irs.IID, nodeGroupReqInfo irs.NodeGroupInfo) (irs.NodeGroupInfo, error) {
	cblogger.Info("Tencent Cloud Driver: called AddNodeGroup()")

	callLogInfo := getCallLogScheme(clusterHandler.RegionInfo.Region, call.CLUSTER, clusterIID.NameId, "AddNodeGroup()")

	// 노드 그룹 생성 요청 변환
	// get cluster info. to get security_group_id

	request, err := getNodeGroupRequest(clusterIID.SystemId, nodeGroupReqInfo)
	if err != nil {
		return irs.NodeGroupInfo{}, err
	}

	start := call.Start()
	response, err := tencent.CreateNodeGroup(clusterHandler.CredentialInfo.ClientId, clusterHandler.CredentialInfo.ClientSecret, clusterHandler.RegionInfo.Region, request)
	loggingInfo(callLogInfo, start)
	if err != nil {
		return irs.NodeGroupInfo{}, err
	}
	printFlattenJSON(response)

	node_group_info, err := getNodeGroupInfo(clusterHandler.CredentialInfo.ClientId, clusterHandler.CredentialInfo.ClientSecret, clusterHandler.RegionInfo.Region, clusterIID.SystemId, *response.Response.NodePoolId)
	if err != nil {
		return irs.NodeGroupInfo{}, err
	}

	return *node_group_info, nil
}

func (clusterHandler *TencentClusterHandler) ListNodeGroup(clusterIID irs.IID) ([]*irs.NodeGroupInfo, error) {
	cblogger.Info("Tencent Cloud Driver: called ListNodeGroup()")
	callLogInfo := getCallLogScheme(clusterHandler.RegionInfo.Region, call.CLUSTER, clusterIID.NameId, "ListNodeGroup()")

	node_group_info_list := []*irs.NodeGroupInfo{}

	start := call.Start()
	res, err := tencent.ListNodeGroup(clusterHandler.CredentialInfo.ClientId, clusterHandler.CredentialInfo.ClientSecret, clusterHandler.RegionInfo.Region, clusterIID.SystemId)
	loggingInfo(callLogInfo, start)
	if err != nil {
		return node_group_info_list, err
	}

	for _, node_group := range res.Response.NodePoolSet {
		node_group_info, err := getNodeGroupInfo(clusterHandler.CredentialInfo.ClientId, clusterHandler.CredentialInfo.ClientSecret, clusterHandler.RegionInfo.Region, clusterIID.SystemId, *node_group.NodePoolId)
		if err != nil {
			return nil, err
		}
		node_group_info_list = append(node_group_info_list, node_group_info)
	}

	return node_group_info_list, nil
}

func (clusterHandler *TencentClusterHandler) GetNodeGroup(clusterIID irs.IID, nodeGroupIID irs.IID) (irs.NodeGroupInfo, error) {
	cblogger.Info("Tencent Cloud Driver: called GetNodeGroup()")
	callLogInfo := getCallLogScheme(clusterHandler.RegionInfo.Region, call.CLUSTER, clusterIID.NameId, "GetNodeGroup()")

	start := call.Start()
	temp, err := getNodeGroupInfo(clusterHandler.CredentialInfo.ClientId, clusterHandler.CredentialInfo.ClientSecret, clusterHandler.RegionInfo.Region, clusterIID.SystemId, nodeGroupIID.SystemId)
	loggingInfo(callLogInfo, start)
	if err != nil {
		return irs.NodeGroupInfo{}, err
	}

	return *temp, nil
}

func (clusterHandler *TencentClusterHandler) SetNodeGroupAutoScaling(clusterIID irs.IID, nodeGroupIID irs.IID, on bool) (bool, error) {
	cblogger.Info("Tencent Cloud Driver: called SetNodeGroupAutoScaling()")

	temp, err := tencent.SetNodeGroupAutoScaling(clusterHandler.CredentialInfo.ClientId, clusterHandler.CredentialInfo.ClientSecret, clusterHandler.RegionInfo.Region, clusterIID.SystemId, nodeGroupIID.SystemId, on)
	if err != nil {
		println(err)
		return false, err
	}
	println(temp.ToJsonString())

	return true, nil
}

func (clusterHandler *TencentClusterHandler) ChangeNodeGroupScaling(clusterIID irs.IID, nodeGroupIID irs.IID, desiredNodeSize int, minNodeSize int, maxNodeSize int) (irs.NodeGroupInfo, error) {
	cblogger.Info("Tencent Cloud Driver: called ChangeNodeGroupScaling()")

	// nodepool.AutoscalingGroupId
	nodegroup, err := tencent.GetNodeGroup(clusterHandler.CredentialInfo.ClientId, clusterHandler.CredentialInfo.ClientSecret, clusterHandler.RegionInfo.Region, clusterIID.SystemId, nodeGroupIID.SystemId)
	if err != nil {
		return irs.NodeGroupInfo{}, err
	}

	temp, err := tencent.ChangeNodeGroupScaling(clusterHandler.CredentialInfo.ClientId, clusterHandler.CredentialInfo.ClientSecret, clusterHandler.RegionInfo.Region, *nodegroup.Response.NodePool.AutoscalingGroupId, uint64(desiredNodeSize), uint64(minNodeSize), uint64(maxNodeSize))
	if err != nil {
		println(err)
	}
	println(temp.ToJsonString())

	node_group_info, err := getNodeGroupInfo(clusterHandler.CredentialInfo.ClientId, clusterHandler.CredentialInfo.ClientSecret, clusterHandler.RegionInfo.Region, clusterIID.SystemId, nodeGroupIID.SystemId)
	if err != nil {
		return irs.NodeGroupInfo{}, err
	}

	return *node_group_info, nil
}

func (clusterHandler *TencentClusterHandler) RemoveNodeGroup(clusterIID irs.IID, nodeGroupIID irs.IID) (bool, error) {
	cblogger.Info("Tencent Cloud Driver: called RemoveNodeGroup()")
	callLogInfo := getCallLogScheme(clusterHandler.RegionInfo.Region, call.CLUSTER, clusterIID.NameId, "RemoveNodeGroup()")

	start := call.Start()
	res, err := tencent.DeleteNodeGroup(clusterHandler.CredentialInfo.ClientId, clusterHandler.CredentialInfo.ClientSecret, clusterHandler.RegionInfo.Region, clusterIID.SystemId, nodeGroupIID.SystemId)
	loggingInfo(callLogInfo, start)
	if err != nil {
		return false, err
	}
	println(res)

	return true, nil
}

func (clusterHandler *TencentClusterHandler) UpgradeCluster(clusterIID irs.IID, newVersion string) (irs.ClusterInfo, error) {
	cblogger.Info("Tencent Cloud Driver: called UpgradeCluster()")
	//callLogInfo := getCallLogScheme(clusterHandler.RegionInfo.Region, call.CLUSTER, clusterIID.NameId, "UpgradeCluster()")

	//version := "1.22.5"
	res, err := tencent.UpgradeCluster(clusterHandler.CredentialInfo.ClientId, clusterHandler.CredentialInfo.ClientSecret, clusterHandler.RegionInfo.Region, clusterIID.SystemId, newVersion)
	if err != nil {
		println(err.Error())
		//[TencentCloudSDKError] Code=InvalidParameter.Param,
		//Message=PARAM_ERROR(unsupported convert 1.20.6 to 1.22.5),
		return irs.ClusterInfo{}, err
	}
	println(res.ToJsonString())

	clusterInfo, err := getClusterInfo(clusterHandler.CredentialInfo.ClientId, clusterHandler.CredentialInfo.ClientSecret, clusterHandler.RegionInfo.Region, clusterIID.SystemId)
	if err != nil {
		return irs.ClusterInfo{}, err
	}

	return *clusterInfo, nil
}

func getClusterInfo(access_key string, access_secret string, region_id string, cluster_id string) (*irs.ClusterInfo, error) {
	var err error = nil
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("getClusterInfo() -> %v", r)
		}
	}()

	res, err := tencent.GetCluster(access_key, access_secret, region_id, cluster_id)
	if err != nil {
		return nil, err
	}

	if *res.Response.TotalCount == 0 {
		return nil, fmt.Errorf("cluster[%s] does not exist", cluster_id)
	}

	printFlattenJSON(res)

	// // k,v 추출
	// // k,v 변환 규칙 작성 [k,v]:[ClusterInfo.k, ClusterInfo.v]
	// // 변환 규칙에 따라 k,v 변환

	// https://intl.cloud.tencent.com/document/api/457/32022#ClusterStatus
	// Cluster status (Running, Creating, Idling or Abnormal)

	health_status := *res.Response.Clusters[0].ClusterStatus
	cluster_status := irs.ClusterActive
	if strings.EqualFold(health_status, "Creating") {
		cluster_status = irs.ClusterCreating
	} else if strings.EqualFold(health_status, "Creating") {
		cluster_status = irs.ClusterUpdating
	} else if strings.EqualFold(health_status, "Abnormal") {
		cluster_status = irs.ClusterInactive
	} else if strings.EqualFold(health_status, "Running") {
		cluster_status = irs.ClusterActive
	}
	// } else if strings.EqualFold(health_status, "") { // tencent has no "delete" state
	// // 	cluster_status = irs.ClusterDeleting
	println(cluster_status)

	// "2022-09-09T13:10:06Z",
	created_at := *res.Response.Clusters[0].CreatedTime // 2022-09-08T09:02:16+08:00,
	datetime, err := time.Parse(time.RFC3339, created_at)
	if err != nil {
		panic(err)
	}

	// "Response.Clusters.0.ClusterName": "cluster-x1",
	// "Response.Clusters.0.ClusterVersion": "1.22.5",
	// "Response.Clusters.0.ClusterNetworkSettings.VpcId": "vpc-q1c6fr9e",
	// "Response.Clusters.0.ClusterStatus": "Creating",
	// "Response.Clusters.0.CreatedTime": "2022-09-09T13:10:06Z",

	cluster_info := &irs.ClusterInfo{
		IId: irs.IID{
			NameId:   *res.Response.Clusters[0].ClusterName,
			SystemId: *res.Response.Clusters[0].ClusterId,
		},
		Version: *res.Response.Clusters[0].ClusterVersion,
		Network: irs.NetworkInfo{
			VpcIID: irs.IID{
				NameId:   "",
				SystemId: *res.Response.Clusters[0].ClusterVersion,
			},
		},
		Status:      cluster_status,
		CreatedTime: datetime,
		// KeyValueList: []irs.KeyValue{}, // flatten data 입력하기
	}
	println(cluster_info)

	// NodeGroups
	res2, err := tencent.ListNodeGroup(access_key, access_secret, region_id, cluster_id)
	if err != nil {
		return nil, err
	}
	print(res.ToJsonString())

	// // k,v 추출
	// // k,v 변환 규칙 작성 [k,v]:[NodeGroup.k, NodeGroup.v]
	// // 변환 규칙에 따라 k,v 변환
	// flat, err = flatten.FlattenString(node_groups_json_str, "", flatten.DotStyle)
	// if err != nil {
	// 	return nil, err
	// }
	// println(flat)

	for _, nodepool := range res2.Response.NodePoolSet {
		node_group_info, err := getNodeGroupInfo(access_key, access_secret, region_id, cluster_id, *nodepool.NodePoolId)
		if err != nil {
			return nil, err
		}
		cluster_info.NodeGroupList = append(cluster_info.NodeGroupList, *node_group_info)
	}

	//return cluster_info, nil

	return cluster_info, nil
}

func printFlattenJSON(json_obj interface{}) {
	temp, err := json.MarshalIndent(json_obj, "", "  ")
	if err != nil {
		println(err)
	} else {
		flat, err := flatten.FlattenString(string(temp), "", flatten.DotStyle)
		if err != nil {
			println(err)
		} else {
			println(flat)
		}
	}
}

func getNodeGroupInfo(access_key, access_secret, region_id, cluster_id, node_group_id string) (*irs.NodeGroupInfo, error) {
	var err error = nil
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("getNodeGroupInfo() -> %v", r)
		}
	}()

	res, err := tencent.GetNodeGroup(access_key, access_secret, region_id, cluster_id, node_group_id)
	if err != nil {
		return nil, err
	}
	printFlattenJSON(res)

	launch_config, err := tencent.GetLaunchConfiguration(access_key, access_secret, region_id, *res.Response.NodePool.LaunchConfigurationId)
	if err != nil {
		return nil, err
	}
	printFlattenJSON(launch_config)

	auto_scaling_group, err := tencent.GetAutoScalingGroup(access_key, access_secret, region_id, *res.Response.NodePool.AutoscalingGroupId)
	if err != nil {
		return nil, err
	}
	printFlattenJSON(auto_scaling_group)

	// nodepool LifeState
	// The lifecycle state of the current node pool.
	// Valid values: creating, normal, updating, deleting, and deleted.
	health_status := *res.Response.NodePool.LifeState
	status := irs.NodeGroupActive
	if strings.EqualFold(health_status, "normal") {
		status = irs.NodeGroupActive
	} else if strings.EqualFold(health_status, "creating") {
		status = irs.NodeGroupUpdating
	} else if strings.EqualFold(health_status, "removing") {
		status = irs.NodeGroupUpdating // removing is a kind of updating?
	} else if strings.EqualFold(health_status, "deleting") {
		status = irs.NodeGroupDeleting
	} else if strings.EqualFold(health_status, "updating") {
		status = irs.NodeGroupUpdating
	}

	println(status)

	auto_scale_enalbed := false
	if strings.EqualFold("Response.AutoScalingGroupSet.0.EnabledStatus", "ENABLED") {
		auto_scale_enalbed = true
	}

	node_group_info := irs.NodeGroupInfo{
		IId: irs.IID{
			NameId:   *res.Response.NodePool.Name,
			SystemId: *res.Response.NodePool.NodePoolId,
		},
		ImageIID: irs.IID{
			NameId:   "",
			SystemId: *launch_config.Response.LaunchConfigurationSet[0].ImageId,
		},
		VMSpecName:      *launch_config.Response.LaunchConfigurationSet[0].InstanceType,
		RootDiskType:    *launch_config.Response.LaunchConfigurationSet[0].SystemDisk.DiskType,
		RootDiskSize:    fmt.Sprintf("%d", *launch_config.Response.LaunchConfigurationSet[0].SystemDisk.DiskSize),
		KeyPairIID:      irs.IID{NameId: "", SystemId: ""}, // not available
		Status:          status,
		OnAutoScaling:   auto_scale_enalbed,
		MinNodeSize:     int(*auto_scaling_group.Response.AutoScalingGroupSet[0].MinSize),
		MaxNodeSize:     int(*auto_scaling_group.Response.AutoScalingGroupSet[0].MaxSize),
		DesiredNodeSize: int(*auto_scaling_group.Response.AutoScalingGroupSet[0].DesiredCapacity),
		NodeList:        []irs.IID{},      // to be implemented
		KeyValueList:    []irs.KeyValue{}, // to be implemented
	}

	return &node_group_info, nil
}

func getCreateClusterRequest(clusterInfo irs.ClusterInfo) (*tke.CreateClusterRequest, error) {

	var err error = nil
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("recovered: %v", r)
		}
	}()

	// clusterInfo := irs.ClusterInfo{
	// 	IId: irs.IID{
	// 		NameId:   "cluster-x1",
	// 		SystemId: "",
	// 	},
	// 	Version: "1.22.5",
	// 	Network: irs.NetworkInfo{
	// 		VpcIID: irs.IID{NameId: "", SystemId: "vpc-q1c6fr9e"},
	// 	},
	// 	KeyValueList: []irs.KeyValue{
	// 		{
	// 			Key:   "cluster_cidr", // 조회가능한 값이면, 내부에서 처리하는 코드 추가
	// 			Value: "172.20.0.0/16",
	// 		},
	// 	},
	// }

	cluster_cidr := "" // 172.X.0.0.16: X Range:16, 17, ... , 31
	for _, v := range clusterInfo.KeyValueList {
		switch v.Key {
		case "cluster_cidr":
			cluster_cidr = v.Value
		}
	}

	request := tke.NewCreateClusterRequest()
	request.ClusterCIDRSettings = &tke.ClusterCIDRSettings{
		ClusterCIDR: common.StringPtr(cluster_cidr), // 172.X.0.0.16: X Range:16, 17, ... , 31
	}
	request.ClusterBasicSettings = &tke.ClusterBasicSettings{
		ClusterName:    common.StringPtr(clusterInfo.IId.NameId),
		VpcId:          common.StringPtr(clusterInfo.Network.VpcIID.SystemId),
		ClusterVersion: common.StringPtr(clusterInfo.Version), //option, version: 1.22.5
	}
	request.ClusterType = common.StringPtr("MANAGED_CLUSTER") //default value

	return request, err
}

func getNodeGroupRequest(cluster_id string, nodeGroupReqInfo irs.NodeGroupInfo) (*tke.CreateClusterNodePoolRequest, error) {
	var err error = nil
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("recovered: %v", r)
		}
	}()

	// KeyValueList: []irs.KeyValue{
	// 	{
	// 		Key:   "security_group_id", // security_group_id는 cluster_info 정보에 있음. 이것을 어떻게 참조할지가 문제.
	// 		Value: "sg-46eef229",
	// 	},
	// 	{
	// 		Key:   "subnet_id", // cluster_info 에서 참조
	// 		Value: "subnet-rl79gxhv",
	// 	},
	// 	{
	// 		Key:   "vpc_id", // cluster_info 정보에 있음. 조회해서 처리 가능
	// 		Value: "vpc-q1c6fr9e",
	// 	},
	// },
	security_group_id := "" // 입력해야함
	subnet_id := ""         // 입력해야함
	vpc_id := ""            // 클러스터 정보에서 조회 가능
	for _, v := range nodeGroupReqInfo.KeyValueList {
		switch v.Key {
		case "security_group_id":
			security_group_id = v.Value
		case "subnet_id":
			subnet_id = v.Value
		case "vpc_id":
			vpc_id = v.Value
		}
	}

	launch_config_json_str := `{
		"InstanceType": "%s",
		"SecurityGroupIds": ["%s"]
	}`

	// launch_config_json_str := `{
	// 	"InstanceType": "%s",
	// 	"SecurityGroupIds": ["%s"],
	// 	"ImageId":"%s"
	// }`

	//"ImageId":""

	// security group id 는 cluster info 에서 지정한다.
	// 그런데 텐센트는 노드그룹 생성에서 지정해야한다.
	// node group info 네는 securityp group id 필드가 없다.
	// 그래서 일단 key/value 에 지정해서 사용한다. // issue

	// 추가
	// https://intl.cloud.tencent.com/document/api/377/30998
	//keypair:
	//image_id: "ImageId": "",

	//launch_config_json_str = fmt.Sprintf(launch_config_json_str, "S3.MEDIUM2", "sg-46eef229")
	// req.setLaunchConfigurationName("name");
	// req.setInstanceType("instance_type");
	// req.setImageId("image_id");
	// SystemDisk systemDisk1 = new SystemDisk();
	// systemDisk1.setDiskType("disk_type");
	// systemDisk1.setDiskSize(50L);
	// req.setSystemDisk(systemDisk1);
	// LoginSettings loginSettings1 = new LoginSettings();
	// loginSettings1.setPassword("password");
	// String[] keyIds1 = {"key_id"};
	// loginSettings1.setKeyIds(keyIds1);
	// req.setLoginSettings(loginSettings1);
	// String[] securityGroupIds1 = {"security_group"};
	// req.setSecurityGroupIds(securityGroupIds1);

	launch_config_json_str = fmt.Sprintf(launch_config_json_str, nodeGroupReqInfo.VMSpecName, security_group_id)

	auto_scaling_group_json_str := `{
		"MinSize": %d,
		"MaxSize": %d,			
		"DesiredCapacity": %d,
		"VpcId": "%s",
		"SubnetIds": ["%s"]
	}`
	auto_scaling_group_json_str = fmt.Sprintf(auto_scaling_group_json_str, 0, 3, 1, vpc_id, subnet_id)
	// auto_scaling_group_json_str = fmt.Sprintf(auto_scaling_group_json_str, nodeGroupReqInfo.MinNodeSize, nodeGroupReqInfo.MaxNodeSize, nodeGroupReqInfo.DesiredNodeSize, nodeGroupReqInfo.Network.VpcIID.SystemId, nodeGroupReqInfo.Network.SubnetIID.SystemId)

	disk_size, _ := strconv.ParseInt(nodeGroupReqInfo.RootDiskSize, 10, 64)

	// cluster_id, "cls-ke0ztn01_nodepool-x", lc_json_str, asc_json_str, true, "CLOUD_PREMIUM", 50
	// Instantiate a request object. You can further set the request parameters according to the API called and actual conditions
	request := tke.NewCreateClusterNodePoolRequest()
	request.Name = common.StringPtr(nodeGroupReqInfo.IId.NameId)
	request.ClusterId = common.StringPtr(cluster_id)
	request.LaunchConfigurePara = common.StringPtr(launch_config_json_str)
	request.AutoScalingGroupPara = common.StringPtr(auto_scaling_group_json_str)
	request.EnableAutoscale = common.BoolPtr(nodeGroupReqInfo.OnAutoScaling)
	request.InstanceAdvancedSettings = &tke.InstanceAdvancedSettings{
		DataDisks: []*tke.DataDisk{
			{
				DiskType: common.StringPtr(nodeGroupReqInfo.RootDiskType), //ex. "CLOUD_PREMIUM"
				DiskSize: common.Int64Ptr(disk_size),                      //ex. 50
			},
		},
	}
	return request, err
}

// getCallLogScheme(clusterHandler.RegionInfo.Region, call.CLUSTER, "ListCluster()", "ListCluster()")
func getCallLogScheme(region string, resourceType call.RES_TYPE, resourceName string, apiName string) call.CLOUDLOGSCHEMA {
	cblogger.Info(fmt.Sprintf("Call %s %s", call.TENCENT, apiName))
	return call.CLOUDLOGSCHEMA{
		CloudOS:      call.TENCENT,
		RegionZone:   region,
		ResourceType: resourceType,
		ResourceName: resourceName,
		CloudOSAPI:   apiName,
	}
}

func loggingError(hiscallInfo call.CLOUDLOGSCHEMA, err error) {
	hiscallInfo.ErrorMSG = err.Error()
	tempCalllogger.Info(call.String(hiscallInfo))
}

func loggingInfo(hiscallInfo call.CLOUDLOGSCHEMA, start time.Time) {
	hiscallInfo.ElapsedTime = call.Elapsed(start)
	tempCalllogger.Info(call.String(hiscallInfo))
}
