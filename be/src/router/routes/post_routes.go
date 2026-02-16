package routes

import (
	"api/src/controllers"
	"net/http"
)

var postRoutes = []Route{
	{
		URI:      "/posts",
		Method:   http.MethodPost,
		Function: controllers.CreatePost,
		NeedAuth: true,
	},
	{
		URI:      "/posts",
		Method:   http.MethodGet,
		Function: controllers.GetPosts,
		NeedAuth: true,
	},
	{
		URI:      "/posts/{postID}",
		Method:   http.MethodGet,
		Function: controllers.GetPost,
		NeedAuth: true,
	},
	{
		URI:      "/posts/{postID}",
		Method:   http.MethodPut,
		Function: controllers.UpdatePost,
		NeedAuth: true,
	},
	{
		URI:      "/posts/{postID}",
		Method:   http.MethodDelete,
		Function: controllers.DeletePost,
		NeedAuth: true,
	},
	{
		URI:      "/users/{userID}/posts",
		Method:   http.MethodGet,
		Function: controllers.GetPostsByUser,
		NeedAuth: true,
	},
	{
		URI:      "/posts/{postID}/like",
		Method:   http.MethodPost,
		Function: controllers.Like,
		NeedAuth: true,
	},
	{
		URI:      "/posts/{postID}/unlike",
		Method:   http.MethodPost,
		Function: controllers.Unlike,
		NeedAuth: true,
	},
}
