# 1. Crear admin inicial
cd localserver
make seed DB_PATH=./barbershop.db
# (ingresa usuario y password)

# 2. Levantar server
make run DB_PATH=./barbershop.db PORT=8080

# 3. En otra terminal, levantar desktopapp
cd ../desktopapp
make run