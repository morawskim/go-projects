package util

import (
	"bytes"
	"context"
	_ "embed"
	"fmt"
	"html/template"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"time"
)

//go:embed template.html
var templateString string

func GenerateIndexFile(data []TraefikItem) (string, error) {
	tmpl, err := template.New("index").Parse(templateString)
	if err != nil {
		return "", fmt.Errorf("error parsing template: %v", err)
	}

	buffer := bytes.NewBufferString("")
	err = tmpl.Execute(buffer, struct {
		Items []TraefikItem
	}{Items: data})

	if err != nil {
		return "", fmt.Errorf("error executing template: %v", err)
	}

	return buffer.String(), nil
}

func CreateConfigMap(k8sClient client.Client, namespace, name string, indexHtml string) error {
	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Data: map[string]string{
			"index.html": indexHtml,
		},
	}
	ctx := context.Background()

	existingConfigMap := &corev1.ConfigMap{}
	err := k8sClient.Get(ctx, client.ObjectKey{Namespace: namespace, Name: name}, existingConfigMap)
	if err != nil && client.IgnoreNotFound(err) != nil {
		// An error occurred, log it
		return fmt.Errorf("failed to get ConfigMap: %w", err)
	}

	if err != nil && client.IgnoreNotFound(err) == nil {
		// If ConfigMap not found, create it
		err = k8sClient.Create(ctx, configMap)
		if err != nil {
			return fmt.Errorf("failed to create ConfigMap: %w", err)
		}
	} else {
		// If ConfigMap exists, update it
		existingConfigMap.Data = configMap.Data
		err = k8sClient.Update(ctx, existingConfigMap)
		if err != nil {
			return fmt.Errorf("failed to update ConfigMap: %w", err)
		}
	}

	return nil
}

func RestartDeployment(k8sClient client.Client, namespace, name string) error {
	deployment := appsv1.Deployment{}
	ctx := context.Background()

	err := k8sClient.Get(ctx, client.ObjectKey{Namespace: namespace, Name: name}, &deployment)
	if err != nil {
		return fmt.Errorf("failed to fetch deployment: %w", err)
	}

	if deployment.Spec.Template.ObjectMeta.Annotations == nil {
		deployment.Spec.Template.ObjectMeta.Annotations = make(map[string]string)
	}
	deployment.Spec.Template.ObjectMeta.Annotations["kubectl.kubernetes.io/restartedAt"] = time.Now().Format("2006-01-02T15:04:05-07:00")
	err = k8sClient.Update(ctx, &deployment)
	if err != nil {
		return fmt.Errorf("failed to update deployment: %w", err)
	}

	return nil
}
