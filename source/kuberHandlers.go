package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ghodss/yaml"
	"github.com/tidwall/gjson"
)

type Pod struct {
	PodName      string
	Age          string
	ImageName    string
	Status       string
	RestartCount string
}

type Job struct {
	JobName    string
	ImageName  string
	PullPolicy string
	Active     string
	StartTime  string
}

type Cronjob struct {
	CronjobName  string
	PullPolicy   string
	ImageName    string
	Scheduler    string
	LastSchedule string
	History      string
	Suspend      bool
}

type Daemonset struct {
	DaemonsetName string
	ImageName     string
	PullPolicy    string
	PodCount      string
}

type Deployment struct {
	DeploymentName string
	ReplicaCount   string
	PullPolicy     string
	ImageName      string
}

type StatefulSet struct {
	StatefulSetName string
	ReplicaCount    string
	PullPolicy      string
	ImageName       string
}

type Node struct {
	NodeName string
	IP       string
	// HostName string
	//      Status        string
	//      Labels  string
}

func MakeReq(url string, method string) string {

	token, err := ioutil.ReadFile("/run/secrets/kubernetes.io/serviceaccount/token")
	if err != nil {
		fmt.Print(err)
	}

	str := string(token)
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	client := &http.Client{}

	req, err := http.NewRequest(method, url, nil)
	TOKEN := fmt.Sprintf("%s %s ", "Bearer", str)
	req.Header.Add("Authorization", TOKEN)

	if err != nil {
		log.Fatalln(err)
	}
	resp, err2 := client.Do(req)

	if err2 != nil {
		log.Fatal("Error reading response. ", err)

	}

	// Read body from response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading response. ", err)
	}
	//return string(body)
	return string(body)
}

func Patch(url string, data string) {

	token, err := ioutil.ReadFile("/run/secrets/kubernetes.io/serviceaccount/token")
	if err != nil {
		fmt.Print(err)
	}

	str := string(token)
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	client := &http.Client{}

	var jsonData = []byte(data)
	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonData))
	TOKEN := fmt.Sprintf("%s %s ", "Bearer", str)

	req.Header.Add("Authorization", TOKEN)
	req.Header.Add("Content-Type", "application/strategic-merge-patch+json")
	//      req.Header.Add("data", data)

	if err != nil {
		log.Fatalln(err)
	}
	resp, err2 := client.Do(req)

	if err2 != nil {
		log.Fatal("Error reading response. ", err)

	}

	// Read body from response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading response. ", err)
	}
	fmt.Println(string(body))
}

func Create(url string, data string) {

	token, err := ioutil.ReadFile("/run/secrets/kubernetes.io/serviceaccount/token")
	if err != nil {
		fmt.Print(err)
	}

	str := string(token)
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	client := &http.Client{}

	var jsonData = []byte(data)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	TOKEN := fmt.Sprintf("%s %s ", "Bearer", str)

	req.Header.Add("Authorization", TOKEN)
	req.Header.Add("Content-Type", "application/json")

	if err != nil {
		log.Fatalln(err)
	}
	resp, err2 := client.Do(req)

	if err2 != nil {
		log.Fatal("Error reading response. ", err)

	}

	// Read body from response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading response. ", err)
	}
	fmt.Println(string(body))
}

func AgeCalculator(startTime string, endTime string) string {

	type Age struct {
		Year   int
		Month  int
		Day    int
		Hour   int
		Minute int
		Second int
	}
	var start, end Age
	var t1, t2 time.Time
	start.Year, _ = strconv.Atoi(startTime[0:4])
	start.Month, _ = strconv.Atoi(startTime[5:7])
	start.Day, _ = strconv.Atoi(startTime[8:10])
	start.Hour, _ = strconv.Atoi(startTime[11:13])
	start.Minute, _ = strconv.Atoi(startTime[14:16])
	start.Second, _ = strconv.Atoi(startTime[17:19])

	end.Year, _ = strconv.Atoi(endTime[0:4])
	end.Month, _ = strconv.Atoi(endTime[5:7])
	end.Day, _ = strconv.Atoi(endTime[8:10])
	end.Hour, _ = strconv.Atoi(endTime[11:13])
	end.Minute, _ = strconv.Atoi(endTime[14:16])
	end.Second, _ = strconv.Atoi(endTime[17:19])

	t1 = time.Date(start.Year, time.Month(start.Month), start.Day, start.Hour, start.Minute, start.Second, 0, time.UTC)

	t2 = time.Date(end.Year, time.Month(end.Month), end.Day, end.Hour, end.Minute, end.Second, 0, time.UTC)

	diff := t2.Sub(t1).String()
	return diff
}

