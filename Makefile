.DEFAULT_GOAL := watch

.PHONY : watch
watch :
	: \
		| printf %s\\n configuration.json main.go \
		| entr sh -c 'go build && ./dnsproxy' \
	;

configuration.json : example-configuration.json
	cp $^ $@

README.html : README.md
	pandoc -i README.md -o README.html
