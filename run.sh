#!/bin/bash
read -p "Would you like to run the concurrent test or nonconcurrent test (c/n)? " answer
case ${answer:0:1} in
    c|C )
        SCRIPT=concurrent.go 
	;;

    n|N )
        SCRIPT=nonconcurrent.go
    ;;
esac

read -p "Would you like to pull files over http or locally (h/l)? " answer
case ${answer:0:1} in
    h|H )
        FILES="http://localhost:8000/data/file1.csv \
                http://localhost:8000/data/file2.csv \
                http://localhost:8000/data/file3.csv \
                http://localhost:8000/data/file4.csv \
                http://localhost:8000/data/file5.csv \
                http://localhost:8000/data/file6_bad.csv \
                http://localhost:8000/data/file7_bad.csv \
                http://localhost:8000/data/file9_bad.csv \
                http://localhost:8000/data/file10_bad.csv" 
	;;

    l|L )
       FILES="file://data/file1.csv \
                file://data/file2.csv \
                file://data/file3.csv \
                file://data/file4.csv \
                file://data/file5.csv \
                file://data/file6_bad.csv \
                file://data/file7_bad.csv \
                file://data/file9_bad.csv \
                file://data/file10_bad.csv"
    ;;
esac

python3 -m http.server &
go run $SCRIPT $FILES
pkill -f http.server 