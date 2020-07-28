
echo "####################################################################"
echo "## Full Test Scripts for CB-Spider IID Working Version - 2020.04.22."
echo "##   1. VPC: Create -> List -> Get"
echo "##   2. SecurityGroup: Create -> List -> Get"
echo "##   3. KeyPair: Create -> List -> Get"
echo "##   4. VM: StartVM -> List -> Get -> ListStatus -> GetStatus -> Suspend -> Resume -> Reboot"
echo "## ---------------------------------"
echo "##   4. VM: Terminate(Delete)"
echo "##   3. KeyPair: Delete"
echo "##   2. SecurityGroup: Delete"
echo "##   1. VPC: Delete"
echo "####################################################################"

echo "####################################################################"
echo "## 1. VPC: Create -> List -> Get"
echo "####################################################################"
$CBSPIDER_ROOT/interface/spider vpc create --config $CBSPIDER_ROOT/interface/grpc_conf.yaml -i json -d \
    '{ 
      "ConnectionName":"'${CONN_CONFIG}'",
      "ReqInfo": { 
        "Name": "vpc-01", 
        "IPv4_CIDR": "'${IPv4_CIDR}'", 
        "SubnetInfoList": [ 
          { 
            "Name": "subnet-01", 
            "IPv4_CIDR": "'${IPv4_CIDR}'"
          } 
        ] 
      } 
    }'

$CBSPIDER_ROOT/interface/spider vpc list --config $CBSPIDER_ROOT/interface/grpc_conf.yaml --cname "${CONN_CONFIG}"

$CBSPIDER_ROOT/interface/spider vpc get --config $CBSPIDER_ROOT/interface/grpc_conf.yaml --cname "${CONN_CONFIG}" -n vpc-01

echo "#-----------------------------"

echo "####################################################################"
echo "## 2. SecurityGroup: Create -> List -> Get"
echo "####################################################################"
$CBSPIDER_ROOT/interface/spider security create --config $CBSPIDER_ROOT/interface/grpc_conf.yaml -i json -d \
    '{ 
      "ConnectionName":"'${CONN_CONFIG}'",
      "ReqInfo": { 
        "Name": "sg-01", 
        "VPCName": "vpc-01", 
        "SecurityRules": [ 
          {
            "FromPort": "1", 
            "ToPort" : "65535", 
            "IPProtocol" : "tcp", 
            "Direction" : "inbound"
          }
        ] 
      } 
    }'

$CBSPIDER_ROOT/interface/spider security list --config $CBSPIDER_ROOT/interface/grpc_conf.yaml --cname "${CONN_CONFIG}"

$CBSPIDER_ROOT/interface/spider security get --config $CBSPIDER_ROOT/interface/grpc_conf.yaml --cname "${CONN_CONFIG}" -n sg-01

echo "#-----------------------------"

echo "####################################################################"
echo "## 3. KeyPair: Create -> List -> Get"
echo "####################################################################"
$CBSPIDER_ROOT/interface/spider keypair create --config $CBSPIDER_ROOT/interface/grpc_conf.yaml -i json -d \
    '{ 
      "ConnectionName":"'${CONN_CONFIG}'",
      "ReqInfo": { 
        "Name": "keypair-01" 
      } 
    }'

$CBSPIDER_ROOT/interface/spider keypair list --config $CBSPIDER_ROOT/interface/grpc_conf.yaml --cname "${CONN_CONFIG}"

$CBSPIDER_ROOT/interface/spider keypair get --config $CBSPIDER_ROOT/interface/grpc_conf.yaml --cname "${CONN_CONFIG}" -n keypair-01

echo "#-----------------------------"

echo "####################################################################"
echo "## 4. VM: StartVM -> List -> Get -> ListStatus -> GetStatus -> Suspend -> Resume -> Reboot"
echo "####################################################################"
$CBSPIDER_ROOT/interface/spider vm start --config $CBSPIDER_ROOT/interface/grpc_conf.yaml -i json -d \
    '{ 
      "ConnectionName":"'${CONN_CONFIG}'",
      "ReqInfo": { 
        "Name": "vm-01", 
        "ImageName": "'${IMAGE_NAME}'", 
        "VPCName": "vpc-01", 
        "SubnetName": "subnet-01", 
        "SecurityGroupNames": [ "sg-01" ], 
        "VMSpecName": "'${SPEC_NAME}'", 
        "KeyPairName": "keypair-01"
      } 
    }'

echo "============== sleep 60 after start VM"
sleep 60

$CBSPIDER_ROOT/interface/spider vm list --config $CBSPIDER_ROOT/interface/grpc_conf.yaml --cname "${CONN_CONFIG}"

$CBSPIDER_ROOT/interface/spider vm get --config $CBSPIDER_ROOT/interface/grpc_conf.yaml --cname "${CONN_CONFIG}" -n vm-01

$CBSPIDER_ROOT/interface/spider vm liststatus --config $CBSPIDER_ROOT/interface/grpc_conf.yaml --cname "${CONN_CONFIG}"

$CBSPIDER_ROOT/interface/spider vm getstatus --config $CBSPIDER_ROOT/interface/grpc_conf.yaml --cname "${CONN_CONFIG}" -n vm-01

$CBSPIDER_ROOT/interface/spider vm control --config $CBSPIDER_ROOT/interface/grpc_conf.yaml --cname "${CONN_CONFIG}" -n vm-01 -a suspend

echo "============== sleep 60 after suspend VM"
sleep 60

$CBSPIDER_ROOT/interface/spider vm control --config $CBSPIDER_ROOT/interface/grpc_conf.yaml --cname "${CONN_CONFIG}" -n vm-01 -a resume
echo "============== sleep 30 after resume VM"
sleep 30

$CBSPIDER_ROOT/interface/spider vm control --config $CBSPIDER_ROOT/interface/grpc_conf.yaml --cname "${CONN_CONFIG}" -n vm-01 -a reboot

echo "============== sleep 60 after reboot VM"
sleep 60 
echo "#-----------------------------"


echo "####################################################################"
echo "####################################################################"
echo "####################################################################"

echo "####################################################################"
echo "## 4. VM: Terminate(Delete)"
echo "####################################################################"
$CBSPIDER_ROOT/interface/spider vm terminate --config $CBSPIDER_ROOT/interface/grpc_conf.yaml --cname "${CONN_CONFIG}" -n vm-01 --force false
echo "============== sleep 70 after delete VM"
sleep 70 

echo "####################################################################"
echo "## 3. KeyPair: Delete"
echo "####################################################################"
$CBSPIDER_ROOT/interface/spider keypair delete --config $CBSPIDER_ROOT/interface/grpc_conf.yaml --cname "${CONN_CONFIG}" -n keypair-01 --force false

echo "####################################################################"
echo "## 2. SecurityGroup: Delete"
echo "####################################################################"
$CBSPIDER_ROOT/interface/spider security delete --config $CBSPIDER_ROOT/interface/grpc_conf.yaml --cname "${CONN_CONFIG}" -n sg-01 --force false

echo "####################################################################"
echo "## 1. VPC: Delete"
echo "####################################################################"
$CBSPIDER_ROOT/interface/spider vpc delete --config $CBSPIDER_ROOT/interface/grpc_conf.yaml --cname "${CONN_CONFIG}" -n vpc-01 --force false


