# clang-libtooling-comp

## TODO
* 下記のような途中からの補完に対応する
```
clang::As{_cursor_}
```

```
cd data
./parse.sh clang

cat > list.txt <<EOF
clang::SourceLocation
clang::SourceManager
# clang::SourceMgr class?
clang::SourceRange
clang::VarDecl
EOF

cat clang.cpp.comp | grep Decl | awk '$2==$4{print "clang::"$2}' >> list.txt
cat clang.cpp.comp | grep Expr | awk '$2==$4{print "clang::"$2}' >> list.txt
cat clang.cpp.comp | grep Stmt | awk '$2==$4{print "clang::"$2}' >> list.txt
cat clang.cpp.comp | grep AST | awk '$2==$4{print "clang::"$2}' >> list.txt
cat list.txt | sort | uniq | grep -v "#" > parse-list.txt

cat parse-list.txt | xargs -L 1 ./parse.sh

go run main.go
```
