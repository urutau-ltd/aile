package resource_test

import (
	"fmt"
	"net/http"

	"codeberg.org/urutau-ltd/aile/v2/x/resource"
)

type collectionExampleHandler struct{}

func (collectionExampleHandler) Index(http.ResponseWriter, *http.Request)  {}
func (collectionExampleHandler) New(http.ResponseWriter, *http.Request)    {}
func (collectionExampleHandler) Create(http.ResponseWriter, *http.Request) {}
func (collectionExampleHandler) Show(http.ResponseWriter, *http.Request)   {}
func (collectionExampleHandler) Edit(http.ResponseWriter, *http.Request)   {}
func (collectionExampleHandler) Update(http.ResponseWriter, *http.Request) {}
func (collectionExampleHandler) Delete(http.ResponseWriter, *http.Request) {}

type singletonExampleHandler struct{}

func (singletonExampleHandler) Show(http.ResponseWriter, *http.Request)   {}
func (singletonExampleHandler) Edit(http.ResponseWriter, *http.Request)   {}
func (singletonExampleHandler) Update(http.ResponseWriter, *http.Request) {}

func ExampleCollection() {
	var h resource.Collection = collectionExampleHandler{}
	fmt.Println(h != nil)
	// Output: true
}

func ExampleSingleton() {
	var h resource.Singleton = singletonExampleHandler{}
	fmt.Println(h != nil)
	// Output: true
}
