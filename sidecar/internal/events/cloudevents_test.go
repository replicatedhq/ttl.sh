package events

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestParseBinaryMode(t *testing.T) {
	body, _ := json.Marshal(ImageUpdatedData{
		Name:      "foo/bar",
		Reference: "1h",
		Digest:    "sha256:deadbeef",
	})
	req := httptest.NewRequest(http.MethodPost, "/events", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Ce-Type", ImageUpdatedType)

	evt, err := Parse(req, body)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if evt.Type != ImageUpdatedType {
		t.Errorf("type = %q want %q", evt.Type, ImageUpdatedType)
	}
	if evt.Data.Name != "foo/bar" || evt.Data.Reference != "1h" {
		t.Errorf("data = %+v", evt.Data)
	}
}

func TestParseStructuredMode(t *testing.T) {
	// Includes the envelope keys this package deliberately does not model
	// (source, id, datacontenttype); they must be ignored, not rejected.
	body := []byte(`{
		"type": "` + ImageUpdatedType + `",
		"source": "zot",
		"id": "abc",
		"datacontenttype": "application/json",
		"data": {"name":"foo/bar","reference":"30m","digest":"sha256:abc"}
	}`)
	req := httptest.NewRequest(http.MethodPost, "/events", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/cloudevents+json; charset=utf-8")

	evt, err := Parse(req, body)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if evt.Type != ImageUpdatedType {
		t.Errorf("type = %q", evt.Type)
	}
	if evt.Data.Reference != "30m" {
		t.Errorf("reference = %q", evt.Data.Reference)
	}
}

func TestParseStructuredEmptyData(t *testing.T) {
	env := structuredEnvelope{Type: "some.other.type"}
	body, _ := json.Marshal(env)
	req := httptest.NewRequest(http.MethodPost, "/events", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/cloudevents+json")
	evt, err := Parse(req, body)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if evt.Type != "some.other.type" {
		t.Errorf("type = %q", evt.Type)
	}
	if evt.Data.Name != "" {
		t.Errorf("expected zero data, got %+v", evt.Data)
	}
}

func TestParseStructuredMalformed(t *testing.T) {
	body := []byte(`{not json`)
	req := httptest.NewRequest(http.MethodPost, "/events", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/cloudevents+json")
	if _, err := Parse(req, body); err == nil {
		t.Fatal("expected error for malformed structured envelope")
	}
}

func TestParseBinaryMalformed(t *testing.T) {
	body := []byte(`{not json`)
	req := httptest.NewRequest(http.MethodPost, "/events", bytes.NewReader(body))
	req.Header.Set("Ce-Type", ImageUpdatedType)
	if _, err := Parse(req, body); err == nil {
		t.Fatal("expected error for malformed binary body")
	}
}

func TestParseStructuredDataDecodeError(t *testing.T) {
	// Valid envelope, but `data` is a JSON array which cannot unmarshal into
	// the ImageUpdatedData struct -> exercises the "decode structured data"
	// error path.
	body := []byte(`{"type":"zotregistry.image.updated","data":[1,2,3]}`)
	req := httptest.NewRequest(http.MethodPost, "/events", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/cloudevents+json")

	if _, err := Parse(req, body); err == nil {
		t.Fatal("expected error decoding structured data array into struct")
	}
}

func TestParseBinaryEmptyBody(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/events", nil)
	req.Header.Set("Ce-Type", ImageUpdatedType)
	evt, err := Parse(req, nil)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if evt.Type != ImageUpdatedType {
		t.Errorf("type = %q", evt.Type)
	}
	if evt.Data.Name != "" {
		t.Errorf("expected zero data, got %+v", evt.Data)
	}
}
