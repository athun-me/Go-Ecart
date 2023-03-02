## Ecom
The application is a Ecommerce application and admin can add product and user can buy that. All the APIs provided here is to alter or modify these datas. The database used here is Postgres and the server is hosted in the AWS EC2 instance. 

## API Specification
You can test the API's using Postman. Use this [postman collection](https://lunar-flare-491559.postman.co/workspace/Team-Workspace~d4586165-a6da-4d40-88d4-f70c842c21ce/collection/25078744-60f0832b-8ffd-4d8a-8c62-7729669956bc?action=share&creator=25078744) to test the API's, all the documentation you needed is provided in the following.<br>
Below is the APIs used in the application and some examples along with it. 

## 👉 Signup as user 
  ### Endpoint :
  ```
  http://54.95.224.42
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
   POST http://54.95.224.42/user/signup \
  -H "Content-Type: application/json" \
  -d '{
      "firstname" : "Athun",
      "lastname": "lal",
      "email" : "Rahulcp5@gmail.com",
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
  http://54.95.224.42/user/signup/otpvalidate
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
   POST  http://54.95.224.42/user/signup/otpvalidate\
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
    http://54.95.224.42/user/login
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
   POST  http://54.95.224.42/user/signup/otpvalidate\
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
  
## 👉 Product add into the cart
  ### Endpoint :
  ```
  http://54.95.224.42/user/profile/addtocart
  ```  
  ### Method:
  `POST`
 
  ### Request Body:
  | Parameter | Type | Description |
  |-----------|------|-------------|
  | `Product_id` | Integer | Id of the product |
  | `Quantity` | string | Quantity of the products |
  
  ### Example Request:
  ```
   POST http://54.95.224.42/user/profile/addtocart \
  -H "Content-Type: application/json" \
  -d '{
      	"Product_id": 3,
	      "Quantity": 1
      }'
  ```
  
  ### Success Response:
  HTTP Code: `200 OK`
  
  ```
  {
    "error": false,
    "message": "Successfully updated the user data"
  } 
  ```
  
## 👉 Checkout the cart products
  ### Endpoint :
  ```
     http://54.95.224.42/user/cart/checkout
  ```  
  ### Method:
  `GET`
  
  ### Example Request:
  ```
  GET http://54.95.224.42/user/cart/checkout
  ```
  
  ### Success Response:
  HTTP Code: `200 OK`
  
  ```
  {
    "Cart Items": [
        {
            "Product_name": "Realme X2",
            "Quantity": 1,
            "Total_price": 20000,
            "Image": "808cbbd6-851c-41b8-a1b2-46f0e34e9bf4.webp",
            "Price": "20000"
        }
    ]
}{
    "Default Address of user": [
        {
            "Name": "",
            "Phoneno": "9605018636",
            "Houseno": "12",
            "Area": "blore",
            "Landmark": "HSR",
            "City": "bangelur",
            "Pincode": "671531",
            "District": "bangelur",
            "State": "karnataka",
            "Country": "india"
        }
    ]
}{
    "Total Price": 20000
}
 lse,
    "message": "Successfully dropped the collection"
  }
  ```


## 👉 To buy the product using COD
  ### Endpoint :
  ```
  http://54.95.224.42/user/payment/cashOnDelivery
  ```  
  ### Method:
  `GET`
  
  ### Example Request:
  ```
   POST http://54.95.224.42/user/payment/cashOnDelivery \
  ```
 ### Success Response:
  HTTP Code: `200 OK`

  ```
 {
    "Message": "Payment Method COD",
    "Status": "True"
 }
 {
    "Message": "Oder Added succesfully"
}
  ```