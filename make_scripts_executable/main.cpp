#include <iostream>

using namespace std;

int main(int argc, char *argv[]) {
    if (argc == 3) {
        string task = argv[1];

        if (task == "test") {
            string file = argv[2];

            cout << "Testing that " << file << " is a script." << endl;

            return 0;
        }
    }

    cerr << "I don't understand!" << endl;

    return -1;
}
