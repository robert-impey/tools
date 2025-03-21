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

fs::path get_home_folder();

fs::path find_autogen_path();

fs::path find_tool_autogen_path(const string&);

fs::path find_local_scripts_path();

fs::path find_config_path();

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
		vector<fs::path> location_paths
	) {
		_locations = std::move(locations);
		_folders = std::move(folders);
		_location_paths = std::move(location_paths);
	}

	void list() {
		list_write(cout);
	}

	void list_write(string managed_folders_file) {
		fs::path managed_folders_path;
		if (managed_folders_file.empty()) {
			auto autogen_path{ find_autogen_path() };
			managed_folders_path = autogen_path / "managed-folders.txt";
		}
		else {
			managed_folders_path = managed_folders_file;
		}		

		cout << "Managed folders file: " << managed_folders_path << endl;

		ofstream managed_folders_file_stream;
		managed_folders_file_stream.open(managed_folders_path, ios::out | ios::trunc);

		write_autogen_header(managed_folders_file_stream);

		list_write(managed_folders_file_stream);
	}

	void generate_synch_scripts() {
		auto synch_autogen_path{ find_tool_autogen_path("synch") };

		generate_synch_location_pair_folders(synch_autogen_path);
	}

	static void generate_synch_windows_config_script() {
		auto synch_autogen_path{ find_tool_autogen_path("synch") };

		auto files_file{ find_windows_config_files_file() };

		generate_synch_windows_config_script(synch_autogen_path, files_file);
	}

private:
	vector<string> _locations, _folders;
	vector<fs::path> _location_paths;

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

	void list_write(ostream& out) {
		auto first{ true };
		for (auto& location : _locations) {
			const fs::path location_path{ location };

			if (first)
				first = false;
			else if (exists(location_path))
				out << endl;

			for (auto& folder : _folders) {
				const fs::path located_folder_path{ location_path / folder };

				try {
					if (exists(located_folder_path)) {
						out << located_folder_path.string() << endl;
					}
				}
				catch (std::filesystem::filesystem_error& e) {
					std::cerr << e.what() << endl;
				}
			}
		}
	}

	void list_pairs_write(ostream& out) const {
		for (const auto pairs = find_pairs(); const auto & [fst, snd]: pairs) {
			out << fst.string() << " <-> " << snd.string() << endl;
		}
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

	static void generate_synch_windows_config_script(const fs::path& synch_autogen_path, const fs::path& files_file) {
		cout << "generate_synch_windows_config_script" << endl;

		auto script_name{ "config-Windows.ps1" };

		const fs::path script_path{ synch_autogen_path / script_name };

		ofstream script_file;
		script_file.open(script_path, ios::out | ios::trunc);

		write_autogen_header(script_file);

		cout << "Reading " << files_file << endl;

		auto home_folder{ get_home_folder() };

		auto files = read_all_non_empty_lines(files_file);

		auto config_path = find_config_path();

		auto common_windows_config_path{ config_path / "_Common" / "Windows" };

		for (auto& file : files) {
			ostringstream command1;
			command1 << "ROBOCOPY " << home_folder << " " << common_windows_config_path << " /xo " << file;
			write_powershell_command(script_file, command1.str());

			ostringstream command2;
			command2 << "ROBOCOPY " << common_windows_config_path << " " << home_folder << " /xo " << file;

			write_powershell_command(script_file, command2.str());
		}
	}

	static fs::path find_windows_config_files_file() {
		auto local_scripts_path{ find_local_scripts_path() };

		const auto computer_name{ getenv("ComputerName") };

		if (computer_name != nullptr) {
			const fs::path computer_specific_files_file{
					local_scripts_path / computer_name / "synch" / "config-Windows" / "files.txt" };

			if (fs::exists(computer_specific_files_file))
				return computer_specific_files_file;
		}

		fs::path common_files_file{ local_scripts_path / "_Common" / "synch" / "config-Windows" / "files.txt" };

		return common_files_file;
	}
};

FolderManager make_folder_manager_from_strings(const string&, const string&);
FolderManager make_folder_manager_from_local_scripts_path(const fs::path&);
FolderManager make_folder_manager_from_paths(const fs::path&, const fs::path&);

fs::path get_home_folder() {
#pragma warning( push )
#pragma warning(disable: 4996)
	const auto user_profile{ getenv("USERPROFILE") };
#pragma warning( pop )

	if (user_profile != nullptr) {
		const fs::path user_profile_folder{ user_profile };

		return user_profile_folder;
	}

#pragma warning( push )
#pragma warning(disable: 4996)
	const auto home{ getenv("HOME") };
#pragma warning( pop )

	if (home != nullptr) {
		const fs::path home_folder{ home };

		return home_folder;
	}

	throw runtime_error{ "Unable to find the home folder!" };
}

