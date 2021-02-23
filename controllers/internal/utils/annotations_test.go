package utils

import (
	"testing"

	"k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/stretchr/testify/assert"
)

func TestAddNewAnnotations(t *testing.T) {
	// arrange
	target, source := provideIngresses()
	// act
	MergeAnnotations(&target.ObjectMeta, &source.ObjectMeta)
	// assert
	assert.Equal(t, 3, len(target.ObjectMeta.Annotations))
	assert.Equal(t, "eu", target.ObjectMeta.Annotations["k8gb.io/primary-geotag"])
	assert.Equal(t, "dummy", target.ObjectMeta.Annotations["field.cattle.io/publicEndpoints"])
}

func TestAddExistingAnnotations(t *testing.T) {
	// arrange
	target, source := provideIngresses()
	for k, v := range source.Annotations {
		target.Annotations[k] = v
	}
	// act
	MergeAnnotations(&target.ObjectMeta, &source.ObjectMeta)
	// assert
	assert.Equal(t, 3, len(target.ObjectMeta.Annotations))
	assert.Equal(t, "eu", target.ObjectMeta.Annotations["k8gb.io/primary-geotag"])
	assert.Equal(t, "dummy", target.ObjectMeta.Annotations["field.cattle.io/publicEndpoints"])
	assert.Equal(t, "failover", target.ObjectMeta.Annotations["k8gb.io/strategy"])
}

func TestUpdateExistingRecords(t *testing.T) {
	// arrange
	target, source := provideIngresses()
	for k, v := range source.Annotations {
		target.Annotations[k] = v
	}
	target.Annotations["k8gb.io/primary-geotag"] = "us"
	// act
	MergeAnnotations(&target.ObjectMeta, &source.ObjectMeta)
	// assert
	assert.Equal(t, 3, len(target.ObjectMeta.Annotations))
	assert.Equal(t, "us", target.ObjectMeta.Annotations["k8gb.io/primary-geotag"])
	assert.Equal(t, "dummy", target.ObjectMeta.Annotations["field.cattle.io/publicEndpoints"])
	assert.Equal(t, "failover", target.ObjectMeta.Annotations["k8gb.io/strategy"])
}

func TestEqualAnnotationsWithEmptyTarget(t *testing.T) {
	// arrange
	_, source := provideIngresses()
	target := &v1beta1.Ingress{}
	// act
	MergeAnnotations(&target.ObjectMeta, &source.ObjectMeta)
	// assert
	assert.True(t, assert.ObjectsAreEqual(source.Annotations, target.ObjectMeta.Annotations))
}

func TestEqualAnnotationsWithEmptySource(t *testing.T) {
	// arrange
	target, _ := provideIngresses()
	source := &v1beta1.Ingress{}
	// act
	MergeAnnotations(&target.ObjectMeta, &source.ObjectMeta)
	// assert
	assert.True(t, assert.ObjectsAreEqual(target.Annotations, target.ObjectMeta.Annotations))
}

func TestEqualAnnotationsWithEmptyInput(t *testing.T) {
	// arrange
	source := &v1beta1.Ingress{}
	target := &v1beta1.Ingress{}
	// act
	MergeAnnotations(&target.ObjectMeta, &source.ObjectMeta)
	// assert
	assert.Equal(t, 0, len(target.ObjectMeta.Annotations))
}

func TestContainsAllAnnotations(t *testing.T) {
	// arrange
	source, target := provideIngresses()
	metav1.SetMetaDataAnnotation(&target.ObjectMeta, "k8gb.io/primary-geotag", "eu")
	metav1.SetMetaDataAnnotation(&target.ObjectMeta, "k8gb.io/strategy", "failover")
	// act
	result := ContainsAnnotations(&target.ObjectMeta, &source.ObjectMeta)
	// assert
	assert.True(t, result)
}

func TestContainsSomeAnnotations(t *testing.T) {
	// arrange
	source, target := provideIngresses()
	metav1.SetMetaDataAnnotation(&target.ObjectMeta, "k8gb.io/primary-geotag", "eu")
	// act
	result := ContainsAnnotations(&target.ObjectMeta, &source.ObjectMeta)
	// assert
	assert.False(t, result)
}

func TestContainsAnnotationsWithDifferentValues(t *testing.T) {
	// arrange
	source, target := provideIngresses()
	// act
	result := ContainsAnnotations(&target.ObjectMeta, &source.ObjectMeta)
	// assert
	assert.False(t, result)
}

func TestContainsSourceIsEmpty(t *testing.T) {
	// arrange
	_, target := provideIngresses()
	source := &v1beta1.Ingress{}
	// act
	result := ContainsAnnotations(&target.ObjectMeta, &source.ObjectMeta)
	// assert
	assert.True(t, result)
}

func TestContainsTargetIsEmpty(t *testing.T) {
	// arrange
	source, _ := provideIngresses()
	target := &v1beta1.Ingress{}
	// act
	result := ContainsAnnotations(&target.ObjectMeta, &source.ObjectMeta)
	// assert
	assert.False(t, result)
}

func TestContainsEmptyInputs(t *testing.T) {
	// arrange
	source := &v1beta1.Ingress{}
	target := &v1beta1.Ingress{}
	// act
	result := ContainsAnnotations(&target.ObjectMeta, &source.ObjectMeta)
	// assert
	assert.True(t, result)
}

func provideIngresses() (isource *v1beta1.Ingress, itarget *v1beta1.Ingress) {
	source := map[string]string{"k8gb.io/primary-geotag": "eu", "k8gb.io/strategy": "failover"}
	target := map[string]string{"field.cattle.io/publicEndpoints": "dummy"}
	isource = &v1beta1.Ingress{}
	itarget = &v1beta1.Ingress{}
	isource.Annotations = source
	itarget.Annotations = target
	return
}
