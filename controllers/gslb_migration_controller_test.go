package controllers

import (
	"context"
	"testing"

	k8gbv1beta1 "github.com/k8gb-io/k8gb/api/v1beta1"
	k8gbv1beta1io "github.com/k8gb-io/k8gb/api/v1beta1io"
	"github.com/k8gb-io/k8gb/controllers/mocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func TestLegacyMigrationCreatesIOAndLabelsLegacy(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	legacy := &k8gbv1beta1.Gslb{ObjectMeta: metav1.ObjectMeta{Name: "demo", Namespace: "default"}}
	cl := mocks.NewMockClient(ctrlMock)

	cl.EXPECT().Get(gomock.Any(), types.NamespacedName{Name: "demo", Namespace: "default"}, gomock.Any()).
		DoAndReturn(func(_ context.Context, _ types.NamespacedName, obj client.Object, _ ...client.GetOption) error {
			*obj.(*k8gbv1beta1.Gslb) = *legacy
			return nil
		})

	cl.EXPECT().Get(gomock.Any(), types.NamespacedName{Name: "demo", Namespace: "default"}, gomock.Any()).
		Return(apierrors.NewNotFound(schema.GroupResource{Group: "k8gb.io", Resource: "gslbs"}, "demo"))

	cl.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, obj client.Object, _ ...client.CreateOption) error {
			_, ok := obj.(*k8gbv1beta1io.Gslb)
			require.True(t, ok)
			return nil
		})

	cl.EXPECT().Patch(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil)

	r := &LegacyGslbMigrationReconciler{Client: cl, Recorder: record.NewFakeRecorder(5)}

	_, err := r.Reconcile(context.Background(), ctrl.Request{NamespacedName: types.NamespacedName{Name: "demo", Namespace: "default"}})
	require.NoError(t, err)
}
