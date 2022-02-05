package services

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"

	"github.com/go-services/gomod/config"
	"github.com/go-services/gomod/resphandle"
	"github.com/minio/minio-go"
)

/** This function handles inital verification of request
  * Then pass the file for encryption
**/
func Upload(w http.ResponseWriter, r *http.Request) {
	//r.ParseMultipartForm(1000)
	file, handler, fileerr := r.FormFile("file")
	if fileerr != nil {
		log.Println("Error getting uploaded file: ", fileerr)
		resphandle.HandleResp(w, http.StatusInternalServerError, false, "Error Retrieving the File")
	}

	bucketName := r.FormValue("userId")
	objectName := r.FormValue("objectName")

	if bucketName == "" {
		resphandle.HandleResp(w, http.StatusOK, false, "userId cannot be null")
		return
	}
	if objectName == "" {
		resphandle.HandleResp(w, http.StatusOK, false, "objectName cannot be null")
		return
	}
	minioClient := config.Minioconfig()

	handlefile, staterr := minioClient.StatObject(bucketName, objectName, minio.StatObjectOptions{})
	if staterr != nil {
		log.Println("Error bucket", staterr)
		if staterr.Error() == "Bucket name contains invalid characters" {
			resphandle.HandleResp(w, http.StatusOK, false, staterr.Error())
			return
		}
	}
	if handlefile.Size > 0 {
		log.Println("File " + objectName + " already exists in bucket " + bucketName)
		resphandle.HandleResp(w, http.StatusOK, false, "File "+objectName+" already exists, please rename the file")
		return
	}
	filename := handler.Filename
	filesize := handler.Size
	log.Println("Got Request for Uploading file:", handler.Filename, " For bucketname", bucketName)
	log.Println("File size:", filesize, "bytes")
	response := encryptFileWithKey(w, file, filename, bucketName, filesize)
	if response {
		resphandle.HandleResp(w, http.StatusOK, true, "File uploaded successfully")
		return
	} else {
		resphandle.HandleResp(w, http.StatusInternalServerError, false, "Error uploading file")
		return
	}
}

/** This method is responsible for encrytion of file
  * @param w http.ResponseWriter
  * @param file multipart.File uploaded file
  * @param filename string filename of the uploaded file
  * @param userId string userId will be treated as bucket name for minio server
  * @param filesize int64 length of the uploaded file
  * @return true if successfull and false if failed
**/
func encryptFileWithKey(w http.ResponseWriter, file multipart.File, filename string, userId string, filesize int64) bool {
	log.Println(filename + " was sent for encryption")
	bytefile, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println(err)
	}
	key := []byte(os.Getenv("ENCRYPT_KEY"))

	// generate a new aes cipher using our 32 byte long key
	cont, err := aes.NewCipher(key)
	if err != nil {
		log.Println(err)
	}
	gcm, err := cipher.NewGCM(cont)
	if err != nil {
		log.Println(err)
	}
	nonce := make([]byte, gcm.NonceSize())
	encryptedfile := gcm.Seal(nonce, nonce, bytefile, nil)
	log.Println("File length after encryption ", len(encryptedfile))
	lenfile := int64(len(encryptedfile))
	log.Println("File successfully encripted... Sending it for upload to bucket")
	return uploadEncryptedFile(w, encryptedfile, filename, userId, lenfile)

}

/** This method is responsible for uploading encrypted file
  * @param w http.ResponseWriter
  * @param file []byte encrypted file in byte form
  * @param filename string filename of the uploaded file
  * @param userId string userId will be treated as bucket name for minio server
  * @param filesize int64 length of the uploaded file
  * @return true if successfull and false if failed
**/
func uploadEncryptedFile(w http.ResponseWriter, file []byte, filename string, userId string, fileSize int64) bool {
	log.Println("Encrypted file " + filename + " was sent to upload")
	minioClient := config.Minioconfig()
	if minioClient == nil {
		return false
	}
	bucketName := userId
	location := "us-east-1"
	log.Println("Checking if bucket exists...")
	exists, err := minioClient.BucketExists(bucketName)
	if err != nil {
		log.Println("error querying for bucket", err)
		return false
	}
	if !exists {
		log.Println(bucketName + " Doesn't exists... creating the new bucket")
		if err := minioClient.MakeBucket(bucketName, location); err != nil {
			log.Println("error creating bucket", bucketName, err)
			resphandle.HandleResp(w, http.StatusInternalServerError, true, "Error creating bucket"+err.Error())
			return false
		}
	}

	contentType := "application/octet-stream"
	reader := bytes.NewReader(file)
	log.Println("Uploading file " + filename + " to the bucket " + bucketName)
	_, errupload := minioClient.PutObject(bucketName, filename, reader, fileSize, minio.PutObjectOptions{ContentType: contentType})
	log.Println("Upload successfull...")
	if errupload != nil {
		log.Println(errupload)
	}

	log.Println("Updating info in db..")
	redisClient := config.RedisClient()
	err4 := redisClient.HSet(bucketName, filename, 0).Err()
	if err4 != nil {
		log.Println("Error updating db ", err)
		return false
	}
	log.Println("Updating successfull..")
	return true
}

