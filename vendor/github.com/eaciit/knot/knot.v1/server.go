package knot

import (
	"fmt"
	"net/http"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/eaciit/toolkit"
)

type Router struct {
	// implements golang own multiplexor
	mux *http.ServeMux

	// map of routes, the value is pointer of interface `http.Handler`.
	// this interface should have method `ServeHTTP(w ResponseWriter, r *Request)`
	// this variable is not used in routing process, but used on `GetHandler()`
	// because golang does not provide function to get handler using specific route
	routes map[string]http.Handler
}

func (r *Router) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	if _, ok := r.routes[pattern]; ok {
		return
	}

	r.mux.HandleFunc(pattern, handler)

	// because the `routes` value must be interface `http.Handler`,
	// we have to create some struct which contains method `ServeHTTP`,
	// and fill the method with closure `handler` (2nd parameter)
	r.routes[pattern] = http.HandlerFunc(handler)
}

func (r *Router) Handle(pattern string, handler http.Handler) {
	if _, ok := r.routes[pattern]; ok {
		return
	}

	r.mux.Handle(pattern, handler)
	r.routes[pattern] = handler
}

func (r *Router) GetHandler(path string) http.Handler {
	if val, ok := r.routes[path]; ok {
		return val
	} else {
		return nil
	}
}

type Server struct {
	Address string

	mxrouter        *Router
	log             *toolkit.LogEngine
	status          chan string
	UseSSL          bool
	CertificatePath string
	PrivateKeyPath  string

	preRequest  FnContent
	postRequest FnContent
}

func (s *Server) Log() *toolkit.LogEngine {
	if s.log == nil {
		s.log, _ = toolkit.NewLog(true, false, "", "", "")
	}
	return s.log
}

type FnContent func(r *WebContext) interface{}

func (s *Server) PreRequest(c FnContent) {
	s.preRequest = c
}

func (s *Server) PostRequest(c FnContent) {
	s.postRequest = c
}

func (s *Server) router() *Router {
	if s.mxrouter == nil {
		s.mxrouter = &Router{mux: http.NewServeMux()}
		s.mxrouter.routes = map[string]http.Handler{}
	}
	return s.mxrouter
}

func (s *Server) Register(c interface{}, prefix string) error {
	return s.RegisterWithConfig(c, prefix, NewResponseConfig())
}

func (s *Server) RegisterWithConfig(c interface{}, prefix string, cfg *ResponseConfig) error {
	var t reflect.Type
	v := reflect.ValueOf(c)
	if v.Kind() != reflect.Ptr {
		return fmt.Errorf("Invalid controller object passed (%s). Controller object should be a pointer", v.Kind())
	}
	t = reflect.TypeOf(c)
	controllerName := reflect.Indirect(v).Type().Name()

	s.Log().Info(fmt.Sprintf("Registering %s", controllerName))
	path := prefix
	fixUrlPath(&path, true, true)
	controllerName = strings.ToLower(controllerName)
	if strings.HasSuffix(controllerName, "controller") {
		controllerName = controllerName[0 : len(controllerName)-len("controller")]
	}
	path += controllerName + "/"

	if t == nil {
	}
	methodCount := t.NumMethod()
	for mi := 0; mi < methodCount; mi++ {
		method := t.Method(mi)

		// validate if this method match FnContent
		isFnContent := false
		tm := method.Type
		if tm.NumIn() == 2 && tm.In(1).String() == "*knot.WebContext" {
			if tm.NumOut() == 1 && tm.Out(0).Kind() == reflect.Interface {
				isFnContent = true
			}
		}

		if isFnContent {
			var fnc FnContent
			fnc = v.MethodByName(method.Name).Interface().(func(*WebContext) interface{})
			methodName := method.Name
			handlerPath := path + strings.ToLower(methodName)
			newcfg := NewResponseConfig()
			*newcfg = *cfg
			newcfg.ControllerName = controllerName
			newcfg.MethodName = methodName
			s.RouteWithConfig(handlerPath, fnc, newcfg)
		}
	}

	return nil
}

func fixUrlPath(urlPath *string, preSlash, postSlash bool) {
	path := strings.ToLower(*urlPath)
	if preSlash && strings.HasPrefix(path, "/") == false {
		path = "/" + path
	}
	if postSlash && strings.HasSuffix(path, "/") == false {
		path += "/"
	}
	*urlPath = path
}

