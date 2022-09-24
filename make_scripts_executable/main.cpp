#include <iostream>
#include <fstream>
#include <filesystem>

using std::cerr;
using std::cout;
using std::endl;
using std::ifstream;
using std::string;

namespace fs = std::filesystem;

bool test_file_starts_with_shebang(const string &file);
void update_path(const fs::path &path);
void search(const string &path);

int main(int argc, char *argv[]) {
    if (argc == 3) {
        string task = argv[1];

        if (task == "test") {
            string file = argv[2];

            cout << "Testing if " << file << " starts with a shebang." << endl;

            if (test_file_starts_with_shebang(file)) {
                cout << "It does!" << endl;
            } else {
                cout << "It does not!" << endl;
            }

            return 0;
        }

        if (task == "update") {
            string file = argv[2];

            if (test_file_starts_with_shebang(file)) {
                fs::path path = file;
                update_path(path);
            } else {
                cout << file << " does not start with a shebang! Skipping." << endl;
            }

            return 0;
        }

        if (task == "search") {
            string path = argv[2];

            search(path);

            return 0;
        }
    }

    cerr << "I don't understand!" << endl;

    return -1;
}

bool test_file_starts_with_shebang(const string &file) {
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

void update_path(const fs::path &path) {
    try {
        fs::permissions(path, fs::perms::owner_all);
    }
    catch (std::exception& e) {
        cerr << e.what() << endl;
    }
}

void search(const string &path) {
    for (auto const& dir_entry : fs::recursive_directory_iterator(path)) {
        if (dir_entry.is_regular_file() && test_file_starts_with_shebang(dir_entry.path().string())) {
            std::cout << "Updating the perms of " << dir_entry << '\n';
            update_path(dir_entry.path());
        }
    }
}