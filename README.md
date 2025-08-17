# Desafio Fullcycle - Deploy no Google Cloud Run

## Objetivo
Desenvolver um sistema em Go que receba um CEP, identifica a cidade e retorna o clima atual (temperatura em graus Celsius, Fahrenheit e Kelvin). Esse sistema deverá ser publicado no Google Cloud Run.

## Requisitos

- O sistema deve receber um CEP válido de 8 dígitos
- O sistema deve realizar a pesquisa do CEP e encontrar o nome da localização, a partir disso, deverá retornar as temperaturas e formatá-las em: Celsius, Fahrenheit, Kelvin.
- O sistema deve responder adequadamente nos seguintes cenários:
	- **Sucesso:**
		- Código HTTP: 200
		- Response Body: `{ "temp_C": 28.5, "temp_F": 28.5, "temp_K": 28.5 }`
	- **CEP inválido (formato incorreto):**
		- Código HTTP: 422
		- Mensagem: `invalid zipcode`
	- **CEP não encontrado:**
		- Código HTTP: 404
		- Mensagem: `can not find zipcode`
- O deploy deve ser realizado no Google Cloud Run.

### Dicas
- Utilize a API [viaCEP](https://viacep.com.br/) para encontrar a localização.
- Utilize a API [WeatherAPI](https://www.weatherapi.com/) para consultar as temperaturas.
- Conversão de Celsius para Fahrenheit: `F = C * 1,8 + 32`
- Conversão de Celsius para Kelvin: `K = C + 273`

### Entrega
- Código-fonte completo da implementação.
- Testes automatizados demonstrando o funcionamento.
- Utilize docker/docker-compose para facilitar os testes da aplicação.
- Deploy realizado no Google Cloud Run (free tier) e endereço ativo para acesso.

---

## Como clonar e rodar a aplicação

### 1. Clone o repositório

```bash
git clone https://github.com/fhsmendes/deploy-cloud-run.git
cd deploy-cloud-run
```

### 2. Configure as variáveis de ambiente

Crie um arquivo `.env` ou exporte as variáveis necessárias, por exemplo:

```bash
export APIKeyWeather=SEU_API_KEY_WEATHERAPI
```

### 3. Build e execute com Docker

```bash
docker build -t deploy-cloud-run .
docker run -p 8080:8080 --env APIKeyWeather=SEU_API_KEY_WEATHERAPI deploy-cloud-run
```

### 4. Rodando os testes

```bash
go test -v -cover ./...
```

### 5. Teste via cloud RUN

```

Você pode testar a aplicação com o seguinte comando, substituindo o CEP se desejar:

```bash
curl "https://deploy-cloud-run-237512667039.us-central1.run.app/temperature?cep=01001000"