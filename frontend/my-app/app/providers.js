"use client";
import { DarkModeProvider } from "../context/darkMod";
import { ProfileProvider } from "../context/profile";
import { WSProvider } from "../context/wsContext";

export default function Providers({ children }) {
  return (
    <DarkModeProvider>
      <WSProvider>
        <ProfileProvider>{children}</ProfileProvider>
      </WSProvider>

    </DarkModeProvider>
  );
}
