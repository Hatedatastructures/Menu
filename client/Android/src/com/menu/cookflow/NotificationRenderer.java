package com.menu.cookflow;

import android.app.Notification;
import android.app.PendingIntent;
import android.content.Context;
import android.content.Intent;

import androidx.core.app.NotificationCompat;
import androidx.core.app.NotificationManagerCompat;
import androidx.core.content.ContextCompat;

final class NotificationRenderer {
    private static final String ACTION_CANCEL_TIMER = "com.menu.cookflow.ACTION_CANCEL_TIMER";
    private static final String EXTRA_ID = "timer_id";
    private static final String TIMER_GROUP = "cookflow_timers";

    private NotificationRenderer() {
    }

    static void showTest(Context context) {
        show(context, "test", "Menu 通知测试", "通知、震动和深浅色由系统设置控制");
    }

    static void show(Context context, String id, String title, String body) {
        if (context == null || !NotificationPermissions.isGranted(context)) {
            return;
        }
        Intent launchIntent = context.getPackageManager()
                .getLaunchIntentForPackage(context.getPackageName());
        PendingIntent contentIntent = launchIntent == null ? null : PendingIntent.getActivity(
                context, requestCode(id), launchIntent,
                PendingIntent.FLAG_UPDATE_CURRENT | PendingIntent.FLAG_IMMUTABLE);
        NotificationCompat.Builder builder = new NotificationCompat.Builder(
                context, NotificationChannels.TIMER)
                .setSmallIcon(R.drawable.ic_menu_notification)
                .setContentTitle(title)
                .setContentText(body)
                .setStyle(new NotificationCompat.BigTextStyle().bigText(body))
                .setCategory(NotificationCompat.CATEGORY_REMINDER)
                .setPriority(NotificationCompat.PRIORITY_HIGH)
                .setVisibility(NotificationCompat.VISIBILITY_PRIVATE)
                .setOnlyAlertOnce(true)
                .setGroup(TIMER_GROUP)
                .setAutoCancel(true)
                .setColor(ContextCompat.getColor(context, R.color.menu_notification_accent));
        if (contentIntent != null) {
            builder.setContentIntent(contentIntent);
        }
        Intent cancelIntent = new Intent(context, NotificationBridge.ReminderReceiver.class)
                .setAction(ACTION_CANCEL_TIMER)
                .putExtra(EXTRA_ID, id);
        PendingIntent cancelPendingIntent = PendingIntent.getBroadcast(
                context, requestCode(id + ".cancel"), cancelIntent,
                PendingIntent.FLAG_UPDATE_CURRENT | PendingIntent.FLAG_IMMUTABLE);
        builder.addAction(R.drawable.ic_menu_notification, "停止计时", cancelPendingIntent);
        Notification notification = builder.build();
        NotificationManagerCompat.from(context).notify(requestCode(id), notification);
    }

    static void cancel(Context context, String id) {
        if (context != null && id != null) {
            NotificationManagerCompat.from(context).cancel(requestCode(id));
        }
    }

    private static int requestCode(String id) {
        return id.hashCode() & 0x7fffffff;
    }
}
