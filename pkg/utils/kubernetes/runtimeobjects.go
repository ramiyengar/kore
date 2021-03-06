package kubernetes

import (
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
)

// UpdateSchema allows any schemas to be added at run time
type UpdateSchema func(*runtime.Scheme) error

// ParseK8sYaml creates multiple runtime.Objects
// - from []byte with yaml defining kubernetes manifests
// - yaml may contain "---" seperators and multiple manifest definitions
// - provide a schema function to add any required schemas at run time
func ParseK8sYaml(fileR []byte, fnUS UpdateSchema) ([]runtime.Object, error) {
	fileAsString := string(fileR[:])
	sepYamlfiles := strings.Split(fileAsString, "---")
	runtimeObjects := make([]runtime.Object, 0, len(sepYamlfiles))
	for _, f := range sepYamlfiles {
		if f == "\n" || f == "" {
			// ignore empty cases
			continue
		}
		// Ensure we know about all types first
		if err := fnUS(scheme.Scheme); err != nil {

			return nil, fmt.Errorf("error loading schemes for decoding - %s", err)
		}
		decode := scheme.Codecs.UniversalDeserializer().Decode
		obj, _, err := decode([]byte(f), nil, nil)
		if err != nil {

			return nil, fmt.Errorf("error while decoding yaml object - %s", err)
		}
		runtimeObjects = append(runtimeObjects, obj)
	}

	return runtimeObjects, nil
}
