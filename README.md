# EncryptFileUpload_GO

This project has the following project structure

    config - Holds the configuration part like connection with minio play server and connection of redis client
             Singleton pattern is used here to create the client only once.  
    model - model holds the normal struct for error response.  
    resphandle - Holds function to send responses.  
    services - This contains the encrytion/decryption and also upload/download of the file logic.  
    main.go - contains endpoint.  
    main_test.go - contains unit tests.  
    Dockerfile - contains code to containarize the application.  

## Pre-Setup before running application
Step 1. Download & install golang :- https://go.dev/dl/  
Step 2. Download and install Redis https://redis.io/download  

 #### Check if Go installed correctly  
  Run command ``$ go version`` 
  If version appears then go is installed correctly.
  
  ![image](https://user-images.githubusercontent.com/6186495/152635648-0e548475-1f7b-4d81-99cc-098e2cd94e78.png)

 #### Check if redis is installed(if in windows) open cmd 
  Run command ``$ redis-cli ``  
  Then run ``ping`` and you will get response ``PONG``  
  
  ![image](https://user-images.githubusercontent.com/6186495/152635734-64685926-7abb-4e78-8153-e50aff429b40.png)

  Redis should be running on 127.0.0.1:6379

## Running project in local 
--------------------------------
Step 1. cd into the go-services-challenge directory ``$ cd go-services-challenge directory``  
Step 2. run command ``$ go mod download``  
Step 3. run command ``$ go run main.go``  
After running above commands the cmd will look like below

![image](https://user-images.githubusercontent.com/6186495/152636431-ac5d6ec7-3a4d-49c3-9c37-fbdce360ab3f.png)


Curl command to upload a file to end point /upload
>curl --location --request PUT localhost:3000/upload --form file=@"<file_path>" --form userId="<userId>" --form objectName="<file_name>"

for example: to upload a file kk.kpg from the file path C:/Users/eleay/Downloads/kk.jpg and objectname as kk.jpg below curl command is used

> curl --location --request PUT localhost:3000/upload --form file=@"C:/Users/eleay/Downloads/kk.jpg" --form userId="ayudxt" --form objectName="kk.jpg"

response will look like this:- 
```json
 {
  "Success": true,
  "Description": "File uploaded successfully"
}
```
Browse to https://play.minio.io:9443/login  
Put username: ``Q3AM3UQ867SPQQA43P2F``  
And password: ``zuf+tfteSlswRu7BJ86wekitnifILbZam1KYY3TG``  
    
Click login
    
Then search the userId(Bucket name) and click on browse and see the uploaded file
    
![image](https://user-images.githubusercontent.com/6186495/152636733-713e8308-3d73-4c5d-84e6-3fdf8151ab66.png)

Download the file from browser and test for the encryption [File will not open since it's the encrypted file]  

    

Curl command to download a file endpoint /download  

> curl --location --request GET localhost:3000/download --form userId="<bucket_name>" --form objectName="<file_name>" -o <file_name>
  
Example: To download file kk.jpg from bucket ayudxt the curl command will look like the below  
  
> curl --location --request GET localhost:3000/download --form userId="ayudxt" --form objectName="kk.jpg"  -o kk.jpg  
  
output: 
  ![image](https://user-images.githubusercontent.com/6186495/152636045-3f4d1330-665e-4692-a0aa-fee07662f227.png)
 
 Now browse to the directory and open the downloaded file and Test

Checking download count in redis  
Step 1.	Go to ``redis-cli``  
Step 2.	Run command ``$ hmget <userId>  <filename/objectname>``  
Example for checking downloads of kk.jpg  
``$ hmget ayudxt kk.jpg``  
  
![image](https://user-images.githubusercontent.com/6186495/152636143-a852916c-4e29-46db-91c2-6d1fff430f09.png)
  
 Download count is 2, since the file is downloaded 2 times

  
