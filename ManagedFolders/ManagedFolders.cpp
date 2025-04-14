#include <CLI/CLI.hpp>
#include <ctime>
#include <filesystem>
#include <fstream>
#include <iostream>
#include <regex>
#include <string>
#include <utility>
#include <vector>

using namespace std;

namespace fs = std::filesystem;

vector<string> read_all_non_empty_lines(const fs::path&);

fs::path find_tool_autogen_path(const fs::path&, const string&);

string clean_path(const string&);

void generate_folder_synch_script(const string&, const fs::path&, const fs::path&, const fs::path&);

void generate_all_folders_synch_script(const vector<string>&, const fs::path&, const fs::path&, const fs::path&);

void write_autogen_header(ostream&);

void write_powershell_command(ostream&, string);

class FolderManager {
public:
	explicit FolderManager(
		vector<string> locations,
		vector<string> folders,
		vector<fs::path> location_paths,
		const fs::path& autogen_dir_path
	) {
		_locations = std::move(locations);
		_folders = std::move(folders);
		_location_paths = std::move(location_paths);
		_autogen_dir_path = autogen_dir_path;
	}

	void generate_synch_scripts() {
		auto synch_autogen_path{ find_tool_autogen_path(_autogen_dir_path, "synch")};

		generate_synch_location_pair_folders(synch_autogen_path);
	}

private:
	vector<string> _locations, _folders;
	vector<fs::path> _location_paths;
	fs::path _autogen_dir_path;

	[[nodiscard]] vector<pair<fs::path, fs::path>> find_pairs() const {
		vector<pair<fs::path, fs::path>> pairs;

		for (auto& folder : _folders) {
			for (auto& location1 : _locations) {
				for (auto& location2 : _locations) {
					if (location1 == location2)
						continue;

					const fs::path location1_path{ location1 };
					const fs::path location2_path{ location2 };

					const fs::path located_folder_path1{ location1_path / folder };
					const fs::path located_folder_path2{ location2_path / folder };

					try {
						if (exists(located_folder_path1)
							&& exists(located_folder_path2)) {
							pair a_pair{ located_folder_path1, located_folder_path2 };

							if (pair reversed_pair{ located_folder_path2, located_folder_path1 };
								ranges::find(pairs, reversed_pair) != pairs.end()) {
								continue;
							}

							pairs.push_back(a_pair);
						}
					}
					catch (std::filesystem::filesystem_error& e) {
						std::cerr << e.what() << endl;
					}
				}
			}
		}

		sort(pairs.begin(), pairs.end());

		return pairs;
	}

	void generate_synch_location_pair_folders(const fs::path& synch_autogen_path) const {
		for (auto& location_path1 : _location_paths) {
			const string clean_path1 = clean_path(location_path1.string());

			const fs::path sub_path1{ synch_autogen_path / clean_path1 };

			if (!exists(sub_path1)) {
				create_directory(sub_path1);
			}

			for (auto& location_path2 : _location_paths) {
				if (location_path1 == location_path2) continue;

				const string clean_path2 = clean_path(location_path2.string());

				const fs::path sub_path2{ sub_path1 / clean_path2 };

				if (!exists(sub_path2)) {
					create_directory(sub_path2);
				}

				vector<string> common_folders;
				for (auto& folder : _folders) {
					auto location_folder_path1{ location_path1 / folder };
					auto location_folder_path2{ location_path2 / folder };

					if (exists(location_folder_path1) && exists(location_folder_path2)) {
						common_folders.push_back(folder);
						generate_folder_synch_script(
							folder,
							sub_path2,
							location_path1,
							location_path2);
					}
				}

				if (common_folders.empty()) {
					if (fs::is_empty(sub_path2))
						fs::remove(sub_path2);
				}
				else
					generate_all_folders_synch_script(
						common_folders,
						sub_path2,
						location_path1,
						location_path2);
			}
		}
	}
};

FolderManager make_folder_manager_from_strings(const string&, const string&, const string&);
FolderManager make_folder_manager_from_paths(const fs::path&, const fs::path&, const fs::path&);

int main(const int argc, char* argv[]) {
	CLI::App app{ "Managed Folders" };
	const auto u8_argv{ app.ensure_utf8(argv) };

	string locations_file;
	string folders_file;
	string autogen_dir;
	
	const auto generate_synch_scripts_sub_command{ app.add_subcommand("generate_synch_scripts", "Generate Synch Scripts") };

	generate_synch_scripts_sub_command->add_option("-l,--locations", locations_file, "Locations File");
	generate_synch_scripts_sub_command->add_option("-f,--folders", folders_file, "Folders File");
	generate_synch_scripts_sub_command->add_option("-a,--autogen", autogen_dir, "Autogen Directory");

	app.require_subcommand();

	CLI11_PARSE(app, argc, u8_argv);

	const string task = app.get_subcommands().back()->get_name();

	try {
		if (task == "generate_synch_scripts") {
			auto folder_manager{ make_folder_manager_from_strings(locations_file, folders_file, autogen_dir) };
			folder_manager.generate_synch_scripts();

			return 0;
		}
	}
	catch (const CLI::Error& e) {
		cerr << "Error: " << e.what() << endl;
		return 1;
	}
	catch (const std::exception& e) {
		cerr << "Error: " << e.what() << endl;
		return 1;
	}
	catch (...) {
		cerr << "Unknown error occurred!" << endl;
		return 1;
	}
}

