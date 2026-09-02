#include <QtTest/QtTest>

#include <QSettings>

#include <ThemeController.hpp>

#include <memory>

namespace {

class FakeSystemUiBackend final : public Menu::Client::SystemUiBackend {
public:
    int Calls = 0;
    bool LastDark = false;

    void ApplyTheme(bool Dark) override {
        ++Calls;
        LastDark = Dark;
    }
};

}  // namespace

class ThemeControllerTest final : public QObject {
    Q_OBJECT

private slots:
    void initTestCase();
    void SupportsSystemLightAndDarkModes();
    void PersistsExplicitMode();
    void UsesInjectedSystemUiBackend();
};

void ThemeControllerTest::initTestCase() {
    QCoreApplication::setOrganizationName(QStringLiteral("MenuTests"));
    QCoreApplication::setApplicationName(QStringLiteral("ThemeController"));
    QSettings Settings;
    Settings.clear();
}

void ThemeControllerTest::SupportsSystemLightAndDarkModes() {
    Menu::Client::ThemeController Controller;
    Controller.SetMode(QStringLiteral("dark"));
    QCOMPARE(Controller.Mode(), QStringLiteral("dark"));
    QVERIFY(Controller.Dark());

    Controller.SetMode(QStringLiteral("light"));
    QCOMPARE(Controller.Mode(), QStringLiteral("light"));
    QVERIFY(!Controller.Dark());

    Controller.SetMode(QStringLiteral("system"));
    QCOMPARE(Controller.Mode(), QStringLiteral("system"));
}

void ThemeControllerTest::PersistsExplicitMode() {
    {
        Menu::Client::ThemeController Controller;
        Controller.SetMode(QStringLiteral("dark"));
    }
    Menu::Client::ThemeController Restored;
    QCOMPARE(Restored.Mode(), QStringLiteral("dark"));
    Restored.SetMode(QStringLiteral("system"));
}

void ThemeControllerTest::UsesInjectedSystemUiBackend() {
    auto Backend = std::make_unique<FakeSystemUiBackend>();
    FakeSystemUiBackend* BackendPointer = Backend.get();
    auto Holder = std::make_shared<std::unique_ptr<FakeSystemUiBackend>>(std::move(Backend));
    Menu::Client::ThemeController Controller(nullptr,
        [Holder]() mutable {
            return std::unique_ptr<Menu::Client::SystemUiBackend>(std::move(*Holder));
        });
    const int InitialCalls = BackendPointer->Calls;
    Controller.SetMode(QStringLiteral("dark"));
    QVERIFY(BackendPointer->Calls > InitialCalls);
    QVERIFY(BackendPointer->LastDark);
}

QTEST_GUILESS_MAIN(ThemeControllerTest)
#include "ThemeControllerTest.moc"
