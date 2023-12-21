.DEFAULT_GOAL := watch

.PHONY : watch
watch :
	: \
		| ( \
			printf %s\\n configuration.json; \
			find . -name '*.go' \
		) \
		| entr make build \
	;

.PHONY : build
build :
	go build
	cd extra-recursive-resolver && go build

configuration.json : example-configuration.json
	cp $^ $@

README.html : README.md
	pandoc -i README.md -o README.html

.PHONY : check
check :
	pgrep dnsproxy

	make -B configuration.json

	test -z "$$(dig +short @localhost github.com)" || exit 1
	test -z "$$(dig +short @localhost alurm.github.io)" || exit 1
	test -z "$$(dig +short @localhost test.nosuchdomain)" || exit 1

	test -n "$$(dig +short @localhost google.com)" || exit 1
	test -n "$$(dig +short @localhost alurm.github.com)" || exit 1
	@ echo dnsproxy ok

	cd extra-recursive-resolver && make check

.PHONY : clean
clean :
	rm -f dnsproxy
	rm -f extra-recursive-resolver/extra-recursive-resolver
	rm -f configuration.json
	rm -f README.html
