#include <filesystem>
#include <fstream>
#include <iostream>
#include <regex>
#include <string>
#include <utility>
#include <vector>
#include <ctime>

using namespace std;

namespace fs = std::filesystem;

vector<string> read_all_non_empty_lines(const fs::path &path);

fs::path find_autogen_path();

string clean_path(const string &);

void generate_folder_synch_script(const string &, const fs::path &, const fs::path &, const fs::path &);

void generate_all_folders_synch_script(const vector<string> &, const fs::path &, const fs::path &, const fs::path &);

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

    void list() const {
        auto first{true};
        for (auto &location: _locations) {
            const fs::path location_path{location};

            if (first)
                first = false;
            else if (exists(location_path))
                cout << endl;

            for (auto &folder: _folders) {
                const fs::path located_folder_path{location_path / folder};

                try {
                    if (exists(located_folder_path)) {
                        cout << located_folder_path.string() << endl;
                    }
                }
                catch (std::filesystem::filesystem_error &e) {
                    std::cerr << e.what() << endl;
                }
            }
        }
    }

    void list_write() {
        auto autogen_path {find_autogen_path()};

        const fs::path managed_folders_path{autogen_path / "managed-folders.txt"};

        ofstream managed_folders_file;
        managed_folders_file.open(managed_folders_path, ios::out | ios::trunc);

        auto now { std::time(0)};
        managed_folders_file << "# AUTOGEN'D at " << std::ctime(&now);
        managed_folders_file << "# DO NOT EDIT!" << endl;

        auto first{true};
        for (auto &location: _locations) {
            const fs::path location_path{location};

            if (first)
                first = false;
            else if (exists(location_path))
                managed_folders_file << endl;

            for (auto &folder: _folders) {
                const fs::path located_folder_path{location_path / folder};

                try {
                    if (exists(located_folder_path)) {
                        managed_folders_file << located_folder_path.string() << endl;
                    }
                }
                catch (std::filesystem::filesystem_error &e) {
                    std::cerr << e.what() << endl;
                }
            }
        }
    }

    void list_pairs() {
        const auto pairs = find_pairs();

        for (const auto &[fst, snd]: pairs) {
            cout << fst.string() << " <-> " << snd.string() << endl;
        }
    }

    void generate_synch_scripts() {
        auto autogen_path { find_autogen_path()};

        auto synch_autogen_path{autogen_path / "synch"};

        if (!exists(synch_autogen_path)) {
            create_directory(synch_autogen_path);
        }

        generate_synch_location_pair_folders(synch_autogen_path);
    }

private:
    vector<string> _locations, _folders;
    vector<fs::path> _location_paths;

    [[nodiscard]] vector<pair<fs::path, fs::path>> find_pairs() const {
        vector<pair<fs::path, fs::path>> pairs;

        for (auto &folder: _folders) {
            for (auto &location1: _locations) {
                for (auto &location2: _locations) {
                    if (location1 == location2)
                        continue;

                    const fs::path location1_path{location1};
                    const fs::path location2_path{location2};

                    const fs::path located_folder_path1{location1_path / folder};
                    const fs::path located_folder_path2{location2_path / folder};

                    try {
                        if (fs::exists(located_folder_path1)
                            && fs::exists(located_folder_path2)) {
                            pair a_pair{located_folder_path1, located_folder_path2};
                            pair reversed_pair{located_folder_path2, located_folder_path1};

                            if (ranges::find(pairs, reversed_pair) != pairs.end()) {
                                continue;
                            }

                            pairs.push_back(a_pair);
                        }
                    }
                    catch (std::filesystem::filesystem_error &e) {
                        std::cerr << e.what() << endl;
                    }
                }
            }
        }

        sort(pairs.begin(), pairs.end());

        return pairs;
    }

    void generate_synch_location_pair_folders(const fs::path &synch_autogen_path) const {
        for (auto &location_path1: _location_paths) {
            const string clean_path1 = clean_path(location_path1.string());

            const fs::path sub_path1{synch_autogen_path / clean_path1};

            if (!exists(sub_path1)) {
                create_directory(sub_path1);
            }

            for (auto &location_path2: _location_paths) {
                if (location_path1 == location_path2) continue;

                const string clean_path2 = clean_path(location_path2.string());

                const fs::path sub_path2{sub_path1 / clean_path2};

                if (!exists(sub_path2)) {
                    create_directory(sub_path2);
                }

                vector<string> common_folders;
                for (auto &folder: _folders) {
                    auto location_folder_path1{location_path1 / folder};
                    auto location_folder_path2{location_path2 / folder};

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
                } else
                    generate_all_folders_synch_script(
                            common_folders,
                            sub_path2,
                            location_path1,
                            location_path2);
            }
        }
    }
};

FolderManager make_folder_manager(const fs::path &);

fs::path get_home_folder() {
#pragma warning( push )
#pragma warning(disable: 4996)
    const auto user_profile{getenv("USERPROFILE")};
#pragma warning( pop )

    if (user_profile != nullptr) {
        const fs::path user_profile_folder{user_profile};

        return user_profile_folder;
    }

#pragma warning( push )
#pragma warning(disable: 4996)
    const auto home{getenv("HOME")};
#pragma warning( pop )

    if (home != nullptr) {
        const fs::path home_folder{home};

        return home_folder;
    }

    throw runtime_error{"Unable to find the home folder!"};
}

