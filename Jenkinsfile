pipeline {
    agent any

    environment {
        ENV              = 'TEST'
        TEST_DB_PORT     = '5432'
        TEST_DB_NAME     = 'tododb_test'
        TEST_DB_USER     = credentials('TEST_DB_USER')
        TEST_DB_PASSWORD = credentials('TEST_DB_PASSWORD')
        JWT_SECRET       = credentials('JWT_SECRET')
        NET_NAME         = "todo-net-${BUILD_NUMBER}"
        HOST_WORKSPACE   = "/var/lib/docker/volumes/jenkins_home/_data/workspace/go-todo-pipeline"
        APP_NAME         = "go-todo-list-api"
        NEXUS_URL        = 'localhost:5000'
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
                    docker network create ${NET_NAME}

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
                sh '''
                    docker run --rm \
                        --network ${NET_NAME} \
                        -v ${HOST_WORKSPACE}:/app \
                        -v ${HOST_WORKSPACE}/.go-cache:/go/pkg/mod \
                        -w /app \
                        -e ENV=TEST \
                        -e TEST_DB_HOST=test-postgres-${BUILD_NUMBER} \
                        -e TEST_DB_PORT=5432 \
                        -e TEST_DB_NAME=${TEST_DB_NAME} \
                        -e TEST_DB_USER=${TEST_DB_USER} \
                        -e TEST_DB_PASSWORD=${TEST_DB_PASSWORD} \
                        -e JWT_SECRET=${JWT_SECRET} \
                        golang:1.24 go test ./... -v -coverprofile=coverage.out
                '''
            }
        }

        stage('Static Analysis & Security (SAST)'){
            steps{
                withCredentials([string(credentialsId: 'SONAR_TOKEN', variable: 'SONAR_TOKEN')]) {
                    sh "sonar-scanner -Dsonar.login=\$'' \
                    -Dsonar.host.url=https://sonarcloud.io/\
                    -Dsonar.projectName=youssef-abbih/go-todo-list-api \
                    -Dsonar.organization=https://sonarcloud.io/organizations/youssef-abbih/ \
                    -Dsonar.projectKey=youssef-abbih_go-todo-list-api \
                    -Dsonar.projectName=go-todo-list-api \
                    -Dsonar.sources=. \
                    -Dsonar.exclusions=docs/** \
                    -Dsonar.go.coverage.reportPaths=coverage.out \
                    -Dsonar.go.version=1.24 "
                }
            }
        }

        stage('Build docker image') {
            steps {  
                sh 'docker build -t $APP_NAME:latest .'
            }
        }

        stage('Scan Docker Image') {
            steps {
                sh '''
                    docker run --rm \
                        -v /var/run/docker.sock:/var/run/docker.sock \
                        aquasec/trivy:latest image \
                        --exit-code 1 \
                        --severity HIGH,CRITICAL \
                        --no-progress \
                        go-todo-api:latest
                '''
            }
        }

        stage('Push to Nexus') {
            steps {
                withCredentials([usernamePassword(credentialsId: 'NEXUS_CREDENTIALS', usernameVariable: 'NEXUS_USER', passwordVariable: 'NEXUS_PASS')]) {
                    sh '''
                        echo $NEXUS_PASS | docker login localhost:5000 --username $NEXUS_USER --password-stdin
                        docker tag go-todo-api:latest localhost:5000/go-todo-api:latest
                        docker push localhost:5000/go-todo-api:latest
                    '''
                }
            }
        }
    }

    post {
        always {
            sh '''
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