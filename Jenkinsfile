pipeline {
    agent any

    tools {
        git 'DefaultGit' // match the name you configured
    }

    stages {
        stage('Build') {
            steps {
                sh 'mkdir -p bin'
                sh 'go build -o bin/app .'
            }
        }

        stage('Test') {
            steps {
                sh 'make test'
            }
        }
    }

    post {
        always {
            echo "CI pipeline finished (test + optional build)."
        }
    }
}

