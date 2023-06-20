## Robinhood test

By Tanatorn Nateesanprasert

---
This application contains 2 main components
1. Golang application container (100% test on services layer)
2. MongoDB container

P.S. MongoDB container have 3 replicas for doing operation in transaction

---
#### How to run?
1. run docker command on terminal: `docker compose up`
2. application runs on port `8080`
3. swagger url: `http://localhost:8080/swagger/index.html`
4. call register api and login to get token then authorize with value `Bearer {token}`

https://github.com/thecuriousbig/robinhood-test/assets/25192525/3c3062a4-83c5-400f-a83a-f28f0e6531cd

6. call other apis





---
#### REST APIS

user related
1. register: `[POST] /api/v1/user/register`
2. login: `[GET] /api/v1/user/login`
3. (required login) update user:  `[PUT] /api/v1/user`

blog related
1. (required login) create blog: `[POST] /api/v1/blog`
2. (required login) list blog: `[GET] /api/v1/blog?page={page}&limit={limit}`
3. (required login) get blog by id: `[GET] /api/v1/blog/:blogId`
4. (required login) update blog status: `[PUT] /api/v1/blog/:blogId`
5. (required login) archive blog: `[DELETE] /api/v1/blog/:blogId`

comment related
1. (required login) create comment: `[POST] /api/v1/comment/:blogId`
2. (required login) list comment: `[GET] /api/v1/comment/:blogId`
