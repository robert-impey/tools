// ManagedFolders.cpp : This file contains the 'main' function. Program execution begins and ends there.
//

#include <iostream>
#include <filesystem>
#include <fstream>
#include <vector>

using namespace std;

namespace fs = std::filesystem;

vector<string> read_all_non_empty_lines(const fs::path& path);

void list_managed_folders(const fs::path& local_scripts_dir);

class FolderManager
{
public:
    FolderManager(string _local_scripts_env)
    {
        local_scripts_env = _local_scripts_env;
        local_scripts_dir = local_scripts_env;
    }
    void list_all_folders()
    {
        list_managed_folders(local_scripts_dir);
    }
private:
    string local_scripts_env;
    fs::path local_scripts_dir;
};


int main(int argc, char* argv[])
{
    auto local_scripts_env{ getenv("LOCAL_SCRIPTS") };

    if (local_scripts_env == NULL)
    {
        cerr << "LOCAL_SCRIPTS env var not set!" << endl;
        return 1;
    }

    if (argc >= 2)
    {
        string task{ argv[1] };

        if (task == "list")
        {
            fs::path local_scripts_dir{ local_scripts_env };
            list_managed_folders(local_scripts_dir);

            return 0;
        }
    }

    cerr << "Please tell me what to do!" << endl;
    return -1;
}

void list_managed_folders(const fs::path &local_scripts_dir)
{
    const auto locations_file_path{ local_scripts_dir / "_Common" / "locations.txt" };
    const auto folders_file_path{ local_scripts_dir / "_Common" / "folders.txt" };

    const auto locations{ read_all_non_empty_lines(locations_file_path) };
    const auto folders{ read_all_non_empty_lines(folders_file_path) };
    
    auto first{ true };
    for (auto &location : locations)
    {
        if (first)
            first = false;
        else
            cout << endl;

        const fs::path location_path{ location };

        for (auto& folder : folders)
        {
            const fs::path located_folder_path{ location_path / folder };

            if (fs::exists(located_folder_path))
            {
                cout << located_folder_path << endl;
            }
        }
    }
}

vector<string> read_all_non_empty_lines(const fs::path& path)
{
    ifstream file_stream{ path };
    if (!file_stream.is_open()) {
        throw std::runtime_error("Failed to open " + path.string() + "'");
    }

    vector<string> lines;
    string line;
    while (getline(file_stream, line))
    {
        if (line.size())
        {
            lines.push_back(line);
        }
    }
    file_stream.close();

    return lines;
}
