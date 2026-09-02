package com.menu.cookflow;

import android.app.AlarmManager;
import android.app.PendingIntent;
import android.content.Context;
import android.content.Intent;
import android.content.SharedPreferences;
import android.os.Build;
import android.os.SystemClock;

import java.util.HashSet;
import java.util.Set;

final class NotificationScheduler {
    static final String ACTION_TIMER = "com.menu.cookflow.ACTION_TIMER";
    static final String ACTION_CANCEL_TIMER = "com.menu.cookflow.ACTION_CANCEL_TIMER";
    private static final String EXTRA_ID = "timer_id";
    private static final String EXTRA_TITLE = "timer_title";
    private static final String EXTRA_BODY = "timer_body";
    private static final String TIMER_STORE = "scheduled_timers";
    private static final String TIMER_IDS = "ids";

    private NotificationScheduler() {
    }

    static void schedule(
            Context context,
            String id,
            long delaySeconds,
            String title,
            String body) {
        if (context == null || id == null || id.isEmpty() || delaySeconds <= 0) {
            return;
        }
        NotificationChannels.ensure(context);
        Intent intent = new Intent(context, NotificationBridge.ReminderReceiver.class)
                .setAction(ACTION_TIMER)
                .putExtra(EXTRA_ID, id)
                .putExtra(EXTRA_TITLE, title)
                .putExtra(EXTRA_BODY, body);
        PendingIntent pendingIntent = PendingIntent.getBroadcast(
                context, requestCode(id), intent,
                PendingIntent.FLAG_UPDATE_CURRENT | PendingIntent.FLAG_IMMUTABLE);
        AlarmManager alarms = (AlarmManager) context.getSystemService(Context.ALARM_SERVICE);
        if (alarms == null) {
            return;
        }
        long triggerAt = SystemClock.elapsedRealtime() + delaySeconds * 1000L;
        if (Build.VERSION.SDK_INT >= 31 && alarms.canScheduleExactAlarms()) {
            alarms.setExactAndAllowWhileIdle(
                    AlarmManager.ELAPSED_REALTIME_WAKEUP, triggerAt, pendingIntent);
        } else if (Build.VERSION.SDK_INT >= 23) {
            alarms.setAndAllowWhileIdle(
                    AlarmManager.ELAPSED_REALTIME_WAKEUP, triggerAt, pendingIntent);
        } else {
            alarms.set(AlarmManager.ELAPSED_REALTIME_WAKEUP, triggerAt, pendingIntent);
        }
        remember(context, id, title, body,
                System.currentTimeMillis() + delaySeconds * 1000L);
    }

    static void cancel(Context context, String id) {
        if (context == null || id == null || id.isEmpty()) {
            return;
        }
        Intent intent = new Intent(context, NotificationBridge.ReminderReceiver.class)
                .setAction(ACTION_TIMER);
        PendingIntent pendingIntent = PendingIntent.getBroadcast(
                context, requestCode(id), intent,
                PendingIntent.FLAG_UPDATE_CURRENT | PendingIntent.FLAG_IMMUTABLE);
        AlarmManager alarms = (AlarmManager) context.getSystemService(Context.ALARM_SERVICE);
        if (alarms != null) {
            alarms.cancel(pendingIntent);
        }
        pendingIntent.cancel();
        forget(context, id);
    }

    static void cancelAll(Context context) {
        if (context == null) {
            return;
        }
        SharedPreferences preferences = context.getSharedPreferences(TIMER_STORE, Context.MODE_PRIVATE);
        Set<String> ids = new HashSet<>(preferences.getStringSet(TIMER_IDS, new HashSet<>()));
        for (String id : ids) {
            cancel(context, id);
        }
    }

    static void receive(Context context, Intent intent) {
        if (context == null || intent == null) {
            return;
        }
        if (Intent.ACTION_BOOT_COMPLETED.equals(intent.getAction())) {
            restore(context);
            return;
        }
        if (ACTION_CANCEL_TIMER.equals(intent.getAction())) {
            String id = intent.getStringExtra(EXTRA_ID);
            cancel(context, id);
            NotificationRenderer.cancel(context, id);
            return;
        }
        if (!ACTION_TIMER.equals(intent.getAction())) {
            return;
        }
        String id = intent.getStringExtra(EXTRA_ID);
        NotificationRenderer.show(
                context,
                id == null ? "timer" : id,
                valueOrDefault(intent.getStringExtra(EXTRA_TITLE), "计时完成"),
                valueOrDefault(intent.getStringExtra(EXTRA_BODY), "可以继续下一步了"));
        if (id != null) {
            forget(context, id);
        }
    }

    private static String valueOrDefault(String value, String fallback) {
        return value == null || value.isEmpty() ? fallback : value;
    }

    private static int requestCode(String id) {
        return id.hashCode() & 0x7fffffff;
    }

    private static void remember(
            Context context, String id, String title, String body, long deadlineMillis) {
        SharedPreferences preferences = context.getSharedPreferences(TIMER_STORE, Context.MODE_PRIVATE);
        Set<String> ids = new HashSet<>(preferences.getStringSet(TIMER_IDS, new HashSet<>()));
        ids.add(id);
        preferences.edit()
                .putStringSet(TIMER_IDS, ids)
                .putString(id + ".title", title == null ? "步骤完成" : title)
                .putString(id + ".body", body == null ? "可以继续下一步了" : body)
                .putLong(id + ".deadline", deadlineMillis)
                .apply();
    }

    private static void forget(Context context, String id) {
        if (id == null) {
            return;
        }
        SharedPreferences preferences = context.getSharedPreferences(TIMER_STORE, Context.MODE_PRIVATE);
        Set<String> ids = new HashSet<>(preferences.getStringSet(TIMER_IDS, new HashSet<>()));
        ids.remove(id);
        preferences.edit()
                .putStringSet(TIMER_IDS, ids)
                .remove(id + ".title")
                .remove(id + ".body")
                .remove(id + ".deadline")
                .apply();
    }

    private static void restore(Context context) {
        SharedPreferences preferences = context.getSharedPreferences(TIMER_STORE, Context.MODE_PRIVATE);
        Set<String> ids = new HashSet<>(preferences.getStringSet(TIMER_IDS, new HashSet<>()));
        long now = System.currentTimeMillis();
        for (String id : ids) {
            long deadline = preferences.getLong(id + ".deadline", 0L);
            if (deadline <= now) {
                forget(context, id);
                continue;
            }
            long delaySeconds = Math.max(1L, (deadline - now + 999L) / 1000L);
            schedule(
                    context,
                    id,
                    delaySeconds,
                    preferences.getString(id + ".title", "步骤完成"),
                    preferences.getString(id + ".body", "可以继续下一步了"));
        }
    }
}
