// ManagedFolders.cpp : This file contains the 'main' function. Program execution begins and ends there.
//

#include <iostream>
#include <filesystem>
#include <fstream>
#include <vector>

using namespace std;

namespace fs = std::filesystem;

void list_managed_folders(const fs::path &local_scripts_dir);

int main(int argc, char* argv[])
{
    char* local_scripts_env = getenv("LOCAL_SCRIPTS");

    if (local_scripts_env == NULL)
    {
        cerr << "LOCAL_SCRIPTS env var not set!" << endl;
        return 1;
    }

    fs::path local_scripts_dir = local_scripts_env;
    list_managed_folders(local_scripts_dir);
}

void list_managed_folders(const fs::path &local_scripts_dir)
{
    const auto locations_file_path = local_scripts_dir / "_Common" / "locations.txt";

    ifstream locations_file_stream(locations_file_path);
    if (!locations_file_stream.is_open()) {
        cerr << "Could not open the locations file - '"
            << locations_file_path << "'" << endl;
        return;
    }

    const auto folders_file_path = local_scripts_dir / "_Common" / "folders.txt";

    ifstream folders_file_stream(folders_file_path);
    if (!folders_file_stream.is_open()) {
        cerr << "Could not open the folders file - '"
            << folders_file_path << "'" << endl;
        return;
    }

    vector<string> folders, locations;

    string folders_line, locations_line;
    while (getline(locations_file_stream, locations_line))
    {
        if (locations_line.size())
        {
            locations.push_back(locations_line);
        }
    }
    locations_file_stream.close();

    while (getline(folders_file_stream, folders_line))
    {
        if (folders_line.size())
        {
            folders.push_back(folders_line);
        }
    }
    folders_file_stream.close();
    
    auto first = true;
    for (auto &location : locations)
    {
        if (first)
            first = false;
        else
            cout << endl;

        const fs::path location_path = location;

        for (auto& folder : folders)
        {
            const fs::path located_folder_path = location_path / folder;

            if (fs::exists(located_folder_path))
            {
                cout << located_folder_path << endl;
            }
        }
    }
}