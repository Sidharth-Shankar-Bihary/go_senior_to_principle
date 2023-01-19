# go_senior_to_principle
Gin — we will use Gin HTTP web framework. 
It is a high-performance framework built on top of net/http that delivers the most essential features, libraries, and functionalities necessary. 
It also has quite clean and comprehensive API.

GORM — is a Golang ORM library developed on top of database/sql. 
It includes all the fancy features like preloading, callbacks, transactions, etc. 
It has some learning curve and not so great documentation, so if you are kind of person who prefers to write raw SQL, then you can just go with sqlx.

db folder contains the code to initialize database connections.
models folder will gold structs


conn.go 
It has InitDB func that passes the DB environment variable and checks, which database driver should be used.
Currently, MySQL, SQLite and Postgres are implemented, but it can be easily extended to all others supported gorm drivers.
After it, connect with a database, it will use gorm to AutoMigrate database structure. T
This is a no-brained if you follow the guides, you can simply put references to all your models into gorm.AutoMigrate().


In models folder, all go files are about database-related logic. 
In other words, it will hold all database structs. And models are only invoked by repos.

In repos folder, all go files are about data logic, it means repo can be used to model, redis, kafka, and so on.
Tests for repos are used the way of Suite Test. It includes BeforeSuite, Describe, AfterSuite.

In services folder, all business logic are in it.
Tests for services are used the way usual tests. It includes init, testClearUp.


JWT_API_SECRET uses "% openssl rand -hex 32" to generate.


Implementation Loopholes

1. The JWT can only be invalidated when it expires. A major limitation to this is: a user can log in, then decide to log out immediately, but the user’s JWT remains valid until the expiration time is reached.
2. The JWT might be hijacked and used by a hacker without the user doing anything about it until the token expires.
3. The user will need to re-login after the token expires, thereby leading to a poor user experience.

Solution

1. Using a persistence storage layer to store JWT metadata. This will enable us to invalidate a JWT the very second the user logs out, thereby improving security.
2. Using the concept of a refresh token to generate a new access token, in the event that the access token expired, thereby improving the user experience.

When a user logs out, I will instantly revoke/invalidate their JWT. This is achieved by deleting the JWT tokenData from our redis store.
