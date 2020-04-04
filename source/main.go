package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
)

var cookieHandler = securecookie.New(securecookie.GenerateRandomKey(64), securecookie.GenerateRandomKey(32))

func main() {

	if os.Getenv("NAMESPACE") == "" {
		os.Setenv("NAMESPACE", "default")
	}
	var router = mux.NewRouter()

	s := http.StripPrefix("/static/", http.FileServer(http.Dir("./ui/static/")))
	router.PathPrefix("/static/").Handler(s)

	router.HandleFunc("/", LoginHandler)
	router.HandleFunc("/Index", IndexHandler)
	router.HandleFunc("/Nodes", NodesHandler)
	router.HandleFunc("/Deployments", DeploymentsHandler)
	router.HandleFunc("/DaemonSets", DaemonSetHandler)
	router.HandleFunc("/StatefulSets", StatefulSetsHandler)
	router.HandleFunc("/CronJobs", CronJobsHandler)
	router.HandleFunc("/CronJobStart", CronJobStartHandler)
	router.HandleFunc("/CronJobHistory", CronJobHistoryHandler)
	router.HandleFunc("/Jobs", JobsHandler)
	router.HandleFunc("/Creator", CreatorHandler)

	router.HandleFunc("/WatchUpdate", UpdateWatchHandler)

	router.HandleFunc("/DeletePod", DeletePodHandler)
	router.HandleFunc("/Scale", ScaleHandler)
	router.HandleFunc("/ChangeScheduler", ChangeSchedulerHandler)
	router.HandleFunc("/ChangeSuspend", CronJobSuspendHandler)
	router.HandleFunc("/UpdateImage", UpdateImageHandler)
	router.HandleFunc("/JobDelete", JobDeleteHandler)
	router.HandleFunc("/PolicyChanger", ChangePolicyHandler)
	router.HandleFunc("/Create", ApplyHandler)

	router.HandleFunc("/GetLog", LogHandler)
	router.HandleFunc("/GetNodeLog", NodeLogHandler)
	router.HandleFunc("/GetDeploymentLog", DeploymentLogHandler)
	router.HandleFunc("/GetStatefulSetLog", StatefulSetLogHandler)
	router.HandleFunc("/GetCronJobLog", CronJobLogHandler)
	router.HandleFunc("/GetDaemonSetLog", DaemonSetLogHandler)

	router.HandleFunc("/Login", Login).Methods("POST")
	router.HandleFunc("/Logout", LogoutHandler)

	http.Handle("/", router)

	port := os.Getenv("PORT")
	if port == "" {
		port = "4000"
	}

	fmt.Println("Serving from: ", port)
	http.ListenAndServe(":"+port, nil)

}
