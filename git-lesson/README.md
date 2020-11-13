
# In-memory Git clone, commit и push с помощью GO
(Статья на Medium: https://ishamaraia.medium.com/in-memory-git-clone-commit-and-push-using-go-2e23fd700ef4)  
(Оригинальная стья: https://ish-ar.io.)  
Перевод: Ердяков Алексей  
30 октября 2020 года  
(Переводчики: deepl.com, translate.google.com, translate.yandex.ru)
## Введение в учебное пособие и требования

Сегодняшняя статья представляет собой учебное пособие по настройке и использованию библиотеки go-git для клонирования и обновления репозитория с файловой системой in-memory.

Эта процедура весьма полезна, если вы хотите сделать push или clone хранилища, не касаясь файловой системы ОС и не иметь дела с разрешениями или временными файлами.

Хотя есть документация о git-go, я нахожу её не совсем понятной и иногда вводящей в заблуждение из-за разных версий и названий библиотеки.

По этой причине я решил поделиться этим учебником, и надеюсь, что он кому-нибудь поможет. В этом руководстве мы сделаем все, что угодно внутри файла main.go, однако вам может понадобиться что-нибудь более сложное для вашего случая использования :)

Требования:

- Https Git-репозиторий (Github, Bitbucket не имеет значения, пока он доступен по https).
- Перейти к установке и настройке.
- Базовые знания о Go.

## Настройка файловой системы в памяти

Для настройки файловой системы in-memory мы будем использовать два хранилища два пакета storage и memphis. Хранилище будет содержать объекты, ссылки и другие метаданные (обычно то, что будет делать каталог .git).

Файловая система memfs будет нашей файловой системой для чтения, создания, удаления любого типа файлов в нашем репозитории.

На первом этапе нам нужно создать два объекта (хранилище и файловую систему).
```go
package main

import (
        "fmt"

        billy "github.com/go-git/go-billy/v5"
        memfs "github.com/go-git/go-billy/v5/memfs"
        git "github.com/go-git/go-git/v5"
        http "github.com/go-git/go-git/v5/plumbing/transport/http"
        memory "github.com/go-git/go-git/v5/storage/memory"
)

var storer *memory.Storage
var fs billy.Filesystem

func main() {
        storer = memory.NewStorage()
        fs = memfs.New()
```

## Настройка Git-объектов и клонирование репы

На втором этапе, для того чтобы наше хранилище оказалось в файловой системе in-memory, нам необходимо клонировать его и создать объект рабочего дерева.

Функция Clone() также вернет интерфейс Repository, который мы затем используем для Push() на удаленном компьютере.

Метод Worktree() вернет объект Worktree, который нам понадобится для Add() и Commit() наших изменений.

Наконец, если наш репозиторий приватный, нам нужно будет настроить базовую аутентификацию, чтобы клонировать его

(Нам нужна базовая аутентификация, чтобы надавить в любом случае, так что лучше настроить ее здесь!).

```go
...

        // Authentication
        auth := &http.BasicAuth{
                Username: "your-git-user",
                Password: "your-git-pass",
        }

        repository := "https://github.com/your-org/your-repo"
        r, err := git.Clone(storer, fs, &git.CloneOptions{
                URL:  repository,
                Auth: auth,
        })
        if err != nil {
                fmt.Printf("%v", err)
                return
        }
        fmt.Println("Repository cloned")

        w, err := r.Worktree()
        if err != nil {
                fmt.Printf("%v", err)
                return
        }
```

## Создайте и зафиксируйте ваши файлы

Теперь, будем использовать объект fs для создания реального файла, добавим и зафиксируем его в Worktree().

ПРИМЕЧАНИЕ:

- По умолчанию репозиторий всегда клонируется в &quot;/&quot;.
- По какой-то причине (неизвестной мне), если имя файла начинается с &quot;/&quot;, он создаст папку, а не файл.

Например: /hello/world.txt (world.txt будет каталогом) hello/world.txt (world.txt будет файлом внутри папки hello)
```go
...

        // Create new file
        filePath := "my-new-ififif.txt"
        newFile, err := fs.Create(filePath)
        if err != nil {
                return
        }
        newFile.Write([]byte("My new file"))
        newFile.Close()

        // Run git status before adding the file to the worktree
        fmt.Println(w.Status())

        // git add $filePath
        w.Add(filePath)

        // Run git status after the file has been added adding to the worktree
        fmt.Println(w.Status())

        // git commit -m $message
        w.Commit("Added my new file", &git.CommitOptions{})
```
## Ошибки при нажатии и проверке

Push() будет выполняться методом интерфейса Репозитария, как показано ниже.

```go
...
        //Push the code to the remote
        err = r.Push(&git.PushOptions{
                RemoteName: "origin",
                Auth:       auth,
        })
        if err != nil {
                return
        }
        fmt.Println("Remote updated.", filePath)
        return
}
```

fmt . ![](RackMultipart20201030-4-1o5a50g_html_8da8cb3d6b1dc24d.jpg)обновлено. &quot;. , filePath) return

Наш последний main.go будет выглядеть следующим образом: 
https://github.com/ish-xyz/ish-ar.io-tutorials/blob/master/tutorial-go-git/main.go

## Установка и запуск нашего модуля

Давайте попробуем наш новый модуль Go. Выполните следующие команды:

```sh
go mod init
go build
go run main.go

    Repository cloned
    ?? my-new-file.txt
    <nil>
    A  my-new-file.txt
    <nil>
    Remote updated. my-new-file.txt
```

Я знаю, что это довольно специфическая тема, но я надеюсь, что этот учебник кому-нибудь поможет!
