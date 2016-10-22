package httputil

import (
	"log"
	"net/http"
	"testing"
)

func TestUnpackURLForm(t *testing.T) {
	var data struct {
		Name  string `http:"n"`
		Age   int    `http:"a"`
		Email string `http:"e"`
	}

	req, err := http.NewRequest("GET", "http://example.com/foo?n=echo&a=18&e=echo@tonnn.com", nil)
	if err != nil {
		log.Fatal(err)
	}

	err = UnpackURLForm(req, &data)
	if err != nil {
		t.Errorf("occur error: %v\n", err)
	}
	if data.Name != "echo" {
		t.Errorf("got %q, want echo", string(data.Name))
	}

}
