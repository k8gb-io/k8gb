package v1beta1io

import (
	"testing"

	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
)

func TestGroupVersionAndScheme(t *testing.T) {
	require.Equal(t, "k8gb.io", GroupVersion.Group)
	require.Equal(t, "v1beta1", GroupVersion.Version)

	scheme := runtime.NewScheme()
	require.NoError(t, AddToScheme(scheme))

	gvk, err := apiutil.GVKForObject(&Gslb{}, scheme)
	require.NoError(t, err)
	require.Equal(t, schema.GroupVersionKind{Group: "k8gb.io", Version: "v1beta1", Kind: "Gslb"}, gvk)
}
