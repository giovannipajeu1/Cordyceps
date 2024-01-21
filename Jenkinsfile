pipeline {
    agent any
    environment {
        SEMGREP_APP_TOKEN = credentials('SEMGREP_APP_TOKEN')
        GIT_DISCOVERY_ACROSS_FILESYSTEM = true
    }
    stages {
        stage('Git') {
            steps {
                script {
                    // Use checkout to fetch the repository
                    checkout scm
                    // Configure safe directory for Git
                    sh 'git config --global --add safe.directory $(pwd)'
                }
            }
        }
        stage('Semgrep-Scan') {
            steps {
                script {
                    // Pull the Semgrep Docker image
                    sh 'docker pull returntocorp/semgrep'
                    // Run Semgrep with the specified environment variables
                    sh '''
                        docker run \
                        -e SEMGREP_APP_TOKEN=$SEMGREP_APP_TOKEN \
                        -e SEMGREP_REPO_URL=$SEMGREP_REPO_URL \
                        -e SEMGREP_REPO_NAME=$SEMGREP_REPO_NAME \
                        -e SEMGREP_BRANCH=$SEMGREP_BRANCH \
                        -e SEMGREP_COMMIT=$SEMGREP_COMMIT \
                        -e SEMGREP_PR_ID=$SEMGREP_PR_ID \
                        -v "$(pwd):$(pwd)" --workdir $(pwd) \
                        returntocorp/semgrep semgrep ci
                    '''
                }
            }
        }
    }
}
