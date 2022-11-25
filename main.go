package main

import (
	"context"
	"flag"
	"fmt"

	autoscalev1 "k8s.io/api/autoscaling/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	// 변수 초기화
	ctx := context.Background()
	ns := flag.String("ns", "", "namespace")
	dep_name := flag.String("dep_name", "", "deployment")
	min_rep_64 := flag.Int("min_rep", 0, "min_rep")
	max_rep_64 := flag.Int("max_rep", 0, "max_rep")
	max_cpu_set_64 := flag.Int("max_cpu_set", 0, "max_cpu_set")
	flag.Parse()
	var min_rep int32 = int32(*min_rep_64)
	var max_rep int32 = int32(*max_rep_64)
	var max_cpu_set int32 = int32(*max_cpu_set_64)
	hpa_name := fmt.Sprintf("%s-%s", *ns, *dep_name)

	if flag.NFlag() < 5 { // 명령줄 옵션의 개수가 0개이면
		flag.Usage()
		return
	}

	// 컨피그 파일 불러옴 ( ~/.kube/config 는 안되서 다른 걸로 함. )
	kubeconfig := flag.String("kubeconfig", "/root/.kube/config", "location to your kubeconfig file")
	// clientcmd -> 쿠버네티스 접근을 위한 접속 정보
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	// handle error
	if err != nil {
		// 클러스터 내부에서 접근해서 데이터를 가져오는 방법
		fmt.Printf("error %s", err)
		config, err = rest.InClusterConfig()
		if err != nil {
			fmt.Printf("InClusterConfig error %s", err)
		}
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Printf("error %s", err)
	}

	_, err = clientset.AutoscalingV1().HorizontalPodAutoscalers(*ns).Get(ctx, hpa_name, metav1.GetOptions{})
	hpa_s := autoscalev1.HorizontalPodAutoscaler{
		ObjectMeta: metav1.ObjectMeta{
			Name: hpa_name,
		},
		Spec: autoscalev1.HorizontalPodAutoscalerSpec{

			ScaleTargetRef: autoscalev1.CrossVersionObjectReference{
				APIVersion: "apps/v1",
				Kind:       "Deployment",
				Name:       *dep_name,
			},
			MinReplicas:                    &min_rep,
			MaxReplicas:                    max_rep,
			TargetCPUUtilizationPercentage: &max_cpu_set,
		},
	}
	// hpa create time
	if apierrors.IsNotFound(err) {
		fmt.Printf("hpa is Not Found %s", err.Error())
		fmt.Printf("hpa create")
		hpa2, err := clientset.AutoscalingV1().HorizontalPodAutoscalers(*ns).Create(ctx, &hpa_s, metav1.CreateOptions{})
		fmt.Print(hpa2)
		if err != nil {
			fmt.Printf("Create error %s", err)
		}
	} else {
		//hpa update
		fmt.Printf("hpa is Update")
		hpaup, err := clientset.AutoscalingV1().HorizontalPodAutoscalers(*ns).Update(ctx, &hpa_s, metav1.UpdateOptions{})
		fmt.Print(hpaup)
		if err != nil {
			fmt.Printf("Create error %s", err)
		}
	}
}
