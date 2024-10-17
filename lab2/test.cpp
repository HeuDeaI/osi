#include <iostream>
#include <sys/wait.h>
#include <unistd.h>
#include <vector>

void fork_processes(const std::vector<int>& relations, int parent, int execIndex) {
    for (int i = 0; i < relations.size(); ++i) {
        if (relations[i] == parent) {
            pid_t pid = fork();

            if (pid < 0) {
                std::cerr << "Fork error!" << std::endl;
                exit(1);
            } else if (pid == 0) {
                std::cout << "Child PID: " << getpid() << ", Parent PID: " << getppid() << std::endl;

                if (i == execIndex) {
                    const char* cmd = "/bin/ps";
                    execlp(cmd, cmd, nullptr);
                } else {
                    fork_processes(relations, i, execIndex);
                }
                exit(0);
            } else {
                waitpid(pid, nullptr, 0);
            }
        }
    }
}

int main() {
    std::vector<int> relations = { -1, 0, 1, 1, 1, 3, 3, 2 };
    int execIndex = 7;
    fork_processes(relations, 0, execIndex);
    return 0;
}