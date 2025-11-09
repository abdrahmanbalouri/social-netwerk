"use client";
import { DarkModeProvider } from "../context/darkMod";
import { ProfileProvider } from "../context/profile";
import { WSProvider } from "../context/wsContext";
import { ToastProvider } from "../context/toastContext";
import GlobalNotification from "../context/GlobalNotification";

export default function Providers({ children }) {
  return (
    <DarkModeProvider>
      <WSProvider>
        <ProfileProvider>
          <ToastProvider>
            {children}
          </ToastProvider>
          <GlobalNotification />
        </ProfileProvider>
      </WSProvider>
    </DarkModeProvider>
  );
}
