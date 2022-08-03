#include <iostream>
#include <fstream>

using std::cout; using std::cerr;
using std::endl; using std::string;
using std::ifstream;

bool test(string file) {
    string line;

    ifstream input_file(file);
    if (!input_file.is_open()) {
        cerr << "Could not open the file - '"
             << file << "'" << endl;
        return false;
    }

    if (getline(input_file, line)) {
        input_file.close();

        return line.rfind("#!", 0) == 0;
    }

    return false;
}

int main(int argc, char *argv[]) {
    if (argc == 3) {
        string task = argv[1];

        if (task == "test") {
            string file = argv[2];

            cout << "Testing that " << file << " is a script." << endl;

            if (test(file)) {
                cout << "It is!" << endl;
            } else {
                cout << "It is not!" << endl;
            }

            return 0;
        }
    }

    cerr << "I don't understand!" << endl;

    return -1;
}