func GetPods() []Pod {
	var pod Pod
	var Pods []Pod
	var startTime, endTime string
	podInfo := MakeReq("https://kubernetes/api/v1/namespaces/"+os.Getenv("NAMESPACE")+"/pods/", "GET")

	items := gjson.Get(podInfo, "items.#")
	itemCount := int(items.Int())
	for i := 0; i < itemCount; i++ {
		count := strconv.Itoa(i)
		pod.PodName = gjson.Get(podInfo, "items."+count+".metadata.name").String()
		pod.Status = gjson.Get(podInfo, "items."+count+".status.phase").String()
		pod.ImageName = gjson.Get(podInfo, "items."+count+".status.containerStatuses.0.image").String()
		pod.RestartCount = gjson.Get(podInfo, "items."+count+".status.containerStatuses.0.restartCount").String()
		startTime = gjson.Get(podInfo, "items."+count+".status.startTime").String()
		endTime = gjson.Get(podInfo, "items."+count+".status.containerStatuses.0.state.terminated.finishedAt").String()
		if endTime == "" {
			endTime = time.Now().String()
		}
		pod.Age = AgeCalculator(startTime, endTime)
		/*	if len(pod.PodName) > 50 {
				pod.PodName = pod.PodName[0:50]
			}
			if len(pod.ImageName) > 100 {
				pod.PodName = pod.PodName[0:100]
			} */
		Pods = append(Pods, pod)

	}

	return Pods

}

func DeletePod(podName string) {

	podInfo := MakeReq("https://kubernetes/api/v1/namespaces/"+os.Getenv("NAMESPACE")+"/pods/"+podName, "DELETE")
	fmt.Println(podInfo)
}

func GetLog(podName string) (string, string) {

	podLogs := MakeReq("https://kubernetes/api/v1/namespaces/"+os.Getenv("NAMESPACE")+"/pods/"+podName+"/log?pretty=true&timestamps=true", "GET")
	podStatus := MakeReq("https://kubernetes/api/v1/namespaces/"+os.Getenv("NAMESPACE")+"/pods/"+podName+"/status?pretty=true&timestamps=true", "GET")

	if len(podLogs) > 10000 {
		podLogs = podLogs[len(podLogs)-10000 : len(podLogs)]
	}

	return podLogs, podStatus
}

func GetNodes() []Node {
	var node Node
	var Nodes []Node

	nodeInfo := MakeReq("https://kubernetes/api/v1/nodes/", "GET")

	items := gjson.Get(nodeInfo, "items.#")
	itemCount := int(items.Int())
	for i := 0; i < itemCount; i++ {
		count := strconv.Itoa(i)
		node.NodeName = gjson.Get(nodeInfo, "items."+count+".metadata.name").String()
		node.IP = gjson.Get(nodeInfo, "items."+count+".status.addresses.0.address").String()
		Nodes = append(Nodes, node)

	}

	return Nodes

}

func GetNodeLog(nodeName string) (string, string, string, string, string, string) {

	nodeInfo := MakeReq("https://kubernetes/api/v1/nodes/"+nodeName, "GET")
	//items := gjson.Get(nodeInfo, "items.#")

	labels := gjson.Get(nodeInfo, "metadata.labels").String()
	capacity := gjson.Get(nodeInfo, "status.capacity").String()
	conditions := gjson.Get(nodeInfo, "status.conditions").String()
	info := gjson.Get(nodeInfo, "status.nodeInfo").String()
	images := gjson.Get(nodeInfo, "status.images").String()

	return labels, capacity, conditions, info, images, nodeInfo
}

func GetDeployments() []Deployment {
	var deployment Deployment
	var Deployments []Deployment

	deploymentInfo := MakeReq("https://kubernetes/apis/apps/v1/namespaces/"+os.Getenv("NAMESPACE")+"/deployments/", "GET")

	items := gjson.Get(deploymentInfo, "items.#")
	itemCount := int(items.Int())
	for i := 0; i < itemCount; i++ {
		count := strconv.Itoa(i)
		deployment.DeploymentName = gjson.Get(deploymentInfo, "items."+count+".metadata.name").String()
		deployment.ImageName = gjson.Get(deploymentInfo, "items."+count+".spec.template.spec.containers.0.image").String()
		deployment.ReplicaCount = gjson.Get(deploymentInfo, "items."+count+".spec.replicas").String()
		deployment.PullPolicy = gjson.Get(deploymentInfo, "items."+count+".spec.template.spec.containers.0.imagePullPolicy").String()

		Deployments = append(Deployments, deployment)

	}

	return Deployments

}

