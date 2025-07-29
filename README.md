# Golang Dependency Graph

Este projeto implementa um sistema simples de gerenciamento de componentes em Go

## Visão Geral

O sistema permite definir componentes com ciclos de vida (inicialização e encerramento) e gerenciar suas dependências. Os componentes são organizados em um grafo acíclico direcionado (DAG), garantindo que sejam inicializados na ordem correta de dependência.

## Características

- Interface `Lifecycle` para componentes com métodos `Start` e `Stop`
- Definição de componentes com suas dependências
- Gerenciamento automático da ordem de inicialização e encerramento
- Detecção de dependências cíclicas
- Encerramento gracioso com tratamento de sinais

## Estrutura do Projeto

```
golang-dependency-graph/
├── component/
│   ├── lifecycle.go    # Interface Lifecycle
│   ├── component.go    # Definição de componentes
│   └── system.go       # Sistema de gerenciamento
├── examples/
│   └── components.go   # Componentes de exemplo
└── cmd/
    └── demo/
        └── main.go     # Aplicação de demonstração
```

## Como Usar

### Definindo um Componente

```go
// Implementar a interface Lifecycle
type MyComponent struct{}

func (c *MyComponent) Start(ctx component.Context) (component.Lifecycle, error) {
    return c, nil
}

func (c *MyComponent) Stop(ctx component.Context) error {
    return nil
}
```

### Criando um Sistema

```go
compA := component.Define("compA", new(MyComponentA))
compB := component.Define("compB", new(MyComponentB), compA.Key())

components := map[string]*component.Component{
    compA.Key(): compA,
    compB.Key(): compB,
}

system := component.CreateSystem(components)
if err := system.Start(); err != nil {
    log.Fatalf("Erro ao iniciar: %v", err)
}

if err := system.Stop(); err != nil {
    log.Fatalf("Erro ao encerrar: %v", err)
}
```

## Exemplo

O projeto inclui um exemplo completo que demonstra o uso do sistema de componentes:

1. `Config`: Fornece valores de configuração
2. `AppRoutes`: Define rotas HTTP
3. `HttpServer`: Configura e executa um servidor HTTP

Para executar o exemplo:

```bash
go run cmd/demo/main.go
```