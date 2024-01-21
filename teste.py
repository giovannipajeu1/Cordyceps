def login(username, password):
    if username == "admin" and password == "admin":
        return True
    else:
        return False

def main():
    print("Hello, World!")

    # Exemplo de uso e senha
    input_username = input("Digite o nome de usuário: ")
    input_password = input("Digite a senha: ")

    if login(input_username, input_password):
        print("Login bem-sucedido!")
    else:
        print("Login falhou. Usuário ou senha incorretos.")

if __name__ == "__main__":
    main()
