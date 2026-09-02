package com.menu.cookflow;

import android.os.Bundle;

import org.qtproject.qt.android.bindings.QtActivity;

public final class MenuActivity extends QtActivity {
    @Override
    public void onCreate(Bundle state) {
        super.onCreate(state);
        SystemUiBridge.applyForActivity(this);
    }

    @Override
    protected void onResume() {
        super.onResume();
        SystemUiBridge.applyForActivity(this);
    }

    @Override
    public void onWindowFocusChanged(boolean hasFocus) {
        super.onWindowFocusChanged(hasFocus);
        if (hasFocus) {
            SystemUiBridge.applyForActivity(this);
        }
    }
}
