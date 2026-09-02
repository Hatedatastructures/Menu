#include <QtTest/QtTest>

#include <DisplayMetrics.hpp>

class DisplayMetricsTest final : public QObject {
    Q_OBJECT

private slots:
    void UsesFallbackForMissingRefreshRate();
    void PreservesReportedRefreshRate();
    void CalculatesFrameBudget();
    void InitializesWithFallbackBudget();
};

void DisplayMetricsTest::UsesFallbackForMissingRefreshRate() {
    QCOMPARE(Menu::Client::EffectiveRefreshRate(0.0), 120.0);
    QCOMPARE(Menu::Client::EffectiveRefreshRate(-1.0), 120.0);
}

void DisplayMetricsTest::PreservesReportedRefreshRate() {
    QCOMPARE(Menu::Client::EffectiveRefreshRate(59.94), 59.94);
    QCOMPARE(Menu::Client::EffectiveRefreshRate(144.0), 144.0);
}

void DisplayMetricsTest::CalculatesFrameBudget() {
    QCOMPARE(Menu::Client::FrameBudgetMilliseconds(120.0), 1000.0 / 120.0);
    QCOMPARE(Menu::Client::FrameBudgetMilliseconds(0.0), 1000.0 / 120.0);
}

void DisplayMetricsTest::InitializesWithFallbackBudget() {
    Menu::Client::DisplayMetrics Metrics;
    QVERIFY(Metrics.RefreshRateHz() >= 1.0);
    QVERIFY(Metrics.ReportedRefreshRateHz() >= 0.0);
}

QTEST_GUILESS_MAIN(DisplayMetricsTest)
#include "DisplayMetricsTest.moc"
