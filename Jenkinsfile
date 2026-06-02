pipeline {
    agent any

    environment {
        ENV              = 'TEST'
        TEST_DB_PORT     = '5432'
        TEST_DB_NAME     = 'tododb_test'
        TEST_DB_USER     = credentials('TEST_DB_USER')
        TEST_DB_PASSWORD = credentials('TEST_DB_PASSWORD')
        JWT_SECRET       = credentials('JWT_SECRET')
        // We define a unique network name for this build
        NET_NAME         = "todo-net-${BUILD_NUMBER}" 
    }

    stages {
        stage('Checkout') {
            steps {
                checkout scm
            }
        }

        stage('Setup Network & DB') {
            steps {
                sh '''
                    # 1. Create a dedicated isolated network
                    docker network create ${NET_NAME}

                    # 2. Start Postgres attached to that network
                    # Note: We name the container deterministically
                    docker run -d \
                        --name test-postgres-${BUILD_NUMBER} \
                        --network ${NET_NAME} \
                        -e POSTGRES_USER=${TEST_DB_USER} \
                        -e POSTGRES_PASSWORD=${TEST_DB_PASSWORD} \
                        -e POSTGRES_DB=${TEST_DB_NAME} \
                        postgres:15

                    echo "Waiting for PostgreSQL to be ready..."
                    until docker exec test-postgres-${BUILD_NUMBER} pg_isready -U ${TEST_DB_USER} -d ${TEST_DB_NAME}; do
                        echo "Postgres is still starting up..."
                        sleep 1
                    done
                '''
            }
        }

        stage('Run Tests') {
            steps {
                // If your Jenkins agent has the Docker CLI, we can spin up a lightweight Go container 
                // on the same network to run the tests cleanly.
                sh '''
                    docker run --rm \
                        --network ${NET_NAME} \
                        -v ${WORKSPACE}:/app \
                        -w /app \
                        -e ENV=TEST \
                        -e TEST_DB_HOST=test-postgres-${BUILD_NUMBER} \
                        -e TEST_DB_PORT=5432 \
                        -e TEST_DB_NAME=${TEST_DB_NAME} \
                        -e TEST_DB_USER=${TEST_DB_USER} \
                        -e TEST_DB_PASSWORD=${TEST_DB_PASSWORD} \
                        golang:1.24 go test ./... -v
                '''
            }
        }
    }

    post {
        always {
            sh '''
                # Clean up everything thoroughly so we don't leak resources
                docker stop test-postgres-${BUILD_NUMBER} || true
                docker rm test-postgres-${BUILD_NUMBER} || true
                docker network rm ${NET_NAME} || true
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