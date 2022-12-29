package main

import (
	"archive/zip"
	"crypto/x509"
	"encoding/csv"
	"encoding/pem"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"gopkg.in/gomail.v2"
)

func demo(f func()) {
	println("In demo function")

	println(f)

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("ap-south-1"),
		Credentials: credentials.NewStaticCredentials("AKIAXCQMVIUIMAM6TREB", "F/saPbQBos8XtggRU6HnBdYLFFdv6JW4x4NulXm3", ""),
	})

	svc := lambda.New(sess, &aws.Config{Region: aws.String("ap-south-1")})

	result, err := svc.ListFunctions(nil)
	if err != nil {
		fmt.Println("Cannot list functions")
		os.Exit(0)
	}

	println("Result from demo ", result)

}

func main() {
	
	println("Env variable"+$(CODECOV_TOKEN))

	empData := [][]string{
		{"Certificate Name", "Lambda Name", "Issued By", "Valid From", "Expiry"},
	}

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("ap-south-1"),
		Credentials: credentials.NewStaticCredentials("AKIAXCQMVIUIMAM6TREB", "F/saPbQBos8XtggRU6HnBdYLFFdv6JW4x4NulXm3", ""),
	})

	svc := lambda.New(sess, &aws.Config{Region: aws.String("ap-south-1")})

	result, err := svc.ListFunctions(nil)
	if err != nil {
		fmt.Println("Cannot list functions")
		os.Exit(0)
	}

	fmt.Println("Functions:")
	for _, f := range result.Functions {
		fmt.Println("Name: " + aws.StringValue(f.FunctionName))
		fmt.Println("")
		input := &lambda.GetFunctionInput{
			FunctionName: aws.String(aws.StringValue(f.FunctionName)),
		}

		result, err := svc.GetFunction(input)

		if err != nil {
			fmt.Println(err)

		}
		//fmt.Println("URL : " + *result.Code.Location)
		specUrl := *result.Code.Location
		resp, err := http.Get(specUrl)
		if err != nil {
			fmt.Printf("err: %s", err)
		}

		defer resp.Body.Close()
		fmt.Println("status", resp.Status)
		if resp.StatusCode != 200 {
			return
		}

		// Create the file
		out, err := os.Create(aws.StringValue(f.FunctionName) + ".zip")
		if err != nil {
			fmt.Printf("err: %s", err)
		}
		defer out.Close()

		// Write the body to file
		_, err = io.Copy(out, resp.Body)
		fmt.Printf("err: %s", err)

		//fmt.Println(result)

		err1 := unzipSource(aws.StringValue(f.FunctionName)+".zip", "./UnzipFunctions/")
		if err1 != nil {
			log.Fatal(err1)
		}

		// List files
		files, err := ioutil.ReadDir("./UnzipFunctions/" + aws.StringValue(f.FunctionName))
		if err != nil {
			log.Fatal(err)
		}

		for _, file := range files {
			ext1 := filepath.Ext(file.Name())

			if ext1 == ".crt" {
				var fname = aws.StringValue(f.FunctionName)
				fmt.Println("File name with extension .crt:", file.Name())
				var f = "./UnzipFunctions/" + aws.StringValue(f.FunctionName) + "/" + file.Name()
				r, _ := ioutil.ReadFile(f)
				block, _ := pem.Decode(r)

				cert, err := x509.ParseCertificate(block.Bytes)
				if err != nil {
					log.Fatal(err)
				}
				//fmt.Println(cert)
				fmt.Printf("Issuer Name: %s\n", cert.Issuer)
				fmt.Printf("Expiry: %s \n", cert.NotAfter.Format("2006-January-02"))
				fmt.Printf("Common Name: %s \n", cert.Issuer.CommonName)

				details := [][]string{
					{file.Name(), fname, cert.Issuer.CommonName, cert.NotBefore.Format("2006-January-02"), cert.NotAfter.Format("2006-January-02")},
				}
				empData = append(empData, details...)
			}
		}

	}

	csvFile, err := os.Create("lambdaFunctions.csv")

	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}

	csvwriter := csv.NewWriter(csvFile)

	for _, empRow := range empData {
		_ = csvwriter.Write(empRow)
	}
	csvwriter.Flush()
	csvFile.Close()
	//log.Println(empData)

	// Email
	abc := gomail.NewMessage()

	abc.SetHeader("From", "vaidyakishor85@gmail.com")
	abc.SetHeader("To", "vaidyakishor14@gmail.com")

	abc.SetHeader("Subject", "Test Email")

	abc.SetHeader("text/plain", "Test body eamil")

	abc.Attach("lambdaFunctions.csv")

	a := gomail.NewDialer("smtp.gmail.com", 587, "vaidyakishor85@gmail.com", "feeesyiqmzenvuse")

	if err := a.DialAndSend(abc); err != nil {
		fmt.Println(err)
		panic(err)
	}

	log.Println("Script execute success")
}

func unzipSource(source, destination string) error {
	// 1. Open the zip file
	reader, err := zip.OpenReader(source)
	if err != nil {
		return err
	}
	defer reader.Close()

	// 2. Get the absolute destination path
	destination, err = filepath.Abs(destination)
	if err != nil {
		return err
	}

	// 3. Iterate over zip files inside the archive and unzip each of them
	for _, f := range reader.File {
		err := unzipFile(f, destination)
		if err != nil {
			return err
		}
	}

	return nil
}

func unzipFile(f *zip.File, destination string) error {
	// 4. Check if file paths are not vulnerable to Zip Slip
	filePath := filepath.Join(destination, f.Name)
	if !strings.HasPrefix(filePath, filepath.Clean(destination)+string(os.PathSeparator)) {
		return fmt.Errorf("invalid file path: %s", filePath)
	}

	// 5. Create directory tree
	if f.FileInfo().IsDir() {
		if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
			return err
		}
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
		return err
	}

	// 6. Create a destination file for unzipped content
	destinationFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	// 7. Unzip the content of a file and copy it to the destination file
	zippedFile, err := f.Open()
	if err != nil {
		return err
	}
	defer zippedFile.Close()

	if _, err := io.Copy(destinationFile, zippedFile); err != nil {
		return err
	}
	return nil
}
