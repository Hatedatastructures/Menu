#include <QColor>
#include <QDir>
#include <QImage>
#include <QPainter>
#include <QPainterPath>
#include <QString>

#include <array>

namespace {

void DrawFoodArtwork(const QString& Path, int Variant) {
    const std::array<QColor, 13> Backgrounds = {
        QColor("#e9e1d3"), QColor("#dfe9df"), QColor("#eee0d8"),
        QColor("#e4e4ed"), QColor("#e8e1cf"), QColor("#dce8e5"),
        QColor("#eee4d6"), QColor("#e2e8dc"), QColor("#eadfdb"),
        QColor("#e1e6ec"), QColor("#e8e5d7"), QColor("#dde8e0"),
        QColor("#ebe0d0")};
    const std::array<QColor, 13> Accents = {
        QColor("#d65e43"), QColor("#b78338"), QColor("#a34935"),
        QColor("#8b6b47"), QColor("#d15d3b"), QColor("#d89a38"),
        QColor("#d46b4e"), QColor("#9d7a34"), QColor("#d55b42"),
        QColor("#b76d36"), QColor("#d16b49"), QColor("#b94c3b"),
        QColor("#ce8439")};
    QImage Image(960, 640, QImage::Format_RGB32);
    Image.fill(Backgrounds.at(Variant));
    QPainter Painter(&Image);
    Painter.setRenderHint(QPainter::Antialiasing, true);
    Painter.setPen(Qt::NoPen);

    Painter.setBrush(QColor(0, 0, 0, 24));
    Painter.drawEllipse(QRectF(112.0, 145.0, 736.0, 380.0));
    Painter.setBrush(QColor("#fffdf8"));
    Painter.drawEllipse(QRectF(96.0, 120.0, 736.0, 380.0));
    Painter.setBrush(QColor("#f1eee5"));
    Painter.drawEllipse(QRectF(124.0, 148.0, 680.0, 324.0));

    const QColor Accent = Accents.at(Variant);
    const QColor AccentDark = Accent.darker(125);
    for (int Index = 0; Index < 7; ++Index) {
        const qreal X = 230.0 + static_cast<qreal>((Index * 83 + Variant * 19) % 390);
        const qreal Y = 216.0 + static_cast<qreal>((Index * 47 + Variant * 13) % 160);
        const qreal Width = 78.0 + static_cast<qreal>((Index + Variant) % 3) * 13.0;
        const qreal Height = 56.0 + static_cast<qreal>((Index * 2 + Variant) % 3) * 11.0;
        Painter.setBrush(Index % 3 == 0 ? Accent : AccentDark);
        Painter.drawEllipse(QRectF(X, Y, Width, Height));
    }

    Painter.setBrush(QColor("#f6c84b"));
    for (int Index = 0; Index < 3; ++Index) {
        const qreal X = 275.0 + static_cast<qreal>((Index * 123 + Variant * 7) % 270);
        const qreal Y = 250.0 + static_cast<qreal>((Index * 31 + Variant * 9) % 100);
        Painter.drawEllipse(QRectF(X, Y, 68.0, 56.0));
    }

    Painter.setBrush(QColor("#4e8b62"));
    for (int Index = 0; Index < 8; ++Index) {
        const qreal X = 205.0 + static_cast<qreal>((Index * 71 + Variant * 23) % 470);
        const qreal Y = 205.0 + static_cast<qreal>((Index * 59 + Variant * 17) % 190);
        Painter.drawEllipse(QRectF(X, Y, 28.0, 18.0));
    }

    Painter.setBrush(QColor("#4b3d35"));
    Painter.drawRoundedRect(QRectF(680.0, 183.0, 22.0, 128.0), 11.0, 11.0);
    Painter.drawRoundedRect(QRectF(714.0, 183.0, 22.0, 128.0), 11.0, 11.0);

    QPainterPath Steam;
    Steam.moveTo(420.0, 165.0);
    Steam.cubicTo(390.0, 130.0, 450.0, 115.0, 420.0, 80.0);
    Steam.moveTo(495.0, 160.0);
    Steam.cubicTo(465.0, 125.0, 525.0, 110.0, 495.0, 74.0);
    Painter.setPen(QPen(QColor(255, 255, 255, 150), 7.0, Qt::SolidLine, Qt::RoundCap));
    Painter.drawPath(Steam);
    Painter.end();

    Image.save(Path, "PNG");
}

}  // namespace

int main(int ArgumentCount, char* Arguments[]) {
    if (ArgumentCount < 2) {
        return 2;
    }
    const QString OutputDirectory = QString::fromLocal8Bit(Arguments[1]);
    if (!QDir().mkpath(OutputDirectory)) {
        return 3;
    }
    const QStringList Names = {
        QStringLiteral("tomato-egg.png"), QStringLiteral("soy-chicken.png"),
        QStringLiteral("beef-onion.png"), QStringLiteral("tofu-spinach.png"),
        QStringLiteral("pasta-tomato.png"), QStringLiteral("lemon-salmon.png"),
        QStringLiteral("cream-spinach.png"), QStringLiteral("oyakodon.png"),
        QStringLiteral("tofu-teriyaki.png"), QStringLiteral("ginger-beef.png"),
        QStringLiteral("tomato-rice.png"), QStringLiteral("tofu-salad.png"),
        QStringLiteral("menu-placeholder.png")};
    for (int Index = 0; Index < Names.size(); ++Index) {
        DrawFoodArtwork(QDir(OutputDirectory).filePath(Names.at(Index)), Index);
    }
    return 0;
}
