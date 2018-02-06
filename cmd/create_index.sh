#!/bin/zsh

cd testingtesting
go run main.go -d /tmp/reuters_xml --classy --classless -o /tmp/index_full.gob.gz
go run main.go -d /tmp/reuters_xml --classy -o /tmp/index.gob.gz
cd ../k_means_index
go run main.go --i /tmp/index.gob.gz -k ${1} -s /tmp/kmeans.gob.gz
go run main.go --i /tmp/index_full.gob.gz -k ${1} -s /tmp/kmeans_full.gob.gz
