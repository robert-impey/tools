cmake_minimum_required(VERSION 3.22)
project(ManagedFolders)

set(CMAKE_CXX_STANDARD 20)

add_executable(ManagedFolders ManagedFolders.cpp)

include(FetchContent)
FetchContent_Declare(
    cli11_proj
    QUIET
    GIT_REPOSITORY https://github.com/CLIUtils/CLI11.git
    GIT_TAG 6c7b07a878ad834957b98d0f9ce1dbe0cb204fc9
)

FetchContent_MakeAvailable(cli11_proj)

target_link_libraries(ManagedFolders PRIVATE CLI11::CLI11)
