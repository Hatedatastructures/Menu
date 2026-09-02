package com.menu.cookflow;

import android.app.Activity;
import android.content.Context;

import java.lang.reflect.Method;

final class AndroidActivityResolver {
    private AndroidActivityResolver() {
    }

    static Activity resolve(Context context) {
        if (context instanceof Activity) {
            return (Activity) context;
        }
        try {
            Class<?> qtNative = Class.forName("org.qtproject.qt.android.QtNative");
            Method method = qtNative.getDeclaredMethod("activity");
            method.setAccessible(true);
            Object value = method.invoke(null);
            return value instanceof Activity ? (Activity) value : null;
        } catch (ReflectiveOperationException | SecurityException ignored) {
            return null;
        }
    }
}
