## Indocoffee REST API 
indocoffee is an online coffee shop application

## To use the application in local
Requirements: 
  - Go 1.24 or newer
  - Swaggo
  - golang-migrate
  - air live reload 
  - mysql
  - redis

## To run the application in local
 - Install all packages
   ``` Go mod download ```
 - Add .env file
 - Run migration using makefile command
   ``` make migrate-up ```
 - Run air live reload
  ``` air ```
- To see the api documentation you can go to [visit](http://localhost:8080/v1/swagger/index.html)

## To run with docker
 - Set environment variables in .env file 
 - Run ```docker compose --build```
 - To see the api documentation you can go to [visit](http://localhost:8080/v1/swagger/index.html)

## Web applications demo version
 - client website to buy coffee:
   - visit : https://indocoffee-website.vercel.app 
 - client website to manage selling coffee:
    - visit : https://indocoffee-web-dashboard.vercel.app
     - credential to login as admin:
        - email: lizzy@example.com
        - password: Password$123


