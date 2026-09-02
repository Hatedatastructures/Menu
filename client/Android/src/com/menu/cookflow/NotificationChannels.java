package com.menu.cookflow;

import android.app.NotificationChannel;
import android.app.NotificationManager;
import android.content.Context;
import android.os.Build;

final class NotificationChannels {
    static final String TIMER = "cooking_timers";
    static final String REMINDER = "meal_reminders";

    private NotificationChannels() {
    }

    static void ensure(Context context) {
        if (Build.VERSION.SDK_INT < 26 || context == null) {
            return;
        }
        NotificationManager manager = context.getSystemService(NotificationManager.class);
        if (manager == null) {
            return;
        }
        NotificationChannel timers = new NotificationChannel(
                TIMER, "做饭计时", NotificationManager.IMPORTANCE_HIGH);
        timers.setDescription("步骤计时完成");
        timers.enableVibration(true);
        timers.setVibrationPattern(new long[]{0, 180, 100, 180});
        NotificationChannel reminders = new NotificationChannel(
                REMINDER, "备菜提醒", NotificationManager.IMPORTANCE_DEFAULT);
        reminders.setDescription("计划和备菜时间提醒");
        reminders.enableVibration(true);
        reminders.setVibrationPattern(new long[]{0, 120});
        manager.createNotificationChannel(timers);
        manager.createNotificationChannel(reminders);
    }
}
