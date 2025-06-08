# TestTaskOzonBank


### Инструкция по установке

1. **Для запуска проекта склонируйте репозиторий:**
```bash
git clone https://github.com/iamstep4ik/TestTaskOzonBank.git
cd TestTaskOzonBank
```

2. **Запуск в Docker**
```bash
docker compose up -d
```

### Пример .env файла
```code
SERVER_HOST=localhost
SERVER_PORT=8000

DB_USER=
DB_PASSWORD=
DB_HOST=
DB_PORT=
DB_NAME=
DB_SSL=


GOOSE_DRIVER=postgres
GOOSE_DBSTRING=postgres://DB_USER:DB_PASSWORD@DB_HOST:DB_PORT/DB_NAME?sslmode=DB_SSL

LOG_LEVEL=info
STORAGE_TYPE=db
```
Для смены типа хранилища на **in-memory**, поменяйте в **.env**  **STORAGE_TYPE** на **memory**

### Для тестирования API 

GraphQL Playground
 * **http://localhost:8080**

### Коллекция запросов Postman


### Примеры запросов:
####  Mutations:
1. Запрос для создания поста
```code
mutation{
  createPost(postInput:{
    authorID: "5c5c6b2a-9655-461a-9415-7e4fc125c1a8"
    title: "bello!"
    content: "yo"
    commentsAllowed: true
  }){
    id
    title
  }
}
```

2. Запрос для создания комментария
```code
mutation{
  createComment(commentInput:{
    authorID: "5c5c6b2a-9655-461a-9415-7e4fc125c1a8"
    content: "nice post"
    parentID: 1
    postID: 3
  }){
    id
    created_at
  }
}
```
3. Запрос для обновления разрешения комментировать посты
```code
mutation{
  updateAllowComments(
    authorID: "5c5c6b2a-9655-461a-9415-7e4fc125c1a8"
    postID: 3
    commentsAllowed: false
  ){
    id
    title
    content
    created_at
  }
}
```

####  Queries:
1. Вывод определенного поста по id
```code
query{
  post(
    postID: 2
  ){
    authorID
    id
    title
    content
    created_at
  }
}
```

2. Вывод всех постов и комментариев к нему
```code
query{
  posts{
    title
    content
    comments{
      depth
      parentID
      content
      replies{
        content
      }
    }
  }
}
```


####  Subscription:

Позволяет подписаться на уведомления по новым комментариям к посту
```code
subscription {
  commentAdded(postID: 2) {
    id
    content
    authorID
    created_at
  }
}
```
