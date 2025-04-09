# gwChanger

## Описание

**gwChanger** – Представляет собой инструмент ...

## Функциональность

### 
- При изменении в /etc/hosts происходит отправка сообщения на почту из конфига
### Командная строка:
-   **-h, --help** - Показать справку.
-   **-v, --version** - Показать версию программы.
-   **-lsd N** - Установить количество дней для сохранения логов (по умолчанию 10).

## Формат конфигурации (JSON)

```json
{
    "hosts": [
        { "hostname": "yandex.ru", "ip": "77.88.55.88" },
        { "hostname": "yandex.ru", "ip": "77.88.55.88" }
    ],
    "target_hostname": "voip.test voip.test2",
    "sipc_path": "",
    "mail": {
        "from": "",
        "to": "",
        "smtp_server": "",
        "smtp_server_port": ""
    },
    "hostname_machine": ""
}
```

- **hosts**:
    **hosts.hostname** – Доменное имя хоста
    **hosts.ip** – IP-адрес
- **target_hostname** – Строка, содержащая через пробел целевые хосты.
- **sipc_path** – Путь к исполняемому файлу sipc (опционально). Если путь не указан, программа попытается запустить sipc из своей директории или из системного пути.
- **mail**:
    **from** – from
    **to** – to
    **smtp_server** – smtp_server
    **smtp_server_port** – smtp_server_port
- **hostname_machine** – hostname_machine

## Логирование

- Логи записываются в `gw_changer_log_%timestamp%.log`.
- Формат записи: `YYYY-MM-DD HH:MM:SS - Сообщение`.
- Если лог-файл отсутствует, он создаётся автоматически.
- Логи старше 10 дней удаляются автоматически.

## Установка

## Linux

Для установки скрипта на Linux выполните следующую команду:

```bash
cd /usr/local && curl https://raw.githubusercontent.com/AsTR0I/gwChanger/refs/heads/main/public/gw_changer_install.sh -o gw_changer_install.sh && chmod +x gw_changer_install.sh && sh gw_changer_install.sh
```
### BSD

Для установки скрипта на BSD выполните следующую команду:

```bash
cd /usr/local && /usr/bin/fetch https://raw.githubusercontent.com/AsTR0I/gwChanger/refs/heads/main/public/gw_changer_install.sh -o gw_changer_install.sh && chmod +x gw_changer_install.sh && sh gw_changer_install.sh
```