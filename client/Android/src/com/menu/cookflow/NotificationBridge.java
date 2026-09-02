package com.menu.cookflow;

import android.content.BroadcastReceiver;
import android.content.Context;
import android.content.Intent;

public final class NotificationBridge {
    private NotificationBridge() {
    }

    public static boolean hasNotificationPermission(Context context) {
        return NotificationPermissions.isGranted(context);
    }

    public static void requestNotificationPermission(Context context) {
        NotificationPermissions.request(context);
    }

    public static void scheduleTimer(
            Context context,
            String id,
            long delaySeconds,
            String title,
            String body) {
        NotificationScheduler.schedule(context, id, delaySeconds, title, body);
    }

    public static void cancelTimer(Context context, String id) {
        NotificationScheduler.cancel(context, id);
    }

    public static void cancelAllTimers(Context context) {
        NotificationScheduler.cancelAll(context);
    }

    public static void showTestNotification(Context context) {
        NotificationRenderer.showTest(context);
    }

    public static void openNotificationSettings(Context context) {
        NotificationPermissions.openSettings(context);
    }

    public static final class ReminderReceiver extends BroadcastReceiver {
        @Override
        public void onReceive(Context context, Intent intent) {
            NotificationScheduler.receive(context, intent);
        }
    }
}
