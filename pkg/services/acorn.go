package services

import (
	v1 "github.com/acorn-io/acorn/pkg/apis/internal.acorn.io/v1"
	"github.com/acorn-io/acorn/pkg/labels"
	ports2 "github.com/acorn-io/acorn/pkg/ports"
	"github.com/acorn-io/baaah/pkg/typed"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kclient "sigs.k8s.io/controller-runtime/pkg/client"
)

func forDefined(appInstance *v1.AppInstance) (result []kclient.Object) {
	for _, entry := range typed.Sorted(appInstance.Status.AppSpec.Services) {
		serviceName, service := entry.Key, entry.Value

		if ports2.IsLinked(appInstance, serviceName) {
			continue
		}

		result = append(result, &v1.ServiceInstance{
			ObjectMeta: metav1.ObjectMeta{
				Name:      serviceName,
				Namespace: appInstance.Status.Namespace,
				Labels:    labels.Managed(appInstance, labels.AcornServiceName, serviceName),
			},
			Spec: v1.ServiceInstanceSpec{
				Labels: labels.Merge(labels.Managed(appInstance, labels.AcornServiceName, serviceName),
					labels.GatherScoped(serviceName, v1.LabelTypeService,
						appInstance.Status.AppSpec.Labels, service.Labels, appInstance.Spec.Labels)),
				Annotations: labels.GatherScoped(serviceName, v1.LabelTypeService,
					appInstance.Status.AppSpec.Annotations, service.Annotations, appInstance.Spec.Annotations),
				Default:    service.Default,
				External:   service.External,
				Address:    service.Address,
				Ports:      service.Ports,
				Container:  service.Container,
				Secrets:    service.Secrets,
				Attributes: service.Attributes,
				Destroy:    service.Destroy,
			},
		})
	}

	return
}

func forRouters(appInstance *v1.AppInstance) (result []kclient.Object) {
	for _, entry := range typed.Sorted(appInstance.Status.AppSpec.Routers) {
		routerName, router := entry.Key, entry.Value

		if ports2.IsLinked(appInstance, routerName) {
			continue
		}

		result = append(result, &v1.ServiceInstance{
			ObjectMeta: metav1.ObjectMeta{
				Name:      routerName,
				Namespace: appInstance.Status.Namespace,
				Labels:    labels.Managed(appInstance, labels.AcornRouterName, routerName),
			},
			Spec: v1.ServiceInstanceSpec{
				Labels: labels.Merge(labels.Managed(appInstance, labels.AcornRouterName, routerName),
					labels.GatherScoped(routerName, v1.LabelTypeRouter,
						appInstance.Status.AppSpec.Labels, router.Labels, appInstance.Spec.Labels)),
				Annotations: labels.GatherScoped(routerName, v1.LabelTypeRouter,
					appInstance.Status.AppSpec.Annotations, router.Annotations, appInstance.Spec.Annotations),
				Ports: []v1.PortDef{
					ports2.RouterPortDef,
				},
				ContainerLabels: labels.Managed(appInstance, labels.AcornRouterName, routerName),
			},
		})
	}

	return
}

func forContainers(appInstance *v1.AppInstance) (result []kclient.Object) {
	for _, entry := range typed.Sorted(appInstance.Status.AppSpec.Containers) {
		containerName, container := entry.Key, entry.Value

		if ports2.IsLinked(appInstance, containerName) {
			continue
		}

		ports := ports2.CollectContainerPorts(&container)
		if len(ports) == 0 {
			continue
		}

		result = append(result, &v1.ServiceInstance{
			ObjectMeta: metav1.ObjectMeta{
				Name:      containerName,
				Namespace: appInstance.Status.Namespace,
				Labels:    labels.Managed(appInstance, labels.AcornContainerName, containerName),
			},
			Spec: v1.ServiceInstanceSpec{
				Labels: labels.Merge(labels.Managed(appInstance, labels.AcornContainerName, containerName),
					labels.GatherScoped(containerName, v1.LabelTypeContainer,
						appInstance.Status.AppSpec.Labels, container.Labels, appInstance.Spec.Labels)),
				Annotations: labels.GatherScoped(containerName, v1.LabelTypeContainer,
					appInstance.Status.AppSpec.Annotations, container.Annotations, appInstance.Spec.Annotations),
				Ports:     ports,
				Container: containerName,
			},
		})
	}

	return
}

func forLinkedServices(app *v1.AppInstance) (result []kclient.Object) {
	for _, link := range app.Spec.Links {
		newService := &v1.ServiceInstance{
			ObjectMeta: metav1.ObjectMeta{
				Name:      link.Target,
				Namespace: app.Status.Namespace,
				Labels: labels.Managed(app,
					labels.AcornLinkName, link.Service),
			},
			Spec: v1.ServiceInstanceSpec{
				External: link.Service,
			},
		}
		result = append(result, newService)
	}

	return
}

func ToAcornServices(appInstance *v1.AppInstance) (result []kclient.Object) {
	result = append(result, forDefined(appInstance)...)
	result = append(result, forContainers(appInstance)...)
	result = append(result, forRouters(appInstance)...)
	// determine default before adding linked services
	if len(result) == 1 {
		result[0].(*v1.ServiceInstance).Spec.Default = true
	}
	result = append(result, forLinkedServices(appInstance)...)
	return
}