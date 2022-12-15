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
	void list()
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
	void list_pairs()
	{
		auto pairs = find_pairs();

		for (auto& a_pair : pairs)
		{
			cout << a_pair.first << " <-> " << a_pair.second << endl;
		}
	}
private:
	fs::path local_scripts_dir, locations_file_path, folders_file_path;
	vector<string> locations, folders;
	vector<pair<fs::path, fs::path>> find_pairs()
	{
		vector<pair<fs::path, fs::path>> pairs;

		for (auto& folder : folders)
		{
			for (auto& location1 : locations)
			{
				for (auto& location2 : locations)
				{
					if (location1 == location2)
						continue;

					const fs::path location1_path{ location1 };
					const fs::path location2_path{ location2 };

					const fs::path located_folder_path1{ location1_path / folder };
					const fs::path located_folder_path2{ location2_path / folder };

					if (fs::exists(located_folder_path1) && fs::exists(located_folder_path2))
					{
						pair<fs::path, fs::path> a_pair{ located_folder_path1 , located_folder_path2 };
						pair<fs::path, fs::path> reversed_pair { located_folder_path2 , located_folder_path1 };

						if (std::find(pairs.begin(), pairs.end(), reversed_pair) != pairs.end())
						{
							continue;
						}

						pairs.push_back(a_pair);
					}
				}
			}
		}

		return pairs;
	}
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

		FolderManager folder_manager(local_scripts_env);

		if (task == "list")
		{
			folder_manager.list();

			return 0;
		}

		if (task == "list_pairs")
		{
			folder_manager.list_pairs();

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
