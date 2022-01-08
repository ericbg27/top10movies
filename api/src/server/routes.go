package server

func (s *Server) routes() {
	s.Router.POST("/login", s.UsersController.Login)
	s.Router.POST("/register", s.UsersController.Create)
	s.Router.POST("/users/:user_id", s.UsersController.Update)
	s.Router.PATCH("/users/:user_id", s.UsersController.Update)
	s.Router.DELETE("/users/:user_id", s.UsersController.Delete)
	s.Router.GET("/users/search", s.UsersController.Search)

	s.Router.GET("/users/:user_id/favorites", s.UsersController.GetFavorites)
	s.Router.POST("/users/:user_id/favorite", s.UsersController.AddFavorite) // TODO: Do we put movie_id in the URL?

	s.Router.GET("/search", s.MoviesController.Search)
}