vector<string> read_all_non_empty_lines(const fs::path& path) {
	ifstream file_stream{ path };
	if (!file_stream.is_open()) {
		throw std::runtime_error("Failed to open " + path.string() + "'");
	}

	vector<string> lines;
	string line;
	while (getline(file_stream, line)) {
		if (!line.empty()) {
			lines.push_back(line);
		}
	}
	file_stream.close();

	return lines;
}

fs::path find_tool_autogen_path(const fs::path& autogen_path, const string& tool) {
	fs::path tool_autogen_path{ autogen_path / tool };

	if (!exists(tool_autogen_path)) {
		create_directory(tool_autogen_path);
	}

	return tool_autogen_path;
}

string clean_path(const string& path_str) {
	const std::regex illegals{ "[:\\\\/ ]+" };
	const string replacement{ "_" };

	const string all_legal{ std::regex_replace(path_str, illegals, replacement) };

	const std::regex trailing_underscore{ "_$" };

	return std::regex_replace(all_legal, trailing_underscore, "");
}

void add_script_file_params(ofstream& script_file) {
	script_file << "param(" << endl;
	script_file << "    [Parameter (Mandatory = $False)]" << endl;
	script_file << "    [switch]$logged = $False" << endl;
	script_file << ")" << endl << endl;
}

void generate_folder_synch_script(
	const std::string& folder,
	const fs::path& script_folder,
	const fs::path& location_path1,
	const fs::path& location_path2) {
	ostringstream script_name;
	script_name << folder << ".ps1";

	const fs::path script_path{ script_folder / script_name.str() };

	cout << "Generating " << script_path << endl;

	ofstream script_file;
	script_file.open(script_path, ios::out | ios::trunc);

	write_autogen_header(script_file);

	add_script_file_params(script_file);

	script_file << R"(Import-Module "$($env:LOCAL_SCRIPTS)\_Common\synch\Synch.psm1")" << endl << endl;

	script_file << "$folder = \"" << folder << "\"" << endl << endl;
	script_file << "$src = " << location_path1 << endl;
	script_file << "$dst = " << location_path2 << endl << endl;

	script_file << "Synch $folder $src $dst $logged" << endl;

	script_file.close();
}

void generate_all_folders_synch_script(
	const vector<string>& folders,
	const fs::path& script_folder,
	const fs::path& location_path1,
	const fs::path& location_path2) {
	const fs::path script_path{ script_folder / "_all.ps1" };

	ofstream script_file;
	script_file.open(script_path, ios::out | ios::trunc);

	cout << "Generating " << script_path << endl;

	write_autogen_header(script_file);

	add_script_file_params(script_file);

	script_file << R"(Import-Module "$($env:LOCAL_SCRIPTS)\_Common\synch\Synch.psm1")" << endl << endl;

	script_file << "$folders = ";

	auto first{ true };
	for (auto& folder : folders) {
		if (first)
			first = false;
		else
			script_file << ", ";

		script_file << "\"" << folder << "\"";
	}

	script_file << endl << endl;

	script_file << "$src = " << location_path1 << endl;
	script_file << "$dst = " << location_path2 << endl << endl;

	script_file << "foreach($folder in $folders)" << endl;
	script_file << "{" << endl;
	script_file << "    Synch $folder $src $dst $logged" << endl;
	script_file << "}" << endl;

	script_file.close();
}

FolderManager make_folder_manager_from_strings(const string& locations_file, const string& folders_file, const string& autogen_dir) {
	if (locations_file.empty() || folders_file.empty()) {
		throw runtime_error{ "Locations folders files must be specified!" };
	}

	if (autogen_dir.empty()) {
		throw runtime_error{ "The autogen directory must be specified!" };
	}

	auto locations_file_path_fs{ fs::path{locations_file} };
	auto folders_file_path_fs{ fs::path{folders_file} };
	auto autogen_dir_path_fs{ fs::path{autogen_dir} };
	return make_folder_manager_from_paths(locations_file_path_fs, folders_file_path_fs, autogen_dir_path_fs);
}

FolderManager make_folder_manager_from_paths(const fs::path& locations_file_path, const fs::path& folders_file_path, const fs::path& autogen_dir_path) {
	cout << "Locations file: " << locations_file_path << endl;
	cout << "Folders file: " << folders_file_path << endl;
	cout << "Autogen Directory: " << autogen_dir_path << endl;

	if (!is_directory(autogen_dir_path)) {
		throw runtime_error{ "The autogen directory must be a directory!" };
	}

	auto locations = read_all_non_empty_lines(locations_file_path);
	auto folders = read_all_non_empty_lines(folders_file_path);

	vector<fs::path> location_paths;

	for (auto& location : locations) {
		const fs::path location_path{ location };

		try {
			if (is_directory(location_path)) {
				location_paths.push_back(location_path);
			}
		}
		catch (std::filesystem::filesystem_error& e) {
			std::cerr << e.what() << endl;
		}
	}

	FolderManager folder_manager(
		locations,
		folders,
		location_paths,
		autogen_dir_path);

	return folder_manager;
}

void write_autogen_header(ostream& out) {
	auto now{ std::time(nullptr) };
	out << "# AUTOGEN'D at " << std::ctime(&now);
	out << "# DO NOT EDIT!" << endl << endl;
}

void write_powershell_command(ostream& out, string command) {
	out << command << endl;
}