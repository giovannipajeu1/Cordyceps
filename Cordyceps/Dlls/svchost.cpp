#include <Windows.h>
#include <iostream>
#include <shlobj.h>
#include <tchar.h>

SERVICE_STATUS g_ServiceStatus;
SERVICE_STATUS_HANDLE g_ServiceStatusHandle;

void WINAPI ServiceCtrlHandler(DWORD dwControl) {
    switch (dwControl) {
        case SERVICE_CONTROL_STOP:
            g_ServiceStatus.dwCurrentState = SERVICE_STOP_PENDING;
            SetServiceStatus(g_ServiceStatusHandle, &g_ServiceStatus);
            // Realize quaisquer tarefas de limpeza ou liberação de recursos necessárias
            g_ServiceStatus.dwWin32ExitCode = 0;
            g_ServiceStatus.dwCurrentState = SERVICE_STOPPED;
            SetServiceStatus(g_ServiceStatusHandle, &g_ServiceStatus);
            break;
    }
}

void WINAPI ServiceMain(DWORD dwArgc, LPTSTR* lpszArgv) {
    g_ServiceStatus.dwServiceType = SERVICE_WIN32;
    g_ServiceStatus.dwCurrentState = SERVICE_START_PENDING;
    g_ServiceStatus.dwControlsAccepted = SERVICE_ACCEPT_STOP;
    g_ServiceStatus.dwWin32ExitCode = 0;
    g_ServiceStatus.dwServiceSpecificExitCode = 0;
    g_ServiceStatus.dwCheckPoint = 0;
    g_ServiceStatus.dwWaitHint = 0;

    g_ServiceStatusHandle = RegisterServiceCtrlHandler(L"MyService", ServiceCtrlHandler);
    if (g_ServiceStatusHandle == NULL) {
        // Lidar com falha ao registrar o manipulador de controle de serviço
        return;
    }

    g_ServiceStatus.dwCurrentState = SERVICE_RUNNING;
    SetServiceStatus(g_ServiceStatusHandle, &g_ServiceStatus);

    // Obtenha o caminho para a pasta Downloads do usuário atual
    TCHAR downloadsPath[MAX_PATH];
    HRESULT result = SHGetFolderPath(NULL, CSIDL_MYDOCUMENTS, NULL, 0, downloadsPath);
    if (SUCCEEDED(result)) {
        std::wstring agenteExePath = downloadsPath;
        agenteExePath = agenteExePath + L"\\Downloads\\agente.exe";

        // Execute o agente.exe aqui
        STARTUPINFO si = { sizeof(si) };
        PROCESS_INFORMATION pi;
        if (CreateProcess(agenteExePath.c_str(), NULL, NULL, NULL, FALSE, 0, NULL, NULL, &si, &pi)) {
            // Aguardar o término do processo, se necessário
            WaitForSingleObject(pi.hProcess, INFINITE);

            // Fechar os handles do processo
            CloseHandle(pi.hProcess);
            CloseHandle(pi.hThread);
        }
    }

    // O serviço será encerrado aqui, pois o agente.exe foi executado e o serviço não tem outra tarefa contínua

    g_ServiceStatus.dwCurrentState = SERVICE_STOPPED;
    SetServiceStatus(g_ServiceStatusHandle, &g_ServiceStatus);
}

int _tmain(int argc, _TCHAR* argv[]) {
    SERVICE_TABLE_ENTRY ServiceTable[] = {
        { (LPWSTR)L"MyService", (LPSERVICE_MAIN_FUNCTION)ServiceMain },
        { NULL, NULL }
    };

    if (!StartServiceCtrlDispatcher(ServiceTable)) {
        // Lidar com a falha ao iniciar o despacho do serviço
    }

    return 0;
}
