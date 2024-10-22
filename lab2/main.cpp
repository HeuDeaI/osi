#include <iostream>
#include <vector>
#include <unistd.h>

int main() {
    std::vector<int> tree = {0, 1, 1, 1, 1, 2, 2, 2, 2, 3, 3, 3, 3, 10, 10, 10, 10};
    std::cout << "Root process PID: " << getpid() << std::endl;
    int u = 0;
    
    for (int i = 0; i < tree.size(); i++) {
        if (u == tree[i] && fork() == 0) {
            u = i + 1; 
        }
    }
    
    getchar();
    return 0;
}