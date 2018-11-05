# university-api

O projeto tem 2 TODOs para o class. Falta verificar se o professor existe antes de criar a turma e terminar os testes.
Deixei isso como exercício, mas pode enviar um email para alyssonggs@hotmail.com em casos de dúvida.

## Collection do postman para testes
https://www.getpostman.com/collections/bca6043d43e55d6ab746

## Possíveis upgrades no projeto

### Trabalhar com um banco de dados relacional ao invés de memória
Para aplicações de mercado, é útil salvar os dados em um banco de dados tal como MySQL, PostgresSQL ou OracleDB.

Eu particularmente uso esta lib: https://github.com/jinzhu/gorm.

A documentação ensina como criar uma struct para trabalhar com um banco de dados de verdade e os primeiros passos de consulta e operações neles.


### Deixar as URLs mais amigáveis
Seria interessante substituir a passagem dos parâmetros na url de localhost:8080/student?code=123 para localhost:8080/student/123.

Para isso, recomendo a lib: https://github.com/gorilla/mux

Ela retira a necessidade da função Handle encontrada nos arquivos class, professor e student do pacote handle, além de deixar a URL bem melhor de ser escrita.

PS.: As duas libs podem ser baixadas executando os seguintes comando no terminal(LINUX)
go get github.com/gorilla/mux 
go get github.com/jinzhu/gorm