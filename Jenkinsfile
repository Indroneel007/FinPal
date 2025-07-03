pipeline {
    agent any

    tools {
        git 'DefaultGit' // match the name you configured
    }

    stages {
        stage('Debug Environment') {
            steps {
                sh '''
                    echo "Go Version:"
                    go version

                    echo "Environment Variables:"
                    printenv

                    echo "Working Directory Contents:"
                    ls -al
                '''
            }
        }

        stage('Checkout') {
            steps {
                checkout scm
            }
        }

        stage('Build') {
            steps {
                sh 'mkdir -p bin'
                sh 'go build -o bin/app .'
            }
        }

        stage('Test') {
            steps {
                sh 'go mod tidy'
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

