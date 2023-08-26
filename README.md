# A Fast API

Esse repositório é um exemplo de como eu penso sobre a arquitetura de um projeto que tem como objetivo primário lidar com muitas requisições concorrentemente.

A ideia veio a partir da [Rinha de Backend](https://github.com/zanfranceschi/rinha-de-backend-2023-q3)

## Fazendo funcionar

A primeira coisa antes de começarmos a aplicarmos estratégias de redução de latência é fazer o código funcionar corretamente.

[Esse commit](https://github.com/mauricioabreu/a-fast-api/commit/cfaad160eb28fd7910c3d2a7046e94e36112a01c) é nosso ponto de partida para escolher como vamos:

* Medir a performance da API (latência, disponibilidade, etc)
* Melhorar a performance da API
