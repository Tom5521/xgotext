test:
  go test -v $(dirname $(find . -name "*_test.go"))
benchmark:
  go test -bench=. $(dirname $(find . -name "*_benchmark_test.go"))
bench path:
  go test -bench=. {{path}}
gen-uml:
  goplantuml ./pkg/go/parse > ./pkg/go/parse/structure.puml
  goplantuml ./pkg/po/compiler/ > ./pkg/po/compiler/structure.puml
  goplantuml ./pkg/po/ > ./pkg/po/structure.puml
  goplantuml ./pkg/po/parse/ > ./pkg/po/parse/structure.puml

  plantuml -theme spacelab ./pkg/po/compiler/structure.puml
  plantuml -theme spacelab ./pkg/po/structure.puml
  plantuml -theme spacelab ./pkg/go/parse/structure.puml
  plantuml -theme spacelab ./pkg/po/parse/structure.puml
clean:
  rm -rf $(find . -name "*.po") \
  $(find . -name "*.mo") \
  $(find . -name "*.log")
gen-diff:
  git diff --staged > diff.log
