pipeline {
    
    agent any

    tools { go '1.19.2' }

    stages {
        stage('Verify Go is installed') {
            steps {
                sh 'go version'
            }
        }
        
        stage('Build') {
            steps {
                sh 'go build -v ./...'
            }
        }

        stage('Test') {
            steps {
                sh 'go test -v ./...'
            }
        }
    }
}