func GetDeploymentLog(deploymentName string) string {

	deploymentInfo := MakeReq("https://kubernetes/apis/apps/v1/namespaces/"+os.Getenv("NAMESPACE")+"/deployments/"+deploymentName, "GET")

	return deploymentInfo
}

func Scaler(kind string, name string, operation string) {
	//TODO : replica sayısını dışarıdan alıp deployment info katmanını kaldır.
	//TODO : kinddan dealıp tek bir request yap Patch("https://kubernetes/apis/apps/v1/namespaces/"+os.Getenv("NAMESPACE")+"/"+kind+"/"+name, data)
	if kind == "Deployment" {

		deploymentInfo := GetDeploymentLog(name)
		ReplicaCount := gjson.Get(deploymentInfo, "spec.replicas").Int()

		if operation == "increase" {
			NewReplicaCount := ReplicaCount + 1
			Count := strconv.FormatInt(NewReplicaCount, 10)
			data := `{"spec":{"replicas":` + Count + `}}`
			Patch("https://kubernetes/apis/apps/v1/namespaces/"+os.Getenv("NAMESPACE")+"/deployments/"+name, data)

		}

		if operation == "decrease" {
			NewReplicaCount := ReplicaCount - 1
			Count := strconv.FormatInt(NewReplicaCount, 10)
			data := `{"spec":{"replicas":` + Count + `}}`
			Patch("https://kubernetes/apis/apps/v1/namespaces/"+os.Getenv("NAMESPACE")+"/deployments/"+name, data)
		}

	}

	if kind == "StatefulSet" {

		statefulsetInfo := GetStatefulSetLog(name)
		ReplicaCount := gjson.Get(statefulsetInfo, "spec.replicas").Int()

		if operation == "increase" {
			NewReplicaCount := ReplicaCount + 1
			Count := strconv.FormatInt(NewReplicaCount, 10)
			data := `{"spec":{"replicas":` + Count + `}}`
			Patch("https://kubernetes/apis/apps/v1/namespaces/"+os.Getenv("NAMESPACE")+"/statefulsets/"+name, data)

		}

		if operation == "decrease" {
			NewReplicaCount := ReplicaCount - 1
			Count := strconv.FormatInt(NewReplicaCount, 10)
			data := `{"spec":{"replicas":` + Count + `}}`
			Patch("https://kubernetes/apis/apps/v1/namespaces/"+os.Getenv("NAMESPACE")+"/statefulsets/"+name, data)
		}
	}

}

func GetStatefulSets() []StatefulSet {
	var statefulset StatefulSet
	var StatefulSets []StatefulSet

	statefulsetInfo := MakeReq("https://kubernetes/apis/apps/v1/namespaces/"+os.Getenv("NAMESPACE")+"/statefulsets/", "GET")

	items := gjson.Get(statefulsetInfo, "items.#")
	itemCount := int(items.Int())
	for i := 0; i < itemCount; i++ {
		count := strconv.Itoa(i)
		statefulset.StatefulSetName = gjson.Get(statefulsetInfo, "items."+count+".metadata.name").String()
		statefulset.ReplicaCount = gjson.Get(statefulsetInfo, "items."+count+".spec.replicas").String()
		statefulset.PullPolicy = gjson.Get(statefulsetInfo, "items."+count+".spec.template.spec.containers.0.imagePullPolicy").String()
		statefulset.ImageName = gjson.Get(statefulsetInfo, "items."+count+".spec.template.spec.containers.0.image").String()

		StatefulSets = append(StatefulSets, statefulset)

	}

	return StatefulSets

}

func GetStatefulSetLog(statefulsetName string) string {

	statefulsetInfo := MakeReq("https://kubernetes/apis/apps/v1/namespaces/"+os.Getenv("NAMESPACE")+"/statefulsets/"+statefulsetName, "GET")

	return statefulsetInfo
}

