package internal // Package internal provides the HTTP handlers for the task service.

import (
	"io"
	"net/http"
)

// Router is a struct that holds the configuration for the task service.
type Router struct {
	TaskServiceURL string
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.URL.Path {
	case "/":
		http.Redirect(w, req, "/health", http.StatusSeeOther)
	case "/health":
		r.healthHandler(w, req)
	case "/create-task":
		r.forwardToTaskService(w, req)
	default:
		http.NotFound(w, req)
	}
}

func (r *Router) healthHandler(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("API Gateway is healthy"))
}

func (r *Router) forwardToTaskService(w http.ResponseWriter, req *http.Request) {
	proxyReq, err := http.NewRequest(req.Method, r.TaskServiceURL+"/create-task", req.Body)
	if err != nil {
		http.Error(w, "Error forwarding request", http.StatusInternalServerError)
		return
	}

	proxyReq.Header = req.Header

	resp, err := http.DefaultClient.Do(proxyReq)
	if err != nil {
		http.Error(w, "Task Service is unavailable", http.StatusServiceUnavailable)
		return
	}

	defer req.Body.Close()

	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}
