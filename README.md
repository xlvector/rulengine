rulengine
=========

Rule Engine

You can try it by following steps

    cd bin
    go build rulengine.go
    #define your rules in rules.txt
    ./rulengine --rules ../rules.txt --data "age=37&gender=male&salary=12345"