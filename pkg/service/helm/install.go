package helm

import (
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
)

// Scheme is the default instance of runtime.Scheme to which types in the Kubernetes API are already registered.
var Scheme = runtime.NewScheme()

// Codecs provides access to encoding and decoding for the scheme
var Codecs = serializer.NewCodecFactory(Scheme)

// ParameterCodec handles versioning of objects that are converted to query parameters.
var ParameterCodec = runtime.NewParameterCodec(Scheme)

func init() {
	ExtensionInstall(Scheme)
	AppInstall(Scheme)
}

func ExtensionInstall(scheme *runtime.Scheme) {
	utilruntime.Must(apiextensions.AddToScheme(scheme))
	utilruntime.Must(apiextensionsv1.AddToScheme(scheme))
	utilruntime.Must(scheme.SetVersionPriority(apiextensionsv1.SchemeGroupVersion))
}

func AppInstall(scheme *runtime.Scheme) {
	utilruntime.Must(appsv1.AddToScheme(scheme))
	utilruntime.Must(scheme.SetVersionPriority(appsv1.SchemeGroupVersion))
}
