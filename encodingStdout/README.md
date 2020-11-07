# Исправляем кодировку в STDOUT
Кракозбры в русском языке можно починить переводом кодировки Windows cp866 в UTF8.
Как правило это требуется только для windows, для Linux вы только все испортите.

1) Что бы Convert() не работал, передаем 0 или 47
2) cp866 кодировка кирилицы, передаем 41
## Использование

```go
import (
	"fmt"
	"github.com/xela07ax/XelaGoDoc/encodingStdout"
)

func main() {
    fmt.Println(string(encodingStdout.Convert(41, []byte("golang �� �����"))))
}
```
## Тестирование
```shell script
go test
```
```shell script
--->cp866:"mvn" �� ����� ����७��� ��� ���譥�
��������, �ᯮ��塞�� �ணࠬ��� ��� �������� 䠩���.

--->newUtf:"mvn" не является внутренней или внешней
командой, исполняемой программой или пакетным файлом.

PASS
ok      _/C_/Users/Tester/Projects/GitHub/XelaGoDoc/encodingStdout      0.017s
```
## Использованные ресурсы
https://studygolang.com/articles/18355

## Поисковые запросы
How to make Unicode charset in cmd.exe by default?
перейти utf-8 в другую кодировку  
golang text encoding 866
golang encoding 866 stdout 
golang 866 to utf-8 
characterset output cmd 
windows cp865 
golang text encoding cyrlic 
text encoding cyrlic 
exec - The Go Programming Language
windows - How to output the output data of `exec.Command` without broken characters
golang exec stdout encoding 
golang stdout encoding 