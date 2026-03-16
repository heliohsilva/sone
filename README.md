# SONE -> SOcial NEtwork

## Overview

Sone é um projeto (ainda em andamento) feito com o âmbito de estudo (build to learn), com o objetivo de fixar o entendimento sobre construção de API's com autenticação JWT em go.
A construção do frontend para consumir essas API's está em andamento. 

## Routes

O projeto conta atualmente com as seguintes rotas:

```
LOGIN:
  - POST /login -> usuário logar no sistema.

POSTS:
  - POST /posts -> usuário criar um post.
  - GET /posts -> exibir todos os posts cadastrados no sistema.
  - GET /posts/{post_id} -> exibir informações de um post específico.
  - PUT /posts/{post_id} -> modificar um post.
  - DELETE /posts/{post_id} -> excluir um post.
  - GET /posts/{user_id}/posts -> exibir todos os posts de um usuário específico.
  - POST /posts/{post_id}/like -> incrementar o contador de likes de um post em 1.
  - POST /posts/{post_id}/unlike -> decrementar o contador de likes de um post em 1.

USERS:
  - POST /users -> cadastrar um usuário.
  - GET /users{user_id} -> buscar um usuário.
  - GET /users -> listar todos os usuários.
  - PUT /users/{user_id} -> modificar um usuário.
  - DELETE /users/{user_id} -> excluir um usuário.
  - POST /users/{user_id}/follow -> seguir um usuário.
  - POST /users/{user_id}/unfollow -> deixar de seguir um usuário.
  - GET /users/{user_id}/followers -> listar todos os seguidores de um usuário.
  - GET /users/{user_id}/following -> listar todos os usuários que um usuário segue.
  - POST /users/{user_id}/update-password -> atualizar a senha de um usuário.
```

## Arquitetura

Este projeto foi construído em uma arquitetura Model-Controller-Repository, que preserva a responsabilidade única em cada camada sem que cada simples alteração no código seja tecnicamente custosa.

