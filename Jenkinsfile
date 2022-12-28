pipeline {
    agent any

    stages {
        stage ('Build') {

            steps {

                git branch: 'main', url: 'https://github.com/vaidyakishor85/go-lambda'

                //Build application 
                sh "go build main.go"
                
            }
        }

        stage ('Deploy/Run') {
            steps {

                //Run application
                sh "nohup go run main.go 2>&1 &"
                
            }
        }
    }
}