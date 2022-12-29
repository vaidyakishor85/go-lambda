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
                def rootDir = pwd()
                config = [
                    host    : '0.0.0.0',
                     user    : 'user1',
                     password: 'pass'
]

                load("main.go").demo(config)

                //Run application
                //bat "go run main.go 2>&1 &"
                
            }
        }
    }
}
