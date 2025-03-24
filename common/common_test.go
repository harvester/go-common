package common

import (
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type mockObject struct {
	metav1.ObjectMeta
}

func TestGetNamespacedName(t *testing.T) {
	testCases := []struct {
		name        string
		input       any
		expectedRet string
		expectedErr bool
	}{
		{
			name: "valid object",
			input: &mockObject{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-name",
					Namespace: "test-ns",
				},
			},
			expectedRet: "test-ns/test-name",
			expectedErr: false,
		},
		{
			name:        "nil object",
			input:       nil,
			expectedRet: "",
			expectedErr: true,
		},
		{
			name: "empty namespace",
			input: &mockObject{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-name",
					Namespace: "",
				},
			},
			expectedRet: "",
			expectedErr: true,
		},
		{
			name: "empty name",
			input: &mockObject{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "",
					Namespace: "test-ns",
				},
			},
			expectedRet: "",
			expectedErr: true,
		},
		{
			name: "empty name and namespace",
			input: &mockObject{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "",
					Namespace: "",
				},
			},
			expectedRet: "",
			expectedErr: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			result, err := GetNamespacedName(testCase.input)
			if (err != nil) != testCase.expectedErr {
				t.Errorf("GetNamespacedName() error = %v, wantErr %v", err, testCase.expectedErr)
				return
			}
			if result != testCase.expectedRet {
				t.Errorf("GetNamespacedName() = %v, want %v", result, testCase.expectedRet)
			}
		})
	}
}
