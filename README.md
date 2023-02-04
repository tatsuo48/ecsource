# ecsource

ECS Task Definition tool to display resources such as CPU, memory, and memory reservations in an easy-to-read format.

```bash
$ ecsource filename
+----------------------+-----+--------+-------------------+
|         NAME         | CPU | MEMORY | MEMORYRESERVATION |
+----------------------+-----+--------+-------------------+
| Task Setting         | 512 |   1024 |              1024 |
| datadog-agent        |   0 |      0 |               100 |
| nginx                |   0 |      0 |               200 |
| webapp               |   0 |      0 |               424 |
| nr_infra_agent       |   0 |      0 |               100 |
| nr_apm_daemon        |   0 |      0 |               100 |
| log_router           |   0 |      0 |               100 |
| sum of all container |   0 |      0 |              1024 |
+----------------------+-----+--------+-------------------+
|       LEFTOVER       | 512 |  1024  |         0         |
+----------------------+-----+--------+-------------------+
```

## how to install

```bash
go install github.com/tatsuo48/ecsource
```
