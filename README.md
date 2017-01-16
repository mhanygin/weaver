# weaver
Процесс-оркестратор
Основные сущности:
 - сенсор
 - триггер
 - шина передачи экземпляров триггеров
 - правило
 - задача
 - шина передачи результатов работы задач
 - оркестратор
 
Сенсоры описаны в пакете `sensors` и имеют интерфейс:
```
type Sensor interface {
   Setup() error
   Run() error
}
```

Сенсоры запускаются оркестратором один раз. При возникновении события сенсор может испустить триггер через интерфейс:

```
type TriggerExchanger interface {
   Init() error
   Dispatch(trigger Trigger) error
   Reception() (Trigger, error)
}
```

Получить экземпляр интерфейса можно через вызов `sensor.GetExchanger`.
Оркестратор обрабатывает полученные триггеры, согласно секции `activeFlows` конфигурационного файла:

``` yaml
activeFlows:
  - sensor: "timeSensor"
    rules: 
      - comparator: neq
        path: name
        value: test
    tasks:
      - timer: {}
      - serve: {}

  - sensor: "hookSensor"
    rules: []
    tasks:
      - timer:
          - interval: 30
      - serve: {}
      - timer: {}
  
  - sensor: "timeSensor"
    rules: []
    tasks:
      - sleep:
          interval: 1
      - sleep:
          interval: 2
      - sleep:
          interval: 3
```

Над полученным в триггере данными применяются правила и в случае успешного выполнения производится вызов задач. 
Первая задача получает на вход все параметры триггера. Последующие задачи получают на вход результаты работы предыдущей задачи. 

Задачи описаны в пакете `tasks`. Каждая задача имеет следующий интерфейс:

```
type Task interface {
	Run(args Args, exchanger ResultExchanger)
}
```

При завершении выполнения задача должны передать результат выполнения через шину с интерфейсом:

```
type ResultExchanger interface {
	Init() error
	Dispatch(result Result) error
	Reception() (Result, error)
}
```

Результаты выполнения задач сохраняются в хранилище.

Каждая из компонент может быть переопределена при сохранении интерфейсов.

Выбор компоненты осуществляется при запуске утилиты:

```
usage: serve-runner [<flags>]

Flags:
  --help                      Show context-sensitive help (also try --help-long and --help-man).
  --config="config.yml"       Path to config.yml file.
  --config-ctrl="default"     Config controller
  --manage-ctrl="default"     Manage controller
  --rule-ctrl="default"       Rule controller
  --sensor-ex-ctrl="default"  Sensor exchange controller
  --task-ex-ctrl="default"    Task exchange controller
  --version                   Show application version.

```
