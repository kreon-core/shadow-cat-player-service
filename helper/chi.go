package helper

import (
	"bytes"
	"encoding/json"
	"net/http"

	tul "github.com/kreon-core/shadow-cat-common"
	"github.com/kreon-core/shadow-cat-common/logc"
)

func PlainText(w http.ResponseWriter, statusCode int, payload string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(statusCode)

	if _, err := w.Write([]byte(payload)); err != nil {
		logc.Warn().Err(err).Msg("Failed to write plain text response")
		return
	}
}

func JSON(w http.ResponseWriter, statusCode int, payload any) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(true)

	if err := enc.Encode(payload); err != nil {
		logc.Error().Err(err).Msg("Failed to encode JSON response")
		http.Error(w, tul.Message(tul.UUnspecifiedError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)

	if _, err := w.Write(buf.Bytes()); err != nil {
		logc.Warn().Err(err).Msg("Failed to write JSON response")
		return
	}
}
