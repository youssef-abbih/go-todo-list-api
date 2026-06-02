pipeline {
    agent any

    environment {
        ENV         = 'TEST'
        TEST_DB_HOST = 'localhost'
        TEST_DB_PORT = '5432'
        TEST_DB_NAME = 'tododb_test'
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
                        -p 5432:5432 \
                        postgres:15

                    echo "Waiting for PostgreSQL to be ready..."
                    sleep 5
                '''
            }
        }

        stage('Run Tests') {
            steps {
                sh '''
                    export ENV=TEST
                    export TEST_DB_HOST=localhost
                    export TEST_DB_PORT=5432
                    export TEST_DB_NAME=tododb_test
                    go test ./... -v
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