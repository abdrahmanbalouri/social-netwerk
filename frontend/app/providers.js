"use client";
import { DarkModeProvider } from "../context/darkMod";
import { ProfileProvider } from "../context/profile";
import { WSProvider } from "../context/wsContext";
import GlobalNotification from "../context/GlobalNotification";

export default function Providers({ children }) {
  return (
    <DarkModeProvider>
      <WSProvider>
        <ProfileProvider>
            {children}
          <GlobalNotification />
        </ProfileProvider>
      </WSProvider>
    </DarkModeProvider>
  );
}
