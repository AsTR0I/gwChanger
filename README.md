# gwChanger

## Описание

**gwChanger** – Представляет собой инструмент для управления записями в файле /etc/hosts.

## Функциональность

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
    "sipc_path": ""
}
```

- **hosts**:
    **hosts.hostname** – Доменное имя хоста
    **hosts.ip** – IP-адрес
- **target_hostname** – Строка, содержащая через пробел целевые хосты.
- **sipc_path** – Путь к исполняемому файлу sipc (опционально). Если путь не указан, программа попытается запустить sipc из своей директории или из системного пути.

## Логирование

- Логи записываются в `gw_changer_log_%timestamp%.log`.
- Формат записи: `YYYY-MM-DD HH:MM:SS - Сообщение`.
- Если лог-файл отсутствует, он создаётся автоматически.
- Логи старше 10 дней удаляются автоматически.

## Установка

## Linux

Для установки скрипта на Linux выполните следующую команду:

```bash
cd /usr/local && curl -sSL https://raw.githubusercontent.com/AsTR0I/gwChanger/refs/heads/main/public/gw_changer_install.sh -o gw_changer_install.sh && chmod +x gw_changer_install.sh && sh gw_changer_install.sh
```
### BSD
```bash
cd /usr/local && /usr/bin/fetch https://raw.githubusercontent.com/AsTR0I/gwChanger/refs/heads/main/public/gw_changer_install.sh -o gw_changer_install.sh && chmod +x gw_changer_install.sh && sh gw_changer_install.sh
```