// ManagedFolders.cpp : This file contains the 'main' function. Program execution begins and ends there.
//

#include <iostream>
#include <filesystem>
#include <fstream>

using namespace std;

namespace fs = std::filesystem;

void manage_local_scripts(const fs::path &local_scripts_dir);

int main(int argc, char* argv[])
{
    if (argc == 2)
    {
        fs::path local_scripts_dir = argv[1];
        manage_local_scripts(local_scripts_dir);

        return 0;
    }

    cerr << "I don't understand!" << endl;
    return 1;
}

void manage_local_scripts(const fs::path &local_scripts_dir)
{
    const auto folders_file_path = local_scripts_dir / "_Common" / "folders.txt";

    ifstream folders_file_stream(folders_file_path);

    if (!folders_file_stream.is_open()) {
        cerr << "Could not open the folders file - '"
            << folders_file_path << "'" << endl;
        return;
    }

    const auto locations_file_path = local_scripts_dir / "_Common" / "locations.txt";

    ifstream locations_file_stream(locations_file_path);

    if (!locations_file_stream.is_open()) {
        cerr << "Could not open the locations file - '"
            << locations_file_path << "'" << endl;
        return;
    }


}