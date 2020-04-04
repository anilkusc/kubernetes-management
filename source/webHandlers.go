package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/tidwall/gjson"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {

	userName := getUserName(r)
	if userName == os.Getenv("USERNAME") {

		fmt.Println("IndexHandler")

		files := []string{
			"./ui/html/index.html",
			"./ui/html/base.html",
		}
		ts, err := template.ParseFiles(files...)
		if err != nil {
			fmt.Println(err)
			return
		}
		data := GetPods()
		err = ts.Execute(w, data)
		if err != nil {
			fmt.Println(err.Error())
		}
	} else {
		http.Redirect(w, r, "/", 302)
	}

}

func DeletePodHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("DeletePodHandler")

	userName := getUserName(r)
	if userName == os.Getenv("USERNAME") {

		podName := r.URL.Query().Get("podName")
		DeletePod(podName)
		http.Redirect(w, r, "/Index", 302)
	} else {
		http.Redirect(w, r, "/", 302)
	}

}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("LoginHandler")

	files := []string{
		"./ui/html/login.html",
	}

	ts, err := template.ParseFiles(files...)

	if err != nil {
		fmt.Println(err)
		return
	}
	err = ts.Execute(w, nil)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func LogHandler(w http.ResponseWriter, r *http.Request) {

	userName := getUserName(r)
	if userName == os.Getenv("USERNAME") {

		fmt.Println("LogHandler")

		files := []string{
			"./ui/html/log.html",
			"./ui/html/base.html",
		}
		ts, err := template.ParseFiles(files...)
		if err != nil {
			fmt.Println(err)
			return
		}

		type Data struct {
			Log    string
			Status string
		}
		podName := r.URL.Query().Get("podName")
		log, status := GetLog(podName)
		log = strings.Replace(log, "\n", "<br>", -1)
		data := &Data{log, status}

		err = ts.Execute(w, data)
		if err != nil {
			fmt.Println(err.Error())
		}
	} else {
		http.Redirect(w, r, "/", 302)
	}

}

func NodesHandler(w http.ResponseWriter, r *http.Request) {

	userName := getUserName(r)
	if userName == os.Getenv("USERNAME") {

		fmt.Println("NodesHandler")

		files := []string{
			"./ui/html/nodes.html",
			"./ui/html/base.html",
		}
		ts, err := template.ParseFiles(files...)
		if err != nil {
			fmt.Println(err)
			return
		}

		data := GetNodes()

		err = ts.Execute(w, data)
		if err != nil {
			fmt.Println(err.Error())
		}
	} else {
		http.Redirect(w, r, "/", 302)
	}

}

func NodeLogHandler(w http.ResponseWriter, r *http.Request) {

	userName := getUserName(r)
	if userName == os.Getenv("USERNAME") {

		fmt.Println("NodeLogHandler")

		files := []string{
			"./ui/html/nodelog.html",
			"./ui/html/base.html",
		}
		ts, err := template.ParseFiles(files...)
		if err != nil {
			fmt.Println(err)
			return
		}
		type Data struct {
			Label      string
			Capacity   string
			Conditions string
			Info       string
			Images     string
			NodeInfo   string
		}
		nodeName := r.URL.Query().Get("nodeName")
		labels, capacity, conditions, info, images, nodeInfo := GetNodeLog(nodeName)
		data := &Data{labels, capacity, conditions, info, images, nodeInfo}

		err = ts.Execute(w, data)
		if err != nil {
			fmt.Println(err.Error())
		}
	} else {
		http.Redirect(w, r, "/", 302)
	}

}

func DeploymentsHandler(w http.ResponseWriter, r *http.Request) {

	userName := getUserName(r)
	if userName == os.Getenv("USERNAME") {

		fmt.Println("DeploymentsHandler")

		files := []string{
			"./ui/html/deployments.html",
			"./ui/html/base.html",
		}
		ts, err := template.ParseFiles(files...)
		if err != nil {
			fmt.Println(err)
			return
		}
		data := GetDeployments()
		err = ts.Execute(w, data)
		if err != nil {
			fmt.Println(err.Error())
		}
	} else {
		http.Redirect(w, r, "/", 302)
	}

}

