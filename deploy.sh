# !/bin/sh

docker build . -t 'gcr.io/hidden-server-318721/barebone-http-1:latest'
docker push gcr.io/hidden-server-318721/barebone-http-1:latest

kubectl apply -k k8s/