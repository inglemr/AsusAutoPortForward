package internal

import (
	autoportforward "autoportforward/pkg"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/davecgh/go-spew/spew"
	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	toolsWatch "k8s.io/client-go/tools/watch"
)

func NewService(config Config) *Service {
	routerClient := autoportforward.NewAsusRouterClient(config.RouterAddress, config.Username, config.Password)
	kConfig, _ := clientcmd.BuildConfigFromFlags("", os.Getenv("KUBECONFIG"))
	clientset, _ := kubernetes.NewForConfig(kConfig)
	return &Service{
		rc:        routerClient,
		mux:       http.NewServeMux(),
		kConfig:   kConfig,
		kClient:   clientset,
		svcConfig: config,
	}
}

type Service struct {
	rc        *autoportforward.AsusRouterClient
	mux       *http.ServeMux
	kConfig   *rest.Config
	kClient   *kubernetes.Clientset
	svcConfig Config
}

func (s *Service) StartService() {
	go s.watchServices()
	s.mux.HandleFunc("/getforwardrules", s.HandleGetPortForwardRules)
	http.ListenAndServe(":8080", s.mux)
}

func (s *Service) HandleGetPortForwardRules(w http.ResponseWriter, r *http.Request) {
	ports := s.rc.GetPortForwardRules()
	output, _ := json.Marshal(ports)
	w.Header().Set("Content-Type", "application/json")
	w.Write(output)
	return
}

func (s *Service) watchServices() {
	log.Info("Starting to watch services")
	watchFunc := func(options metav1.ListOptions) (watch.Interface, error) {
		timeOut := int64(60)
		labelSelector := fmt.Sprintf("%s=%s", "service.kubernetes.io/autoportforward", "true")
		return s.kClient.CoreV1().Services(v1.NamespaceAll).Watch(context.Background(), metav1.ListOptions{TimeoutSeconds: &timeOut, LabelSelector: labelSelector})
	}

	watcher, _ := toolsWatch.NewRetryWatcher("1", &cache.ListWatch{WatchFunc: watchFunc})
	for event := range watcher.ResultChan() {
		item := event.Object.(*corev1.Service)
		log.Debugf("Event: %v", event.Type)
		log.Debugf("Item: %v", item.Name)
		if item.Spec.Type == "NodePort" {
			switch event.Type {
			case watch.Modified:
				rules := autoportforward.NewPortForwardRulesFromK8sService(*item, s.svcConfig.DefaultTargetAddress)
				spew.Dump(rules)
			case watch.Deleted:
				fmt.Println(item.Name)
			case watch.Added:
				fmt.Println(item.Name)
			}
		}
	}
}

// service.kubernetes.io/autoportforward: "true"