func DeploymentLogHandler(w http.ResponseWriter, r *http.Request) {

	userName := getUserName(r)
	if userName == os.Getenv("USERNAME") {

		fmt.Println("DeploymentLogHandler")

		files := []string{
			"./ui/html/deploymentlog.html",
			"./ui/html/base.html",
		}
		ts, err := template.ParseFiles(files...)
		if err != nil {
			fmt.Println(err)
			return
		}
		type Data struct {
			Info string
		}
		deploymentName := r.URL.Query().Get("deploymentName")
		info := GetDeploymentLog(deploymentName)
		data := &Data{info}

		err = ts.Execute(w, data)
		if err != nil {
			fmt.Println(err.Error())
		}
	} else {
		http.Redirect(w, r, "/", 302)
	}

}

func ScaleHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("ScaleHandler")

	userName := getUserName(r)
	if userName == os.Getenv("USERNAME") {

		Kind := r.URL.Query().Get("Kind")
		Name := r.URL.Query().Get("Name")
		Operation := r.URL.Query().Get("Operation")

		Scaler(Kind, Name, Operation)

		switch Kind {

		case "Deployment":
			http.Redirect(w, r, "/Deployments", 302)
		case "StatefulSet":
			http.Redirect(w, r, "/StatefulSets", 302)
		default:
			fmt.Println("No such a resource type")
		}

	} else {
		http.Redirect(w, r, "/", 302)
	}

}

func StatefulSetsHandler(w http.ResponseWriter, r *http.Request) {

	userName := getUserName(r)
	if userName == os.Getenv("USERNAME") {

		fmt.Println("StatefulSetsHandler")

		files := []string{
			"./ui/html/statefulsets.html",
			"./ui/html/base.html",
		}
		ts, err := template.ParseFiles(files...)
		if err != nil {
			fmt.Println(err)
			return
		}
		data := GetStatefulSets()
		err = ts.Execute(w, data)
		if err != nil {
			fmt.Println(err.Error())
		}
	} else {
		http.Redirect(w, r, "/", 302)
	}

}

func StatefulSetLogHandler(w http.ResponseWriter, r *http.Request) {

	userName := getUserName(r)
	if userName == os.Getenv("USERNAME") {

		fmt.Println("StatefulSetLogHandler")

		files := []string{
			"./ui/html/statefulsetlog.html",
			"./ui/html/base.html",
		}
		ts, err := template.ParseFiles(files...)
		if err != nil {
			fmt.Println(err)
			return
		}
		type Data struct {
			Info string
		}
		statefulsetName := r.URL.Query().Get("statefulsetName")
		info := GetStatefulSetLog(statefulsetName)
		data := &Data{info}

		err = ts.Execute(w, data)
		if err != nil {
			fmt.Println(err.Error())
		}
	} else {
		http.Redirect(w, r, "/", 302)
	}

}

func CronJobsHandler(w http.ResponseWriter, r *http.Request) {

	userName := getUserName(r)
	if userName == os.Getenv("USERNAME") {

		fmt.Println("CronJobHandler")

		files := []string{
			"./ui/html/cronjob.html",
			"./ui/html/base.html",
		}
		ts, err := template.ParseFiles(files...)
		if err != nil {
			fmt.Println(err)
			return
		}
		data := GetCronjobs()
		err = ts.Execute(w, data)
		if err != nil {
			fmt.Println(err.Error())
		}
	} else {
		http.Redirect(w, r, "/", 302)
	}

}

func CronJobLogHandler(w http.ResponseWriter, r *http.Request) {

	userName := getUserName(r)
	if userName == os.Getenv("USERNAME") {

		fmt.Println("CronJobLogHandler")

		files := []string{
			"./ui/html/cronjoblog.html",
			"./ui/html/base.html",
		}
		ts, err := template.ParseFiles(files...)
		if err != nil {
			fmt.Println(err)
			return
		}
		type Data struct {
			Info string
		}
		cronjobName := r.URL.Query().Get("cronjobName")
		info := GetCronJobLog(cronjobName)
		data := &Data{info}

		err = ts.Execute(w, data)
		if err != nil {
			fmt.Println(err.Error())
		}
	} else {
		http.Redirect(w, r, "/", 302)
	}

}

func CronJobStartHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("CronJobStartHandler")

	userName := getUserName(r)
	if userName == os.Getenv("USERNAME") {

		cronjobName := r.URL.Query().Get("cronjobName")

		CronJobStarter(cronjobName)

		http.Redirect(w, r, "/CronJobs", 302)
	} else {
		http.Redirect(w, r, "/", 302)
	}

}

