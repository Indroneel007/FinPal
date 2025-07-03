pipeline {
    agent any

    stages {
        stage('Tidy Modules') {
            steps {
                bat 'make tidy'
            }
        }

        stage('Run Tests') {
            steps {
                bat 'make test'
            }
        }

        stage('Build (optional)') {
            when {
                expression { fileExists('main.go') }
            }
            steps {
                bat 'make build'
            }
        }
    }

    post {
        always {
            echo "CI pipeline finished (test + optional build)."
        }
    }
}
