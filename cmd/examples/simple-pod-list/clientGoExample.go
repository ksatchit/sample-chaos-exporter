package main

import (
        "fmt"
        "log"
        "flag"
        "path/filepath"
        //"k8s.io/apimachinery/pkg/api/errors"
        metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
        "k8s.io/client-go/kubernetes"
        "k8s.io/client-go/rest"
        "k8s.io/client-go/util/homedir"
        "k8s.io/client-go/tools/clientcmd"
        _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"

)

var kubeconfig *string

func init() {
        if home := homedir.HomeDir(); home != ""  {
                kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
        } else {
               kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
        }
        flag.Parse()
}

func main() {

        var config *rest.Config
        var err error

        if *kubeconfig == "" {
               log.Printf("using the in-cluster config")
               config, err = rest.InClusterConfig()
        } else {
               log.Printf("using configuration from '%s'", *kubeconfig)
               config, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)
        }

        if err != nil {
		panic(err.Error())
	}

       // creates the clientset
       clientset, err := kubernetes.NewForConfig(config)
       if err != nil {
                panic(err.Error())
        }

       for {
               pods, err := clientset.CoreV1().Pods("").List(metav1.ListOptions{})
               if err != nil {
			panic(err.Error())
	       }
	       fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))
}
}
