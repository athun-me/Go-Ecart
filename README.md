# Ecommerce-application
This project is written purely in Go Language. Gin (Http Web Frame Work) is used in this project. PostgreSQL Database is used to manage the data.
## Framework Used
Gin-Gonic: This whole project is built on Gin frame work. Its is a popular http web frame work. 
```
go get -u github.com/gin-gonic/gin
```
## Database used:
PostgreSQL: PostgreSQL is a powerful, open source object-relational database. The data managment in this project is done using PostgreSQL. ORM tool named GORM is also been used to simplify the forms of queries for better understanding.

```
go get -u gorm.io/gorm
go get -u gorm.io/driver/postgres
```
## External Packages Used
#### Razorpay
For Payment I have used the test case of Razorpay.
```
github.com/razorpay/razorpay-go
```
#### Validator
Package validator implements value validations for structs and individual fields based on tags.
```
github.com/go-playground/validator/v10
```
#### Twilio
The twilio-go helper library lets you write Go code to make HTTP requests to the Twilio API and get the OTP. This is open source library.
```
github.com/twilio/twilio-go/rest/api/v2010
```
#### Gomail
Gomail is a simple and efficient package to send emails. It is well tested and documented.
```
gopkg.in/mail.v2
```
#### JWT 
JSON Web Tokens are an open, industry standard RFC 7519 method for representing claims securely between two parties.
```
github.com/golang-jwt/jwt/v4
```
#### Commands to run project:
```
go run main.go
```

## API Platform Used
API platforms Postman is used to run all the API's Provided by this project

### API Documentation
```
```


## API Specification
You can test the API's using Postman. Use this [postman collection](https://lunar-flare-491559.postman.co/workspace/Team-Workspace~d4586165-a6da-4d40-88d4-f70c842c21ce/collection/25078744-60f0832b-8ffd-4d8a-8c62-7729669956bc?action=share&creator=25078744) to test the API's, all the documentation you needed is provided in the following.<br>
Below is the APIs used in the application and some examples along with it. 

## 👉 Signup as user 
  ### Endpoint :
  ```
  http://43.207.185.37:5000/user/signup
  ```  
  ### Method:
  `POST`
  
  ### Request Body:
  | Parameter     | Type    | Description              |
  |---------------|---------|--------------------------|
  | `firstname`   | string  | First name of the user   |
  | `lastname`    | string  | Last name of the user    |
  | `email`       | string  | Email ID of the user     |
  | `password`    | string  | Password of the user     |
  | `phonenumber` | Intiger | Phone number of the user |
  
  ### Example Request:
  ```
   POST http://43.207.185.37/user/signup 
  -H "Content-Type: application/json" 
  -d '{
      "firstname" : "Tony",
      "lastname": "Stark",
      "email" : "tony@yopmail.com",
      "password" : "12345",
      "phonenumber": 9087867817
  }'
  ```
  
  ### Success Response:
  HTTP Code: `200 OK`
  
  ```
  {
    "message": "Go to /signup/otpvalidate"
  }
  ```
  
## 👉 To verify the otp
  ### Endpoint :
  ```
  http://localhost:5000/user/signup/otpvalidate
  ```  
  ### Method:
  `POST`

   
  ### Request Body:
 ### Request Body:
  | Parameter   | Type     | Description       |
  |-------------|----------|-------------------|
  | `otp`       | Intiger  | Otp of the user   |
  | `email`     | string   | Email of the user |
  
 
  ### Example Request:
  ```
   POST  http://43.207.185.37:5000/user/signup/otpvalidate\
  -H "Content-Type: application/json" \
  -d '{
         "otp" : "1904",
         "email" : "tony@yopmail.com"
      }'
  ``` 
  ### Success Response:
  HTTP Code: `200 OK`
  
  ```
  {
    "Message": "New User Successfully Registered"
  }
  ```
  
## 👉 To login as a user
  ### Endpoint :
  ```
    http://43.207.185.37:5000/user/login
  ```  
  ### Method:
  `POST`
 
   ### Request Body:
  | Parameter     | Type    | Description              |
  |---------------|---------|--------------------------|
  | `email`       | string  | Email ID of the user     |
  | `password`    | string  | Password of the user     |
  
 ### Example Request:
  ```
   POST  http://43.207.185.37:5000/user/signup/otpvalidate\
  -H "Content-Type: application/json" \
  -d '{
        "email" : "tony@yopmail.com",
        "password" :"12345"
      }'
  ``` 

  ### Success Response:
  HTTP Code: `200 OK`
  
  ```
  {
    "message": "User login successfully"
  }
  ```
 