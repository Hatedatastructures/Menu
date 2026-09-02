package com.menu.cookflow;

import android.content.Context;
import android.graphics.Color;
import android.view.View;
import android.view.Window;
import android.view.WindowInsetsController;

import androidx.core.view.WindowCompat;
import androidx.core.view.WindowInsetsControllerCompat;


public final class SystemUiBridge {
    private SystemUiBridge() {
    }

    private static volatile boolean lastDark;

    public static int statusBarHeight(Context context) {
        if (context == null) {
            return 0;
        }
        int resourceId = context.getResources().getIdentifier(
                "status_bar_height", "dimen", "android");
        return resourceId > 0
                ? context.getResources().getDimensionPixelSize(resourceId) : 0;
    }

    public static void setKeepScreenOn(Context context, boolean enabled) {
        android.app.Activity activity = AndroidActivityResolver.resolve(context);
        if (activity == null) {
            return;
        }
        activity.runOnUiThread(() -> {
            if (enabled) {
                activity.getWindow().addFlags(
                        android.view.WindowManager.LayoutParams.FLAG_KEEP_SCREEN_ON);
            } else {
                activity.getWindow().clearFlags(
                        android.view.WindowManager.LayoutParams.FLAG_KEEP_SCREEN_ON);
            }
        });
    }

    public static void applySystemBars(Context context, boolean dark) {
        lastDark = dark;
        android.app.Activity activity = AndroidActivityResolver.resolve(context);
        if (activity == null) {
            return;
        }
        applyForActivity(activity);
    }

    static void applyForActivity(android.app.Activity activity) {
        if (activity == null) {
            return;
        }
        activity.runOnUiThread(() -> {
            Window window = activity.getWindow();
            WindowCompat.setDecorFitsSystemWindows(window, false);
            applyAppearance(window, lastDark);
            View decor = window.getDecorView();
            decor.setOnSystemUiVisibilityChangeListener(ignored -> applyAppearance(window, lastDark));
            decor.postDelayed(() -> applyAppearance(window, lastDark), 250L);
            decor.postDelayed(() -> applyAppearance(window, lastDark), 1000L);
        });
    }

    private static void applyAppearance(Window window, boolean dark) {
        int statusBarColor = dark ? Color.parseColor("#0b100d")
                : Color.parseColor("#e7f0ea");
        window.setStatusBarColor(statusBarColor);
        window.setNavigationBarColor(Color.parseColor(dark ? "#0b0f0c" : "#f6f7f3"));
        View decor = window.getDecorView();
        int flags = decor.getSystemUiVisibility();
        if (android.os.Build.VERSION.SDK_INT >= 23) {
            flags &= ~View.SYSTEM_UI_FLAG_LIGHT_STATUS_BAR;
        }
        if (android.os.Build.VERSION.SDK_INT >= 26) {
            flags &= ~View.SYSTEM_UI_FLAG_LIGHT_NAVIGATION_BAR;
        }
        if (!dark) {
            if (android.os.Build.VERSION.SDK_INT >= 23) {
                flags |= View.SYSTEM_UI_FLAG_LIGHT_STATUS_BAR;
            }
            if (android.os.Build.VERSION.SDK_INT >= 26) {
                flags |= View.SYSTEM_UI_FLAG_LIGHT_NAVIGATION_BAR;
            }
        }
        decor.setSystemUiVisibility(flags);
        WindowInsetsControllerCompat compatController =
                WindowCompat.getInsetsController(window, decor);
        if (compatController != null) {
            compatController.setAppearanceLightStatusBars(!dark);
            compatController.setAppearanceLightNavigationBars(!dark);
        }
        if (android.os.Build.VERSION.SDK_INT >= 30) {
            WindowInsetsController controller = window.getInsetsController();
            if (controller != null) {
                int appearance = dark ? 0
                        : WindowInsetsController.APPEARANCE_LIGHT_STATUS_BARS
                                | WindowInsetsController.APPEARANCE_LIGHT_NAVIGATION_BARS;
                controller.setSystemBarsAppearance(
                        appearance,
                        WindowInsetsController.APPEARANCE_LIGHT_STATUS_BARS
                                | WindowInsetsController.APPEARANCE_LIGHT_NAVIGATION_BARS);
            }
        }
        if (android.os.Build.VERSION.SDK_INT >= 29) {
            window.setStatusBarContrastEnforced(false);
            window.setNavigationBarContrastEnforced(false);
        }
    }

}