func GetCronjobs() []Cronjob {
	var cronjob Cronjob
	var Cronjobs []Cronjob

	cronjobInfo := MakeReq("https://kubernetes/apis/batch/v1beta1/namespaces/"+os.Getenv("NAMESPACE")+"/cronjobs", "GET")

	items := gjson.Get(cronjobInfo, "items.#")
	itemCount := int(items.Int())
	for i := 0; i < itemCount; i++ {
		count := strconv.Itoa(i)
		cronjob.CronjobName = gjson.Get(cronjobInfo, "items."+count+".metadata.name").String()
		cronjob.ImageName = gjson.Get(cronjobInfo, "items."+count+".spec.jobTemplate.spec.template.spec.containers.0.image").String()
		cronjob.Scheduler = gjson.Get(cronjobInfo, "items."+count+".spec.schedule").String()
		cronjob.LastSchedule = gjson.Get(cronjobInfo, "items."+count+".status.lastScheduleTime").String()
		cronjob.History = gjson.Get(cronjobInfo, "items."+count+".spec.successfulJobsHistoryLimit").String()
		cronjob.Suspend = gjson.Get(cronjobInfo, "items."+count+".spec.suspend").Bool()
		cronjob.PullPolicy = gjson.Get(cronjobInfo, "items."+count+".spec.jobTemplate.spec.template.spec.containers.0.imagePullPolicy").String()
		Cronjobs = append(Cronjobs, cronjob)

	}

	return Cronjobs

}

func GetCronJobLog(cronjobName string) string {

	cronjobNameInfo := MakeReq("https://kubernetes/apis/batch/v1beta1/namespaces/"+os.Getenv("NAMESPACE")+"/cronjobs/"+cronjobName, "GET")

	return cronjobNameInfo
}

func CronJobStarter(cronjobName string) {
	cronjobInfo := GetCronJobLog(cronjobName)
	name := gjson.Get(cronjobInfo, "metadata.name").String() + "-" + time.Now().String()[len(time.Now().String())-8:]
	spec := gjson.Get(cronjobInfo, "spec.jobTemplate.spec.template.spec").String()
	data := `{"kind": "Job","apiVersion": "batch/v1","metadata": {"name": "` + name + `","namespace": "` + os.Getenv("NAMESPACE") + `"},"spec": {"parallelism": 1,"completions": 1,"template": {"metadata": {"name": "` + name + `"},"spec":` + spec + `}}}`
	Create("https://kubernetes/apis/batch/v1/namespaces/"+os.Getenv("NAMESPACE")+"/jobs", data)

}

func ChangeScheduler(cronjobName string, schedule string) {

	data := `{"spec":{"schedule":"` + schedule + `"}}`
	Patch("https://kubernetes/apis/batch/v1beta1/namespaces/"+os.Getenv("NAMESPACE")+"/cronjobs/"+cronjobName, data)

}

func ChangeCronJobHistory(cronjobName string, history string) {

	data := `{"spec":{"successfulJobsHistoryLimit":` + history + `}}`
	fmt.Println(data)
	Patch("https://kubernetes/apis/batch/v1beta1/namespaces/"+os.Getenv("NAMESPACE")+"/cronjobs/"+cronjobName, data)

}

func CronJobSuspender(cronjobName string, suspend string) {

	data := `{"spec":{"suspend":` + suspend + `}}`
	Patch("https://kubernetes/apis/batch/v1beta1/namespaces/"+os.Getenv("NAMESPACE")+"/cronjobs/"+cronjobName, data)

}

func GetDaemonSets() []Daemonset {
	var daemonset Daemonset
	var Daemonsets []Daemonset

	daemonsetInfo := MakeReq("https://kubernetes/apis/apps/v1/namespaces/"+os.Getenv("NAMESPACE")+"/daemonsets", "GET")

	items := gjson.Get(daemonsetInfo, "items.#")
	itemCount := int(items.Int())
	for i := 0; i < itemCount; i++ {
		count := strconv.Itoa(i)
		daemonset.DaemonsetName = gjson.Get(daemonsetInfo, "items."+count+".metadata.name").String()
		daemonset.ImageName = gjson.Get(daemonsetInfo, "items."+count+".spec.template.spec.containers.0.image").String()
		daemonset.PullPolicy = gjson.Get(daemonsetInfo, "items."+count+".spec.template.spec.containers.0.imagePullPolicy").String()
		daemonset.PodCount = gjson.Get(daemonsetInfo, "items."+count+".status.numberAvailable").String()
		Daemonsets = append(Daemonsets, daemonset)
	}

	return Daemonsets

}

func GetDaemonSetLog(daemonsetName string) string {

	daemonsetInfo := MakeReq("https://kubernetes/apis/apps/v1/namespaces/"+os.Getenv("NAMESPACE")+"/daemonsets/"+daemonsetName, "GET")

	return daemonsetInfo
}

func UpdateImage(imageName string) {

	data := `{"spec":{"template":{"spec":{"initContainers":[{"command":["docker","pull","hello-world"],"image":"docker","imagePullPolicy":"Always","name":"prepull"}]}}}}`

	Patch("https://kubernetes/apis/apps/v1/namespaces/default/daemonsets/prepull", data)

	time.Sleep(3 * time.Second)

	data = `{"spec":{"template":{"spec":{"initContainers":[{"command":["docker","pull","` + imageName + `"],"image":"docker","imagePullPolicy":"Always","name":"prepull"}]}}}}`

	Patch("https://kubernetes/apis/apps/v1/namespaces/default/daemonsets/prepull", data)

}

