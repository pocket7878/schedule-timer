# schedule-timer

[![asciicast](https://asciinema.org/a/2ZaIUz1V2s0ebtwUUyZ00mE6M.svg)](https://asciinema.org/a/2ZaIUz1V2s0ebtwUUyZ00mE6M)

## Usage

```shell
schedule-timer example/repeat.yaml
```

## YAML Format

```yaml
name: Pomodoro
repeat: true
tasks:
  - name: Work 1
    duration: 1500
  - name: Short Break
    duration: 300
  - name: Work 2
    duration: 1500
  - name: Short Break
    duration: 300
  - name: Work 3
    duration: 1500
  - name: Short Break
    duration: 300
  - name: Work 4
    duration: 1500
  - name: Long Break
    duration: 900
```