fs::path find_local_scripts_path() {
#pragma warning( push )
#pragma warning(disable: 4996)
    const auto local_scripts_env_var{getenv("LOCAL_SCRIPTS")};
#pragma warning( pop )

    if (local_scripts_env_var != nullptr) {
        const fs::path local_scripts_path{local_scripts_env_var};

        return local_scripts_path;
    }

    const auto home_folder_path{get_home_folder()};

    fs::path local_scripts_path{home_folder_path / "local-scripts"};

    return local_scripts_path;
}

int main(const int argc, char *argv[]) {
    const auto local_scripts_path{find_local_scripts_path()};

    if (argc >= 2) {
        const string task{argv[1]};

        auto folder_manager = make_folder_manager(local_scripts_path);

        if (task == "list") {
            folder_manager.list();

            return 0;
        }

        if (task == "list_write") {
            folder_manager.list_write();

            return 0;
        }

        if (task == "list_pairs") {
            folder_manager.list_pairs();

            return 0;
        }

        if (task == "generate_synch_scripts") {
            folder_manager.generate_synch_scripts();

            return 0;
        }
    }

    cerr << "Please tell me what to do!" << endl;
    return -1;
}

vector<string> read_all_non_empty_lines(const fs::path &path) {
    ifstream file_stream{path};
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
    const auto home_folder_path{get_home_folder()};

    fs::path autogen_path{home_folder_path / "autogen"};

    if (!exists(autogen_path)) {
        create_directory(autogen_path);
    }

    return autogen_path;
}

string clean_path(const string &path_str) {
    const std::regex illegals{"[:\\\\/ ]+"};
    const string replacement{"_"};

    const string all_legal{std::regex_replace(path_str, illegals, replacement)};

    const std::regex trailing_underscore{"_$"};

    return std::regex_replace(all_legal, trailing_underscore, "");
}

void add_script_file_params(ofstream &script_file) {
    script_file << "param(" << endl;
    script_file << "    [Parameter (Mandatory = $False)]" << endl;
    script_file << "    [switch]$logged = $False" << endl;
    script_file << ")" << endl << endl;
}

void generate_folder_synch_script(
        const std::string &folder,
        const fs::path &script_folder,
        const fs::path &location_path1,
        const fs::path &location_path2) {
    ostringstream script_name;
    script_name << folder << ".ps1";

    const fs::path script_path{script_folder / script_name.str()};

    ofstream script_file;
    script_file.open(script_path, ios::out | ios::trunc);

    auto now {std::time(0)};
    script_file << "# AUTOGEN'D at " << std::ctime(&now) << endl;
    script_file << "# DO NOT EDIT!" << endl << endl;

    add_script_file_params(script_file);

    script_file << R"(Import-Module "$($env:LOCAL_SCRIPTS)\_Common\synch\Synch.psm1")" << endl << endl;

    script_file << "$folder = \"" << folder << "\"" << endl << endl;
    script_file << "$src = " << location_path1 << endl;
    script_file << "$dst = " << location_path2 << endl << endl;

    script_file << "Synch $folder $src $dst $logged" << endl;

    script_file.close();
}

void generate_all_folders_synch_script(
        const vector<string> &folders,
        const fs::path &script_folder,
        const fs::path &location_path1,
        const fs::path &location_path2) {
    const fs::path script_path{script_folder / "_all.ps1"};

    ofstream script_file;
    script_file.open(script_path, ios::out | ios::trunc);

    auto now {std::time(0)};
    script_file << "# AUTOGEN'D at " << std::ctime(&now) << endl;
    script_file << "# DO NOT EDIT!" << endl << endl;

    add_script_file_params(script_file);

    script_file << R"(Import-Module "$($env:LOCAL_SCRIPTS)\_Common\synch\Synch.psm1")" << endl << endl;

    script_file << "$folders = ";

    auto first{true};
    for (auto &folder: folders) {
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
    auto os_folder { "Windows"};
#elif _WIN64
    auto os_folder { "Windows" };
#elif __linux__
    auto os_folder { "linux"};
#elif __unix
    auto os_folder { "linux"};
#elif __unix__
    auto os_folder { "linux"};
#else
    auto os_folder = "Other";
#endif

    auto locations_file_path{ local_scripts_dir / "_Common" / os_folder / "locations.txt" };

    return locations_file_path;
}

FolderManager make_folder_manager(const fs::path &local_scripts_dir) {
    auto locations_file_path{ find_locations_file_path(local_scripts_dir) };
    auto folders_file_path{local_scripts_dir / "_Common" / "folders.txt"};

    auto locations = read_all_non_empty_lines(locations_file_path);
    auto folders = read_all_non_empty_lines(folders_file_path);

    vector<fs::path> location_paths;

    for (auto &location: locations) {
        const fs::path location_path{location};

        try {
            if (is_directory(location_path)) {
                location_paths.push_back(location_path);
            }
        }
        catch (std::filesystem::filesystem_error &e) {
            std::cerr << e.what() << endl;
        }
    }

    FolderManager folder_manager(
            locations,
            folders,
            location_paths);

    return folder_manager;
}