/** This method is responsible for decrypting and  downloading requested file
  * @param w http.ResponseWriter
  * @param filename string filename of the uploaded file
  * @param bucketname string bucketname will be treated as userId name for minio server
  * @return decryted file if the file exists
**/
func decryptFile(w http.ResponseWriter, bucketName string, filename string) []byte {

	log.Println("Got download request for bucket " + bucketName + " and filename " + filename)

	minioClient := config.Minioconfig()
	_, bucketexistserr := minioClient.BucketExists(bucketName)
	if bucketexistserr != nil {
		log.Println("Error in fetching bucket:", bucketexistserr)
		resphandle.HandleResp(w, http.StatusOK, false, bucketName+" Doesn't exists")
		return nil
	}
	_, staterr := minioClient.StatObject(bucketName, filename, minio.StatObjectOptions{})
	if staterr != nil {
		log.Println("Error in fetching file:", staterr)
		resphandle.HandleResp(w, http.StatusOK, false, filename+" doesn't exist on bucket "+bucketName)
		return nil
	}

	getfile, getfileerr := minioClient.GetObject(bucketName, filename, minio.GetObjectOptions{})
	if getfileerr != nil {
		log.Println("Error in fetching file:", getfileerr)
		resphandle.HandleResp(w, http.StatusInternalServerError, false, "Error Fetching file contact support team")
		return nil
	}
	/** Updating Downloadfilecount before encrypting the file **/
	redisClient := config.RedisClient()
	val, geterr := redisClient.HGet(bucketName, filename).Result()
	if geterr != nil {
		log.Println("Redis error", geterr)
	}
	intVar, interr := strconv.Atoi(val)
	if interr != nil {
		log.Println("Redis error", interr)
	}
	err4 := redisClient.HSet(bucketName, filename, intVar+1).Err()

	if err4 != nil {
		log.Println(err4)
	}

	dat, err6 := io.ReadAll(getfile)
	key := []byte(os.Getenv("ENCRYPT_KEY"))
	if err6 != nil {
		log.Println(string(dat))
	}
	c, err := aes.NewCipher(key)
	if err != nil {
		log.Println(err)
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		log.Println(err)
	}

	nonceSize := gcm.NonceSize()
	if len(dat) < nonceSize {
		log.Println(err)
	}

	nonce, dat := dat[:nonceSize], dat[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, dat, nil)
	if err != nil {
		log.Println(err)
	}

	log.Println("Finished")
	return plaintext

}

/** This method handles the initial request for downloading a file
  * calls decryptfile method for processing
  * retunrn the file or nil
**/
func Download(w http.ResponseWriter, r *http.Request) {

	bucketName := r.FormValue("userId")
	objectName := r.FormValue("objectName")
	if bucketName == "" {
		resphandle.HandleResp(w, http.StatusOK, false, "bucketName cannot be empty")
		return
	}
	if objectName == "" {
		resphandle.HandleResp(w, http.StatusOK, false, "objectName cannot be empty")
		return
	}
	file := decryptFile(w, bucketName, objectName)
	w.Header().Add("Content-Disposition", "attachment ; filename="+objectName)
	w.Write(file)

}

/* -------------------------------------------- Methods for another type of Encription & Decription  -------------------------------- */
/*
func encryptfileByPBKDF(w http.ResponseWriter, file multipart.File, filename string, bucketName string, filesize int64) {
	bytefile, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println(err)
	}
	reader := bytes.NewReader(bytefile)
	encryption := encrypt.DefaultPBKDF([]byte("password"), []byte(bucketName+filename))
	minioClient := config.Minioconfig()
	minioClient.PutObject(bucketName, filename, reader, filesize, minio.PutObjectOptions{ServerSideEncryption: encryption})
	// getfile, err2 := minioClient.GetObject(bucketName, filename, minio.GetObjectOptions{ServerSideEncryption: encryption})
	// //getfile, err2 := minioClient.GetObject(bucketName, filename, minio.GetObjectOptions{})
	// if err2 != nil {
	// 	log.Fatalln(err)
	// }

	// dat, err6 := io.ReadAll(getfile)
	// if err6 != nil {
	// 	fmt.Println(string(dat))
	// }
	// fmt.Println(string(dat))
	//decryptFile(bucketName, filename)
	//dataJson := &model.FileData{ObjectName: filename, DownloadCount: 0} //model.InitFileData(filename, 0)
	//dat, err := json.Marshal(dataJson)
}
*/
