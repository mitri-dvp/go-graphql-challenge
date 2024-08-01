package schema

import (
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/graphql-go/graphql"
)

var PostList []Post = []Post{
	{
		ID:          uuid.New().String(),
		Title:       "Post 1 title",
		Description: "Lorem, ipsum dolor sit amet consectetur adipisicing elit. Quaerat ea odit eaque amet dicta consequuntur eum dolore commodi error exercitationem, dolorum corporis accusamus esse assumenda obcaecati qui nam illo dolores.",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	},
	{
		ID:          uuid.New().String(),
		Title:       "Post 2 title",
		Description: "Dolore nesciunt aspernatur debitis porro ullam impedit, doloremque deleniti delectus perferendis tempora earum velit dignissimos quam minus voluptate beatae nulla. Natus, molestiae officia fugiat dolor asperiores ex vel. Incidunt, perspiciatis!",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	},
	{
		ID:          uuid.New().String(),
		Title:       "Post 3 title",
		Description: "Minima facere optio cupiditate quisquam, asperiores, voluptatem alias, ducimus quos eum magnam possimus suscipit accusamus. Vero, est nemo! Obcaecati cumque ipsa deleniti laboriosam quaerat doloremque dolores. Maxime deserunt dolores quidem!",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	},
}

type Post struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// Define custom GraphQL ObjectType `postType` for our Golang struct `Post`
// Note that the fields in our postType map with the json tags for the fields in our struct
// and the field type matches the field type in our struct
var postType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Post",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.String,
		},
		"title": &graphql.Field{
			Type: graphql.String,
		},
		"description": &graphql.Field{
			Type: graphql.String,
		},
		"createdAt": &graphql.Field{
			Type: graphql.DateTime,
		},
		"updatedAt": &graphql.Field{
			Type: graphql.DateTime,
		},
	},
})

// Define an input type for Post
var postCreateInputType = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "PostCreateInput",
	Fields: graphql.InputObjectConfigFieldMap{
		"title": &graphql.InputObjectFieldConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
		"description": &graphql.InputObjectFieldConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
	},
})

var postUpdateInputType = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "PostUpdateInput",
	Fields: graphql.InputObjectConfigFieldMap{
		"id": &graphql.InputObjectFieldConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
		"title": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
		"description": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
	},
})

// Root query
var query = graphql.NewObject(graphql.ObjectConfig{
	Name: "query",
	Fields: graphql.Fields{
		"post": &graphql.Field{
			Type:        postType,
			Description: "Get single post",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				id, isOK := params.Args["id"].(string)

				if isOK {
					// Search for element with id
					for _, post := range PostList {
						if post.ID == id {
							return post, nil
						}
					}
				}

				return nil, nil
			},
		},
		"lastPost": &graphql.Field{
			Type:        postType,
			Description: "Last post added",
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				// Sort PostList by CreatedAt in descending order
				sort.Slice(PostList, func(i, j int) bool {
					return PostList[i].CreatedAt.After(PostList[j].CreatedAt)
				})

				// Return the first element (the most recent post)
				return PostList[0], nil
			},
		},
		"postList": &graphql.Field{
			Type:        graphql.NewList(postType),
			Description: "List of posts",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return PostList, nil
			},
		},
	},
})

// Root mutation
var mutation = graphql.NewObject(graphql.ObjectConfig{
	Name: "mutation",
	Fields: graphql.Fields{
		"createPost": &graphql.Field{
			Type:        postType,
			Description: "Create new post",
			Args: graphql.FieldConfigArgument{
				"post": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(postCreateInputType),
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				post := params.Args["post"].(map[string]interface{})
				title := post["title"].(string)
				description := post["description"].(string)

				newPost := Post{
					ID:          uuid.New().String(),
					Title:       title,
					Description: description,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}

				PostList = append(PostList, newPost)

				return newPost, nil
			},
		},
		"updatePost": &graphql.Field{
			Type:        postType,
			Description: "Update existing post",
			Args: graphql.FieldConfigArgument{
				"post": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(postUpdateInputType),
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				// Marshal and cast the argument value
				post := params.Args["post"].(map[string]interface{})
				id := post["id"].(string)
				title, titleIsOk := post["title"].(string)
				description, descriptionIsOk := post["description"].(string)

				// Search list for post with id and change the done variable
				for i := 0; i < len(PostList); i++ {
					if PostList[i].ID == id {
						if titleIsOk {
							PostList[i].Title = title
						}
						if descriptionIsOk {
							PostList[i].Description = description
						}

						PostList[i].UpdatedAt = time.Now()

						return PostList[i], nil
					}
				}
				// Return affected post
				return nil, nil
			},
		},
	},
})

// Define schema, with our query and mutation
var PostSchema, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query:    query,
	Mutation: mutation,
})
