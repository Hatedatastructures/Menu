#include <QtTest/QtTest>

#include <ScreenAwakeController.hpp>

#include <memory>

namespace {

class FakeScreenAwakeBackend final : public Menu::Client::ScreenAwakeBackend {
public:
    int ApplyCount = 0;
    bool LastValue = false;

    void Apply(bool Enabled) override {
        ++ApplyCount;
        LastValue = Enabled;
    }
};

}  // namespace

class ScreenAwakeControllerTest final : public QObject {
    Q_OBJECT

private slots:
    void AppliesOnlyChangedState();
    void ClearsAwakeStateOnDestruction();
};

void ScreenAwakeControllerTest::AppliesOnlyChangedState() {
    auto Backend = std::make_unique<FakeScreenAwakeBackend>();
    auto* BackendPointer = Backend.get();
    auto Holder = std::make_shared<std::unique_ptr<FakeScreenAwakeBackend>>(
        std::move(Backend));
    Menu::Client::ScreenAwakeController Controller(nullptr,
        [Holder]() mutable {
            return std::move(*Holder);
        });

    Controller.SetEnabled(true);
    Controller.SetEnabled(true);
    Controller.SetEnabled(false);

    QCOMPARE(BackendPointer->ApplyCount, 2);
    QVERIFY(!BackendPointer->LastValue);
}

void ScreenAwakeControllerTest::ClearsAwakeStateOnDestruction() {
    auto Backend = std::make_unique<FakeScreenAwakeBackend>();
    auto* BackendPointer = Backend.get();
    auto Holder = std::make_shared<std::unique_ptr<FakeScreenAwakeBackend>>(
        std::move(Backend));
    {
        Menu::Client::ScreenAwakeController Controller(nullptr,
            [Holder]() mutable {
                return std::move(*Holder);
            });
        Controller.SetEnabled(true);
    }
    QCOMPARE(BackendPointer->ApplyCount, 2);
    QVERIFY(!BackendPointer->LastValue);
}

QTEST_GUILESS_MAIN(ScreenAwakeControllerTest)
#include "ScreenAwakeControllerTest.moc"
