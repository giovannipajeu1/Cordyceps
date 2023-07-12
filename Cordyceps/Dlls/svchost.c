#include <stdio.h>
#include <windows.h>
#include <knownfolders.h>
#include <shlobj.h>

int main() {
    // Obter o caminho completo do diretório "Downloads" do usuário atual
    PWSTR downloadsPath;
    HRESULT result = SHGetKnownFolderPath(&FOLDERID_Downloads, 0, NULL, &downloadsPath);
    if (SUCCEEDED(result)) {
        // Concatenar o caminho completo do arquivo .exe
        wchar_t exePath[MAX_PATH];
        swprintf(exePath, MAX_PATH, L"%s\\agente.exe", downloadsPath);

        // Injetar o comando na memória RAM
        wchar_t cmd[MAX_PATH * 2];
        swprintf(cmd, MAX_PATH * 2, L"cmd.exe /C \"%s\"", exePath);

        // Resto do código de injeção de memória aqui...

        // Liberar a memória alocada
        CoTaskMemFree(downloadsPath);
    }

    return 0;
}
