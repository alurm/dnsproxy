.PHONY : check
check :
	./extra-recursive-resolver example.com >/dev/null || exit 1
	./extra-recursive-resolver example.com. >/dev/null || exit 1
	./extra-recursive-resolver test.nosuchdomain >/dev/null && exit 1 || :
	./extra-recursive-resolver test.nosuchdomain. >/dev/null && exit 1 || :
	@ echo extra-recursive-resolver ok
