// ManagedFolders.cpp : This file contains the 'main' function. Program execution begins and ends there.
//

#include <iostream>
#include <filesystem>
#include <fstream>
#include <vector>

using namespace std;

namespace fs = std::filesystem;

vector<string> read_all_non_empty_lines(const fs::path& path);

class FolderManager
{
public:
	FolderManager(string _local_scripts_env)
	{
		local_scripts_dir = _local_scripts_env;

		locations_file_path = local_scripts_dir / "_Common" / "locations.txt";
		folders_file_path = local_scripts_dir / "_Common" / "folders.txt";

		locations = read_all_non_empty_lines(locations_file_path);
		folders = read_all_non_empty_lines(folders_file_path);
	}
	void list_all_folders()
	{
		auto first{ true };
		for (auto& location : locations)
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
private:
	fs::path local_scripts_dir, locations_file_path, folders_file_path;
	vector<string> locations, folders;
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
			FolderManager folder_manager(local_scripts_env);
			folder_manager.list_all_folders();

			return 0;
		}
	}

	cerr << "Please tell me what to do!" << endl;
	return -1;
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
