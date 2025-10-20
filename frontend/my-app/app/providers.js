"use client";
import { DarkModeProvider } from "../context/darkMod";
import { ProfileProvider } from "../context/profile";
import { WSProvider } from "../context/wsContext";
import { ChatProvider } from "../context/chatContext";

export default function Providers({ children }) {
  return (
    <DarkModeProvider>
      <WSProvider>
        <ProfileProvider>
          <ChatProvider>{children}</ChatProvider>
        </ProfileProvider>
      </WSProvider>
    </DarkModeProvider>
  );
}
