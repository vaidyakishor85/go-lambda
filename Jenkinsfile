pipeline {
    agent any

    stages {
        stage ('Build') {

            steps {

                git branch: 'main', url: 'https://github.com/vaidyakishor85/go-lambda'

                //Build application 
                bat "go build main.go"
                
            }
        }

        stage ('Deploy/Run') {
            environment {
                CODECOV_TOKEN = "codecov_token111"
            }
            steps {          

                //Run application
                bat "go run main.go 2>&1 &"
                
            }
        }
    }
}
