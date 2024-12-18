package controller

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/scheme"
)

// +k8s:deepcopy-gen=true

// IngressRouteSpec defines the desired state of IngressRoute.
type IngressRouteSpec struct {
	// Routes defines the list of routes.
	Routes []Route `json:"routes"`
	// EntryPoints defines the list of entry point names to bind to.
	// Entry points have to be configured in the static configuration.
	// More info: https://doc.traefik.io/traefik/v3.3/routing/entrypoints/
	// Default: all.
	EntryPoints []string `json:"entryPoints,omitempty"`
	// TLS defines the TLS configuration.
	// More info: https://doc.traefik.io/traefik/v3.3/routing/routers/#tls
	TLS *TLS `json:"tls,omitempty"`
}

// Route holds the HTTP route configuration.
type Route struct {
	// Match defines the router's rule.
	// More info: https://doc.traefik.io/traefik/v3.3/routing/routers/#rule
	Match string `json:"match"`
	// Kind defines the kind of the route.
	// Rule is the only supported kind.
	// If not defined, defaults to Rule.
	// +kubebuilder:validation:Enum=Rule
	Kind string `json:"kind,omitempty"`
	// Priority defines the router's priority.
	// More info: https://doc.traefik.io/traefik/v3.3/routing/routers/#priority
	Priority int `json:"priority,omitempty"`
	// Syntax defines the router's rule syntax.
	// More info: https://doc.traefik.io/traefik/v3.3/routing/routers/#rulesyntax
	Syntax string `json:"syntax,omitempty"`
}

// TLS holds the TLS configuration.
// More info: https://doc.traefik.io/traefik/v3.3/routing/routers/#tls
type TLS struct {
	// SecretName is the name of the referenced Kubernetes Secret to specify the certificate details.
	SecretName string `json:"secretName,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// IngressRoute is the CRD implementation of a Traefik HTTP Router.
type IngressRoute struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	metav1.ObjectMeta `json:"metadata"`

	Spec IngressRouteSpec `json:"spec"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// IngressRouteList is a collection of IngressRoute.
type IngressRouteList struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	metav1.ListMeta `json:"metadata"`

	// Items is the list of IngressRoute.
	Items []IngressRoute `json:"items"`
}

// Register these types with the scheme
func AddToScheme(scheme2 *runtime.Scheme) error {
	GroupVersion := schema.GroupVersion{Group: "traefik.io", Version: "v1alpha1"}
	// SchemeBuilder is used to add go types to the GroupVersionKind scheme
	SchemeBuilder := &scheme.Builder{GroupVersion: GroupVersion}
	SchemeBuilder.Register(&IngressRoute{}, &IngressRouteList{})

	return SchemeBuilder.AddToScheme(scheme2)
}
