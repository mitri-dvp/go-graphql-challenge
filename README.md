# Go Graphql Challenge

This is a simple Graphql API written in Go to manage Posts. This was based on the [graphql-go todo example](https://github.com/graphql-go/graphql/tree/master/examples/todo).

I've adjusted the HTTP handler to handle HTTP POST request to be compatible with how GraphQL operations are requested. Making this server compatible with GraphQL Query Clients.

## Features

- Post queries
  - post: Get a Post by ID.
  - lastPost: Get last createdAt Post.
  - postList: Get list of Posts.
- Post mutations
  - createPost: Create a new Post.
  - updatePost: Update an existing Post.

## Sample queries

```gql
query GetPost($id: String) {
  post(id: $id) {
    id
    title
    description
    createdAt
    updatedAt
  }
}
```

```gql
query GetAllPosts {
  postList {
    id
  }
}
```

```gql
mutation CreatePost($post: PostCreateInput!) {
  createPost(post: $post) {
    id
    title
    description
    createdAt
    updatedAt
  }
}
```

```gql
mutation UpdatePost($post: PostUpdateInput!) {
  updatePost(post: $post) {
    id
    title
    description
    createdAt
    updatedAt
  }
}
```
