// ManagedFolders.cpp : This file contains the 'main' function. Program execution begins and ends there.
//

#include <iostream>

using namespace std;

void manage_local_scripts(string);

int main(int argc, char* argv[])
{
    if (argc == 2)
    {
        string local_scripts_dir = argv[1];

        manage_local_scripts(local_scripts_dir);

        return 0;
    }

    cerr << "I don't understand!" << endl;
    return 1;
}

void manage_local_scripts(string local_scripts_dir)
{

}