func fixLogicalPath(logicalPath *string, preSlash, postSlash bool) {
	path := strings.ToLower(*logicalPath)
	if preSlash && strings.HasPrefix(path, "/") == false {
		path = "/" + path
	}
	if postSlash && strings.HasSuffix(path, "/") == false {
		path += "/"
	}
	*logicalPath = path
}

func (s *Server) RouteStatic(pathUrl, path string) {
	_, ePath := os.Stat(path)
	if ePath != nil {
		s.Log().Error(fmt.Sprintf("Unable to add static %s from %s : %s", pathUrl, path, ePath.Error()))
		return
	}

	fixUrlPath(&pathUrl, true, true)
	s.Log().Info(fmt.Sprintf("Add static %s from %s", pathUrl, path))
	fsHandler := http.StripPrefix(pathUrl, http.FileServer(http.Dir(path)))
	s.router().Handle(pathUrl, GzipHandler(fsHandler))
}

func (s *Server) Route(path string, fnc FnContent) {
	s.RouteWithConfig(path, fnc, NewResponseConfig())
}

func (s *Server) RouteWithConfig(path string, fnc FnContent, cfg *ResponseConfig) {
	fixUrlPath(&path, true, false)
	s.Log().Info(fmt.Sprintf("Registering handler for %s", path))
	s.router().HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				s.Log().Errorf("Panic error detected: %v", rec)
			}
		}()

		if fnc != nil {
			rcfg := NewResponseConfig()
			*rcfg = *cfg

			app := rcfg.App
			kr := new(WebContext)
			kr.cookieStore = new(CookieStore)
			kr.Server = s
			kr.Request = r
			kr.Writer = w
			kr.Config = rcfg
			if app != nil && int(rcfg.OutputType) == 0 {
				rcfg.OutputType = app.DefaultOutputType
			}

			if s.preRequest != nil {
				s.preRequest(kr)
			}

			v := fnc(kr)

			if app != nil && app.requireValidation && !kr.Config.IgnoreValidation {
				valid := true

				if app.fnRedirectUrl == nil {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte("500 - Require session info is not complete"))
					return
				}

				valid = app.fnValidate(kr)

				if !valid {
					url := ""
					if app.fnRedirectUrl != nil {
						url = app.fnRedirectUrl(kr)
					}
					if url != "" {
						http.Redirect(w, r, url, http.StatusTemporaryRedirect)
					} else {
						w.WriteHeader(http.StatusInternalServerError)
						w.Write([]byte("500 - Require session info is not complete"))
						return
					}
				}
			}

			if s.postRequest != nil {
				s.postRequest(kr)
			}

			if kr.Config.NoLog == false {
				s.Log().Info(fmt.Sprintf("%s%s %s",
					s.Address, r.URL.String(), r.RemoteAddr))
			}

			kr.WriteCookie()
			kr.Write(v)
		} else {
			w.Write([]byte(""))
		}
	})
}

func (s *Server) GetHandler(path string) http.Handler {
	return s.GetHandler(path)
}

func (s *Server) GetAddress() string {
	address := s.Address

	// when using SSL enabled `http://` and `https://` need to be stripped,
	// because it will make the routes unacessable
	if strings.Contains(address, "https") {
		return strings.Replace(address, "https://", "", -1)
	} else if strings.Contains(address, "http") {
		return strings.Replace(address, "http://", "", -1)
	}

	return address
}

func (s *Server) isReadyForSSL() bool {
	if s.CertificatePath == "" || s.PrivateKeyPath == "" {
		s.Log().Error("Both certificate.pem and privatekey.pem full path should be defined when using SSL")
		return false
	}

	return true
}

func (s *Server) Listen() {
	s.start()
	s.listen()
}

func (s *Server) start() error {
	addr := s.GetAddress()
	s.status = make(chan string)
	s.Log().Info("Start listening on server " + addr)

	go func() {
		if s.UseSSL {
			if !s.isReadyForSSL() {
				return
			}

			http.ListenAndServeTLS(addr, s.CertificatePath, s.PrivateKeyPath, s.router().mux)
		} else {
			http.ListenAndServe(addr, s.router().mux)
		}
	}()
	return nil
}

func (s *Server) Stop() {
	s.Log().Info(fmt.Sprintf("Stopping server %s", s.Address))
	go func() {
		time.Sleep(1 * time.Second)
		s.status <- "Stop"
	}()
}

func (s *Server) listen() {
	running := true
	for running {
		select {
		case status := <-s.status:
			if status == "Stop" {
				running = false
			}
		}
	}
}
