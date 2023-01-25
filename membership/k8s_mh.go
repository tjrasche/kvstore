package membership

import (
	"context"
	"hash/fnv"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"os"
	"strconv"
)

type K8sHandler struct {
	replicas int
	k8s      *kubernetes.Clientset
}

func NewK8sHandler() *K8sHandler {
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	clientset, err := kubernetes.NewForConfig(config)
	replicas, err := strconv.Atoi(os.Getenv("REPLICAS"))
	return &K8sHandler{
		k8s:      clientset,
		replicas: replicas,
	}
}

// Members queries the current members from the k8s API
func (k K8sHandler) Members() ([]Node, error) {
	// we get all pods belonging to the statefulset by their pod labels for now, needs nicer solution
	podresult, err := k.k8s.CoreV1().Pods(os.Getenv("POD_NAMESPACE")).List(context.Background(), v1.ListOptions{
		LabelSelector: os.Getenv("POD_LABEL"),
	})
	if err != nil {
		return nil, err
	}

	nodes := make([]Node, k.replicas)

	// create nodes for every pod found belonging to the stateful set
	for _, pod := range podresult.Items {
		nodes = append(nodes, *NewNode(pod.Name, "http://"+pod.Status.PodIP+":8080"))
	}
	return nodes, nil

}

// CalculateReplication returns 3 id s of members where
func (k K8sHandler) CalculateReplication(key string) []Node {
	keyHash := hash(key)
	// using modulo to determine which node to write to should be ok when keys are aprox random
	nodeIndex := keyHash % uint32(k.replicas)
	allNodes, _ := k.Members()
	nodesToReplicateTo := make([]Node, 3)

	// add node calculated thorugh hash modulo

	nodesToReplicateTo = append(nodesToReplicateTo, allNodes[nodeIndex])

	// were building a sort off "ring topology light" here
	// choose a value before and after the determined node index to write replicas too
	// first find out id before the determined node in the ring
	var nodeLtIndex int
	if nodeIndex == 0 {
		nodeLtIndex = len(allNodes)
	} else {
		nodeLtIndex = int(nodeIndex - 1)
	}
	// add ltNode to returnValue
	nodesToReplicateTo = append(nodesToReplicateTo, allNodes[nodeLtIndex])

	// now find out id after the determined node in the ring
	var nodeGtIndex int
	if int(nodeIndex) == len(allNodes) {
		nodeGtIndex = 0
	} else {
		nodeGtIndex = int(nodeIndex + 1)
	}
	// add gtNode to returnValue
	nodesToReplicateTo = append(nodesToReplicateTo, allNodes[nodeGtIndex])

	return nodesToReplicateTo

}

// little helper function to get deterministic  integer hashes of strings
func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}
