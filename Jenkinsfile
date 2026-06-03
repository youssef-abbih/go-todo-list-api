pipeline {
    agent any

    environment {
        ENV              = 'TEST'
        TEST_DB_HOST     = 'host.docker.internal'
        TEST_DB_PORT     = '5432'
        TEST_DB_NAME     = 'tododb_test'
        TEST_DB_USER     = credentials('TEST_DB_USER')
        TEST_DB_PASSWORD = credentials('TEST_DB_PASSWORD')
        JWT_SECRET       = credentials('JWT_SECRET')
    }

    stages {

        stage('Checkout') {
            steps {
                checkout scm
            }
        }

        stage('Start Test Database') {
            steps {
                sh '''
                    docker run -d \
                        --name test-postgres-${BUILD_NUMBER} \
                        -e POSTGRES_USER=${TEST_DB_USER} \
                        -e POSTGRES_PASSWORD=${TEST_DB_PASSWORD} \
                        -e POSTGRES_DB=${TEST_DB_NAME} \
                        -p 5433:5432 \
                        postgres:15

                    echo "Waiting for PostgreSQL to be ready..."
                    sleep 5
                '''
            }
        }

        stage('Run Tests') {
            steps {
                sh '''
                    docker run --rm \
                        --network ${NET_NAME} \
                        -v ${WORKSPACE}:/app \
                        -v ${WORKSPACE}/.go-cache:/go/pkg/mod \
                        -w /app \
                        -e ENV=TEST \
                        -e TEST_DB_HOST=test-postgres-${BUILD_NUMBER} \
                        -e TEST_DB_PORT=5432 \
                        -e TEST_DB_NAME=${TEST_DB_NAME} \
                        -e TEST_DB_USER=${TEST_DB_USER} \
                        -e TEST_DB_PASSWORD=${TEST_DB_PASSWORD} \
                        -e JWT_SECRET=${JWT_SECRET} \
                        golang:1.24 go test ./... -v
                '''
            }
        }
    }

    post {
        always {
            sh '''
                docker stop test-postgres-${BUILD_NUMBER} || true
                docker rm test-postgres-${BUILD_NUMBER} || true
            '''
        }
        success {
            echo 'All tests passed!'
        }
        failure {
            echo 'Tests failed!'
        }
    }
}