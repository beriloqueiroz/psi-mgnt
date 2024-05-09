# Plain

## objetivo

- O objetivo desse projeto é registrar sessões/atendimentos de forma a conseguir consultar posteriormente as notas da sessão e o valor, bem como o paciente.

## MVP

### casos de uso

- Cadastro de sessão/atendimento
  - Nome e Sobrenome do paciente (listar as opções com mesmo nome e cadastrar caso não exista)
  - incluir valor
  - incluir nota
  - Data
  - Data pagamento
  - Tempo
- Deletar uma sessão/atendimento

- Relatórios
  - sessões por paciente em um período de pagamento
    - com valor total
    - tempo total
    - média de tempo por sessão
  - sessões por paciente em um período de sessão
    - com valor total
    - tempo total
    - média de tempo por sessão
  - notas por paciente em um período de sessão em ordem de data de sessão
  - sessões em um período
    - data de pagamento
      - com valor total
      - tempo total
      - média de tempo por sessão
    - data da sessão
      - com valor total
      - tempo total
      - média de tempo por sessão

### regras de negócio

- não pode haver nome do paciente repetido
- a sessão não pode ter tempo <0
- a nota tem que ser obrigatória