func ChangeSchedulerHandler(w http.ResponseWriter, r *http.Request) {

	userName := getUserName(r)
	if userName == os.Getenv("USERNAME") {
		fmt.Println("ChangeSchedulerHandler")

		cronjobName := r.URL.Query().Get("name")
		schedule := r.URL.Query().Get("schedule")
		ChangeScheduler(cronjobName, schedule)

		http.Redirect(w, r, "/CronJobs", 302)
	} else {
		http.Redirect(w, r, "/", 302)
	}

}

func CronJobHistoryHandler(w http.ResponseWriter, r *http.Request) {

	userName := getUserName(r)
	if userName == os.Getenv("USERNAME") {
		fmt.Println("CronJobHistoryHandler")

		cronjobName := r.URL.Query().Get("name")
		history := r.URL.Query().Get("history")
		ChangeCronJobHistory(cronjobName, history)

		http.Redirect(w, r, "/CronJobs", 302)
	} else {
		http.Redirect(w, r, "/", 302)
	}

}

func CronJobSuspendHandler(w http.ResponseWriter, r *http.Request) {

	userName := getUserName(r)
	if userName == os.Getenv("USERNAME") {
		fmt.Println("CronJobSuspendHandler")

		cronjobName := r.URL.Query().Get("cronjobName")
		suspend := r.URL.Query().Get("suspend")
		CronJobSuspender(cronjobName, suspend)

		http.Redirect(w, r, "/CronJobs", 302)
	} else {
		http.Redirect(w, r, "/", 302)
	}

}

func DaemonSetHandler(w http.ResponseWriter, r *http.Request) {

	userName := getUserName(r)
	if userName == os.Getenv("USERNAME") {

		fmt.Println("GetDaemonSetHandler")

		files := []string{
			"./ui/html/daemonsets.html",
			"./ui/html/base.html",
		}
		ts, err := template.ParseFiles(files...)
		if err != nil {
			fmt.Println(err)
			return
		}
		data := GetDaemonSets()
		err = ts.Execute(w, data)
		if err != nil {
			fmt.Println(err.Error())
		}
	} else {
		http.Redirect(w, r, "/", 302)
	}

}

func DaemonSetLogHandler(w http.ResponseWriter, r *http.Request) {

	userName := getUserName(r)
	if userName == os.Getenv("USERNAME") {

		fmt.Println("DaemonSetLogHandler")

		files := []string{
			"./ui/html/daemonsetlog.html",
			"./ui/html/base.html",
		}
		ts, err := template.ParseFiles(files...)
		if err != nil {
			fmt.Println(err)
			return
		}
		type Data struct {
			Info string
		}
		daemonsetName := r.URL.Query().Get("daemonsetName")
		info := GetDaemonSetLog(daemonsetName)
		data := &Data{info}

		err = ts.Execute(w, data)
		if err != nil {
			fmt.Println(err.Error())
		}
	} else {
		http.Redirect(w, r, "/", 302)
	}

}

func UpdateImageHandler(w http.ResponseWriter, r *http.Request) {

	userName := getUserName(r)
	if userName == os.Getenv("USERNAME") {
		fmt.Println("UpdateImageHandler")

		imageName := r.URL.Query().Get("imageName")
		UpdateImage(imageName)

		http.Redirect(w, r, "/WatchUpdate", 302)
	} else {
		http.Redirect(w, r, "/", 302)
	}

}

func UpdateWatchHandler(w http.ResponseWriter, r *http.Request) {

	userName := getUserName(r)
	if userName == os.Getenv("USERNAME") {

		fmt.Println("UpdateWatchHandler")

		files := []string{
			"./ui/html/updating.html",
			"./ui/html/base.html",
		}
		ts, err := template.ParseFiles(files...)
		if err != nil {
			fmt.Println(err)
			return
		}
		type Data struct {
			Info string
		}
		info := GetUpdateStatus()
		data := &Data{info}
		err = ts.Execute(w, data)
		if err != nil {
			fmt.Println(err.Error())
		}
	} else {
		http.Redirect(w, r, "/", 302)
	}

}

func JobsHandler(w http.ResponseWriter, r *http.Request) {

	userName := getUserName(r)
	if userName == os.Getenv("USERNAME") {

		fmt.Println("JobsHandler")

		files := []string{
			"./ui/html/job.html",
			"./ui/html/base.html",
		}
		ts, err := template.ParseFiles(files...)
		if err != nil {
			fmt.Println(err)
			return
		}
		data := GetJobs()
		err = ts.Execute(w, data)
		if err != nil {
			fmt.Println(err.Error())
		}
	} else {
		http.Redirect(w, r, "/", 302)
	}

}

