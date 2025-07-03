pipeline {
    agent any

    tools {
        git 'DefaultGit' // match the name you configured
    }


        stage('Build') {
            steps {
                bat 'make build'
            }
        }

        stage('Run Tests') {
            steps {
                bat 'make test'
            }
        }
    

    post {
        always {
            echo "CI pipeline finished (test + optional build)."
        }
    }
}
