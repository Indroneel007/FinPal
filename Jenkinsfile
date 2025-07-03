pipeline {
    agent any

    tools {
        git 'DefaultGit' // match the name you configured
    }

    stages {
                
        stage('Build') {
            steps {
                bat 'go build -o bin/app.exe .'
            }
        }

        stage('Test') {
            steps {
                bat 'make test'
            }
        }
    }
    

    post {
        always {
            echo "CI pipeline finished (test + optional build)."
        }
    }
}
