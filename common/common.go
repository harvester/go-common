package common

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetNamespacedName returns the namespaced name of the object.
// If the object is not supported, it returns "NOSUPPORTED_OBJECT_TYPE".
func GetNamespacedName[T any](t T) (string, error) {
	switch obj := any(t).(type) {
	case metav1.Object:
		if obj.GetNamespace() == "" {
			return "", fmt.Errorf("input object type: %T should have namespaced", t)
		}
		if obj.GetName() == "" {
			return "", fmt.Errorf("input object type: %T should have name", t)
		}
		return fmt.Sprintf("%s/%s", obj.GetNamespace(), obj.GetName()), nil
	default:
		return "", fmt.Errorf("unsupported object type: %T", t)
	}
}
