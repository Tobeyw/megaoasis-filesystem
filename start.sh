#!/bin/sh
echo shut down existed docker service
echo you env is $1
if [ $1 == "TEST" ]
then
    export RUNTIME="test"
    docker stop megaoasis_filesystem_test
    docker container rm megaoasis_filesystem_test
    docker rmi test_megaoasis_filesystem -f
    docker-compose -p "test" up -d
fi

if [ $1 == "STAGING" ]
then
    export RUNTIME="staging"
    docker stop megaoasis_filesystem_staging
    docker container rm megaoasis_filesystem_staging
    docker rmi staging_megaoasis_filesystem -f
    docker-compose -p "staging" up -d
fi


