<meta charset="UTF-8">

# DNS blocklist proxy

## About

This is a simple DNS proxy accepting a blocklist made as a test task for a company.

There are two programs: dnsproxy and extra-recursive-resolver.

Documentation is available for both of them in Go comments, suitable for godoc.

## Direct dependencies

- Go.
- Parts of Go's standard library.
- DNS library: <https://github.com/miekg/dns>.
- GNU make.
- Entr command for rerunning commands when files change.

## Make targets

Make check does some test queries if dnsproxy is running. It also performs some checks on extra-recursive-resolver program.

Make configuration.json copies example-configuration.json as the configuration.

Make build builds both programs.

Make clean removes some of the build artifacts.

Make watch provides a simple compile-edit loop using entr.

Make README.html generates this document in HTML form.

## License

MIT.

## Resources used during development

<https://jvns.ca/blog/2022/02/01/a-dns-resolver-in-80-lines-of-go/>

<https://reintech.io/blog/implementing-a-dns-server-in-go/>

## Original task description in Russian

Написать DNS прокси-сервер с поддержкой "черного" списка доменных имен.

1. Для параметров используется конфигурационный файл, считывающийся при запуске сервера;
2. "Черный" список доменных имен находится в конфигурационном файле;
3. Адрес вышестоящего сервера также находится в конфигурационном файле;
4. Сервер принимает запросы DNS-клиентов на стандартном порту;
5. Если запрос содержит доменное имя, включенное в "черный" список, сервер возвращает клиенту ответ, заданный конфигурационным файлом (варианты: not resolved, адрес в локальной сети, ...);
6. Если запрос содержит доменное имя, не входящее в "черный" список, сервер перенаправляет запрос вышестоящему серверу, дожидается ответа и возвращает его клиенту.

Язык разработки: Go. Использование готовых библиотек: без ограничений. Использованный чужой код должен быть помечен соответствующими копирайтами, нарушать авторские права запрещено.

Остальные условия/допущения, не затронутые в тестовом задании - по собственному усмотрению.
