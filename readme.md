
# URL Shortener (Go Basic)

This app is built when i am a bit bored in my weekend night.  
How to use this app is:  

    1. Go to docker-compose.yaml file and up the compose 
    2. Go to api folder and go run main.go 
    3. You can hit the get and post to test it. 



# Curl for create and get the short URL
## Create short url
### Curl Command
    curl --request POST \
    --url http://localhost:5000/api/v1/short-url \
    --header 'Content-Type: application/json' \
    --data '{
        "long-url":"https://shamirhusein.my.id", 
        "email": "shamirhusein@gmail.com",
        "short-url":"ssss"
    }'
### Response
    {
        "error": false,
        "msg": null,
        "shortUrl": "000000b"
    }


## Get by short url
### Curl Command
    curl --request GET \
    --url http://localhost:5000/api/v1/short-url/000000a

### Response
    {
        "message": "record not found"
    }




## Tech Stack

**api:** Go and Fiber Go
**Server:** Postgre and Etcd



## Tech Stack

**Client:** React, Redux, TailwindCSS

**Server:** Node, Express


## 🔗 Links
[![portfolio and blog](https://img.shields.io/badge/my_portfolio-000?style=for-the-badge&logo=ko-fi&logoColor=white)](https://shamirhusein.my.id/)
[![linkedin](https://img.shields.io/badge/linkedin-0A66C2?style=for-the-badge&logo=linkedin&logoColor=white)](https://www.linkedin.com/in/shamirhusein)


