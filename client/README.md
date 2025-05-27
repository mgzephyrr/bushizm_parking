## Установка зависимостей

### Установка Docker и Make

#### macOS

1. **Установка Docker Desktop**:
   - Перейдите на [официальный сайт Docker](https://www.docker.com/products/docker-desktop) и скачайте установщик Docker Desktop для macOS.
   - Откройте скачанный файл `.dmg` и перетащите иконку Docker в папку "Программы".
   - Откройте Docker из папки "Программы" и дождитесь, пока он полностью загрузится.

2. **Установка Make**:
   - Откройте терминал и выполните следующую команду для установки Homebrew (если еще не установлен):

   ```bash
   /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
   ```

   - Установите Make:

   ```bash
   brew install make
   ```

#### Linux

1. **Установка Docker**:
   - Откройте терминал и выполните следующие команды:

   ```bash
   sudo apt-get update
   sudo apt-get install -y apt-transport-https ca-certificates curl software-properties-common
   curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
   sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"
   sudo apt-get update
   sudo apt-get install -y docker-ce
   ```

2. **Проверьте установку**:
   - Убедитесь, что Docker установлен, выполнив команду:

   ```bash
   docker --version
   ```

3. **Установка Make**:
   - Выполните следующую команду:

   ```bash
   sudo apt-get install -y make
   ```

#### Windows

Рекомендуется делать действия для Linux через [WSL](https://learn.microsoft.com/ru-ru/windows/wsl/install). Альтернатива представлена ниже

1. **Установка Docker Desktop**:
   - Перейдите на [официальный сайт Docker](https://www.docker.com/products/docker-desktop) и скачайте установщик Docker Desktop для Windows.
   - Запустите скачанный установщик и следуйте инструкциям на экране.
   - После установки запустите Docker Desktop и дождитесь, пока он полностью загрузится.

2. **Установка Make**:
   - Откройте PowerShell от имени администратора и выполните следующую команду для установки Chocolatey (если еще не установлен):

   ```powershell
   Set-ExecutionPolicy Bypass -Scope Process -Force; [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.SecurityProtocolType]::Tls12; iex ((New-Object System.Net.WebClient).DownloadString('https://chocolatey.org/install.ps1'))
   ```

   или https://chocolatey.org/install

   - Установите Make:

   ```powershell
   choco install make
   ```

### Проверка установки

После установки Docker и Make, вы можете проверить их работоспособность, выполнив следующие команды в терминале или командной строке:

```bash
docker --version
make --version
```

Если обе команды возвращают версии, значит, установка прошла успешно!

## Управление Docker-контейнером

Для управления Docker-контейнером приложения доступны следующие команды:

### Основные команды

- `make initial-start` - Полная перезапуск приложения (очистка + сборка + запуск)
- `make build` - Собрать Docker-образ приложения
- `make run` - Запустить контейнер с приложением
- `make stop` - Остановить контейнер
- `make clean` - Остановить и удалить контейнер
- `make clean-image` - Удалить Docker-образ
- `make clean-all` - Полная очистка (удаление контейнера и образа)

### Пример использования

1. Первый запуск приложения:
   ```bash
   make initial-start
   ```

2. Остановка приложения:
   ```bash
   make stop
   ```

3. Удаление контейнера:
   ```bash
   make clean
   ```

4. Полная очистка:
   ```bash
   make clean-all
   ```


## Запуск проекта

### 1. Запуск через Docker

Для запуска приложения в контейнере потребуется установленный Docker и предоставленный Makefile.

Убедитесь, что Docker установлен на вашем компьютере.
Выполните команду:
```bash
    make initial-start
   ```

Откройте приложение в браузере по адресу: http://localhost:5173.

### 2. Запуск в режиме разработчика
   Для запуска в режиме разработчика потребуется установленный Node.js (версии 18 или выше) и пакетный менеджер npm.

Установите зависимости:

```bash
    npm install
```

Запустите сервер разработки:
```bash
    npm run dev
```
