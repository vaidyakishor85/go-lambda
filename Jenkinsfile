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
            steps {
               
               environment {
                CODECOV_TOKEN = credentials('codecov_token111')
            }

                //Run application
                bat "go run main.go 2>&1 &"
                
            }
        }
    }
}
