package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/kubernetes"

	"github.com/gorilla/mux"
	"github.com/TimWoolford/straw-pod/internal/status"
)

type server struct {
	status    status.Status
	clientSet *kubernetes.Clientset
}

func Start() {
	k8sConfig, err := rest.InClusterConfig()

	if err != nil {
		panic(err.Error())
	}

	clientSet, err := kubernetes.NewForConfig(k8sConfig)

	if err != nil {
		panic(err.Error())
	}

	server := server{
		status:    status.MakeMe(),
		clientSet: clientSet,
	}

	r := mux.NewRouter()

	r.HandleFunc("/status", server.StatusHandler)
	r.HandleFunc("/setStatus/{status}", server.SetStatusHandler)
	r.HandleFunc("/pods", server.Pods).Methods("GET")
	r.HandleFunc("/pods", server.UpdatePods).Methods("PUT")

	r.PathPrefix("/").Handler( http.FileServer(http.Dir("/html")))

	http.Handle("/", r)
	http.ListenAndServe(fmt.Sprintf(":%d", server.status.Port()), nil)
}

func (s *server) StatusHandler(w http.ResponseWriter, r *http.Request) {
	w.Header()["Content-Type"] = []string{"text/json"}

	bytes, _ := json.Marshal(s.status)

	w.Write(bytes)
}

func (s *server) SetStatusHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	s.status.OverallStatus = vars["status"]
}

func (s *server) Pods(w http.ResponseWriter, r *http.Request) {
	pods, err := s.clientSet.CoreV1().Pods(s.status.Namespace()).List(metaV1.ListOptions{})
	if err != nil {
		panic(err)
	}

	podStatus := make(map[string]string)

	for _, pod := range pods.Items {
		if pod.Labels["app_name"] == "straw-pod" {
			resp, getErr := http.Get(fmt.Sprintf("http://%s:%d/status", pod.Status.PodIP, s.status.Port()))
			if getErr != nil {
				panic(getErr)
			}
			theStatus := &status.Status{}
			closer, _ := ioutil.ReadAll(resp.Body)
			json.Unmarshal(closer, theStatus)
			podStatus[pod.Name] = theStatus.OverallStatus
		}
	}

	w.Header()["Content-Type"] = []string{"text/json"}
	json.NewEncoder(w).Encode(podStatus)
}

func (s *server) UpdatePods(w http.ResponseWriter, r *http.Request) {
	respMap := make(map[string]string)
	responseContent, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(responseContent, respMap)

	for pod, newStatus := range respMap {
		if pod == s.status.Hostname {
			s.status.OverallStatus = newStatus
		} else {
			thePod, _ := s.clientSet.CoreV1().Pods(s.status.Namespace()).Get(pod, metaV1.GetOptions{})
			request, _ := http.NewRequest("PUT", fmt.Sprintf("http://%s:%d/setStatus/%s", s.status.Port(), thePod.Status.PodIP, newStatus), nil)

			client := &http.Client{}

			client.Do(request)
		}
	}
}
