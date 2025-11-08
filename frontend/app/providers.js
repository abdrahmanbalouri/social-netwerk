"use client";
import { DarkModeProvider } from "../context/darkMod";
import { ProfileProvider } from "../context/profile";
import { WSProvider } from "../context/wsContext";
import { ChatProvider } from "../context/chatContext";
import GlobalNotification from "../context/GlobalNotification";

export default function Providers({ children }) {
  return (
    <DarkModeProvider>
      <WSProvider>
        <ProfileProvider>
          <ChatProvider>
            {children}
            <GlobalNotification />
          </ChatProvider>
        </ProfileProvider>
      </WSProvider>
    </DarkModeProvider>
  );
}
