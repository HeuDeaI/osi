#include <iostream>
#include <vector>
#include <map>
#include <unistd.h>
#include <sys/wait.h>

using namespace std;

void createProcessTree(vector<int>& parent, map<int, int>& countMap) {
    for (int i = 0; i < parent.size(); i++) {
        if (i < countMap[parent[i]]) {
            pid_t pid = fork();
            if (pid == 0) {
                printf("Process %d with PID %d created by parent %d (PPID %d)\n", i + 1, getpid(), parent[i], getppid());
                if (countMap[parent[i]] == 0) {
                    countMap[parent[i]]--;
                } else {
                    return;
                }
            }
        }
    }
}

int main() {
    vector<int> parent = {1, 1, 2, 3};
    map<int, int> countMap;
    for (int num : parent) {
        countMap[num]++;
    }
    createProcessTree(parent, countMap);
    getchar();
    return 0;
}