func JobDeleteHandler(w http.ResponseWriter, r *http.Request) {

	userName := getUserName(r)
	if userName == os.Getenv("USERNAME") {

		fmt.Println("JobDeleteHandler")
		jobName := r.URL.Query().Get("jobName")
		DeleteJob(jobName)
		time.Sleep(1 * time.Second)
		http.Redirect(w, r, "/Jobs", 302)

	} else {
		http.Redirect(w, r, "/", 302)
	}

}

func ChangePolicyHandler(w http.ResponseWriter, r *http.Request) {

	userName := getUserName(r)
	if userName == os.Getenv("USERNAME") {

		kind := r.URL.Query().Get("kind")
		name := r.URL.Query().Get("name")
		policy := r.URL.Query().Get("policy")

		PullPolicyChanger(kind, name, policy)

		switch kind {

		case "deployment":
			http.Redirect(w, r, "/Deployments", 302)
		case "statefulset":
			http.Redirect(w, r, "/StatefulSets", 302)
		case "cronjob":
			http.Redirect(w, r, "/CronJobs", 302)
		case "daemonset":
			http.Redirect(w, r, "/DaemonSets", 302)
		default:
			fmt.Println("No such a resource type")
		}

	} else {
		http.Redirect(w, r, "/", 302)
	}

}

func CreatorHandler(w http.ResponseWriter, r *http.Request) {

	userName := getUserName(r)
	if userName == os.Getenv("USERNAME") {

		fmt.Println("CreatorHandler")

		files := []string{
			"./ui/html/creator.html",
			"./ui/html/base.html",
		}
		ts, err := template.ParseFiles(files...)
		if err != nil {
			fmt.Println(err)
			return
		}
		err = ts.Execute(w, nil)
		if err != nil {
			fmt.Println(err.Error())
		}
	} else {
		http.Redirect(w, r, "/", 302)
	}

}

func ApplyHandler(w http.ResponseWriter, r *http.Request) {

	userName := getUserName(r)
	if userName == os.Getenv("USERNAME") {

		fmt.Println("ApplyHandler")
		yaml := r.URL.Query().Get("yaml")
		data := Yaml2Json(yaml)
		kind := gjson.Get(data, "kind").String()
		// TODO if already exist use Patch instead of Create
		switch kind {
		case "Deployment":
			Create("https://kubernetes/api/v1/namespaces/"+os.Getenv("NAMESPACE")+"/deployments", data)
		case "StatefulSet":
			Create("https://kubernetes/api/v1/namespaces/"+os.Getenv("NAMESPACE")+"/statefulsets", data)
		case "Service":
			Create("https://kubernetes/api/v1/namespaces/"+os.Getenv("NAMESPACE")+"/services", data)
		case "ConfigMap":
			Create("https://kubernetes/api/v1/namespaces/"+os.Getenv("NAMESPACE")+"/configmaps", data)
		case "Cronjob":
			Create("https://kubernetes/apis/batch/v1beta1/namespaces/"+os.Getenv("NAMESPACE")+"/cronjobs", data)
		}

		http.Redirect(w, r, "/Index", 302)
	} else {
		http.Redirect(w, r, "/", 302)
	}

}

///////////////////////////////////////////////////////////////////////////////////////////////////////

func Login(response http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	name := request.FormValue("username")
	pass := request.FormValue("password")
	redirectTarget := "/"
	if name == os.Getenv("USERNAME") && pass == os.Getenv("PASSWORD") {
		// .. check credentials ..
		setSession(name, response)
		redirectTarget = "/Index"
	}
	http.Redirect(response, request, redirectTarget, 302)
}

func LogoutHandler(response http.ResponseWriter, request *http.Request) {
	clearSession(response)
	http.Redirect(response, request, "/", 302)
}

func setSession(userName string, response http.ResponseWriter) {
	value := map[string]string{
		"name": userName,
	}
	if encoded, err := cookieHandler.Encode("session", value); err == nil {
		cookie := &http.Cookie{
			Name:  "session",
			Value: encoded,
			Path:  "/",
		}
		http.SetCookie(response, cookie)
	}
}

func getUserName(request *http.Request) (userName string) {
	if cookie, err := request.Cookie("session"); err == nil {
		cookieValue := make(map[string]string)
		if err = cookieHandler.Decode("session", cookie.Value, &cookieValue); err == nil {
			userName = cookieValue["name"]
		}
	}
	return userName
}

func clearSession(response http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(response, cookie)
}
