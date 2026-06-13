pipeline {
    agent any

    environment {
        IMAGE_NAME = "url-shortener-app"
        CONTAINER_NAME = "url_shortener_app"
        DB_CONTAINER_NAME = "url_shortener_db"
    }

    stages {

        stage('Checkout') {
            steps {
                git branch: 'main',
                    url: 'https://github.com/abhishekbhonde/url-shortner-to-1m.git'
            }
        }

        stage('Build Docker Image') {
            steps {
                sh '''
                docker build -t $IMAGE_NAME .
                '''
            }
        }

        stage('Stop Existing Containers') {
            steps {
                sh '''
                docker stop $CONTAINER_NAME || true
                docker rm $CONTAINER_NAME || true
                docker stop $DB_CONTAINER_NAME || true
                docker rm $DB_CONTAINER_NAME || true
                '''
            }
        }

        stage('Run with Docker Compose') {
            steps {
                sh '''
                docker compose up -d --build
                '''
            }
        }
    }

    post {
        success {
            echo 'URL Shortener deployed successfully on port 8080!'
        }
        failure {
            echo 'Deployment failed!'
        }
    }
}
