function setDataBase() {
    whiptail --title "Pager" --msgbox "Перед началом установки убедитесь что у вас установлена База Данных MySQL и у вас есть данные для подключения к ней. Для продолжения установки Pager нажмите - ОК" 10 60
      # License and server name
      DB_HOST=$(whiptail --title "DataBase IP" --inputbox "Введите IP сервера DataBase. Пример: 127.0.0.1" 10 60 3>&1 1>&2 2>&3)
      exitstatus=$?
      if [ $exitstatus = 0 ]; then
        echo "DataBase IP" "$DB_HOST"
      else
        echo "Exit"
        exit 1
      fi
      DB_PORT=$(whiptail --title "DataBase Port" --inputbox "Введите PORT сервера DataBase. Пример: 5432" 10 60 3>&1 1>&2 2>&3)
      exitstatus=$?
      if [ $exitstatus = 0 ]; then
        echo "DataBase Port" "$DB_PORT"
      else
        echo "Exit"
        exit 1
      fi
      DB_USER=$(whiptail --title "DataBase User" --inputbox "Введите username пользователя DataBase. Пример: username" 10 60 3>&1 1>&2 2>&3)
      exitstatus=$?
      if [ $exitstatus = 0 ]; then
        echo "DataBase User" "$DB_USER"
      else
        echo "Exit"
        exit 1
      fi
      DB_PASSWORD=$(whiptail --title "DataBase User Password" --inputbox "Введите пароль пользователя DataBase. Пример: qwerty123" 10 60 3>&1 1>&2 2>&3)
      exitstatus=$?
      if [ $exitstatus = 0 ]; then
        echo "DataBase User Password" "$DB_PASSWORD"
      else
        echo "Exit"
        exit 1
      fi
      DB_NAME=$(whiptail --title "DataBase Name" --inputbox "Введите имя Базы данных. Пример: postgres" 10 60 3>&1 1>&2 2>&3)
      exitstatus=$?
      if [ $exitstatus = 0 ]; then
        echo "DataBase Name" "$DB_NAME"
      else
        echo "Exit"
        exit 1
      fi
      if (whiptail --title "DataBase Config" --yes-button "Да" --no-button "Нет" --yesno "Создать таблицу mflash в Базе Данных? Если вы устанавливаете кластер, то данную операцию нужно выполнять только на одном сервере WEB." 10 60) then
                echo "Create mflash in DataBase"
                DB_SSLMODE="enable"
                echo $DB_SSLMODE
              else
                echo "Create mflash in Database: Skip."
                DB_SSLMODE="disable"
                echo $DB_SSLMODE
              fi
      APP_PORT=$(whiptail --title "Application Port" --inputbox "Введите порт, на котором требуется запустить API сервис. Пример: 433" 10 60 3>&1 1>&2 2>&3)
      exitstatus=$?
      if [ $exitstatus = 0 ]; then
        echo "Application Port" "$APP_PORT"
      else
        echo "Exit"
        exit 1
      fi
      LANGUAGE=$(whiptail --title "API Language" --inputbox "Введите язык API. Пример: en" 10 60 3>&1 1>&2 2>&3)
      exitstatus=$?
      if [ $exitstatus = 0 ]; then
        echo "API Language" "$LANGUAGE"
      else
        echo "Exit"
        exit 1
      fi
}

function GoInstall() {
    mkdir /opt/pager/
    mkdir /var/log/pager
    sudo wget https://go.dev/dl/go1.21.4.linux-amd64.tar.gz
    sudo tar -C /usr/local -xzf go1.21.4.linux-amd64.tar.gz
    export PATH=$PATH:/usr/local/go/bin
    rm -f go1.21.4.linux-amd64.tar.gz
}


#function ApiInstall() {
#    cd /opt/pager
#    git init
#    git remote add origin git@github.com:masya2066/ApiMessanger.git
#    git fetch
#    git pull origin main
#
#    rm -f .env
#    file_name=".env"
#    text_content=$(cat <<EOL
#    DB_HOST: $1
#    DB_PORT: $2
#    DB_USER: $3
#    DB_PASSWORD: $4
#    DB_NAME: $5
#    DB_SSLMODE: $6
#    APP_PORT: $7
#    LANGUAGE: $8
#    TOKEN_LIFE_TIME: $9
#    DATE_FORMAT: $10
#    EOL
#    )
#
#    echo "$text_content" > "$file_name"
#    chown -R connector:connector /opt/pager/*
#    chown -R connector:connector /var/log/pager
#    go build
#}

#function MakeService() {
#    cd /etc/systemd/system/
#
#    systemctl stop pager-api.service
#    systemctl disable pager-api.service
#    rm pager-api.service
#
#    service_file="pager-api.service"
#    service_content=$(cat <<EOL
#    [Unit]
#    Description=Api service Connector
#    After=network.target
#
#    [Service]
#    ExecStart=/opt/pager/ApiMessenger
#    WorkingDirectory=/opt/pager/
#    Restart=always
#    User=root
#    Group=root
#    StandardOutput=append:/var/log/pager/api.log
#    StandardError=append:/var/log/pager/api-error.log
#
#    [Install]
#    WantedBy=multi-user.target
#    EOL
#    )
#
#    echo "$service_content" > "$service_file"
#
#    chmod +x pager-api.service
#
#    systemctl enable pager-api.service
#    systemctl start pager-api.service
#    systemctl status pager-api.service
#    chown -R connector:connector /var/log/pager/*
#}

function test() {
    echo $DB_HOST
    echo $DB_PORT
    echo $DB_USER
    echo $DB_PASSWORD
    echo $DB_NAME
    echo $DB_SSLMODE
    echo $APP_PORT
    echo $LANGUAGE
    echo '86400s'
    echo '2006-01-02 15:04:05'
}

setDataBase
test