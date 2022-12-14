package server

import (
	"context"
	"github.com/Barton0403/btgo-pkg/oauth2"
	"github.com/Barton0403/btgo-pkg/oauth2/store"
	goauth2 "github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/errors"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/server"
	"log"
)

type Server struct {
	*server.Server
	tokenStore *store.TokenStore
}

func (s *Server) SetClientStorage(store goauth2.ClientStore) {
	s.Server.Manager.(*manage.Manager).MapClientStorage(store)
}

func (s *Server) SetTokenStorage(gostore goauth2.TokenStore) {
	s.Server.Manager.(*manage.Manager).MapTokenStorage(gostore)
	s.tokenStore = gostore.(*store.TokenStore)
}

func (s *Server) GetManager() goauth2.Manager {
	return s.Manager
}

func (s *Server) Logout(ctx context.Context, id string) {
	s.tokenStore.RemoveAll(ctx, id)
}

func NewServer() oauth2.Server {
	manager := manage.NewDefaultManager()
	manager.MapTokenStorage(&store.TokenStore{})

	srv := server.NewDefaultServer(manager)
	srv.SetAllowGetAccessRequest(true)

	srv.SetInternalErrorHandler(func(err error) (re *errors.Response) {
		log.Println("Oauth2 Internal Error:", err.Error())
		return
	})

	srv.SetResponseErrorHandler(func(re *errors.Response) {
		log.Println("Oauth2 Response Error:", re.Error.Error())
	})

	return &Server{Server: srv}
}
