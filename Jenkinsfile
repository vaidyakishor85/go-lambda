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

                //Run application
                bat "nohup go run main.go 2>&1 &"
                
            }
        }
    }
}