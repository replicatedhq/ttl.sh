// Package events parses zot's events-extension CloudEvents.
//
// zot uses the cloudevents/sdk-go HTTP protocol, which defaults to binary
// content mode: event metadata travels in Ce-* headers and the JSON data field
// is the request body. This package also accepts structured mode (Content-Type
// application/cloudevents+json), where the whole envelope is the body.
package events

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

const (
	ImageUpdatedType = "zotregistry.image.updated"
)

// ImageUpdatedData is the subset of zot's event payload this service acts on.
// Fields it does not read are deliberately absent: unmarshal ignores unknown
// keys, so declaring them would only add ways for a type mismatch to reject an
// otherwise usable event.
type ImageUpdatedData struct {
	Name      string `json:"name"`
	Reference string `json:"reference"`
	Digest    string `json:"digest"`
}

type structuredEnvelope struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

// Event is the minimum we need from any incoming POST.
type Event struct {
	Type string
	Data ImageUpdatedData
}

func Parse(r *http.Request, body []byte) (Event, error) {
	contentType := r.Header.Get("Content-Type")

	// Structured mode: full envelope is the body.
	if strings.HasPrefix(contentType, "application/cloudevents+json") {
		var env structuredEnvelope
		if err := json.Unmarshal(body, &env); err != nil {
			return Event{}, fmt.Errorf("decode structured envelope: %w", err)
		}
		var data ImageUpdatedData
		if len(env.Data) > 0 {
			if err := json.Unmarshal(env.Data, &data); err != nil {
				return Event{}, fmt.Errorf("decode structured data: %w", err)
			}
		}
		return Event{Type: env.Type, Data: data}, nil
	}

	// Binary mode: Ce-* headers + JSON body. Header.Get canonicalizes the key,
	// so a lowercased ce-type from an intermediary resolves here too.
	t := r.Header.Get("Ce-Type")
	var data ImageUpdatedData
	if len(body) > 0 {
		if err := json.Unmarshal(body, &data); err != nil {
			return Event{}, fmt.Errorf("decode binary data: %w", err)
		}
	}
	return Event{Type: t, Data: data}, nil
}