fs::path find_local_scripts_path() {
#pragma warning( push )
#pragma warning(disable: 4996)
	const auto local_scripts_env_var{ getenv("LOCAL_SCRIPTS") };
#pragma warning( pop )

	if (local_scripts_env_var != nullptr) {
		const fs::path local_scripts_path{ local_scripts_env_var };

		return local_scripts_path;
	}

	const auto home_folder_path{ get_home_folder() };

	fs::path local_scripts_path{ home_folder_path / "local-scripts" };

	return local_scripts_path;
}

fs::path find_config_path() {
#pragma warning( push )
#pragma warning(disable: 4996)
	const auto config_env_var{ getenv("CONFIG") };
#pragma warning( pop )

	if (config_env_var != nullptr) {
		const fs::path config_path{ config_env_var };

		return config_path;
	}

	const auto home_folder_path{ get_home_folder() };

	fs::path config_path{ home_folder_path / "config" };

	return config_path;
}

int main(const int argc, char* argv[]) {
	CLI::App app{ "Managed Folders" };
	const auto u8_argv{ app.ensure_utf8(argv) };

	auto write {false};

	const auto list_sub_command{
		app.add_subcommand("list", "List managed folders") };
	list_sub_command->add_flag("-w,--write", write, "Write list to file");
	string locations_file;
	list_sub_command->add_option("-l,--locations", locations_file, "Locations File");
	string folders_file;
	list_sub_command->add_option("-f,--folders", folders_file, "Folders File");
	string managed_folders_file;
	list_sub_command->add_option("-m,--managed-folders", managed_folders_file, "Managed Folders File");

	CLI::App* list_pairs_sub_command{ app.add_subcommand("list_pairs", "List pairs of managed folders") };
	CLI::App* list_pairs_write_sub_command{ app.add_subcommand("list_pairs_write", "Write list of managed folders to file") };

	const auto generate_synch_scripts_sub_command{ app.add_subcommand("generate_synch_scripts", "Generate Synch Scripts") };
	CLI::App* generate_synch_windows_config_script_sub_command{
	   app.add_subcommand("generate_synch_windows_config_script", "Generate Synch Scripts for Windows Config") };

	app.require_subcommand();

	CLI11_PARSE(app, argc, u8_argv);

	const string task = app.get_subcommands().back()->get_name();

	if (task == "list") {
		auto folder_manager{ make_folder_manager_from_strings(locations_file, folders_file) };

		if (write) {
			folder_manager.list_write(managed_folders_file);
		}
		else {
			folder_manager.list();
		}

		return 0;
	}
	
	const auto local_scripts_path{ find_local_scripts_path() };
	auto folder_manager{ make_folder_manager_from_local_scripts_path(local_scripts_path) };

	if (task == "generate_synch_scripts") {
		folder_manager.generate_synch_scripts();

		return 0;
	}

	if (task == "generate_synch_windows_config_script") {
		FolderManager::generate_synch_windows_config_script();

		return 0;
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

fs::path find_autogen_path() {
	const auto home_folder_path{ get_home_folder() };

	fs::path autogen_path{ home_folder_path / "autogen" };

	if (!exists(autogen_path)) {
		create_directory(autogen_path);
	}

	return autogen_path;
}

fs::path find_tool_autogen_path(const string& tool) {
	const auto autogen_path{ find_autogen_path() };

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

fs::path find_locations_file_path(const fs::path& local_scripts_dir) {
#ifdef _WIN32
	auto os_folder{ "Windows" };
#elif _WIN64
	auto os_folder{ "Windows" };
#elif __linux__
	auto os_folder{ "linux" };
#elif __unix
	auto os_folder{ "linux" };
#elif __unix__
	auto os_folder{ "linux" };
#else
	auto os_folder = "Other";
#endif

	auto locations_file_path{ local_scripts_dir / "_Common" / os_folder / "locations.txt" };

	return locations_file_path;
}

fs::path find_folders_file_path(const fs::path& local_scripts_dir) {
	return local_scripts_dir / "_Common" / "folders.txt";
}

FolderManager make_folder_manager_from_strings(const string& locations_file, const string& folders_file) {
	auto locations_file_path_fs{ fs::path{locations_file} };
	auto folders_file_path_fs{ fs::path{folders_file} };
	return make_folder_manager_from_paths(locations_file_path_fs, folders_file_path_fs);
}

FolderManager make_folder_manager_from_local_scripts_path(const fs::path& local_scripts_dir) {
	auto locations_file_path{ find_locations_file_path(local_scripts_dir) };
	auto folders_file_path{ find_folders_file_path(local_scripts_dir) };

	return make_folder_manager_from_paths(locations_file_path, folders_file_path);
}

FolderManager make_folder_manager_from_paths(const fs::path& locations_file_path, const fs::path& folders_file_path) {
	cout << "Locations file: " << locations_file_path << endl;
	cout << "Folders file: " << folders_file_path << endl;

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
		location_paths);

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