func GetUpdateStatus() string {

	info := MakeReq("https://kubernetes/apis/apps/v1/namespaces/default/daemonsets/prepull", "GET")
	desired := gjson.Get(info, "status.desiredNumberScheduled").String()
	available := gjson.Get(info, "status.numberAvailable").String()

	if desired == available {
		return "Updated!"
	} else {
		return "Updating.Not Updated Yet!"
	}

}

func GetJobs() []Job {
	var job Job
	var Jobs []Job

	jobInfo := MakeReq("https://kubernetes/apis/batch/v1/namespaces/"+os.Getenv("NAMESPACE")+"/jobs/", "GET")
	items := gjson.Get(jobInfo, "items.#")
	itemCount := int(items.Int())
	for i := 0; i < itemCount; i++ {
		count := strconv.Itoa(i)
		job.JobName = gjson.Get(jobInfo, "items."+count+".metadata.name").String()
		job.ImageName = gjson.Get(jobInfo, "items."+count+".spec.template.spec.containers.0.image").String()
		job.PullPolicy = gjson.Get(jobInfo, "items."+count+".spec.template.spec.containers.0.imagePullPolicy").String()
		active := gjson.Get(jobInfo, "items."+count+".status.active").Int()
		job.Active = strconv.FormatInt(active, 10)
		job.StartTime = gjson.Get(jobInfo, "items."+count+".status.startTime").String()

		Jobs = append(Jobs, job)

	}

	return Jobs

}

func DeleteJob(jobName string) {

	MakeReq("https://kubernetes/apis/batch/v1/namespaces/"+os.Getenv("NAMESPACE")+"/jobs/"+jobName, "DELETE")

}

func PullPolicyChanger(kind string, name string, policy string) {

	switch kind {

	case "deployment":
		data := GetDeploymentLog(name)

		currentpolicy := gjson.Get(data, "spec.template.spec.containers.0.imagePullPolicy").String()

		data = strings.Replace(data, currentpolicy, policy, -1)

		data = `{"spec":{"template":{"spec":{"containers":` + gjson.Get(data, "spec.template.spec.containers").String() + `}}}}`

		Patch("https://kubernetes/apis/apps/v1/namespaces/"+os.Getenv("NAMESPACE")+"/deployments/"+name, data)

	case "statefulset":

		data := GetStatefulSetLog(name)

		currentpolicy := gjson.Get(data, "spec.template.spec.containers.0.imagePullPolicy").String()

		data = strings.Replace(data, currentpolicy, policy, -1)

		data = `{"spec":{"template":{"spec":{"containers":` + gjson.Get(data, "spec.template.spec.containers").String() + `}}}}`

		Patch("https://kubernetes/apis/apps/v1/namespaces/"+os.Getenv("NAMESPACE")+"/statefulsets/"+name, data)

	case "daemonset":

		data := GetDaemonSetLog(name)

		currentpolicy := gjson.Get(data, "spec.template.spec.containers.0.imagePullPolicy").String()

		data = strings.Replace(data, currentpolicy, policy, -1)

		data = `{"spec":{"template":{"spec":{"containers":` + gjson.Get(data, "spec.template.spec.containers").String() + `}}}}`

		Patch("https://kubernetes/apis/apps/v1/namespaces/"+os.Getenv("NAMESPACE")+"/daemonsets/"+name, data)

	case "cronjob":

		data := GetCronJobLog(name)

		currentpolicy := gjson.Get(data, "spec.jobTemplate.spec.template.spec.containers.0.imagePullPolicy").String()

		data = strings.Replace(data, currentpolicy, policy, -1)

		data = `{"spec":{"jobTemplate":{"spec":{"template":{"spec":{"containers":` + gjson.Get(data, "spec.jobTemplate.spec.template.spec.containers").String() + `}}}}}}`

		Patch("https://kubernetes/apis/batch/v1beta1/namespaces/"+os.Getenv("NAMESPACE")+"/cronjobs/"+name, data)
	default:
		fmt.Println("No such a resource type")

	}

}

///////////////////////////////
func ToJSON(data []byte) ([]byte, error) {
	return yaml.YAMLToJSON(data)
}

func Yaml2Json(yaml string) string {
	j, _ := ToJSON([]byte(yaml))
	return string(j)

}

////////////////////////////////
