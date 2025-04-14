# Gateway de Pagamento - Imersão Full Cycle

## Sobre o Projeto

Este projeto foi desenvolvido durante a [Imersão Full Stack & Full Cycle](https://imersao.fullcycle.com.br/evento/), onde construímos um Gateway de Pagamento completo utilizando arquitetura de microsserviços.

O objetivo é demonstrar a construção de um sistema distribuído moderno, com separação de responsabilidades, comunicação assíncrona e análise de fraudes em tempo real.

## Arquitetura

[Visualize a arquitetura completa aqui](https://link.excalidraw.com/readonly/Nrz6WjyTrn7IY8ZkrZHy)

### Componentes do Sistema

- **Frontend (Next.js)**
  - Interface do usuário para gerenciamento de contas e processamento de pagamentos
  - Desenvolvido com Next.js para garantir performance e boa experiência do usuário

- **Gateway (Go)**
  - Sistema principal de processamento de pagamentos
  - Gerencia contas, transações e coordena o fluxo de pagamentos
  - Publica eventos de transação no Kafka para análise de fraude

- **Apache Kafka**
  - Responsável pela comunicação assíncrona entre API Gateway e Antifraude
  - Garante confiabilidade na troca de mensagens entre os serviços
  - Tópicos específicos para transações e resultados de análise

- **Antifraude (Nest.js)**
  - Consome eventos de transação do Kafka
  - Realiza análise em tempo real para identificar possíveis fraudes
  - Publica resultados da análise de volta no Kafka

## Fluxo de Comunicação

1. Frontend realiza requisições para a API Gateway via REST
2. Gateway processa as requisições e publica eventos de transação no Kafka
3. Serviço Antifraude consome os eventos e realiza análise em tempo real
4. Resultados das análises são publicados de volta no Kafka
5. Gateway consome os resultados e finaliza o processamento da transação

## Ordem de Execução dos Serviços

Para executar o projeto completo, os serviços devem ser iniciados na seguinte ordem:

1. **API Gateway (Go)** - Deve ser executado primeiro pois configura a rede Docker
2. **Serviço Antifraude (Nest.js)** - Depende do Kafka configurado pelo Gateway
3. **Frontend (Next.js)** - Interface de usuário que se comunica com a API Gateway

## Instruções Detalhadas

Cada componente do sistema possui instruções específicas de instalação e configuração em suas respectivas pastas:

- **API Gateway**: Consulte o README na pasta `/go-gateway`
- **Serviço Antifraude**: Consulte o README na pasta `/nestjs-antifraude` 
- **Frontend**: Consulte o README na pasta `/next-frontend`

> **Importante**: É fundamental seguir a ordem de execução mencionada acima, pois cada serviço depende dos anteriores para funcionar corretamente.

## Pré-requisitos Gerais

Para executar todos os componentes do projeto, você precisará ter instalado:

- [Docker](https://www.docker.com/get-started)
- [Docker Compose](https://docs.docker.com/compose/install/)
- Git

## Regras de Negócio Importantes

- Transações acima de R$ 10.000 são automaticamente enviadas para análise e ficam com status "pendente"
- Transações menores são processadas imediatamente
- A interface mostra status diferenciados por cores: verde (aprovado), amarelo (pendente), vermelho (rejeitado)