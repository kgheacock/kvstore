from tools import NodeContainer, NodeNetwork 
import subprocess
net1 = NodeNetwork("kv_subnet","10.10.0.0/16")
nodeList = {}
ipList = ["10.0.0."+str(x)+":3800" for x in range(0,20)]
repFactor = 3
shards = 10
for i in range(0,repFactor*shards):
    node = NodeContainer("10.0.0."+str(i),ipList[i],repFactor,"kv-store:4.0")
    net1.connect(node,ipList[i])
    print(net1.containers)
