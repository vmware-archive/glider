package handler

import (
	"net/http"

	"github.com/pivotal-golang/lager"
)

func (handler *Handler) AbortBuild(w http.ResponseWriter, r *http.Request) {
	guid := r.FormValue(":guid")

	handler.buildsMutex.RLock()
	build, found := handler.builds[guid]
	handler.buildsMutex.RUnlock()

	if !found {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	log := handler.logger.Session("abort", lager.Data{
		"build": build,
	})

	log.Info("aborting")

	req, err := http.NewRequest(r.Method, build.AbortURL, r.Body)
	if err != nil {
		log.Error("failed-to-create-request", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Error("failed-to-abort", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if resp.StatusCode != http.StatusOK {
		log.Info("bad-abort-response", lager.Data{
			"status": resp.Status,
		})

		resp.Write(w)
		return
	}

	w.WriteHeader(http.StatusOK)

	log.Info("aborted")
}
