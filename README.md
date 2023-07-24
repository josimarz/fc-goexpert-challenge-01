# Pós Go Expert

## Desafio Client-Server API

### Iniciando o servidor

Para iniciar o servidor, através do terminal de comandos acesse o diretório raiz deste projeto e execute:

```sh
$ go run server/main.go
```

Se tudo correr bem, você verá a mensagem "Server listening on port 8080".

O servidor está preparado para escutar requisições no *endpoint* `/cotacao`. Sempre que o servidor receber uma requisição nesse *endpoint*, as seguintes ações serão executadas:

* Consultar a API de cotação do Dólar.
* Salvar a resposta da API no banco de dados.
* Retornar o valor da cotação para o cliente no formato JSON.

O formato de retorno do *endpoint* é:

```json
{
    "bid": "{valor}"
}
```

onde `{valor}` é o valor da cotação atual.

### Consultando o banco de dados

Para consultar o banco de dados, através do terminal de comandos acesse o diretório raiz deste projeto e execute (requer o SQLite 3 instalado):

```sh
$ sqlite3 db.sqlite3
```

Após acessar o banco de dados, execute a seguinte instrução SQL para visualizar o histórico de consultas:

```sql
select * from quotation;
```

### Executando o cliente

Para executar o cliente, através do terminal de comandos acesse o diretório raiz deste projeto e execute:

```sh
$ go run client/main.go
```

Se tudo correr bem, um arquivo denominado `cotacao.txt` será criado na raiz do projeto. Esse arquivo contém a cotação do Dólar que foi informada pelo servidor.