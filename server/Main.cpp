#include <Composition/ServerApplication.hpp>
#include <Runtime/ServerConfiguration.hpp>

#include <iostream>
#include <filesystem>
#include <string_view>
#include <utility>

namespace {

void ReportFailure(const Menu::Foundation::Error& ErrorValue) {
    std::cerr << "MenuServer startup failed: "
              << ErrorValue.MessageValue() << '\n';
}

}  // namespace

int main(int ArgumentCount, char** Arguments) {
    std::filesystem::path ConfigPath = "config/Menu.example.json";
    for (int Index = 1; Index < ArgumentCount; ++Index) {
        if (std::string_view(Arguments[Index]) == "--config" && Index + 1 < ArgumentCount) {
            ConfigPath = Arguments[++Index];
        } else {
            std::cerr << "Usage: MenuServer [--config path]\n";
            return 2;
        }
    }

    const auto ConfigurationResult = Menu::Runtime::LoadServerConfiguration(ConfigPath);
    if (!ConfigurationResult.HasValue()) {
        ReportFailure(ConfigurationResult.ErrorValue());
        return 2;
    }
    auto ApplicationResult = Menu::Server::ServerApplication::Create(
        std::move(ConfigurationResult).Value());
    if (!ApplicationResult.HasValue()) {
        ReportFailure(ApplicationResult.ErrorValue());
        return 2;
    }
    auto Application = std::move(ApplicationResult).Value();
    const auto StartResult = Application->Start();
    if (!StartResult.HasValue()) {
        ReportFailure(StartResult.ErrorValue());
        return 2;
    }
    Application->Run();
    return 0;
}
