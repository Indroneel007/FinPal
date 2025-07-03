pipeline {
    agent any

    tools {
        git 'DefaultGit' // match the name you configured
    }

    stages {
                
        stage('Build') {
            steps {
                bat 'make build'
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
